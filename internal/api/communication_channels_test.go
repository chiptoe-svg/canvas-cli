package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommunicationChannelsService_List(t *testing.T) {
	want := []CommunicationChannel{
		{ID: 1, Address: "test@example.com", Type: "email", UserID: 10},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/10/communication_channels" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	got, err := svc.List(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d channels, want %d", len(got), len(want))
	}
	if got[0].ID != want[0].ID {
		t.Errorf("got ID %d, want %d", got[0].ID, want[0].ID)
	}
	if got[0].Address != want[0].Address {
		t.Errorf("got Address %q, want %q", got[0].Address, want[0].Address)
	}
}

func TestCommunicationChannelsService_Create(t *testing.T) {
	want := &CommunicationChannel{ID: 2, Address: "new@example.com", Type: "email", UserID: 10}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/10/communication_channels" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if _, ok := body["communication_channel"]; !ok {
			t.Error("expected communication_channel key in body")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	params := CreateCommunicationChannelParams{
		Address:          "new@example.com",
		Type:             "email",
		SkipConfirmation: true,
	}
	got, err := svc.Create(context.Background(), 10, params)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != want.ID {
		t.Errorf("got ID %d, want %d", got.ID, want.ID)
	}
	if got.Address != want.Address {
		t.Errorf("got Address %q, want %q", got.Address, want.Address)
	}
}

func TestCommunicationChannelsService_Delete(t *testing.T) {
	want := &CommunicationChannel{ID: 5, Address: "old@example.com", Type: "email", UserID: 10}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/10/communication_channels/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	got, err := svc.Delete(context.Background(), 10, 5)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != want.ID {
		t.Errorf("got ID %d, want %d", got.ID, want.ID)
	}
}

func TestCommunicationChannelsService_GetNotificationPreferences(t *testing.T) {
	want := &NotificationPreferences{
		NotificationPreferences: []NotificationPreference{
			{Frequency: "immediately", Notification: "new_announcement", Category: "announcement"},
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/10/communication_channels/5/notification_preferences" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommunicationChannelsService(client)
	got, err := svc.GetNotificationPreferences(context.Background(), 10, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.NotificationPreferences) != len(want.NotificationPreferences) {
		t.Fatalf("got %d prefs, want %d", len(got.NotificationPreferences), len(want.NotificationPreferences))
	}
	if got.NotificationPreferences[0].Frequency != want.NotificationPreferences[0].Frequency {
		t.Errorf("got Frequency %q, want %q",
			got.NotificationPreferences[0].Frequency,
			want.NotificationPreferences[0].Frequency)
	}
}

func TestNewCommunicationChannelsService(t *testing.T) {
	client := &Client{}
	svc := NewCommunicationChannelsService(client)
	if svc == nil {
		t.Fatal("NewCommunicationChannelsService returned nil")
	}
	if svc.client != client {
		t.Error("NewCommunicationChannelsService did not set client correctly")
	}
}
