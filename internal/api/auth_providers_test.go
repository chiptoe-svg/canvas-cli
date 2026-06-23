package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthProvidersService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/authentication_providers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		providers := []AuthenticationProvider{
			{ID: 10, AuthType: "saml", Position: 1, WorkflowState: "active"},
			{ID: 11, AuthType: "ldap", Position: 2, WorkflowState: "active"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(providers)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAuthProvidersService(client)

	providers, err := svc.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(providers) != 2 {
		t.Errorf("expected 2 providers, got %d", len(providers))
	}
	if providers[0].AuthType != "saml" {
		t.Errorf("expected auth_type saml, got %s", providers[0].AuthType)
	}
}

func TestAuthProvidersService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/authentication_providers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		provider := AuthenticationProvider{ID: 20, AuthType: "saml", WorkflowState: "active"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(provider)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAuthProvidersService(client)

	params := &AuthProviderCreateParams{AuthType: "saml", ClientID: "client123"}
	provider, err := svc.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if provider.ID != 20 {
		t.Errorf("expected ID 20, got %d", provider.ID)
	}
}

func TestAuthProvidersService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/authentication_providers/10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		provider := AuthenticationProvider{ID: 10, AuthType: "saml", WorkflowState: "active"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(provider)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAuthProvidersService(client)

	provider, err := svc.Get(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if provider.ID != 10 {
		t.Errorf("expected ID 10, got %d", provider.ID)
	}
}

func TestAuthProvidersService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/authentication_providers/10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		provider := AuthenticationProvider{ID: 10, AuthType: "saml", WorkflowState: "active"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(provider)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAuthProvidersService(client)

	params := &AuthProviderCreateParams{AuthType: "saml"}
	provider, err := svc.Update(context.Background(), 1, 10, params)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if provider.ID != 10 {
		t.Errorf("expected ID 10, got %d", provider.ID)
	}
}

func TestAuthProvidersService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/authentication_providers/10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAuthProvidersService(client)

	err := svc.Delete(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestAuthProvidersService_Restore(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/authentication_providers/10/restore" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		provider := AuthenticationProvider{ID: 10, AuthType: "saml", WorkflowState: "active"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(provider)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAuthProvidersService(client)

	provider, err := svc.Restore(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if provider.ID != 10 {
		t.Errorf("expected ID 10, got %d", provider.ID)
	}
}

func TestAuthProvidersService_ForcePasswordReset(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/authentication_providers/force_password_reset" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAuthProvidersService(client)

	err := svc.ForcePasswordReset(context.Background(), 1)
	if err != nil {
		t.Fatalf("ForcePasswordReset: %v", err)
	}
}

func TestAuthProvidersService_GetSSOSettings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/sso_settings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		settings := SSOSettings{
			LoginHandleName:   "Email",
			ChangePasswordURL: "https://example.com/change",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAuthProvidersService(client)

	settings, err := svc.GetSSOSettings(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetSSOSettings: %v", err)
	}
	if settings.LoginHandleName != "Email" {
		t.Errorf("expected LoginHandleName Email, got %s", settings.LoginHandleName)
	}
}

func TestAuthProvidersService_UpdateSSOSettings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/sso_settings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		settings := SSOSettings{
			LoginHandleName:  "Username",
			AuthDiscoveryURL: "https://example.com/discover",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAuthProvidersService(client)

	body := &SSOSettings{LoginHandleName: "Username"}
	settings, err := svc.UpdateSSOSettings(context.Background(), 1, body)
	if err != nil {
		t.Fatalf("UpdateSSOSettings: %v", err)
	}
	if settings.LoginHandleName != "Username" {
		t.Errorf("expected LoginHandleName Username, got %s", settings.LoginHandleName)
	}
}
