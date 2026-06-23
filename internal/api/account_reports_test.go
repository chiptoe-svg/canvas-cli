package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountReportsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/reports" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		reports := []AccountReport{
			{ID: "course_storage_csv", Title: "Course Storage"},
			{ID: "student_assignment_outcome_map_csv", Title: "Student Competency"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reports)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountReportsService(client)

	reports, err := svc.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(reports) != 2 {
		t.Errorf("expected 2 reports, got %d", len(reports))
	}
	if reports[0].ID != "course_storage_csv" {
		t.Errorf("expected report ID 'course_storage_csv', got %s", reports[0].ID)
	}
}

func TestAccountReportsService_ListRuns(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/reports/course_storage_csv" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		runs := []AccountReportRun{
			{ID: 10, Report: "course_storage_csv", Status: "complete", Progress: 100},
			{ID: 11, Report: "course_storage_csv", Status: "running", Progress: 50},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(runs)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountReportsService(client)

	runs, err := svc.ListRuns(context.Background(), 1, "course_storage_csv")
	if err != nil {
		t.Fatalf("ListRuns: %v", err)
	}
	if len(runs) != 2 {
		t.Errorf("expected 2 runs, got %d", len(runs))
	}
	if runs[0].Status != "complete" {
		t.Errorf("expected status 'complete', got %s", runs[0].Status)
	}
}

func TestAccountReportsService_Start(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/reports/course_storage_csv" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		run := AccountReportRun{
			ID:     20,
			Report: "course_storage_csv",
			Status: "created",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(run)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountReportsService(client)

	run, err := svc.Start(context.Background(), 1, "course_storage_csv", nil)
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	if run.ID != 20 {
		t.Errorf("expected ID 20, got %d", run.ID)
	}
	if run.Status != "created" {
		t.Errorf("expected status 'created', got %s", run.Status)
	}
}

func TestAccountReportsService_GetRun(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/reports/course_storage_csv/10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		run := AccountReportRun{
			ID:       10,
			Report:   "course_storage_csv",
			Status:   "complete",
			Progress: 100,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(run)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountReportsService(client)

	run, err := svc.GetRun(context.Background(), 1, "course_storage_csv", 10)
	if err != nil {
		t.Fatalf("GetRun: %v", err)
	}
	if run.ID != 10 {
		t.Errorf("expected ID 10, got %d", run.ID)
	}
}

func TestAccountReportsService_DeleteRun(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/reports/course_storage_csv/10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountReportsService(client)

	err := svc.DeleteRun(context.Background(), 1, "course_storage_csv", 10)
	if err != nil {
		t.Fatalf("DeleteRun: %v", err)
	}
}

func TestAccountReportsService_AbortRun(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/reports/course_storage_csv/10/abort" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		run := AccountReportRun{
			ID:     10,
			Report: "course_storage_csv",
			Status: "aborted",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(run)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountReportsService(client)

	run, err := svc.AbortRun(context.Background(), 1, "course_storage_csv", 10)
	if err != nil {
		t.Fatalf("AbortRun: %v", err)
	}
	if run.Status != "aborted" {
		t.Errorf("expected status 'aborted', got %s", run.Status)
	}
}
