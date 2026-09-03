package options

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Assignment kinds accepted by --types.
const (
	MissingTypeAssignment = "assignment"
	MissingTypeQuiz       = "quiz"
	MissingTypeDiscussion = "discussion"
)

// DefaultMissingTypes is the --types default: every submittable kind.
var DefaultMissingTypes = []string{MissingTypeAssignment, MissingTypeQuiz, MissingTypeDiscussion}

// SubmissionsMissingOptions contains options for the read-only
// `submissions missing` report.
type SubmissionsMissingOptions struct {
	CourseIDs []int64
	AllActive bool

	AssignmentIDs   []int64
	AssignmentMatch string

	PublishedOnly     bool
	Types             []string
	ExcludeZeroPoints bool
	DueBefore         string
	DueAfter          string
	IncludeUndated    bool

	IncludeInactive bool

	ZeroOnly   bool
	MinMissing int
}

// Validate validates the options.
func (o *SubmissionsMissingOptions) Validate() error {
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
	if len(o.AssignmentIDs) > 0 && o.AssignmentMatch != "" {
		return fmt.Errorf("--assignment-id and --assignment-match are mutually exclusive")
	}
	for _, id := range o.AssignmentIDs {
		if id <= 0 {
			return fmt.Errorf("assignment-id must be greater than 0, got %d", id)
		}
	}
	if _, err := o.TypeSet(); err != nil {
		return err
	}
	if _, err := CompileAssignmentMatch(o.AssignmentMatch); err != nil {
		return err
	}
	now := time.Now()
	if _, err := ParseDueBound(o.DueBefore, now); err != nil {
		return fmt.Errorf("invalid --due-before: %w", err)
	}
	if o.DueAfter != "" {
		if _, err := ParseDueBound(o.DueAfter, now); err != nil {
			return fmt.Errorf("invalid --due-after: %w", err)
		}
	}
	if o.MinMissing < 0 {
		return fmt.Errorf("--min-missing must be 0 or greater")
	}
	return nil
}

// TypeSet returns the normalized --types selection as a set, rejecting
// unknown kinds.
func (o *SubmissionsMissingOptions) TypeSet() (map[string]bool, error) {
	set := map[string]bool{}
	for _, raw := range o.Types {
		for _, t := range strings.Split(raw, ",") {
			t = strings.ToLower(strings.TrimSpace(t))
			if t == "" {
				continue
			}
			switch t {
			case MissingTypeAssignment, MissingTypeQuiz, MissingTypeDiscussion:
				set[t] = true
			default:
				return nil, fmt.Errorf("unknown --types value %q (expected assignment, quiz, discussion)", t)
			}
		}
	}
	if len(set) == 0 {
		return nil, fmt.Errorf("--types must name at least one of assignment, quiz, discussion")
	}
	return set, nil
}

// ParseDueBound parses a --due-before / --due-after value: "now" (or empty)
// means now, "YYYY-MM-DD" means local midnight at the start of that day, and
// anything else must be RFC 3339.
func ParseDueBound(value string, now time.Time) (time.Time, error) {
	v := strings.TrimSpace(value)
	if v == "" || strings.EqualFold(v, "now") {
		return now, nil
	}
	if t, err := time.ParseInLocation("2006-01-02", v, time.Local); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC3339, v); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("%q is not \"now\", YYYY-MM-DD, or RFC 3339", value)
}

// CompileAssignmentMatch turns an --assignment-match value into a predicate on
// assignment names. "/…/" is a Go regular expression (case-sensitive unless
// the pattern says otherwise); anything else is a case-insensitive substring.
// An empty value matches everything.
func CompileAssignmentMatch(value string) (func(string) bool, error) {
	if value == "" {
		return func(string) bool { return true }, nil
	}
	if len(value) >= 2 && strings.HasPrefix(value, "/") && strings.HasSuffix(value, "/") {
		re, err := regexp.Compile(value[1 : len(value)-1])
		if err != nil {
			return nil, fmt.Errorf("invalid --assignment-match regex: %w", err)
		}
		return re.MatchString, nil
	}
	needle := strings.ToLower(value)
	return func(name string) bool { return strings.Contains(strings.ToLower(name), needle) }, nil
}
