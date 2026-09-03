package options

import (
	"strings"
	"testing"
	"time"
)

func validMissingOptions() *SubmissionsMissingOptions {
	return &SubmissionsMissingOptions{
		CourseIDs:     []int64{1},
		PublishedOnly: true,
		Types:         DefaultMissingTypes,
		DueBefore:     "now",
	}
}

func TestSubmissionsMissingOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(o *SubmissionsMissingOptions)
		wantErr string
	}{
		{"valid course id", func(o *SubmissionsMissingOptions) {}, ""},
		{"valid all-active", func(o *SubmissionsMissingOptions) { o.CourseIDs = nil; o.AllActive = true }, ""},
		{"neither scope", func(o *SubmissionsMissingOptions) { o.CourseIDs = nil }, "one of --course-id or --all-active"},
		{"both scopes", func(o *SubmissionsMissingOptions) { o.AllActive = true }, "mutually exclusive"},
		{"bad course id", func(o *SubmissionsMissingOptions) { o.CourseIDs = []int64{0} }, "course-id must be greater than 0"},
		{"assignment id and match", func(o *SubmissionsMissingOptions) { o.AssignmentIDs = []int64{5}; o.AssignmentMatch = "x" }, "mutually exclusive"},
		{"bad assignment id", func(o *SubmissionsMissingOptions) { o.AssignmentIDs = []int64{-1} }, "assignment-id must be greater than 0"},
		{"unknown type", func(o *SubmissionsMissingOptions) { o.Types = []string{"assignment", "page"} }, "unknown --types value"},
		{"empty types", func(o *SubmissionsMissingOptions) { o.Types = []string{""} }, "at least one"},
		{"bad regex", func(o *SubmissionsMissingOptions) { o.AssignmentMatch = "/(/" }, "invalid --assignment-match regex"},
		{"bad due-before", func(o *SubmissionsMissingOptions) { o.DueBefore = "yesterday" }, "invalid --due-before"},
		{"bad due-after", func(o *SubmissionsMissingOptions) { o.DueAfter = "03/01/2026" }, "invalid --due-after"},
		{"negative min-missing", func(o *SubmissionsMissingOptions) { o.MinMissing = -1 }, "--min-missing must be 0 or greater"},
		{"comma-joined types", func(o *SubmissionsMissingOptions) { o.Types = []string{"quiz,discussion"} }, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			o := validMissingOptions()
			tc.mutate(o)
			err := o.Validate()
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error = %v, want containing %q", err, tc.wantErr)
			}
		})
	}
}

func TestParseDueBound(t *testing.T) {
	now := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)

	if got, err := ParseDueBound("now", now); err != nil || !got.Equal(now) {
		t.Errorf("now: got %v, %v", got, err)
	}
	if got, err := ParseDueBound("", now); err != nil || !got.Equal(now) {
		t.Errorf("empty: got %v, %v", got, err)
	}
	got, err := ParseDueBound("2026-03-10", now)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 3, 10, 0, 0, 0, 0, time.Local)
	if !got.Equal(want) {
		t.Errorf("date: got %v, want local midnight %v", got, want)
	}
	got, err = ParseDueBound("2026-03-10T23:59:00Z", now)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Equal(time.Date(2026, 3, 10, 23, 59, 0, 0, time.UTC)) {
		t.Errorf("rfc3339: got %v", got)
	}
	if _, err := ParseDueBound("next tuesday", now); err == nil {
		t.Error("expected an error for free text")
	}
}

func TestCompileAssignmentMatch(t *testing.T) {
	any, err := CompileAssignmentMatch("")
	if err != nil || !any("whatever") {
		t.Errorf("empty pattern should match everything: %v", err)
	}
	sub, err := CompileAssignmentMatch("QUIZ")
	if err != nil {
		t.Fatal(err)
	}
	if !sub("Weekly quiz 3") || sub("Essay") {
		t.Error("substring match should be case-insensitive")
	}
	re, err := CompileAssignmentMatch("/^Quiz \\d+$/")
	if err != nil {
		t.Fatal(err)
	}
	if !re("Quiz 3") || re("quiz 3") || re("Quiz three") {
		t.Error("regex match should be applied verbatim")
	}
	if _, err := CompileAssignmentMatch("/[/"); err == nil {
		t.Error("expected an error for an invalid regex")
	}
}

func TestSubmissionsMissingOptions_TypeSet(t *testing.T) {
	o := &SubmissionsMissingOptions{Types: []string{" Quiz ", "discussion,assignment"}}
	set, err := o.TypeSet()
	if err != nil {
		t.Fatal(err)
	}
	if len(set) != 3 || !set["quiz"] || !set["discussion"] || !set["assignment"] {
		t.Errorf("set = %v", set)
	}
}
