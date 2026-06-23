package api

import (
	"context"
	"fmt"
	"io"
)

// AccountTemporaryEnrollmentPairingsService handles temporary enrollment pairing API calls
type AccountTemporaryEnrollmentPairingsService struct {
	client *Client
}

// NewAccountTemporaryEnrollmentPairingsService creates a new temporary enrollment pairings service
func NewAccountTemporaryEnrollmentPairingsService(client *Client) *AccountTemporaryEnrollmentPairingsService {
	return &AccountTemporaryEnrollmentPairingsService{client: client}
}

// TemporaryEnrollmentPairing represents a Canvas temporary enrollment pairing
type TemporaryEnrollmentPairing struct {
	ID                      int64  `json:"id"`
	RootAccountID           int64  `json:"root_account_id"`
	WorkflowState           string `json:"workflow_state"`
	CreatedAt               string `json:"created_at"`
	UpdatedAt               string `json:"updated_at"`
	StartingEnrollmentState string `json:"starting_enrollment_state"`
}

// TemporaryEnrollmentPairingParams holds parameters for creating a temporary enrollment pairing
type TemporaryEnrollmentPairingParams struct {
	StartingEnrollmentState string `json:"starting_enrollment_state,omitempty"`
	RoleID                  int64  `json:"role_id,omitempty"`
}

// List retrieves all temporary enrollment pairings for an account
func (s *AccountTemporaryEnrollmentPairingsService) List(ctx context.Context, accountID int64) ([]TemporaryEnrollmentPairing, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/temporary_enrollment_pairings", accountID)

	var pairings []TemporaryEnrollmentPairing
	if err := s.client.GetAllPages(ctx, path, &pairings); err != nil {
		return nil, fmt.Errorf("failed to list temporary enrollment pairings: %w", err)
	}

	return pairings, nil
}

// Get retrieves a single temporary enrollment pairing
func (s *AccountTemporaryEnrollmentPairingsService) Get(ctx context.Context, accountID, id int64) (*TemporaryEnrollmentPairing, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/temporary_enrollment_pairings/%d", accountID, id)

	var pairing TemporaryEnrollmentPairing
	if err := s.client.GetJSON(ctx, path, &pairing); err != nil {
		return nil, fmt.Errorf("failed to get temporary enrollment pairing: %w", err)
	}

	return &pairing, nil
}

// GetNew returns a new (unsaved) temporary enrollment pairing template
func (s *AccountTemporaryEnrollmentPairingsService) GetNew(ctx context.Context, accountID int64) (*TemporaryEnrollmentPairing, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/temporary_enrollment_pairings/new", accountID)

	var pairing TemporaryEnrollmentPairing
	if err := s.client.GetJSON(ctx, path, &pairing); err != nil {
		return nil, fmt.Errorf("failed to get new temporary enrollment pairing: %w", err)
	}

	return &pairing, nil
}

// Create creates a new temporary enrollment pairing for an account
func (s *AccountTemporaryEnrollmentPairingsService) Create(ctx context.Context, accountID int64, body *TemporaryEnrollmentPairingParams) (*TemporaryEnrollmentPairing, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/temporary_enrollment_pairings", accountID)

	var pairing TemporaryEnrollmentPairing
	if err := s.client.PostJSON(ctx, path, body, &pairing); err != nil {
		return nil, fmt.Errorf("failed to create temporary enrollment pairing: %w", err)
	}

	return &pairing, nil
}

// Delete deletes a temporary enrollment pairing
func (s *AccountTemporaryEnrollmentPairingsService) Delete(ctx context.Context, accountID, id int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/temporary_enrollment_pairings/%d", accountID, id)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete temporary enrollment pairing: %w", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}
