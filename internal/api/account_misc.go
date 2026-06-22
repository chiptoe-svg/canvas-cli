package api

import (
	"context"
	"fmt"
	"net/url"
)

// AccountMiscService handles miscellaneous account-level Canvas API endpoints
type AccountMiscService struct {
	client *Client
}

// NewAccountMiscService creates a new account misc service
func NewAccountMiscService(client *Client) *AccountMiscService {
	return &AccountMiscService{client: client}
}

// Search searches for accounts by name
func (s *AccountMiscService) Search(ctx context.Context, searchTerm string, accountID int64) ([]Account, error) {
	query := url.Values{}
	if searchTerm != "" {
		query.Set("name", searchTerm)
	}
	if accountID > 0 {
		query.Set("account_id", fmt.Sprintf("%d", accountID))
	}

	path := "/api/v1/accounts/search"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var accounts []Account
	if err := s.client.GetAllPages(ctx, path, &accounts); err != nil {
		return nil, fmt.Errorf("failed to search accounts: %w", err)
	}

	return accounts, nil
}

// GetBrandVariables returns the brand variables for an account
func (s *AccountMiscService) GetBrandVariables(ctx context.Context, accountID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/brand_variables", accountID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get brand variables: %w", err)
	}

	return result, nil
}

// GetHelpLinks returns the help links for an account
func (s *AccountMiscService) GetHelpLinks(ctx context.Context, accountID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/help_links", accountID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get help links: %w", err)
	}

	return result, nil
}

// GetTermsOfService returns the terms of service for an account
func (s *AccountMiscService) GetTermsOfService(ctx context.Context, accountID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/terms_of_service", accountID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get terms of service: %w", err)
	}

	return result, nil
}

// GetGradingStandards returns the grading standards for an account
func (s *AccountMiscService) GetGradingStandards(ctx context.Context, accountID int64) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/grading_standards", accountID)

	var result []map[string]interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get grading standards: %w", err)
	}

	return result, nil
}
