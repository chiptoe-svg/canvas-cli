package localtime

import (
	"strings"
	"testing"
	"time"
)

const ny = "America/New_York"

// fixtureNow is Thursday 2026-09-03 10:00 in New York (EDT, UTC-4).
var fixtureNow = time.Date(2026, 9, 3, 10, 0, 0, 0, mustLoc(ny))

func mustLoc(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

func parseNY(t *testing.T, input string) *Parsed {
	t.Helper()
	p, err := Parse(input, Options{Timezone: ny, Now: fixtureNow})
	if err != nil {
		t.Fatalf("Parse(%q): %v", input, err)
	}
	return p
}

func TestParse_AcceptedForms(t *testing.T) {
	cases := []struct {
		in      string
		wantUTC string
		date    bool
		clock   bool
	}{
		// times combine with today (2026-09-03, EDT = UTC-4)
		{"4:00pm", "2026-09-03T20:00:00Z", false, true},
		{"4pm", "2026-09-03T20:00:00Z", false, true},
		{"4:50PM", "2026-09-03T20:50:00Z", false, true},
		{"4:50 pm", "2026-09-03T20:50:00Z", false, true},
		{"4:50 p.m.", "2026-09-03T20:50:00Z", false, true},
		{"16:50", "2026-09-03T20:50:00Z", false, true},
		{"00:30", "2026-09-03T04:30:00Z", false, true},
		{"12am", "2026-09-03T04:00:00Z", false, true},
		{"12pm", "2026-09-03T16:00:00Z", false, true},
		{"12:30am", "2026-09-03T04:30:00Z", false, true},
		{"noon", "2026-09-03T16:00:00Z", false, true},
		{"midnight", "2026-09-03T04:00:00Z", false, true},
		// dates resolve to local midnight
		{"2026-09-09", "2026-09-09T04:00:00Z", true, false},
		{"9/9/26", "2026-09-09T04:00:00Z", true, false},
		{"9/9/2026", "2026-09-09T04:00:00Z", true, false},
		{"09/09/2026", "2026-09-09T04:00:00Z", true, false},
		{"today", "2026-09-03T04:00:00Z", true, false},
		{"tomorrow", "2026-09-04T04:00:00Z", true, false},
		{"yesterday", "2026-09-02T04:00:00Z", true, false},
		{"this sunday", "2026-09-06T04:00:00Z", true, false},
		{"sunday", "2026-09-06T04:00:00Z", true, false},
		{"Sun", "2026-09-06T04:00:00Z", true, false},
		{"next sunday", "2026-09-13T04:00:00Z", true, false},
		{"this thursday", "2026-09-03T04:00:00Z", true, false}, // today is Thursday
		{"next thursday", "2026-09-10T04:00:00Z", true, false},
		{"next monday", "2026-09-14T04:00:00Z", true, false},
		{"this monday", "2026-09-07T04:00:00Z", true, false},
		// both
		{"2026-09-09 4:50pm", "2026-09-09T20:50:00Z", true, true},
		{"9/9/26 4:50pm", "2026-09-09T20:50:00Z", true, true},
		{"tomorrow 9am", "2026-09-04T13:00:00Z", true, true},
		{"this sunday 11:59pm", "2026-09-07T03:59:00Z", true, true},
		{"4pm tomorrow", "2026-09-04T20:00:00Z", true, true},
		{"  Today   4:00 PM ", "2026-09-03T20:00:00Z", true, true},
		// exact
		{"2026-09-09T16:50", "2026-09-09T20:50:00Z", true, true},
		{"2026-09-09T16:50:00", "2026-09-09T20:50:00Z", true, true},
		{"2026-09-09T20:50:00Z", "2026-09-09T20:50:00Z", true, true},
		{"2026-09-09T16:50:00-04:00", "2026-09-09T20:50:00Z", true, true},
		{"2026-09-09T16:50:00+02:00", "2026-09-09T14:50:00Z", true, true},
		// a winter date is EST (UTC-5)
		{"2026-12-01 4pm", "2026-12-01T21:00:00Z", true, true},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			p := parseNY(t, c.in)
			if got := p.Time.Format(time.RFC3339); got != c.wantUTC {
				t.Errorf("Parse(%q) = %s, want %s", c.in, got, c.wantUTC)
			}
			if p.Time.Location() != time.UTC {
				t.Errorf("Parse(%q).Time is not UTC", c.in)
			}
			if !p.Local.Equal(p.Time) || p.Local.Location().String() != ny {
				t.Errorf("Parse(%q).Local = %v (%s), want the same instant in %s", c.in, p.Local, p.Local.Location(), ny)
			}
			if p.HasDate != c.date || p.HasTime != c.clock {
				t.Errorf("Parse(%q) HasDate/HasTime = %v/%v, want %v/%v", c.in, p.HasDate, p.HasTime, c.date, c.clock)
			}
			if p.Input != c.in {
				t.Errorf("Input = %q, want %q", p.Input, c.in)
			}
		})
	}
}

func TestParse_Refusals(t *testing.T) {
	cases := []struct {
		in   string
		want string // substring of the error
	}{
		{"", "empty"},
		{"   ", "empty"},
		{"4:50", "ambiguous"},
		{"12:30", "ambiguous"},
		{"16", "write a time as"},
		{"25:00", "hour must be 0–23"},
		{"4:60pm", "minutes must be 00–59"},
		{"13pm", "hour must be 1–12"},
		{"soon", `could not read "soon"`},
		{"2026-09-09 at 4pm", `could not read "at 4pm"`},
		{"4pm 5pm", `could not read "5pm"`},
		{"this", "needs a weekday"},
		{"next week", "not a weekday"},
		{"2026-02-30", "not a calendar date"},
		{"13/9/26", "month/day/year"},
		{"2026-13-01", "not a calendar date"},
		{"4pm 2026-09-09 tomorrow", `could not read "tomorrow"`},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			_, err := Parse(c.in, Options{Timezone: ny, Now: fixtureNow})
			if err == nil {
				t.Fatalf("Parse(%q) succeeded, want an error containing %q", c.in, c.want)
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("Parse(%q) error = %q, want it to contain %q", c.in, err, c.want)
			}
			if c.in != "" && strings.TrimSpace(c.in) != "" && !strings.Contains(err.Error(), "accepted forms:") {
				t.Errorf("Parse(%q) error does not show the accepted forms:\n%s", c.in, err)
			}
		})
	}
}

// DST in New York, 2026: clocks spring forward 2026-03-08 02:00→03:00 and
// fall back 2026-11-01 02:00→01:00.
func TestParse_DSTEdges(t *testing.T) {
	t.Run("gap is refused", func(t *testing.T) {
		_, err := Parse("2026-03-08 2:30am", Options{Timezone: ny, Now: fixtureNow})
		if err == nil || !strings.Contains(err.Error(), "does not exist") {
			t.Fatalf("want a spring-forward refusal, got %v", err)
		}
	})
	t.Run("fold is refused", func(t *testing.T) {
		_, err := Parse("2026-11-01 1:30am", Options{Timezone: ny, Now: fixtureNow})
		if err == nil || !strings.Contains(err.Error(), "occurs twice") {
			t.Fatalf("want a fall-back refusal, got %v", err)
		}
	})
	t.Run("either side of the gap is fine", func(t *testing.T) {
		if got := parseNY(t, "2026-03-08 1:59am").Time.Format(time.RFC3339); got != "2026-03-08T06:59:00Z" {
			t.Errorf("1:59am EST = %s", got)
		}
		if got := parseNY(t, "2026-03-08 3:00am").Time.Format(time.RFC3339); got != "2026-03-08T07:00:00Z" {
			t.Errorf("3:00am EDT = %s", got)
		}
	})
	t.Run("either side of the fold is fine", func(t *testing.T) {
		if got := parseNY(t, "2026-11-01 12:59am").Time.Format(time.RFC3339); got != "2026-11-01T04:59:00Z" {
			t.Errorf("12:59am EDT = %s", got)
		}
		if got := parseNY(t, "2026-11-01 2:00am").Time.Format(time.RFC3339); got != "2026-11-01T07:00:00Z" {
			t.Errorf("2:00am EST = %s", got)
		}
	})
	t.Run("an offset resolves the fold", func(t *testing.T) {
		if got := parseNY(t, "2026-11-01T01:30:00-05:00").Time.Format(time.RFC3339); got != "2026-11-01T06:30:00Z" {
			t.Errorf("1:30am EST = %s", got)
		}
	})
	t.Run("day arithmetic across the changeover keeps the calendar day", func(t *testing.T) {
		sat := time.Date(2026, 3, 7, 23, 30, 0, 0, mustLoc(ny))
		p, err := Parse("tomorrow 4pm", Options{Timezone: ny, Now: sat})
		if err != nil {
			t.Fatal(err)
		}
		if got := p.Local.Format("2006-01-02 15:04 MST"); got != "2026-03-08 16:00 EDT" {
			t.Errorf("tomorrow 4pm from Saturday night = %s", got)
		}
	})
}

func TestParse_DateContext(t *testing.T) {
	ctx := time.Date(2026, 9, 9, 23, 30, 0, 0, time.UTC) // 7:30 PM on 2026-09-09 in New York
	p, err := Parse("4:50pm", Options{Timezone: ny, Now: fixtureNow, DateContext: ctx})
	if err != nil {
		t.Fatal(err)
	}
	if got := p.Time.Format(time.RFC3339); got != "2026-09-09T20:50:00Z" {
		t.Errorf("time-only with DateContext = %s, want 2026-09-09T20:50:00Z", got)
	}
	// The context's day is read in the target zone, not UTC.
	late := time.Date(2026, 9, 10, 2, 0, 0, 0, time.UTC) // still 2026-09-09 10 PM in New York
	p, err = Parse("4:50pm", Options{Timezone: ny, Now: fixtureNow, DateContext: late})
	if err != nil {
		t.Fatal(err)
	}
	if got := p.Time.Format(time.RFC3339); got != "2026-09-09T20:50:00Z" {
		t.Errorf("DateContext day taken in UTC: got %s", got)
	}
	// An explicit date wins over the context.
	p, err = Parse("2026-09-10 4:50pm", Options{Timezone: ny, Now: fixtureNow, DateContext: ctx})
	if err != nil {
		t.Fatal(err)
	}
	if got := p.Time.Format(time.RFC3339); got != "2026-09-10T20:50:00Z" {
		t.Errorf("explicit date with DateContext = %s", got)
	}
}

func TestParse_ZoneResolution(t *testing.T) {
	t.Run("explicit zone", func(t *testing.T) {
		p, err := Parse("2026-09-09 4:50pm", Options{Timezone: "Europe/Berlin", Now: fixtureNow})
		if err != nil {
			t.Fatal(err)
		}
		if got := p.Time.Format(time.RFC3339); got != "2026-09-09T14:50:00Z" {
			t.Errorf("Berlin 4:50pm = %s", got)
		}
	})
	t.Run("TZ environment", func(t *testing.T) {
		t.Setenv("TZ", "Asia/Tokyo")
		p, err := Parse("2026-09-09 4:50pm", Options{Now: fixtureNow})
		if err != nil {
			t.Fatal(err)
		}
		if got := p.Time.Format(time.RFC3339); got != "2026-09-09T07:50:00Z" {
			t.Errorf("Tokyo 4:50pm = %s", got)
		}
		if p.Location.String() != "Asia/Tokyo" {
			t.Errorf("Location = %s", p.Location)
		}
	})
	t.Run("explicit zone beats TZ", func(t *testing.T) {
		t.Setenv("TZ", "Asia/Tokyo")
		p, err := Parse("2026-09-09 4:50pm", Options{Timezone: "UTC", Now: fixtureNow})
		if err != nil {
			t.Fatal(err)
		}
		if got := p.Time.Format(time.RFC3339); got != "2026-09-09T16:50:00Z" {
			t.Errorf("UTC 4:50pm = %s", got)
		}
	})
	t.Run("system zone when nothing is set", func(t *testing.T) {
		t.Setenv("TZ", "")
		loc, err := Location("")
		if err != nil {
			t.Fatal(err)
		}
		if loc != time.Local {
			t.Errorf("Location(\"\") = %s, want time.Local", loc)
		}
	})
	t.Run("unknown zone", func(t *testing.T) {
		_, err := Parse("4pm", Options{Timezone: "Mars/Olympus"})
		if err == nil || !strings.Contains(err.Error(), "unknown time zone") {
			t.Errorf("want unknown-zone error, got %v", err)
		}
	})
	t.Run("relative days follow the zone, not UTC", func(t *testing.T) {
		// 2026-09-04 02:00 UTC is still Thursday 2026-09-03 10 PM in New York.
		now := time.Date(2026, 9, 4, 2, 0, 0, 0, time.UTC)
		p, err := Parse("today", Options{Timezone: ny, Now: now})
		if err != nil {
			t.Fatal(err)
		}
		if got := p.Local.Format("2006-01-02"); got != "2026-09-03" {
			t.Errorf("today = %s, want 2026-09-03", got)
		}
	})
}

func TestEndOfDayAndDescribe(t *testing.T) {
	p := parseNY(t, "9/9/26")
	e := p.EndOfDay()
	if got := e.Time.Format(time.RFC3339); got != "2026-09-10T03:59:59Z" {
		t.Errorf("EndOfDay = %s", got)
	}
	if !e.HasTime || !e.HasDate {
		t.Errorf("EndOfDay HasDate/HasTime = %v/%v", e.HasDate, e.HasTime)
	}
	if got := parseNY(t, "2026-09-09 4:50pm").Describe(); got != "Wed 2026-09-09 4:50 PM EDT (2026-09-09T20:50:00Z)" {
		t.Errorf("Describe = %q", got)
	}
	if got := Describe(time.Date(2026, 12, 1, 21, 0, 0, 0, time.UTC), mustLoc(ny)); got != "Tue 2026-12-01 4:00 PM EST (2026-12-01T21:00:00Z)" {
		t.Errorf("Describe winter = %q", got)
	}
	if got := Describe(time.Date(2026, 12, 1, 21, 0, 0, 0, time.UTC), nil); got != "Tue 2026-12-01 9:00 PM UTC (2026-12-01T21:00:00Z)" {
		t.Errorf("Describe nil loc = %q", got)
	}
	if got := FormatUTC(time.Date(2026, 12, 1, 16, 0, 0, 0, mustLoc(ny))); got != "2026-12-01T21:00:00Z" {
		t.Errorf("FormatUTC = %q", got)
	}
}
