package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserFeaturesService_List(t *testing.T) {
	want := []Feature{
		{Feature: "new_gradebook", AppliesTo: "User", FeatureFlag: FeatureFlag{Feature: "new_gradebook", State: "allowed"}},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/features" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUserFeaturesService(client)
	got, err := svc.List(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d features, want %d", len(got), len(want))
	}
	if got[0].Feature != want[0].Feature {
		t.Errorf("got Feature %q, want %q", got[0].Feature, want[0].Feature)
	}
}

func TestUserFeaturesService_ListEnabled(t *testing.T) {
	want := []string{"student_planner", "new_gradebook"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/features/enabled" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUserFeaturesService(client)
	got, err := svc.ListEnabled(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d features, want %d", len(got), len(want))
	}
	if got[0] != want[0] {
		t.Errorf("got feature name %q, want %q", got[0], want[0])
	}
}

func TestUserFeaturesService_GetFlag(t *testing.T) {
	want := &FeatureFlag{Feature: "new_gradebook", State: "allowed", Locked: false}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/features/flags/new_gradebook" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUserFeaturesService(client)
	got, err := svc.GetFlag(context.Background(), 42, "new_gradebook")
	if err != nil {
		t.Fatal(err)
	}
	if got.Feature != want.Feature {
		t.Errorf("got Feature %q, want %q", got.Feature, want.Feature)
	}
	if got.State != want.State {
		t.Errorf("got State %q, want %q", got.State, want.State)
	}
}

func TestUserFeaturesService_SetFlag(t *testing.T) {
	want := &FeatureFlag{Feature: "new_gradebook", State: "on"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/features/flags/new_gradebook" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body["state"] != "on" {
			t.Errorf("expected state=on, got %q", body["state"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUserFeaturesService(client)
	got, err := svc.SetFlag(context.Background(), 42, "new_gradebook", "on")
	if err != nil {
		t.Fatal(err)
	}
	if got.State != want.State {
		t.Errorf("got State %q, want %q", got.State, want.State)
	}
}

func TestUserFeaturesService_DeleteFlag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/features/flags/new_gradebook" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}")) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUserFeaturesService(client)
	err := svc.DeleteFlag(context.Background(), 42, "new_gradebook")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewUserFeaturesService(t *testing.T) {
	client := &Client{}
	svc := NewUserFeaturesService(client)
	if svc == nil {
		t.Fatal("NewUserFeaturesService returned nil")
	}
	if svc.client != client {
		t.Error("NewUserFeaturesService did not set client correctly")
	}
}
