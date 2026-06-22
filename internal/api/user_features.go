package api

import (
	"context"
	"fmt"
)

// UserFeaturesService handles user feature flag API calls
type UserFeaturesService struct {
	client *Client
}

// NewUserFeaturesService creates a new user features service
func NewUserFeaturesService(client *Client) *UserFeaturesService {
	return &UserFeaturesService{client: client}
}

// List retrieves all features for a user
func (s *UserFeaturesService) List(ctx context.Context, userID int64) ([]Feature, error) {
	path := fmt.Sprintf("/api/v1/users/%d/features", userID)
	var result []Feature
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing features for user %d: %w", userID, err)
	}
	return result, nil
}

// ListEnabled retrieves all enabled features for a user
func (s *UserFeaturesService) ListEnabled(ctx context.Context, userID int64) ([]Feature, error) {
	path := fmt.Sprintf("/api/v1/users/%d/features/enabled", userID)
	var result []Feature
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing enabled features for user %d: %w", userID, err)
	}
	return result, nil
}

// GetFlag retrieves the feature flag state for a specific feature and user
func (s *UserFeaturesService) GetFlag(ctx context.Context, userID int64, feature string) (*FeatureFlag, error) {
	path := fmt.Sprintf("/api/v1/users/%d/features/flags/%s", userID, feature)
	var result FeatureFlag
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting feature flag %s for user %d: %w", feature, userID, err)
	}
	return &result, nil
}

// SetFlag sets the state of a feature flag for a user
func (s *UserFeaturesService) SetFlag(ctx context.Context, userID int64, feature, state string) (*FeatureFlag, error) {
	path := fmt.Sprintf("/api/v1/users/%d/features/flags/%s", userID, feature)
	body := map[string]string{"state": state}
	var result FeatureFlag
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("setting feature flag %s for user %d: %w", feature, userID, err)
	}
	return &result, nil
}

// DeleteFlag removes a feature flag override for a user, reverting to the inherited state
func (s *UserFeaturesService) DeleteFlag(ctx context.Context, userID int64, feature string) error {
	path := fmt.Sprintf("/api/v1/users/%d/features/flags/%s", userID, feature)
	_, err := s.client.Delete(ctx, path)
	return err
}
