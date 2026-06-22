package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountLoginsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/logins" {
			t.Errorf("expected path /api/v1/accounts/1/logins, got %s", r.URL.Path)
		}

		logins := []Login{
			{ID: 1, UserID: 100, UniqueID: "alice@example.com", AccountID: 1},
			{ID: 2, UserID: 101, UniqueID: "bob@example.com", AccountID: 1},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logins)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewAccountLoginsService(client)

	logins, err := service.List(context.Background(), 1, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(logins) != 2 {
		t.Fatalf("expected 2 logins, got %d", len(logins))
	}

	if logins[0].UniqueID != "alice@example.com" {
		t.Errorf("expected unique_id 'alice@example.com', got '%s'", logins[0].UniqueID)
	}
}

func TestAccountLoginsService_ListWithUserID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/logins" {
			t.Errorf("expected path /api/v1/accounts/1/logins, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("user_id") != "100" {
			t.Errorf("expected user_id=100, got %s", r.URL.Query().Get("user_id"))
		}

		logins := []Login{
			{ID: 1, UserID: 100, UniqueID: "alice@example.com", AccountID: 1},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logins)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewAccountLoginsService(client)

	logins, err := service.List(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("List with user_id failed: %v", err)
	}

	if len(logins) != 1 {
		t.Fatalf("expected 1 login, got %d", len(logins))
	}
}

func TestAccountLoginsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/logins" {
			t.Errorf("expected path /api/v1/accounts/1/logins, got %s", r.URL.Path)
		}

		login := Login{ID: 10, UserID: 200, UniqueID: "charlie@example.com", AccountID: 1}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(login)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewAccountLoginsService(client)

	params := &LoginParams{UserID: 200, UniqueID: "charlie@example.com"}
	login, err := service.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if login.UniqueID != "charlie@example.com" {
		t.Errorf("expected unique_id 'charlie@example.com', got '%s'", login.UniqueID)
	}
}

func TestAccountLoginsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/logins/10" {
			t.Errorf("expected path /api/v1/accounts/1/logins/10, got %s", r.URL.Path)
		}

		login := Login{ID: 10, UserID: 200, UniqueID: "updated@example.com", AccountID: 1}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(login)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewAccountLoginsService(client)

	params := &LoginParams{UniqueID: "updated@example.com"}
	login, err := service.Update(context.Background(), 1, 10, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if login.UniqueID != "updated@example.com" {
		t.Errorf("expected unique_id 'updated@example.com', got '%s'", login.UniqueID)
	}
}
