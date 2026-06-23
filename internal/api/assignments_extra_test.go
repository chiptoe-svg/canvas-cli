package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAssignmentsService_GetDateDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/assignments/2/date_details" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"due_at": nil})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewAssignmentsService(client)
	result, err := svc.GetDateDetails(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("GetDateDetails failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestAssignmentsService_UpdateDateDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/assignments/2/date_details" || r.Method != http.MethodPut {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"due_at": "2025-01-01T00:00:00Z"})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewAssignmentsService(client)
	result, err := svc.UpdateDateDetails(context.Background(), 1, 2, map[string]interface{}{"due_at": "2025-01-01T00:00:00Z"})
	if err != nil {
		t.Fatalf("UpdateDateDetails failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestAssignmentsService_Duplicate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/assignments/2/duplicate" || r.Method != http.MethodPost {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 99, "name": "Copy of Assignment"})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewAssignmentsService(client)
	result, err := svc.Duplicate(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("Duplicate failed: %v", err)
	}
	if result.ID != 99 {
		t.Errorf("expected ID 99, got %d", result.ID)
	}
}

func TestAssignmentsService_ListGradeableStudents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/assignments/2/gradeable_students" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"id": 10, "display_name": "Alice"}})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewAssignmentsService(client)
	students, err := svc.ListGradeableStudents(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("ListGradeableStudents failed: %v", err)
	}
	if len(students) != 1 {
		t.Errorf("expected 1 student, got %d", len(students))
	}
	if students[0].ID != 10 {
		t.Errorf("expected student ID 10, got %d", students[0].ID)
	}
}

func TestAssignmentsService_GetSubmissionSummary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/assignments/2/submission_summary" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"graded": 5, "ungraded": 3, "not_submitted": 2})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewAssignmentsService(client)
	result, err := svc.GetSubmissionSummary(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("GetSubmissionSummary failed: %v", err)
	}
	if result.Graded != 5 {
		t.Errorf("expected Graded 5, got %d", result.Graded)
	}
	if result.NotSubmitted != 2 {
		t.Errorf("expected NotSubmitted 2, got %d", result.NotSubmitted)
	}
}

func TestAssignmentsService_ListOverrides(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/assignments/overrides" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"id": 7, "title": "Section Override"}})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewAssignmentsService(client)
	overrides, err := svc.ListOverrides(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListOverrides failed: %v", err)
	}
	if len(overrides) != 1 {
		t.Errorf("expected 1 override, got %d", len(overrides))
	}
}

func TestAssignmentGroupsService_ListAssignments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/assignment_groups/5/assignments" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"id": 42, "name": "Group Assignment"}})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewAssignmentGroupsService(client)
	assignments, err := svc.ListAssignments(context.Background(), 1, 5, nil)
	if err != nil {
		t.Fatalf("ListAssignments failed: %v", err)
	}
	if len(assignments) != 1 {
		t.Errorf("expected 1 assignment, got %d", len(assignments))
	}
	if assignments[0].ID != 42 {
		t.Errorf("expected ID 42, got %d", assignments[0].ID)
	}
}
