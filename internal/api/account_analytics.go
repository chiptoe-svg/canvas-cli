package api

import (
	"context"
	"fmt"
)

// AccountAnalyticsService handles account-level analytics API calls
type AccountAnalyticsService struct {
	client *Client
}

// NewAccountAnalyticsService creates a new account analytics service
func NewAccountAnalyticsService(client *Client) *AccountAnalyticsService {
	return &AccountAnalyticsService{client: client}
}

// GetTermActivity retrieves activity data for an account term
func (s *AccountAnalyticsService) GetTermActivity(ctx context.Context, accountID, termID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/analytics/terms/%d/activity", accountID, termID)

	var result []interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting term activity for account %d term %d: %w", accountID, termID, err)
	}

	return result, nil
}

// GetTermGrades retrieves grade data for an account term
func (s *AccountAnalyticsService) GetTermGrades(ctx context.Context, accountID, termID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/analytics/terms/%d/grades", accountID, termID)

	var result []interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting term grades for account %d term %d: %w", accountID, termID, err)
	}

	return result, nil
}

// GetTermStatistics retrieves statistics for an account term
func (s *AccountAnalyticsService) GetTermStatistics(ctx context.Context, accountID, termID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/analytics/terms/%d/statistics", accountID, termID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting term statistics for account %d term %d: %w", accountID, termID, err)
	}

	return result, nil
}

// GetTermStatsBySubaccount retrieves per-subaccount statistics for an account term
func (s *AccountAnalyticsService) GetTermStatsBySubaccount(ctx context.Context, accountID, termID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/analytics/terms/%d/statistics_by_subaccount", accountID, termID)

	var result []interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting term stats by subaccount for account %d term %d: %w", accountID, termID, err)
	}

	return result, nil
}

// GetCompletedActivity retrieves activity data for completed courses in an account
func (s *AccountAnalyticsService) GetCompletedActivity(ctx context.Context, accountID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/analytics/completed/activity", accountID)

	var result []interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting completed activity for account %d: %w", accountID, err)
	}

	return result, nil
}

// GetCompletedGrades retrieves grade data for completed courses in an account
func (s *AccountAnalyticsService) GetCompletedGrades(ctx context.Context, accountID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/analytics/completed/grades", accountID)

	var result []interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting completed grades for account %d: %w", accountID, err)
	}

	return result, nil
}

// GetCompletedStatistics retrieves statistics for completed courses in an account
func (s *AccountAnalyticsService) GetCompletedStatistics(ctx context.Context, accountID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/analytics/completed/statistics", accountID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting completed statistics for account %d: %w", accountID, err)
	}

	return result, nil
}

// GetCompletedStatsBySubaccount retrieves per-subaccount statistics for completed courses in an account
func (s *AccountAnalyticsService) GetCompletedStatsBySubaccount(ctx context.Context, accountID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/analytics/completed/statistics_by_subaccount", accountID)

	var result []interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting completed stats by subaccount for account %d: %w", accountID, err)
	}

	return result, nil
}

// GetCurrentStatsBySubaccount retrieves per-subaccount statistics for current courses in an account
func (s *AccountAnalyticsService) GetCurrentStatsBySubaccount(ctx context.Context, accountID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/analytics/current/statistics_by_subaccount", accountID)

	var result []interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting current stats by subaccount for account %d: %w", accountID, err)
	}

	return result, nil
}
