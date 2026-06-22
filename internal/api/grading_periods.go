package api

import (
	"context"
	"fmt"
)

// GradingPeriodsService handles grading period API calls for courses.
type GradingPeriodsService struct {
	client *Client
}

// NewGradingPeriodsService creates a new grading periods service.
func NewGradingPeriodsService(client *Client) *GradingPeriodsService {
	return &GradingPeriodsService{client: client}
}

// GradingPeriod represents a Canvas grading period.
type GradingPeriod struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
	CloseDate string  `json:"close_date,omitempty"`
	Weight    float64 `json:"weight,omitempty"`
	IsLast    bool    `json:"is_last,omitempty"`
	IsClosed  bool    `json:"is_closed,omitempty"`
}

// gradingPeriodsEnvelope wraps the Canvas envelope for grading period responses.
// Canvas always wraps grading periods in {"grading_periods": [...]} even for single objects.
type gradingPeriodsEnvelope struct {
	GradingPeriods []GradingPeriod `json:"grading_periods"`
}

// GradingPeriodParams holds mutable fields for create/update operations.
type GradingPeriodParams struct {
	Title     string  `json:"title,omitempty"`
	StartDate string  `json:"start_date,omitempty"`
	EndDate   string  `json:"end_date,omitempty"`
	CloseDate string  `json:"close_date,omitempty"`
	Weight    float64 `json:"weight,omitempty"`
}

// ListGradingPeriods lists all grading periods for a course.
// GET /api/v1/courses/:course_id/grading_periods
func (s *GradingPeriodsService) List(ctx context.Context, courseID int64) ([]GradingPeriod, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/grading_periods", courseID)
	var envelope gradingPeriodsEnvelope
	if err := s.client.GetJSON(ctx, path, &envelope); err != nil {
		return nil, fmt.Errorf("listing grading periods for course %d: %w", courseID, err)
	}
	return envelope.GradingPeriods, nil
}

// GetGradingPeriod retrieves a single grading period.
// GET /api/v1/courses/:course_id/grading_periods/:id
func (s *GradingPeriodsService) Get(ctx context.Context, courseID, id int64) (*GradingPeriod, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/grading_periods/%d", courseID, id)
	var envelope gradingPeriodsEnvelope
	if err := s.client.GetJSON(ctx, path, &envelope); err != nil {
		return nil, fmt.Errorf("getting grading period %d for course %d: %w", id, courseID, err)
	}
	if len(envelope.GradingPeriods) == 0 {
		return nil, fmt.Errorf("grading period %d not found", id)
	}
	return &envelope.GradingPeriods[0], nil
}

// UpdateGradingPeriod updates a grading period.
// PUT /api/v1/courses/:course_id/grading_periods/:id
func (s *GradingPeriodsService) Update(ctx context.Context, courseID, id int64, params GradingPeriodParams) (*GradingPeriod, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/grading_periods/%d", courseID, id)
	body := map[string]interface{}{
		"grading_periods": []GradingPeriodParams{params},
	}
	var envelope gradingPeriodsEnvelope
	if err := s.client.PutJSON(ctx, path, body, &envelope); err != nil {
		return nil, fmt.Errorf("updating grading period %d for course %d: %w", id, courseID, err)
	}
	if len(envelope.GradingPeriods) == 0 {
		return nil, fmt.Errorf("unexpected empty response updating grading period %d", id)
	}
	return &envelope.GradingPeriods[0], nil
}

// DeleteGradingPeriod deletes a grading period.
// DELETE /api/v1/courses/:course_id/grading_periods/:id
func (s *GradingPeriodsService) Delete(ctx context.Context, courseID, id int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/grading_periods/%d", courseID, id)
	if _, err := s.client.Delete(ctx, path); err != nil {
		return fmt.Errorf("deleting grading period %d for course %d: %w", id, courseID, err)
	}
	return nil
}

// BatchUpdateGradingPeriods updates multiple grading periods in a single request.
// PATCH /api/v1/courses/:course_id/grading_periods/batch_update
func (s *GradingPeriodsService) BatchUpdate(ctx context.Context, courseID int64, periods []GradingPeriodParams) ([]GradingPeriod, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/grading_periods/batch_update", courseID)
	body := map[string]interface{}{
		"grading_periods": periods,
	}
	var envelope gradingPeriodsEnvelope
	if err := s.client.PatchJSON(ctx, path, body, &envelope); err != nil {
		return nil, fmt.Errorf("batch-updating grading periods for course %d: %w", courseID, err)
	}
	return envelope.GradingPeriods, nil
}
