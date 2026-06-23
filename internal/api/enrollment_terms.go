package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// EnrollmentTermsService handles enrollment-term-related API calls.
type EnrollmentTermsService struct {
	client *Client
}

// NewEnrollmentTermsService creates a new EnrollmentTermsService.
func NewEnrollmentTermsService(client *Client) *EnrollmentTermsService {
	return &EnrollmentTermsService{client: client}
}

// EnrollmentTerm represents a Canvas enrollment term.
type EnrollmentTerm struct {
	ID                   int64  `json:"id"`
	SISTermID            string `json:"sis_term_id"`
	SISImportID          int64  `json:"sis_import_id"`
	Name                 string `json:"name"`
	StartAt              string `json:"start_at"`
	EndAt                string `json:"end_at"`
	WorkflowState        string `json:"workflow_state"`
	GradingPeriodGroupID int64  `json:"grading_period_group_id"`
}

// EnrollmentTermFields holds the nested fields for creating/updating an enrollment term.
type EnrollmentTermFields struct {
	Name      string `json:"name,omitempty"`
	StartAt   string `json:"start_at,omitempty"`
	EndAt     string `json:"end_at,omitempty"`
	SISTermID string `json:"sis_term_id,omitempty"`
}

// EnrollmentTermParams wraps the nested enrollment_term envelope Canvas expects.
// Canvas expects {"enrollment_term": {...}} rather than flat bracket-style JSON keys.
type EnrollmentTermParams struct {
	EnrollmentTerm EnrollmentTermFields `json:"enrollment_term"`
}

// EnrollmentTermsResponse is the wrapper Canvas uses for the list endpoint.
type EnrollmentTermsResponse struct {
	EnrollmentTerms []EnrollmentTerm `json:"enrollment_terms"`
}

// List retrieves all enrollment terms for an account.
func (s *EnrollmentTermsService) List(ctx context.Context, accountID int64) ([]EnrollmentTerm, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/terms", accountID)

	var wrapper EnrollmentTermsResponse
	if err := s.client.GetJSON(ctx, path, &wrapper); err != nil {
		return nil, fmt.Errorf("list enrollment terms: %w", err)
	}

	return wrapper.EnrollmentTerms, nil
}

// Get retrieves a single enrollment term.
func (s *EnrollmentTermsService) Get(ctx context.Context, accountID, termID int64) (*EnrollmentTerm, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/terms/%d", accountID, termID)

	var term EnrollmentTerm
	if err := s.client.GetJSON(ctx, path, &term); err != nil {
		return nil, fmt.Errorf("get enrollment term: %w", err)
	}

	return &term, nil
}

// Create creates a new enrollment term.
func (s *EnrollmentTermsService) Create(ctx context.Context, accountID int64, body *EnrollmentTermParams) (*EnrollmentTerm, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/terms", accountID)

	var term EnrollmentTerm
	if err := s.client.PostJSON(ctx, path, body, &term); err != nil {
		return nil, fmt.Errorf("create enrollment term: %w", err)
	}

	return &term, nil
}

// Update updates an existing enrollment term.
func (s *EnrollmentTermsService) Update(ctx context.Context, accountID, termID int64, body *EnrollmentTermParams) (*EnrollmentTerm, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/terms/%d", accountID, termID)

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal enrollment term params: %w", err)
	}

	resp, err := s.client.Put(ctx, path, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("update enrollment term: %w", err)
	}
	defer resp.Body.Close()

	var term EnrollmentTerm
	if err := json.NewDecoder(resp.Body).Decode(&term); err != nil {
		return nil, fmt.Errorf("decode enrollment term: %w", err)
	}

	return &term, nil
}

// Delete deletes an enrollment term.
func (s *EnrollmentTermsService) Delete(ctx context.Context, accountID, termID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/terms/%d", accountID, termID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("delete enrollment term: %w", err)
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}
