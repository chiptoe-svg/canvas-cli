package api

import (
	"context"
	"fmt"
)

// CoursePacingService handles course pacing (pace) API calls.
type CoursePacingService struct {
	client *Client
}

// NewCoursePacingService creates a new course pacing service.
func NewCoursePacingService(client *Client) *CoursePacingService {
	return &CoursePacingService{client: client}
}

// CoursePace represents a Canvas course pace object.
type CoursePace struct {
	ID                    int64            `json:"id"`
	CourseID              int64            `json:"course_id"`
	WorkflowState         string           `json:"workflow_state"`
	ExcludeWeekends       bool             `json:"exclude_weekends"`
	HardEndDates          bool             `json:"hard_end_dates"`
	CreatedAt             string           `json:"created_at,omitempty"`
	UpdatedAt             string           `json:"updated_at,omitempty"`
	PublishedAt           string           `json:"published_at,omitempty"`
	EndDate               string           `json:"end_date,omitempty"`
	CoursePaceModuleItems []CoursePaceItem `json:"course_pace_module_items,omitempty"`
}

// CoursePaceItem represents a module item within a course pace.
type CoursePaceItem struct {
	ID            int64 `json:"id"`
	CoursePaceID  int64 `json:"course_pace_id"`
	ModuleItemID  int64 `json:"module_item_id"`
	DurationDays  int   `json:"duration"`
	RootAccountID int64 `json:"root_account_id,omitempty"`
}

// CoursePaceParams holds mutable fields for creating/updating a course pace.
// ExcludeWeekends and HardEndDates use *bool so that an explicit false
// is not silently dropped by json's omitempty zero-value rule.
type CoursePaceParams struct {
	ExcludeWeekends *bool  `json:"exclude_weekends,omitempty"`
	HardEndDates    *bool  `json:"hard_end_dates,omitempty"`
	EndDate         string `json:"end_date,omitempty"`
}

// coursePaceEnvelope unwraps the Canvas {"course_pace": {...}} envelope.
type coursePaceEnvelope struct {
	CoursePace CoursePace `json:"course_pace"`
}

// GetCoursePace retrieves a course pace by ID.
// GET /api/v1/courses/:course_id/course_pacing/:id
func (s *CoursePacingService) Get(ctx context.Context, courseID, paceID int64) (*CoursePace, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/course_pacing/%d", courseID, paceID)
	var envelope coursePaceEnvelope
	if err := s.client.GetJSON(ctx, path, &envelope); err != nil {
		return nil, fmt.Errorf("getting course pace %d for course %d: %w", paceID, courseID, err)
	}
	return &envelope.CoursePace, nil
}

// CreateCoursePace creates a new course pace.
// POST /api/v1/courses/:course_id/course_pacing
func (s *CoursePacingService) Create(ctx context.Context, courseID int64, params CoursePaceParams) (*CoursePace, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/course_pacing", courseID)
	body := map[string]interface{}{
		"course_pace": params,
	}
	var envelope coursePaceEnvelope
	if err := s.client.PostJSON(ctx, path, body, &envelope); err != nil {
		return nil, fmt.Errorf("creating course pace for course %d: %w", courseID, err)
	}
	return &envelope.CoursePace, nil
}

// UpdateCoursePace updates an existing course pace.
// PUT /api/v1/courses/:course_id/course_pacing/:id
func (s *CoursePacingService) Update(ctx context.Context, courseID, paceID int64, params CoursePaceParams) (*CoursePace, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/course_pacing/%d", courseID, paceID)
	body := map[string]interface{}{
		"course_pace": params,
	}
	var envelope coursePaceEnvelope
	if err := s.client.PutJSON(ctx, path, body, &envelope); err != nil {
		return nil, fmt.Errorf("updating course pace %d for course %d: %w", paceID, courseID, err)
	}
	return &envelope.CoursePace, nil
}

// DeleteCoursePace deletes a course pace.
// DELETE /api/v1/courses/:course_id/course_pacing/:id
func (s *CoursePacingService) Delete(ctx context.Context, courseID, paceID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/course_pacing/%d", courseID, paceID)
	if _, err := s.client.Delete(ctx, path); err != nil {
		return fmt.Errorf("deleting course pace %d for course %d: %w", paceID, courseID, err)
	}
	return nil
}
