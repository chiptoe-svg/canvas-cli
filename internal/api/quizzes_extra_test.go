package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuizzesService_ValidateAccessCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/validate_access_code" || r.Method != http.MethodPost {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"valid": true})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewQuizzesService(client)
	valid, err := svc.ValidateAccessCode(context.Background(), 1, 2, "secret")
	if err != nil {
		t.Fatalf("ValidateAccessCode failed: %v", err)
	}
	if !valid {
		t.Error("expected valid=true")
	}
}

func TestQuizzesService_MessageSubmissionUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/submission_users/message" || r.Method != http.MethodPost {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewQuizzesService(client)
	err := svc.MessageSubmissionUsers(context.Background(), 1, 2, &QuizMessageParams{Body: "hello", RecipientGroup: "submitted"})
	if err != nil {
		t.Fatalf("MessageSubmissionUsers failed: %v", err)
	}
}

func TestQuizzesService_GetCurrentUserSubmission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/submission" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"quiz_submissions": []map[string]interface{}{{"id": 55, "quiz_id": 2}}})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewQuizzesService(client)
	sub, err := svc.GetCurrentUserSubmission(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("GetCurrentUserSubmission failed: %v", err)
	}
	if sub.ID != 55 {
		t.Errorf("expected ID 55, got %d", sub.ID)
	}
}

func TestQuizzesService_ListGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/groups" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"id": 3, "name": "Group 1", "pick_count": 2}})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewQuizzesService(client)
	groups, err := svc.ListGroups(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("ListGroups failed: %v", err)
	}
	if len(groups) != 1 || groups[0].ID != 3 {
		t.Errorf("unexpected groups: %+v", groups)
	}
}

func TestQuizzesService_GetDateDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/date_details" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"due_at": nil})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewQuizzesService(client)
	result, err := svc.GetDateDetails(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("GetDateDetails failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestQuizzesService_UpdateDateDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/date_details" || r.Method != http.MethodPut {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		dueAt := "2025-01-01T00:00:00Z"
		json.NewEncoder(w).Encode(map[string]interface{}{"due_at": dueAt})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewQuizzesService(client)
	dueAt := "2025-01-01T00:00:00Z"
	result, err := svc.UpdateDateDetails(context.Background(), 1, 2, &QuizDateDetailsParams{DueAt: &dueAt})
	if err != nil {
		t.Fatalf("UpdateDateDetails failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}
