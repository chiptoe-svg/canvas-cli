package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuizQuestionGroupsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/groups/3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizQuestionGroupResponse{
			QuizGroups: []QuizQuestionGroup{
				{ID: 3, QuizID: 2, Name: "Group A", PickCount: 5, QuestionPoints: 2.0},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizQuestionGroupsService(client)

	grp, err := svc.Get(context.Background(), 1, 2, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if grp.ID != 3 {
		t.Errorf("expected ID 3, got %d", grp.ID)
	}
	if grp.Name != "Group A" {
		t.Errorf("expected name 'Group A', got %q", grp.Name)
	}
}

func TestQuizQuestionGroupsService_Get_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizQuestionGroupResponse{})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizQuestionGroupsService(client)

	_, err := svc.Get(context.Background(), 1, 2, 3)
	if err == nil {
		t.Fatal("expected error for empty response")
	}
}

func TestQuizQuestionGroupsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		groups, ok := body["quiz_groups"].([]interface{})
		if !ok || len(groups) == 0 {
			t.Errorf("expected quiz_groups array in body")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizQuestionGroupResponse{
			QuizGroups: []QuizQuestionGroup{
				{ID: 10, QuizID: 2, Name: "New Group", PickCount: 3, QuestionPoints: 1.5},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizQuestionGroupsService(client)

	grp, err := svc.Create(context.Background(), 1, 2, &CreateQuizQuestionGroupParams{
		Name:           "New Group",
		PickCount:      3,
		QuestionPoints: 1.5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if grp.ID != 10 {
		t.Errorf("expected ID 10, got %d", grp.ID)
	}
}

func TestQuizQuestionGroupsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/groups/3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizQuestionGroupResponse{
			QuizGroups: []QuizQuestionGroup{
				{ID: 3, QuizID: 2, Name: "Updated Group"},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizQuestionGroupsService(client)

	name := "Updated Group"
	pick := 4
	pts := 2.5
	grp, err := svc.Update(context.Background(), 1, 2, 3, &UpdateQuizQuestionGroupParams{
		Name:           &name,
		PickCount:      &pick,
		QuestionPoints: &pts,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if grp.Name != "Updated Group" {
		t.Errorf("expected 'Updated Group', got %q", grp.Name)
	}
}

func TestQuizQuestionGroupsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/groups/3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizQuestionGroupsService(client)

	if err := svc.Delete(context.Background(), 1, 2, 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestQuizQuestionGroupsService_ReorderItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/groups/3/reorder" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		order, ok := body["order"].([]interface{})
		if !ok || len(order) != 2 {
			t.Errorf("expected 2 items in order, got %v", body["order"])
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizQuestionGroupsService(client)

	order := []ReorderItem{{ID: 10, Type: "question"}, {ID: 11, Type: "question"}}
	if err := svc.ReorderItems(context.Background(), 1, 2, 3, order); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
