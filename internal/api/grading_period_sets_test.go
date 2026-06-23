package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGradingPeriodSetsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/grading_period_sets" {
			t.Errorf("expected path /api/v1/accounts/1/grading_period_sets, got %s", r.URL.Path)
		}

		sets := []GradingPeriodSet{
			{ID: 1, Title: "2024-2025"},
			{ID: 2, Title: "2025-2026"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sets)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewGradingPeriodSetsService(client)

	sets, err := service.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(sets) != 2 {
		t.Fatalf("expected 2 grading period sets, got %d", len(sets))
	}

	if sets[0].Title != "2024-2025" {
		t.Errorf("expected title '2024-2025', got '%s'", sets[0].Title)
	}
}

func TestGradingPeriodSetsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/grading_period_sets" {
			t.Errorf("expected path /api/v1/accounts/1/grading_period_sets, got %s", r.URL.Path)
		}

		set := GradingPeriodSet{ID: 10, Title: "New Set"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(set)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewGradingPeriodSetsService(client)

	params := &GradingPeriodSetParams{GradingPeriodSet: GradingPeriodSetFields{Title: "New Set"}}
	set, err := service.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if set.ID != 10 {
		t.Errorf("expected ID 10, got %d", set.ID)
	}
}

func TestGradingPeriodSetsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/grading_period_sets/10" {
			t.Errorf("expected path /api/v1/accounts/1/grading_period_sets/10, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewGradingPeriodSetsService(client)

	if err := service.Delete(context.Background(), 1, 10); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestGradingPeriodSetsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/grading_period_sets/10" {
			t.Errorf("expected path /api/v1/accounts/1/grading_period_sets/10, got %s", r.URL.Path)
		}

		set := GradingPeriodSet{ID: 10, Title: "Updated Set"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(set)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewGradingPeriodSetsService(client)

	params := &GradingPeriodSetParams{GradingPeriodSet: GradingPeriodSetFields{Title: "Updated Set"}}
	set, err := service.Update(context.Background(), 1, 10, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if set.Title != "Updated Set" {
		t.Errorf("expected title 'Updated Set', got '%s'", set.Title)
	}
}

func TestGradingPeriodSetsService_ListPeriods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/grading_periods" {
			t.Errorf("expected path /api/v1/accounts/1/grading_periods, got %s", r.URL.Path)
		}

		periods := []GradingPeriod{
			{ID: 1, Title: "Q1", IsClosed: true},
			{ID: 2, Title: "Q2", IsClosed: false},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(periods)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewGradingPeriodSetsService(client)

	periods, err := service.ListPeriods(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListPeriods failed: %v", err)
	}

	if len(periods) != 2 {
		t.Fatalf("expected 2 grading periods, got %d", len(periods))
	}
}

func TestGradingPeriodSetsService_DeletePeriod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/grading_periods/5" {
			t.Errorf("expected path /api/v1/accounts/1/grading_periods/5, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewGradingPeriodSetsService(client)

	if err := service.DeletePeriod(context.Background(), 1, 5); err != nil {
		t.Fatalf("DeletePeriod failed: %v", err)
	}
}
