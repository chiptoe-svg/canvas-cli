package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCoursePacingService_Get(t *testing.T) {
	want := coursePaceEnvelope{CoursePace: CoursePace{ID: 1, CourseID: 10, WorkflowState: "active"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/10/course_pacing/1" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursePacingService(client)
	got, err := svc.Get(context.Background(), 10, 1)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.WorkflowState != "active" {
		t.Errorf("want 'active', got %q", got.WorkflowState)
	}
}

func TestCoursePacingService_Create(t *testing.T) {
	want := coursePaceEnvelope{CoursePace: CoursePace{ID: 2, CourseID: 10, ExcludeWeekends: true}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/courses/10/course_pacing" {
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
	svc := NewCoursePacingService(client)
	t1 := true
	got, err := svc.Create(context.Background(), 10, CoursePaceParams{ExcludeWeekends: &t1})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if !got.ExcludeWeekends {
		t.Error("want ExcludeWeekends=true")
	}
}

func TestCoursePacingService_Update(t *testing.T) {
	want := coursePaceEnvelope{CoursePace: CoursePace{ID: 2, CourseID: 10, HardEndDates: true}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/courses/10/course_pacing/2" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursePacingService(client)
	t2 := true
	got, err := svc.Update(context.Background(), 10, 2, CoursePaceParams{HardEndDates: &t2})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if !got.HardEndDates {
		t.Error("want HardEndDates=true")
	}
}

func TestCoursePacingService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/courses/10/course_pacing/2" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCoursePacingService(client)
	if err := svc.Delete(context.Background(), 10, 2); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}
