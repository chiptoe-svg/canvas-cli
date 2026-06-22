package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGradingStandardsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/grading_standards" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]GradingStandard{
			{ID: 1, Title: "Letter Grade", ContextType: "Account", ContextID: 1},
			{ID: 2, Title: "Pass/Fail", ContextType: "Account", ContextID: 1},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingStandardsService(client)

	standards, err := svc.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(standards) != 2 {
		t.Errorf("expected 2 standards, got %d", len(standards))
	}
}

func TestGradingStandardsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/grading_standards" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(GradingStandard{
			ID:    10,
			Title: "My Standard",
			GradingScheme: []GradingSchemeEntry{
				{Name: "A", Value: 0.94},
				{Name: "B", Value: 0.84},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingStandardsService(client)

	params := &GradingStandardParams{
		Title: "My Standard",
		GradingScheme: []GradingSchemeEntry{
			{Name: "A", Value: 0.94},
		},
	}
	standard, err := svc.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if standard.ID != 10 {
		t.Errorf("expected ID 10, got %d", standard.ID)
	}
}

func TestGradingStandardsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/grading_standards/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GradingStandard{
			ID:    5,
			Title: "Letter Grade",
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingStandardsService(client)

	standard, err := svc.Get(context.Background(), 1, 5)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if standard.ID != 5 {
		t.Errorf("expected ID 5, got %d", standard.ID)
	}
}

func TestGradingStandardsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/grading_standards/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GradingStandard{
			ID:    5,
			Title: "Updated Standard",
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingStandardsService(client)

	params := &GradingStandardParams{Title: "Updated Standard"}
	standard, err := svc.Update(context.Background(), 1, 5, params)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if standard.Title != "Updated Standard" {
		t.Errorf("expected 'Updated Standard', got %q", standard.Title)
	}
}

func TestGradingStandardsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/grading_standards/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingStandardsService(client)

	err := svc.Delete(context.Background(), 1, 5)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
}
