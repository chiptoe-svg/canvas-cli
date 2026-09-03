package options

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// AssignmentsUpcomingOptions contains options for the read-only
// `assignments upcoming` report.
type AssignmentsUpcomingOptions struct {
	CourseIDs []int64
	AllActive bool

	Within string // 10d, 2w, 36h
	By     string // a local date/time; date-only means end of that day

	IncludeUndated bool
	PublishedOnly  bool
	Timezone       string
}

// Validate validates the options.
func (o *AssignmentsUpcomingOptions) Validate() error {
	if len(o.CourseIDs) == 0 && !o.AllActive {
		return fmt.Errorf("one of --course-id or --all-active is required")
	}
	if len(o.CourseIDs) > 0 && o.AllActive {
		return fmt.Errorf("--course-id and --all-active are mutually exclusive")
	}
	for _, id := range o.CourseIDs {
		if id <= 0 {
			return fmt.Errorf("course-id must be greater than 0, got %d", id)
		}
	}
	if o.Within == "" && o.By == "" {
		return fmt.Errorf("one of --within or --by is required (e.g. --within 10d, --by \"this sunday\")")
	}
	if o.Within != "" && o.By != "" {
		return fmt.Errorf("--within and --by are mutually exclusive")
	}
	if o.Within != "" {
		if _, err := ParseWithin(o.Within); err != nil {
			return err
		}
	}
	return nil
}

var withinRe = regexp.MustCompile(`^(\d+)\s*(h|d|w)$`)

// ParseWithin reads a --within span: <n>h (hours), <n>d (days) or <n>w
// (weeks), e.g. 36h, 10d, 2w. Days are 24 hours; the limit is an instant,
// not a calendar day.
func ParseWithin(value string) (time.Duration, error) {
	m := withinRe.FindStringSubmatch(value)
	if m == nil {
		return 0, fmt.Errorf("invalid --within %q: use <n>h, <n>d or <n>w (e.g. 36h, 10d, 2w)", value)
	}
	n, err := strconv.Atoi(m[1])
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("invalid --within %q: the number must be greater than 0", value)
	}
	switch m[2] {
	case "h":
		return time.Duration(n) * time.Hour, nil
	case "d":
		return time.Duration(n) * 24 * time.Hour, nil
	default:
		return time.Duration(n) * 7 * 24 * time.Hour, nil
	}
}
