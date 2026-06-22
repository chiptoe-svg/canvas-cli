package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewGroupsLTIService(t *testing.T) {
	client := newTestClient(t, "https://canvas.example.com")
	svc := NewGroupsLTIService(client)
	if svc == nil {
		t.Fatal("expected non-nil LTI service")
	}
	if svc.client != client {
		t.Error("expected client to be set")
	}
}

func TestGroupsLTIService_GetNamesAndRoles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/lti/groups/5/names_and_roles" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LTINamesAndRolesResponse{
			ID: "https://canvas.example.com/api/lti/groups/5/names_and_roles",
			Members: []LTIGroupMember{
				{Status: "Active", Name: "Alice", UserID: "1"},
				{Status: "Active", Name: "Bob", UserID: "2"},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsLTIService(client)

	resp, err := svc.GetNamesAndRoles(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetNamesAndRoles: %v", err)
	}
	if len(resp.Members) != 2 {
		t.Errorf("expected 2 members, got %d", len(resp.Members))
	}
	if resp.Members[0].Name != "Alice" {
		t.Errorf("expected first member 'Alice', got %q", resp.Members[0].Name)
	}
}

func TestGroupsLTIService_ListLTIUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/lti/groups/5/users" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]LTIUser{
			{ID: "u1", Name: "Alice"},
			{ID: "u2", Name: "Bob"},
			{ID: "u3", Name: "Charlie"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsLTIService(client)

	users, err := svc.ListLTIUsers(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListLTIUsers: %v", err)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}
