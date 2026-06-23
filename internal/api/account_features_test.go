package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountFeaturesService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/features" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]AccountFeature{
			{Feature: "peer_reviews", DisplayName: "Peer Reviews", AppliesTo: "Course"},
			{Feature: "analytics_2", DisplayName: "Analytics 2.0", AppliesTo: "Course"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountFeaturesService(client)

	features, err := svc.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(features) != 2 {
		t.Errorf("expected 2 features, got %d", len(features))
	}
	if features[0].Feature != "peer_reviews" {
		t.Errorf("expected peer_reviews, got %s", features[0].Feature)
	}
}

func TestAccountFeaturesService_ListEnabled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/features/enabled" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{"analytics_2", "new_gradebook"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountFeaturesService(client)

	features, err := svc.ListEnabled(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListEnabled: %v", err)
	}
	if len(features) != 2 {
		t.Errorf("expected 2 features, got %d", len(features))
	}
}

func TestAccountFeaturesService_GetFlag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/features/flags/analytics_2" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AccountFeatureFlag{
			Feature: "analytics_2",
			State:   "on",
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountFeaturesService(client)

	flag, err := svc.GetFlag(context.Background(), 1, "analytics_2")
	if err != nil {
		t.Fatalf("GetFlag: %v", err)
	}
	if flag.Feature != "analytics_2" {
		t.Errorf("expected analytics_2, got %s", flag.Feature)
	}
	if flag.State != "on" {
		t.Errorf("expected on, got %s", flag.State)
	}
}

func TestAccountFeaturesService_SetFlag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/features/flags/analytics_2" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AccountFeatureFlag{
			Feature: "analytics_2",
			State:   "off",
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountFeaturesService(client)

	flag, err := svc.SetFlag(context.Background(), 1, "analytics_2", "off")
	if err != nil {
		t.Fatalf("SetFlag: %v", err)
	}
	if flag.State != "off" {
		t.Errorf("expected off, got %s", flag.State)
	}
}

func TestAccountFeaturesService_DeleteFlag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/features/flags/analytics_2" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountFeaturesService(client)

	err := svc.DeleteFlag(context.Background(), 1, "analytics_2")
	if err != nil {
		t.Fatalf("DeleteFlag: %v", err)
	}
}

func TestAccountFeaturesService_GetSettings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/settings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AccountSettings{
			RestrictStudentPastView: true,
			LockAllAnnouncements:    false,
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountFeaturesService(client)

	settings, err := svc.GetSettings(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetSettings: %v", err)
	}
	if !settings.RestrictStudentPastView {
		t.Error("expected RestrictStudentPastView to be true")
	}
}

func TestAccountFeaturesService_GetPermissions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/permissions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{
			"manage_courses": true,
			"manage_users":   false,
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountFeaturesService(client)

	perms, err := svc.GetPermissions(context.Background(), 1, []string{"manage_courses", "manage_users"})
	if err != nil {
		t.Fatalf("GetPermissions: %v", err)
	}
	if !perms["manage_courses"] {
		t.Error("expected manage_courses to be true")
	}
	if perms["manage_users"] {
		t.Error("expected manage_users to be false")
	}
}

func TestAccountFeaturesService_GetScopes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/scopes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{
			map[string]interface{}{"resource": "courses", "verb": "GET"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountFeaturesService(client)

	scopes, err := svc.GetScopes(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetScopes: %v", err)
	}
	if len(scopes) != 1 {
		t.Errorf("expected 1 scope, got %d", len(scopes))
	}
}
