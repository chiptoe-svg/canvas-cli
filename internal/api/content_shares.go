package api

import (
	"context"
	"fmt"
)

// ContentSharesService handles content sharing API calls
type ContentSharesService struct {
	client *Client
}

// NewContentSharesService creates a new content shares service
func NewContentSharesService(client *Client) *ContentSharesService {
	return &ContentSharesService{client: client}
}

// ContentShareSender represents the sender or receiver of a content share
type ContentShareSender struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	HTMLUrl     string `json:"html_url"`
}

// ContentShare represents a shared content item between Canvas users
type ContentShare struct {
	ID          int64                `json:"id"`
	Name        string               `json:"name"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
	UserID      int64                `json:"user_id"`
	ReadState   string               `json:"read_state"`
	Type        string               `json:"type"`
	ContentType string               `json:"content_type"`
	CourseID    int64                `json:"course_id"`
	Sender      *ContentShareSender  `json:"sender,omitempty"`
	Receivers   []ContentShareSender `json:"receivers,omitempty"`
}

// CreateContentShareParams holds parameters for creating a content share
type CreateContentShareParams struct {
	ReceiverIDs []int64 `json:"receiver_ids"`
	ContentType string  `json:"content_type"`
	ContentID   int64   `json:"content_id"`
}

// UpdateContentShareParams holds parameters for updating a content share
type UpdateContentShareParams struct {
	ReadState string `json:"read_state,omitempty"`
}

// AddContentShareUsersParams holds parameters for adding users to a content share
type AddContentShareUsersParams struct {
	ReceiverIDs []int64 `json:"receiver_ids"`
}

// ContentShareUnreadCount holds the count of unread content shares
type ContentShareUnreadCount struct {
	UnreadCount int `json:"unread_count"`
}

// Create sends content to one or more users
func (s *ContentSharesService) Create(ctx context.Context, userID int64, params CreateContentShareParams) (*ContentShare, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_shares", userID)
	var result ContentShare
	if err := s.client.PostJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("creating content share for user %d: %w", userID, err)
	}
	return &result, nil
}

// ListSent retrieves content shares sent by a user
func (s *ContentSharesService) ListSent(ctx context.Context, userID int64) ([]ContentShare, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_shares/sent", userID)
	var result []ContentShare
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing sent content shares for user %d: %w", userID, err)
	}
	return result, nil
}

// ListReceived retrieves content shares received by a user
func (s *ContentSharesService) ListReceived(ctx context.Context, userID int64) ([]ContentShare, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_shares/received", userID)
	var result []ContentShare
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing received content shares for user %d: %w", userID, err)
	}
	return result, nil
}

// GetUnreadCount retrieves the number of unread received content shares for a user
func (s *ContentSharesService) GetUnreadCount(ctx context.Context, userID int64) (*ContentShareUnreadCount, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_shares/unread_count", userID)
	var result ContentShareUnreadCount
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting unread count for user %d: %w", userID, err)
	}
	return &result, nil
}

// Get retrieves a specific content share for a user
func (s *ContentSharesService) Get(ctx context.Context, userID, id int64) (*ContentShare, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_shares/%d", userID, id)
	var result ContentShare
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting content share %d for user %d: %w", id, userID, err)
	}
	return &result, nil
}

// Update updates a content share (e.g. marks it as read)
func (s *ContentSharesService) Update(ctx context.Context, userID, id int64, params UpdateContentShareParams) (*ContentShare, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_shares/%d", userID, id)
	var result ContentShare
	if err := s.client.PutJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("updating content share %d for user %d: %w", id, userID, err)
	}
	return &result, nil
}

// Delete removes a content share for a user
func (s *ContentSharesService) Delete(ctx context.Context, userID, id int64) error {
	path := fmt.Sprintf("/api/v1/users/%d/content_shares/%d", userID, id)
	_, err := s.client.Delete(ctx, path)
	return err
}

// AddUsers adds additional receivers to an existing content share
func (s *ContentSharesService) AddUsers(ctx context.Context, userID, id int64, params AddContentShareUsersParams) (*ContentShare, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_shares/%d/add_users", userID, id)
	var result ContentShare
	if err := s.client.PostJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("adding users to content share %d for user %d: %w", id, userID, err)
	}
	return &result, nil
}
