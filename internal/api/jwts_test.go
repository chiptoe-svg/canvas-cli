package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJWTsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/jwts" || r.Method != http.MethodPost {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		jwt := JWT{Token: "eyJ.test.token"}
		if err := json.NewEncoder(w).Encode(jwt); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewJWTsService(client)

	jwt, err := svc.Create(context.Background(), "")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if jwt.Token != "eyJ.test.token" {
		t.Errorf("unexpected token: %s", jwt.Token)
	}
}

func TestJWTsService_Refresh(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/jwts/refresh" || r.Method != http.MethodPost {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		jwt := JWT{Token: "eyJ.refreshed.token"}
		if err := json.NewEncoder(w).Encode(jwt); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewJWTsService(client)

	jwt, err := svc.Refresh(context.Background(), "old-token")
	if err != nil {
		t.Fatalf("Refresh failed: %v", err)
	}
	if jwt.Token != "eyJ.refreshed.token" {
		t.Errorf("unexpected token: %s", jwt.Token)
	}
}
