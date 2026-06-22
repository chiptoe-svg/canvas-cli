package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCSPSettingsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/csp_settings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		settings := CSPSettings{
			Status:  "enabled",
			Domains: []string{"example.com", "canvas.com"},
			Locked:  false,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCSPSettingsService(client)

	settings, err := svc.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if settings.Status != "enabled" {
		t.Errorf("expected status enabled, got %s", settings.Status)
	}
	if len(settings.Domains) != 2 {
		t.Errorf("expected 2 domains, got %d", len(settings.Domains))
	}
}

func TestCSPSettingsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/csp_settings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		settings := CSPSettings{
			Status: "disabled",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCSPSettingsService(client)

	body := &CSPSettings{Status: "disabled"}
	settings, err := svc.Update(context.Background(), 1, body)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if settings.Status != "disabled" {
		t.Errorf("expected status disabled, got %s", settings.Status)
	}
}

func TestCSPSettingsService_RemoveDomains(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/csp_settings/domains" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		settings := CSPSettings{
			Status:  "enabled",
			Domains: []string{"canvas.com"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCSPSettingsService(client)

	settings, err := svc.RemoveDomains(context.Background(), 1, []string{"example.com"})
	if err != nil {
		t.Fatalf("RemoveDomains: %v", err)
	}
	if len(settings.Domains) != 1 {
		t.Errorf("expected 1 domain, got %d", len(settings.Domains))
	}
}

func TestCSPSettingsService_AddDomains(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/csp_settings/domains" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		settings := CSPSettings{
			Status:  "enabled",
			Domains: []string{"example.com", "canvas.com"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCSPSettingsService(client)

	settings, err := svc.AddDomains(context.Background(), 1, []string{"example.com"})
	if err != nil {
		t.Fatalf("AddDomains: %v", err)
	}
	if len(settings.Domains) != 2 {
		t.Errorf("expected 2 domains, got %d", len(settings.Domains))
	}
}

func TestCSPSettingsService_BatchAddDomains(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/csp_settings/domains/batch_create" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		settings := CSPSettings{
			Status:  "enabled",
			Domains: []string{"a.com", "b.com", "c.com"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCSPSettingsService(client)

	settings, err := svc.BatchAddDomains(context.Background(), 1, []string{"a.com", "b.com", "c.com"})
	if err != nil {
		t.Fatalf("BatchAddDomains: %v", err)
	}
	if len(settings.Domains) != 3 {
		t.Errorf("expected 3 domains, got %d", len(settings.Domains))
	}
}

func TestCSPSettingsService_Lock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/csp_settings/lock" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		settings := CSPSettings{
			Status: "enabled",
			Locked: true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCSPSettingsService(client)

	settings, err := svc.Lock(context.Background(), 1, true)
	if err != nil {
		t.Fatalf("Lock: %v", err)
	}
	if !settings.Locked {
		t.Error("expected locked to be true")
	}
}
