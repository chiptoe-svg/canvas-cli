package api

import (
	"context"
	"fmt"
	"net/url"
)

// AccountCalendarsService handles account calendar API calls
type AccountCalendarsService struct {
	client *Client
}

// NewAccountCalendarsService creates a new account calendars service
func NewAccountCalendarsService(client *Client) *AccountCalendarsService {
	return &AccountCalendarsService{client: client}
}

// AccountCalendar represents a Canvas account calendar
type AccountCalendar struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	ParentAccountID int64  `json:"parent_account_id"`
	RootAccountID   int64  `json:"root_account_id"`
	Visible         bool   `json:"visible"`
	AutoSubscribe   bool   `json:"auto_subscribe"`
	SubContextType  string `json:"sub_context_type"`
	SubContextID    int64  `json:"sub_context_id"`
}

// AccountCalendarParams holds parameters for creating/updating an account calendar
type AccountCalendarParams struct {
	Visible       *bool `json:"visible,omitempty"`
	AutoSubscribe *bool `json:"auto_subscribe,omitempty"`
}

// accountCalendarsResponse wraps the Canvas API envelope for account calendar list endpoints.
// Canvas returns {"account_calendars":[...],"total_results":N} rather than a bare array.
type accountCalendarsResponse struct {
	AccountCalendars []AccountCalendar `json:"account_calendars"`
	TotalResults     int               `json:"total_results"`
}

// ListAll retrieves all account calendars visible to the current user
func (s *AccountCalendarsService) ListAll(ctx context.Context, searchTerm string) ([]AccountCalendar, error) {
	path := "/api/v1/account_calendars"

	if searchTerm != "" {
		query := url.Values{}
		query.Set("search_term", searchTerm)
		path += "?" + query.Encode()
	}

	var resp accountCalendarsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("listing account calendars: %w", err)
	}

	return resp.AccountCalendars, nil
}

// Get retrieves a single account calendar by account ID
func (s *AccountCalendarsService) Get(ctx context.Context, accountID int64) (*AccountCalendar, error) {
	path := fmt.Sprintf("/api/v1/account_calendars/%d", accountID)

	var calendar AccountCalendar
	if err := s.client.GetJSON(ctx, path, &calendar); err != nil {
		return nil, fmt.Errorf("getting account calendar %d: %w", accountID, err)
	}

	return &calendar, nil
}

// Update updates an account calendar
func (s *AccountCalendarsService) Update(ctx context.Context, accountID int64, body *AccountCalendarParams) (*AccountCalendar, error) {
	path := fmt.Sprintf("/api/v1/account_calendars/%d", accountID)

	var calendar AccountCalendar
	if err := s.client.PutJSON(ctx, path, body, &calendar); err != nil {
		return nil, fmt.Errorf("updating account calendar %d: %w", accountID, err)
	}

	return &calendar, nil
}

// ListForAccount retrieves account calendars for a specific account
func (s *AccountCalendarsService) ListForAccount(ctx context.Context, accountID int64, searchTerm string) ([]AccountCalendar, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/account_calendars", accountID)

	if searchTerm != "" {
		query := url.Values{}
		query.Set("search_term", searchTerm)
		path += "?" + query.Encode()
	}

	var resp accountCalendarsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("listing calendars for account %d: %w", accountID, err)
	}

	return resp.AccountCalendars, nil
}

// UpdateForAccount updates account calendars for a specific account
func (s *AccountCalendarsService) UpdateForAccount(ctx context.Context, accountID int64, body *AccountCalendarParams) (*AccountCalendar, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/account_calendars", accountID)

	var calendar AccountCalendar
	if err := s.client.PutJSON(ctx, path, body, &calendar); err != nil {
		return nil, fmt.Errorf("updating calendars for account %d: %w", accountID, err)
	}

	return &calendar, nil
}

// GetVisibleCount retrieves the count of visible calendars for an account
func (s *AccountCalendarsService) GetVisibleCount(ctx context.Context, accountID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/visible_calendars_count", accountID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting visible calendars count for account %d: %w", accountID, err)
	}

	return result, nil
}
