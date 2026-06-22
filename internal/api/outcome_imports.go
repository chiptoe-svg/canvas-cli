package api

import (
	"context"
	"fmt"
)

// AccountOutcomeImportsService handles outcome import API calls for accounts
type AccountOutcomeImportsService struct {
	client *Client
}

// NewAccountOutcomeImportsService creates a new account outcome imports service
func NewAccountOutcomeImportsService(client *Client) *AccountOutcomeImportsService {
	return &AccountOutcomeImportsService{client: client}
}

// OutcomeImport represents a Canvas outcome import
type OutcomeImport struct {
	ID                 int64                  `json:"id"`
	CreatedAt          string                 `json:"created_at"`
	EndedAt            string                 `json:"ended_at,omitempty"`
	UpdatedAt          string                 `json:"updated_at"`
	WorkflowState      string                 `json:"workflow_state"`
	Data               map[string]interface{} `json:"data,omitempty"`
	Progress           float64                `json:"progress"`
	UserID             int64                  `json:"user_id"`
	ProcessingErrors   [][]string             `json:"processing_errors,omitempty"`
	ProcessingWarnings [][]string             `json:"processing_warnings,omitempty"`
}

// Create initiates an outcome import for an account
func (s *AccountOutcomeImportsService) Create(ctx context.Context, accountID int64, body map[string]interface{}) (*OutcomeImport, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_imports", accountID)

	var result OutcomeImport
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to create outcome import: %w", err)
	}

	return &result, nil
}

// Get retrieves an outcome import by ID
func (s *AccountOutcomeImportsService) Get(ctx context.Context, accountID, id int64) (*OutcomeImport, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_imports/%d", accountID, id)

	var result OutcomeImport
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get outcome import: %w", err)
	}

	return &result, nil
}

// GetCreatedGroupIDs returns the IDs of outcome groups created during an import
func (s *AccountOutcomeImportsService) GetCreatedGroupIDs(ctx context.Context, accountID, id int64) ([]int64, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_imports/%d/created_group_ids", accountID, id)

	var result []int64
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get created group IDs: %w", err)
	}

	return result, nil
}

// GetOutcomeGroupLinks returns outcome group links for an account
func (s *AccountOutcomeImportsService) GetOutcomeGroupLinks(ctx context.Context, accountID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_group_links", accountID)

	var result []interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get outcome group links: %w", err)
	}

	return result, nil
}

// GetRootOutcomeGroup returns the root outcome group for an account
func (s *AccountOutcomeImportsService) GetRootOutcomeGroup(ctx context.Context, accountID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/root_outcome_group", accountID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get root outcome group: %w", err)
	}

	return result, nil
}

// GetOutcomeGroupSubgroups returns subgroups of an outcome group
func (s *AccountOutcomeImportsService) GetOutcomeGroupSubgroups(ctx context.Context, accountID, groupID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_groups/%d/subgroups", accountID, groupID)

	var result []interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get outcome group subgroups: %w", err)
	}

	return result, nil
}

// CreateOutcomeGroupSubgroup creates a subgroup within an outcome group
func (s *AccountOutcomeImportsService) CreateOutcomeGroupSubgroup(ctx context.Context, accountID, groupID int64, body map[string]interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_groups/%d/subgroups", accountID, groupID)

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to create outcome group subgroup: %w", err)
	}

	return result, nil
}

// ImportOutcomeGroup imports an outcome group into another group
func (s *AccountOutcomeImportsService) ImportOutcomeGroup(ctx context.Context, accountID, groupID int64, body map[string]interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_groups/%d/import", accountID, groupID)

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to import outcome group: %w", err)
	}

	return result, nil
}

// GetOutcomeProficiency returns the outcome proficiency settings for an account
func (s *AccountOutcomeImportsService) GetOutcomeProficiency(ctx context.Context, accountID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_proficiency", accountID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get outcome proficiency: %w", err)
	}

	return result, nil
}

// UpdateOutcomeProficiency creates or updates outcome proficiency settings for an account
func (s *AccountOutcomeImportsService) UpdateOutcomeProficiency(ctx context.Context, accountID int64, body map[string]interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_proficiency", accountID)

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to update outcome proficiency: %w", err)
	}

	return result, nil
}
