package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContentSharesService_Create(t *testing.T) {
	want := &ContentShare{ID: 1, Name: "Assignment 1", ContentType: "assignment", UserID: 10}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/10/content_shares" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentSharesService(client)
	params := CreateContentShareParams{
		ReceiverIDs: []int64{20, 30},
		ContentType: "assignment",
		ContentID:   100,
	}
	got, err := svc.Create(context.Background(), 10, params)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != want.ID {
		t.Errorf("got ID %d, want %d", got.ID, want.ID)
	}
	if got.ContentType != want.ContentType {
		t.Errorf("got ContentType %q, want %q", got.ContentType, want.ContentType)
	}
}

func TestContentSharesService_ListSent(t *testing.T) {
	want := []ContentShare{{ID: 1, Name: "Sent Share", UserID: 10}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/10/content_shares/sent" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentSharesService(client)
	got, err := svc.ListSent(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d shares, want %d", len(got), len(want))
	}
	if got[0].ID != want[0].ID {
		t.Errorf("got ID %d, want %d", got[0].ID, want[0].ID)
	}
}

func TestContentSharesService_ListReceived(t *testing.T) {
	want := []ContentShare{{ID: 2, Name: "Received Share", UserID: 10, ReadState: "unread"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/10/content_shares/received" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentSharesService(client)
	got, err := svc.ListReceived(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d shares, want %d", len(got), len(want))
	}
	if got[0].ReadState != want[0].ReadState {
		t.Errorf("got ReadState %q, want %q", got[0].ReadState, want[0].ReadState)
	}
}

func TestContentSharesService_Get(t *testing.T) {
	want := &ContentShare{ID: 5, Name: "Specific Share", UserID: 10}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/10/content_shares/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentSharesService(client)
	got, err := svc.Get(context.Background(), 10, 5)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != want.ID {
		t.Errorf("got ID %d, want %d", got.ID, want.ID)
	}
}

func TestContentSharesService_Update(t *testing.T) {
	want := &ContentShare{ID: 5, Name: "Updated Share", ReadState: "read", UserID: 10}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/10/content_shares/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentSharesService(client)
	params := UpdateContentShareParams{ReadState: "read"}
	got, err := svc.Update(context.Background(), 10, 5, params)
	if err != nil {
		t.Fatal(err)
	}
	if got.ReadState != want.ReadState {
		t.Errorf("got ReadState %q, want %q", got.ReadState, want.ReadState)
	}
}

func TestContentSharesService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/10/content_shares/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}")) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentSharesService(client)
	err := svc.Delete(context.Background(), 10, 5)
	if err != nil {
		t.Fatal(err)
	}
}

func TestContentSharesService_GetUnreadCount(t *testing.T) {
	want := &ContentShareUnreadCount{UnreadCount: 3}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/10/content_shares/unread_count" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentSharesService(client)
	got, err := svc.GetUnreadCount(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if got.UnreadCount != want.UnreadCount {
		t.Errorf("got UnreadCount %d, want %d", got.UnreadCount, want.UnreadCount)
	}
}

func TestNewContentSharesService(t *testing.T) {
	client := &Client{}
	svc := NewContentSharesService(client)
	if svc == nil {
		t.Fatal("NewContentSharesService returned nil")
	}
	if svc.client != client {
		t.Error("NewContentSharesService did not set client correctly")
	}
}
