package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBookmarksService_List(t *testing.T) {
	want := []Bookmark{{ID: 1, Name: "Test Bookmark", URL: "https://example.com"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/bookmarks" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBookmarksService(client)
	got, err := svc.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d bookmarks, want %d", len(got), len(want))
	}
	if got[0].ID != want[0].ID {
		t.Errorf("got ID %d, want %d", got[0].ID, want[0].ID)
	}
	if got[0].Name != want[0].Name {
		t.Errorf("got Name %q, want %q", got[0].Name, want[0].Name)
	}
}

func TestBookmarksService_Create(t *testing.T) {
	want := &Bookmark{ID: 2, Name: "New Bookmark", URL: "https://canvas.example.com", Position: 1}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/bookmarks" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBookmarksService(client)
	params := CreateBookmarkParams{Name: "New Bookmark", URL: "https://canvas.example.com", Position: 1}
	got, err := svc.Create(context.Background(), params)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != want.ID {
		t.Errorf("got ID %d, want %d", got.ID, want.ID)
	}
	if got.Name != want.Name {
		t.Errorf("got Name %q, want %q", got.Name, want.Name)
	}
}

func TestBookmarksService_Get(t *testing.T) {
	want := &Bookmark{ID: 5, Name: "Specific Bookmark", URL: "https://example.com/5"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/bookmarks/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBookmarksService(client)
	got, err := svc.Get(context.Background(), 5)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != want.ID {
		t.Errorf("got ID %d, want %d", got.ID, want.ID)
	}
}

func TestBookmarksService_Update(t *testing.T) {
	want := &Bookmark{ID: 3, Name: "Updated Bookmark", URL: "https://updated.example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/bookmarks/3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBookmarksService(client)
	params := UpdateBookmarkParams{Name: "Updated Bookmark", URL: "https://updated.example.com"}
	got, err := svc.Update(context.Background(), 3, params)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != want.ID {
		t.Errorf("got ID %d, want %d", got.ID, want.ID)
	}
	if got.Name != want.Name {
		t.Errorf("got Name %q, want %q", got.Name, want.Name)
	}
}

func TestBookmarksService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/bookmarks/7" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}")) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBookmarksService(client)
	err := svc.Delete(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewBookmarksService(t *testing.T) {
	client := &Client{}
	svc := NewBookmarksService(client)
	if svc == nil {
		t.Fatal("NewBookmarksService returned nil")
	}
	if svc.client != client {
		t.Error("NewBookmarksService did not set client correctly")
	}
}
