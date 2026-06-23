package api

import (
	"context"
	"fmt"
)

// ObserveesService handles observer/observee relationship API calls
type ObserveesService struct {
	client *Client
}

// NewObserveesService creates a new observees service
func NewObserveesService(client *Client) *ObserveesService {
	return &ObserveesService{client: client}
}

// ObserverPairingCode represents a pairing code used to link an observer to a student
type ObserverPairingCode struct {
	UserID        int64  `json:"user_id"`
	Code          string `json:"code"`
	ExpiresAt     string `json:"expires_at"`
	WorkflowState string `json:"workflow_state"`
}

// AddObserveeParams holds parameters for adding an observee
type AddObserveeParams struct {
	ObserveeID    int64  `json:"observee_id,omitempty"`
	AccessToken   string `json:"access_token,omitempty"`
	PairingCode   string `json:"pairing_code,omitempty"`
	RootAccountID int64  `json:"root_account_id,omitempty"`
}

// ListObservees retrieves the students being observed by a user
func (s *ObserveesService) ListObservees(ctx context.Context, userID int64) ([]User, error) {
	path := fmt.Sprintf("/api/v1/users/%d/observees", userID)
	var result []User
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing observees for user %d: %w", userID, err)
	}
	return result, nil
}

// AddObservee links a student to be observed by the given user
func (s *ObserveesService) AddObservee(ctx context.Context, userID int64, params AddObserveeParams) (*User, error) {
	path := fmt.Sprintf("/api/v1/users/%d/observees", userID)
	var result User
	if err := s.client.PostJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("adding observee for user %d: %w", userID, err)
	}
	return &result, nil
}

// GetObservee retrieves a specific observee for a user
func (s *ObserveesService) GetObservee(ctx context.Context, userID, observeeID int64) (*User, error) {
	path := fmt.Sprintf("/api/v1/users/%d/observees/%d", userID, observeeID)
	var result User
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting observee %d for user %d: %w", observeeID, userID, err)
	}
	return &result, nil
}

// UpdateObservee updates the observee relationship (e.g. re-links via root account)
func (s *ObserveesService) UpdateObservee(ctx context.Context, userID, observeeID int64) (*User, error) {
	path := fmt.Sprintf("/api/v1/users/%d/observees/%d", userID, observeeID)
	var result User
	if err := s.client.PutJSON(ctx, path, nil, &result); err != nil {
		return nil, fmt.Errorf("updating observee %d for user %d: %w", observeeID, userID, err)
	}
	return &result, nil
}

// RemoveObservee unlinks an observee from a user
func (s *ObserveesService) RemoveObservee(ctx context.Context, userID, observeeID int64) (*User, error) {
	path := fmt.Sprintf("/api/v1/users/%d/observees/%d", userID, observeeID)
	var result User
	if err := s.client.DeleteJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("removing observee %d from user %d: %w", observeeID, userID, err)
	}
	return &result, nil
}

// CreatePairingCode generates a new observer pairing code for a user
func (s *ObserveesService) CreatePairingCode(ctx context.Context, userID int64) (*ObserverPairingCode, error) {
	path := fmt.Sprintf("/api/v1/users/%d/observer_pairing_codes", userID)
	var result ObserverPairingCode
	if err := s.client.PostJSON(ctx, path, nil, &result); err != nil {
		return nil, fmt.Errorf("creating pairing code for user %d: %w", userID, err)
	}
	return &result, nil
}

// ListObservers retrieves the observers watching a given user
func (s *ObserveesService) ListObservers(ctx context.Context, userID int64) ([]User, error) {
	path := fmt.Sprintf("/api/v1/users/%d/observers", userID)
	var result []User
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing observers for user %d: %w", userID, err)
	}
	return result, nil
}

// GetObserver retrieves a specific observer for a user
func (s *ObserveesService) GetObserver(ctx context.Context, userID, observerID int64) (*User, error) {
	path := fmt.Sprintf("/api/v1/users/%d/observers/%d", userID, observerID)
	var result User
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting observer %d for user %d: %w", observerID, userID, err)
	}
	return &result, nil
}
