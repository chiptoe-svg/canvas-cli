package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContentMigrationsService_ListForAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/10/content_migrations" {
			t.Errorf("Expected path /api/v1/accounts/10/content_migrations, got %s", r.URL.Path)
		}

		migrations := []ContentMigration{
			{ID: 1, MigrationType: "course_copy_importer", WorkflowState: "completed"},
			{ID: 2, MigrationType: "common_cartridge_importer", WorkflowState: "running"},
		}
		json.NewEncoder(w).Encode(migrations)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewContentMigrationsService(client)

	migrations, err := service.ListForAccount(context.Background(), 10)
	if err != nil {
		t.Fatalf("ListForAccount failed: %v", err)
	}

	if len(migrations) != 2 {
		t.Errorf("Expected 2 migrations, got %d", len(migrations))
	}

	if migrations[0].ID != 1 {
		t.Errorf("Expected migration ID 1, got %d", migrations[0].ID)
	}
}

func TestContentMigrationsService_GetForAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/10/content_migrations/42" {
			t.Errorf("Expected path /api/v1/accounts/10/content_migrations/42, got %s", r.URL.Path)
		}

		migration := ContentMigration{
			ID:            42,
			MigrationType: "course_copy_importer",
			WorkflowState: "completed",
		}
		json.NewEncoder(w).Encode(migration)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewContentMigrationsService(client)

	migration, err := service.GetForAccount(context.Background(), 10, 42)
	if err != nil {
		t.Fatalf("GetForAccount failed: %v", err)
	}

	if migration.ID != 42 {
		t.Errorf("Expected ID 42, got %d", migration.ID)
	}

	if migration.WorkflowState != "completed" {
		t.Errorf("Expected state 'completed', got '%s'", migration.WorkflowState)
	}
}

func TestContentMigrationsService_CreateForAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/accounts/10/content_migrations" {
			t.Errorf("Expected path /api/v1/accounts/10/content_migrations, got %s", r.URL.Path)
		}

		migration := ContentMigration{
			ID:            99,
			MigrationType: "course_copy_importer",
			WorkflowState: "queued",
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(migration)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewContentMigrationsService(client)

	params := &CreateContentMigrationParams{
		MigrationType: "course_copy_importer",
	}

	migration, err := service.CreateForAccount(context.Background(), 10, params)
	if err != nil {
		t.Fatalf("CreateForAccount failed: %v", err)
	}

	if migration.ID != 99 {
		t.Errorf("Expected ID 99, got %d", migration.ID)
	}

	if migration.MigrationType != "course_copy_importer" {
		t.Errorf("Expected type 'course_copy_importer', got '%s'", migration.MigrationType)
	}
}

func TestContentMigrationsService_UpdateForAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/accounts/10/content_migrations/42" {
			t.Errorf("Expected path /api/v1/accounts/10/content_migrations/42, got %s", r.URL.Path)
		}

		migration := ContentMigration{
			ID:            42,
			WorkflowState: "failed",
		}
		json.NewEncoder(w).Encode(migration)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewContentMigrationsService(client)

	params := &UpdateContentMigrationParams{WorkflowState: "failed"}

	migration, err := service.UpdateForAccount(context.Background(), 10, 42, params)
	if err != nil {
		t.Fatalf("UpdateForAccount failed: %v", err)
	}

	if migration.WorkflowState != "failed" {
		t.Errorf("Expected state 'failed', got '%s'", migration.WorkflowState)
	}
}

func TestContentMigrationsService_GetMigratorsForAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/10/content_migrations/migrators" {
			t.Errorf("Expected path /api/v1/accounts/10/content_migrations/migrators, got %s", r.URL.Path)
		}

		migrators := []Migrator{
			{Type: "course_copy_importer", Name: "Course Copy", RequiresFileUpload: false},
			{Type: "common_cartridge_importer", Name: "Common Cartridge", RequiresFileUpload: true},
		}
		json.NewEncoder(w).Encode(migrators)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewContentMigrationsService(client)

	migrators, err := service.GetMigratorsForAccount(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetMigratorsForAccount failed: %v", err)
	}

	if len(migrators) != 2 {
		t.Errorf("Expected 2 migrators, got %d", len(migrators))
	}

	if migrators[0].Type != "course_copy_importer" {
		t.Errorf("Expected type 'course_copy_importer', got '%s'", migrators[0].Type)
	}
}

func TestContentMigrationsService_GetMigrationIssuesForAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/10/content_migrations/42/migration_issues" {
			t.Errorf("Expected path /api/v1/accounts/10/content_migrations/42/migration_issues, got %s", r.URL.Path)
		}

		issues := []MigrationIssue{
			{ID: 1, IssueType: "warning", WorkflowState: "active", Description: "Missing asset"},
		}
		json.NewEncoder(w).Encode(issues)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewContentMigrationsService(client)

	issues, err := service.GetMigrationIssuesForAccount(context.Background(), 10, 42)
	if err != nil {
		t.Fatalf("GetMigrationIssuesForAccount failed: %v", err)
	}

	if len(issues) != 1 {
		t.Fatalf("Expected 1 issue, got %d", len(issues))
	}

	if issues[0].IssueType != "warning" {
		t.Errorf("Expected issue type 'warning', got '%s'", issues[0].IssueType)
	}
}

func TestContentMigrationsService_UpdateMigrationIssueForAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/accounts/10/content_migrations/42/migration_issues/7" {
			t.Errorf("Expected path /api/v1/accounts/10/content_migrations/42/migration_issues/7, got %s", r.URL.Path)
		}

		issue := MigrationIssue{
			ID:            7,
			WorkflowState: "resolved",
			IssueType:     "warning",
		}
		json.NewEncoder(w).Encode(issue)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewContentMigrationsService(client)

	params := &MigrationIssueUpdateParams{WorkflowState: "resolved"}

	issue, err := service.UpdateMigrationIssueForAccount(context.Background(), 10, 42, 7, params)
	if err != nil {
		t.Fatalf("UpdateMigrationIssueForAccount failed: %v", err)
	}

	if issue.WorkflowState != "resolved" {
		t.Errorf("Expected state 'resolved', got '%s'", issue.WorkflowState)
	}
}
