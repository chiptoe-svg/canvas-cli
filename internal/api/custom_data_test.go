package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCustomDataService_Get(t *testing.T) {
	responseData := map[string]interface{}{"key": "value", "count": float64(42)}
	want := CustomDataResult{Data: responseData}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/123/custom_data" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		ns := r.URL.Query().Get("ns")
		if ns != "com.example.app" {
			t.Errorf("unexpected ns: %s", ns)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCustomDataService(client)
	got, err := svc.Get(context.Background(), 123, "com.example.app")
	if err != nil {
		t.Fatal(err)
	}
	if got["key"] != responseData["key"] {
		t.Errorf("got key %v, want %v", got["key"], responseData["key"])
	}
}

func TestCustomDataService_Get_NoNamespace(t *testing.T) {
	responseData := map[string]interface{}{"foo": "bar"}
	want := CustomDataResult{Data: responseData}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/456/custom_data" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("ns") != "" {
			t.Errorf("expected no ns param, got %q", r.URL.Query().Get("ns"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCustomDataService(client)
	got, err := svc.Get(context.Background(), 456, "")
	if err != nil {
		t.Fatal(err)
	}
	if got["foo"] != responseData["foo"] {
		t.Errorf("got foo %v, want %v", got["foo"], responseData["foo"])
	}
}

func TestCustomDataService_Set(t *testing.T) {
	data := map[string]interface{}{"setting": "enabled"}
	responseData := map[string]interface{}{"setting": "enabled"}
	want := CustomDataResult{Data: responseData}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/789/custom_data" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if body["ns"] != "com.example.ns" {
			t.Errorf("unexpected ns: %v", body["ns"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCustomDataService(client)
	got, err := svc.Set(context.Background(), 789, "com.example.ns", data)
	if err != nil {
		t.Fatal(err)
	}
	if got["setting"] != responseData["setting"] {
		t.Errorf("got setting %v, want %v", got["setting"], responseData["setting"])
	}
}

func TestCustomDataService_Delete(t *testing.T) {
	responseData := map[string]interface{}{"deleted": true}
	want := CustomDataResult{Data: responseData}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/321/custom_data" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCustomDataService(client)
	got, err := svc.Delete(context.Background(), 321, "com.example.ns")
	if err != nil {
		t.Fatal(err)
	}
	if got["deleted"] != responseData["deleted"] {
		t.Errorf("got deleted %v, want %v", got["deleted"], responseData["deleted"])
	}
}

func TestNewCustomDataService(t *testing.T) {
	client := &Client{}
	svc := NewCustomDataService(client)
	if svc == nil {
		t.Fatal("NewCustomDataService returned nil")
	}
	if svc.client != client {
		t.Error("NewCustomDataService did not set client correctly")
	}
}
