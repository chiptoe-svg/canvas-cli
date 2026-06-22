package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCollaborationsService_ListForCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/5/collaborations" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		colls := []Collaboration{
			{ID: 1, Title: "Group Doc", CollaborationType: "google_docs"},
		}
		if err := json.NewEncoder(w).Encode(colls); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCollaborationsService(client)

	colls, err := svc.ListForCourse(context.Background(), 5, nil)
	if err != nil {
		t.Fatalf("ListForCourse failed: %v", err)
	}
	if len(colls) != 1 || colls[0].ID != 1 {
		t.Errorf("unexpected result: %+v", colls)
	}
}

func TestCollaborationsService_ListForGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/3/collaborations" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		colls := []Collaboration{{ID: 7, Title: "Shared Sheet"}}
		if err := json.NewEncoder(w).Encode(colls); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCollaborationsService(client)

	colls, err := svc.ListForGroup(context.Background(), 3, nil)
	if err != nil {
		t.Fatalf("ListForGroup failed: %v", err)
	}
	if len(colls) != 1 || colls[0].ID != 7 {
		t.Errorf("unexpected result: %+v", colls)
	}
}

func TestCollaborationsService_ListMembers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/collaborations/11/members" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		members := []Collaborator{
			{ID: 100, Type: "user", Name: "Alice"},
			{ID: 101, Type: "user", Name: "Bob"},
		}
		if err := json.NewEncoder(w).Encode(members); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCollaborationsService(client)

	members, err := svc.ListMembers(context.Background(), 11, nil)
	if err != nil {
		t.Fatalf("ListMembers failed: %v", err)
	}
	if len(members) != 2 {
		t.Errorf("expected 2 members, got %d", len(members))
	}
}
