package api

import (
	"context"
	"fmt"
	"io"
)

// AccountLTIRegistrationsService handles LTI 1.3 registration API calls for accounts
type AccountLTIRegistrationsService struct {
	client *Client
}

// NewAccountLTIRegistrationsService creates a new LTI registrations service
func NewAccountLTIRegistrationsService(client *Client) *AccountLTIRegistrationsService {
	return &AccountLTIRegistrationsService{client: client}
}

// LTIRegistration represents a Canvas LTI 1.3 registration
type LTIRegistration struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	AdminNickname string `json:"admin_nickname,omitempty"`
	ClientID      string `json:"client_id"`
	WorkflowState string `json:"workflow_state"`
	AccountID     int64  `json:"account_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// List retrieves all LTI registrations for an account
func (s *AccountLTIRegistrationsService) List(ctx context.Context, accountID int64) ([]LTIRegistration, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations", accountID)

	var registrations []LTIRegistration
	if err := s.client.GetAllPages(ctx, path, &registrations); err != nil {
		return nil, fmt.Errorf("failed to list LTI registrations: %w", err)
	}

	return registrations, nil
}

// Get retrieves a single LTI registration
func (s *AccountLTIRegistrationsService) Get(ctx context.Context, accountID, id int64) (*LTIRegistration, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d", accountID, id)

	var registration LTIRegistration
	if err := s.client.GetJSON(ctx, path, &registration); err != nil {
		return nil, fmt.Errorf("failed to get LTI registration: %w", err)
	}

	return &registration, nil
}

// Create creates a new LTI registration for an account
func (s *AccountLTIRegistrationsService) Create(ctx context.Context, accountID int64, body interface{}) (*LTIRegistration, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations", accountID)

	var registration LTIRegistration
	if err := s.client.PostJSON(ctx, path, body, &registration); err != nil {
		return nil, fmt.Errorf("failed to create LTI registration: %w", err)
	}

	return &registration, nil
}

// Update updates an existing LTI registration
func (s *AccountLTIRegistrationsService) Update(ctx context.Context, accountID, id int64, body interface{}) (*LTIRegistration, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d", accountID, id)

	var registration LTIRegistration
	if err := s.client.PutJSON(ctx, path, body, &registration); err != nil {
		return nil, fmt.Errorf("failed to update LTI registration: %w", err)
	}

	return &registration, nil
}

// Delete deletes an LTI registration
func (s *AccountLTIRegistrationsService) Delete(ctx context.Context, accountID, id int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d", accountID, id)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete LTI registration: %w", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}

// Bind binds an LTI registration to an account
func (s *AccountLTIRegistrationsService) Bind(ctx context.Context, accountID, id int64, body interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/bind", accountID, id)

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to bind LTI registration: %w", err)
	}

	return result, nil
}

// Unbind removes the binding of an LTI registration from an account
func (s *AccountLTIRegistrationsService) Unbind(ctx context.Context, accountID, id int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/bind", accountID, id)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to unbind LTI registration: %w", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}

// GetHistory retrieves the history of an LTI registration
func (s *AccountLTIRegistrationsService) GetHistory(ctx context.Context, accountID, id int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/history", accountID, id)

	var result []interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get LTI registration history: %w", err)
	}

	return result, nil
}

// Reset resets an LTI registration to its default state
func (s *AccountLTIRegistrationsService) Reset(ctx context.Context, accountID, id int64) (*LTIRegistration, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/reset", accountID, id)

	var registration LTIRegistration
	if err := s.client.PutJSON(ctx, path, nil, &registration); err != nil {
		return nil, fmt.Errorf("failed to reset LTI registration: %w", err)
	}

	return &registration, nil
}

// GetByClientID retrieves an LTI registration by its client ID
func (s *AccountLTIRegistrationsService) GetByClientID(ctx context.Context, accountID int64, clientID string) (*LTIRegistration, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registration_by_client_id/%s", accountID, clientID)

	var registration LTIRegistration
	if err := s.client.GetJSON(ctx, path, &registration); err != nil {
		return nil, fmt.Errorf("failed to get LTI registration by client ID: %w", err)
	}

	return &registration, nil
}

// GetLaunchDefinitions retrieves LTI app launch definitions for an account
func (s *AccountLTIRegistrationsService) GetLaunchDefinitions(ctx context.Context, accountID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_apps/launch_definitions", accountID)

	var result []interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get LTI launch definitions: %w", err)
	}

	return result, nil
}

// InstallFromTemplate installs an LTI registration from a template
func (s *AccountLTIRegistrationsService) InstallFromTemplate(ctx context.Context, accountID, id int64, body interface{}) (*LTIRegistration, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/install_from_template", accountID, id)

	var result LTIRegistration
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to install LTI registration from template: %w", err)
	}

	return &result, nil
}

// GetLatestUpdateRequest retrieves the latest update request for an LTI registration
func (s *AccountLTIRegistrationsService) GetLatestUpdateRequest(ctx context.Context, accountID, id int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/latest_update_request", accountID, id)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get latest update request: %w", err)
	}

	return result, nil
}

// GetOverlayHistory retrieves the overlay history of an LTI registration
func (s *AccountLTIRegistrationsService) GetOverlayHistory(ctx context.Context, accountID, id int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/overlay_history", accountID, id)

	var result []interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get overlay history: %w", err)
	}

	return result, nil
}

// GetUpdateRequest retrieves a specific update request for an LTI registration
func (s *AccountLTIRegistrationsService) GetUpdateRequest(ctx context.Context, accountID, id, updateRequestID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/update_requests/%d", accountID, id, updateRequestID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get update request: %w", err)
	}

	return result, nil
}

// ApplyUpdateRequest applies a specific update request for an LTI registration
func (s *AccountLTIRegistrationsService) ApplyUpdateRequest(ctx context.Context, accountID, id, updateRequestID int64) (*LTIRegistration, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/update_requests/%d/apply", accountID, id, updateRequestID)

	var result LTIRegistration
	if err := s.client.PutJSON(ctx, path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to apply update request: %w", err)
	}

	return &result, nil
}

// ListControls retrieves controls for an LTI registration
func (s *AccountLTIRegistrationsService) ListControls(ctx context.Context, accountID, registrationID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/controls", accountID, registrationID)

	var result []interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to list controls: %w", err)
	}

	return result, nil
}

// CreateControl creates a control for an LTI registration
func (s *AccountLTIRegistrationsService) CreateControl(ctx context.Context, accountID, registrationID int64, body interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/controls", accountID, registrationID)

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to create control: %w", err)
	}

	return result, nil
}

// GetControl retrieves a specific control for an LTI registration
func (s *AccountLTIRegistrationsService) GetControl(ctx context.Context, accountID, registrationID, controlID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/controls/%d", accountID, registrationID, controlID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get control: %w", err)
	}

	return result, nil
}

// UpdateControl updates a specific control for an LTI registration
func (s *AccountLTIRegistrationsService) UpdateControl(ctx context.Context, accountID, registrationID, controlID int64, body interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/controls/%d", accountID, registrationID, controlID)

	var result map[string]interface{}
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to update control: %w", err)
	}

	return result, nil
}

// DeleteControl deletes a specific control for an LTI registration
func (s *AccountLTIRegistrationsService) DeleteControl(ctx context.Context, accountID, registrationID, controlID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/controls/%d", accountID, registrationID, controlID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete control: %w", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}

// BulkCreateControls bulk-creates controls for an LTI registration
func (s *AccountLTIRegistrationsService) BulkCreateControls(ctx context.Context, accountID, registrationID int64, body interface{}) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/controls/bulk", accountID, registrationID)

	var result []interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to bulk create controls: %w", err)
	}

	return result, nil
}

// GetDeploymentContextSearch searches deployment contexts for an LTI registration
func (s *AccountLTIRegistrationsService) GetDeploymentContextSearch(ctx context.Context, accountID, registrationID, deploymentID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/deployments/%d/context_search", accountID, registrationID, deploymentID)

	var result []interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to search deployment contexts: %w", err)
	}

	return result, nil
}

// GetByUTID retrieves an LTI registration by its UTID
func (s *AccountLTIRegistrationsService) GetByUTID(ctx context.Context, accountID int64, utid string) (*LTIRegistration, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/by_utid/%s", accountID, utid)

	var result LTIRegistration
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get LTI registration by UTID: %w", err)
	}

	return &result, nil
}

// GetInstallStatus retrieves the install status of an LTI registration by client ID
func (s *AccountLTIRegistrationsService) GetInstallStatus(ctx context.Context, accountID int64, clientID string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/install_status/%s", accountID, clientID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get install status: %w", err)
	}

	return result, nil
}

// GetLTIAccount retrieves an account via the LTI API path
// This is used by LTI tools to get account information.
func (s *AccountLTIRegistrationsService) GetLTIAccount(ctx context.Context, accountID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/lti/accounts/%d", accountID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get LTI account: %w", err)
	}

	return result, nil
}

// CreateControlForCurrentAccount creates a control for an LTI registration on the current account
func (s *AccountLTIRegistrationsService) CreateControlForCurrentAccount(ctx context.Context, currentAccountID, registrationID int64, body interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/lti_registrations/%d/controls", currentAccountID, registrationID)

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to create control for current account: %w", err)
	}

	return result, nil
}
