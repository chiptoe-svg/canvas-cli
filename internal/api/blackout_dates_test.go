package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountBlackoutDatesService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/blackout_dates" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]BlackoutDate{
			{ID: 1, StartDate: "2024-01-01", EndDate: "2024-01-02", EventTitle: "Holiday"},
			{ID: 2, StartDate: "2024-02-01", EndDate: "2024-02-01", EventTitle: "Break"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountBlackoutDatesService(client)

	dates, err := svc.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(dates) != 2 {
		t.Errorf("expected 2 dates, got %d", len(dates))
	}
	if dates[0].EventTitle != "Holiday" {
		t.Errorf("expected EventTitle 'Holiday', got %q", dates[0].EventTitle)
	}
}

func TestAccountBlackoutDatesService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/blackout_dates/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(BlackoutDate{
			ID:         42,
			StartDate:  "2024-03-01",
			EndDate:    "2024-03-05",
			EventTitle: "Spring Break",
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountBlackoutDatesService(client)

	date, err := svc.Get(context.Background(), 1, 42)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if date.ID != 42 {
		t.Errorf("expected ID 42, got %d", date.ID)
	}
	if date.EventTitle != "Spring Break" {
		t.Errorf("expected EventTitle 'Spring Break', got %q", date.EventTitle)
	}
}

func TestAccountBlackoutDatesService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/blackout_dates" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(BlackoutDate{
			ID:         10,
			StartDate:  "2024-04-01",
			EndDate:    "2024-04-05",
			EventTitle: "Exam Week",
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountBlackoutDatesService(client)

	params := &BlackoutDateParams{
		StartDate:  "2024-04-01",
		EndDate:    "2024-04-05",
		EventTitle: "Exam Week",
	}

	date, err := svc.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if date.ID != 10 {
		t.Errorf("expected ID 10, got %d", date.ID)
	}
}

func TestAccountBlackoutDatesService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/blackout_dates/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(BlackoutDate{
			ID:         42,
			StartDate:  "2024-03-01",
			EndDate:    "2024-03-10",
			EventTitle: "Extended Break",
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountBlackoutDatesService(client)

	params := &BlackoutDateParams{
		EndDate:    "2024-03-10",
		EventTitle: "Extended Break",
	}

	date, err := svc.Update(context.Background(), 1, 42, params)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if date.EventTitle != "Extended Break" {
		t.Errorf("expected EventTitle 'Extended Break', got %q", date.EventTitle)
	}
}

func TestAccountBlackoutDatesService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/blackout_dates/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountBlackoutDatesService(client)

	if err := svc.Delete(context.Background(), 1, 42); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}
