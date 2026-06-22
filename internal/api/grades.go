package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// GradesService handles grade-related API calls
type GradesService struct {
	client *Client
}

// NewGradesService creates a new grades service
func NewGradesService(client *Client) *GradesService {
	return &GradesService{client: client}
}

// GradebookHistoryGrader represents a grader in gradebook history
type GradebookHistoryGrader struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GradebookHistoryDay represents a day in gradebook history
type GradebookHistoryDay struct {
	Date    string                   `json:"date"`
	Graders []GradebookHistoryGrader `json:"graders"`
}

// GradebookHistoryEntry represents an entry in gradebook history feed
type GradebookHistoryEntry struct {
	ID                int64      `json:"id"`
	UserID            int64      `json:"user_id"`
	UserName          string     `json:"user_name"`
	AssignmentID      int64      `json:"assignment_id"`
	AssignmentName    string     `json:"assignment_name"`
	CurrentGrade      string     `json:"current_grade"`
	CurrentGradedAt   *time.Time `json:"current_graded_at,omitempty"`
	NewGrade          string     `json:"new_grade"`
	NewGradedAt       *time.Time `json:"new_graded_at,omitempty"`
	GraderID          int64      `json:"grader_id"`
	GraderName        string     `json:"grader_name"`
	Excused           bool       `json:"excused"`
	GradedAnonymously bool       `json:"graded_anonymously"`
	PointsPossible    float64    `json:"points_possible"`
	CurrentScore      float64    `json:"current_score"`
	NewScore          float64    `json:"new_score"`
	GradeBefore       string     `json:"grade_before"`
	ScoreBefore       float64    `json:"score_before"`
	ExcusedBefore     bool       `json:"excused_before"`
	RequestID         string     `json:"request_id,omitempty"`
}

// CustomGradebookColumn represents a custom gradebook column
type CustomGradebookColumn struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Position     int    `json:"position"`
	TeacherNotes bool   `json:"teacher_notes"`
	ReadOnly     bool   `json:"read_only"`
	Hidden       bool   `json:"hidden"`
}

// CustomColumnDatum represents data in a custom column for a user
type CustomColumnDatum struct {
	ColumnID int64  `json:"column_id,omitempty"`
	UserID   int64  `json:"user_id"`
	Content  string `json:"content"`
}

// ListGradebookHistoryOptions holds options for listing gradebook history
type ListGradebookHistoryOptions struct {
	CourseID  int64
	StartDate string
	EndDate   string
	Page      int
	PerPage   int
}

// GetHistory retrieves gradebook history days
func (s *GradesService) GetHistory(ctx context.Context, courseID int64, opts *ListGradebookHistoryOptions) ([]GradebookHistoryDay, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/gradebook_history/days", courseID)

	if opts != nil {
		query := url.Values{}

		if opts.StartDate != "" {
			query.Add("start_date", opts.StartDate)
		}

		if opts.EndDate != "" {
			query.Add("end_date", opts.EndDate)
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var days []GradebookHistoryDay
	if err := s.client.GetAllPages(ctx, path, &days); err != nil {
		return nil, err
	}

	return days, nil
}

// ListGradebookFeedOptions holds options for listing gradebook feed
type ListGradebookFeedOptions struct {
	UserID       int64
	AssignmentID int64
	StartDate    string
	EndDate      string
	Page         int
	PerPage      int
}

// GetFeed retrieves gradebook history feed
func (s *GradesService) GetFeed(ctx context.Context, courseID int64, opts *ListGradebookFeedOptions) ([]GradebookHistoryEntry, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/gradebook_history/feed", courseID)

	if opts != nil {
		query := url.Values{}

		if opts.UserID > 0 {
			query.Add("user_id", strconv.FormatInt(opts.UserID, 10))
		}

		if opts.AssignmentID > 0 {
			query.Add("assignment_id", strconv.FormatInt(opts.AssignmentID, 10))
		}

		if opts.StartDate != "" {
			query.Add("start_date", opts.StartDate)
		}

		if opts.EndDate != "" {
			query.Add("end_date", opts.EndDate)
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var entries []GradebookHistoryEntry
	if err := s.client.GetAllPages(ctx, path, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// GradebookHistorySubmission is a submission record returned by the
// per-date/grader/assignment gradebook history endpoint.
type GradebookHistorySubmission struct {
	SubmissionID int64   `json:"submission_id"`
	UserID       int64   `json:"user_id"`
	UserName     string  `json:"user_name"`
	GradedAt     string  `json:"graded_at"`
	Grade        string  `json:"grade,omitempty"`
	Score        float64 `json:"score,omitempty"`
	AssignmentID int64   `json:"assignment_id"`
}

// GetGradebookHistoryDay returns graders active on a specific date.
// GET /api/v1/courses/:course_id/gradebook_history/:date
func (s *GradesService) GetGradebookHistoryDay(ctx context.Context, courseID int64, date string) ([]GradebookHistoryGrader, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/gradebook_history/%s", courseID, date)
	var out []GradebookHistoryGrader
	if err := s.client.GetAllPages(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("getting gradebook history for course %d on %s: %w", courseID, date, err)
	}
	return out, nil
}

// GetGradebookHistorySubmissions lists submissions graded on a date by a grader for an assignment.
// GET /api/v1/courses/:course_id/gradebook_history/:date/graders/:grader_id/assignments/:assignment_id/submissions
func (s *GradesService) GetGradebookHistorySubmissions(ctx context.Context, courseID int64, date string, graderID, assignmentID int64) ([]GradebookHistorySubmission, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/gradebook_history/%s/graders/%d/assignments/%d/submissions", courseID, date, graderID, assignmentID)
	var out []GradebookHistorySubmission
	if err := s.client.GetAllPages(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("getting gradebook history submissions for course %d: %w", courseID, err)
	}
	return out, nil
}

// ListCustomColumnsOptions holds options for listing custom columns
type ListCustomColumnsOptions struct {
	IncludeHidden bool
	Page          int
	PerPage       int
}

// ListCustomColumns retrieves custom gradebook columns
func (s *GradesService) ListCustomColumns(ctx context.Context, courseID int64, opts *ListCustomColumnsOptions) ([]CustomGradebookColumn, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/custom_gradebook_columns", courseID)

	if opts != nil {
		query := url.Values{}

		if opts.IncludeHidden {
			query.Add("include_hidden", "true")
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var columns []CustomGradebookColumn
	if err := s.client.GetAllPages(ctx, path, &columns); err != nil {
		return nil, err
	}

	return columns, nil
}

// GetCustomColumn retrieves a single custom column
func (s *GradesService) GetCustomColumn(ctx context.Context, courseID, columnID int64) (*CustomGradebookColumn, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/custom_gradebook_columns/%d", courseID, columnID)

	var column CustomGradebookColumn
	if err := s.client.GetJSON(ctx, path, &column); err != nil {
		return nil, err
	}

	return &column, nil
}

// CreateCustomColumnParams holds parameters for creating a custom column
type CreateCustomColumnParams struct {
	Title        string
	Position     int
	Hidden       bool
	TeacherNotes bool
	ReadOnly     bool
}

// CreateCustomColumn creates a new custom gradebook column
func (s *GradesService) CreateCustomColumn(ctx context.Context, courseID int64, params *CreateCustomColumnParams) (*CustomGradebookColumn, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/custom_gradebook_columns", courseID)

	body := map[string]interface{}{
		"column": make(map[string]interface{}),
	}

	columnData, ok := body["column"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid column data structure")
	}

	if params.Title != "" {
		columnData["title"] = params.Title
	}

	if params.Position > 0 {
		columnData["position"] = params.Position
	}

	if params.Hidden {
		columnData["hidden"] = params.Hidden
	}

	if params.TeacherNotes {
		columnData["teacher_notes"] = params.TeacherNotes
	}

	if params.ReadOnly {
		columnData["read_only"] = params.ReadOnly
	}

	var column CustomGradebookColumn
	if err := s.client.PostJSON(ctx, path, body, &column); err != nil {
		return nil, err
	}

	return &column, nil
}

// UpdateCustomColumnParams holds parameters for updating a custom column
type UpdateCustomColumnParams struct {
	Title        *string
	Position     *int
	Hidden       *bool
	TeacherNotes *bool
	ReadOnly     *bool
}

// UpdateCustomColumn updates an existing custom column
func (s *GradesService) UpdateCustomColumn(ctx context.Context, courseID, columnID int64, params *UpdateCustomColumnParams) (*CustomGradebookColumn, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/custom_gradebook_columns/%d", courseID, columnID)

	body := map[string]interface{}{
		"column": make(map[string]interface{}),
	}

	columnData, ok := body["column"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid column data structure")
	}

	if params.Title != nil {
		columnData["title"] = *params.Title
	}

	if params.Position != nil {
		columnData["position"] = *params.Position
	}

	if params.Hidden != nil {
		columnData["hidden"] = *params.Hidden
	}

	if params.TeacherNotes != nil {
		columnData["teacher_notes"] = *params.TeacherNotes
	}

	if params.ReadOnly != nil {
		columnData["read_only"] = *params.ReadOnly
	}

	var column CustomGradebookColumn
	if err := s.client.PutJSON(ctx, path, body, &column); err != nil {
		return nil, err
	}

	return &column, nil
}

// DeleteCustomColumn deletes a custom column
func (s *GradesService) DeleteCustomColumn(ctx context.Context, courseID, columnID int64) (*CustomGradebookColumn, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/custom_gradebook_columns/%d", courseID, columnID)

	var column CustomGradebookColumn
	if err := s.client.DeleteJSON(ctx, path, &column); err != nil {
		return nil, err
	}

	return &column, nil
}

// GetCustomColumnData retrieves data for a custom column
func (s *GradesService) GetCustomColumnData(ctx context.Context, courseID, columnID int64) ([]CustomColumnDatum, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/custom_gradebook_columns/%d/data", courseID, columnID)

	var data []CustomColumnDatum
	if err := s.client.GetAllPages(ctx, path, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// SetCustomColumnData sets data for a user in a custom column
func (s *GradesService) SetCustomColumnData(ctx context.Context, courseID, columnID, userID int64, content string) (*CustomColumnDatum, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/custom_gradebook_columns/%d/data/%d", courseID, columnID, userID)

	body := map[string]interface{}{
		"column_data": map[string]interface{}{
			"content": content,
		},
	}

	var datum CustomColumnDatum
	if err := s.client.PutJSON(ctx, path, body, &datum); err != nil {
		return nil, err
	}

	// Set column_id since the Canvas API doesn't return it
	if datum.ColumnID == 0 {
		datum.ColumnID = columnID
	}

	return &datum, nil
}

// BulkUpdateGrade represents a grade update in bulk operations
type BulkUpdateGrade struct {
	StudentID    int64
	AssignmentID int64
	Grade        string
	Excused      bool
}

// BulkUpdateGrades updates multiple grades at once
func (s *GradesService) BulkUpdateGrades(ctx context.Context, courseID int64, grades []BulkUpdateGrade) error {
	path := fmt.Sprintf("/api/v1/courses/%d/submissions/update_grades", courseID)

	gradeUpdates := make(map[string]map[string]interface{})

	for _, g := range grades {
		studentKey := strconv.FormatInt(g.StudentID, 10)
		if _, exists := gradeUpdates[studentKey]; !exists {
			gradeUpdates[studentKey] = make(map[string]interface{})
		}

		assignmentKey := strconv.FormatInt(g.AssignmentID, 10)
		updateData := map[string]interface{}{}

		if g.Grade != "" {
			updateData["posted_grade"] = g.Grade
		}
		if g.Excused {
			updateData["excuse"] = true
		}

		gradeUpdates[studentKey][assignmentKey] = updateData
	}

	body := map[string]interface{}{
		"grade_data": gradeUpdates,
	}

	return s.client.PostJSON(ctx, path, body, nil)
}
