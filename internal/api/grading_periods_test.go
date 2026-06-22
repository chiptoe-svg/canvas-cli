package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGradingPeriodsService_List(t *testing.T) {
	envelope := gradingPeriodsEnvelope{GradingPeriods: []GradingPeriod{
		{ID: 1, Title: "Q1", StartDate: "2024-08-01", EndDate: "2024-10-15"},
		{ID: 2, Title: "Q2", StartDate: "2024-10-16", EndDate: "2025-01-15"},
	}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/10/grading_periods" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(envelope); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingPeriodsService(client)
	got, err := svc.List(context.Background(), 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2, got %d", len(got))
	}
}

func TestGradingPeriodsService_Get(t *testing.T) {
	envelope := gradingPeriodsEnvelope{GradingPeriods: []GradingPeriod{
		{ID: 1, Title: "Q1", StartDate: "2024-08-01", EndDate: "2024-10-15"},
	}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/10/grading_periods/1" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(envelope); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingPeriodsService(client)
	got, err := svc.Get(context.Background(), 10, 1)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Title != "Q1" {
		t.Errorf("want 'Q1', got %q", got.Title)
	}
}

func TestGradingPeriodsService_Update(t *testing.T) {
	envelope := gradingPeriodsEnvelope{GradingPeriods: []GradingPeriod{
		{ID: 1, Title: "Q1 Updated", StartDate: "2024-08-01", EndDate: "2024-10-20"},
	}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/courses/10/grading_periods/1" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(envelope); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingPeriodsService(client)
	got, err := svc.Update(context.Background(), 10, 1, GradingPeriodParams{Title: "Q1 Updated"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got.Title != "Q1 Updated" {
		t.Errorf("want 'Q1 Updated', got %q", got.Title)
	}
}

func TestGradingPeriodsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/courses/10/grading_periods/1" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingPeriodsService(client)
	if err := svc.Delete(context.Background(), 10, 1); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestGradingPeriodsService_BatchUpdate(t *testing.T) {
	envelope := gradingPeriodsEnvelope{GradingPeriods: []GradingPeriod{
		{ID: 1, Title: "Semester 1"},
		{ID: 2, Title: "Semester 2"},
	}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/api/v1/courses/10/grading_periods/batch_update" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(envelope); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGradingPeriodsService(client)
	got, err := svc.BatchUpdate(context.Background(), 10, []GradingPeriodParams{
		{Title: "Semester 1"},
		{Title: "Semester 2"},
	})
	if err != nil {
		t.Fatalf("BatchUpdate: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2, got %d", len(got))
	}
}
