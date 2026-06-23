package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCoursesService_GetActivityStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/activity_stream" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]ActivityStreamItem{{ID: 1, Type: "DiscussionTopic", Title: "Test Discussion"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	items, err := svc.GetActivityStream(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetActivityStream: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Type != "DiscussionTopic" {
		t.Errorf("expected DiscussionTopic, got %s", items[0].Type)
	}
}

func TestCoursesService_GetActivityStream_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	_, err := svc.GetActivityStream(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCoursesService_GetActivityStreamSummary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/activity_stream/summary" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]ActivityStreamSummary{{Type: "DiscussionTopic", Count: 3, UnreadCount: 1}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	summary, err := svc.GetActivityStreamSummary(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetActivityStreamSummary: %v", err)
	}
	if len(summary) != 1 {
		t.Fatalf("expected 1 summary item, got %d", len(summary))
	}
	if summary[0].Count != 3 {
		t.Errorf("expected count=3, got %d", summary[0].Count)
	}
}

func TestCoursesService_GetStudents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/students" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]User{{ID: 1, Name: "Student A"}, {ID: 2, Name: "Student B"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	students, err := svc.GetStudents(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetStudents: %v", err)
	}
	if len(students) != 2 {
		t.Fatalf("expected 2 students, got %d", len(students))
	}
	if students[0].Name != "Student A" {
		t.Errorf("expected Student A, got %s", students[0].Name)
	}
}

func TestCoursesService_GetUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/users/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{ID: 42, Name: "Specific User"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	user, err := svc.GetUser(context.Background(), 10, 42, nil)
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if user.ID != 42 {
		t.Errorf("expected ID 42, got %d", user.ID)
	}
	if user.Name != "Specific User" {
		t.Errorf("expected 'Specific User', got %s", user.Name)
	}
}

func TestCoursesService_GetUser_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	_, err := svc.GetUser(context.Background(), 10, 42, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCoursesService_ListUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/users" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]User{{ID: 1, Name: "User 1"}, {ID: 2, Name: "User 2"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	users, err := svc.ListUsers(context.Background(), 10, nil)
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestCoursesService_ListUsers_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("search_term") != "alice" {
			t.Errorf("expected search_term=alice, got %s", q.Get("search_term"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]User{{ID: 1, Name: "Alice"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	opts := &ListCourseUsersOptions{SearchTerm: "alice", Include: []string{"enrollments"}, EnrollmentType: []string{"student"}}
	users, err := svc.ListUsers(context.Background(), 10, opts)
	if err != nil {
		t.Fatalf("ListUsers with opts: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
}

func TestCoursesService_GetUserProgress(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/users/42/progress" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CourseProgress{RequirementCount: 10, RequirementCompletedCount: 7})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	progress, err := svc.GetUserProgress(context.Background(), 10, 42)
	if err != nil {
		t.Fatalf("GetUserProgress: %v", err)
	}
	if progress.RequirementCount != 10 {
		t.Errorf("expected RequirementCount=10, got %d", progress.RequirementCount)
	}
	if progress.RequirementCompletedCount != 7 {
		t.Errorf("expected RequirementCompletedCount=7, got %d", progress.RequirementCompletedCount)
	}
}

func TestCoursesService_SearchUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/search_users" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("search_term") != "bob" {
			t.Errorf("expected search_term=bob, got %s", r.URL.Query().Get("search_term"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]User{{ID: 3, Name: "Bob"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	users, err := svc.SearchUsers(context.Background(), 10, "bob")
	if err != nil {
		t.Fatalf("SearchUsers: %v", err)
	}
	if len(users) != 1 || users[0].Name != "Bob" {
		t.Errorf("expected [Bob], got %v", users)
	}
}

func TestCoursesService_GetStudentViewStudent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/student_view_student" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{ID: 999, Name: "Test Student"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	user, err := svc.GetStudentViewStudent(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetStudentViewStudent: %v", err)
	}
	if user.ID != 999 {
		t.Errorf("expected ID 999, got %d", user.ID)
	}
}

func TestCoursesService_ResetContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/reset_content" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Course{ID: 10, Name: "Reset Course"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	course, err := svc.ResetContent(context.Background(), 10)
	if err != nil {
		t.Fatalf("ResetContent: %v", err)
	}
	if course.ID != 10 {
		t.Errorf("expected ID 10, got %d", course.ID)
	}
}

func TestCoursesService_GetContentShareUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/content_share_users" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]ContentShareUser{{ID: 1, DisplayName: "Content User"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	users, err := svc.GetContentShareUsers(context.Background(), 10, "")
	if err != nil {
		t.Fatalf("GetContentShareUsers: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if users[0].DisplayName != "Content User" {
		t.Errorf("expected 'Content User', got %s", users[0].DisplayName)
	}
}

func TestCoursesService_GetPotentialCollaborators(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/potential_collaborators" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]User{{ID: 5, Name: "Collaborator"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	users, err := svc.GetPotentialCollaborators(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetPotentialCollaborators: %v", err)
	}
	if len(users) != 1 || users[0].Name != "Collaborator" {
		t.Errorf("unexpected result: %v", users)
	}
}

func TestCoursesService_GetCSPSettings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/csp_settings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"enabled": true, "inherited": false})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	result, err := svc.GetCSPSettings(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetCSPSettings: %v", err)
	}
	if result["enabled"] != true {
		t.Errorf("expected enabled=true, got %v", result["enabled"])
	}
}

func TestCoursesService_UpdateCSPSettings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/csp_settings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"enabled": body["status"] == "enabled"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursesService(client)
	result, err := svc.UpdateCSPSettings(context.Background(), 10, map[string]interface{}{"status": "enabled"})
	if err != nil {
		t.Fatalf("UpdateCSPSettings: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}
