package api

import (
	"context"
	"fmt"
	"io"
	"net/url"
)

// AccountFeaturesService handles account feature flag API calls
type AccountFeaturesService struct {
	client *Client
}

// NewAccountFeaturesService creates a new account features service
func NewAccountFeaturesService(client *Client) *AccountFeaturesService {
	return &AccountFeaturesService{client: client}
}

// AccountFeature represents a Canvas account feature
type AccountFeature struct {
	Feature     string                 `json:"feature"`
	DisplayName string                 `json:"display_name"`
	AppliesTo   string                 `json:"applies_to"`
	FeatureFlag *AccountFeatureFlag    `json:"feature_flag"`
	Transitions map[string]interface{} `json:"transitions,omitempty"`
}

// AccountFeatureFlag represents the flag state for an account feature
type AccountFeatureFlag struct {
	Feature          string `json:"feature"`
	State            string `json:"state"` // "on", "off", "allowed", "allowed_on"
	LockedAt         string `json:"locked_at,omitempty"`
	TransitionLocked bool   `json:"transition_locked"`
}

// AccountSettings represents Canvas account settings
type AccountSettings struct {
	RestrictStudentPastView   bool `json:"restrict_student_past_view"`
	RestrictStudentFutureView bool `json:"restrict_student_future_view"`
	HideDistributionGraphs    bool `json:"hide_distribution_graphs"`
	LockAllAnnouncements      bool `json:"lock_all_announcements"`
	UsageRightsRequired       bool `json:"usage_rights_required"`
	DefaultDueDateRestricted  bool `json:"default_due_date_restricted"`
}

// List retrieves all features for an account
func (s *AccountFeaturesService) List(ctx context.Context, accountID int64) ([]AccountFeature, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/features", accountID)

	var features []AccountFeature
	if err := s.client.GetAllPages(ctx, path, &features); err != nil {
		return nil, fmt.Errorf("listing account features: %w", err)
	}

	return features, nil
}

// ListEnabled retrieves the names of enabled features for an account
func (s *AccountFeaturesService) ListEnabled(ctx context.Context, accountID int64) ([]string, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/features/enabled", accountID)

	var features []string
	if err := s.client.GetJSON(ctx, path, &features); err != nil {
		return nil, fmt.Errorf("listing enabled account features: %w", err)
	}

	return features, nil
}

// GetFlag retrieves the feature flag for a specific feature on an account
func (s *AccountFeaturesService) GetFlag(ctx context.Context, accountID int64, feature string) (*AccountFeatureFlag, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/features/flags/%s", accountID, url.PathEscape(feature))

	var flag AccountFeatureFlag
	if err := s.client.GetJSON(ctx, path, &flag); err != nil {
		return nil, fmt.Errorf("getting account feature flag %q: %w", feature, err)
	}

	return &flag, nil
}

// SetFlag sets the state of a feature flag for an account
func (s *AccountFeaturesService) SetFlag(ctx context.Context, accountID int64, feature, state string) (*AccountFeatureFlag, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/features/flags/%s", accountID, url.PathEscape(feature))

	body := map[string]interface{}{
		"state": state,
	}

	var flag AccountFeatureFlag
	if err := s.client.PutJSON(ctx, path, body, &flag); err != nil {
		return nil, fmt.Errorf("setting account feature flag %q to %q: %w", feature, state, err)
	}

	return &flag, nil
}

// DeleteFlag removes a feature flag override for an account, reverting to the parent value
func (s *AccountFeaturesService) DeleteFlag(ctx context.Context, accountID int64, feature string) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/features/flags/%s", accountID, url.PathEscape(feature))

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("deleting account feature flag %q: %w", feature, err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}

// GetSettings retrieves settings for an account
func (s *AccountFeaturesService) GetSettings(ctx context.Context, accountID int64) (*AccountSettings, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/settings", accountID)

	var settings AccountSettings
	if err := s.client.GetJSON(ctx, path, &settings); err != nil {
		return nil, fmt.Errorf("getting account settings: %w", err)
	}

	return &settings, nil
}

// GetPermissions retrieves permissions for the current user in an account
func (s *AccountFeaturesService) GetPermissions(ctx context.Context, accountID int64, permissions []string) (map[string]bool, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/permissions", accountID)

	if len(permissions) > 0 {
		query := url.Values{}
		for _, p := range permissions {
			query.Add("permissions[]", p)
		}
		path += "?" + query.Encode()
	}

	var result map[string]bool
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting account permissions: %w", err)
	}

	return result, nil
}

// GetScopes retrieves available API scopes for an account
func (s *AccountFeaturesService) GetScopes(ctx context.Context, accountID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/scopes", accountID)

	var scopes []interface{}
	if err := s.client.GetJSON(ctx, path, &scopes); err != nil {
		return nil, fmt.Errorf("getting account scopes: %w", err)
	}

	return scopes, nil
}
