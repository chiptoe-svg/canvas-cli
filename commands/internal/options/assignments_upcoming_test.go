package options

import (
	"strings"
	"testing"
	"time"
)

func TestAssignmentsUpcomingOptions_Validate(t *testing.T) {
	valid := func() *AssignmentsUpcomingOptions {
		return &AssignmentsUpcomingOptions{CourseIDs: []int64{1}, Within: "10d", PublishedOnly: true}
	}
	cases := []struct {
		name    string
		mutate  func(o *AssignmentsUpcomingOptions)
		wantErr string
	}{
		{"valid within", func(o *AssignmentsUpcomingOptions) {}, ""},
		{"valid by", func(o *AssignmentsUpcomingOptions) { o.Within = ""; o.By = "this sunday" }, ""},
		{"valid all-active", func(o *AssignmentsUpcomingOptions) { o.CourseIDs = nil; o.AllActive = true }, ""},
		{"no scope", func(o *AssignmentsUpcomingOptions) { o.CourseIDs = nil }, "one of --course-id or --all-active"},
		{"both scopes", func(o *AssignmentsUpcomingOptions) { o.AllActive = true }, "mutually exclusive"},
		{"bad course", func(o *AssignmentsUpcomingOptions) { o.CourseIDs = []int64{0} }, "course-id must be greater than 0"},
		{"no window", func(o *AssignmentsUpcomingOptions) { o.Within = "" }, "one of --within or --by"},
		{"both windows", func(o *AssignmentsUpcomingOptions) { o.By = "tomorrow" }, "mutually exclusive"},
		{"bad within", func(o *AssignmentsUpcomingOptions) { o.Within = "10 days" }, "invalid --within"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			o := valid()
			c.mutate(o)
			err := o.Validate()
			if c.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), c.wantErr) {
				t.Fatalf("error = %v, want %q", err, c.wantErr)
			}
		})
	}
}

func TestParseWithin(t *testing.T) {
	for in, want := range map[string]time.Duration{
		"36h": 36 * time.Hour,
		"10d": 240 * time.Hour,
		"2w":  14 * 24 * time.Hour,
		"1 d": 24 * time.Hour,
	} {
		got, err := ParseWithin(in)
		if err != nil || got != want {
			t.Errorf("ParseWithin(%q) = %v, %v; want %v", in, got, err, want)
		}
	}
	for _, in := range []string{"", "0d", "-1d", "10", "10m", "d10", "10days"} {
		if _, err := ParseWithin(in); err == nil {
			t.Errorf("ParseWithin(%q) accepted", in)
		}
	}
}
