package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsersService_ListUserContentMigrations(t *testing.T) {
	want := []UserContentMigration{{ID: 1, MigrationType: "course_copy_importer"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/content_migrations" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.ListUserContentMigrations(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != 1 {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestUsersService_CreateUserContentMigration(t *testing.T) {
	want := &UserContentMigration{ID: 5, MigrationType: "course_copy_importer"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/users/42/content_migrations" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.CreateUserContentMigration(context.Background(), 42, map[string]interface{}{"migration_type": "course_copy_importer"})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 5 {
		t.Errorf("got ID %d, want 5", got.ID)
	}
}

func TestUsersService_GetUserContentMigration(t *testing.T) {
	want := &UserContentMigration{ID: 3, WorkflowState: "completed"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/content_migrations/3" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetUserContentMigration(context.Background(), 42, 3)
	if err != nil {
		t.Fatal(err)
	}
	if got.WorkflowState != "completed" {
		t.Errorf("got state %q, want completed", got.WorkflowState)
	}
}

func TestUsersService_UpdateUserContentMigration(t *testing.T) {
	want := &UserContentMigration{ID: 3, WorkflowState: "running"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/users/42/content_migrations/3" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.UpdateUserContentMigration(context.Background(), 42, 3, map[string]interface{}{"workflow_state": "running"})
	if err != nil {
		t.Fatal(err)
	}
	if got.WorkflowState != "running" {
		t.Errorf("got state %q, want running", got.WorkflowState)
	}
}

func TestUsersService_GetUserContentMigrationSelectiveData(t *testing.T) {
	want := []interface{}{map[string]interface{}{"type": "course_settings"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/content_migrations/3/selective_data" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetUserContentMigrationSelectiveData(context.Background(), 42, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Errorf("got %d items, want 1", len(got))
	}
}

func TestUsersService_ListUserContentMigrationMigrators(t *testing.T) {
	want := []interface{}{map[string]interface{}{"type": "course_copy_importer"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/content_migrations/migrators" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.ListUserContentMigrationMigrators(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Errorf("got %d items, want 1", len(got))
	}
}

func TestUsersService_ListUserContentMigrationIssues(t *testing.T) {
	want := []UserMigrationIssue{{ID: 1, IssueType: "warning"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/content_migrations/3/migration_issues" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.ListUserContentMigrationIssues(context.Background(), 42, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].IssueType != "warning" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestUsersService_GetUserContentMigrationIssue(t *testing.T) {
	want := &UserMigrationIssue{ID: 7, Description: "some issue"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/content_migrations/3/migration_issues/7" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetUserContentMigrationIssue(context.Background(), 42, 3, 7)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 7 {
		t.Errorf("got ID %d, want 7", got.ID)
	}
}

func TestUsersService_UpdateUserContentMigrationIssue(t *testing.T) {
	want := &UserMigrationIssue{ID: 7, WorkflowState: "resolved"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/users/42/content_migrations/3/migration_issues/7" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.UpdateUserContentMigrationIssue(context.Background(), 42, 3, 7, map[string]interface{}{"workflow_state": "resolved"})
	if err != nil {
		t.Fatal(err)
	}
	if got.WorkflowState != "resolved" {
		t.Errorf("got state %q, want resolved", got.WorkflowState)
	}
}

func TestUsersService_ListUserContentLicenses(t *testing.T) {
	want := []ContentLicense{{ID: "cc_by", Name: "CC Attribution"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/content_licenses" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.ListUserContentLicenses(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "cc_by" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestUsersService_GetUserFile(t *testing.T) {
	want := &Attachment{ID: 99, DisplayName: "doc.pdf"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/files/99" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetUserFile(context.Background(), 42, 99)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 99 {
		t.Errorf("got ID %d, want 99", got.ID)
	}
}

func TestUsersService_ListUserFolders(t *testing.T) {
	want := []Folder{{ID: 1, Name: "Root"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/folders" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.ListUserFolders(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "Root" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestUsersService_CreateUserFolder(t *testing.T) {
	want := &Folder{ID: 10, Name: "New Folder"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/users/42/folders" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.CreateUserFolder(context.Background(), 42, map[string]interface{}{"name": "New Folder"})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 10 {
		t.Errorf("got ID %d, want 10", got.ID)
	}
}

func TestUsersService_GetUserFolder(t *testing.T) {
	want := &Folder{ID: 10, Name: "Subfolder"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/folders/10" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetUserFolder(context.Background(), 42, 10)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Subfolder" {
		t.Errorf("got name %q, want Subfolder", got.Name)
	}
}

func TestUsersService_GetUserFoldersByPath(t *testing.T) {
	want := []Folder{{ID: 1, Name: "my files"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/folders/by_path" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetUserFoldersByPath(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Errorf("got %d folders, want 1", len(got))
	}
}

func TestUsersService_SetUserUsageRights(t *testing.T) {
	want := &UsageRights{UseJustification: "own_copyright"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/users/42/usage_rights" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.SetUserUsageRights(context.Background(), 42, map[string]interface{}{"usage_rights[use_justification]": "own_copyright"})
	if err != nil {
		t.Fatal(err)
	}
	if got.UseJustification != "own_copyright" {
		t.Errorf("got %q, want own_copyright", got.UseJustification)
	}
}

func TestUsersService_DeleteUserUsageRights(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/users/42/usage_rights" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}")) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	if err := svc.DeleteUserUsageRights(context.Background(), 42); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("handler not called")
	}
}

func TestUsersService_CreatePandataEventsToken(t *testing.T) {
	want := &PandataEventsToken{Token: "abc123"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/users/self/pandata_events_token" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.CreatePandataEventsToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got.Token != "abc123" {
		t.Errorf("got token %q, want abc123", got.Token)
	}
}

func TestUsersService_CreateEducatorAccessibilityCourseScan(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/users/42/educator_accessibility_course_scan" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		called = true
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`)) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	_, err := svc.CreateEducatorAccessibilityCourseScan(context.Background(), 42, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("handler not called")
	}
}

func TestUsersService_GetEducatorAccessibilityCourseStatistics(t *testing.T) {
	want := &EducatorAccessibilityCourseStats{CourseID: 100, ErrorCount: 3}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/educator_accessibility_course_statistics" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetEducatorAccessibilityCourseStatistics(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if got.ErrorCount != 3 {
		t.Errorf("got ErrorCount %d, want 3", got.ErrorCount)
	}
}
