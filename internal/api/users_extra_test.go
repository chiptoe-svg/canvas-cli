package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsersService_GetProfile(t *testing.T) {
	want := &UserProfile{ID: 42, Name: "Test User", PrimaryEmail: "test@example.com", LoginID: "testuser"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/profile" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	got, err := svc.GetProfile(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != want.ID {
		t.Errorf("got ID %d, want %d", got.ID, want.ID)
	}
	if got.Name != want.Name {
		t.Errorf("got Name %q, want %q", got.Name, want.Name)
	}
	if got.PrimaryEmail != want.PrimaryEmail {
		t.Errorf("got PrimaryEmail %q, want %q", got.PrimaryEmail, want.PrimaryEmail)
	}
}

func TestUsersService_GetSettings(t *testing.T) {
	want := &UserSettings{ManualMarkAsRead: true, CollapseGlobalNav: false}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/settings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	got, err := svc.GetSettings(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if got.ManualMarkAsRead != want.ManualMarkAsRead {
		t.Errorf("got ManualMarkAsRead %v, want %v", got.ManualMarkAsRead, want.ManualMarkAsRead)
	}
}

func TestUsersService_UpdateSettings(t *testing.T) {
	trueVal := true
	want := &UserSettings{ManualMarkAsRead: true, CollapseGlobalNav: true}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/settings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	params := UpdateUserSettingsParams{
		ManualMarkAsRead:  &trueVal,
		CollapseGlobalNav: &trueVal,
	}
	got, err := svc.UpdateSettings(context.Background(), 42, params)
	if err != nil {
		t.Fatal(err)
	}
	if got.ManualMarkAsRead != want.ManualMarkAsRead {
		t.Errorf("got ManualMarkAsRead %v, want %v", got.ManualMarkAsRead, want.ManualMarkAsRead)
	}
	if got.CollapseGlobalNav != want.CollapseGlobalNav {
		t.Errorf("got CollapseGlobalNav %v, want %v", got.CollapseGlobalNav, want.CollapseGlobalNav)
	}
}

func TestUsersService_ListLogins(t *testing.T) {
	want := []UserLogin{
		{ID: 1, UserID: 42, UniqueID: "testuser@example.com", WorkflowState: "active"},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/logins" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	got, err := svc.ListLogins(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d logins, want %d", len(got), len(want))
	}
	if got[0].UniqueID != want[0].UniqueID {
		t.Errorf("got UniqueID %q, want %q", got[0].UniqueID, want[0].UniqueID)
	}
}

func TestUsersService_GetColors(t *testing.T) {
	want := &UserColors{CustomColors: map[string]string{"course_1": "#E66000", "course_2": "#008EE2"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/colors" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	got, err := svc.GetColors(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if got.CustomColors["course_1"] != want.CustomColors["course_1"] {
		t.Errorf("got course_1 color %q, want %q", got.CustomColors["course_1"], want.CustomColors["course_1"])
	}
}

func TestUsersService_GetActivityStream(t *testing.T) {
	want := []ActivityStreamItem{
		{ID: 1, Type: "Message", Title: "New Message", ReadState: false},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/activity_stream" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	got, err := svc.GetActivityStream(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d items, want %d", len(got), len(want))
	}
	if got[0].ID != want[0].ID {
		t.Errorf("got ID %d, want %d", got[0].ID, want[0].ID)
	}
	if got[0].Type != want[0].Type {
		t.Errorf("got Type %q, want %q", got[0].Type, want[0].Type)
	}
}

func TestUsersService_GetTodo(t *testing.T) {
	want := []TodoItem{
		{Type: "submitting", ContextType: "Course", CourseID: 10},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/todo" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	got, err := svc.GetTodo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d todo items, want %d", len(got), len(want))
	}
	if got[0].Type != want[0].Type {
		t.Errorf("got Type %q, want %q", got[0].Type, want[0].Type)
	}
	if got[0].CourseID != want[0].CourseID {
		t.Errorf("got CourseID %d, want %d", got[0].CourseID, want[0].CourseID)
	}
}
