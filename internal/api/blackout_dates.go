package api

import (
	"context"
	"fmt"
	"io"
)

// AccountBlackoutDatesService handles blackout date API calls for accounts
type AccountBlackoutDatesService struct {
	client *Client
}

// NewAccountBlackoutDatesService creates a new account blackout dates service
func NewAccountBlackoutDatesService(client *Client) *AccountBlackoutDatesService {
	return &AccountBlackoutDatesService{client: client}
}

// BlackoutDate represents a Canvas account blackout date
type BlackoutDate struct {
	ID          int64  `json:"id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	EventTitle  string `json:"event_title"`
	ContextType string `json:"context_type"`
	ContextID   int64  `json:"context_id"`
}

// BlackoutDateParams holds parameters for creating or updating a blackout date
type BlackoutDateParams struct {
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	EventTitle string `json:"event_title"`
}

// List retrieves all blackout dates for an account
func (s *AccountBlackoutDatesService) List(ctx context.Context, accountID int64) ([]BlackoutDate, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/blackout_dates", accountID)

	var dates []BlackoutDate
	if err := s.client.GetAllPages(ctx, path, &dates); err != nil {
		return nil, fmt.Errorf("failed to list blackout dates: %w", err)
	}

	return dates, nil
}

// Get retrieves a single blackout date
func (s *AccountBlackoutDatesService) Get(ctx context.Context, accountID, id int64) (*BlackoutDate, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/blackout_dates/%d", accountID, id)

	var date BlackoutDate
	if err := s.client.GetJSON(ctx, path, &date); err != nil {
		return nil, fmt.Errorf("failed to get blackout date: %w", err)
	}

	return &date, nil
}

// GetNew returns a new (unsaved) blackout date template for an account
func (s *AccountBlackoutDatesService) GetNew(ctx context.Context, accountID int64) (*BlackoutDate, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/blackout_dates/new", accountID)

	var date BlackoutDate
	if err := s.client.GetJSON(ctx, path, &date); err != nil {
		return nil, fmt.Errorf("failed to get new blackout date: %w", err)
	}

	return &date, nil
}

// Create creates a new blackout date for an account
func (s *AccountBlackoutDatesService) Create(ctx context.Context, accountID int64, body *BlackoutDateParams) (*BlackoutDate, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/blackout_dates", accountID)

	var date BlackoutDate
	if err := s.client.PostJSON(ctx, path, body, &date); err != nil {
		return nil, fmt.Errorf("failed to create blackout date: %w", err)
	}

	return &date, nil
}

// Update updates an existing blackout date
func (s *AccountBlackoutDatesService) Update(ctx context.Context, accountID, id int64, body *BlackoutDateParams) (*BlackoutDate, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/blackout_dates/%d", accountID, id)

	var date BlackoutDate
	if err := s.client.PutJSON(ctx, path, body, &date); err != nil {
		return nil, fmt.Errorf("failed to update blackout date: %w", err)
	}

	return &date, nil
}

// Delete deletes a blackout date
func (s *AccountBlackoutDatesService) Delete(ctx context.Context, accountID, id int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/blackout_dates/%d", accountID, id)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete blackout date: %w", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}
