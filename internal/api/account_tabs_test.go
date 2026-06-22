package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountTabsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/tabs" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]AccountTab{
			{ID: "courses", Label: "Courses", Type: "internal", Position: 1},
			{ID: "users", Label: "People", Type: "internal", Position: 2},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountTabsService(client)

	tabs, err := svc.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tabs) != 2 {
		t.Errorf("expected 2 tabs, got %d", len(tabs))
	}
	if tabs[0].ID != "courses" {
		t.Errorf("expected courses, got %s", tabs[0].ID)
	}
}

func TestAccountTabsService_List_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"not found"}]}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountTabsService(client)

	_, err := svc.List(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
