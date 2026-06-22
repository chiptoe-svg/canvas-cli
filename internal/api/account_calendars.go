package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// AccountCalendar represents a Canvas account-level calendar (for visibility settings).
type AccountCalendar struct {
	ID              int64  `json:"id"`
	Name            string `json:"name,omitempty"`
	ParentAccountID *int64 `json:"parent_account_id,omitempty"`
	RootAccountID   int64  `json:"root_account_id,omitempty"`
	Visible         bool   `json:"visible"`
	AutoSubscribe   bool   `json:"auto_subscribe"`
	Type            string `json:"type,omitempty"`
	SubAccountCount int    `json:"sub_account_count,omitempty"`
}

// AccountCalendarsService handles account-calendar API calls.
type AccountCalendarsService struct {
	client *Client
}

// NewAccountCalendarsService creates a new AccountCalendarsService.
func NewAccountCalendarsService(client *Client) *AccountCalendarsService {
	return &AccountCalendarsService{client: client}
}

// ListAccountCalendarsOptions holds query parameters.
type ListAccountCalendarsOptions struct {
	Filter  string // search term
	PerPage int
}

func buildAccountCalendarQuery(opts *ListAccountCalendarsOptions) string {
	if opts == nil {
		return ""
	}
	q := url.Values{}
	if opts.Filter != "" {
		q.Set("filter", opts.Filter)
	}
	if opts.PerPage > 0 {
		q.Set("per_page", strconv.Itoa(opts.PerPage))
	}
	if len(q) > 0 {
		return "?" + q.Encode()
	}
	return ""
}

// List retrieves all visible account calendars for the current user.
func (s *AccountCalendarsService) List(ctx context.Context, opts *ListAccountCalendarsOptions) ([]AccountCalendar, error) {
	path := "/api/v1/account_calendars" + buildAccountCalendarQuery(opts)

	var cals []AccountCalendar
	if err := s.client.GetAllPages(ctx, path, &cals); err != nil {
		return nil, err
	}

	return cals, nil
}

// Get retrieves a single account calendar by account ID.
func (s *AccountCalendarsService) Get(ctx context.Context, accountID int64) (*AccountCalendar, error) {
	path := fmt.Sprintf("/api/v1/account_calendars/%d", accountID)

	var cal AccountCalendar
	if err := s.client.GetJSON(ctx, path, &cal); err != nil {
		return nil, err
	}

	return &cal, nil
}

// UpdateAccountCalendarParams holds update parameters for an account calendar.
type UpdateAccountCalendarParams struct {
	Visible       *bool
	AutoSubscribe *bool
}

// Update updates the visibility / auto-subscribe setting of an account calendar.
func (s *AccountCalendarsService) Update(ctx context.Context, accountID int64, params *UpdateAccountCalendarParams) (*AccountCalendar, error) {
	path := fmt.Sprintf("/api/v1/account_calendars/%d", accountID)

	body := map[string]interface{}{}
	if params.Visible != nil {
		body["visible"] = *params.Visible
	}
	if params.AutoSubscribe != nil {
		body["auto_subscribe"] = *params.AutoSubscribe
	}

	var cal AccountCalendar
	if err := s.client.PutJSON(ctx, path, body, &cal); err != nil {
		return nil, err
	}

	return &cal, nil
}

// ListForAccount retrieves all account calendars under a given account.
func (s *AccountCalendarsService) ListForAccount(ctx context.Context, accountID int64, opts *ListAccountCalendarsOptions) ([]AccountCalendar, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/account_calendars", accountID) + buildAccountCalendarQuery(opts)

	var cals []AccountCalendar
	if err := s.client.GetAllPages(ctx, path, &cals); err != nil {
		return nil, err
	}

	return cals, nil
}
