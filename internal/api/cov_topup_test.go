package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestTopup_GroupsService_ListCourseWithOpts exercises the opts != nil branch in ListCourse.
func TestTopup_GroupsService_ListCourseWithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/api/v1/courses/1/groups") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Group{{ID: 5, Name: "G5"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	opts := &ListGroupsOptions{Include: []string{"users"}, Page: 1, PerPage: 10}
	groups, err := svc.ListCourse(context.Background(), 1, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 || groups[0].ID != 5 {
		t.Errorf("expected 1 group with ID 5, got %v", groups)
	}
}

// TestTopup_GroupsService_ListAccountWithOpts exercises the opts != nil branch in ListAccount.
func TestTopup_GroupsService_ListAccountWithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/api/v1/accounts/2/groups") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Group{{ID: 7, Name: "G7"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	opts := &ListGroupsOptions{PerPage: 5}
	groups, err := svc.ListAccount(context.Background(), 2, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 || groups[0].ID != 7 {
		t.Errorf("expected 1 group with ID 7, got %v", groups)
	}
}

// TestTopup_GroupsService_ListCategoriesCourseWithOpts exercises the opts != nil branch.
func TestTopup_GroupsService_ListCategoriesCourseWithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/api/v1/courses/3/group_categories") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]GroupCategory{{ID: 1, Name: "Cat1"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)

	opts := &ListCategoriesOptions{Page: 2, PerPage: 20}
	cats, err := svc.ListCategoriesCourse(context.Background(), 3, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cats) != 1 || cats[0].ID != 1 {
		t.Errorf("expected 1 category, got %v", cats)
	}
}

// TestTopup_GroupsService_ListCourseNilOpts exercises the nil opts path in ListCourse.
func TestTopup_GroupsService_ListCourseNilOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Group{})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewGroupsService(client)
	_, err := svc.ListCourse(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestTopup_OutcomesService_ListGroupsAccountWithOpts exercises the opts != nil branch.
func TestTopup_OutcomesService_ListGroupsAccountWithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/api/v1/accounts/1/outcome_groups") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]OutcomeGroup{{ID: 10, Title: "Root"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)

	opts := &ListOutcomeGroupsOptions{Page: 1, PerPage: 10}
	groups, err := svc.ListGroupsAccount(context.Background(), 1, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 || groups[0].ID != 10 {
		t.Errorf("expected 1 group, got %v", groups)
	}
}

// TestTopup_OutcomesService_ListGroupsCourseWithOpts exercises the opts != nil branch.
func TestTopup_OutcomesService_ListGroupsCourseWithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/api/v1/courses/5/outcome_groups") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]OutcomeGroup{{ID: 20, Title: "CourseGroup"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)

	opts := &ListOutcomeGroupsOptions{PerPage: 5}
	groups, err := svc.ListGroupsCourse(context.Background(), 5, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 || groups[0].ID != 20 {
		t.Errorf("expected 1 group, got %v", groups)
	}
}

// TestTopup_OutcomesService_GetResults exercises GetResults with and without options.
func TestTopup_OutcomesService_GetResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/api/v1/courses/7/outcome_results") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OutcomeResultsResponse{
			OutcomeResults: []OutcomeResult{{ID: 1}},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)

	// nil opts path
	resp, err := svc.GetResults(context.Background(), 7, nil)
	if err != nil {
		t.Fatalf("unexpected error (nil opts): %v", err)
	}
	if len(resp.OutcomeResults) != 1 {
		t.Errorf("expected 1 result, got %d", len(resp.OutcomeResults))
	}

	// opts with values path
	opts := &OutcomeResultsOptions{
		UserIDs:       []int64{1, 2},
		OutcomeIDs:    []int64{10},
		Include:       []string{"alignments"},
		IncludeHidden: true,
		Page:          1,
		PerPage:       10,
	}
	resp2, err := svc.GetResults(context.Background(), 7, opts)
	if err != nil {
		t.Fatalf("unexpected error (with opts): %v", err)
	}
	if len(resp2.OutcomeResults) != 1 {
		t.Errorf("expected 1 result, got %d", len(resp2.OutcomeResults))
	}
}

// TestTopup_OutcomesService_ListGroupsAccountNilOpts exercises the nil opts path.
func TestTopup_OutcomesService_ListGroupsAccountNilOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]OutcomeGroup{})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	_, err := svc.ListGroupsAccount(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestTopup_OverridesService_ListWithOpts exercises the opts != nil branch in overrides List.
func TestTopup_OverridesService_ListWithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/api/v1/courses/1/assignments/2/overrides") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]AssignmentOverride{{ID: 1, Title: "Override1"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOverridesService(client)

	opts := &AssignmentOverrideListOptions{Page: 1, PerPage: 10}
	overrides, err := svc.List(context.Background(), 1, 2, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(overrides) != 1 || overrides[0].ID != 1 {
		t.Errorf("expected 1 override, got %v", overrides)
	}
}

// TestTopup_OverridesService_ListNilOpts exercises the nil opts path.
func TestTopup_OverridesService_ListNilOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]AssignmentOverride{})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOverridesService(client)
	_, err := svc.List(context.Background(), 1, 2, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestTopup_QuizQuestionsService_ListWithOpts exercises the opts != nil branch.
func TestTopup_QuizQuestionsService_ListWithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/api/v1/courses/1/quizzes/2/questions") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]QuizQuestion{{ID: 1, QuestionName: "Q1"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizQuestionsService(client)

	opts := &ListQuizQuestionsOptions{Page: 1, PerPage: 25}
	questions, err := svc.List(context.Background(), 1, 2, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(questions) != 1 || questions[0].ID != 1 {
		t.Errorf("expected 1 question, got %v", questions)
	}
}

// TestTopup_OutcomesService_ListOutcomesInGroupAccountWithOpts exercises the opts != nil branch.
func TestTopup_OutcomesService_ListOutcomesInGroupAccountWithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/api/v1/accounts/1/outcome_groups/5/outcomes") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]OutcomeLink{{OutcomeGroup: &OutcomeGroup{ID: 5}}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)

	opts := &ListOutcomesInGroupOptions{Page: 1, PerPage: 10}
	links, err := svc.ListOutcomesInGroupAccount(context.Background(), 1, 5, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 1 {
		t.Errorf("expected 1 outcome link, got %d", len(links))
	}
}
