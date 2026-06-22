package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuditLogsService_ListAuthenticationForAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/audit/authentication/accounts/1" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		events := []AuditLogEvent{{ID: "evt1", EventType: "login"}}
		if err := json.NewEncoder(w).Encode(events); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAuditLogsService(client)

	events, err := svc.ListAuthenticationForAccount(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("ListAuthenticationForAccount failed: %v", err)
	}
	if len(events) != 1 || events[0].EventType != "login" {
		t.Errorf("unexpected result: %+v", events)
	}
}

func TestAuditLogsService_ListAuthenticationForUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/audit/authentication/users/5" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		events := []AuditLogEvent{{ID: "e2", EventType: "logout"}}
		if err := json.NewEncoder(w).Encode(events); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAuditLogsService(client)

	events, err := svc.ListAuthenticationForUser(context.Background(), 5, nil)
	if err != nil {
		t.Fatalf("ListAuthenticationForUser failed: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

func TestAuditLogsService_ListAuthenticationForLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/audit/authentication/logins/10" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		events := []AuditLogEvent{{ID: "e3", EventType: "login"}}
		if err := json.NewEncoder(w).Encode(events); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAuditLogsService(client)

	events, err := svc.ListAuthenticationForLogin(context.Background(), 10, nil)
	if err != nil {
		t.Fatalf("ListAuthenticationForLogin failed: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

func TestAuditLogsService_ListGradeChangeEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/audit/grade_change" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		events := []AuditLogEvent{{ID: "gc1", EventType: "grade_change"}}
		if err := json.NewEncoder(w).Encode(events); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAuditLogsService(client)

	events, err := svc.ListGradeChangeEvents(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListGradeChangeEvents failed: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

func TestAuditLogsService_ListGradeChangeForCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/audit/grade_change/courses/20" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		events := []AuditLogEvent{{ID: "gc2"}}
		if err := json.NewEncoder(w).Encode(events); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAuditLogsService(client)

	events, err := svc.ListGradeChangeForCourse(context.Background(), 20, nil)
	if err != nil {
		t.Fatalf("ListGradeChangeForCourse failed: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}
