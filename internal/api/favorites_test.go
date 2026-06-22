package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFavoritesService_ListCourses(t *testing.T) {
	want := []Course{{ID: 10, Name: "Favorite Course"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/favorites/courses" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFavoritesService(client)
	got, err := svc.ListCourses(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d courses, want %d", len(got), len(want))
	}
	if got[0].ID != want[0].ID {
		t.Errorf("got ID %d, want %d", got[0].ID, want[0].ID)
	}
}

func TestFavoritesService_AddCourse(t *testing.T) {
	want := &Course{ID: 20, Name: "Added Course"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/favorites/courses/20" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFavoritesService(client)
	got, err := svc.AddCourse(context.Background(), 20)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != want.ID {
		t.Errorf("got ID %d, want %d", got.ID, want.ID)
	}
}

func TestFavoritesService_RemoveCourse(t *testing.T) {
	want := &Course{ID: 30, Name: "Removed Course"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/favorites/courses/30" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFavoritesService(client)
	got, err := svc.RemoveCourse(context.Background(), 30)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != want.ID {
		t.Errorf("got ID %d, want %d", got.ID, want.ID)
	}
}

func TestFavoritesService_ResetCourses(t *testing.T) {
	want := []Course{{ID: 1, Name: "Default Course 1"}, {ID: 2, Name: "Default Course 2"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/favorites/courses" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFavoritesService(client)
	got, err := svc.ResetCourses(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d courses, want %d", len(got), len(want))
	}
	if got[0].ID != want[0].ID {
		t.Errorf("got ID %d, want %d", got[0].ID, want[0].ID)
	}
}

func TestFavoritesService_ListGroups(t *testing.T) {
	want := []Group{{ID: 100, Name: "Favorite Group"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/favorites/groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFavoritesService(client)
	got, err := svc.ListGroups(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d groups, want %d", len(got), len(want))
	}
	if got[0].ID != want[0].ID {
		t.Errorf("got ID %d, want %d", got[0].ID, want[0].ID)
	}
}

func TestFavoritesService_AddGroup(t *testing.T) {
	want := &Group{ID: 200, Name: "Added Group"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/favorites/groups/200" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFavoritesService(client)
	got, err := svc.AddGroup(context.Background(), 200)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != want.ID {
		t.Errorf("got ID %d, want %d", got.ID, want.ID)
	}
}

func TestFavoritesService_RemoveGroup(t *testing.T) {
	want := &Group{ID: 300, Name: "Removed Group"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/favorites/groups/300" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFavoritesService(client)
	got, err := svc.RemoveGroup(context.Background(), 300)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != want.ID {
		t.Errorf("got ID %d, want %d", got.ID, want.ID)
	}
}

func TestFavoritesService_ResetGroups(t *testing.T) {
	want := []Group{{ID: 1, Name: "Default Group 1"}, {ID: 2, Name: "Default Group 2"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/favorites/groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFavoritesService(client)
	got, err := svc.ResetGroups(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d groups, want %d", len(got), len(want))
	}
	if got[0].ID != want[0].ID {
		t.Errorf("got ID %d, want %d", got[0].ID, want[0].ID)
	}
}

func TestNewFavoritesService(t *testing.T) {
	client := &Client{}
	svc := NewFavoritesService(client)
	if svc == nil {
		t.Fatal("NewFavoritesService returned nil")
	}
	if svc.client != client {
		t.Error("NewFavoritesService did not set client correctly")
	}
}
