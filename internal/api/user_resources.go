package api

import (
	"context"
	"fmt"
)

// UserContentMigration is a lightweight migration struct for user-scoped content migrations.
// The full canvas migration type lives in content_migrations.go (course-scoped).
type UserContentMigration struct {
	ID                 int64       `json:"id"`
	MigrationType      string      `json:"migration_type"`
	MigrationTypeTitle string      `json:"migration_type_title"`
	WorkflowState      string      `json:"workflow_state"`
	CreatedAt          string      `json:"created_at"`
	UpdatedAt          string      `json:"updated_at"`
	FinishedAt         string      `json:"finished_at"`
	PreAttachment      interface{} `json:"pre_attachment,omitempty"`
}

// UserMigrationIssue represents an issue encountered during a content migration.
type UserMigrationIssue struct {
	ID                  int64       `json:"id"`
	ContentMigrationURL string      `json:"content_migration_url"`
	Description         string      `json:"description"`
	WorkflowState       string      `json:"workflow_state"`
	FixIssueURL         string      `json:"fix_issue_html_url"`
	IssueType           string      `json:"issue_type"`
	ErrorMessage        string      `json:"error_message"`
	ErrorReport         interface{} `json:"error_report,omitempty"`
}

// EducatorAccessibilityCourseStats holds accessibility statistics for a course scan.
type EducatorAccessibilityCourseStats struct {
	CourseID     int64  `json:"course_id"`
	ErrorCount   int    `json:"error_count"`
	FixableCount int    `json:"fixable_count"`
	CheckedAt    string `json:"checked_at"`
}

// PandataEventsToken holds a pandata events token for analytics.
type PandataEventsToken struct {
	Token        string `json:"token"`
	PropertiesID string `json:"props_token"`
	ExpiresAt    string `json:"expires_at"`
}

// ListUserContentMigrations retrieves content migrations for a user.
// GET /api/v1/users/:user_id/content_migrations
func (s *UsersService) ListUserContentMigrations(ctx context.Context, userID int64) ([]UserContentMigration, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_migrations", userID)
	var result []UserContentMigration
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing content migrations for user %d: %w", userID, err)
	}
	return result, nil
}

// CreateUserContentMigration creates a new content migration for a user.
// POST /api/v1/users/:user_id/content_migrations
func (s *UsersService) CreateUserContentMigration(ctx context.Context, userID int64, params map[string]interface{}) (*UserContentMigration, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_migrations", userID)
	var result UserContentMigration
	if err := s.client.PostJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("creating content migration for user %d: %w", userID, err)
	}
	return &result, nil
}

// GetUserContentMigration retrieves a single content migration for a user.
// GET /api/v1/users/:user_id/content_migrations/:id
func (s *UsersService) GetUserContentMigration(ctx context.Context, userID, migrationID int64) (*UserContentMigration, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_migrations/%d", userID, migrationID)
	var result UserContentMigration
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting content migration %d for user %d: %w", migrationID, userID, err)
	}
	return &result, nil
}

// UpdateUserContentMigration updates a content migration for a user.
// PUT /api/v1/users/:user_id/content_migrations/:id
func (s *UsersService) UpdateUserContentMigration(ctx context.Context, userID, migrationID int64, params map[string]interface{}) (*UserContentMigration, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_migrations/%d", userID, migrationID)
	var result UserContentMigration
	if err := s.client.PutJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("updating content migration %d for user %d: %w", migrationID, userID, err)
	}
	return &result, nil
}

// GetUserContentMigrationSelectiveData retrieves selective data for a content migration.
// GET /api/v1/users/:user_id/content_migrations/:id/selective_data
func (s *UsersService) GetUserContentMigrationSelectiveData(ctx context.Context, userID, migrationID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_migrations/%d/selective_data", userID, migrationID)
	var result []interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting selective data for migration %d of user %d: %w", migrationID, userID, err)
	}
	return result, nil
}

// ListUserContentMigrationMigrators retrieves available migrators for a user.
// GET /api/v1/users/:user_id/content_migrations/migrators
func (s *UsersService) ListUserContentMigrationMigrators(ctx context.Context, userID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_migrations/migrators", userID)
	var result []interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing migrators for user %d: %w", userID, err)
	}
	return result, nil
}

// ListUserContentMigrationIssues lists migration issues for a specific content migration.
// GET /api/v1/users/:user_id/content_migrations/:content_migration_id/migration_issues
func (s *UsersService) ListUserContentMigrationIssues(ctx context.Context, userID, migrationID int64) ([]UserMigrationIssue, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_migrations/%d/migration_issues", userID, migrationID)
	var result []UserMigrationIssue
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing migration issues for migration %d of user %d: %w", migrationID, userID, err)
	}
	return result, nil
}

// GetUserContentMigrationIssue retrieves a single migration issue.
// GET /api/v1/users/:user_id/content_migrations/:content_migration_id/migration_issues/:id
func (s *UsersService) GetUserContentMigrationIssue(ctx context.Context, userID, migrationID, issueID int64) (*UserMigrationIssue, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_migrations/%d/migration_issues/%d", userID, migrationID, issueID)
	var result UserMigrationIssue
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting migration issue %d for migration %d of user %d: %w", issueID, migrationID, userID, err)
	}
	return &result, nil
}

// UpdateUserContentMigrationIssue updates a migration issue.
// PUT /api/v1/users/:user_id/content_migrations/:content_migration_id/migration_issues/:id
func (s *UsersService) UpdateUserContentMigrationIssue(ctx context.Context, userID, migrationID, issueID int64, params map[string]interface{}) (*UserMigrationIssue, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_migrations/%d/migration_issues/%d", userID, migrationID, issueID)
	var result UserMigrationIssue
	if err := s.client.PutJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("updating migration issue %d for migration %d of user %d: %w", issueID, migrationID, userID, err)
	}
	return &result, nil
}

// ListUserContentLicenses retrieves content licenses available to a user.
// GET /api/v1/users/:user_id/content_licenses
func (s *UsersService) ListUserContentLicenses(ctx context.Context, userID int64) ([]ContentLicense, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_licenses", userID)
	var result []ContentLicense
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing content licenses for user %d: %w", userID, err)
	}
	return result, nil
}

// GetUserFile retrieves metadata for a file in a user's context.
// GET /api/v1/users/:user_id/files/:id
func (s *UsersService) GetUserFile(ctx context.Context, userID, fileID int64) (*Attachment, error) {
	path := fmt.Sprintf("/api/v1/users/%d/files/%d", userID, fileID)
	var result Attachment
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting file %d for user %d: %w", fileID, userID, err)
	}
	return &result, nil
}

// ListUserFolders retrieves folders for a user.
// GET /api/v1/users/:user_id/folders
func (s *UsersService) ListUserFolders(ctx context.Context, userID int64) ([]Folder, error) {
	path := fmt.Sprintf("/api/v1/users/%d/folders", userID)
	var result []Folder
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing folders for user %d: %w", userID, err)
	}
	return result, nil
}

// CreateUserFolder creates a new folder for a user.
// POST /api/v1/users/:user_id/folders
func (s *UsersService) CreateUserFolder(ctx context.Context, userID int64, params map[string]interface{}) (*Folder, error) {
	path := fmt.Sprintf("/api/v1/users/%d/folders", userID)
	var result Folder
	if err := s.client.PostJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("creating folder for user %d: %w", userID, err)
	}
	return &result, nil
}

// GetUserFolder retrieves a single folder for a user.
// GET /api/v1/users/:user_id/folders/:id
func (s *UsersService) GetUserFolder(ctx context.Context, userID, folderID int64) (*Folder, error) {
	path := fmt.Sprintf("/api/v1/users/%d/folders/%d", userID, folderID)
	var result Folder
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting folder %d for user %d: %w", folderID, userID, err)
	}
	return &result, nil
}

// GetUserFoldersByPath retrieves user folders by path (root level).
// GET /api/v1/users/:user_id/folders/by_path
func (s *UsersService) GetUserFoldersByPath(ctx context.Context, userID int64) ([]Folder, error) {
	path := fmt.Sprintf("/api/v1/users/%d/folders/by_path", userID)
	var result []Folder
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting folders by path for user %d: %w", userID, err)
	}
	return result, nil
}

// CreateEducatorAccessibilityCourseScan initiates an accessibility scan for a user.
// POST /api/v1/users/:user_id/educator_accessibility_course_scan
func (s *UsersService) CreateEducatorAccessibilityCourseScan(ctx context.Context, userID int64, params map[string]interface{}) (interface{}, error) {
	path := fmt.Sprintf("/api/v1/users/%d/educator_accessibility_course_scan", userID)
	var result interface{}
	if err := s.client.PostJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("creating educator accessibility scan for user %d: %w", userID, err)
	}
	return result, nil
}

// GetEducatorAccessibilityCourseStatistics retrieves accessibility statistics for a user.
// GET /api/v1/users/:user_id/educator_accessibility_course_statistics
func (s *UsersService) GetEducatorAccessibilityCourseStatistics(ctx context.Context, userID int64) (*EducatorAccessibilityCourseStats, error) {
	path := fmt.Sprintf("/api/v1/users/%d/educator_accessibility_course_statistics", userID)
	var result EducatorAccessibilityCourseStats
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting educator accessibility stats for user %d: %w", userID, err)
	}
	return &result, nil
}

// SetUserUsageRights sets usage rights on files for a user.
// PUT /api/v1/users/:user_id/usage_rights
func (s *UsersService) SetUserUsageRights(ctx context.Context, userID int64, params map[string]interface{}) (*UsageRights, error) {
	path := fmt.Sprintf("/api/v1/users/%d/usage_rights", userID)
	var result UsageRights
	if err := s.client.PutJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("setting usage rights for user %d: %w", userID, err)
	}
	return &result, nil
}

// DeleteUserUsageRights removes usage rights from files for a user.
// DELETE /api/v1/users/:user_id/usage_rights
func (s *UsersService) DeleteUserUsageRights(ctx context.Context, userID int64) error {
	path := fmt.Sprintf("/api/v1/users/%d/usage_rights", userID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// CreatePandataEventsToken generates a pandata events token for the current user.
// POST /api/v1/users/self/pandata_events_token
func (s *UsersService) CreatePandataEventsToken(ctx context.Context) (*PandataEventsToken, error) {
	path := "/api/v1/users/self/pandata_events_token"
	var result PandataEventsToken
	if err := s.client.PostJSON(ctx, path, nil, &result); err != nil {
		return nil, fmt.Errorf("creating pandata events token: %w", err)
	}
	return &result, nil
}
