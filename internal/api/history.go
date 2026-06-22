package api

import (
	"context"
	"fmt"
)

// HistoryEntry represents a single page-view history entry for a user.
type HistoryEntry struct {
	AssetCode             string `json:"asset_code,omitempty"`
	AssetName             string `json:"asset_name,omitempty"`
	AssetIcon             string `json:"asset_icon,omitempty"`
	AssetReadableCategory string `json:"asset_readable_category,omitempty"`
	ContextType           string `json:"context_type,omitempty"`
	ContextID             int64  `json:"context_id,omitempty"`
	ContextName           string `json:"context_name,omitempty"`
	VisitedAt             string `json:"visited_at,omitempty"`
	VisitedURL            string `json:"visited_url,omitempty"`
}

// HistoryService handles user page-view history API calls.
type HistoryService struct {
	client *Client
}

// NewHistoryService creates a new HistoryService.
func NewHistoryService(client *Client) *HistoryService {
	return &HistoryService{client: client}
}

// List retrieves the page-view history for a user.
func (s *HistoryService) List(ctx context.Context, userID int64) ([]HistoryEntry, error) {
	path := fmt.Sprintf("/api/v1/users/%d/history", userID)

	var entries []HistoryEntry
	if err := s.client.GetAllPages(ctx, path, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}
