package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEpubExportsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/epub_exports" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		exports := []EpubExport{
			{ID: 1, WorkflowState: "generated"},
		}
		if err := json.NewEncoder(w).Encode(exports); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewEpubExportsService(client)

	exports, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(exports) != 1 {
		t.Errorf("expected 1 export, got %d", len(exports))
	}
}

func TestEpubExportsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/42/epub_exports" || r.Method != http.MethodPost {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		export := EpubExport{ID: 5, WorkflowState: "created"}
		if err := json.NewEncoder(w).Encode(export); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewEpubExportsService(client)

	export, err := svc.Create(context.Background(), 42)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if export.ID != 5 {
		t.Errorf("expected ID 5, got %d", export.ID)
	}
}

func TestEpubExportsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/42/epub_exports/5" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		export := EpubExport{ID: 5, WorkflowState: "generated"}
		if err := json.NewEncoder(w).Encode(export); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewEpubExportsService(client)

	export, err := svc.Get(context.Background(), 42, 5)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if export.WorkflowState != "generated" {
		t.Errorf("unexpected state: %s", export.WorkflowState)
	}
}
