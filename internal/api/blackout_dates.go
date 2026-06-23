package api

import (
	"context"
	"fmt"
	"io"
)

// BlackoutDatesService handles blackout date API calls for courses.
type BlackoutDatesService struct {
	client *Client
}

// NewBlackoutDatesService creates a new blackout dates service.
func NewBlackoutDatesService(client *Client) *BlackoutDatesService {
	return &BlackoutDatesService{client: client}
}

// BlackoutDate represents a Canvas blackout date (a blocked period excluded from course pacing).
type BlackoutDate struct {
	ID          int64  `json:"id"`
	ContextID   int64  `json:"context_id"`
	ContextType string `json:"context_type"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	EventTitle  string `json:"event_title"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// BlackoutDateParams holds mutable fields for create/update operations.
type BlackoutDateParams struct {
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	EventTitle string `json:"event_title"`
}

// ListBlackoutDates lists all blackout dates for a course.
// GET /api/v1/courses/:course_id/blackout_dates
func (s *BlackoutDatesService) List(ctx context.Context, courseID int64) ([]BlackoutDate, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/blackout_dates", courseID)
	var out []BlackoutDate
	if err := s.client.GetAllPages(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("listing blackout dates for course %d: %w", courseID, err)
	}
	return out, nil
}

// GetBlackoutDate retrieves a single blackout date.
// GET /api/v1/courses/:course_id/blackout_dates/:id
func (s *BlackoutDatesService) Get(ctx context.Context, courseID, id int64) (*BlackoutDate, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/blackout_dates/%d", courseID, id)
	var out BlackoutDate
	if err := s.client.GetJSON(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("getting blackout date %d for course %d: %w", id, courseID, err)
	}
	return &out, nil
}

// CreateBlackoutDate creates a new blackout date for a course.
// POST /api/v1/courses/:course_id/blackout_dates
func (s *BlackoutDatesService) Create(ctx context.Context, courseID int64, params BlackoutDateParams) (*BlackoutDate, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/blackout_dates", courseID)
	body := map[string]interface{}{
		"blackout_date": params,
	}
	var out BlackoutDate
	if err := s.client.PostJSON(ctx, path, body, &out); err != nil {
		return nil, fmt.Errorf("creating blackout date for course %d: %w", courseID, err)
	}
	return &out, nil
}

// UpdateBlackoutDate updates a single blackout date.
// PUT /api/v1/courses/:course_id/blackout_dates/:id
func (s *BlackoutDatesService) Update(ctx context.Context, courseID, id int64, params BlackoutDateParams) (*BlackoutDate, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/blackout_dates/%d", courseID, id)
	body := map[string]interface{}{
		"blackout_date": params,
	}
	var out BlackoutDate
	if err := s.client.PutJSON(ctx, path, body, &out); err != nil {
		return nil, fmt.Errorf("updating blackout date %d for course %d: %w", id, courseID, err)
	}
	return &out, nil
}

// BulkUpdateBlackoutDates replaces all blackout dates for a course in one request.
// PUT /api/v1/courses/:course_id/blackout_dates
func (s *BlackoutDatesService) BulkUpdate(ctx context.Context, courseID int64, dates []BlackoutDateParams) ([]BlackoutDate, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/blackout_dates", courseID)
	body := map[string]interface{}{
		"blackout_dates": dates,
	}
	var out []BlackoutDate
	if err := s.client.PutJSON(ctx, path, body, &out); err != nil {
		return nil, fmt.Errorf("bulk-updating blackout dates for course %d: %w", courseID, err)
	}
	return out, nil
}

// DeleteBlackoutDate deletes a blackout date.
// DELETE /api/v1/courses/:course_id/blackout_dates/:id
func (s *BlackoutDatesService) Delete(ctx context.Context, courseID, id int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/blackout_dates/%d", courseID, id)
	if _, err := s.client.Delete(ctx, path); err != nil {
		return fmt.Errorf("deleting blackout date %d for course %d: %w", id, courseID, err)
	}
	return nil
}

type AccountBlackoutDatesService struct {
	client *Client
}

// NewAccountBlackoutDatesService creates a new account blackout dates service
func NewAccountBlackoutDatesService(client *Client) *AccountBlackoutDatesService {
	return &AccountBlackoutDatesService{client: client}
}

// BlackoutDate represents a Canvas account blackout date

// BlackoutDateParams holds parameters for creating or updating a blackout date

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
