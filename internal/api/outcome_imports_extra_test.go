package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCourseOutcomeImportsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/outcome_imports" || r.Method != http.MethodPost {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 100, "workflow_state": "created"})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewCourseOutcomeImportsService(client)
	result, err := svc.Create(context.Background(), 1, map[string]interface{}{"import_type": "instructure_csv"})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if result.ID != 100 {
		t.Errorf("expected ID 100, got %d", result.ID)
	}
}

func TestCourseOutcomeImportsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/outcome_imports/100" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 100, "workflow_state": "succeeded"})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewCourseOutcomeImportsService(client)
	result, err := svc.Get(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if result.WorkflowState != "succeeded" {
		t.Errorf("expected succeeded, got %s", result.WorkflowState)
	}
}

func TestCourseOutcomeImportsService_GetCreatedGroupIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/outcome_imports/100/created_group_ids" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]int64{10, 20, 30})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewCourseOutcomeImportsService(client)
	ids, err := svc.GetCreatedGroupIDs(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("GetCreatedGroupIDs failed: %v", err)
	}
	if len(ids) != 3 {
		t.Errorf("expected 3 IDs, got %d", len(ids))
	}
}
