package api

import (
	"context"
	"fmt"
)

// Eportfolio represents a Canvas ePortfolio.
type Eportfolio struct {
	ID            int64  `json:"id"`
	UserID        int64  `json:"user_id,omitempty"`
	Name          string `json:"name,omitempty"`
	Public        bool   `json:"public,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
	WorkflowState string `json:"workflow_state,omitempty"`
}

// EportfolioPage represents a page inside a Canvas ePortfolio.
type EportfolioPage struct {
	ID            int64  `json:"id"`
	EportfolioID  int64  `json:"eportfolio_id,omitempty"`
	Position      int    `json:"position,omitempty"`
	Name          string `json:"name,omitempty"`
	WorkflowState string `json:"workflow_state,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
	Content       string `json:"content,omitempty"`
}

// EportfoliosService handles ePortfolio API calls.
type EportfoliosService struct {
	client *Client
}

// NewEportfoliosService creates a new EportfoliosService.
func NewEportfoliosService(client *Client) *EportfoliosService {
	return &EportfoliosService{client: client}
}

// ListForUser retrieves all ePortfolios for a user.
func (s *EportfoliosService) ListForUser(ctx context.Context, userID int64) ([]Eportfolio, error) {
	path := fmt.Sprintf("/api/v1/users/%d/eportfolios", userID)

	var eps []Eportfolio
	if err := s.client.GetAllPages(ctx, path, &eps); err != nil {
		return nil, err
	}

	return eps, nil
}

// Get retrieves a single ePortfolio by ID.
func (s *EportfoliosService) Get(ctx context.Context, id int64) (*Eportfolio, error) {
	path := fmt.Sprintf("/api/v1/eportfolios/%d", id)

	var ep Eportfolio
	if err := s.client.GetJSON(ctx, path, &ep); err != nil {
		return nil, err
	}

	return &ep, nil
}

// Delete deletes an ePortfolio by ID.
func (s *EportfoliosService) Delete(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/api/v1/eportfolios/%d", id)
	_, err := s.client.Delete(ctx, path)
	return err
}

// Moderate sets the moderation state of an ePortfolio (spam / allow / etc.).
func (s *EportfoliosService) Moderate(ctx context.Context, id int64, spamStatus string) (*Eportfolio, error) {
	path := fmt.Sprintf("/api/v1/eportfolios/%d/moderate", id)

	body := map[string]interface{}{
		"eportfolio": map[string]interface{}{
			"spam_status": spamStatus,
		},
	}

	var ep Eportfolio
	if err := s.client.PutJSON(ctx, path, body, &ep); err != nil {
		return nil, err
	}

	return &ep, nil
}

// Restore restores a deleted ePortfolio.
func (s *EportfoliosService) Restore(ctx context.Context, id int64) (*Eportfolio, error) {
	path := fmt.Sprintf("/api/v1/eportfolios/%d/restore", id)

	var ep Eportfolio
	if err := s.client.PutJSON(ctx, path, nil, &ep); err != nil {
		return nil, err
	}

	return &ep, nil
}

// ListPages retrieves all pages in an ePortfolio.
func (s *EportfoliosService) ListPages(ctx context.Context, eportfolioID int64) ([]EportfolioPage, error) {
	path := fmt.Sprintf("/api/v1/eportfolios/%d/pages", eportfolioID)

	var pages []EportfolioPage
	if err := s.client.GetAllPages(ctx, path, &pages); err != nil {
		return nil, err
	}

	return pages, nil
}
