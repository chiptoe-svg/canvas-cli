package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSectionsService_Get_WithInclude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		includes := q["include[]"]
		if len(includes) != 2 {
			t.Errorf("expected 2 includes, got %d", len(includes))
		}
		json.NewEncoder(w).Encode(Section{ID: 10, Name: "Section With Includes", CourseID: 5, TotalStudents: 30})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewSectionsService(client)
	section, err := svc.Get(context.Background(), 10, []string{"students", "total_students"})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if section.TotalStudents != 30 {
		t.Errorf("expected TotalStudents 30, got %d", section.TotalStudents)
	}
}

func TestSectionsService_Update_AllParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		sectionData, ok := body["course_section"].(map[string]interface{})
		if !ok {
			t.Fatal("expected course_section in body")
		}
		if sectionData["name"] != "All Params Section" {
			t.Errorf("expected name 'All Params Section', got %v", sectionData["name"])
		}
		if sectionData["sis_section_id"] != "SIS999" {
			t.Errorf("expected sis_section_id 'SIS999', got %v", sectionData["sis_section_id"])
		}
		if sectionData["integration_id"] != "INT001" {
			t.Errorf("expected integration_id 'INT001', got %v", sectionData["integration_id"])
		}
		if sectionData["start_at"] != "2024-01-01T00:00:00Z" {
			t.Errorf("expected start_at, got %v", sectionData["start_at"])
		}
		if sectionData["end_at"] != "2024-06-30T00:00:00Z" {
			t.Errorf("expected end_at, got %v", sectionData["end_at"])
		}
		if sectionData["restrict_enrollments_to_section_dates"] != true {
			t.Errorf("expected restrict_enrollments_to_section_dates=true")
		}
		if body["override_sis_stickiness"] != true {
			t.Errorf("expected override_sis_stickiness=true, got %v", body["override_sis_stickiness"])
		}
		json.NewEncoder(w).Encode(Section{ID: 10, Name: "All Params Section"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewSectionsService(client)
	name := "All Params Section"
	sis := "SIS999"
	integID := "INT001"
	startAt := "2024-01-01T00:00:00Z"
	endAt := "2024-06-30T00:00:00Z"
	restrict := true
	params := &UpdateSectionParams{
		Name:                              &name,
		SISSectionID:                      &sis,
		IntegrationID:                     &integID,
		StartAt:                           &startAt,
		EndAt:                             &endAt,
		RestrictEnrollmentsToSectionDates: &restrict,
		OverrideSISStickiness:             true,
	}
	section, err := svc.Update(context.Background(), 10, params)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if section.Name != "All Params Section" {
		t.Errorf("expected 'All Params Section', got %s", section.Name)
	}
}

func TestSectionsService_Create_WithIntegrationAndEndAt(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		sectionData, ok := body["course_section"].(map[string]interface{})
		if !ok {
			t.Fatal("expected course_section in body")
		}
		if sectionData["integration_id"] != "INTEG42" {
			t.Errorf("expected integration_id 'INTEG42', got %v", sectionData["integration_id"])
		}
		if sectionData["end_at"] != "2024-12-31T00:00:00Z" {
			t.Errorf("expected end_at, got %v", sectionData["end_at"])
		}
		if sectionData["enable_sis_reactivation"] != true {
			t.Errorf("expected enable_sis_reactivation=true, got %v", sectionData["enable_sis_reactivation"])
		}
		json.NewEncoder(w).Encode(Section{ID: 55, Name: "Integration Section"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewSectionsService(client)
	params := &CreateSectionParams{
		Name:                  "Integration Section",
		IntegrationID:         "INTEG42",
		EndAt:                 "2024-12-31T00:00:00Z",
		EnableSISReactivation: true,
	}
	section, err := svc.Create(context.Background(), 5, params)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if section.ID != 55 {
		t.Errorf("expected ID 55, got %d", section.ID)
	}
}

func TestSectionsService_Crosslist_WithOverride(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Query().Get("override_sis_stickiness") != "true" {
			t.Errorf("expected override_sis_stickiness=true, got %q", r.URL.Query().Get("override_sis_stickiness"))
		}
		json.NewEncoder(w).Encode(Section{ID: 5, Name: "Crosslisted", CourseID: 99})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewSectionsService(client)
	section, err := svc.Crosslist(context.Background(), 5, 99, true)
	if err != nil {
		t.Fatalf("Crosslist with override: %v", err)
	}
	if section.CourseID != 99 {
		t.Errorf("expected CourseID 99, got %d", section.CourseID)
	}
}

func TestSectionsService_Uncrosslist_WithOverride(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Query().Get("override_sis_stickiness") != "true" {
			t.Errorf("expected override_sis_stickiness=true, got %q", r.URL.Query().Get("override_sis_stickiness"))
		}
		json.NewEncoder(w).Encode(Section{ID: 5, Name: "Uncrosslisted", CourseID: 1})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewSectionsService(client)
	section, err := svc.Uncrosslist(context.Background(), 5, true)
	if err != nil {
		t.Fatalf("Uncrosslist with override: %v", err)
	}
	if section.CourseID != 1 {
		t.Errorf("expected CourseID 1, got %d", section.CourseID)
	}
}

// TestSectionsService_GradeSubmission_FullBody asserts that GradeSubmission
// sends the full body including SecondsLateOverride, RubricAssessment, and all
// comment fields — not just PostedGrade/Excuse/LatePolicyStatus/TextComment.
func TestSectionsService_GradeSubmission_FullBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/sections/5/assignments/10/submissions/99" || r.Method != http.MethodPut {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		sub, ok := body["submission"].(map[string]interface{})
		if !ok {
			t.Fatal("expected submission object")
		}
		if sub["posted_grade"] != "A" {
			t.Errorf("expected posted_grade=A, got %v", sub["posted_grade"])
		}
		if sub["seconds_late_override"] == nil {
			t.Error("expected seconds_late_override set")
		}
		comment, ok := body["comment"].(map[string]interface{})
		if !ok {
			t.Fatal("expected comment object")
		}
		if comment["text_comment"] != "Great work" {
			t.Errorf("expected text_comment, got %v", comment["text_comment"])
		}
		if comment["group_comment"] != true {
			t.Errorf("expected group_comment=true, got %v", comment["group_comment"])
		}
		rubric, ok := body["rubric_assessment"].(map[string]interface{})
		if !ok {
			t.Fatal("expected rubric_assessment object")
		}
		criterion, ok := rubric["criterion_1"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected criterion_1 in rubric_assessment, got %v", rubric)
		}
		if criterion["points"] == nil {
			t.Error("expected points in criterion")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Submission{ID: 999})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewSectionsService(client)
	secondsLate := 3600
	sub, err := svc.GradeSubmission(context.Background(), 5, 10, 99, &GradeSubmissionParams{
		PostedGrade:         "A",
		SecondsLateOverride: &secondsLate,
		Comment: &SubmissionCommentParams{
			TextComment:  "Great work",
			GroupComment: true,
		},
		RubricAssessment: map[string]RubricAssessmentParams{
			"criterion_1": {Points: 5.0, Rating: "full_marks", Comments: "Excellent"},
		},
	})
	if err != nil {
		t.Fatalf("GradeSubmission full body: %v", err)
	}
	if sub.ID != 999 {
		t.Errorf("expected submission ID 999, got %d", sub.ID)
	}
}

func TestSectionsService_ListCourse_WithPageOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("page") != "2" {
			t.Errorf("expected page=2, got %q", q.Get("page"))
		}
		if q.Get("per_page") != "15" {
			t.Errorf("expected per_page=15, got %q", q.Get("per_page"))
		}
		json.NewEncoder(w).Encode([]Section{{ID: 1, Name: "Paged Section", CourseID: 10}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewSectionsService(client)
	opts := &ListSectionsOptions{Page: 2, PerPage: 15}
	sections, err := svc.ListCourse(context.Background(), 10, opts)
	if err != nil {
		t.Fatalf("ListCourse: %v", err)
	}
	if len(sections) != 1 {
		t.Errorf("expected 1, got %d", len(sections))
	}
}
