package options

import "testing"

// The rubric flags are the only route to rubric scores that works on every
// Canvas instance: posting a rubric assessment on its own needs an association
// id, and Canvas has no endpoint that lists associations.
func TestSubmissionsGradeOptions_RubricAssessment(t *testing.T) {
	opts := SubmissionsGradeOptions{
		CourseID: 1, AssignmentID: 2, UserID: 3, Score: 14,
		Rubric:        []string{"_1234=4", "_5678=0", " _9012 = 2.5 "},
		RubricComment: []string{"_5678=Missing the third review = the biggest gap"},
	}
	got, err := opts.RubricAssessment()
	if err != nil {
		t.Fatalf("RubricAssessment() error = %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("want 3 criteria, got %d: %#v", len(got), got)
	}

	// A criterion scored ZERO must be present with a real zero, not dropped.
	zero, ok := got["_5678"]
	if !ok {
		t.Fatal("criterion scored 0 is missing from the assessment")
	}
	if zero.Points == nil {
		t.Fatal("criterion scored 0 has no points set — Canvas would leave it unchanged")
	}
	if *zero.Points != 0 {
		t.Errorf("zero criterion points = %v, want 0", *zero.Points)
	}
	// Only the FIRST '=' splits, so a comment may contain one.
	if want := "Missing the third review = the biggest gap"; zero.Comments != want {
		t.Errorf("comment = %q, want %q", zero.Comments, want)
	}
	if got["_9012"].Points == nil || *got["_9012"].Points != 2.5 {
		t.Errorf("decimal points not parsed: %#v", got["_9012"])
	}
	if got["_1234"].Comments != "" {
		t.Errorf("criterion without a comment got one: %q", got["_1234"].Comments)
	}
}

func TestSubmissionsGradeOptions_RubricRejects(t *testing.T) {
	base := func() SubmissionsGradeOptions {
		return SubmissionsGradeOptions{CourseID: 1, AssignmentID: 2, UserID: 3, Score: 10}
	}
	tests := []struct {
		name string
		mut  func(*SubmissionsGradeOptions)
	}{
		{"no equals sign", func(o *SubmissionsGradeOptions) { o.Rubric = []string{"_1234"} }},
		{"empty points", func(o *SubmissionsGradeOptions) { o.Rubric = []string{"_1234="} }},
		{"non-numeric points", func(o *SubmissionsGradeOptions) { o.Rubric = []string{"_1234=full"} }},
		{"negative points", func(o *SubmissionsGradeOptions) { o.Rubric = []string{"_1234=-1"} }},
		{"duplicate criterion", func(o *SubmissionsGradeOptions) { o.Rubric = []string{"_1234=1", "_1234=2"} }},
		{"comment without a score", func(o *SubmissionsGradeOptions) {
			o.Rubric = []string{"_1234=1"}
			o.RubricComment = []string{"_9999=orphan"}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := base()
			tt.mut(&opts)
			if _, err := opts.RubricAssessment(); err == nil {
				t.Fatal("want an error, got none")
			}
			if err := opts.Validate(); err == nil {
				t.Fatal("Validate must reject it too, before any request is sent")
			}
		})
	}
}

// A rubric alone is a complete grading request — Canvas computes the total.
func TestSubmissionsGradeOptions_RubricAloneIsEnough(t *testing.T) {
	opts := SubmissionsGradeOptions{CourseID: 1, AssignmentID: 2, UserID: 3, Rubric: []string{"_1234=4"}}
	if err := opts.Validate(); err != nil {
		t.Fatalf("rubric-only grade rejected: %v", err)
	}
}
