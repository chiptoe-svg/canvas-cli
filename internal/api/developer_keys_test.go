package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeveloperKeysService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/developer_keys" {
			t.Errorf("expected path /api/v1/accounts/1/developer_keys, got %s", r.URL.Path)
		}

		keys := []DeveloperKey{
			{ID: 1, Name: "My App", Email: "dev@example.com"},
			{ID: 2, Name: "Other App", Email: "other@example.com"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(keys)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewDeveloperKeysService(client)

	keys, err := service.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(keys) != 2 {
		t.Fatalf("expected 2 developer keys, got %d", len(keys))
	}

	if keys[0].Name != "My App" {
		t.Errorf("expected name 'My App', got '%s'", keys[0].Name)
	}
}

func TestDeveloperKeysService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/developer_keys" {
			t.Errorf("expected path /api/v1/accounts/1/developer_keys, got %s", r.URL.Path)
		}

		key := DeveloperKey{ID: 10, Name: "New Key", Email: "new@example.com"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(key)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewDeveloperKeysService(client)

	params := &DeveloperKeyParams{DeveloperKey: DeveloperKeyFields{Name: "New Key", Email: "new@example.com"}}
	key, err := service.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if key.ID != 10 {
		t.Errorf("expected ID 10, got %d", key.ID)
	}
}

func TestDeveloperKeysService_CreateBinding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		expected := "/api/v1/accounts/1/developer_keys/10/developer_key_account_bindings"
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != expected {
			t.Errorf("expected path %s, got %s", expected, r.URL.Path)
		}

		binding := DeveloperKeyBinding{ID: 5, AccountID: 1, DeveloperKeyID: 10, WorkflowState: "on"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(binding)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewDeveloperKeysService(client)

	params := &DeveloperKeyBindingParams{DeveloperKeyAccountBinding: DeveloperKeyBindingFields{WorkflowState: "on"}}
	binding, err := service.CreateBinding(context.Background(), 1, 10, params)
	if err != nil {
		t.Fatalf("CreateBinding failed: %v", err)
	}

	if binding.WorkflowState != "on" {
		t.Errorf("expected workflow_state 'on', got '%s'", binding.WorkflowState)
	}
}
