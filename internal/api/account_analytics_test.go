package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountAnalyticsService_GetTermActivity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/analytics/terms/2/activity" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{
			map[string]interface{}{"date": "2024-01-01", "views": 100},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountAnalyticsService(client)

	result, err := svc.GetTermActivity(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("GetTermActivity: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 item, got %d", len(result))
	}
}

func TestAccountAnalyticsService_GetTermGrades(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/analytics/terms/2/grades" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{
			map[string]interface{}{"score": 85.0, "count": 10},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountAnalyticsService(client)

	result, err := svc.GetTermGrades(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("GetTermGrades: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 item, got %d", len(result))
	}
}

func TestAccountAnalyticsService_GetTermStatistics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/analytics/terms/2/statistics" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"courses":  42,
			"students": 1200,
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountAnalyticsService(client)

	result, err := svc.GetTermStatistics(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("GetTermStatistics: %v", err)
	}
	if result["courses"] == nil {
		t.Error("expected courses key in result")
	}
}

func TestAccountAnalyticsService_GetTermStatsBySubaccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/analytics/terms/2/statistics_by_subaccount" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{
			map[string]interface{}{"id": 10, "courses": 5},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountAnalyticsService(client)

	result, err := svc.GetTermStatsBySubaccount(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("GetTermStatsBySubaccount: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 item, got %d", len(result))
	}
}

func TestAccountAnalyticsService_GetCompletedActivity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/analytics/completed/activity" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{
			map[string]interface{}{"date": "2024-01-01", "views": 50},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountAnalyticsService(client)

	result, err := svc.GetCompletedActivity(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCompletedActivity: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 item, got %d", len(result))
	}
}

func TestAccountAnalyticsService_GetCompletedStatistics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/analytics/completed/statistics" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"courses":  20,
			"students": 500,
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountAnalyticsService(client)

	result, err := svc.GetCompletedStatistics(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCompletedStatistics: %v", err)
	}
	if result["courses"] == nil {
		t.Error("expected courses key in result")
	}
}

func TestAccountAnalyticsService_GetCurrentStatsBySubaccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/analytics/current/statistics_by_subaccount" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{
			map[string]interface{}{"id": 10, "courses": 3},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountAnalyticsService(client)

	result, err := svc.GetCurrentStatsBySubaccount(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCurrentStatsBySubaccount: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 item, got %d", len(result))
	}
}
