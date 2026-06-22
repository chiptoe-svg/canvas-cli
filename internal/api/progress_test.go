package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProgressService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/progress/42" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		p := ProgressJob{ID: 42, WorkflowState: "running", Completion: 50.0}
		if err := json.NewEncoder(w).Encode(p); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewProgressService(client)

	p, err := svc.Get(context.Background(), 42)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if p.ID != 42 {
		t.Errorf("expected ID 42, got %d", p.ID)
	}
	if p.WorkflowState != "running" {
		t.Errorf("expected state running, got %s", p.WorkflowState)
	}
}

func TestProgressService_Cancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/progress/42/cancel" || r.Method != http.MethodPost {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		p := ProgressJob{ID: 42, WorkflowState: "failed", Message: "cancelled"}
		if err := json.NewEncoder(w).Encode(p); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewProgressService(client)

	p, err := svc.Cancel(context.Background(), 42)
	if err != nil {
		t.Fatalf("Cancel failed: %v", err)
	}
	if p.WorkflowState != "failed" {
		t.Errorf("expected state failed, got %s", p.WorkflowState)
	}
}
