package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRubricAssociationsService_Update(t *testing.T) {
	want := RubricAssociation{ID: 3, RubricID: 1, AssociationID: 20, UseForGrading: true}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/courses/10/rubric_associations/3" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewRubricAssociationsService(client)
	got, err := svc.Update(context.Background(), 10, 3, RubricAssociationUpdateParams{UseForGrading: true})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if !got.UseForGrading {
		t.Error("want UseForGrading=true")
	}
}

func TestRubricAssociationsService_Delete(t *testing.T) {
	want := RubricAssociation{ID: 3}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/courses/10/rubric_associations/3" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewRubricAssociationsService(client)
	got, err := svc.Delete(context.Background(), 10, 3)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if got.ID != 3 {
		t.Errorf("want ID 3, got %d", got.ID)
	}
}

func TestRubricAssociationsService_CreateAssessment(t *testing.T) {
	want := RubricAssessmentRecord{ID: 50, RubricID: 1, RubricAssociationID: 3}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/courses/10/rubric_associations/3/rubric_assessments" {
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
	svc := NewRubricAssociationsService(client)
	got, err := svc.CreateAssessment(context.Background(), 10, 3, RubricAssessmentRecord{RubricID: 1})
	if err != nil {
		t.Fatalf("CreateAssessment: %v", err)
	}
	if got.ID != 50 {
		t.Errorf("want ID 50, got %d", got.ID)
	}
}

func TestRubricAssociationsService_UpdateAssessment(t *testing.T) {
	want := RubricAssessmentRecord{ID: 50, Score: 10.0}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/courses/10/rubric_associations/3/rubric_assessments/50" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewRubricAssociationsService(client)
	got, err := svc.UpdateAssessment(context.Background(), 10, 3, 50, RubricAssessmentRecord{Score: 10.0})
	if err != nil {
		t.Fatalf("UpdateAssessment: %v", err)
	}
	if got.Score != 10.0 {
		t.Errorf("want Score 10.0, got %v", got.Score)
	}
}

func TestRubricAssociationsService_DeleteAssessment(t *testing.T) {
	want := RubricAssessmentRecord{ID: 50}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/courses/10/rubric_associations/3/rubric_assessments/50" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewRubricAssociationsService(client)
	got, err := svc.DeleteAssessment(context.Background(), 10, 3, 50)
	if err != nil {
		t.Fatalf("DeleteAssessment: %v", err)
	}
	if got.ID != 50 {
		t.Errorf("want ID 50, got %d", got.ID)
	}
}
