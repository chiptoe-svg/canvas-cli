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
func TestGradingStandardsService_ListForCourse(t *testing.T) {
	standards := []GradingStandard{
		{ID: 1, Title: "Standard Grading", ContextType: "Course", ContextID: 10},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/10/grading_standards" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(standards); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingStandardsService(client)
	got, err := svc.ListForCourse(context.Background(), 10)
	if err != nil {
		t.Fatalf("ListForCourse: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1, got %d", len(got))
	}
}

func TestGradingStandardsService_GetForCourse(t *testing.T) {
	want := GradingStandard{ID: 1, Title: "Standard Grading", ContextType: "Course", ContextID: 10}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/10/grading_standards/1" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingStandardsService(client)
	got, err := svc.GetForCourse(context.Background(), 10, 1)
	if err != nil {
		t.Fatalf("GetForCourse: %v", err)
	}
	if got.Title != "Standard Grading" {
		t.Errorf("want 'Standard Grading', got %q", got.Title)
	}
}

func TestGradingStandardsService_CreateForCourse(t *testing.T) {
	want := GradingStandard{ID: 2, Title: "Custom Scheme", ContextType: "Course", ContextID: 10}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/courses/10/grading_standards" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingStandardsService(client)
	got, err := svc.CreateForCourse(context.Background(), 10, GradingStandardParams{Title: "Custom Scheme"})
	if err != nil {
		t.Fatalf("CreateForCourse: %v", err)
	}
	if got.ID != 2 {
		t.Errorf("want ID 2, got %d", got.ID)
	}
}

func TestGradingStandardsService_UpdateForCourse(t *testing.T) {
	want := GradingStandard{ID: 1, Title: "Updated Scheme", ContextType: "Course", ContextID: 10}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/courses/10/grading_standards/1" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingStandardsService(client)
	got, err := svc.UpdateForCourse(context.Background(), 10, 1, GradingStandardParams{Title: "Updated Scheme"})
	if err != nil {
		t.Fatalf("UpdateForCourse: %v", err)
	}
	if got.Title != "Updated Scheme" {
		t.Errorf("want 'Updated Scheme', got %q", got.Title)
	}
}

func TestGradingStandardsService_DeleteForCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/courses/10/grading_standards/1" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingStandardsService(client)
	if err := svc.DeleteForCourse(context.Background(), 10, 1); err != nil {
		t.Fatalf("DeleteForCourse: %v", err)
	}
}
