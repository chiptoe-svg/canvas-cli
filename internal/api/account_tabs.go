package api

import (
	"context"
	"fmt"
)

// AccountTabsService handles account navigation tab API calls
type AccountTabsService struct {
	client *Client
}

// NewAccountTabsService creates a new account tabs service
func NewAccountTabsService(client *Client) *AccountTabsService {
	return &AccountTabsService{client: client}
}

// AccountTab represents a Canvas account navigation tab
type AccountTab struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Type     string `json:"type"`
	HTMLURL  string `json:"html_url"`
	Position int    `json:"position"`
	Hidden   bool   `json:"hidden"`
	Unused   bool   `json:"unused"`
}

// List retrieves navigation tabs for an account
func (s *AccountTabsService) List(ctx context.Context, accountID int64) ([]AccountTab, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/tabs", accountID)

	var tabs []AccountTab
	if err := s.client.GetAllPages(ctx, path, &tabs); err != nil {
		return nil, fmt.Errorf("listing account tabs for account %d: %w", accountID, err)
	}

	return tabs, nil
}
