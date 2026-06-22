package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountCalendarsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/account_calendars" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		cals := []AccountCalendar{
			{ID: 1, Name: "Root Account", Visible: true},
			{ID: 2, Name: "Sub Account", Visible: false},
		}
		if err := json.NewEncoder(w).Encode(cals); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountCalendarsService(client)

	cals, err := svc.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(cals) != 2 {
		t.Errorf("expected 2 calendars, got %d", len(cals))
	}
}

func TestAccountCalendarsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/account_calendars/1" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		cal := AccountCalendar{ID: 1, Name: "Root Account", Visible: true}
		if err := json.NewEncoder(w).Encode(cal); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountCalendarsService(client)

	cal, err := svc.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if cal.ID != 1 {
		t.Errorf("expected ID 1, got %d", cal.ID)
	}
}

func TestAccountCalendarsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/account_calendars/1" || r.Method != http.MethodPut {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		visible := true
		cal := AccountCalendar{ID: 1, Visible: visible}
		if err := json.NewEncoder(w).Encode(cal); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountCalendarsService(client)

	vis := true
	cal, err := svc.Update(context.Background(), 1, &UpdateAccountCalendarParams{Visible: &vis})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if !cal.Visible {
		t.Error("expected Visible=true")
	}
}

func TestAccountCalendarsService_ListForAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/5/account_calendars" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		cals := []AccountCalendar{{ID: 10, Name: "Sub"}}
		if err := json.NewEncoder(w).Encode(cals); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountCalendarsService(client)

	cals, err := svc.ListForAccount(context.Background(), 5, nil)
	if err != nil {
		t.Fatalf("ListForAccount failed: %v", err)
	}
	if len(cals) != 1 || cals[0].ID != 10 {
		t.Errorf("unexpected result: %+v", cals)
	}
}
