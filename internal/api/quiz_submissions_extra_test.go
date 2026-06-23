package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuizSubmissionsService_ListEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/submissions/3/events" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"quiz_submission_events": []map[string]interface{}{
				{"event_type": "question_answered", "created_at": "2025-01-01T00:00:00Z"},
			},
		})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewQuizSubmissionsService(client)
	events, err := svc.ListEvents(context.Background(), 1, 2, 3)
	if err != nil {
		t.Fatalf("ListEvents failed: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
	if events[0].EventType != "question_answered" {
		t.Errorf("unexpected event type: %s", events[0].EventType)
	}
}

func TestQuizSubmissionsService_CreateEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/submissions/3/events" || r.Method != http.MethodPost {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewQuizSubmissionsService(client)
	err := svc.CreateEvents(context.Background(), 1, 2, 3, []QuizSubmissionEvent{
		{EventType: "question_answered", CreatedAt: "2025-01-01T00:00:00Z"},
	})
	if err != nil {
		t.Fatalf("CreateEvents failed: %v", err)
	}
}

func TestQuizSubmissionsService_GetTime(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/submissions/3/time" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"time_left": 300})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewQuizSubmissionsService(client)
	result, err := svc.GetTime(context.Background(), 1, 2, 3)
	if err != nil {
		t.Fatalf("GetTime failed: %v", err)
	}
	if result.TimeLeft != 300 {
		t.Errorf("expected TimeLeft 300, got %d", result.TimeLeft)
	}
}
