// Package localtime parses the dates and times people type — "4:50pm",
// "9/9/26", "this sunday 11:59pm" — in a local time zone and returns the
// UTC instant Canvas expects. Every command that takes a wall-clock time
// from the user goes through Parse so that the accepted forms, the zone
// resolution and the ambiguity rules are the same everywhere.
//
// Zone resolution: Options.Timezone (an IANA name such as
// America/New_York), else the TZ environment variable, else the system
// zone. Commands put the config file's settings.timezone into
// Options.Timezone before calling Parse.
package localtime

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// AcceptedForms lists the inputs Parse understands. Error messages quote it
// so the user sees what to type instead.
const AcceptedForms = `  times:      4pm, 4:00pm, 4:00 PM, 16:50, noon, midnight
  dates:      2026-09-09, 9/9/26, 9/9/2026, today, tomorrow, yesterday,
              sunday, this sunday, next monday
  both:       2026-09-09 4:50pm, tomorrow 9am, this sunday 11:59pm
  exact:      2026-09-09T16:50 (local), 2026-09-09T20:50:00Z, 2026-09-09T16:50:00-04:00`

// Options control how an input is interpreted.
type Options struct {
	// Timezone is an IANA zone name. Empty means $TZ, then the system zone.
	Timezone string
	// Now anchors today/tomorrow/this sunday. Zero means time.Now().
	Now time.Time
	// DateContext supplies the calendar day for a time-only input such as
	// "4:50pm". Zero means Now's day in the resolved zone.
	DateContext time.Time
}

// Parsed is the result of Parse.
type Parsed struct {
	// Time is the instant in UTC — the value to send to Canvas.
	Time time.Time
	// Local is the same instant in the resolved zone.
	Local time.Time
	// Location is the zone the input was read in.
	Location *time.Location
	// HasDate is true when the input named a calendar day (as opposed to
	// taking it from DateContext).
	HasDate bool
	// HasTime is true when the input named a clock time. A date-only input
	// resolves to 00:00 local; callers that want end-of-day use EndOfDay.
	HasTime bool
	// Input is the original text.
	Input string
}

// Location resolves a zone name: name if given, else $TZ, else the system
// zone. An unknown name is an error that shows the expected form.
func Location(name string) (*time.Location, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		name = strings.TrimSpace(os.Getenv("TZ"))
	}
	if name == "" {
		return time.Local, nil
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, fmt.Errorf("unknown time zone %q (use an IANA name such as America/New_York or UTC)", name)
	}
	return loc, nil
}

// ErrEmpty is returned for blank input.
var ErrEmpty = errors.New("empty date/time")

// Parse reads input in the zone resolved from opts and returns the UTC
// instant. Time-only inputs combine with opts.DateContext (default: today).
// Ambiguous inputs — "4:50" without am/pm, a clock time that does not exist
// or occurs twice on a DST changeover — are refused with a message that
// shows the accepted forms.
func Parse(input string, opts Options) (*Parsed, error) {
	loc, err := Location(opts.Timezone)
	if err != nil {
		return nil, err
	}
	raw := strings.TrimSpace(input)
	if raw == "" {
		return nil, ErrEmpty
	}
	now := opts.Now
	if now.IsZero() {
		now = time.Now()
	}
	now = now.In(loc)

	// Exact forms first: RFC 3339 (with zone) is passed through untouched;
	// the same shape without a zone is read in loc.
	if p, ok := parseExact(raw, loc); ok {
		p.Input = input
		return p, nil
	}

	words := strings.Fields(strings.ToLower(raw))
	day, rest, hasDate, err := parseDate(words, now)
	if err != nil {
		return nil, refuse(input, err)
	}
	clock, rest, hasTime, err := parseClock(rest)
	if err != nil {
		return nil, refuse(input, err)
	}
	if !hasDate && hasTime && len(rest) > 0 {
		// "4pm tomorrow": the date may follow the time.
		day, rest, hasDate, err = parseDate(rest, now)
		if err != nil {
			return nil, refuse(input, err)
		}
	}
	if len(rest) > 0 {
		return nil, refuse(input, fmt.Errorf("could not read %q", strings.Join(rest, " ")))
	}
	if !hasDate && !hasTime {
		return nil, refuse(input, errors.New("no date or time found"))
	}
	if !hasDate {
		ctx := opts.DateContext
		if ctx.IsZero() {
			ctx = now
		}
		ctx = ctx.In(loc)
		day = civilDay{ctx.Year(), ctx.Month(), ctx.Day()}
	}

	local, err := resolveWallClock(day, clock, loc)
	if err != nil {
		return nil, refuse(input, err)
	}
	return &Parsed{
		Time:     local.UTC(),
		Local:    local,
		Location: loc,
		HasDate:  hasDate,
		HasTime:  hasTime,
		Input:    input,
	}, nil
}

// EndOfDay returns the last second (23:59:59) of p's local day. It is what
// a date-only "due 9/9/26" or "by this sunday" usually means.
func (p *Parsed) EndOfDay() *Parsed {
	y, m, d := p.Local.Date()
	local := time.Date(y, m, d, 23, 59, 59, 0, p.Location)
	return &Parsed{Time: local.UTC(), Local: local, Location: p.Location, HasDate: p.HasDate, HasTime: true, Input: p.Input}
}

// Describe renders the instant in both the local zone and UTC:
// "Wed 2026-09-09 4:50 PM EDT (2026-09-09T20:50:00Z)".
func (p *Parsed) Describe() string {
	return Describe(p.Time, p.Location)
}

// Describe renders t in loc and in UTC. A nil loc means UTC.
func Describe(t time.Time, loc *time.Location) string {
	return FormatLocal(t, loc) + " (" + FormatUTC(t) + ")"
}

// FormatLocal renders t in loc as "Wed 2026-09-09 4:50 PM EDT".
func FormatLocal(t time.Time, loc *time.Location) string {
	if loc == nil {
		loc = time.UTC
	}
	return t.In(loc).Format("Mon 2006-01-02 3:04 PM MST")
}

// FormatUTC renders t as RFC 3339 in UTC.
func FormatUTC(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// --- internals ---------------------------------------------------------

type civilDay struct {
	year  int
	month time.Month
	day   int
}

type wallClock struct {
	hour, minute int
}

func refuse(input string, reason error) error {
	return fmt.Errorf("cannot read date/time %q: %v\naccepted forms:\n%s", input, reason, AcceptedForms)
}

var exactLayouts = []struct {
	layout string
	zoned  bool
}{
	{time.RFC3339, true},
	{"2006-01-02T15:04:05.999999999Z07:00", true},
	{"2006-01-02T15:04Z07:00", true},
	{"2006-01-02T15:04:05", false},
	{"2006-01-02T15:04", false},
}

func parseExact(raw string, loc *time.Location) (*Parsed, bool) {
	for _, l := range exactLayouts {
		var t time.Time
		var err error
		if l.zoned {
			t, err = time.Parse(l.layout, raw)
		} else {
			t, err = time.ParseInLocation(l.layout, raw, loc)
		}
		if err != nil {
			continue
		}
		return &Parsed{Time: t.UTC(), Local: t.In(loc), Location: loc, HasDate: true, HasTime: true}, true
	}
	return nil, false
}

var weekdays = map[string]time.Weekday{
	"sunday": time.Sunday, "sun": time.Sunday,
	"monday": time.Monday, "mon": time.Monday,
	"tuesday": time.Tuesday, "tue": time.Tuesday, "tues": time.Tuesday,
	"wednesday": time.Wednesday, "wed": time.Wednesday,
	"thursday": time.Thursday, "thu": time.Thursday, "thur": time.Thursday, "thurs": time.Thursday,
	"friday": time.Friday, "fri": time.Friday,
	"saturday": time.Saturday, "sat": time.Saturday,
}

var (
	isoDateRe = regexp.MustCompile(`^(\d{4})-(\d{1,2})-(\d{1,2})$`)
	usDateRe  = regexp.MustCompile(`^(\d{1,2})/(\d{1,2})/(\d{2}|\d{4})$`)
)

// parseDate consumes a date phrase from the front of words. "this sunday"
// is the first Sunday on or after today (today itself when today is a
// Sunday); "next sunday" is the one a week later; a bare "sunday" means
// "this sunday".
func parseDate(words []string, now time.Time) (civilDay, []string, bool, error) {
	today := civilDay{now.Year(), now.Month(), now.Day()}
	if len(words) == 0 {
		return today, words, false, nil
	}
	w := words[0]
	switch w {
	case "today":
		return today, words[1:], true, nil
	case "tomorrow":
		return today.add(1, now.Location()), words[1:], true, nil
	case "yesterday":
		return today.add(-1, now.Location()), words[1:], true, nil
	case "this", "next":
		if len(words) < 2 {
			return today, words, false, fmt.Errorf("%q needs a weekday after it (this sunday, next monday)", w)
		}
		wd, ok := weekdays[words[1]]
		if !ok {
			return today, words, false, fmt.Errorf("%q is not a weekday", words[1])
		}
		d := nextWeekday(today, wd, now.Location())
		if w == "next" {
			d = d.add(7, now.Location())
		}
		return d, words[2:], true, nil
	}
	if wd, ok := weekdays[w]; ok {
		return nextWeekday(today, wd, now.Location()), words[1:], true, nil
	}
	if m := isoDateRe.FindStringSubmatch(w); m != nil {
		d, err := makeDay(atoi(m[1]), atoi(m[2]), atoi(m[3]))
		return d, words[1:], true, err
	}
	if m := usDateRe.FindStringSubmatch(w); m != nil {
		year := atoi(m[3])
		if len(m[3]) == 2 {
			year += 2000
		}
		month, day := atoi(m[1]), atoi(m[2])
		if month > 12 {
			return today, words, false, fmt.Errorf("%q: dates are month/day/year; use 2026-09-09 to be explicit", w)
		}
		d, err := makeDay(year, month, day)
		return d, words[1:], true, err
	}
	return today, words, false, nil
}

func makeDay(year, month, day int) (civilDay, error) {
	if month < 1 || month > 12 || day < 1 || day > 31 {
		return civilDay{}, fmt.Errorf("%04d-%02d-%02d is not a calendar date", year, month, day)
	}
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if t.Day() != day {
		return civilDay{}, fmt.Errorf("%04d-%02d-%02d is not a calendar date", year, month, day)
	}
	return civilDay{year, time.Month(month), day}, nil
}

func (d civilDay) add(days int, loc *time.Location) civilDay {
	t := time.Date(d.year, d.month, d.day+days, 12, 0, 0, 0, loc)
	return civilDay{t.Year(), t.Month(), t.Day()}
}

func nextWeekday(from civilDay, wd time.Weekday, loc *time.Location) civilDay {
	t := time.Date(from.year, from.month, from.day, 12, 0, 0, 0, loc)
	delta := (int(wd) - int(t.Weekday()) + 7) % 7
	return from.add(delta, loc)
}

var (
	clockRe    = regexp.MustCompile(`^(\d{1,2})(?::(\d{2}))?(am|pm|a\.m\.|p\.m\.)?$`)
	meridiemRe = regexp.MustCompile(`^(am|pm|a\.m\.|p\.m\.)$`)
)

// parseClock consumes a clock phrase from the front of words: "4pm",
// "4:00pm", "4:00 pm", "16:50", "noon", "midnight". A 1–12 o'clock value
// without am/pm is refused as ambiguous; 24-hour values (0 or 13–23) are
// unambiguous and accepted.
func parseClock(words []string) (wallClock, []string, bool, error) {
	if len(words) == 0 {
		return wallClock{}, words, false, nil
	}
	w := words[0]
	switch w {
	case "noon":
		return wallClock{12, 0}, words[1:], true, nil
	case "midnight":
		return wallClock{0, 0}, words[1:], true, nil
	}
	m := clockRe.FindStringSubmatch(w)
	if m == nil {
		return wallClock{}, words, false, nil
	}
	rest := words[1:]
	meridiem := m[3]
	if meridiem == "" && len(rest) > 0 && meridiemRe.MatchString(rest[0]) {
		meridiem = rest[0]
		rest = rest[1:]
	}
	hour := atoi(m[1])
	minute := 0
	if m[2] != "" {
		minute = atoi(m[2])
	}
	if minute > 59 {
		return wallClock{}, words, false, fmt.Errorf("%q: minutes must be 00–59", w)
	}
	switch {
	case meridiem != "":
		if hour < 1 || hour > 12 {
			return wallClock{}, words, false, fmt.Errorf("%q: with am/pm the hour must be 1–12", w)
		}
		pm := strings.HasPrefix(meridiem, "p")
		if hour == 12 {
			hour = 0
		}
		if pm {
			hour += 12
		}
	case m[2] == "":
		// A bare number is not a time: "16" could be a date fragment or a typo.
		return wallClock{}, words, false, fmt.Errorf("%q: write a time as 4pm, 4:00pm or 16:00", w)
	case hour >= 1 && hour <= 12:
		return wallClock{}, words, false, fmt.Errorf("%q is ambiguous: write %d:%02dam, %d:%02dpm or the 24-hour %02d:%02d",
			w, hour, minute, hour, minute, hour+12, minute)
	case hour > 23:
		return wallClock{}, words, false, fmt.Errorf("%q: the hour must be 0–23", w)
	}
	return wallClock{hour, minute}, rest, true, nil
}

// resolveWallClock turns a civil day + clock in loc into an instant,
// refusing wall clocks that do not exist (spring-forward gap) or occur
// twice (fall-back fold) so that a DST edge never silently shifts a due
// time by an hour.
func resolveWallClock(d civilDay, c wallClock, loc *time.Location) (time.Time, error) {
	t := time.Date(d.year, d.month, d.day, c.hour, c.minute, 0, 0, loc)
	if t.Hour() != c.hour || t.Minute() != c.minute {
		return time.Time{}, fmt.Errorf("%02d:%02d does not exist on %04d-%02d-%02d in %s (clocks spring forward); give the time as RFC 3339 with an offset",
			c.hour, c.minute, d.year, d.month, d.day, loc)
	}
	sameWall := func(u time.Time) bool {
		return u.Year() == d.year && u.Month() == d.month && u.Day() == d.day && u.Hour() == c.hour && u.Minute() == c.minute
	}
	if sameWall(t.Add(time.Hour)) || sameWall(t.Add(-time.Hour)) {
		return time.Time{}, fmt.Errorf("%02d:%02d occurs twice on %04d-%02d-%02d in %s (clocks fall back); give the time as RFC 3339 with an offset",
			c.hour, c.minute, d.year, d.month, d.day, loc)
	}
	return t, nil
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
