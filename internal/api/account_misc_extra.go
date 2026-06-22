package api

import (
	"context"
	"fmt"
	"io"
)

// AccountMiscExtraService handles additional miscellaneous account-level Canvas API endpoints
type AccountMiscExtraService struct {
	client *Client
}

// NewAccountMiscExtraService creates a new account misc extra service
func NewAccountMiscExtraService(client *Client) *AccountMiscExtraService {
	return &AccountMiscExtraService{client: client}
}

// AccountRestoreUser restores a deleted user in an account
func (s *AccountMiscExtraService) AccountRestoreUser(ctx context.Context, accountID, userID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/users/%d/restore", accountID, userID)

	var result interface{}
	if err := s.client.PutJSON(ctx, path, nil, &result); err != nil {
		return fmt.Errorf("failed to restore user: %w", err)
	}

	return nil
}

// AccountBulkUpdateUsers performs a bulk update of users in an account
func (s *AccountMiscExtraService) AccountBulkUpdateUsers(ctx context.Context, accountID int64, body interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/users/bulk_update", accountID)

	var result map[string]interface{}
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to bulk update users: %w", err)
	}

	return result, nil
}

// AccountBulkEnrollment performs a bulk enrollment operation in an account
func (s *AccountMiscExtraService) AccountBulkEnrollment(ctx context.Context, accountID int64, body interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/bulk_enrollment", accountID)

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to bulk enroll: %w", err)
	}

	return result, nil
}

// AccountDeleteSubAccount deletes a sub-account
func (s *AccountMiscExtraService) AccountDeleteSubAccount(ctx context.Context, accountID, subID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/sub_accounts/%d", accountID, subID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete sub-account: %w", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}

// AccountGetEnrollment retrieves a specific enrollment in an account
func (s *AccountMiscExtraService) AccountGetEnrollment(ctx context.Context, accountID, enrollmentID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/enrollments/%d", accountID, enrollmentID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}

	return result, nil
}

// AccountGetCourse retrieves a specific course in an account
func (s *AccountMiscExtraService) AccountGetCourse(ctx context.Context, accountID, courseID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/courses/%d", accountID, courseID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	return result, nil
}

// AccountSelfRegistration enables or configures self-registration for an account
func (s *AccountMiscExtraService) AccountSelfRegistration(ctx context.Context, accountID int64, body interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/self_registration", accountID)

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to configure self registration: %w", err)
	}

	return result, nil
}

// AccountGetRolesPermissions retrieves permissions for all roles in an account
func (s *AccountMiscExtraService) AccountGetRolesPermissions(ctx context.Context, accountID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/roles/permissions", accountID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get roles permissions: %w", err)
	}

	return result, nil
}

// AccountCreateSharedBrandConfig creates a shared brand configuration for an account
func (s *AccountMiscExtraService) AccountCreateSharedBrandConfig(ctx context.Context, accountID int64, body interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/shared_brand_configs", accountID)

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to create shared brand config: %w", err)
	}

	return result, nil
}

// AccountUpdateSharedBrandConfig updates a shared brand configuration
func (s *AccountMiscExtraService) AccountUpdateSharedBrandConfig(ctx context.Context, accountID, id int64, body interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/shared_brand_configs/%d", accountID, id)

	var result map[string]interface{}
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to update shared brand config: %w", err)
	}

	return result, nil
}

// AccountGetRubricUsedLocations retrieves locations where a rubric is used
func (s *AccountMiscExtraService) AccountGetRubricUsedLocations(ctx context.Context, accountID, rubricID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/rubrics/%d/used_locations", accountID, rubricID)

	var result []interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get rubric used locations: %w", err)
	}

	return result, nil
}

// AccountUploadRubric uploads a rubric to an account
func (s *AccountMiscExtraService) AccountUploadRubric(ctx context.Context, accountID int64, body interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/rubrics/upload", accountID)

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to upload rubric: %w", err)
	}

	return result, nil
}

// AccountGetRubricUpload retrieves the status of a rubric upload
func (s *AccountMiscExtraService) AccountGetRubricUpload(ctx context.Context, accountID, id int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/rubrics/upload/%d", accountID, id)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get rubric upload: %w", err)
	}

	return result, nil
}

// AccountCreateFolder creates a folder in an account
func (s *AccountMiscExtraService) AccountCreateFolder(ctx context.Context, accountID int64, body interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/folders", accountID)

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	return result, nil
}
