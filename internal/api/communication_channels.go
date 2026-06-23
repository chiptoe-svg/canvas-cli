package api

import (
	"context"
	"fmt"
)

// CommunicationChannelsService handles communication channel API calls
type CommunicationChannelsService struct {
	client *Client
}

// NewCommunicationChannelsService creates a new communication channels service
func NewCommunicationChannelsService(client *Client) *CommunicationChannelsService {
	return &CommunicationChannelsService{client: client}
}

// CommunicationChannel represents a user's communication channel (e.g. email, SMS)
type CommunicationChannel struct {
	ID            int64  `json:"id"`
	Address       string `json:"address"`
	Type          string `json:"type"`
	Position      int    `json:"position"`
	WorkflowState string `json:"workflow_state"`
	UserID        int64  `json:"user_id"`
	CreatedAt     string `json:"created_at"`
}

// CreateCommunicationChannelParams holds parameters for creating a communication channel
type CreateCommunicationChannelParams struct {
	Address          string `json:"address"`
	Type             string `json:"type"`
	SkipConfirmation bool   `json:"skip_confirmation,omitempty"`
}

// NotificationPreference represents a single notification preference setting
type NotificationPreference struct {
	Frequency    string `json:"frequency"`
	Notification string `json:"notification"`
	Category     string `json:"category"`
}

// NotificationPreferences wraps a list of notification preferences (read response from Canvas).
type NotificationPreferences struct {
	NotificationPreferences []NotificationPreference `json:"notification_preferences"`
}

// NotificationPreferenceUpdate holds the frequency update for a single notification event.
type NotificationPreferenceUpdate struct {
	Frequency string `json:"frequency"`
}

// NotificationPreferencesUpdateBody is the request body for bulk notification preference updates.
// Canvas expects {"notification_preferences": {event_name: {"frequency": "..."}, ...}}.
type NotificationPreferencesUpdateBody struct {
	NotificationPreferences map[string]NotificationPreferenceUpdate `json:"notification_preferences"`
}

// List retrieves all communication channels for a user
func (s *CommunicationChannelsService) List(ctx context.Context, userID int64) ([]CommunicationChannel, error) {
	path := fmt.Sprintf("/api/v1/users/%d/communication_channels", userID)
	var result []CommunicationChannel
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing communication channels for user %d: %w", userID, err)
	}
	return result, nil
}

// Create creates a new communication channel for a user
func (s *CommunicationChannelsService) Create(ctx context.Context, userID int64, params CreateCommunicationChannelParams) (*CommunicationChannel, error) {
	path := fmt.Sprintf("/api/v1/users/%d/communication_channels", userID)
	body := map[string]interface{}{
		"communication_channel": params,
	}
	var result CommunicationChannel
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("creating communication channel for user %d: %w", userID, err)
	}
	return &result, nil
}

// Delete deletes a communication channel by ID
func (s *CommunicationChannelsService) Delete(ctx context.Context, userID, channelID int64) (*CommunicationChannel, error) {
	path := fmt.Sprintf("/api/v1/users/%d/communication_channels/%d", userID, channelID)
	var result CommunicationChannel
	if err := s.client.DeleteJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("deleting communication channel %d for user %d: %w", channelID, userID, err)
	}
	return &result, nil
}

// DeleteByTypeAddress deletes a communication channel identified by type and address
func (s *CommunicationChannelsService) DeleteByTypeAddress(ctx context.Context, userID int64, channelType, address string) (*CommunicationChannel, error) {
	path := fmt.Sprintf("/api/v1/users/%d/communication_channels/%s/%s", userID, channelType, address)
	var result CommunicationChannel
	if err := s.client.DeleteJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("deleting communication channel %s/%s for user %d: %w", channelType, address, userID, err)
	}
	return &result, nil
}

// GetNotificationPreferences retrieves notification preferences for a channel
func (s *CommunicationChannelsService) GetNotificationPreferences(ctx context.Context, userID, channelID int64) (*NotificationPreferences, error) {
	path := fmt.Sprintf("/api/v1/users/%d/communication_channels/%d/notification_preferences", userID, channelID)
	var result NotificationPreferences
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting notification preferences for channel %d: %w", channelID, err)
	}
	return &result, nil
}

// GetNotificationPreferenceCategories retrieves notification preference categories for a channel
func (s *CommunicationChannelsService) GetNotificationPreferenceCategories(ctx context.Context, userID, channelID int64) ([]string, error) {
	path := fmt.Sprintf("/api/v1/users/%d/communication_channels/%d/notification_preference_categories", userID, channelID)
	var result []string
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting notification preference categories for channel %d: %w", channelID, err)
	}
	return result, nil
}

// GetNotificationPreference retrieves a single notification preference for a channel
func (s *CommunicationChannelsService) GetNotificationPreference(ctx context.Context, userID, channelID int64, notification string) (*NotificationPreference, error) {
	path := fmt.Sprintf("/api/v1/users/%d/communication_channels/%d/notification_preferences/%s", userID, channelID, notification)
	var result NotificationPreference
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting notification preference %s for channel %d: %w", notification, channelID, err)
	}
	return &result, nil
}

// GetNotificationPreferencesByType retrieves notification preferences for a channel identified by type/address
func (s *CommunicationChannelsService) GetNotificationPreferencesByType(ctx context.Context, userID int64, channelType, address string) (*NotificationPreferences, error) {
	path := fmt.Sprintf("/api/v1/users/%d/communication_channels/%s/%s/notification_preferences", userID, channelType, address)
	var result NotificationPreferences
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting notification preferences for channel %s/%s: %w", channelType, address, err)
	}
	return &result, nil
}

// GetNotificationPreferenceByType retrieves a single notification preference for a channel by type/address
func (s *CommunicationChannelsService) GetNotificationPreferenceByType(ctx context.Context, userID int64, channelType, address, notification string) (*NotificationPreference, error) {
	path := fmt.Sprintf("/api/v1/users/%d/communication_channels/%s/%s/notification_preferences/%s", userID, channelType, address, notification)
	var result NotificationPreference
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting notification preference %s for channel %s/%s: %w", notification, channelType, address, err)
	}
	return &result, nil
}

// UpdateNotificationPreferences updates all notification preferences for a channel.
// Canvas expects {"notification_preferences": {event_name: {"frequency": "..."}, ...}}.
func (s *CommunicationChannelsService) UpdateNotificationPreferences(ctx context.Context, channelID int64, prefs NotificationPreferencesUpdateBody) (*NotificationPreferences, error) {
	path := fmt.Sprintf("/api/v1/users/self/communication_channels/%d/notification_preferences", channelID)
	var result NotificationPreferences
	if err := s.client.PutJSON(ctx, path, prefs, &result); err != nil {
		return nil, fmt.Errorf("updating notification preferences for channel %d: %w", channelID, err)
	}
	return &result, nil
}

// UpdateNotificationPreference updates a single notification preference for a channel
func (s *CommunicationChannelsService) UpdateNotificationPreference(ctx context.Context, channelID int64, notification string, pref NotificationPreference) (*NotificationPreference, error) {
	path := fmt.Sprintf("/api/v1/users/self/communication_channels/%d/notification_preferences/%s", channelID, notification)
	var result NotificationPreference
	if err := s.client.PutJSON(ctx, path, pref, &result); err != nil {
		return nil, fmt.Errorf("updating notification preference %s for channel %d: %w", notification, channelID, err)
	}
	return &result, nil
}

// UpdateNotificationPreferenceCategory updates all notification preferences in a category for a channel
func (s *CommunicationChannelsService) UpdateNotificationPreferenceCategory(ctx context.Context, channelID int64, category string, prefs NotificationPreferences) (*NotificationPreferences, error) {
	path := fmt.Sprintf("/api/v1/users/self/communication_channels/%d/notification_preference_categories/%s", channelID, category)
	var result NotificationPreferences
	if err := s.client.PutJSON(ctx, path, prefs, &result); err != nil {
		return nil, fmt.Errorf("updating notification preference category %s for channel %d: %w", category, channelID, err)
	}
	return &result, nil
}

// UpdateNotificationPreferencesByType updates notification preferences for a channel identified by type/address.
// Canvas expects {"notification_preferences": {event_name: {"frequency": "..."}, ...}}.
func (s *CommunicationChannelsService) UpdateNotificationPreferencesByType(ctx context.Context, channelType, address string, prefs NotificationPreferencesUpdateBody) (*NotificationPreferences, error) {
	path := fmt.Sprintf("/api/v1/users/self/communication_channels/%s/%s/notification_preferences", channelType, address)
	var result NotificationPreferences
	if err := s.client.PutJSON(ctx, path, prefs, &result); err != nil {
		return nil, fmt.Errorf("updating notification preferences for channel %s/%s: %w", channelType, address, err)
	}
	return &result, nil
}

// UpdateNotificationPreferenceByType updates a single notification preference for a channel by type/address
func (s *CommunicationChannelsService) UpdateNotificationPreferenceByType(ctx context.Context, channelType, address, notification string, pref NotificationPreference) (*NotificationPreference, error) {
	path := fmt.Sprintf("/api/v1/users/self/communication_channels/%s/%s/notification_preferences/%s", channelType, address, notification)
	var result NotificationPreference
	if err := s.client.PutJSON(ctx, path, pref, &result); err != nil {
		return nil, fmt.Errorf("updating notification preference %s for channel %s/%s: %w", notification, channelType, address, err)
	}
	return &result, nil
}

// DeletePushChannel removes the push notification channel for the current user
func (s *CommunicationChannelsService) DeletePushChannel(ctx context.Context) error {
	path := "/api/v1/users/self/communication_channels/push"
	_, err := s.client.Delete(ctx, path)
	return err
}
