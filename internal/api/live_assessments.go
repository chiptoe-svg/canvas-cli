package api

import (
	"context"
	"fmt"
)

// LiveAssessmentsService handles live assessment API calls for courses.
type LiveAssessmentsService struct {
	client *Client
}

// NewLiveAssessmentsService creates a new live assessments service.
func NewLiveAssessmentsService(client *Client) *LiveAssessmentsService {
	return &LiveAssessmentsService{client: client}
}

// LiveAssessment represents a Canvas live assessment.
type LiveAssessment struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Title string `json:"title"`
}

// LiveAssessmentResult represents a result for a live assessment.
type LiveAssessmentResult struct {
	Passed     bool   `json:"passed"`
	AssessedAt string `json:"assessed_at,omitempty"`
	Links      struct {
		User       string `json:"user,omitempty"`
		Assessor   string `json:"assessor,omitempty"`
		Assessment string `json:"assessment,omitempty"`
	} `json:"links,omitempty"`
}

// liveAssessmentsEnvelope wraps the Canvas live assessments response envelope.
type liveAssessmentsEnvelope struct {
	LiveAssessments []LiveAssessment `json:"live_assessments"`
}

// liveAssessmentResultsEnvelope wraps the Canvas live assessment results envelope.
type liveAssessmentResultsEnvelope struct {
	Results []LiveAssessmentResult `json:"results"`
}

// ListLiveAssessments lists live assessments for a course.
// GET /api/v1/courses/:course_id/live_assessments
func (s *LiveAssessmentsService) List(ctx context.Context, courseID int64) ([]LiveAssessment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/live_assessments", courseID)
	var envelope liveAssessmentsEnvelope
	if err := s.client.GetJSON(ctx, path, &envelope); err != nil {
		return nil, fmt.Errorf("listing live assessments for course %d: %w", courseID, err)
	}
	return envelope.LiveAssessments, nil
}

// CreateLiveAssessment creates a live assessment for a course.
// POST /api/v1/courses/:course_id/live_assessments
func (s *LiveAssessmentsService) Create(ctx context.Context, courseID int64, assessments []LiveAssessment) ([]LiveAssessment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/live_assessments", courseID)
	body := liveAssessmentsEnvelope{LiveAssessments: assessments}
	var envelope liveAssessmentsEnvelope
	if err := s.client.PostJSON(ctx, path, body, &envelope); err != nil {
		return nil, fmt.Errorf("creating live assessments for course %d: %w", courseID, err)
	}
	return envelope.LiveAssessments, nil
}

// ListResults lists results for a specific live assessment.
// GET /api/v1/courses/:course_id/live_assessments/:assessment_id/results
func (s *LiveAssessmentsService) ListResults(ctx context.Context, courseID int64, assessmentID string) ([]LiveAssessmentResult, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/live_assessments/%s/results", courseID, assessmentID)
	var envelope liveAssessmentResultsEnvelope
	if err := s.client.GetJSON(ctx, path, &envelope); err != nil {
		return nil, fmt.Errorf("listing results for live assessment %s in course %d: %w", assessmentID, courseID, err)
	}
	return envelope.Results, nil
}

// CreateResults posts results for a live assessment.
// POST /api/v1/courses/:course_id/live_assessments/:assessment_id/results
func (s *LiveAssessmentsService) CreateResults(ctx context.Context, courseID int64, assessmentID string, results []LiveAssessmentResult) ([]LiveAssessmentResult, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/live_assessments/%s/results", courseID, assessmentID)
	body := liveAssessmentResultsEnvelope{Results: results}
	var envelope liveAssessmentResultsEnvelope
	if err := s.client.PostJSON(ctx, path, body, &envelope); err != nil {
		return nil, fmt.Errorf("creating results for live assessment %s in course %d: %w", assessmentID, courseID, err)
	}
	return envelope.Results, nil
}
