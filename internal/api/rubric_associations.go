package api

import (
	"context"
	"fmt"
)

// RubricAssociationsService handles rubric association CRUD and assessments under courses.
// Note: the existing rubrics.go already defines the RubricAssociation type and
// associate/list operations on /courses/:id/rubrics/:rubric_id/associations.
// This service covers the *association-scoped* sub-endpoints that were missing:
// PUT/DELETE on /courses/:id/rubric_associations/:id and the rubric_assessments sub-resources.
type RubricAssociationsService struct {
	client *Client
}

// NewRubricAssociationsService creates a new rubric associations service.
func NewRubricAssociationsService(client *Client) *RubricAssociationsService {
	return &RubricAssociationsService{client: client}
}

// RubricAssessment represents a submitted rubric assessment.
type RubricAssessmentRecord struct {
	ID                  int64                       `json:"id"`
	RubricID            int64                       `json:"rubric_id"`
	RubricAssociationID int64                       `json:"rubric_association_id"`
	Score               float64                     `json:"score,omitempty"`
	ArtifactType        string                      `json:"artifact_type,omitempty"`
	ArtifactID          int64                       `json:"artifact_id,omitempty"`
	AssessorID          int64                       `json:"assessor_id,omitempty"`
	Data                []RubricAssessmentCriterion `json:"data,omitempty"`
}

// RubricAssessmentCriterion is one criterion row inside a rubric assessment.
type RubricAssessmentCriterion struct {
	CriterionID string  `json:"criterion_id"`
	Points      float64 `json:"points,omitempty"`
	Comments    string  `json:"comments,omitempty"`
	RatingID    string  `json:"rating_id,omitempty"`
}

// RubricAssociationUpdateParams holds editable fields for a rubric association.
type RubricAssociationUpdateParams struct {
	RubricID        int64  `json:"rubric_id,omitempty"`
	AssociationID   int64  `json:"association_id,omitempty"`
	AssociationType string `json:"association_type,omitempty"`
	UseForGrading   bool   `json:"use_for_grading,omitempty"`
	Purpose         string `json:"purpose,omitempty"`
}

// UpdateRubricAssociation updates a rubric association.
// PUT /api/v1/courses/:course_id/rubric_associations/:id
func (s *RubricAssociationsService) Update(ctx context.Context, courseID, associationID int64, params RubricAssociationUpdateParams) (*RubricAssociation, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubric_associations/%d", courseID, associationID)
	body := map[string]interface{}{
		"rubric_association": params,
	}
	var out RubricAssociation
	if err := s.client.PutJSON(ctx, path, body, &out); err != nil {
		return nil, fmt.Errorf("updating rubric association %d for course %d: %w", associationID, courseID, err)
	}
	return &out, nil
}

// DeleteRubricAssociation deletes a rubric association.
// DELETE /api/v1/courses/:course_id/rubric_associations/:id
func (s *RubricAssociationsService) Delete(ctx context.Context, courseID, associationID int64) (*RubricAssociation, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubric_associations/%d", courseID, associationID)
	var out RubricAssociation
	if err := s.client.DeleteJSON(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("deleting rubric association %d for course %d: %w", associationID, courseID, err)
	}
	return &out, nil
}

// CreateRubricAssessment posts a rubric assessment for an association.
// POST /api/v1/courses/:course_id/rubric_associations/:rubric_association_id/rubric_assessments
func (s *RubricAssociationsService) CreateAssessment(ctx context.Context, courseID, associationID int64, assessment RubricAssessmentRecord) (*RubricAssessmentRecord, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubric_associations/%d/rubric_assessments", courseID, associationID)
	body := map[string]interface{}{
		"rubric_assessment": assessment,
	}
	var out RubricAssessmentRecord
	if err := s.client.PostJSON(ctx, path, body, &out); err != nil {
		return nil, fmt.Errorf("creating rubric assessment for association %d in course %d: %w", associationID, courseID, err)
	}
	return &out, nil
}

// UpdateRubricAssessment updates a rubric assessment.
// PUT /api/v1/courses/:course_id/rubric_associations/:rubric_association_id/rubric_assessments/:id
func (s *RubricAssociationsService) UpdateAssessment(ctx context.Context, courseID, associationID, assessmentID int64, assessment RubricAssessmentRecord) (*RubricAssessmentRecord, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubric_associations/%d/rubric_assessments/%d", courseID, associationID, assessmentID)
	body := map[string]interface{}{
		"rubric_assessment": assessment,
	}
	var out RubricAssessmentRecord
	if err := s.client.PutJSON(ctx, path, body, &out); err != nil {
		return nil, fmt.Errorf("updating rubric assessment %d for association %d in course %d: %w", assessmentID, associationID, courseID, err)
	}
	return &out, nil
}

// DeleteRubricAssessment deletes a rubric assessment.
// DELETE /api/v1/courses/:course_id/rubric_associations/:rubric_association_id/rubric_assessments/:id
func (s *RubricAssociationsService) DeleteAssessment(ctx context.Context, courseID, associationID, assessmentID int64) (*RubricAssessmentRecord, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubric_associations/%d/rubric_assessments/%d", courseID, associationID, assessmentID)
	var out RubricAssessmentRecord
	if err := s.client.DeleteJSON(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("deleting rubric assessment %d for association %d in course %d: %w", assessmentID, associationID, courseID, err)
	}
	return &out, nil
}
