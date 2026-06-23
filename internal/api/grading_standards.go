package api

import (
	"context"
	"fmt"
	"io"
)

// GradingStandardsService handles grading standards API calls
type GradingStandardsService struct {
	client *Client
}

// NewGradingStandardsService creates a new grading standards service
func NewGradingStandardsService(client *Client) *GradingStandardsService {
	return &GradingStandardsService{client: client}
}

// GradingStandard represents a Canvas grading standard
type GradingStandard struct {
	ID            int64                `json:"id"`
	Title         string               `json:"title"`
	ContextType   string               `json:"context_type"`
	ContextID     int64                `json:"context_id"`
	GradingScheme []GradingSchemeEntry `json:"grading_scheme"`
}

// GradingSchemeEntry represents a single entry in a grading scheme
type GradingSchemeEntry struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// GradingStandardParams holds parameters for creating/updating grading standards
type GradingStandardParams struct {
	Title         string               `json:"title"`
	GradingScheme []GradingSchemeEntry `json:"grading_scheme_entry"`
}

// List retrieves grading standards for an account
func (s *GradingStandardsService) List(ctx context.Context, accountID int64) ([]GradingStandard, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/grading_standards", accountID)

	var standards []GradingStandard
	if err := s.client.GetAllPages(ctx, path, &standards); err != nil {
		return nil, fmt.Errorf("listing grading standards for account %d: %w", accountID, err)
	}

	return standards, nil
}

// Create creates a new grading standard for an account
func (s *GradingStandardsService) Create(ctx context.Context, accountID int64, body *GradingStandardParams) (*GradingStandard, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/grading_standards", accountID)

	var standard GradingStandard
	if err := s.client.PostJSON(ctx, path, body, &standard); err != nil {
		return nil, fmt.Errorf("creating grading standard for account %d: %w", accountID, err)
	}

	return &standard, nil
}

// Get retrieves a specific grading standard for an account
func (s *GradingStandardsService) Get(ctx context.Context, accountID, standardID int64) (*GradingStandard, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/grading_standards/%d", accountID, standardID)

	var standard GradingStandard
	if err := s.client.GetJSON(ctx, path, &standard); err != nil {
		return nil, fmt.Errorf("getting grading standard %d for account %d: %w", standardID, accountID, err)
	}

	return &standard, nil
}

// Update updates a grading standard for an account
func (s *GradingStandardsService) Update(ctx context.Context, accountID, standardID int64, body *GradingStandardParams) (*GradingStandard, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/grading_standards/%d", accountID, standardID)

	var standard GradingStandard
	if err := s.client.PutJSON(ctx, path, body, &standard); err != nil {
		return nil, fmt.Errorf("updating grading standard %d for account %d: %w", standardID, accountID, err)
	}

	return &standard, nil
}

// Delete deletes a grading standard for an account
func (s *GradingStandardsService) Delete(ctx context.Context, accountID, standardID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/grading_standards/%d", accountID, standardID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("deleting grading standard %d for account %d: %w", standardID, accountID, err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}

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
