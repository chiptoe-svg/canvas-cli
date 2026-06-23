package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsersService_GetAvatars(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/avatars" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]UserAvatar{
			{Type: "attachment", URL: "https://example.com/avatar1.jpg", DisplayName: "Avatar 1"},
			{Type: "gravatar", URL: "https://example.com/gravatar.jpg", DisplayName: "Gravatar"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	avatars, err := svc.GetAvatars(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetAvatars: %v", err)
	}
	if len(avatars) != 2 {
		t.Fatalf("expected 2 avatars, got %d", len(avatars))
	}
	if avatars[0].Type != "attachment" {
		t.Errorf("expected type=attachment, got %s", avatars[0].Type)
	}
	if avatars[1].DisplayName != "Gravatar" {
		t.Errorf("expected DisplayName=Gravatar, got %s", avatars[1].DisplayName)
	}
}

func TestUsersService_GetAvatars_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	_, err := svc.GetAvatars(context.Background(), 42)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersService_GetHistory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/history" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{
			map[string]interface{}{"item_type": "Course", "item_id": float64(10)},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	history, err := svc.GetHistory(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetHistory: %v", err)
	}
	if len(history) != 1 {
		t.Fatalf("expected 1 history item, got %d", len(history))
	}
}

func TestUsersService_GetMissingSubmissions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/missing_submissions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]MissingSubmission{{ID: 1, Name: "Missing Assignment"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	missing, err := svc.GetMissingSubmissions(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetMissingSubmissions: %v", err)
	}
	if len(missing) != 1 {
		t.Fatalf("expected 1 missing submission, got %d", len(missing))
	}
	if missing[0].Name != "Missing Assignment" {
		t.Errorf("expected 'Missing Assignment', got %s", missing[0].Name)
	}
}

func TestUsersService_GetPageViews(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/page_views" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]PageView{{ID: "abc123", URL: "/courses/10", HTTPMethod: "GET"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	views, err := svc.GetPageViews(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetPageViews: %v", err)
	}
	if len(views) != 1 {
		t.Fatalf("expected 1 page view, got %d", len(views))
	}
	if views[0].HTTPMethod != "GET" {
		t.Errorf("expected HTTPMethod=GET, got %s", views[0].HTTPMethod)
	}
}

func TestUsersService_GetTabs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/tabs" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Tab{
			{ID: "courses", Label: "Courses", Type: "internal"},
			{ID: "groups", Label: "Groups", Type: "internal"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	tabs, err := svc.GetTabs(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetTabs: %v", err)
	}
	if len(tabs) != 2 {
		t.Fatalf("expected 2 tabs, got %d", len(tabs))
	}
	if tabs[0].Label != "Courses" {
		t.Errorf("expected Courses, got %s", tabs[0].Label)
	}
}

func TestUsersService_DeleteLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/logins/7" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(UserLogin{ID: 7, UserID: 42, UniqueID: "testuser"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	login, err := svc.DeleteLogin(context.Background(), 42, 7)
	if err != nil {
		t.Fatalf("DeleteLogin: %v", err)
	}
	if login.ID != 7 {
		t.Errorf("expected ID 7, got %d", login.ID)
	}
	if login.UniqueID != "testuser" {
		t.Errorf("expected UniqueID=testuser, got %s", login.UniqueID)
	}
}

func TestUsersService_DeleteLogin_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	_, err := svc.DeleteLogin(context.Background(), 42, 7)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersService_ListUserCourses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/courses" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Course{{ID: 1, Name: "Course A"}, {ID: 2, Name: "Course B"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	courses, err := svc.ListUserCourses(context.Background(), 42)
	if err != nil {
		t.Fatalf("ListUserCourses: %v", err)
	}
	if len(courses) != 2 {
		t.Fatalf("expected 2 courses, got %d", len(courses))
	}
	if courses[0].Name != "Course A" {
		t.Errorf("expected 'Course A', got %s", courses[0].Name)
	}
}

func TestUsersService_GetColor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/colors/course_10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(UserColor{HexCode: "#FF5733"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	color, err := svc.GetColor(context.Background(), 42, "course_10")
	if err != nil {
		t.Fatalf("GetColor: %v", err)
	}
	if color.HexCode != "#FF5733" {
		t.Errorf("expected #FF5733, got %s", color.HexCode)
	}
}

func TestUsersService_SetColor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/colors/course_10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["hexcode"] != "#00AABB" {
			t.Errorf("expected hexcode=#00AABB, got %s", body["hexcode"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(UserColor{HexCode: "#00AABB"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	color, err := svc.SetColor(context.Background(), 42, "course_10", "#00AABB")
	if err != nil {
		t.Fatalf("SetColor: %v", err)
	}
	if color.HexCode != "#00AABB" {
		t.Errorf("expected #00AABB, got %s", color.HexCode)
	}
}

func TestUsersService_GetDashboardPositions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/dashboard_positions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DashboardPositions{DashboardPositions: map[string]int{"course_10": 1, "course_11": 2}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	positions, err := svc.GetDashboardPositions(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetDashboardPositions: %v", err)
	}
	if positions.DashboardPositions["course_10"] != 1 {
		t.Errorf("expected course_10=1, got %d", positions.DashboardPositions["course_10"])
	}
}

func TestUsersService_SetDashboardPositions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/dashboard_positions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DashboardPositions{DashboardPositions: map[string]int{"course_10": 3}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	positions, err := svc.SetDashboardPositions(context.Background(), 42, map[string]int{"course_10": 3})
	if err != nil {
		t.Fatalf("SetDashboardPositions: %v", err)
	}
	if positions.DashboardPositions["course_10"] != 3 {
		t.Errorf("expected course_10=3, got %d", positions.DashboardPositions["course_10"])
	}
}

func TestUsersService_MergeInto(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/merge_into/99" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{ID: 99, Name: "Destination User"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	user, err := svc.MergeInto(context.Background(), 42, 99)
	if err != nil {
		t.Fatalf("MergeInto: %v", err)
	}
	if user.ID != 99 {
		t.Errorf("expected ID 99, got %d", user.ID)
	}
}

func TestUsersService_MergeIntoAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/merge_into/accounts/1/users/99" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{ID: 99, Name: "Merged User"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	user, err := svc.MergeIntoAccount(context.Background(), 42, 1, 99)
	if err != nil {
		t.Fatalf("MergeIntoAccount: %v", err)
	}
	if user.ID != 99 {
		t.Errorf("expected ID 99, got %d", user.ID)
	}
}

func TestUsersService_Split(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/42/split" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]User{{ID: 42, Name: "User A"}, {ID: 43, Name: "User B"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	users, err := svc.Split(context.Background(), 42)
	if err != nil {
		t.Fatalf("Split: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestUsersService_DeleteActivityStream(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/activity_stream" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	err := svc.DeleteActivityStream(context.Background())
	if err != nil {
		t.Fatalf("DeleteActivityStream: %v", err)
	}
	if !called {
		t.Error("expected DELETE request")
	}
}

func TestUsersService_DeleteActivityStreamItem(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/activity_stream/77" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	err := svc.DeleteActivityStreamItem(context.Background(), 77)
	if err != nil {
		t.Fatalf("DeleteActivityStreamItem: %v", err)
	}
	if !called {
		t.Error("expected DELETE request")
	}
}

func TestUsersService_GetActivityStreamSummary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/self/activity_stream/summary" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]ActivityStreamSummary{{Type: "Message", Count: 5, UnreadCount: 2}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	summaries, err := svc.GetActivityStreamSummary(context.Background())
	if err != nil {
		t.Fatalf("GetActivityStreamSummary: %v", err)
	}
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
	if summaries[0].Count != 5 {
		t.Errorf("expected count=5, got %d", summaries[0].Count)
	}
	if summaries[0].UnreadCount != 2 {
		t.Errorf("expected unread_count=2, got %d", summaries[0].UnreadCount)
	}
}

func TestUsersService_GetTodoItemCount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/self/todo_item_count" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TodoItemCount{NeedsGradingCount: 7, AssignmentsNeedingSubmitting: 3})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	count, err := svc.GetTodoItemCount(context.Background())
	if err != nil {
		t.Fatalf("GetTodoItemCount: %v", err)
	}
	if count.NeedsGradingCount != 7 {
		t.Errorf("expected NeedsGradingCount=7, got %d", count.NeedsGradingCount)
	}
	if count.AssignmentsNeedingSubmitting != 3 {
		t.Errorf("expected AssignmentsNeedingSubmitting=3, got %d", count.AssignmentsNeedingSubmitting)
	}
}

func TestUsersService_GetUpcomingEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/self/upcoming_events" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{
			map[string]interface{}{"type": "quiz", "title": "Quiz 1"},
			map[string]interface{}{"type": "assignment", "title": "HW 1"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewUsersService(client)
	events, err := svc.GetUpcomingEvents(context.Background())
	if err != nil {
		t.Fatalf("GetUpcomingEvents: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}
