package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
