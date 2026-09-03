package options

import (
	"strings"
	"testing"
)

func TestSubmissionsExcuseOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *SubmissionsExcuseOptions
		wantErr string
	}{
		{"valid", &SubmissionsExcuseOptions{CourseID: 1, Student: "ada", Assignment: "quiz 3"}, ""},
		{"valid ids", &SubmissionsExcuseOptions{CourseID: 1, Student: "10", Assignment: "456", Unexcuse: true}, ""},
		{"no course", &SubmissionsExcuseOptions{Student: "ada", Assignment: "quiz 3"}, "course-id is required"},
		{"no student", &SubmissionsExcuseOptions{CourseID: 1, Student: "  ", Assignment: "quiz 3"}, "--student is required"},
		{"no assignment", &SubmissionsExcuseOptions{CourseID: 1, Student: "ada"}, "--assignment is required"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.opts.Validate()
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
	if !(&SubmissionsExcuseOptions{}).Excused() || (&SubmissionsExcuseOptions{Unexcuse: true}).Excused() {
		t.Error("Excused")
	}
}
