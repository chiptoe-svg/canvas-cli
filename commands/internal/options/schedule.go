package options

import (
	"fmt"
	"regexp"
	"strings"
)

// Item kinds accepted by `schedule --type`.
const (
	ScheduleTypeQuiz       = "quiz"
	ScheduleTypeAssignment = "assignment"
	ScheduleTypeAll        = "all"
)

// The three schedule fields, in the names the command uses. They map to
// Canvas unlock_at / due_at / lock_at.
const (
	ScheduleFieldAvailable = "available"
	ScheduleFieldDue       = "due"
	ScheduleFieldClosed    = "closed"
)

// ScheduleOptions contains options for `canvas schedule`.
type ScheduleOptions struct {
	CourseID int64
	ID       int64  // one quiz or assignment by id
	Match    string // or every quiz/assignment whose title matches
	Type     string // quiz | assignment | all

	Available string // new unlock_at, in local time
	Due       string // new due_at
	Closed    string // new lock_at
	Date      string // calendar day for time-only values (default today)
	Clear     []string

	Timezone string
	Force    bool
	DryRun   bool
}

// Validate validates the options.
func (o *ScheduleOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.ID == 0 && o.Match == "" {
		return fmt.Errorf("one of --id or --match is required")
	}
	if o.ID != 0 && o.Match != "" {
		return fmt.Errorf("--id and --match are mutually exclusive")
	}
	if o.ID < 0 {
		return fmt.Errorf("id must be greater than 0, got %d", o.ID)
	}
	switch o.Type {
	case "", ScheduleTypeQuiz, ScheduleTypeAssignment, ScheduleTypeAll:
	default:
		return fmt.Errorf("unknown --type %q (expected quiz, assignment or all)", o.Type)
	}
	if _, err := CompileNameMatch("--match", o.Match); err != nil {
		return err
	}
	clear, err := o.ClearSet()
	if err != nil {
		return err
	}
	for field, value := range map[string]string{
		ScheduleFieldAvailable: o.Available,
		ScheduleFieldDue:       o.Due,
		ScheduleFieldClosed:    o.Closed,
	} {
		if value != "" && clear[field] {
			return fmt.Errorf("--%s and --clear %s are mutually exclusive", field, field)
		}
	}
	if !o.HasChanges() {
		return fmt.Errorf("nothing to change: give at least one of --available, --due, --closed or --clear")
	}
	return nil
}

// ItemType is the normalized --type (default all).
func (o *ScheduleOptions) ItemType() string {
	if o.Type == "" {
		return ScheduleTypeAll
	}
	return o.Type
}

// HasChanges reports whether any field is set or cleared.
func (o *ScheduleOptions) HasChanges() bool {
	return o.Available != "" || o.Due != "" || o.Closed != "" || len(o.Clear) > 0
}

// ClearSet returns the normalized --clear selection as a set, rejecting
// unknown field names. Values may be repeated flags or comma-separated.
func (o *ScheduleOptions) ClearSet() (map[string]bool, error) {
	set := map[string]bool{}
	for _, raw := range o.Clear {
		for _, f := range strings.Split(raw, ",") {
			f = strings.ToLower(strings.TrimSpace(f))
			if f == "" {
				continue
			}
			switch f {
			case ScheduleFieldAvailable, ScheduleFieldDue, ScheduleFieldClosed:
				set[f] = true
			default:
				return nil, fmt.Errorf("unknown --clear value %q (expected available, due or closed)", f)
			}
		}
	}
	return set, nil
}

// CompileNameMatch turns a title-match flag value into a predicate. "/…/"
// is a Go regular expression (case-sensitive unless the pattern says
// otherwise, e.g. "/(?i)quiz/"); anything else is a case-insensitive
// substring. An empty value matches everything.
func CompileNameMatch(flag, value string) (func(string) bool, error) {
	if value == "" {
		return func(string) bool { return true }, nil
	}
	if len(value) >= 2 && strings.HasPrefix(value, "/") && strings.HasSuffix(value, "/") {
		re, err := regexp.Compile(value[1 : len(value)-1])
		if err != nil {
			return nil, fmt.Errorf("invalid %s regex: %w", flag, err)
		}
		return re.MatchString, nil
	}
	needle := strings.ToLower(value)
	return func(name string) bool { return strings.Contains(strings.ToLower(name), needle) }, nil
}
