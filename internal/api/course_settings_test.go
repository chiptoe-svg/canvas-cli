package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCourseSettingsService_GetSettings(t *testing.T) {
	want := CourseSettings{HideFinalGrades: true, LockAllAnnouncements: true}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/7/settings" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseSettingsService(client)
	got, err := svc.GetSettings(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetSettings: %v", err)
	}
	if !got.HideFinalGrades {
		t.Error("want HideFinalGrades=true")
	}
}

func TestCourseSettingsService_UpdateSettings(t *testing.T) {
	want := CourseSettings{HideFinalGrades: false}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/courses/7/settings" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseSettingsService(client)
	got, err := svc.UpdateSettings(context.Background(), 7, CourseSettings{HideFinalGrades: false})
	if err != nil {
		t.Fatalf("UpdateSettings: %v", err)
	}
	if got.HideFinalGrades {
		t.Error("want HideFinalGrades=false")
	}
}

func TestCourseSettingsService_GetTodo(t *testing.T) {
	todos := []CourseTodo{{Type: "submitting"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/7/todo" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(todos); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseSettingsService(client)
	got, err := svc.GetTodo(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetTodo: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1, got %d", len(got))
	}
}

func TestCourseSettingsService_ListTabs(t *testing.T) {
	tabs := []CourseTab{
		{ID: "home", Label: "Home", Type: "internal", Position: 1},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/7/tabs" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(tabs); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseSettingsService(client)
	got, err := svc.ListTabs(context.Background(), 7)
	if err != nil {
		t.Fatalf("ListTabs: %v", err)
	}
	if len(got) != 1 || got[0].ID != "home" {
		t.Errorf("unexpected tabs: %+v", got)
	}
}

func TestCourseSettingsService_UpdateTab(t *testing.T) {
	want := CourseTab{ID: "grades", Label: "Grades", Hidden: false, Position: 2}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/courses/7/tabs/grades" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseSettingsService(client)
	got, err := svc.UpdateTab(context.Background(), 7, "grades", CourseTab{Hidden: false, Position: 2})
	if err != nil {
		t.Fatalf("UpdateTab: %v", err)
	}
	if got.ID != "grades" {
		t.Errorf("want 'grades', got %q", got.ID)
	}
}

func TestCourseSettingsService_GetPermissions(t *testing.T) {
	perms := map[string]bool{"read_course_content": true, "manage_grades": false}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/7/permissions" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(perms); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseSettingsService(client)
	got, err := svc.GetPermissions(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetPermissions: %v", err)
	}
	if !got["read_course_content"] {
		t.Error("want read_course_content=true")
	}
}

func TestCourseSettingsService_GetEffectiveDueDates(t *testing.T) {
	dueDates := EffectiveDueDates{
		"10": {"5": map[string]interface{}{"due_at": "2024-12-01T23:59:00Z"}},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/7/effective_due_dates" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(dueDates); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseSettingsService(client)
	got, err := svc.GetEffectiveDueDates(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetEffectiveDueDates: %v", err)
	}
	if _, ok := got["10"]; !ok {
		t.Error("want assignment 10 in effective due dates")
	}
}

func TestCourseSettingsService_GetLatePolicy(t *testing.T) {
	want := latePolicyEnvelope{LatePolicy: LatePolicy{ID: 1, CourseID: 7, LateSubmissionDeductionEnabled: true}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/7/late_policy" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseSettingsService(client)
	got, err := svc.GetLatePolicy(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetLatePolicy: %v", err)
	}
	if !got.LateSubmissionDeductionEnabled {
		t.Error("want LateSubmissionDeductionEnabled=true")
	}
}

func TestCourseSettingsService_CreateLatePolicy(t *testing.T) {
	want := latePolicyEnvelope{LatePolicy: LatePolicy{ID: 2, CourseID: 7}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/courses/7/late_policy" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseSettingsService(client)
	got, err := svc.CreateLatePolicy(context.Background(), 7, LatePolicy{})
	if err != nil {
		t.Fatalf("CreateLatePolicy: %v", err)
	}
	if got.ID != 2 {
		t.Errorf("want ID 2, got %d", got.ID)
	}
}

func TestCourseSettingsService_UpdateLatePolicy(t *testing.T) {
	want := latePolicyEnvelope{LatePolicy: LatePolicy{ID: 1, LateSubmissionDeduction: 10.0}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/api/v1/courses/7/late_policy" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseSettingsService(client)
	got, err := svc.UpdateLatePolicy(context.Background(), 7, LatePolicy{LateSubmissionDeduction: 10.0})
	if err != nil {
		t.Fatalf("UpdateLatePolicy: %v", err)
	}
	if got.LateSubmissionDeduction != 10.0 {
		t.Errorf("want 10.0, got %v", got.LateSubmissionDeduction)
	}
}

func TestCourseSettingsService_GetRecentStudents(t *testing.T) {
	students := []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/7/recent_students" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(students); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseSettingsService(client)
	got, err := svc.GetRecentStudents(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetRecentStudents: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2, got %d", len(got))
	}
}
