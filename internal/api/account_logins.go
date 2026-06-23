package api

import (
	"context"
	"fmt"
)

// AccountLoginsService handles account login-related API calls.
type AccountLoginsService struct {
	client *Client
}

// NewAccountLoginsService creates a new AccountLoginsService.
func NewAccountLoginsService(client *Client) *AccountLoginsService {
	return &AccountLoginsService{client: client}
}

// Login represents a Canvas user login.
type Login struct {
	ID                         int64  `json:"id"`
	UserID                     int64  `json:"user_id"`
	UniqueID                   string `json:"unique_id"`
	SISUserID                  string `json:"sis_user_id"`
	AccountID                  int64  `json:"account_id"`
	WorkflowState              string `json:"workflow_state"`
	AuthenticationProviderType string `json:"authentication_provider_type"`
	AuthenticationProviderID   int64  `json:"authentication_provider_id"`
}

// UserIDFields holds the user id nested under "user" for login creation.
type UserIDFields struct {
	ID int64 `json:"id,omitempty"`
}

// LoginFields holds the nested login fields Canvas expects.
type LoginFields struct {
	UniqueID                   string `json:"unique_id,omitempty"`
	Password                   string `json:"password,omitempty"`
	SISUserID                  string `json:"sis_user_id,omitempty"`
	AuthenticationProviderType string `json:"authentication_provider_type,omitempty"`
	AuthenticationProviderID   int64  `json:"authentication_provider_id,omitempty"`
}

// LoginParams wraps the nested envelopes Canvas expects.
// Canvas expects {"user": {"id": ...}, "login": {...}} rather than flat bracket-style JSON keys.
type LoginParams struct {
	User  UserIDFields `json:"user,omitempty"`
	Login LoginFields  `json:"login"`
}

// List retrieves logins for an account. When userID is non-zero only that
// user's logins are returned; zero means all logins in the account.
func (s *AccountLoginsService) List(ctx context.Context, accountID int64, userID int64) ([]Login, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/logins", accountID)
	if userID != 0 {
		path = fmt.Sprintf("%s?user_id=%d", path, userID)
	}

	var logins []Login
	if err := s.client.GetAllPages(ctx, path, &logins); err != nil {
		return nil, fmt.Errorf("list account logins: %w", err)
	}

	return logins, nil
}

// Create creates a new login for a user in an account.
func (s *AccountLoginsService) Create(ctx context.Context, accountID int64, body *LoginParams) (*Login, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/logins", accountID)

	var login Login
	if err := s.client.PostJSON(ctx, path, body, &login); err != nil {
		return nil, fmt.Errorf("create account login: %w", err)
	}

	return &login, nil
}

// Update updates an existing login.
func (s *AccountLoginsService) Update(ctx context.Context, accountID, loginID int64, body *LoginParams) (*Login, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/logins/%d", accountID, loginID)

	var login Login
	if err := s.client.PutJSON(ctx, path, body, &login); err != nil {
		return nil, fmt.Errorf("update account login: %w", err)
	}

	return &login, nil
}
