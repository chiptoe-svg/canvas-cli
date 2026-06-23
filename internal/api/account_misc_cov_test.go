package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountMiscService_NewService(t *testing.T) {
	client := newTestClient(t, "http://localhost")
	svc := NewAccountMiscService(client)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
	if svc.client != client {
		t.Error("expected client to be set")
	}
}

func TestAccountMiscService_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/search" {
			t.Errorf("expected /api/v1/accounts/search, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("name") != "Test" {
			t.Errorf("expected name=Test, got %s", r.URL.Query().Get("name"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Account{{ID: 1, Name: "Test Account"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountMiscService(client)
	accounts, err := svc.Search(context.Background(), "Test", 0)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(accounts) != 1 {
		t.Fatalf("expected 1 account, got %d", len(accounts))
	}
	if accounts[0].Name != "Test Account" {
		t.Errorf("expected 'Test Account', got %s", accounts[0].Name)
	}
}

func TestAccountMiscService_Search_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountMiscService(client)
	_, err := svc.Search(context.Background(), "test", 0)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountMiscService_GetBrandVariables(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/brand_variables" {
			t.Errorf("expected /api/v1/accounts/1/brand_variables, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"ic-link-color": "#0077CC"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountMiscService(client)
	result, err := svc.GetBrandVariables(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetBrandVariables: %v", err)
	}
	if result["ic-link-color"] != "#0077CC" {
		t.Errorf("expected color #0077CC, got %v", result["ic-link-color"])
	}
}

func TestAccountMiscService_GetBrandVariables_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountMiscService(client)
	_, err := svc.GetBrandVariables(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountMiscService_GetHelpLinks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/help_links" {
			t.Errorf("expected /api/v1/accounts/1/help_links, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"help_link_name": "Help"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountMiscService(client)
	result, err := svc.GetHelpLinks(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetHelpLinks: %v", err)
	}
	if result["help_link_name"] != "Help" {
		t.Errorf("expected help_link_name=Help, got %v", result["help_link_name"])
	}
}

func TestAccountMiscService_GetTermsOfService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/terms_of_service" {
			t.Errorf("expected /api/v1/accounts/1/terms_of_service, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"passive": false})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountMiscService(client)
	result, err := svc.GetTermsOfService(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTermsOfService: %v", err)
	}
	if result["passive"] != false {
		t.Errorf("unexpected passive value: %v", result["passive"])
	}
}

func TestAccountMiscService_GetGradingStandards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/grading_standards" {
			t.Errorf("expected /api/v1/accounts/1/grading_standards, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"id": float64(1), "title": "Letter Grade"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountMiscService(client)
	result, err := svc.GetGradingStandards(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetGradingStandards: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 grading standard, got %d", len(result))
	}
	if result[0]["title"] != "Letter Grade" {
		t.Errorf("expected title=Letter Grade, got %v", result[0]["title"])
	}
}
