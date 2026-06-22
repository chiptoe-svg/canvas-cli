package api

import (
	"context"
	"fmt"
)

// BookmarksService handles bookmark-related API calls
type BookmarksService struct {
	client *Client
}

// NewBookmarksService creates a new bookmarks service
func NewBookmarksService(client *Client) *BookmarksService {
	return &BookmarksService{client: client}
}

// Bookmark represents a Canvas user bookmark
type Bookmark struct {
	ID       int64                  `json:"id"`
	Name     string                 `json:"name"`
	URL      string                 `json:"url"`
	Position int                    `json:"position"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

// CreateBookmarkParams holds parameters for creating a bookmark
type CreateBookmarkParams struct {
	Name     string                 `json:"name"`
	URL      string                 `json:"url"`
	Position int                    `json:"position,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

// UpdateBookmarkParams holds parameters for updating a bookmark
type UpdateBookmarkParams struct {
	Name     string                 `json:"name,omitempty"`
	URL      string                 `json:"url,omitempty"`
	Position int                    `json:"position,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

// List retrieves all bookmarks for the current user
func (s *BookmarksService) List(ctx context.Context) ([]Bookmark, error) {
	path := "/api/v1/users/self/bookmarks"
	var result []Bookmark
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing bookmarks: %w", err)
	}
	return result, nil
}

// Create creates a new bookmark for the current user
func (s *BookmarksService) Create(ctx context.Context, params CreateBookmarkParams) (*Bookmark, error) {
	path := "/api/v1/users/self/bookmarks"
	var result Bookmark
	if err := s.client.PostJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("creating bookmark: %w", err)
	}
	return &result, nil
}

// Get retrieves a single bookmark by ID
func (s *BookmarksService) Get(ctx context.Context, id int64) (*Bookmark, error) {
	path := fmt.Sprintf("/api/v1/users/self/bookmarks/%d", id)
	var result Bookmark
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting bookmark %d: %w", id, err)
	}
	return &result, nil
}

// Update updates an existing bookmark
func (s *BookmarksService) Update(ctx context.Context, id int64, params UpdateBookmarkParams) (*Bookmark, error) {
	path := fmt.Sprintf("/api/v1/users/self/bookmarks/%d", id)
	var result Bookmark
	if err := s.client.PutJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("updating bookmark %d: %w", id, err)
	}
	return &result, nil
}

// Delete deletes a bookmark by ID
func (s *BookmarksService) Delete(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/api/v1/users/self/bookmarks/%d", id)
	_, err := s.client.Delete(ctx, path)
	return err
}
