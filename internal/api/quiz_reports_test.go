package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuizReportsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/reports" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]QuizReport{
			{ID: 5, QuizID: 2, ReportType: "student_analysis"},
			{ID: 6, QuizID: 2, ReportType: "item_analysis"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizReportsService(client)

	reports, err := svc.List(context.Background(), 1, 2, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reports) != 2 {
		t.Errorf("expected 2 reports, got %d", len(reports))
	}
}

func TestQuizReportsService_List_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Query().Get("includes_all_versions") != "true" {
			t.Errorf("expected includes_all_versions=true")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]QuizReport{{ID: 7, QuizID: 2, ReportType: "student_analysis"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizReportsService(client)

	reports, err := svc.List(context.Background(), 1, 2, &ListQuizReportsOptions{IncludesAllVersions: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reports) != 1 {
		t.Errorf("expected 1 report, got %d", len(reports))
	}
}

func TestQuizReportsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/reports/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizReport{ID: 5, QuizID: 2, ReportType: "student_analysis"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizReportsService(client)

	report, err := svc.Get(context.Background(), 1, 2, 5, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.ID != 5 {
		t.Errorf("expected ID 5, got %d", report.ID)
	}
}

func TestQuizReportsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/reports" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		report, ok := body["quiz_report"].(map[string]interface{})
		if !ok {
			t.Errorf("expected quiz_report in body")
		}
		if report["report_type"] != "student_analysis" {
			t.Errorf("expected report_type 'student_analysis', got %v", report["report_type"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizReport{ID: 8, QuizID: 2, ReportType: "student_analysis"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizReportsService(client)

	report, err := svc.Create(context.Background(), 1, 2, &CreateQuizReportParams{
		ReportType: "student_analysis",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.ID != 8 {
		t.Errorf("expected ID 8, got %d", report.ID)
	}
}

func TestQuizReportsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/reports/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizReportsService(client)

	if err := svc.Delete(context.Background(), 1, 2, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
