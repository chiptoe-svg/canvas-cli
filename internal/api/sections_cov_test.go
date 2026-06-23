package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSectionsService_GetCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/sections/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Section{ID: 5, Name: "Section Alpha", CourseID: 10})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	section, err := svc.GetCourse(context.Background(), 10, 5, nil)
	if err != nil {
		t.Fatalf("GetCourse: %v", err)
	}
	if section.ID != 5 {
		t.Errorf("expected ID 5, got %d", section.ID)
	}
	if section.Name != "Section Alpha" {
		t.Errorf("expected name Section Alpha, got %s", section.Name)
	}
}

func TestSectionsService_GetCourse_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	_, err := svc.GetCourse(context.Background(), 10, 5, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSectionsService_ListUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/sections/5/users" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	users, err := svc.ListUsers(context.Background(), 5, nil)
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].Name != "Alice" {
		t.Errorf("expected Alice, got %s", users[0].Name)
	}
}

func TestSectionsService_CreateEnrollment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/sections/5/enrollments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": float64(99), "type": "StudentEnrollment"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	result, err := svc.CreateEnrollment(context.Background(), 5, map[string]interface{}{"enrollment": map[string]interface{}{"user_id": 42}})
	if err != nil {
		t.Fatalf("CreateEnrollment: %v", err)
	}
	if result["type"] != "StudentEnrollment" {
		t.Errorf("expected StudentEnrollment, got %v", result["type"])
	}
}

func TestSectionsService_ListSubmissions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/sections/5/assignments/10/submissions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Submission{{ID: 1, Score: 95.0}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	submissions, err := svc.ListSubmissions(context.Background(), 5, 10, nil)
	if err != nil {
		t.Fatalf("ListSubmissions: %v", err)
	}
	if len(submissions) != 1 {
		t.Fatalf("expected 1 submission, got %d", len(submissions))
	}
	if submissions[0].Score != 95.0 {
		t.Errorf("expected score 95, got %f", submissions[0].Score)
	}
}

func TestSectionsService_GetSubmission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/sections/5/assignments/10/submissions/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Submission{ID: 1, Score: 88.5, UserID: 42})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	sub, err := svc.GetSubmission(context.Background(), 5, 10, 42, nil)
	if err != nil {
		t.Fatalf("GetSubmission: %v", err)
	}
	if sub.Score != 88.5 {
		t.Errorf("expected score 88.5, got %f", sub.Score)
	}
}

func TestSectionsService_GetSubmission_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	_, err := svc.GetSubmission(context.Background(), 5, 10, 42, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSectionsService_ListSubmissionsForStudents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/sections/5/students/submissions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("student_ids[]") == "" {
			t.Error("expected student_ids[] query param")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Submission{{ID: 1}, {ID: 2}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	submissions, err := svc.ListSubmissionsForStudents(context.Background(), 5, []int64{42, 43}, nil)
	if err != nil {
		t.Fatalf("ListSubmissionsForStudents: %v", err)
	}
	if len(submissions) != 2 {
		t.Errorf("expected 2 submissions, got %d", len(submissions))
	}
}

func TestSectionsService_GradeSubmission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/sections/5/assignments/10/submissions/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		sub, ok := body["submission"].(map[string]interface{})
		if !ok {
			t.Fatal("expected submission key in body")
		}
		if sub["posted_grade"] != "95" {
			t.Errorf("expected posted_grade=95, got %v", sub["posted_grade"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Submission{ID: 1, Score: 95.0})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	params := &GradeSubmissionParams{PostedGrade: "95"}
	sub, err := svc.GradeSubmission(context.Background(), 5, 10, 42, params)
	if err != nil {
		t.Fatalf("GradeSubmission: %v", err)
	}
	if sub.Score != 95.0 {
		t.Errorf("expected score 95, got %f", sub.Score)
	}
}

func TestSectionsService_GetAssignmentOverride(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/sections/5/assignments/10/override" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AssignmentOverride{ID: 1, AssignmentID: 10, Title: "Section Override"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	override, err := svc.GetAssignmentOverride(context.Background(), 5, 10)
	if err != nil {
		t.Fatalf("GetAssignmentOverride: %v", err)
	}
	if override.Title != "Section Override" {
		t.Errorf("expected 'Section Override', got %s", override.Title)
	}
}

func TestSectionsService_GetSubmissionSummary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/sections/5/assignments/10/submission_summary" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SubmissionSummary{Graded: 10, Ungraded: 5, NotSubmitted: 2})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	summary, err := svc.GetSubmissionSummary(context.Background(), 5, 10)
	if err != nil {
		t.Fatalf("GetSubmissionSummary: %v", err)
	}
	if summary.Graded != 10 {
		t.Errorf("expected graded=10, got %d", summary.Graded)
	}
	if summary.Ungraded != 5 {
		t.Errorf("expected ungraded=5, got %d", summary.Ungraded)
	}
}

func TestSectionsService_SubmitAssignment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/sections/5/assignments/10/submissions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Submission{ID: 99, UserID: 42})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	sub, err := svc.SubmitAssignment(context.Background(), 5, 10, map[string]interface{}{"submission_type": "online_text_entry"})
	if err != nil {
		t.Fatalf("SubmitAssignment: %v", err)
	}
	if sub.ID != 99 {
		t.Errorf("expected ID 99, got %d", sub.ID)
	}
}

func TestSectionsService_MarkSubmissionAsRead(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/sections/5/assignments/10/submissions/42/read" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	err := svc.MarkSubmissionAsRead(context.Background(), 5, 10, 42)
	if err != nil {
		t.Fatalf("MarkSubmissionAsRead: %v", err)
	}
	if !called {
		t.Error("expected PUT request to be made")
	}
}

func TestSectionsService_MarkSubmissionAsUnread(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/sections/5/assignments/10/submissions/42/read" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	err := svc.MarkSubmissionAsUnread(context.Background(), 5, 10, 42)
	if err != nil {
		t.Fatalf("MarkSubmissionAsUnread: %v", err)
	}
	if !called {
		t.Error("expected DELETE request to be made")
	}
}

func TestSectionsService_BulkMarkRead(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/sections/5/submissions/bulk_mark_read" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	err := svc.BulkMarkRead(context.Background(), 5)
	if err != nil {
		t.Fatalf("BulkMarkRead: %v", err)
	}
	if !called {
		t.Error("expected PUT request to be made")
	}
}

func TestSectionsService_UpdateGrades(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/sections/5/assignments/10/submissions/update_grades" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"progress": map[string]interface{}{"workflow_state": "queued"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	result, err := svc.UpdateGrades(context.Background(), 5, 10, map[string]interface{}{"grade_data": map[string]interface{}{"42": map[string]interface{}{"posted_grade": "A"}}})
	if err != nil {
		t.Fatalf("UpdateGrades: %v", err)
	}
	if _, ok := result["progress"]; !ok {
		t.Error("expected progress key in result")
	}
}

func TestSectionsService_UpdateSubmissionsGrades(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/sections/5/submissions/update_grades" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"completed": true})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	result, err := svc.UpdateSubmissionsGrades(context.Background(), 5, map[string]interface{}{})
	if err != nil {
		t.Fatalf("UpdateSubmissionsGrades: %v", err)
	}
	if result["completed"] != true {
		t.Errorf("expected completed=true, got %v", result["completed"])
	}
}

func TestSectionsService_ClearSubmissionUnread(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/sections/5/submissions/42/clear_unread" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	err := svc.ClearSubmissionUnread(context.Background(), 5, 42)
	if err != nil {
		t.Fatalf("ClearSubmissionUnread: %v", err)
	}
	if !called {
		t.Error("expected PUT request")
	}
}

func TestSectionsService_ListPeerReviews(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/sections/5/assignments/10/peer_reviews" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"assessor_id": float64(1), "asset_id": float64(2)}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewSectionsService(client)
	reviews, err := svc.ListPeerReviews(context.Background(), 5, 10)
	if err != nil {
		t.Fatalf("ListPeerReviews: %v", err)
	}
	if len(reviews) != 1 {
		t.Fatalf("expected 1 peer review, got %d", len(reviews))
	}
	if reviews[0]["assessor_id"] != float64(1) {
		t.Errorf("expected assessor_id=1, got %v", reviews[0]["assessor_id"])
	}
}
