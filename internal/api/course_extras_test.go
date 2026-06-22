package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCourseExtrasService_CreateQuizExtensions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/courses/5/quiz_extensions" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseExtrasService(client)
	err := svc.CreateQuizExtensions(context.Background(), 5, []QuizExtensionParams{
		{UserID: 10, ExtraTime: 30},
	})
	if err != nil {
		t.Fatalf("CreateQuizExtensions: %v", err)
	}
}

func TestCourseExtrasService_CreateAssignmentExtensions(t *testing.T) {
	want := assignmentExtensionsBody{AssignmentExtensions: []AssignmentExtension{
		{AssignmentID: 20, UserID: 10, ExtraAttempts: 2},
	}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/courses/5/assignments/20/extensions" {
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
	svc := NewCourseExtrasService(client)
	got, err := svc.CreateAssignmentExtensions(context.Background(), 5, 20, []AssignmentExtension{
		{AssignmentID: 20, UserID: 10, ExtraAttempts: 2},
	})
	if err != nil {
		t.Fatalf("CreateAssignmentExtensions: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1, got %d", len(got))
	}
}

func TestGradesService_GetGradebookHistoryDay(t *testing.T) {
	graders := []GradebookHistoryGrader{{ID: 1, Name: "Teacher A"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/5/gradebook_history/2024-11-01" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(graders); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradesService(client)
	got, err := svc.GetGradebookHistoryDay(context.Background(), 5, "2024-11-01")
	if err != nil {
		t.Fatalf("GetGradebookHistoryDay: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Teacher A" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestGradesService_GetGradebookHistorySubmissions(t *testing.T) {
	subs := []GradebookHistorySubmission{{UserID: 3, Grade: "A"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/5/gradebook_history/2024-11-01/graders/1/assignments/20/submissions" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(subs); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradesService(client)
	got, err := svc.GetGradebookHistorySubmissions(context.Background(), 5, "2024-11-01", 1, 20)
	if err != nil {
		t.Fatalf("GetGradebookHistorySubmissions: %v", err)
	}
	if len(got) != 1 || got[0].Grade != "A" {
		t.Errorf("unexpected: %+v", got)
	}
}
