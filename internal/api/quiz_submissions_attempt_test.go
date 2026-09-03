package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuizSubmissionsService_GetAttempt(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/123/quizzes/456/submissions/789" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(QuizSubmissionsResponse{
			QuizSubmissions: []QuizSubmission{{ID: 789, QuizID: 456, Attempt: 1, Score: 2, WorkflowState: "complete"}},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	sub, err := NewQuizSubmissionsService(client).GetAttempt(context.Background(), 123, 456, 789, 1)
	if err != nil {
		t.Fatalf("GetAttempt() error = %v", err)
	}
	if gotQuery != "attempt=1" {
		t.Errorf("query = %q, want attempt=1", gotQuery)
	}
	if sub.Attempt != 1 || sub.Score != 2 {
		t.Errorf("unexpected submission: %+v", sub)
	}
}
