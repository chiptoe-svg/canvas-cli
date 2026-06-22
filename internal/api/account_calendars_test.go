package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountCalendarsService_ListAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/account_calendars" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]AccountCalendar{
			{ID: 1, Name: "Main Calendar", Visible: true},
			{ID: 2, Name: "Sub Calendar", Visible: false},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountCalendarsService(client)

	calendars, err := svc.ListAll(context.Background(), "")
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}
	if len(calendars) != 2 {
		t.Errorf("expected 2 calendars, got %d", len(calendars))
	}
}

func TestAccountCalendarsService_ListAll_WithSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Query().Get("search_term") != "main" {
			t.Errorf("expected search_term=main, got %q", r.URL.Query().Get("search_term"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]AccountCalendar{
			{ID: 1, Name: "Main Calendar", Visible: true},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountCalendarsService(client)

	calendars, err := svc.ListAll(context.Background(), "main")
	if err != nil {
		t.Fatalf("ListAll with search: %v", err)
	}
	if len(calendars) != 1 {
		t.Errorf("expected 1 calendar, got %d", len(calendars))
	}
}

func TestAccountCalendarsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/account_calendars/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AccountCalendar{
			ID:      42,
			Name:    "Test Calendar",
			Visible: true,
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountCalendarsService(client)

	cal, err := svc.Get(context.Background(), 42)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if cal.ID != 42 {
		t.Errorf("expected ID 42, got %d", cal.ID)
	}
}

func TestAccountCalendarsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/account_calendars/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		visible := true
		json.NewEncoder(w).Encode(AccountCalendar{
			ID:      42,
			Visible: visible,
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountCalendarsService(client)

	visible := true
	cal, err := svc.Update(context.Background(), 42, &AccountCalendarParams{Visible: &visible})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if !cal.Visible {
		t.Error("expected calendar to be visible")
	}
}

func TestAccountCalendarsService_ListForAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/account_calendars" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]AccountCalendar{
			{ID: 1, Name: "Cal A"},
			{ID: 2, Name: "Cal B"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountCalendarsService(client)

	cals, err := svc.ListForAccount(context.Background(), 1, "")
	if err != nil {
		t.Fatalf("ListForAccount: %v", err)
	}
	if len(cals) != 2 {
		t.Errorf("expected 2 calendars, got %d", len(cals))
	}
}

func TestAccountCalendarsService_GetVisibleCount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/visible_calendars_count" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"count": 5,
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountCalendarsService(client)

	count, err := svc.GetVisibleCount(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetVisibleCount: %v", err)
	}
	if count["count"] == nil {
		t.Error("expected count key in result")
	}
}
