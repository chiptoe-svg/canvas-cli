package api

import (
	"context"
	"fmt"
)

// MigrationIssueUpdateParams holds parameters for updating a migration issue
type MigrationIssueUpdateParams struct {
	WorkflowState string `json:"workflow_state"`
}

// ListForAccount retrieves content migrations for an account
func (s *ContentMigrationsService) ListForAccount(ctx context.Context, accountID int64) ([]ContentMigration, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/content_migrations", accountID)

	var migrations []ContentMigration
	if err := s.client.GetAllPages(ctx, path, &migrations); err != nil {
		return nil, fmt.Errorf("failed to list account content migrations: %w", err)
	}

	return migrations, nil
}

// GetForAccount retrieves a single content migration for an account
func (s *ContentMigrationsService) GetForAccount(ctx context.Context, accountID, migrationID int64) (*ContentMigration, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/content_migrations/%d", accountID, migrationID)

	var migration ContentMigration
	if err := s.client.GetJSON(ctx, path, &migration); err != nil {
		return nil, fmt.Errorf("failed to get account content migration: %w", err)
	}

	return &migration, nil
}

// CreateForAccount creates a new content migration for an account
func (s *ContentMigrationsService) CreateForAccount(ctx context.Context, accountID int64, params *CreateContentMigrationParams) (*ContentMigration, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/content_migrations", accountID)

	body := make(map[string]interface{})
	body["migration_type"] = params.MigrationType

	settings := make(map[string]interface{})

	if params.SourceCourseID != nil {
		settings["source_course_id"] = *params.SourceCourseID
	}

	if params.FileURL != "" {
		settings["file_url"] = params.FileURL
	}

	if params.ContentExportID != nil {
		settings["content_export_id"] = *params.ContentExportID
	}

	if len(settings) > 0 {
		body["settings"] = settings
	}

	if params.SelectiveImport != nil {
		body["selective_import"] = *params.SelectiveImport
	}

	var migration ContentMigration
	if err := s.client.PostJSON(ctx, path, body, &migration); err != nil {
		return nil, fmt.Errorf("failed to create account content migration: %w", err)
	}

	return &migration, nil
}

// UpdateForAccount updates a content migration for an account
func (s *ContentMigrationsService) UpdateForAccount(ctx context.Context, accountID, migrationID int64, params *UpdateContentMigrationParams) (*ContentMigration, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/content_migrations/%d", accountID, migrationID)

	body := make(map[string]interface{})
	if params.WorkflowState != "" {
		body["workflow_state"] = params.WorkflowState
	}

	var migration ContentMigration
	if err := s.client.PutJSON(ctx, path, body, &migration); err != nil {
		return nil, fmt.Errorf("failed to update account content migration: %w", err)
	}

	return &migration, nil
}

// GetSelectiveDataForAccount retrieves selective data for a content migration in an account
func (s *ContentMigrationsService) GetSelectiveDataForAccount(ctx context.Context, accountID, migrationID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/content_migrations/%d/selective_data", accountID, migrationID)

	var items []interface{}
	if err := s.client.GetAllPages(ctx, path, &items); err != nil {
		return nil, fmt.Errorf("failed to get selective data: %w", err)
	}

	return items, nil
}

// GetMigratorsForAccount retrieves available migrator types for an account
func (s *ContentMigrationsService) GetMigratorsForAccount(ctx context.Context, accountID int64) ([]Migrator, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/content_migrations/migrators", accountID)

	var migrators []Migrator
	if err := s.client.GetJSON(ctx, path, &migrators); err != nil {
		return nil, fmt.Errorf("failed to get account migrators: %w", err)
	}

	return migrators, nil
}

// GetMigrationIssuesForAccount retrieves issues for a content migration in an account
func (s *ContentMigrationsService) GetMigrationIssuesForAccount(ctx context.Context, accountID, migrationID int64) ([]MigrationIssue, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/content_migrations/%d/migration_issues", accountID, migrationID)

	var issues []MigrationIssue
	if err := s.client.GetAllPages(ctx, path, &issues); err != nil {
		return nil, fmt.Errorf("failed to get account migration issues: %w", err)
	}

	return issues, nil
}

// GetMigrationIssueForAccount retrieves a single migration issue for an account
func (s *ContentMigrationsService) GetMigrationIssueForAccount(ctx context.Context, accountID, migrationID, issueID int64) (*MigrationIssue, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/content_migrations/%d/migration_issues/%d", accountID, migrationID, issueID)

	var issue MigrationIssue
	if err := s.client.GetJSON(ctx, path, &issue); err != nil {
		return nil, fmt.Errorf("failed to get account migration issue: %w", err)
	}

	return &issue, nil
}

// UpdateMigrationIssueForAccount updates a migration issue for an account
func (s *ContentMigrationsService) UpdateMigrationIssueForAccount(ctx context.Context, accountID, migrationID, issueID int64, params *MigrationIssueUpdateParams) (*MigrationIssue, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/content_migrations/%d/migration_issues/%d", accountID, migrationID, issueID)

	body := map[string]interface{}{
		"workflow_state": params.WorkflowState,
	}

	var issue MigrationIssue
	if err := s.client.PutJSON(ctx, path, body, &issue); err != nil {
		return nil, fmt.Errorf("failed to update account migration issue: %w", err)
	}

	return &issue, nil
}
