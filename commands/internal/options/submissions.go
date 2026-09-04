package options

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// SubmissionsListOptions contains options for listing submissions
type SubmissionsListOptions struct {
	CourseID      int64
	AssignmentID  int64
	WorkflowState string
	GradedSince   string
	Include       []string
}

// Validate validates the options
func (o *SubmissionsListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	return nil
}

// SubmissionsGetOptions contains options for getting a submission
type SubmissionsGetOptions struct {
	CourseID     int64
	AssignmentID int64
	UserID       int64
	Include      []string
}

// SubmissionsDownloadOptions contains options for downloading every file
// attached to an assignment's submissions.
type SubmissionsDownloadOptions struct {
	CourseID     int64
	AssignmentID int64
	Destination  string
	Overwrite    bool
}

// Validate validates the options.
func (o *SubmissionsDownloadOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	if o.Destination == "" {
		return fmt.Errorf("destination is required")
	}
	return nil
}

// Validate validates the options
func (o *SubmissionsGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// SubmissionsGradeOptions contains options for grading a submission
type SubmissionsGradeOptions struct {
	CourseID     int64
	AssignmentID int64
	UserID       int64
	Score        float64
	Comment      string
	Excuse       bool
	PostedGrade  string
	// Rubric rows, as repeated `<criterion_id>=<points>` flags, and their
	// optional per-criterion comments as `<criterion_id>=<text>`.
	Rubric        []string
	RubricComment []string
}

// Validate validates the options
func (o *SubmissionsGradeOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	// At least one grading parameter is required
	if o.Score == 0 && o.Comment == "" && !o.Excuse && o.PostedGrade == "" && len(o.Rubric) == 0 {
		return fmt.Errorf("at least one grading parameter is required: score, comment, excuse, posted-grade, or rubric")
	}
	if _, err := o.RubricAssessment(); err != nil {
		return err
	}
	return nil
}

// RubricAssessment turns the repeated --rubric / --rubric-comment flags into
// the per-criterion map Canvas expects on the submission update.
//
// Canvas accepts rubric rows on the same request as the score, which is the
// only route that works everywhere: creating a rubric assessment directly
// needs a rubric-association id, and Canvas exposes no endpoint that lists
// associations, so on many instances that id cannot be discovered at all.
//
// `--rubric <criterion_id>=<points>` (points may be 0, or a decimal)
// `--rubric-comment <criterion_id>=<text>` (text may contain '='; only the
// first '=' splits, and the criterion must also carry a --rubric entry)
func (o *SubmissionsGradeOptions) RubricAssessment() (map[string]api.RubricAssessmentParams, error) {
	if len(o.Rubric) == 0 && len(o.RubricComment) == 0 {
		return nil, nil
	}
	assessment := make(map[string]api.RubricAssessmentParams, len(o.Rubric))
	for _, raw := range o.Rubric {
		id, value, ok := strings.Cut(raw, "=")
		id, value = strings.TrimSpace(id), strings.TrimSpace(value)
		if !ok || id == "" || value == "" {
			return nil, fmt.Errorf("--rubric %q must be <criterion-id>=<points>", raw)
		}
		if _, exists := assessment[id]; exists {
			return nil, fmt.Errorf("--rubric names criterion %q twice", id)
		}
		points, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("--rubric %q: points must be a number", raw)
		}
		if points < 0 {
			return nil, fmt.Errorf("--rubric %q: points cannot be negative", raw)
		}
		scored := points
		assessment[id] = api.RubricAssessmentParams{Points: &scored}
	}
	for _, raw := range o.RubricComment {
		id, text, ok := strings.Cut(raw, "=")
		id = strings.TrimSpace(id)
		if !ok || id == "" || strings.TrimSpace(text) == "" {
			return nil, fmt.Errorf("--rubric-comment %q must be <criterion-id>=<text>", raw)
		}
		entry, exists := assessment[id]
		if !exists {
			// Sending a comment alone would blank that criterion's score on
			// some instances; refuse rather than write a surprise.
			return nil, fmt.Errorf("--rubric-comment names criterion %q, which has no --rubric entry", id)
		}
		entry.Comments = text
		assessment[id] = entry
	}
	return assessment, nil
}

// SubmissionsBulkGradeOptions contains options for bulk grading submissions
type SubmissionsBulkGradeOptions struct {
	CourseID int64
	CSV      string
	DryRun   bool
}

// Validate validates the options
func (o *SubmissionsBulkGradeOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.CSV == "" {
		return fmt.Errorf("csv file is required")
	}
	return nil
}

// SubmissionsCommentsOptions contains options for listing submission comments
type SubmissionsCommentsOptions struct {
	CourseID     int64
	AssignmentID int64
	UserID       int64
}

// Validate validates the options
func (o *SubmissionsCommentsOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// SubmissionsAddCommentOptions contains options for adding a comment to a submission
type SubmissionsAddCommentOptions struct {
	CourseID     int64
	AssignmentID int64
	UserID       int64
	Text         string
	GroupShare   bool
}

// Validate validates the options
func (o *SubmissionsAddCommentOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	if o.Text == "" {
		return fmt.Errorf("comment text is required")
	}
	return nil
}

// SubmissionsDeleteCommentOptions contains options for deleting a submission comment
type SubmissionsDeleteCommentOptions struct {
	CourseID     int64
	AssignmentID int64
	UserID       int64
	CommentID    int64
	Force        bool
}

// Validate validates the options
func (o *SubmissionsDeleteCommentOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	if o.CommentID <= 0 {
		return fmt.Errorf("comment-id is required and must be greater than 0")
	}
	return nil
}
