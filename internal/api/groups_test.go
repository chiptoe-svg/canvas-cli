package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGroupsService_ListCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/groups" {
			t.Errorf("expected /api/v1/courses/123/groups, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Group{
			{ID: 1, Name: "Group 1", MembersCount: 5},
			{ID: 2, Name: "Group 2", MembersCount: 3},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	groups, err := service.ListCourse(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}

	if groups[0].Name != "Group 1" {
		t.Errorf("expected 'Group 1', got %s", groups[0].Name)
	}
}

func TestGroupsService_ListAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/groups" {
			t.Errorf("expected /api/v1/accounts/1/groups, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Group{
			{ID: 1, Name: "Account Group"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	groups, err := service.ListAccount(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
}

func TestGroupsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/456" {
			t.Errorf("expected /api/v1/groups/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Group{
			ID:           456,
			Name:         "Test Group",
			Description:  "Test description",
			MembersCount: 10,
			JoinLevel:    "invitation_only",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	group, err := service.Get(context.Background(), 456, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.ID != 456 {
		t.Errorf("expected ID 456, got %d", group.ID)
	}

	if group.Name != "Test Group" {
		t.Errorf("expected 'Test Group', got %s", group.Name)
	}
}

func TestGroupsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["name"] != "New Group" {
			t.Errorf("expected name 'New Group', got %v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Group{
			ID:   789,
			Name: "New Group",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	params := &CreateGroupParams{
		Name:      "New Group",
		JoinLevel: "invitation_only",
	}

	group, err := service.Create(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.ID != 789 {
		t.Errorf("expected ID 789, got %d", group.ID)
	}
}

func TestGroupsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["name"] != "Updated Group" {
			t.Errorf("expected name 'Updated Group', got %v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Group{
			ID:   456,
			Name: "Updated Group",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	name := "Updated Group"
	params := &UpdateGroupParams{
		Name: &name,
	}

	group, err := service.Update(context.Background(), 456, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.Name != "Updated Group" {
		t.Errorf("expected 'Updated Group', got %s", group.Name)
	}
}

func TestGroupsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/groups/456" {
			t.Errorf("expected /api/v1/groups/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Group{
			ID:   456,
			Name: "Deleted Group",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	group, err := service.Delete(context.Background(), 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.ID != 456 {
		t.Errorf("expected ID 456, got %d", group.ID)
	}
}

func TestGroupsService_ListMembers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/456/users" {
			t.Errorf("expected /api/v1/groups/456/users, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]User{
			{ID: 100, Name: "User 1"},
			{ID: 101, Name: "User 2"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	users, err := service.ListMembers(context.Background(), 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestGroupsService_AddMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["user_id"].(float64) != 100 {
			t.Errorf("expected user_id 100, got %v", body["user_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(GroupMembership{
			ID:      999,
			GroupID: 456,
			UserID:  100,
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	membership, err := service.AddMember(context.Background(), 456, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if membership.UserID != 100 {
		t.Errorf("expected user ID 100, got %d", membership.UserID)
	}
}

func TestGroupsService_ListCategoriesCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/group_categories" {
			t.Errorf("expected /api/v1/courses/123/group_categories, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]GroupCategory{
			{ID: 1, Name: "Category 1"},
			{ID: 2, Name: "Category 2"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	categories, err := service.ListCategoriesCourse(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(categories))
	}
}

func TestGroupsService_GetCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/group_categories/456" {
			t.Errorf("expected /api/v1/group_categories/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GroupCategory{
			ID:         456,
			Name:       "Test Category",
			SelfSignup: "enabled",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	category, err := service.GetCategory(context.Background(), 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if category.ID != 456 {
		t.Errorf("expected ID 456, got %d", category.ID)
	}

	if category.Name != "Test Category" {
		t.Errorf("expected 'Test Category', got %s", category.Name)
	}
}

func TestGroupsService_CreateCategoryCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["name"] != "New Category" {
			t.Errorf("expected name 'New Category', got %v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(GroupCategory{
			ID:   789,
			Name: "New Category",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	params := &CreateCategoryParams{
		Name:       "New Category",
		SelfSignup: "enabled",
	}

	category, err := service.CreateCategoryCourse(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if category.ID != 789 {
		t.Errorf("expected ID 789, got %d", category.ID)
	}
}

func TestGroupsService_UpdateCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["name"] != "Updated Category" {
			t.Errorf("expected name 'Updated Category', got %v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GroupCategory{
			ID:   456,
			Name: "Updated Category",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	name := "Updated Category"
	params := &UpdateCategoryParams{
		Name: &name,
	}

	category, err := service.UpdateCategory(context.Background(), 456, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if category.Name != "Updated Category" {
		t.Errorf("expected 'Updated Category', got %s", category.Name)
	}
}

func TestGroupsService_DeleteCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/group_categories/456" {
			t.Errorf("expected /api/v1/group_categories/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GroupCategory{
			ID:   456,
			Name: "Deleted Category",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	category, err := service.DeleteCategory(context.Background(), 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if category.ID != 456 {
		t.Errorf("expected ID 456, got %d", category.ID)
	}
}

func TestGroupsService_ListUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Group{{ID: 10, Name: "Group A"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGroupsService(client)
	groups, err := svc.ListUser(context.Background(), 42, nil)
	if err != nil {
		t.Fatalf("ListUser: %v", err)
	}
	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
}

func TestGroupsService_ListUser_Self(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/self/groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Group{{ID: 11}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGroupsService(client)
	opts := &ListGroupsOptions{Include: []string{"context_info"}, Page: 1, PerPage: 5}
	groups, err := svc.ListUser(context.Background(), 0, opts)
	if err != nil {
		t.Fatalf("ListUser self: %v", err)
	}
	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
}

func TestGroupsService_RemoveMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/memberships/99" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGroupsService(client)
	if err := svc.RemoveMember(context.Background(), 5, 99); err != nil {
		t.Fatalf("RemoveMember: %v", err)
	}
}

func TestGroupsService_ListCategoriesAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/3/group_categories" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]GroupCategory{{ID: 20, Name: "Cat A"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGroupsService(client)
	cats, err := svc.ListCategoriesAccount(context.Background(), 3, nil)
	if err != nil {
		t.Fatalf("ListCategoriesAccount: %v", err)
	}
	if len(cats) != 1 {
		t.Errorf("expected 1 category, got %d", len(cats))
	}
}

func TestGroupsService_ListCategoriesAccount_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("per_page") != "10" {
			t.Errorf("expected per_page=10, got %q", q.Get("per_page"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]GroupCategory{{ID: 21}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGroupsService(client)
	cats, err := svc.ListCategoriesAccount(context.Background(), 3, &ListCategoriesOptions{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("ListCategoriesAccount with opts: %v", err)
	}
	if len(cats) != 1 {
		t.Errorf("expected 1 category, got %d", len(cats))
	}
}

func TestGroupsService_CreateCategoryAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/3/group_categories" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GroupCategory{ID: 30, Name: "New Category"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGroupsService(client)
	params := &CreateCategoryParams{
		Name:               "New Category",
		SelfSignup:         "enabled",
		AutoLeader:         "first",
		GroupLimit:         5,
		SISGroupCategoryID: "sis-123",
	}
	cat, err := svc.CreateCategoryAccount(context.Background(), 3, params)
	if err != nil {
		t.Fatalf("CreateCategoryAccount: %v", err)
	}
	if cat.ID != 30 {
		t.Errorf("expected ID 30, got %d", cat.ID)
	}
}

func TestGroupsService_ListGroupsInCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/group_categories/5/groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Group{{ID: 100, Name: "Group X"}, {ID: 101, Name: "Group Y"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGroupsService(client)
	groups, err := svc.ListGroupsInCategory(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListGroupsInCategory: %v", err)
	}
	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}
}

func TestNewGroupsService(t *testing.T) {
	client, err := NewClient(ClientConfig{
		BaseURL: "https://canvas.example.com",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	if service == nil {
		t.Fatal("expected non-nil service")
		return
	}
	if service.client != client {
		t.Error("expected client to be set")
	}
}

func TestGroupsService_CreateStandalone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups" {
			t.Errorf("expected /api/v1/groups, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Group{ID: 99, Name: "Standalone Group"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	group, err := svc.CreateStandalone(context.Background(), &CreateGroupParams{Name: "Standalone Group"})
	if err != nil {
		t.Fatalf("CreateStandalone: %v", err)
	}
	if group.ID != 99 {
		t.Errorf("expected ID 99, got %d", group.ID)
	}
}

func TestGroupsService_ListMemberships(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/memberships" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]GroupMembership{
			{ID: 1, GroupID: 5, UserID: 10, WorkflowState: "accepted"},
			{ID: 2, GroupID: 5, UserID: 11, WorkflowState: "accepted"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	memberships, err := svc.ListMemberships(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListMemberships: %v", err)
	}
	if len(memberships) != 2 {
		t.Errorf("expected 2 memberships, got %d", len(memberships))
	}
}

func TestGroupsService_GetMembership(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/memberships/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GroupMembership{ID: 42, GroupID: 5, UserID: 10, WorkflowState: "accepted"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	m, err := svc.GetMembership(context.Background(), 5, 42)
	if err != nil {
		t.Fatalf("GetMembership: %v", err)
	}
	if m.ID != 42 {
		t.Errorf("expected ID 42, got %d", m.ID)
	}
}

func TestGroupsService_UpdateMembership(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/5/memberships/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GroupMembership{ID: 42, GroupID: 5, UserID: 10, Moderator: true})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	mod := true
	m, err := svc.UpdateMembership(context.Background(), 5, 42, &UpdateMembershipParams{Moderator: &mod})
	if err != nil {
		t.Fatalf("UpdateMembership: %v", err)
	}
	if !m.Moderator {
		t.Error("expected moderator=true")
	}
}

func TestGroupsService_RemoveUserBySelf(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/5/users" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	if err := svc.RemoveUserBySelf(context.Background(), 5); err != nil {
		t.Fatalf("RemoveUserBySelf: %v", err)
	}
}

func TestGroupsService_GetUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/users/100" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{ID: 100, Name: "Test User"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	u, err := svc.GetUser(context.Background(), 5, 100)
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if u.ID != 100 {
		t.Errorf("expected ID 100, got %d", u.ID)
	}
}

func TestGroupsService_UpdateUserMembership(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/5/users/100" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GroupMembership{GroupID: 5, UserID: 100, WorkflowState: "accepted"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	ws := "accepted"
	m, err := svc.UpdateUserMembership(context.Background(), 5, 100, &UpdateMembershipParams{WorkflowState: &ws})
	if err != nil {
		t.Fatalf("UpdateUserMembership: %v", err)
	}
	if m.UserID != 100 {
		t.Errorf("expected UserID 100, got %d", m.UserID)
	}
}

func TestGroupsService_RemoveUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/5/users/100" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	if err := svc.RemoveUser(context.Background(), 5, 100); err != nil {
		t.Fatalf("RemoveUser: %v", err)
	}
}

func TestGroupsService_GetActivityStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/activity_stream" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"type": "Message"}, {"type": "Submission"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	items, err := svc.GetActivityStream(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetActivityStream: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestGroupsService_GetActivityStreamSummary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/activity_stream/summary" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"type": "Message", "count": 3}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	items, err := svc.GetActivityStreamSummary(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetActivityStreamSummary: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}

func TestGroupsService_GetPermissions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/permissions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"create_discussion_topic": true})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	perms, err := svc.GetPermissions(context.Background(), 5, nil)
	if err != nil {
		t.Fatalf("GetPermissions: %v", err)
	}
	if perms["create_discussion_topic"] != true {
		t.Error("expected create_discussion_topic=true")
	}
}

func TestGroupsService_Invite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/5/invite" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]GroupMembership{{ID: 1, GroupID: 5, WorkflowState: "invited"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	memberships, err := svc.Invite(context.Background(), 5, []string{"user@example.com"})
	if err != nil {
		t.Fatalf("Invite: %v", err)
	}
	if len(memberships) != 1 {
		t.Errorf("expected 1 membership, got %d", len(memberships))
	}
}

func TestGroupsService_ListTabs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/tabs" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]GroupTab{{ID: "home", Label: "Home"}, {ID: "announcements", Label: "Announcements"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	tabs, err := svc.ListTabs(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListTabs: %v", err)
	}
	if len(tabs) != 2 {
		t.Errorf("expected 2 tabs, got %d", len(tabs))
	}
}

func TestGroupsService_ListCollaborations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/collaborations" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"id": 1, "title": "Doc 1"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	items, err := svc.ListCollaborations(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListCollaborations: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}

func TestGroupsService_ListConferences(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/conferences" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"id": 1, "title": "Conference 1"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	items, err := svc.ListConferences(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListConferences: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}

func TestGroupsService_ListExternalFeeds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/external_feeds" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]GroupExternalFeed{{ID: 1, URL: "https://example.com/feed"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	feeds, err := svc.ListExternalFeeds(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListExternalFeeds: %v", err)
	}
	if len(feeds) != 1 {
		t.Errorf("expected 1 feed, got %d", len(feeds))
	}
}

func TestGroupsService_CreateExternalFeed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/5/external_feeds" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(GroupExternalFeed{ID: 10, URL: "https://example.com/feed"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	feed, err := svc.CreateExternalFeed(context.Background(), 5, &CreateExternalFeedParams{
		URL:       "https://example.com/feed",
		Verbosity: "full",
	})
	if err != nil {
		t.Fatalf("CreateExternalFeed: %v", err)
	}
	if feed.ID != 10 {
		t.Errorf("expected ID 10, got %d", feed.ID)
	}
}

func TestGroupsService_DeleteExternalFeed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/5/external_feeds/10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GroupExternalFeed{ID: 10, URL: "https://example.com/feed"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	feed, err := svc.DeleteExternalFeed(context.Background(), 5, 10)
	if err != nil {
		t.Fatalf("DeleteExternalFeed: %v", err)
	}
	if feed.ID != 10 {
		t.Errorf("expected ID 10, got %d", feed.ID)
	}
}

func TestGroupsService_ListExternalTools(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/external_tools" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"id": 1, "name": "Tool 1"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	tools, err := svc.ListExternalTools(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListExternalTools: %v", err)
	}
	if len(tools) != 1 {
		t.Errorf("expected 1 tool, got %d", len(tools))
	}
}

func TestGroupsService_ListContentExports(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/content_exports" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]ContentExport{{ID: 1, ExportType: "common_cartridge", WorkflowState: "exported"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	exports, err := svc.ListContentExports(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListContentExports: %v", err)
	}
	if len(exports) != 1 {
		t.Errorf("expected 1 export, got %d", len(exports))
	}
}

func TestGroupsService_CreateContentExport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/5/content_exports" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ContentExport{ID: 20, ExportType: "common_cartridge", WorkflowState: "created"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	export, err := svc.CreateContentExport(context.Background(), 5, "common_cartridge")
	if err != nil {
		t.Fatalf("CreateContentExport: %v", err)
	}
	if export.ID != 20 {
		t.Errorf("expected ID 20, got %d", export.ID)
	}
}

func TestGroupsService_GetContentExport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/content_exports/20" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ContentExport{ID: 20, ExportType: "common_cartridge", WorkflowState: "exported"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	export, err := svc.GetContentExport(context.Background(), 5, 20)
	if err != nil {
		t.Fatalf("GetContentExport: %v", err)
	}
	if export.ID != 20 {
		t.Errorf("expected ID 20, got %d", export.ID)
	}
}

func TestGroupsService_ListContentLicenses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/content_licenses" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]ContentLicenseItem{{ID: "public_domain", Name: "Public Domain"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	licenses, err := svc.ListContentLicenses(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListContentLicenses: %v", err)
	}
	if len(licenses) != 1 {
		t.Errorf("expected 1 license, got %d", len(licenses))
	}
}

func TestGroupsService_ListMediaAttachments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/media_attachments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]MediaAttachment{{ID: 1, DisplayName: "video.mp4"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	items, err := svc.ListMediaAttachments(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListMediaAttachments: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}

func TestGroupsService_ListMediaObjects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/media_objects" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]MediaObject{{MediaID: "m-abc123", Title: "My Video"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	items, err := svc.ListMediaObjects(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListMediaObjects: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}

func TestGroupsService_ListPotentialCollaborators(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/potential_collaborators" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]User{{ID: 100, Name: "User A"}, {ID: 101, Name: "User B"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	users, err := svc.ListPotentialCollaborators(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListPotentialCollaborators: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestGroupsService_PreviewHTML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/5/preview_html" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"html": "<p>Hello</p>"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	html, err := svc.PreviewHTML(context.Background(), 5, "<p>Hello</p>")
	if err != nil {
		t.Fatalf("PreviewHTML: %v", err)
	}
	if html != "<p>Hello</p>" {
		t.Errorf("expected '<p>Hello</p>', got %q", html)
	}
}

func TestGroupsService_DeleteUsageRights(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/5/usage_rights" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(UsageRightsResult{UseJustification: "public_domain", Message: "deleted"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	result, err := svc.DeleteUsageRights(context.Background(), 5, []int64{1, 2})
	if err != nil {
		t.Fatalf("DeleteUsageRights: %v", err)
	}
	if result.UseJustification != "public_domain" {
		t.Errorf("unexpected use_justification: %s", result.UseJustification)
	}
}

func TestGroupsService_SetUsageRights(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/5/usage_rights" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(UsageRightsResult{UseJustification: "own_copyright"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	result, err := svc.SetUsageRights(context.Background(), 5, &SetUsageRightsGroupParams{
		FileIDs:          []int64{1},
		UseJustification: "own_copyright",
	})
	if err != nil {
		t.Fatalf("SetUsageRights: %v", err)
	}
	if result.UseJustification != "own_copyright" {
		t.Errorf("unexpected use_justification: %s", result.UseJustification)
	}
}

func TestGroupsService_GetAssignmentOverride(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/assignments/10/override" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 77, "assignment_id": 10})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	result, err := svc.GetAssignmentOverride(context.Background(), 5, 10)
	if err != nil {
		t.Fatalf("GetAssignmentOverride: %v", err)
	}
	if result["id"].(float64) != 77 {
		t.Errorf("expected id=77, got %v", result["id"])
	}
}

func TestGroupsService_AssignUnassignedMembers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/group_categories/7/assign_unassigned_members" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "running"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	result, err := svc.AssignUnassignedMembers(context.Background(), 7, false)
	if err != nil {
		t.Fatalf("AssignUnassignedMembers: %v", err)
	}
	if result["status"] != "running" {
		t.Errorf("unexpected status: %v", result["status"])
	}
}

func TestGroupsService_ListUsersInCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/group_categories/7/users" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	users, err := svc.ListUsersInCategory(context.Background(), 7)
	if err != nil {
		t.Fatalf("ListUsersInCategory: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestGroupsService_ExportCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/group_categories/7/export" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": "csv_content"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	result, err := svc.ExportCategory(context.Background(), 7)
	if err != nil {
		t.Fatalf("ExportCategory: %v", err)
	}
	if result["data"] != "csv_content" {
		t.Errorf("unexpected data: %v", result["data"])
	}
}

func TestGroupsService_ImportCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/group_categories/7/import" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "processing"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	result, err := svc.ImportCategory(context.Background(), 7, "csv data")
	if err != nil {
		t.Fatalf("ImportCategory: %v", err)
	}
	if result["status"] != "processing" {
		t.Errorf("unexpected status: %v", result["status"])
	}
}

func TestGroupsService_CreateGroupInCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/group_categories/7/groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Group{ID: 88, Name: "Cat Group"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	group, err := svc.CreateGroupInCategory(context.Background(), 7, &CreateGroupParams{Name: "Cat Group"})
	if err != nil {
		t.Fatalf("CreateGroupInCategory: %v", err)
	}
	if group.ID != 88 {
		t.Errorf("expected ID 88, got %d", group.ID)
	}
}
