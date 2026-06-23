package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// QuizAssignmentOverridesService handles quiz assignment override-related API calls.
type QuizAssignmentOverridesService struct {
	client *Client
}

// NewQuizAssignmentOverridesService creates a new quiz assignment overrides service.
func NewQuizAssignmentOverridesService(client *Client) *QuizAssignmentOverridesService {
	return &QuizAssignmentOverridesService{client: client}
}

// QuizAssignmentOverrideSet holds per-quiz assignment overrides.
type QuizAssignmentOverrideSet struct {
	QuizID    string                    `json:"quiz_id"`
	DueDate   *AssignmentDueDate        `json:"due_dates,omitempty"`
	Overrides []AssignmentOverrideEntry `json:"overrides,omitempty"`
}

// AssignmentDueDate holds base due-date information.
type AssignmentDueDate struct {
	DueAt    string `json:"due_at,omitempty"`
	LockAt   string `json:"lock_at,omitempty"`
	UnlockAt string `json:"unlock_at,omitempty"`
}

// AssignmentOverrideEntry holds a single override entry.
type AssignmentOverrideEntry struct {
	ID               int64    `json:"id,omitempty"`
	AssignmentID     int64    `json:"assignment_id,omitempty"`
	Title            string   `json:"title,omitempty"`
	StudentIDs       []int64  `json:"student_ids,omitempty"`
	GroupID          int64    `json:"group_id,omitempty"`
	CourseSectionID  int64    `json:"course_section_id,omitempty"`
	DueAt            string   `json:"due_at,omitempty"`
	LockAt           string   `json:"lock_at,omitempty"`
	UnlockAt         string   `json:"unlock_at,omitempty"`
	AllDay           bool     `json:"all_day,omitempty"`
	AllDayDate       string   `json:"all_day_date,omitempty"`
	CourseSectionIDs []int64  `json:"course_section_ids,omitempty"`
	StudentNames     []string `json:"student_names,omitempty"`
}

// QuizAssignmentOverridesResponse wraps the Canvas API response.
type QuizAssignmentOverridesResponse struct {
	QuizAssignmentOverrides []QuizAssignmentOverrideSet `json:"quiz_assignment_overrides"`
}

// List retrieves assignment overrides for a set of quizzes.
// Canvas API: GET /api/v1/courses/:course_id/quizzes/assignment_overrides
func (s *QuizAssignmentOverridesService) List(ctx context.Context, courseID int64, quizIDs []int64) ([]QuizAssignmentOverrideSet, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/assignment_overrides", courseID)

	if len(quizIDs) > 0 {
		query := url.Values{}
		for _, id := range quizIDs {
			query.Add("quiz_assignment_overrides[0][quiz_ids][]", strconv.FormatInt(id, 10))
		}
		path += "?" + query.Encode()
	}

	var resp QuizAssignmentOverridesResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.QuizAssignmentOverrides, nil
}

// SetQuizAssignmentOverridesParams holds parameters for setting quiz assignment overrides.
type SetQuizAssignmentOverridesParams struct {
	// QuizAssignmentOverrides maps quiz IDs to their override entries.
	QuizAssignmentOverrides []QuizAssignmentOverrideSetInput `json:"quiz_assignment_overrides"`
}

// QuizAssignmentOverrideSetInput holds overrides for a single quiz ID.
type QuizAssignmentOverrideSetInput struct {
	QuizID    string                    `json:"quiz_id"`
	Overrides []AssignmentOverrideEntry `json:"overrides,omitempty"`
}

// Set creates or updates assignment overrides for quizzes.
// Canvas API: POST /api/v1/courses/:course_id/quizzes/assignment_overrides
func (s *QuizAssignmentOverridesService) Set(ctx context.Context, courseID int64, params *SetQuizAssignmentOverridesParams) ([]QuizAssignmentOverrideSet, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/assignment_overrides", courseID)

	var resp QuizAssignmentOverridesResponse
	if err := s.client.PostJSON(ctx, path, params, &resp); err != nil {
		return nil, err
	}

	return resp.QuizAssignmentOverrides, nil
}
