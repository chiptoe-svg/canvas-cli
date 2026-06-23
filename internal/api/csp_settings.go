package api

import (
	"context"
	"fmt"
	"net/url"
)

// CSPSettingsService handles Content Security Policy settings API calls
type CSPSettingsService struct {
	client *Client
}

// NewCSPSettingsService creates a new CSP settings service
func NewCSPSettingsService(client *Client) *CSPSettingsService {
	return &CSPSettingsService{client: client}
}

// CSPSettings represents Canvas CSP settings for an account
type CSPSettings struct {
	Status    string   `json:"status"`
	Domains   []string `json:"domains"`
	Locked    bool     `json:"locked"`
	LockedBy  string   `json:"locked_by,omitempty"`
	Effective []string `json:"effective"`
}

// Get retrieves CSP settings for an account
func (s *CSPSettingsService) Get(ctx context.Context, accountID int64) (*CSPSettings, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/csp_settings", accountID)

	var settings CSPSettings
	if err := s.client.GetJSON(ctx, path, &settings); err != nil {
		return nil, fmt.Errorf("get CSP settings: %w", err)
	}

	return &settings, nil
}

// Update updates CSP settings for an account
func (s *CSPSettingsService) Update(ctx context.Context, accountID int64, body *CSPSettings) (*CSPSettings, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/csp_settings", accountID)

	var settings CSPSettings
	if err := s.client.PutJSON(ctx, path, body, &settings); err != nil {
		return nil, fmt.Errorf("update CSP settings: %w", err)
	}

	return &settings, nil
}

// RemoveDomains removes domains from the CSP allowlist.
// Canvas expects the domains as repeated "domains[]" query parameters on the DELETE request.
func (s *CSPSettingsService) RemoveDomains(ctx context.Context, accountID int64, domains []string) (*CSPSettings, error) {
	q := url.Values{}
	for _, d := range domains {
		q.Add("domains[]", d)
	}
	path := fmt.Sprintf("/api/v1/accounts/%d/csp_settings/domains?%s", accountID, q.Encode())

	var settings CSPSettings
	if err := s.client.DeleteJSON(ctx, path, &settings); err != nil {
		return nil, fmt.Errorf("remove CSP domains: %w", err)
	}

	return &settings, nil
}

// AddDomains adds domains to the CSP allowlist
func (s *CSPSettingsService) AddDomains(ctx context.Context, accountID int64, domains []string) (*CSPSettings, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/csp_settings/domains", accountID)

	body := map[string]interface{}{
		"domains": domains,
	}

	var settings CSPSettings
	if err := s.client.PostJSON(ctx, path, body, &settings); err != nil {
		return nil, fmt.Errorf("add CSP domains: %w", err)
	}

	return &settings, nil
}

// BatchAddDomains bulk adds domains to the CSP allowlist
func (s *CSPSettingsService) BatchAddDomains(ctx context.Context, accountID int64, domains []string) (*CSPSettings, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/csp_settings/domains/batch_create", accountID)

	body := map[string]interface{}{
		"domains": domains,
	}

	var settings CSPSettings
	if err := s.client.PostJSON(ctx, path, body, &settings); err != nil {
		return nil, fmt.Errorf("batch add CSP domains: %w", err)
	}

	return &settings, nil
}

// Lock locks or unlocks CSP settings for sub-accounts
func (s *CSPSettingsService) Lock(ctx context.Context, accountID int64, lock bool) (*CSPSettings, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/csp_settings/lock", accountID)

	body := map[string]interface{}{
		"locked": lock,
	}

	var settings CSPSettings
	if err := s.client.PutJSON(ctx, path, body, &settings); err != nil {
		return nil, fmt.Errorf("lock CSP settings: %w", err)
	}

	return &settings, nil
}
