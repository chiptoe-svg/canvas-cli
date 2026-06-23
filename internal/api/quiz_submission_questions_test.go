package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuizSubmissionQuestionsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/quiz_submissions/77/questions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizSubmissionQuestionsResponse{
			QuizSubmissionQuestions: []QuizSubmissionQuestion{
				{ID: 1, QuestionType: "multiple_choice_question", QuestionText: "What is Go?"},
				{ID: 2, QuestionType: "true_false_question", QuestionText: "Is Go compiled?"},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizSubmissionQuestionsService(client)

	questions, err := svc.List(context.Background(), 77, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(questions) != 2 {
		t.Errorf("expected 2 questions, got %d", len(questions))
	}
}

func TestQuizSubmissionQuestionsService_List_WithInclude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		inc := r.URL.Query()["include[]"]
		if len(inc) == 0 || inc[0] != "quiz_question" {
			t.Errorf("expected include[]=quiz_question, got %v", inc)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizSubmissionQuestionsResponse{
			QuizSubmissionQuestions: []QuizSubmissionQuestion{{ID: 3}},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizSubmissionQuestionsService(client)

	questions, err := svc.List(context.Background(), 77, &ListQuizSubmissionQuestionsOptions{
		Include: []string{"quiz_question"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(questions) != 1 {
		t.Errorf("expected 1 question, got %d", len(questions))
	}
}

func TestQuizSubmissionQuestionsService_Answer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/quiz_submissions/77/questions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body AnswerQuizSubmissionQuestionsParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Attempt != 1 {
			t.Errorf("expected attempt 1, got %d", body.Attempt)
		}
		if body.ValidationToken != "tok123" {
			t.Errorf("expected validation token 'tok123', got %q", body.ValidationToken)
		}
		if len(body.QuizQuestions) != 1 {
			t.Errorf("expected 1 quiz question answer, got %d", len(body.QuizQuestions))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizSubmissionQuestionsResponse{
			QuizSubmissionQuestions: []QuizSubmissionQuestion{{ID: 1}},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizSubmissionQuestionsService(client)

	params := &AnswerQuizSubmissionQuestionsParams{
		Attempt:         1,
		ValidationToken: "tok123",
		QuizQuestions: []QuizSubmissionAnswerParams{
			{ID: 1, Answer: "choice_a"},
		},
	}
	questions, err := svc.Answer(context.Background(), 77, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(questions) != 1 {
		t.Errorf("expected 1 question, got %d", len(questions))
	}
}

func TestQuizSubmissionQuestionsService_Flag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/quiz_submissions/77/questions/1/flag" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizSubmissionQuestionsResponse{QuizSubmissionQuestions: []QuizSubmissionQuestion{{ID: 1}}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizSubmissionQuestionsService(client)

	q, err := svc.Flag(context.Background(), 77, 1, 1, "tok123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.ID != 1 {
		t.Errorf("expected ID 1, got %d", q.ID)
	}
}

func TestQuizSubmissionQuestionsService_Unflag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/quiz_submissions/77/questions/1/unflag" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizSubmissionQuestionsResponse{QuizSubmissionQuestions: []QuizSubmissionQuestion{{ID: 1}}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizSubmissionQuestionsService(client)

	q, err := svc.Unflag(context.Background(), 77, 1, 1, "tok123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.ID != 1 {
		t.Errorf("expected ID 1, got %d", q.ID)
	}
}
