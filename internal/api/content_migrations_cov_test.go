package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContentMigrationsService_GetMigrationIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/content_migrations/5/migration_issues/3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MigrationIssue{ID: 3, WorkflowState: "active", Description: "Missing resource"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentMigrationsService(client)
	issue, err := svc.GetMigrationIssue(context.Background(), 10, 5, 3)
	if err != nil {
		t.Fatalf("GetMigrationIssue: %v", err)
	}
	if issue.ID != 3 {
		t.Errorf("expected ID 3, got %d", issue.ID)
	}
	if issue.Description != "Missing resource" {
		t.Errorf("expected description='Missing resource', got %s", issue.Description)
	}
}

func TestContentMigrationsService_GetMigrationIssue_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentMigrationsService(client)
	_, err := svc.GetMigrationIssue(context.Background(), 10, 5, 3)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestContentMigrationsService_UpdateMigrationIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/content_migrations/5/migration_issues/3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var params MigrationIssueUpdateParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if params.WorkflowState != "resolved" {
			t.Errorf("expected workflow_state=resolved, got %s", params.WorkflowState)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MigrationIssue{ID: 3, WorkflowState: "resolved"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentMigrationsService(client)
	issue, err := svc.UpdateMigrationIssue(context.Background(), 10, 5, 3, &MigrationIssueUpdateParams{WorkflowState: "resolved"})
	if err != nil {
		t.Fatalf("UpdateMigrationIssue: %v", err)
	}
	if issue.WorkflowState != "resolved" {
		t.Errorf("expected resolved, got %s", issue.WorkflowState)
	}
}

func TestContentMigrationsService_GetAssetIDMapping(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/content_migrations/5/asset_id_mapping" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"Assignment_1": float64(2), "Quiz_3": float64(4)})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentMigrationsService(client)
	mapping, err := svc.GetAssetIDMapping(context.Background(), 10, 5)
	if err != nil {
		t.Fatalf("GetAssetIDMapping: %v", err)
	}
	if mapping["Assignment_1"] != float64(2) {
		t.Errorf("expected Assignment_1=2, got %v", mapping["Assignment_1"])
	}
}

func TestContentMigrationsService_ListMigrationIssuesForGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/7/content_migrations/5/migration_issues" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]MigrationIssue{{ID: 1, WorkflowState: "active"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentMigrationsService(client)
	issues, err := svc.ListMigrationIssuesForGroup(context.Background(), 7, 5)
	if err != nil {
		t.Fatalf("ListMigrationIssuesForGroup: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestContentMigrationsService_GetMigrationIssueForGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/7/content_migrations/5/migration_issues/3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MigrationIssue{ID: 3, WorkflowState: "active"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentMigrationsService(client)
	issue, err := svc.GetMigrationIssueForGroup(context.Background(), 7, 5, 3)
	if err != nil {
		t.Fatalf("GetMigrationIssueForGroup: %v", err)
	}
	if issue.ID != 3 {
		t.Errorf("expected ID 3, got %d", issue.ID)
	}
}

func TestContentMigrationsService_UpdateMigrationIssueForGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/7/content_migrations/5/migration_issues/3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MigrationIssue{ID: 3, WorkflowState: "resolved"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentMigrationsService(client)
	issue, err := svc.UpdateMigrationIssueForGroup(context.Background(), 7, 5, 3, &MigrationIssueUpdateParams{WorkflowState: "resolved"})
	if err != nil {
		t.Fatalf("UpdateMigrationIssueForGroup: %v", err)
	}
	if issue.WorkflowState != "resolved" {
		t.Errorf("expected resolved, got %s", issue.WorkflowState)
	}
}

func TestContentMigrationsService_ListForGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/7/content_migrations" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]ContentMigration{{ID: 1, MigrationType: "course_copy_importer"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentMigrationsService(client)
	migrations, err := svc.ListForGroup(context.Background(), 7)
	if err != nil {
		t.Fatalf("ListForGroup: %v", err)
	}
	if len(migrations) != 1 {
		t.Fatalf("expected 1 migration, got %d", len(migrations))
	}
	if migrations[0].MigrationType != "course_copy_importer" {
		t.Errorf("expected course_copy_importer, got %s", migrations[0].MigrationType)
	}
}

func TestContentMigrationsService_CreateForGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/7/content_migrations" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["migration_type"] != "course_copy_importer" {
			t.Errorf("expected migration_type=course_copy_importer, got %v", body["migration_type"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ContentMigration{ID: 50, MigrationType: "course_copy_importer"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentMigrationsService(client)
	migration, err := svc.CreateForGroup(context.Background(), 7, &CreateContentMigrationParams{MigrationType: "course_copy_importer"})
	if err != nil {
		t.Fatalf("CreateForGroup: %v", err)
	}
	if migration.ID != 50 {
		t.Errorf("expected ID 50, got %d", migration.ID)
	}
}

func TestContentMigrationsService_GetForGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/7/content_migrations/50" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ContentMigration{ID: 50, MigrationType: "course_copy_importer", WorkflowState: "running"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentMigrationsService(client)
	migration, err := svc.GetForGroup(context.Background(), 7, 50)
	if err != nil {
		t.Fatalf("GetForGroup: %v", err)
	}
	if migration.WorkflowState != "running" {
		t.Errorf("expected running, got %s", migration.WorkflowState)
	}
}

func TestContentMigrationsService_UpdateForGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/7/content_migrations/50" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ContentMigration{ID: 50, WorkflowState: "completed"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentMigrationsService(client)
	migration, err := svc.UpdateForGroup(context.Background(), 7, 50, &UpdateContentMigrationParams{})
	if err != nil {
		t.Fatalf("UpdateForGroup: %v", err)
	}
	if migration.WorkflowState != "completed" {
		t.Errorf("expected completed, got %s", migration.WorkflowState)
	}
}

func TestContentMigrationsService_GetSelectiveDataForGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/7/content_migrations/50/selective_data" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		// Assert type query param is forwarded when provided.
		if got := r.URL.Query().Get("type"); got != "assignments" {
			t.Errorf("expected type=assignments, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]ContentListItem{{Type: "assignment", Property: "copy[all_assignments]", Title: "HW1"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentMigrationsService(client)
	data, err := svc.GetSelectiveDataForGroup(context.Background(), 7, 50, "assignments")
	if err != nil {
		t.Fatalf("GetSelectiveDataForGroup: %v", err)
	}
	if len(data) != 1 {
		t.Fatalf("expected 1 item, got %d", len(data))
	}
	// Assert typed return: Title field should be decoded correctly.
	if data[0].Title != "HW1" {
		t.Errorf("expected Title HW1, got %q", data[0].Title)
	}
}

func TestContentMigrationsService_GetMigratorsForGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/7/content_migrations/migrators" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Migrator{{Type: "course_copy_importer", Name: "Course Copy"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentMigrationsService(client)
	migrators, err := svc.GetMigratorsForGroup(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetMigratorsForGroup: %v", err)
	}
	if len(migrators) != 1 {
		t.Fatalf("expected 1 migrator, got %d", len(migrators))
	}
	if migrators[0].Type != "course_copy_importer" {
		t.Errorf("expected course_copy_importer, got %s", migrators[0].Type)
	}
}
