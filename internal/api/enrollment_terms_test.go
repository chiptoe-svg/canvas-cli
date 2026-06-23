package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEnrollmentTermsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/terms" {
			t.Errorf("expected path /api/v1/accounts/1/terms, got %s", r.URL.Path)
		}

		wrapper := EnrollmentTermsResponse{
			EnrollmentTerms: []EnrollmentTerm{
				{ID: 1, Name: "Fall 2024", WorkflowState: "active"},
				{ID: 2, Name: "Spring 2025", WorkflowState: "active"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wrapper)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewEnrollmentTermsService(client)

	terms, err := service.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(terms) != 2 {
		t.Fatalf("expected 2 terms, got %d", len(terms))
	}

	if terms[0].Name != "Fall 2024" {
		t.Errorf("expected name 'Fall 2024', got '%s'", terms[0].Name)
	}
}

func TestEnrollmentTermsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/terms/42" {
			t.Errorf("expected path /api/v1/accounts/1/terms/42, got %s", r.URL.Path)
		}

		term := EnrollmentTerm{ID: 42, Name: "Fall 2024", WorkflowState: "active"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(term)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewEnrollmentTermsService(client)

	term, err := service.Get(context.Background(), 1, 42)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if term.ID != 42 {
		t.Errorf("expected ID 42, got %d", term.ID)
	}
}

func TestEnrollmentTermsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/terms" {
			t.Errorf("expected path /api/v1/accounts/1/terms, got %s", r.URL.Path)
		}

		term := EnrollmentTerm{ID: 10, Name: "Summer 2025", WorkflowState: "active"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(term)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewEnrollmentTermsService(client)

	params := &EnrollmentTermParams{EnrollmentTerm: EnrollmentTermFields{Name: "Summer 2025"}}
	term, err := service.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if term.ID != 10 {
		t.Errorf("expected ID 10, got %d", term.ID)
	}
}

func TestEnrollmentTermsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/terms/42" {
			t.Errorf("expected path /api/v1/accounts/1/terms/42, got %s", r.URL.Path)
		}

		term := EnrollmentTerm{ID: 42, Name: "Updated Term", WorkflowState: "active"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(term)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewEnrollmentTermsService(client)

	params := &EnrollmentTermParams{EnrollmentTerm: EnrollmentTermFields{Name: "Updated Term"}}
	term, err := service.Update(context.Background(), 1, 42, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if term.Name != "Updated Term" {
		t.Errorf("expected name 'Updated Term', got '%s'", term.Name)
	}
}

func TestEnrollmentTermsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/terms/42" {
			t.Errorf("expected path /api/v1/accounts/1/terms/42, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewEnrollmentTermsService(client)

	if err := service.Delete(context.Background(), 1, 42); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}
