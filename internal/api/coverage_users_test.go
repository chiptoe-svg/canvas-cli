package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestUsersService_GetCurrentUser_Error covers the error path in GetCurrentUser
func TestUsersService_GetCurrentUser_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"errors":[{"message":"unauthorized"}]}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewUsersService(client)

	_, err := service.GetCurrentUser(context.Background())
	if err == nil {
		t.Error("Expected error from GetCurrentUser on 401, got nil")
	}
}

// TestUsersService_Get_Error covers the error path in Get
func TestUsersService_Get_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"user not found"}]}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewUsersService(client)

	_, err := service.Get(context.Background(), 9999, nil)
	if err == nil {
		t.Error("Expected error from Get on 404, got nil")
	}
}

// TestUsersService_List_Error covers the error path in List
func TestUsersService_List_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors":[{"message":"server error"}]}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewUsersService(client)

	_, err := service.List(context.Background(), 1, nil)
	if err == nil {
		t.Error("Expected error from List on 500, got nil")
	}
}

// TestUsersService_List_WithNilOptions ensures nil opts doesn't panic
func TestUsersService_List_WithNilOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewUsersService(client)

	users, err := service.List(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(users) != 0 {
		t.Errorf("Expected empty slice, got %d users", len(users))
	}
}

// TestUsersService_ListCourseUsers_WithEnrollmentState exercises EnrollmentState path
func TestUsersService_ListCourseUsers_WithEnrollmentState(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		state := r.URL.Query().Get("enrollment_state[]")
		if state != "active" {
			t.Errorf("Expected enrollment_state[] 'active', got %q", state)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": 1, "name": "Active Student"}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewUsersService(client)

	opts := &ListUsersOptions{
		EnrollmentState: "active",
	}
	users, err := service.ListCourseUsers(context.Background(), 10, opts)
	if err != nil {
		t.Fatalf("ListCourseUsers failed: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}
}

// TestUsersService_ListCourseUsers_WithPageOptions exercises Page and PerPage paths
func TestUsersService_ListCourseUsers_WithPageOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		page := r.URL.Query().Get("page")
		if page != "2" {
			t.Errorf("Expected page '2', got %q", page)
		}
		perPage := r.URL.Query().Get("per_page")
		if perPage != "25" {
			t.Errorf("Expected per_page '25', got %q", perPage)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": 5, "name": "Paged User"}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewUsersService(client)

	opts := &ListUsersOptions{
		Page:    2,
		PerPage: 25,
	}
	users, err := service.ListCourseUsers(context.Background(), 10, opts)
	if err != nil {
		t.Fatalf("ListCourseUsers failed: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}
}

// TestUsersService_ListCourseUsers_WithIncludes exercises the Include slice path
func TestUsersService_ListCourseUsers_WithIncludes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		includes := r.URL.Query()["include[]"]
		if len(includes) == 0 {
			t.Error("Expected include[] parameters")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": 3, "name": "User With Email", "email": "user@example.com"}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewUsersService(client)

	opts := &ListUsersOptions{
		Include: []string{"email", "avatar_url"},
	}
	users, err := service.ListCourseUsers(context.Background(), 10, opts)
	if err != nil {
		t.Fatalf("ListCourseUsers failed: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}
}

// TestUsersService_ListCourseUsers_Error covers the error path in ListCourseUsers
func TestUsersService_ListCourseUsers_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"errors":[{"message":"forbidden"}]}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewUsersService(client)

	_, err := service.ListCourseUsers(context.Background(), 10, nil)
	if err == nil {
		t.Error("Expected error from ListCourseUsers on 403, got nil")
	}
}

// TestUsersService_ListCourseUsers_NilOptions confirms nil opts works
func TestUsersService_ListCourseUsers_NilOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewUsersService(client)

	users, err := service.ListCourseUsers(context.Background(), 10, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(users) != 0 {
		t.Errorf("Expected empty slice, got %d users", len(users))
	}
}

// TestUsersService_Search_Error covers the error path in Search
func TestUsersService_Search_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors":[{"message":"server error"}]}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewUsersService(client)

	_, err := service.Search(context.Background(), "john")
	if err == nil {
		t.Error("Expected error from Search on 500, got nil")
	}
}

// TestUsersService_Create_Error covers the error path in Create
func TestUsersService_Create_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"errors":[{"message":"invalid user data"}]}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewUsersService(client)

	params := &CreateUserParams{Name: "Test User"}
	_, err := service.Create(context.Background(), 1, params)
	if err == nil {
		t.Error("Expected error from Create on 422, got nil")
	}
}

// TestUsersService_Update_Error covers the error path in Update
func TestUsersService_Update_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"user not found"}]}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewUsersService(client)

	params := &UpdateUserParams{Name: "New Name"}
	_, err := service.Update(context.Background(), 9999, params)
	if err == nil {
		t.Error("Expected error from Update on 404, got nil")
	}
}

// TestUsersService_Update_AvatarURLOnly exercises the avatar.URL-only path
func TestUsersService_Update_AvatarURLOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 55, "name": "Avatar User"}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewUsersService(client)

	// Token is empty, URL is set — exercises the URL-only branch inside avatar
	params := &UpdateUserParams{
		Name:   "Avatar User",
		Avatar: &AvatarParams{URL: "https://example.com/pic.jpg"},
	}
	user, err := service.Update(context.Background(), 55, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if user.ID != 55 {
		t.Errorf("Expected user ID 55, got %d", user.ID)
	}
}
