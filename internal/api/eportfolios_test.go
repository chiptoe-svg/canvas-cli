package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEportfoliosService_ListForUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/1/eportfolios" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		eps := []Eportfolio{
			{ID: 10, Name: "My Portfolio", Public: true},
		}
		if err := json.NewEncoder(w).Encode(eps); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewEportfoliosService(client)

	eps, err := svc.ListForUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListForUser failed: %v", err)
	}
	if len(eps) != 1 || eps[0].Name != "My Portfolio" {
		t.Errorf("unexpected result: %+v", eps)
	}
}

func TestEportfoliosService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/eportfolios/10" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		ep := Eportfolio{ID: 10, Name: "My Portfolio"}
		if err := json.NewEncoder(w).Encode(ep); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewEportfoliosService(client)

	ep, err := svc.Get(context.Background(), 10)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if ep.ID != 10 {
		t.Errorf("expected ID 10, got %d", ep.ID)
	}
}

func TestEportfoliosService_Delete(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/eportfolios/10" || r.Method != http.MethodDelete {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewEportfoliosService(client)

	if err := svc.Delete(context.Background(), 10); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if !called {
		t.Error("DELETE was not called")
	}
}

func TestEportfoliosService_ListPages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/eportfolios/10/pages" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		pages := []EportfolioPage{
			{ID: 1, Name: "Introduction"},
			{ID: 2, Name: "Projects"},
		}
		if err := json.NewEncoder(w).Encode(pages); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewEportfoliosService(client)

	pages, err := svc.ListPages(context.Background(), 10)
	if err != nil {
		t.Fatalf("ListPages failed: %v", err)
	}
	if len(pages) != 2 {
		t.Errorf("expected 2 pages, got %d", len(pages))
	}
}

func TestEportfoliosService_Moderate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/eportfolios/10/moderate" || r.Method != http.MethodPut {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		ep := Eportfolio{ID: 10, WorkflowState: "spam"}
		if err := json.NewEncoder(w).Encode(ep); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewEportfoliosService(client)

	ep, err := svc.Moderate(context.Background(), 10, "spam")
	if err != nil {
		t.Fatalf("Moderate failed: %v", err)
	}
	if ep.WorkflowState != "spam" {
		t.Errorf("expected state spam, got %s", ep.WorkflowState)
	}
}
