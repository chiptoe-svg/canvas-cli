package api

import (
	"context"
	"fmt"
	"io"
)

// AccountNotificationsService handles account notification API calls
type AccountNotificationsService struct {
	client *Client
}

// NewAccountNotificationsService creates a new account notifications service
func NewAccountNotificationsService(client *Client) *AccountNotificationsService {
	return &AccountNotificationsService{client: client}
}

// AccountNotification represents a Canvas account notification
type AccountNotification struct {
	ID      int64    `json:"id"`
	Subject string   `json:"subject"`
	Message string   `json:"message"`
	StartAt string   `json:"start_at"`
	EndAt   string   `json:"end_at"`
	Icon    string   `json:"icon"`
	Roles   []string `json:"roles"`
	RoleIDs []int64  `json:"role_ids"`
}

// AccountNotificationParams holds parameters for creating or updating an account notification
type AccountNotificationParams struct {
	Subject string   `json:"account_notification[subject]"`
	Message string   `json:"account_notification[message]"`
	StartAt string   `json:"account_notification[start_at]"`
	EndAt   string   `json:"account_notification[end_at]"`
	Icon    string   `json:"account_notification[icon],omitempty"`
	Roles   []string `json:"account_notification_roles,omitempty"`
}

// List retrieves all account notifications for an account
func (s *AccountNotificationsService) List(ctx context.Context, accountID int64) ([]AccountNotification, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/account_notifications", accountID)

	var notifications []AccountNotification
	if err := s.client.GetAllPages(ctx, path, &notifications); err != nil {
		return nil, fmt.Errorf("list account notifications: %w", err)
	}

	return notifications, nil
}

// Create creates a new account notification
func (s *AccountNotificationsService) Create(ctx context.Context, accountID int64, body *AccountNotificationParams) (*AccountNotification, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/account_notifications", accountID)

	var notification AccountNotification
	if err := s.client.PostJSON(ctx, path, body, &notification); err != nil {
		return nil, fmt.Errorf("create account notification: %w", err)
	}

	return &notification, nil
}

// Get retrieves a single account notification
func (s *AccountNotificationsService) Get(ctx context.Context, accountID, id int64) (*AccountNotification, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/account_notifications/%d", accountID, id)

	var notification AccountNotification
	if err := s.client.GetJSON(ctx, path, &notification); err != nil {
		return nil, fmt.Errorf("get account notification: %w", err)
	}

	return &notification, nil
}

// Update updates an account notification
func (s *AccountNotificationsService) Update(ctx context.Context, accountID, id int64, body *AccountNotificationParams) (*AccountNotification, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/account_notifications/%d", accountID, id)

	var notification AccountNotification
	if err := s.client.PutJSON(ctx, path, body, &notification); err != nil {
		return nil, fmt.Errorf("update account notification: %w", err)
	}

	return &notification, nil
}

// Delete deletes an account notification
func (s *AccountNotificationsService) Delete(ctx context.Context, accountID, id int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/account_notifications/%d", accountID, id)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("delete account notification: %w", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}
