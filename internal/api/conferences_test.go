package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConferencesService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/conferences" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		confs := []Conference{
			{ID: 1, Title: "Lecture 1"},
			{ID: 2, Title: "Office Hours"},
		}
		if err := json.NewEncoder(w).Encode(confs); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewConferencesService(client)

	confs, err := svc.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(confs) != 2 {
		t.Errorf("expected 2 conferences, got %d", len(confs))
	}
}

func TestConferencesService_ListForCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/conferences" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		confs := []Conference{{ID: 5, Title: "Course Conf"}}
		if err := json.NewEncoder(w).Encode(confs); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewConferencesService(client)

	confs, err := svc.ListForCourse(context.Background(), 10, nil)
	if err != nil {
		t.Fatalf("ListForCourse failed: %v", err)
	}
	if len(confs) != 1 || confs[0].ID != 5 {
		t.Errorf("unexpected result: %+v", confs)
	}
}

func TestConferencesService_ListForGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/7/conferences" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		confs := []Conference{{ID: 9, Title: "Group Conf"}}
		if err := json.NewEncoder(w).Encode(confs); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewConferencesService(client)

	confs, err := svc.ListForGroup(context.Background(), 7, nil)
	if err != nil {
		t.Fatalf("ListForGroup failed: %v", err)
	}
	if len(confs) != 1 || confs[0].ID != 9 {
		t.Errorf("unexpected result: %+v", confs)
	}
}
