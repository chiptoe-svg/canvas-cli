package api

import (
	"context"
	"fmt"
)

// CourseFeaturesService handles feature flag API calls for courses.
type CourseFeaturesService struct {
	client *Client
}

// NewCourseFeaturesService creates a new course features service.
func NewCourseFeaturesService(client *Client) *CourseFeaturesService {
	return &CourseFeaturesService{client: client}
}

// Feature represents a Canvas feature definition.
type Feature struct {
	Feature     string      `json:"feature"`
	DisplayName string      `json:"display_name,omitempty"`
	AppliesTo   string      `json:"applies_to,omitempty"`
	FeatureFlag FeatureFlag `json:"feature_flag,omitempty"`
}

// FeatureFlag represents the flag state for a Canvas feature.
type FeatureFlag struct {
	Feature     string `json:"feature"`
	State       string `json:"state"` // "off", "allowed", "on"
	Locked      bool   `json:"locked,omitempty"`
	ContextID   int64  `json:"context_id,omitempty"`
	ContextType string `json:"context_type,omitempty"`
}

// ListCourseFeatures lists all available features for a course.
// GET /api/v1/courses/:course_id/features
func (s *CourseFeaturesService) List(ctx context.Context, courseID int64) ([]Feature, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/features", courseID)
	var out []Feature
	if err := s.client.GetAllPages(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("listing features for course %d: %w", courseID, err)
	}
	return out, nil
}

// ListEnabledCourseFeatures lists only enabled features for a course.
// GET /api/v1/courses/:course_id/features/enabled
func (s *CourseFeaturesService) ListEnabled(ctx context.Context, courseID int64) ([]string, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/features/enabled", courseID)
	var out []string
	if err := s.client.GetAllPages(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("listing enabled features for course %d: %w", courseID, err)
	}
	return out, nil
}

// GetCourseFeatureFlag retrieves a specific feature flag for a course.
// GET /api/v1/courses/:course_id/features/flags/:feature
func (s *CourseFeaturesService) GetFlag(ctx context.Context, courseID int64, feature string) (*FeatureFlag, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/features/flags/%s", courseID, feature)
	var out FeatureFlag
	if err := s.client.GetJSON(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("getting feature flag %q for course %d: %w", feature, courseID, err)
	}
	return &out, nil
}

// SetCourseFeatureFlag sets the state of a feature flag for a course.
// PUT /api/v1/courses/:course_id/features/flags/:feature
func (s *CourseFeaturesService) SetFlag(ctx context.Context, courseID int64, feature, state string) (*FeatureFlag, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/features/flags/%s", courseID, feature)
	body := map[string]string{"state": state}
	var out FeatureFlag
	if err := s.client.PutJSON(ctx, path, body, &out); err != nil {
		return nil, fmt.Errorf("setting feature flag %q for course %d: %w", feature, courseID, err)
	}
	return &out, nil
}

// DeleteCourseFeatureFlag removes a feature flag override for a course (reverts to default).
// DELETE /api/v1/courses/:course_id/features/flags/:feature
func (s *CourseFeaturesService) DeleteFlag(ctx context.Context, courseID int64, feature string) (*FeatureFlag, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/features/flags/%s", courseID, feature)
	var out FeatureFlag
	if err := s.client.DeleteJSON(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("deleting feature flag %q for course %d: %w", feature, courseID, err)
	}
	return &out, nil
}
