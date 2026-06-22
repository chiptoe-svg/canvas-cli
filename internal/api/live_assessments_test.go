package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLiveAssessmentsService_List(t *testing.T) {
	envelope := liveAssessmentsEnvelope{LiveAssessments: []LiveAssessment{
		{ID: "abc", Key: "outcome-1", Title: "Comprehension Check"},
	}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/5/live_assessments" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(envelope); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewLiveAssessmentsService(client)
	got, err := svc.List(context.Background(), 5)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1, got %d", len(got))
	}
	if got[0].Key != "outcome-1" {
		t.Errorf("want key 'outcome-1', got %q", got[0].Key)
	}
}

func TestLiveAssessmentsService_Create(t *testing.T) {
	input := []LiveAssessment{{Key: "outcome-2", Title: "New Check"}}
	envelope := liveAssessmentsEnvelope{LiveAssessments: input}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/courses/5/live_assessments" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(envelope); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewLiveAssessmentsService(client)
	got, err := svc.Create(context.Background(), 5, input)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1, got %d", len(got))
	}
}

func TestLiveAssessmentsService_ListResults(t *testing.T) {
	envelope := liveAssessmentResultsEnvelope{Results: []LiveAssessmentResult{
		{Passed: true},
	}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/5/live_assessments/abc/results" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(envelope); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewLiveAssessmentsService(client)
	got, err := svc.ListResults(context.Background(), 5, "abc")
	if err != nil {
		t.Fatalf("ListResults: %v", err)
	}
	if len(got) != 1 || !got[0].Passed {
		t.Errorf("unexpected results: %+v", got)
	}
}

func TestLiveAssessmentsService_CreateResults(t *testing.T) {
	input := []LiveAssessmentResult{{Passed: true}}
	envelope := liveAssessmentResultsEnvelope{Results: input}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/courses/5/live_assessments/abc/results" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(envelope); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewLiveAssessmentsService(client)
	got, err := svc.CreateResults(context.Background(), 5, "abc", input)
	if err != nil {
		t.Fatalf("CreateResults: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1, got %d", len(got))
	}
}
