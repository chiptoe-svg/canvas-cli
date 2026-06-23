package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBlueprintService_ListSubscriptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/blueprint_subscriptions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]BlueprintSubscription{{ID: 1, TemplateID: 5}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBlueprintService(client)
	subs, err := svc.ListSubscriptions(context.Background(), 10)
	if err != nil {
		t.Fatalf("ListSubscriptions: %v", err)
	}
	if len(subs) != 1 {
		t.Fatalf("expected 1 subscription, got %d", len(subs))
	}
	if subs[0].TemplateID != 5 {
		t.Errorf("expected TemplateID=5, got %d", subs[0].TemplateID)
	}
}

func TestBlueprintService_ListSubscriptions_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBlueprintService(client)
	_, err := svc.ListSubscriptions(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBlueprintService_ListSubscriptionMigrations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/blueprint_subscriptions/1/migrations" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]BlueprintMigration{{ID: 100, TemplateID: 5}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBlueprintService(client)
	migrations, err := svc.ListSubscriptionMigrations(context.Background(), 10, 1)
	if err != nil {
		t.Fatalf("ListSubscriptionMigrations: %v", err)
	}
	if len(migrations) != 1 {
		t.Fatalf("expected 1 migration, got %d", len(migrations))
	}
	if migrations[0].ID != 100 {
		t.Errorf("expected ID=100, got %d", migrations[0].ID)
	}
}

func TestBlueprintService_GetSubscriptionMigration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/blueprint_subscriptions/1/migrations/100" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(BlueprintMigration{ID: 100, TemplateID: 5})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBlueprintService(client)
	migration, err := svc.GetSubscriptionMigration(context.Background(), 10, 1, 100)
	if err != nil {
		t.Fatalf("GetSubscriptionMigration: %v", err)
	}
	if migration.ID != 100 {
		t.Errorf("expected ID=100, got %d", migration.ID)
	}
}

func TestBlueprintService_GetSubscriptionMigration_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBlueprintService(client)
	_, err := svc.GetSubscriptionMigration(context.Background(), 10, 1, 100)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBlueprintService_GetMigrationDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/blueprint_templates/default/migrations/100/details" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"asset_type": "Assignment", "asset_id": float64(1)}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBlueprintService(client)
	details, err := svc.GetMigrationDetails(context.Background(), 10, "", 100)
	if err != nil {
		t.Fatalf("GetMigrationDetails: %v", err)
	}
	if len(details) != 1 {
		t.Fatalf("expected 1 detail, got %d", len(details))
	}
	if details[0]["asset_type"] != "Assignment" {
		t.Errorf("expected asset_type=Assignment, got %v", details[0]["asset_type"])
	}
}

func TestBlueprintService_GetSubscriptionMigrationDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/blueprint_subscriptions/1/migrations/100/details" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"asset_type": "Quiz", "asset_id": float64(2)}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBlueprintService(client)
	details, err := svc.GetSubscriptionMigrationDetails(context.Background(), 10, 1, 100)
	if err != nil {
		t.Fatalf("GetSubscriptionMigrationDetails: %v", err)
	}
	if len(details) != 1 {
		t.Fatalf("expected 1 detail, got %d", len(details))
	}
	if details[0]["asset_type"] != "Quiz" {
		t.Errorf("expected asset_type=Quiz, got %v", details[0]["asset_type"])
	}
}
