package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestModulesService_ListAssignmentOverrides(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/modules/2/assignment_overrides" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"id": 10, "title": "Override"}})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewModulesService(client)
	overrides, err := svc.ListAssignmentOverrides(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("ListAssignmentOverrides failed: %v", err)
	}
	if len(overrides) != 1 {
		t.Errorf("expected 1 override, got %d", len(overrides))
	}
}

func TestModulesService_UpdateAssignmentOverrides(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/modules/2/assignment_overrides" || r.Method != http.MethodPut {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"id": 10, "title": "Updated Override"}})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewModulesService(client)
	overrides, err := svc.UpdateAssignmentOverrides(context.Background(), 1, 2, []ModuleOverrideParams{{CourseSectionID: 5}})
	if err != nil {
		t.Fatalf("UpdateAssignmentOverrides failed: %v", err)
	}
	if len(overrides) != 1 {
		t.Errorf("expected 1 override, got %d", len(overrides))
	}
}
