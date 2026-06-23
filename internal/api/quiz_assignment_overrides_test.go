package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuizAssignmentOverridesService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/assignment_overrides" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizAssignmentOverridesResponse{
			QuizAssignmentOverrides: []QuizAssignmentOverrideSet{
				{QuizID: "10", Overrides: []AssignmentOverrideEntry{{ID: 1, Title: "Section A"}}},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizAssignmentOverridesService(client)

	overrides, err := svc.List(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(overrides) != 1 {
		t.Fatalf("expected 1 override set, got %d", len(overrides))
	}
	if overrides[0].QuizID != "10" {
		t.Errorf("expected QuizID '10', got %q", overrides[0].QuizID)
	}
}

func TestQuizAssignmentOverridesService_List_WithIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		// Verify that quiz IDs were appended as query params.
		if r.URL.RawQuery == "" {
			t.Errorf("expected query params with quiz IDs")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizAssignmentOverridesResponse{
			QuizAssignmentOverrides: []QuizAssignmentOverrideSet{
				{QuizID: "5"},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizAssignmentOverridesService(client)

	overrides, err := svc.List(context.Background(), 1, []int64{5, 6})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(overrides) != 1 {
		t.Errorf("expected 1 override set, got %d", len(overrides))
	}
}

func TestQuizAssignmentOverridesService_Set(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/assignment_overrides" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body SetQuizAssignmentOverridesParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(body.QuizAssignmentOverrides) == 0 {
			t.Errorf("expected at least one override in body")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizAssignmentOverridesResponse{
			QuizAssignmentOverrides: []QuizAssignmentOverrideSet{
				{QuizID: "10"},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizAssignmentOverridesService(client)

	params := &SetQuizAssignmentOverridesParams{
		QuizAssignmentOverrides: []QuizAssignmentOverrideSetInput{
			{
				QuizID: "10",
				Overrides: []AssignmentOverrideEntry{
					{Title: "Section A Override", CourseSectionID: 3},
				},
			},
		},
	}

	overrides, err := svc.Set(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(overrides) != 1 {
		t.Fatalf("expected 1 override set, got %d", len(overrides))
	}
}
