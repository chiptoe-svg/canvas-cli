package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGroupsService_ListCourse_WithOptions exercises the query-string branch
// of ListCourse that was previously at 33.3%.
func TestGroupsService_ListCourse_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("per_page") != "5" {
			t.Errorf("expected per_page=5, got %q", q.Get("per_page"))
		}
		if q.Get("page") != "2" {
			t.Errorf("expected page=2, got %q", q.Get("page"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Group{{ID: 1, Name: "G1"}})
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	opts := &ListGroupsOptions{Include: []string{"users"}, Page: 2, PerPage: 5}
	groups, err := svc.ListCourse(context.Background(), 10, opts)
	if err != nil {
		t.Fatalf("ListCourse with opts: %v", err)
	}
	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
}

// TestGroupsService_ListCourse_Error exercises the error path in ListCourse.
func TestGroupsService_ListCourse_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors":[{"message":"server error"}]}`))
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	_, err := svc.ListCourse(context.Background(), 10, nil)
	if err == nil {
		t.Fatal("expected error from ListCourse, got nil")
	}
}

// TestGroupsService_ListAccount_WithOptions exercises the query-string branch
// of ListAccount (was at 33.3%).
func TestGroupsService_ListAccount_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/2/groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("per_page") != "10" {
			t.Errorf("expected per_page=10, got %q", q.Get("per_page"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Group{{ID: 20, Name: "Account Group"}})
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	opts := &ListGroupsOptions{Include: []string{"permissions"}, Page: 1, PerPage: 10}
	groups, err := svc.ListAccount(context.Background(), 2, opts)
	if err != nil {
		t.Fatalf("ListAccount with opts: %v", err)
	}
	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
}

// TestGroupsService_ListAccount_Error exercises the error path in ListAccount.
func TestGroupsService_ListAccount_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"errors":[{"message":"unauthorized"}]}`))
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	_, err := svc.ListAccount(context.Background(), 2, nil)
	if err == nil {
		t.Fatal("expected error from ListAccount, got nil")
	}
}

// TestGroupsService_Get_Error exercises the error path in Get.
func TestGroupsService_Get_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"not found"}]}`))
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	_, err := svc.Get(context.Background(), 999, nil)
	if err == nil {
		t.Fatal("expected error from Get, got nil")
	}
}

// TestGroupsService_Update_AllFields exercises the Update path with every
// optional field set, raising coverage from 65%.
func TestGroupsService_Update_AllFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		// Verify all fields arrived
		if body["name"] != "Full Update" {
			t.Errorf("name mismatch: %v", body["name"])
		}
		if body["description"] != "A description" {
			t.Errorf("description mismatch: %v", body["description"])
		}
		if body["is_public"] != true {
			t.Errorf("is_public mismatch: %v", body["is_public"])
		}
		if body["join_level"] != "invitation_only" {
			t.Errorf("join_level mismatch: %v", body["join_level"])
		}
		if body["sis_group_id"] != "sis-abc" {
			t.Errorf("sis_group_id mismatch: %v", body["sis_group_id"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Group{ID: 7, Name: "Full Update"})
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	name := "Full Update"
	desc := "A description"
	pub := true
	jl := "invitation_only"
	sisID := "sis-abc"
	var avatarID int64 = 99
	var quota int64 = 512
	params := &UpdateGroupParams{
		Name:           &name,
		Description:    &desc,
		IsPublic:       &pub,
		JoinLevel:      &jl,
		AvatarID:       &avatarID,
		StorageQuotaMb: &quota,
		SISGroupID:     &sisID,
	}
	group, err := svc.Update(context.Background(), 7, params)
	if err != nil {
		t.Fatalf("Update all fields: %v", err)
	}
	if group.Name != "Full Update" {
		t.Errorf("expected 'Full Update', got %q", group.Name)
	}
}

// TestGroupsService_Update_Error exercises the error path of Update.
func TestGroupsService_Update_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errors":[{"message":"bad request"}]}`))
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	name := "x"
	_, err := svc.Update(context.Background(), 1, &UpdateGroupParams{Name: &name})
	if err == nil {
		t.Fatal("expected error from Update, got nil")
	}
}

// TestGroupsService_Delete_Error exercises the error path of Delete.
func TestGroupsService_Delete_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"not found"}]}`))
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	_, err := svc.Delete(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error from Delete, got nil")
	}
}

// TestGroupsService_ListMembers_Error exercises the error path of ListMembers.
func TestGroupsService_ListMembers_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"errors":[{"message":"forbidden"}]}`))
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	_, err := svc.ListMembers(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error from ListMembers, got nil")
	}
}

// TestGroupsService_AddMember_Error exercises the error path of AddMember.
func TestGroupsService_AddMember_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"errors":[{"message":"already a member"}]}`))
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	_, err := svc.AddMember(context.Background(), 1, 99)
	if err == nil {
		t.Fatal("expected error from AddMember, got nil")
	}
}

// TestGroupsService_ListCategoriesCourse_WithOptions exercises the pagination
// query branch of ListCategoriesCourse (was at 38.5%).
func TestGroupsService_ListCategoriesCourse_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/7/group_categories" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("page") != "1" {
			t.Errorf("expected page=1, got %q", q.Get("page"))
		}
		if q.Get("per_page") != "20" {
			t.Errorf("expected per_page=20, got %q", q.Get("per_page"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]GroupCategory{{ID: 10, Name: "Cat A"}, {ID: 11, Name: "Cat B"}})
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	cats, err := svc.ListCategoriesCourse(context.Background(), 7, &ListCategoriesOptions{Page: 1, PerPage: 20})
	if err != nil {
		t.Fatalf("ListCategoriesCourse with opts: %v", err)
	}
	if len(cats) != 2 {
		t.Errorf("expected 2 categories, got %d", len(cats))
	}
}

// TestGroupsService_ListCategoriesCourse_Error exercises the error branch.
func TestGroupsService_ListCategoriesCourse_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"not found"}]}`))
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	_, err := svc.ListCategoriesCourse(context.Background(), 99, nil)
	if err == nil {
		t.Fatal("expected error from ListCategoriesCourse, got nil")
	}
}

// TestGroupsService_GetCategory_Error exercises the error branch of GetCategory.
func TestGroupsService_GetCategory_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"not found"}]}`))
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	_, err := svc.GetCategory(context.Background(), 99)
	if err == nil {
		t.Fatal("expected error from GetCategory, got nil")
	}
}

// TestGroupsService_CreateCategoryCourse_AllFields exercises all optional
// parameters of CreateCategoryCourse (was at 70%).
func TestGroupsService_CreateCategoryCourse_AllFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["self_signup"] != "enabled" {
			t.Errorf("self_signup mismatch: %v", body["self_signup"])
		}
		if body["auto_leader"] != "random" {
			t.Errorf("auto_leader mismatch: %v", body["auto_leader"])
		}
		if body["create_group_count"].(float64) != 3 {
			t.Errorf("create_group_count mismatch: %v", body["create_group_count"])
		}
		if body["split_group_count"].(float64) != 4 {
			t.Errorf("split_group_count mismatch: %v", body["split_group_count"])
		}
		if body["sis_group_category_id"] != "sis-cat-1" {
			t.Errorf("sis_group_category_id mismatch: %v", body["sis_group_category_id"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GroupCategory{ID: 50, Name: "Full Cat"})
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	params := &CreateCategoryParams{
		Name:               "Full Cat",
		SelfSignup:         "enabled",
		AutoLeader:         "random",
		GroupLimit:         5,
		CreateGroupCount:   3,
		SplitGroupCount:    4,
		SISGroupCategoryID: "sis-cat-1",
	}
	cat, err := svc.CreateCategoryCourse(context.Background(), 5, params)
	if err != nil {
		t.Fatalf("CreateCategoryCourse all fields: %v", err)
	}
	if cat.ID != 50 {
		t.Errorf("expected ID 50, got %d", cat.ID)
	}
}

// TestGroupsService_CreateCategoryCourse_Error exercises the error path.
func TestGroupsService_CreateCategoryCourse_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errors":[{"message":"bad request"}]}`))
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	_, err := svc.CreateCategoryCourse(context.Background(), 1, &CreateCategoryParams{Name: "x"})
	if err == nil {
		t.Fatal("expected error from CreateCategoryCourse, got nil")
	}
}

// TestGroupsService_UpdateCategory_AllFields exercises all optional parameters
// of UpdateCategory (was at 68.8%).
func TestGroupsService_UpdateCategory_AllFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["self_signup"] != "restricted" {
			t.Errorf("self_signup mismatch: %v", body["self_signup"])
		}
		if body["auto_leader"] != "first" {
			t.Errorf("auto_leader mismatch: %v", body["auto_leader"])
		}
		if body["group_limit"].(float64) != 8 {
			t.Errorf("group_limit mismatch: %v", body["group_limit"])
		}
		if body["sis_group_category_id"] != "sis-upd" {
			t.Errorf("sis_group_category_id mismatch: %v", body["sis_group_category_id"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GroupCategory{ID: 60, Name: "Updated"})
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	name := "Updated"
	ss := "restricted"
	al := "first"
	gl := 8
	sis := "sis-upd"
	params := &UpdateCategoryParams{
		Name:               &name,
		SelfSignup:         &ss,
		AutoLeader:         &al,
		GroupLimit:         &gl,
		SISGroupCategoryID: &sis,
	}
	cat, err := svc.UpdateCategory(context.Background(), 60, params)
	if err != nil {
		t.Fatalf("UpdateCategory all fields: %v", err)
	}
	if cat.ID != 60 {
		t.Errorf("expected ID 60, got %d", cat.ID)
	}
}

// TestGroupsService_UpdateCategory_Error exercises the error path of UpdateCategory.
func TestGroupsService_UpdateCategory_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"not found"}]}`))
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	name := "x"
	_, err := svc.UpdateCategory(context.Background(), 999, &UpdateCategoryParams{Name: &name})
	if err == nil {
		t.Fatal("expected error from UpdateCategory, got nil")
	}
}

// TestGroupsService_DeleteCategory_Error exercises the error path of DeleteCategory.
func TestGroupsService_DeleteCategory_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"errors":[{"message":"forbidden"}]}`))
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	_, err := svc.DeleteCategory(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error from DeleteCategory, got nil")
	}
}

// TestGroupsService_ListGroupsInCategory_Error exercises the error path of
// ListGroupsInCategory.
func TestGroupsService_ListGroupsInCategory_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"not found"}]}`))
	}))
	defer server.Close()

	svc := NewGroupsService(newTestClient(t, server.URL))
	_, err := svc.ListGroupsInCategory(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error from ListGroupsInCategory, got nil")
	}
}
