package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewAccountOutcomeImportsService(t *testing.T) {
	client := newTestClient(t, "http://localhost")
	svc := NewAccountOutcomeImportsService(client)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
	if svc.client != client {
		t.Error("expected client to be set")
	}
}

func TestAccountOutcomeImportsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_imports" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OutcomeImport{ID: 42, WorkflowState: "created", Progress: 0.0})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountOutcomeImportsService(client)
	result, err := svc.Create(context.Background(), 1, map[string]interface{}{"import_type": "instructure_csv"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if result.ID != 42 {
		t.Errorf("expected ID 42, got %d", result.ID)
	}
	if result.WorkflowState != "created" {
		t.Errorf("expected workflow_state=created, got %s", result.WorkflowState)
	}
}

func TestAccountOutcomeImportsService_Create_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountOutcomeImportsService(client)
	_, err := svc.Create(context.Background(), 1, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountOutcomeImportsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_imports/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OutcomeImport{ID: 42, WorkflowState: "succeeded", Progress: 100.0})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountOutcomeImportsService(client)
	result, err := svc.Get(context.Background(), 1, 42)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if result.WorkflowState != "succeeded" {
		t.Errorf("expected succeeded, got %s", result.WorkflowState)
	}
	if result.Progress != 100.0 {
		t.Errorf("expected progress=100, got %f", result.Progress)
	}
}

func TestAccountOutcomeImportsService_Get_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountOutcomeImportsService(client)
	_, err := svc.Get(context.Background(), 1, 42)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountOutcomeImportsService_GetCreatedGroupIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_imports/42/created_group_ids" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]int64{100, 200, 300})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountOutcomeImportsService(client)
	ids, err := svc.GetCreatedGroupIDs(context.Background(), 1, 42)
	if err != nil {
		t.Fatalf("GetCreatedGroupIDs: %v", err)
	}
	if len(ids) != 3 {
		t.Fatalf("expected 3 IDs, got %d", len(ids))
	}
	if ids[0] != 100 {
		t.Errorf("expected first ID=100, got %d", ids[0])
	}
}

func TestAccountOutcomeImportsService_GetOutcomeGroupLinks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_group_links" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{map[string]interface{}{"id": float64(1)}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountOutcomeImportsService(client)
	links, err := svc.GetOutcomeGroupLinks(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetOutcomeGroupLinks: %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(links))
	}
}

func TestAccountOutcomeImportsService_GetRootOutcomeGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/root_outcome_group" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": float64(1), "title": "Root"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountOutcomeImportsService(client)
	result, err := svc.GetRootOutcomeGroup(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetRootOutcomeGroup: %v", err)
	}
	if result["title"] != "Root" {
		t.Errorf("expected title=Root, got %v", result["title"])
	}
}

func TestAccountOutcomeImportsService_GetOutcomeGroupSubgroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_groups/5/subgroups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{map[string]interface{}{"id": float64(10)}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountOutcomeImportsService(client)
	result, err := svc.GetOutcomeGroupSubgroups(context.Background(), 1, 5)
	if err != nil {
		t.Fatalf("GetOutcomeGroupSubgroups: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 subgroup, got %d", len(result))
	}
}

func TestAccountOutcomeImportsService_CreateOutcomeGroupSubgroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_groups/5/subgroups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": float64(20), "title": "New Subgroup"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountOutcomeImportsService(client)
	result, err := svc.CreateOutcomeGroupSubgroup(context.Background(), 1, 5, map[string]interface{}{"title": "New Subgroup"})
	if err != nil {
		t.Fatalf("CreateOutcomeGroupSubgroup: %v", err)
	}
	if result["title"] != "New Subgroup" {
		t.Errorf("expected title=New Subgroup, got %v", result["title"])
	}
}

func TestAccountOutcomeImportsService_ImportOutcomeGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_groups/5/import" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": float64(5), "imported": true})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountOutcomeImportsService(client)
	result, err := svc.ImportOutcomeGroup(context.Background(), 1, 5, map[string]interface{}{"source_outcome_group_id": 99})
	if err != nil {
		t.Fatalf("ImportOutcomeGroup: %v", err)
	}
	if result["imported"] != true {
		t.Errorf("expected imported=true, got %v", result["imported"])
	}
}

func TestAccountOutcomeImportsService_GetOutcomeProficiency(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_proficiency" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"ratings": []interface{}{}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountOutcomeImportsService(client)
	result, err := svc.GetOutcomeProficiency(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetOutcomeProficiency: %v", err)
	}
	if _, ok := result["ratings"]; !ok {
		t.Error("expected ratings key in result")
	}
}

func TestAccountOutcomeImportsService_UpdateOutcomeProficiency(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_proficiency" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"ratings": []interface{}{map[string]interface{}{"description": "Excellent"}}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountOutcomeImportsService(client)
	result, err := svc.UpdateOutcomeProficiency(context.Background(), 1, map[string]interface{}{"ratings": []interface{}{}})
	if err != nil {
		t.Fatalf("UpdateOutcomeProficiency: %v", err)
	}
	if _, ok := result["ratings"]; !ok {
		t.Error("expected ratings key in result")
	}
}
