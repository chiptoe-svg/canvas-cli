package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCourseFeaturesService_List(t *testing.T) {
	features := []Feature{
		{Feature: "new_quizzes", DisplayName: "New Quizzes"},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/5/features" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(features); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseFeaturesService(client)
	got, err := svc.List(context.Background(), 5)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 1 || got[0].Feature != "new_quizzes" {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestCourseFeaturesService_ListEnabled(t *testing.T) {
	enabled := []string{"new_quizzes", "lti_advantage"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/5/features/enabled" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(enabled); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseFeaturesService(client)
	got, err := svc.ListEnabled(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListEnabled: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2, got %d", len(got))
	}
}

func TestCourseFeaturesService_GetFlag(t *testing.T) {
	want := FeatureFlag{Feature: "new_quizzes", State: "on"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/5/features/flags/new_quizzes" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseFeaturesService(client)
	got, err := svc.GetFlag(context.Background(), 5, "new_quizzes")
	if err != nil {
		t.Fatalf("GetFlag: %v", err)
	}
	if got.State != "on" {
		t.Errorf("want state 'on', got %q", got.State)
	}
}

func TestCourseFeaturesService_SetFlag(t *testing.T) {
	want := FeatureFlag{Feature: "new_quizzes", State: "off"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/courses/5/features/flags/new_quizzes" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseFeaturesService(client)
	got, err := svc.SetFlag(context.Background(), 5, "new_quizzes", "off")
	if err != nil {
		t.Fatalf("SetFlag: %v", err)
	}
	if got.State != "off" {
		t.Errorf("want 'off', got %q", got.State)
	}
}

func TestCourseFeaturesService_DeleteFlag(t *testing.T) {
	want := FeatureFlag{Feature: "new_quizzes", State: "allowed"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/courses/5/features/flags/new_quizzes" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseFeaturesService(client)
	got, err := svc.DeleteFlag(context.Background(), 5, "new_quizzes")
	if err != nil {
		t.Fatalf("DeleteFlag: %v", err)
	}
	if got.State != "allowed" {
		t.Errorf("want 'allowed', got %q", got.State)
	}
}
