package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommunicationChannelsService_DeleteByTypeAddress(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/communication_channels/email/test@example.com" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CommunicationChannel{ID: 7, Address: "test@example.com", Type: "email"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	ch, err := svc.DeleteByTypeAddress(context.Background(), 42, "email", "test@example.com")
	if err != nil {
		t.Fatalf("DeleteByTypeAddress: %v", err)
	}
	if ch.Address != "test@example.com" {
		t.Errorf("expected address=test@example.com, got %s", ch.Address)
	}
	if ch.Type != "email" {
		t.Errorf("expected type=email, got %s", ch.Type)
	}
}

func TestCommunicationChannelsService_DeleteByTypeAddress_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	_, err := svc.DeleteByTypeAddress(context.Background(), 42, "email", "missing@example.com")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCommunicationChannelsService_GetNotificationPreferenceCategories(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/communication_channels/7/notification_preference_categories" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{"Course Activities", "Discussions", "Conversations"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	cats, err := svc.GetNotificationPreferenceCategories(context.Background(), 42, 7)
	if err != nil {
		t.Fatalf("GetNotificationPreferenceCategories: %v", err)
	}
	if len(cats) != 3 {
		t.Fatalf("expected 3 categories, got %d", len(cats))
	}
	if cats[0] != "Course Activities" {
		t.Errorf("expected 'Course Activities', got %s", cats[0])
	}
}

func TestCommunicationChannelsService_GetNotificationPreference(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/communication_channels/7/notification_preferences/new_assignment" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(NotificationPreference{Frequency: "immediately", Notification: "new_assignment"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	pref, err := svc.GetNotificationPreference(context.Background(), 42, 7, "new_assignment")
	if err != nil {
		t.Fatalf("GetNotificationPreference: %v", err)
	}
	if pref.Frequency != "immediately" {
		t.Errorf("expected frequency=immediately, got %s", pref.Frequency)
	}
}

func TestCommunicationChannelsService_GetNotificationPreferencesByType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/communication_channels/email/user@example.com/notification_preferences" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(NotificationPreferences{
			NotificationPreferences: []NotificationPreference{{Frequency: "daily", Notification: "new_discussion"}},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	prefs, err := svc.GetNotificationPreferencesByType(context.Background(), 42, "email", "user@example.com")
	if err != nil {
		t.Fatalf("GetNotificationPreferencesByType: %v", err)
	}
	if len(prefs.NotificationPreferences) != 1 {
		t.Fatalf("expected 1 preference, got %d", len(prefs.NotificationPreferences))
	}
	if prefs.NotificationPreferences[0].Frequency != "daily" {
		t.Errorf("expected frequency=daily, got %s", prefs.NotificationPreferences[0].Frequency)
	}
}

func TestCommunicationChannelsService_GetNotificationPreferenceByType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/communication_channels/email/user@example.com/notification_preferences/new_assignment" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(NotificationPreference{Frequency: "weekly", Notification: "new_assignment"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	pref, err := svc.GetNotificationPreferenceByType(context.Background(), 42, "email", "user@example.com", "new_assignment")
	if err != nil {
		t.Fatalf("GetNotificationPreferenceByType: %v", err)
	}
	if pref.Frequency != "weekly" {
		t.Errorf("expected frequency=weekly, got %s", pref.Frequency)
	}
}

func TestCommunicationChannelsService_UpdateNotificationPreferences(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/communication_channels/7/notification_preferences" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body NotificationPreferencesUpdateBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if _, ok := body.NotificationPreferences["new_assignment"]; !ok {
			t.Error("expected new_assignment in preferences")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(NotificationPreferences{
			NotificationPreferences: []NotificationPreference{{Frequency: "never", Notification: "new_assignment"}},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	body := NotificationPreferencesUpdateBody{
		NotificationPreferences: map[string]NotificationPreferenceUpdate{
			"new_assignment": {Frequency: "never"},
		},
	}
	prefs, err := svc.UpdateNotificationPreferences(context.Background(), 7, body)
	if err != nil {
		t.Fatalf("UpdateNotificationPreferences: %v", err)
	}
	if len(prefs.NotificationPreferences) != 1 {
		t.Fatalf("expected 1 preference, got %d", len(prefs.NotificationPreferences))
	}
}

func TestCommunicationChannelsService_UpdateNotificationPreference(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/communication_channels/7/notification_preferences/new_assignment" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(NotificationPreference{Frequency: "immediately", Notification: "new_assignment"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	pref, err := svc.UpdateNotificationPreference(context.Background(), 7, "new_assignment", NotificationPreference{Frequency: "immediately"})
	if err != nil {
		t.Fatalf("UpdateNotificationPreference: %v", err)
	}
	if pref.Frequency != "immediately" {
		t.Errorf("expected frequency=immediately, got %s", pref.Frequency)
	}
}

func TestCommunicationChannelsService_UpdateNotificationPreferenceCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/self/communication_channels/7/notification_preference_categories/Course+Activities" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(NotificationPreferences{})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	result, err := svc.UpdateNotificationPreferenceCategory(context.Background(), 7, "Course+Activities", NotificationPreferences{})
	if err != nil {
		t.Fatalf("UpdateNotificationPreferenceCategory: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestCommunicationChannelsService_UpdateNotificationPreferencesByType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/self/communication_channels/email/user@example.com/notification_preferences" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(NotificationPreferences{})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	body := NotificationPreferencesUpdateBody{
		NotificationPreferences: map[string]NotificationPreferenceUpdate{
			"graded": {Frequency: "immediately"},
		},
	}
	result, err := svc.UpdateNotificationPreferencesByType(context.Background(), "email", "user@example.com", body)
	if err != nil {
		t.Fatalf("UpdateNotificationPreferencesByType: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestCommunicationChannelsService_UpdateNotificationPreferenceByType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/self/communication_channels/email/user@example.com/notification_preferences/new_assignment" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(NotificationPreference{Frequency: "daily", Notification: "new_assignment"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	pref, err := svc.UpdateNotificationPreferenceByType(context.Background(), "email", "user@example.com", "new_assignment", NotificationPreference{Frequency: "daily"})
	if err != nil {
		t.Fatalf("UpdateNotificationPreferenceByType: %v", err)
	}
	if pref.Frequency != "daily" {
		t.Errorf("expected daily, got %s", pref.Frequency)
	}
}

func TestCommunicationChannelsService_DeletePushChannel(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/communication_channels/push" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	err := svc.DeletePushChannel(context.Background())
	if err != nil {
		t.Fatalf("DeletePushChannel: %v", err)
	}
	if !called {
		t.Error("expected DELETE request")
	}
}
