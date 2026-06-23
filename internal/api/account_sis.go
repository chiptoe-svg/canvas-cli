package api

import (
	"context"
	"fmt"
)

// AccountSISService handles account-level SIS import management endpoints
// that go beyond the per-import operations in SISImportsService.
type AccountSISService struct {
	client *Client
}

// NewAccountSISService creates a new account SIS service
func NewAccountSISService(client *Client) *AccountSISService {
	return &AccountSISService{client: client}
}

// sisImportErrorsResponse wraps the Canvas API envelope for SIS import errors.
type sisImportErrorsResponse struct {
	SISImportErrors []interface{} `json:"sis_import_errors"`
}

// GetSISImportErrors returns all SIS import errors for an account (not scoped to a specific import)
func (s *AccountSISService) GetSISImportErrors(ctx context.Context, accountID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/sis_import_errors", accountID)

	var resp sisImportErrorsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get SIS import errors: %w", err)
	}

	return resp.SISImportErrors, nil
}

// AbortAllPendingSISImports aborts all pending SIS imports for an account
func (s *AccountSISService) AbortAllPendingSISImports(ctx context.Context, accountID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/sis_imports/abort_all_pending", accountID)

	var result map[string]interface{}
	if err := s.client.PutJSON(ctx, path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to abort all pending SIS imports: %w", err)
	}

	return result, nil
}

// sisImportsResponse wraps the Canvas API envelope for SIS imports.
type sisImportsResponse struct {
	SISImports []interface{} `json:"sis_imports"`
}

// GetImportingSISImports returns SIS imports that are currently being imported for an account
func (s *AccountSISService) GetImportingSISImports(ctx context.Context, accountID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/sis_imports/importing", accountID)

	var resp sisImportsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get importing SIS imports: %w", err)
	}

	return resp.SISImports, nil
}
