package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountExternalToolsService_AddRCEFavorite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/accounts/10/external_tools/rce_favorites/5" {
			t.Errorf("Expected path /api/v1/accounts/10/external_tools/rce_favorites/5, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewAccountExternalToolsService(client)

	err := service.AddRCEFavorite(context.Background(), 10, 5)
	if err != nil {
		t.Fatalf("AddRCEFavorite failed: %v", err)
	}
}

func TestAccountExternalToolsService_RemoveRCEFavorite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/accounts/10/external_tools/rce_favorites/5" {
			t.Errorf("Expected path /api/v1/accounts/10/external_tools/rce_favorites/5, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewAccountExternalToolsService(client)

	err := service.RemoveRCEFavorite(context.Background(), 10, 5)
	if err != nil {
		t.Fatalf("RemoveRCEFavorite failed: %v", err)
	}
}

func TestAccountExternalToolsService_AddTopNavFavorite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/accounts/10/external_tools/top_nav_favorites/7" {
			t.Errorf("Expected path /api/v1/accounts/10/external_tools/top_nav_favorites/7, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewAccountExternalToolsService(client)

	err := service.AddTopNavFavorite(context.Background(), 10, 7)
	if err != nil {
		t.Fatalf("AddTopNavFavorite failed: %v", err)
	}
}

func TestAccountExternalToolsService_RemoveTopNavFavorite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/accounts/10/external_tools/top_nav_favorites/7" {
			t.Errorf("Expected path /api/v1/accounts/10/external_tools/top_nav_favorites/7, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewAccountExternalToolsService(client)

	err := service.RemoveTopNavFavorite(context.Background(), 10, 7)
	if err != nil {
		t.Fatalf("RemoveTopNavFavorite failed: %v", err)
	}
}
