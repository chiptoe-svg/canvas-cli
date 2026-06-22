package api

import (
	"context"
	"fmt"
	"io"
)

// AuthProvidersService handles authentication provider API calls
type AuthProvidersService struct {
	client *Client
}

// NewAuthProvidersService creates a new auth providers service
func NewAuthProvidersService(client *Client) *AuthProvidersService {
	return &AuthProvidersService{client: client}
}

// AuthenticationProvider represents a Canvas authentication provider
type AuthenticationProvider struct {
	ID              int64  `json:"id"`
	AuthType        string `json:"auth_type"`
	Position        int    `json:"position"`
	WorkflowState   string `json:"workflow_state"`
	JITProvisioning bool   `json:"jit_provisioning"`
}

// AuthProviderCreateParams holds parameters for creating an authentication provider
type AuthProviderCreateParams struct {
	AuthType            string                 `json:"auth_type"`
	ClientID            string                 `json:"client_id,omitempty"`
	ClientSecret        string                 `json:"client_secret,omitempty"`
	LoginAttribute      string                 `json:"login_attribute,omitempty"`
	FederatedAttributes map[string]interface{} `json:"federated_attributes,omitempty"`
}

// SSOSettings represents Canvas SSO settings
type SSOSettings struct {
	LoginHandleName   string `json:"login_handle_name"`
	ChangePasswordURL string `json:"change_password_url"`
	AuthDiscoveryURL  string `json:"auth_discovery_url"`
	UnknownUserURL    string `json:"unknown_user_url"`
}

// List retrieves all authentication providers for an account
func (s *AuthProvidersService) List(ctx context.Context, accountID int64) ([]AuthenticationProvider, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/authentication_providers", accountID)

	var providers []AuthenticationProvider
	if err := s.client.GetAllPages(ctx, path, &providers); err != nil {
		return nil, fmt.Errorf("list auth providers: %w", err)
	}

	return providers, nil
}

// Create creates a new authentication provider for an account
func (s *AuthProvidersService) Create(ctx context.Context, accountID int64, body *AuthProviderCreateParams) (*AuthenticationProvider, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/authentication_providers", accountID)

	var provider AuthenticationProvider
	if err := s.client.PostJSON(ctx, path, body, &provider); err != nil {
		return nil, fmt.Errorf("create auth provider: %w", err)
	}

	return &provider, nil
}

// Get retrieves a single authentication provider
func (s *AuthProvidersService) Get(ctx context.Context, accountID, id int64) (*AuthenticationProvider, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/authentication_providers/%d", accountID, id)

	var provider AuthenticationProvider
	if err := s.client.GetJSON(ctx, path, &provider); err != nil {
		return nil, fmt.Errorf("get auth provider: %w", err)
	}

	return &provider, nil
}

// Update updates an authentication provider
func (s *AuthProvidersService) Update(ctx context.Context, accountID, id int64, body *AuthProviderCreateParams) (*AuthenticationProvider, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/authentication_providers/%d", accountID, id)

	var provider AuthenticationProvider
	if err := s.client.PutJSON(ctx, path, body, &provider); err != nil {
		return nil, fmt.Errorf("update auth provider: %w", err)
	}

	return &provider, nil
}

// Delete deletes an authentication provider
func (s *AuthProvidersService) Delete(ctx context.Context, accountID, id int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/authentication_providers/%d", accountID, id)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("delete auth provider: %w", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}

// Restore restores a deleted authentication provider
func (s *AuthProvidersService) Restore(ctx context.Context, accountID, id int64) (*AuthenticationProvider, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/authentication_providers/%d/restore", accountID, id)

	var provider AuthenticationProvider
	if err := s.client.PutJSON(ctx, path, nil, &provider); err != nil {
		return nil, fmt.Errorf("restore auth provider: %w", err)
	}

	return &provider, nil
}

// ForcePasswordReset forces a password reset for all users via all authentication providers
func (s *AuthProvidersService) ForcePasswordReset(ctx context.Context, accountID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/authentication_providers/force_password_reset", accountID)

	var result interface{}
	if err := s.client.PostJSON(ctx, path, nil, &result); err != nil {
		return fmt.Errorf("force password reset: %w", err)
	}

	return nil
}

// GetSSOSettings retrieves SSO settings for an account
func (s *AuthProvidersService) GetSSOSettings(ctx context.Context, accountID int64) (*SSOSettings, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/sso_settings", accountID)

	var settings SSOSettings
	if err := s.client.GetJSON(ctx, path, &settings); err != nil {
		return nil, fmt.Errorf("get SSO settings: %w", err)
	}

	return &settings, nil
}

// UpdateSSOSettings updates SSO settings for an account
func (s *AuthProvidersService) UpdateSSOSettings(ctx context.Context, accountID int64, body *SSOSettings) (*SSOSettings, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/sso_settings", accountID)

	var settings SSOSettings
	if err := s.client.PutJSON(ctx, path, body, &settings); err != nil {
		return nil, fmt.Errorf("update SSO settings: %w", err)
	}

	return &settings, nil
}
