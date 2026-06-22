package api

import (
	"context"
	"fmt"
)

// DeveloperKeysService handles developer key-related API calls.
type DeveloperKeysService struct {
	client *Client
}

// NewDeveloperKeysService creates a new DeveloperKeysService.
func NewDeveloperKeysService(client *Client) *DeveloperKeysService {
	return &DeveloperKeysService{client: client}
}

// DeveloperKey represents a Canvas developer key.
type DeveloperKey struct {
	ID                        int64  `json:"id"`
	APIKey                    string `json:"api_key"`
	CreatedAt                 string `json:"created_at"`
	LastUsedAt                string `json:"last_used_at"`
	UserID                    int64  `json:"user_id"`
	Name                      string `json:"name"`
	Email                     string `json:"email"`
	RedirectURI               string `json:"redirect_uri"`
	RedirectURIs              string `json:"redirect_uris"`
	VendorCode                string `json:"vendor_code"`
	IconURL                   string `json:"icon_url"`
	Notes                     string `json:"notes"`
	AccessTokenCount          int    `json:"access_token_count"`
	ClientCredentialsAudience string `json:"client_credentials_audience"`
}

// DeveloperKeyParams holds create parameters for a developer key.
type DeveloperKeyParams struct {
	Name        string `json:"developer_key[name],omitempty"`
	Email       string `json:"developer_key[email],omitempty"`
	RedirectURI string `json:"developer_key[redirect_uri],omitempty"`
	Notes       string `json:"developer_key[notes],omitempty"`
}

// DeveloperKeyBinding represents a developer key account binding.
type DeveloperKeyBinding struct {
	ID             int64  `json:"id"`
	AccountID      int64  `json:"account_id"`
	DeveloperKeyID int64  `json:"developer_key_id"`
	WorkflowState  string `json:"workflow_state"`
}

// DeveloperKeyBindingParams holds parameters for creating/updating a binding.
type DeveloperKeyBindingParams struct {
	WorkflowState string `json:"developer_key_account_binding[workflow_state]"`
}

// List retrieves all developer keys for an account.
func (s *DeveloperKeysService) List(ctx context.Context, accountID int64) ([]DeveloperKey, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/developer_keys", accountID)

	var keys []DeveloperKey
	if err := s.client.GetAllPages(ctx, path, &keys); err != nil {
		return nil, fmt.Errorf("list developer keys: %w", err)
	}

	return keys, nil
}

// Create creates a new developer key in an account.
func (s *DeveloperKeysService) Create(ctx context.Context, accountID int64, body *DeveloperKeyParams) (*DeveloperKey, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/developer_keys", accountID)

	var key DeveloperKey
	if err := s.client.PostJSON(ctx, path, body, &key); err != nil {
		return nil, fmt.Errorf("create developer key: %w", err)
	}

	return &key, nil
}

// CreateBinding creates a developer key account binding.
func (s *DeveloperKeysService) CreateBinding(ctx context.Context, accountID, developerKeyID int64, body *DeveloperKeyBindingParams) (*DeveloperKeyBinding, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/developer_keys/%d/developer_key_account_bindings", accountID, developerKeyID)

	var binding DeveloperKeyBinding
	if err := s.client.PostJSON(ctx, path, body, &binding); err != nil {
		return nil, fmt.Errorf("create developer key binding: %w", err)
	}

	return &binding, nil
}
