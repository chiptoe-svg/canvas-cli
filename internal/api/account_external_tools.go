package api

import (
	"context"
	"fmt"
	"io"
)

// AccountExternalToolsService handles account-level external tool (LTI) favorites
type AccountExternalToolsService struct {
	client *Client
}

// NewAccountExternalToolsService creates a new account external tools service
func NewAccountExternalToolsService(client *Client) *AccountExternalToolsService {
	return &AccountExternalToolsService{client: client}
}

// AddRCEFavorite marks an external tool as an RCE favorite for an account
func (s *AccountExternalToolsService) AddRCEFavorite(ctx context.Context, accountID, toolID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/external_tools/rce_favorites/%d", accountID, toolID)

	resp, err := s.client.Post(ctx, path, nil)
	if err != nil {
		return fmt.Errorf("failed to add RCE favorite: %w", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	return nil
}

// RemoveRCEFavorite removes an external tool from the RCE favorites for an account
func (s *AccountExternalToolsService) RemoveRCEFavorite(ctx context.Context, accountID, toolID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/external_tools/rce_favorites/%d", accountID, toolID)

	// Delete already drains the body internally
	if _, err := s.client.Delete(ctx, path); err != nil {
		return fmt.Errorf("failed to remove RCE favorite: %w", err)
	}
	return nil
}

// AddTopNavFavorite marks an external tool as a top-nav favorite for an account
func (s *AccountExternalToolsService) AddTopNavFavorite(ctx context.Context, accountID, toolID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/external_tools/top_nav_favorites/%d", accountID, toolID)

	resp, err := s.client.Post(ctx, path, nil)
	if err != nil {
		return fmt.Errorf("failed to add top-nav favorite: %w", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	return nil
}

// RemoveTopNavFavorite removes an external tool from the top-nav favorites for an account
func (s *AccountExternalToolsService) RemoveTopNavFavorite(ctx context.Context, accountID, toolID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/external_tools/top_nav_favorites/%d", accountID, toolID)

	// Delete already drains the body internally
	if _, err := s.client.Delete(ctx, path); err != nil {
		return fmt.Errorf("failed to remove top-nav favorite: %w", err)
	}
	return nil
}
