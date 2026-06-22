package api

import (
	"context"
	"fmt"
)

// GradingStandardsService handles grading standard API calls for courses.
type GradingStandardsService struct {
	client *Client
}

// NewGradingStandardsService creates a new grading standards service.
func NewGradingStandardsService(client *Client) *GradingStandardsService {
	return &GradingStandardsService{client: client}
}

// GradingStandard represents a Canvas grading standard (grading scheme).
type GradingStandard struct {
	ID            int64                `json:"id"`
	Title         string               `json:"title"`
	ContextType   string               `json:"context_type"`
	ContextID     int64                `json:"context_id"`
	GradingScheme []GradingSchemeEntry `json:"grading_scheme,omitempty"`
}

// GradingSchemeEntry is one letter-grade row in a grading standard.
type GradingSchemeEntry struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// GradingStandardParams holds mutable fields for creating/updating a grading standard.
type GradingStandardParams struct {
	Title         string               `json:"title"`
	GradingScheme []GradingSchemeEntry `json:"grading_scheme_entry,omitempty"`
}

// ListGradingStandards lists all grading standards available to a course.
// GET /api/v1/courses/:course_id/grading_standards
func (s *GradingStandardsService) ListForCourse(ctx context.Context, courseID int64) ([]GradingStandard, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/grading_standards", courseID)
	var out []GradingStandard
	if err := s.client.GetAllPages(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("listing grading standards for course %d: %w", courseID, err)
	}
	return out, nil
}

// GetGradingStandard retrieves a single grading standard from a course context.
// GET /api/v1/courses/:course_id/grading_standards/:grading_standard_id
func (s *GradingStandardsService) GetForCourse(ctx context.Context, courseID, standardID int64) (*GradingStandard, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/grading_standards/%d", courseID, standardID)
	var out GradingStandard
	if err := s.client.GetJSON(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("getting grading standard %d for course %d: %w", standardID, courseID, err)
	}
	return &out, nil
}

// CreateGradingStandardForCourse creates a new grading standard in a course context.
// POST /api/v1/courses/:course_id/grading_standards
func (s *GradingStandardsService) CreateForCourse(ctx context.Context, courseID int64, params GradingStandardParams) (*GradingStandard, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/grading_standards", courseID)
	var out GradingStandard
	if err := s.client.PostJSON(ctx, path, params, &out); err != nil {
		return nil, fmt.Errorf("creating grading standard for course %d: %w", courseID, err)
	}
	return &out, nil
}

// UpdateGradingStandardForCourse updates a grading standard in a course context.
// PUT /api/v1/courses/:course_id/grading_standards/:grading_standard_id
func (s *GradingStandardsService) UpdateForCourse(ctx context.Context, courseID, standardID int64, params GradingStandardParams) (*GradingStandard, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/grading_standards/%d", courseID, standardID)
	var out GradingStandard
	if err := s.client.PutJSON(ctx, path, params, &out); err != nil {
		return nil, fmt.Errorf("updating grading standard %d for course %d: %w", standardID, courseID, err)
	}
	return &out, nil
}

// DeleteGradingStandardForCourse deletes a grading standard from a course.
// DELETE /api/v1/courses/:course_id/grading_standards/:grading_standard_id
func (s *GradingStandardsService) DeleteForCourse(ctx context.Context, courseID, standardID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/grading_standards/%d", courseID, standardID)
	if _, err := s.client.Delete(ctx, path); err != nil {
		return fmt.Errorf("deleting grading standard %d for course %d: %w", standardID, courseID, err)
	}
	return nil
}
