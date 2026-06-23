package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsersService_GetGradedSubmissions(t *testing.T) {
	want := []GradedSubmission{{ID: 1, AssignmentID: 10, Score: 95.5}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/graded_submissions" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetGradedSubmissions(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Score != 95.5 {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestUsersService_DeleteSessions(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/users/7/sessions" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}")) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	if err := svc.DeleteSessions(context.Background(), 7); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("handler not called")
	}
}

func TestUsersService_DeleteMobileSessions(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/users/7/mobile_sessions" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}")) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	if err := svc.DeleteMobileSessions(context.Background(), 7); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("handler not called")
	}
}

func TestUsersService_DeleteAllMobileSessions(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/users/mobile_sessions" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}")) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	if err := svc.DeleteAllMobileSessions(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("handler not called")
	}
}

func TestUsersService_GetActivityStreamAll(t *testing.T) {
	want := []ActivityStreamItem{{ID: 1, Type: "Message", Title: "Hello"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/activity_stream" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetActivityStreamAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d items, want 1", len(got))
	}
}

func TestUsersService_ListEportfolios(t *testing.T) {
	want := []Eportfolio{{ID: 1, Name: "My Portfolio", Public: true}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/eportfolios" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.ListEportfolios(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "My Portfolio" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestUsersService_GetTemporaryEnrollmentStatus(t *testing.T) {
	want := &TemporaryEnrollmentStatus{IsProvider: true}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/temporary_enrollment_status" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetTemporaryEnrollmentStatus(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if !got.IsProvider {
		t.Error("expected IsProvider=true")
	}
}

func TestUsersService_CreatePageViewQuery(t *testing.T) {
	want := &PageViewQuery{ID: "query-123", Status: "pending"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/users/42/page_views/query" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.CreatePageViewQuery(context.Background(), 42, map[string]interface{}{"start_time": "2024-01-01"})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "query-123" {
		t.Errorf("got ID %q, want %q", got.ID, "query-123")
	}
}

func TestUsersService_GetPageViewQueryStatus(t *testing.T) {
	want := &PageViewQuery{ID: "query-123", Status: "completed"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/page_views/query/query-123" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetPageViewQueryStatus(context.Background(), 42, "query-123")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "completed" {
		t.Errorf("got Status %q, want completed", got.Status)
	}
}

func TestUsersService_GetPageViewQueryResults(t *testing.T) {
	want := &PageViewQueryResults{Results: []PageView{{ID: "pv1", URL: "/courses/1"}}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/page_views/query/query-123/results" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetPageViewQueryResults(context.Background(), 42, "query-123")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Results) != 1 {
		t.Errorf("got %d results, want 1", len(got.Results))
	}
}

func TestUsersService_CreateSelfPageViewQuery(t *testing.T) {
	want := &PageViewQuery{ID: "self-query-1", Status: "pending"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/users/page_views/query" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.CreateSelfPageViewQuery(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "self-query-1" {
		t.Errorf("got ID %q, want self-query-1", got.ID)
	}
}

func TestUsersService_GetSelfPageViewQueryStatus(t *testing.T) {
	want := &PageViewQuery{ID: "sq1", Status: "pending"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/page_views/query/sq1" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetSelfPageViewQueryStatus(context.Background(), "sq1")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "sq1" {
		t.Errorf("got ID %q, want sq1", got.ID)
	}
}

func TestUsersService_GetSelfPageViewQueryResults(t *testing.T) {
	want := &PageViewQueryResults{Results: []PageView{{ID: "pv2"}}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/page_views/query/sq1/results" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetSelfPageViewQueryResults(context.Background(), "sq1")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Results) != 1 {
		t.Errorf("got %d results, want 1", len(got.Results))
	}
}

func TestUsersService_ListContentExports(t *testing.T) {
	want := []ContentExport{{ID: 1, ExportType: "common_cartridge"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/content_exports" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.ListContentExports(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ExportType != "common_cartridge" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestUsersService_CreateContentExport(t *testing.T) {
	want := &ContentExport{ID: 10, ExportType: "zip", WorkflowState: "created"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/users/42/content_exports" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.CreateContentExport(context.Background(), 42, ContentExportParams{ExportType: "zip"})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 10 {
		t.Errorf("got ID %d, want 10", got.ID)
	}
}

func TestUsersService_GetContentExport(t *testing.T) {
	want := &ContentExport{ID: 10, ExportType: "zip"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/content_exports/10" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetContentExport(context.Background(), 42, 10)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 10 {
		t.Errorf("got ID %d, want 10", got.ID)
	}
}

func TestUsersService_SetFilesUIVersionPreference(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/users/42/files_ui_version_preference" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}")) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	if err := svc.SetFilesUIVersionPreference(context.Background(), 42, "v2"); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("handler not called")
	}
}

func TestUsersService_SetTextEditorPreference(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/users/42/text_editor_preference" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}")) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	if err := svc.SetTextEditorPreference(context.Background(), 42, "rce"); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("handler not called")
	}
}

func TestUsersService_GetPlannerItems(t *testing.T) {
	want := []interface{}{map[string]interface{}{"type": "assignment"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/users/42/planner/items" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	got, err := svc.GetPlannerItems(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Errorf("got %d items, want 1", len(got))
	}
}

func TestUsersService_ResetPassword(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/users/reset_password" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}")) // #nosec G104
	}))
	defer server.Close()
	svc := NewUsersService(newTestClient(t, server.URL))
	if err := svc.ResetPassword(context.Background(), "user@example.com"); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("handler not called")
	}
}
