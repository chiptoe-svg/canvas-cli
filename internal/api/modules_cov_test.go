package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestModulesService_Publish(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/modules/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		module, ok := body["module"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected 'module' key in body")
		}
		if module["published"] != true {
			t.Errorf("expected published=true, got %v", module["published"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Module{ID: 5, Name: "Published Module", Published: true})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewModulesService(client)
	module, err := svc.Publish(context.Background(), 10, 5)
	if err != nil {
		t.Fatalf("Publish: %v", err)
	}
	if module.ID != 5 {
		t.Errorf("expected ID 5, got %d", module.ID)
	}
	if !module.Published {
		t.Error("expected module to be published")
	}
}

func TestModulesService_Publish_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewModulesService(client)
	_, err := svc.Publish(context.Background(), 10, 5)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestModulesService_Unpublish(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/modules/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		module, ok := body["module"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected 'module' key in body")
		}
		if module["published"] != false {
			t.Errorf("expected published=false, got %v", module["published"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Module{ID: 5, Name: "Unpublished Module", Published: false})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewModulesService(client)
	module, err := svc.Unpublish(context.Background(), 10, 5)
	if err != nil {
		t.Fatalf("Unpublish: %v", err)
	}
	if module.ID != 5 {
		t.Errorf("expected ID 5, got %d", module.ID)
	}
	if module.Published {
		t.Error("expected module to be unpublished")
	}
}
