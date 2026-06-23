package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountNotificationsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/account_notifications" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		notifications := []AccountNotification{
			{ID: 1, Subject: "Test Notification 1", Message: "Hello", Icon: "warning"},
			{ID: 2, Subject: "Test Notification 2", Message: "World", Icon: "information"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notifications)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountNotificationsService(client)

	notifications, err := svc.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(notifications) != 2 {
		t.Errorf("expected 2 notifications, got %d", len(notifications))
	}
	if notifications[0].Subject != "Test Notification 1" {
		t.Errorf("expected subject 'Test Notification 1', got %s", notifications[0].Subject)
	}
}

func TestAccountNotificationsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/account_notifications" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		notification := AccountNotification{
			ID:      3,
			Subject: "New Notification",
			Message: "Created!",
			StartAt: "2026-01-01T00:00:00Z",
			EndAt:   "2026-12-31T23:59:59Z",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notification)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountNotificationsService(client)

	params := &AccountNotificationParams{
		AccountNotification: AccountNotificationFields{
			Subject: "New Notification",
			Message: "Created!",
			StartAt: "2026-01-01T00:00:00Z",
			EndAt:   "2026-12-31T23:59:59Z",
		},
	}
	notification, err := svc.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if notification.ID != 3 {
		t.Errorf("expected ID 3, got %d", notification.ID)
	}
}

func TestAccountNotificationsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/account_notifications/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		notification := AccountNotification{
			ID:      5,
			Subject: "Fetched Notification",
			Message: "Hello World",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notification)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountNotificationsService(client)

	notification, err := svc.Get(context.Background(), 1, 5)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if notification.ID != 5 {
		t.Errorf("expected ID 5, got %d", notification.ID)
	}
}

func TestAccountNotificationsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/account_notifications/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		notification := AccountNotification{
			ID:      5,
			Subject: "Updated Notification",
			Message: "Updated!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notification)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountNotificationsService(client)

	params := &AccountNotificationParams{
		AccountNotification: AccountNotificationFields{
			Subject: "Updated Notification",
			Message: "Updated!",
			StartAt: "2026-01-01T00:00:00Z",
			EndAt:   "2026-12-31T23:59:59Z",
		},
	}
	notification, err := svc.Update(context.Background(), 1, 5, params)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if notification.Subject != "Updated Notification" {
		t.Errorf("expected subject 'Updated Notification', got %s", notification.Subject)
	}
}

func TestAccountNotificationsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/account_notifications/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountNotificationsService(client)

	err := svc.Delete(context.Background(), 1, 5)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
}
