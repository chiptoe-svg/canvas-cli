package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBlackoutDatesService_List(t *testing.T) {
	dates := []BlackoutDate{
		{ID: 1, StartDate: "2024-12-24", EndDate: "2024-12-26", EventTitle: "Winter Break"},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/5/blackout_dates" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(dates); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBlackoutDatesService(client)
	got, err := svc.List(context.Background(), 5)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1, got %d", len(got))
	}
}

func TestBlackoutDatesService_Get(t *testing.T) {
	want := BlackoutDate{ID: 1, StartDate: "2024-12-24", EndDate: "2024-12-26", EventTitle: "Winter Break"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/5/blackout_dates/1" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBlackoutDatesService(client)
	got, err := svc.Get(context.Background(), 5, 1)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.EventTitle != "Winter Break" {
		t.Errorf("want 'Winter Break', got %q", got.EventTitle)
	}
}

func TestBlackoutDatesService_Create(t *testing.T) {
	want := BlackoutDate{ID: 2, StartDate: "2024-11-28", EndDate: "2024-11-29", EventTitle: "Thanksgiving"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/courses/5/blackout_dates" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBlackoutDatesService(client)
	got, err := svc.Create(context.Background(), 5, BlackoutDateParams{
		StartDate:  "2024-11-28",
		EndDate:    "2024-11-29",
		EventTitle: "Thanksgiving",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != 2 {
		t.Errorf("want ID 2, got %d", got.ID)
	}
}

func TestBlackoutDatesService_Update(t *testing.T) {
	want := BlackoutDate{ID: 1, StartDate: "2024-12-23", EndDate: "2024-12-27", EventTitle: "Extended Winter Break"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/courses/5/blackout_dates/1" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBlackoutDatesService(client)
	got, err := svc.Update(context.Background(), 5, 1, BlackoutDateParams{
		StartDate:  "2024-12-23",
		EndDate:    "2024-12-27",
		EventTitle: "Extended Winter Break",
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got.EventTitle != "Extended Winter Break" {
		t.Errorf("want 'Extended Winter Break', got %q", got.EventTitle)
	}
}

func TestBlackoutDatesService_BulkUpdate(t *testing.T) {
	want := []BlackoutDate{{ID: 10, StartDate: "2025-01-01", EndDate: "2025-01-01", EventTitle: "New Year"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/courses/5/blackout_dates" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBlackoutDatesService(client)
	got, err := svc.BulkUpdate(context.Background(), 5, []BlackoutDateParams{
		{StartDate: "2025-01-01", EndDate: "2025-01-01", EventTitle: "New Year"},
	})
	if err != nil {
		t.Fatalf("BulkUpdate: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1, got %d", len(got))
	}
}

func TestBlackoutDatesService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/courses/5/blackout_dates/1" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBlackoutDatesService(client)
	if err := svc.Delete(context.Background(), 5, 1); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}
