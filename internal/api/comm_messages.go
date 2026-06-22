package api

import (
	"context"
	"net/url"
	"strconv"
)

// CommMessage represents a Canvas communication message sent to a user.
type CommMessage struct {
	ID            int64  `json:"id"`
	CreatedAt     string `json:"created_at,omitempty"`
	SentAt        string `json:"sent_at,omitempty"`
	WorkflowState string `json:"workflow_state,omitempty"`
	From          string `json:"from,omitempty"`
	FromName      string `json:"from_name,omitempty"`
	To            string `json:"to,omitempty"`
	ReplyTo       string `json:"reply_to,omitempty"`
	Subject       string `json:"subject,omitempty"`
	Body          string `json:"body,omitempty"`
	HTMLBody      string `json:"html_body,omitempty"`
	Author        *User  `json:"author,omitempty"`
}

// CommMessagesService handles communication-message API calls.
type CommMessagesService struct {
	client *Client
}

// NewCommMessagesService creates a new CommMessagesService.
func NewCommMessagesService(client *Client) *CommMessagesService {
	return &CommMessagesService{client: client}
}

// ListCommMessagesOptions holds query parameters for listing comm messages.
type ListCommMessagesOptions struct {
	UserID  int64
	PerPage int
}

// List retrieves communication messages for a user.
func (s *CommMessagesService) List(ctx context.Context, opts *ListCommMessagesOptions) ([]CommMessage, error) {
	path := "/api/v1/comm_messages"

	query := url.Values{}
	if opts != nil {
		if opts.UserID > 0 {
			query.Set("user_id", strconv.FormatInt(opts.UserID, 10))
		}
		if opts.PerPage > 0 {
			query.Set("per_page", strconv.Itoa(opts.PerPage))
		}
	}
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var msgs []CommMessage
	if err := s.client.GetAllPages(ctx, path, &msgs); err != nil {
		return nil, err
	}

	return msgs, nil
}
