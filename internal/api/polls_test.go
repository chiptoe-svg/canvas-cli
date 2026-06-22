package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPollsService_ListPolls(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls" {
			t.Errorf("expected path /api/v1/polls, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":1,"question":"What is 2+2?","user_id":10},{"id":2,"question":"Favourite colour?","user_id":10}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	polls, err := svc.ListPolls(context.Background())
	if err != nil {
		t.Fatalf("ListPolls failed: %v", err)
	}
	if len(polls) != 2 {
		t.Fatalf("expected 2 polls, got %d", len(polls))
	}
	if polls[0].Question != "What is 2+2?" {
		t.Errorf("unexpected question: %s", polls[0].Question)
	}
}

func TestPollsService_GetPoll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1" {
			t.Errorf("expected path /api/v1/polls/1, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":1,"question":"What is 2+2?","description":"Math check","user_id":10}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	poll, err := svc.GetPoll(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetPoll failed: %v", err)
	}
	if poll.ID != 1 {
		t.Errorf("expected poll ID 1, got %d", poll.ID)
	}
	if poll.Description != "Math check" {
		t.Errorf("expected description 'Math check', got %s", poll.Description)
	}
}

func TestPollsService_CreatePoll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls" {
			t.Errorf("expected path /api/v1/polls, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		polls, ok := body["polls"].([]interface{})
		if !ok || len(polls) == 0 {
			t.Fatal("expected polls array in body")
		}
		pollData := polls[0].(map[string]interface{})
		if pollData["question"] != "What is Go?" {
			t.Errorf("expected question 'What is Go?', got %v", pollData["question"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":5,"question":"What is Go?","description":"A Go question","user_id":10}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	poll, err := svc.CreatePoll(context.Background(), &CreatePollParams{
		Question:    "What is Go?",
		Description: "A Go question",
	})
	if err != nil {
		t.Fatalf("CreatePoll failed: %v", err)
	}
	if poll.ID != 5 {
		t.Errorf("expected poll ID 5, got %d", poll.ID)
	}
}

func TestPollsService_UpdatePoll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1" {
			t.Errorf("expected path /api/v1/polls/1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":1,"question":"Updated question","user_id":10}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	poll, err := svc.UpdatePoll(context.Background(), 1, &UpdatePollParams{
		Question: "Updated question",
	})
	if err != nil {
		t.Fatalf("UpdatePoll failed: %v", err)
	}
	if poll.Question != "Updated question" {
		t.Errorf("expected 'Updated question', got %s", poll.Question)
	}
}

func TestPollsService_DeletePoll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1" {
			t.Errorf("expected path /api/v1/polls/1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	if err := svc.DeletePoll(context.Background(), 1); err != nil {
		t.Fatalf("DeletePoll failed: %v", err)
	}
}

func TestPollsService_ListPollChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1/poll_choices" {
			t.Errorf("expected path /api/v1/polls/1/poll_choices, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":10,"poll_id":1,"text":"Option A","is_correct":true,"position":1}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	choices, err := svc.ListPollChoices(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListPollChoices failed: %v", err)
	}
	if len(choices) != 1 {
		t.Fatalf("expected 1 choice, got %d", len(choices))
	}
	if choices[0].Text != "Option A" {
		t.Errorf("expected 'Option A', got %s", choices[0].Text)
	}
}

func TestPollsService_GetPollChoice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1/poll_choices/10" {
			t.Errorf("expected path /api/v1/polls/1/poll_choices/10, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":10,"poll_id":1,"text":"Option A","is_correct":true,"position":1}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	choice, err := svc.GetPollChoice(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetPollChoice failed: %v", err)
	}
	if choice.ID != 10 {
		t.Errorf("expected choice ID 10, got %d", choice.ID)
	}
}

func TestPollsService_CreatePollChoice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1/poll_choices" {
			t.Errorf("expected path /api/v1/polls/1/poll_choices, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":20,"poll_id":1,"text":"Paris","is_correct":true,"position":1}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	isCorrect := true
	_ = isCorrect
	choice, err := svc.CreatePollChoice(context.Background(), 1, &CreatePollChoiceParams{
		Text:      "Paris",
		IsCorrect: true,
		Position:  1,
	})
	if err != nil {
		t.Fatalf("CreatePollChoice failed: %v", err)
	}
	if choice.Text != "Paris" {
		t.Errorf("expected 'Paris', got %s", choice.Text)
	}
}

func TestPollsService_UpdatePollChoice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1/poll_choices/10" {
			t.Errorf("expected path /api/v1/polls/1/poll_choices/10, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":10,"poll_id":1,"text":"Updated Option","is_correct":false,"position":2}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	isCorrect := false
	choice, err := svc.UpdatePollChoice(context.Background(), 1, 10, &UpdatePollChoiceParams{
		Text:      "Updated Option",
		IsCorrect: &isCorrect,
		Position:  2,
	})
	if err != nil {
		t.Fatalf("UpdatePollChoice failed: %v", err)
	}
	if choice.Text != "Updated Option" {
		t.Errorf("expected 'Updated Option', got %s", choice.Text)
	}
}

func TestPollsService_DeletePollChoice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1/poll_choices/10" {
			t.Errorf("expected path /api/v1/polls/1/poll_choices/10, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	if err := svc.DeletePollChoice(context.Background(), 1, 10); err != nil {
		t.Fatalf("DeletePollChoice failed: %v", err)
	}
}

func TestPollsService_ListPollSessions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1/poll_sessions" {
			t.Errorf("expected path /api/v1/polls/1/poll_sessions, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":100,"poll_id":1,"course_id":999,"is_published":true}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	sessions, err := svc.ListPollSessions(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListPollSessions failed: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].CourseID != 999 {
		t.Errorf("expected course_id 999, got %d", sessions[0].CourseID)
	}
}

func TestPollsService_GetPollSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1/poll_sessions/100" {
			t.Errorf("expected path /api/v1/polls/1/poll_sessions/100, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":100,"poll_id":1,"course_id":999,"is_published":false,"has_public_results":true}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	session, err := svc.GetPollSession(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("GetPollSession failed: %v", err)
	}
	if !session.HasPublicResults {
		t.Error("expected has_public_results to be true")
	}
}

func TestPollsService_CreatePollSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1/poll_sessions" {
			t.Errorf("expected path /api/v1/polls/1/poll_sessions, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":101,"poll_id":1,"course_id":999,"is_published":false}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	session, err := svc.CreatePollSession(context.Background(), 1, &CreatePollSessionParams{
		CourseID:         999,
		HasPublicResults: true,
	})
	if err != nil {
		t.Fatalf("CreatePollSession failed: %v", err)
	}
	if session.ID != 101 {
		t.Errorf("expected session ID 101, got %d", session.ID)
	}
}

func TestPollsService_UpdatePollSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1/poll_sessions/100" {
			t.Errorf("expected path /api/v1/polls/1/poll_sessions/100, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":100,"poll_id":1,"course_id":999,"is_published":false,"has_public_results":false}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	hasPublic := false
	session, err := svc.UpdatePollSession(context.Background(), 1, 100, &UpdatePollSessionParams{
		HasPublicResults: &hasPublic,
	})
	if err != nil {
		t.Fatalf("UpdatePollSession failed: %v", err)
	}
	if session.HasPublicResults {
		t.Error("expected has_public_results to be false")
	}
}

func TestPollsService_DeletePollSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1/poll_sessions/100" {
			t.Errorf("expected path /api/v1/polls/1/poll_sessions/100, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	if err := svc.DeletePollSession(context.Background(), 1, 100); err != nil {
		t.Fatalf("DeletePollSession failed: %v", err)
	}
}

func TestPollsService_OpenPollSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1/poll_sessions/100/open" {
			t.Errorf("expected path /api/v1/polls/1/poll_sessions/100/open, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":100,"poll_id":1,"course_id":999,"is_published":true}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	session, err := svc.OpenPollSession(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("OpenPollSession failed: %v", err)
	}
	if !session.IsPublished {
		t.Error("expected is_published to be true")
	}
}

func TestPollsService_ClosePollSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1/poll_sessions/100/close" {
			t.Errorf("expected path /api/v1/polls/1/poll_sessions/100/close, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":100,"poll_id":1,"course_id":999,"is_published":false}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	session, err := svc.ClosePollSession(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("ClosePollSession failed: %v", err)
	}
	if session.IsPublished {
		t.Error("expected is_published to be false")
	}
}

func TestPollsService_ListOpenedPollSessions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/poll_sessions/opened" {
			t.Errorf("expected path /api/v1/poll_sessions/opened, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":100,"poll_id":1,"course_id":999,"is_published":true}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	sessions, err := svc.ListOpenedPollSessions(context.Background())
	if err != nil {
		t.Fatalf("ListOpenedPollSessions failed: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
}

func TestPollsService_ListClosedPollSessions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/poll_sessions/closed" {
			t.Errorf("expected path /api/v1/poll_sessions/closed, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":50,"poll_id":2,"course_id":888,"is_published":false}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	sessions, err := svc.ListClosedPollSessions(context.Background())
	if err != nil {
		t.Fatalf("ListClosedPollSessions failed: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
}

func TestPollsService_GetPollSubmission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1/poll_sessions/100/poll_submissions/200" {
			t.Errorf("expected path /api/v1/polls/1/poll_sessions/100/poll_submissions/200, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":200,"poll_choice_id":10,"user_id":42}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	sub, err := svc.GetPollSubmission(context.Background(), 1, 100, 200)
	if err != nil {
		t.Fatalf("GetPollSubmission failed: %v", err)
	}
	if sub.PollChoiceID != 10 {
		t.Errorf("expected poll_choice_id 10, got %d", sub.PollChoiceID)
	}
}

func TestPollsService_CreatePollSubmission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/polls/1/poll_sessions/100/poll_submissions" {
			t.Errorf("expected path /api/v1/polls/1/poll_sessions/100/poll_submissions, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":201,"poll_choice_id":10,"user_id":42}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPollsService(client)

	sub, err := svc.CreatePollSubmission(context.Background(), 1, 100, &CreatePollSubmissionParams{
		PollChoiceID: 10,
	})
	if err != nil {
		t.Fatalf("CreatePollSubmission failed: %v", err)
	}
	if sub.ID != 201 {
		t.Errorf("expected submission ID 201, got %d", sub.ID)
	}
}

func TestNewPollsService(t *testing.T) {
	client := &Client{}
	svc := NewPollsService(client)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
	if svc.client != client {
		t.Error("expected service client to match input client")
	}
}
