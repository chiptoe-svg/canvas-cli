package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOutcomesService_ListGlobalGroupSubgroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/global/outcome_groups/5/subgroups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]OutcomeGroup{{ID: 10, Title: "Subgroup A"}, {ID: 11, Title: "Subgroup B"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	groups, err := svc.ListGlobalGroupSubgroups(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListGlobalGroupSubgroups: %v", err)
	}
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Title != "Subgroup A" {
		t.Errorf("expected 'Subgroup A', got %s", groups[0].Title)
	}
}

func TestOutcomesService_ListGlobalGroupSubgroups_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	_, err := svc.ListGlobalGroupSubgroups(context.Background(), 5)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOutcomesService_CreateGlobalGroupSubgroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/global/outcome_groups/5/subgroups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body OutcomeGroupParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Title != "New Global Subgroup" {
			t.Errorf("expected title 'New Global Subgroup', got %s", body.Title)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OutcomeGroup{ID: 20, Title: "New Global Subgroup"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	group, err := svc.CreateGlobalGroupSubgroup(context.Background(), 5, &OutcomeGroupParams{Title: "New Global Subgroup"})
	if err != nil {
		t.Fatalf("CreateGlobalGroupSubgroup: %v", err)
	}
	if group.ID != 20 {
		t.Errorf("expected ID 20, got %d", group.ID)
	}
}

func TestOutcomesService_ListGlobalGroupOutcomes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/global/outcome_groups/5/outcomes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]OutcomeLink{{Outcome: &Outcome{ID: 1, Title: "Outcome 1"}}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	links, err := svc.ListGlobalGroupOutcomes(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListGlobalGroupOutcomes: %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(links))
	}
	if links[0].Outcome.Title != "Outcome 1" {
		t.Errorf("expected Outcome 1, got %s", links[0].Outcome.Title)
	}
}

func TestOutcomesService_CreateGlobalGroupOutcome(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/global/outcome_groups/5/outcomes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OutcomeLink{Outcome: &Outcome{ID: 30, Title: "New Outcome"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	link, err := svc.CreateGlobalGroupOutcome(context.Background(), 5, &CreateOutcomeParams{Title: "New Outcome"})
	if err != nil {
		t.Fatalf("CreateGlobalGroupOutcome: %v", err)
	}
	if link.Outcome.ID != 30 {
		t.Errorf("expected ID 30, got %d", link.Outcome.ID)
	}
}

func TestOutcomesService_DeleteGlobalGroupOutcome(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/global/outcome_groups/5/outcomes/10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OutcomeLink{Outcome: &Outcome{ID: 10}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	link, err := svc.DeleteGlobalGroupOutcome(context.Background(), 5, 10)
	if err != nil {
		t.Fatalf("DeleteGlobalGroupOutcome: %v", err)
	}
	if link.Outcome.ID != 10 {
		t.Errorf("expected ID 10, got %d", link.Outcome.ID)
	}
}

func TestOutcomesService_LinkGlobalGroupOutcome(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/global/outcome_groups/5/outcomes/10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OutcomeLink{Outcome: &Outcome{ID: 10}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	link, err := svc.LinkGlobalGroupOutcome(context.Background(), 5, 10)
	if err != nil {
		t.Fatalf("LinkGlobalGroupOutcome: %v", err)
	}
	if link.Outcome.ID != 10 {
		t.Errorf("expected ID 10, got %d", link.Outcome.ID)
	}
}

func TestOutcomesService_ImportGlobalGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/global/outcome_groups/5/import" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OutcomeGroup{ID: 5, Title: "Imported Group"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	group, err := svc.ImportGlobalGroup(context.Background(), 5, &OutcomeGroupImportParams{SourceOutcomeGroupID: 99})
	if err != nil {
		t.Fatalf("ImportGlobalGroup: %v", err)
	}
	if group.ID != 5 {
		t.Errorf("expected ID 5, got %d", group.ID)
	}
}

func TestOutcomesService_UpdateGroupAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_groups/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OutcomeGroup{ID: 5, Title: "Updated Account Group"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	group, err := svc.UpdateGroupAccount(context.Background(), 1, 5, &OutcomeGroupParams{Title: "Updated Account Group"})
	if err != nil {
		t.Fatalf("UpdateGroupAccount: %v", err)
	}
	if group.Title != "Updated Account Group" {
		t.Errorf("expected 'Updated Account Group', got %s", group.Title)
	}
}

func TestOutcomesService_ImportGroupAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_groups/5/import" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OutcomeGroup{ID: 5})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	group, err := svc.ImportGroupAccount(context.Background(), 1, 5, &OutcomeGroupImportParams{SourceOutcomeGroupID: 3})
	if err != nil {
		t.Fatalf("ImportGroupAccount: %v", err)
	}
	if group.ID != 5 {
		t.Errorf("expected ID 5, got %d", group.ID)
	}
}

func TestOutcomesService_ListGroupSubgroupsAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_groups/5/subgroups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]OutcomeGroup{{ID: 10, Title: "Account Subgroup"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	groups, err := svc.ListGroupSubgroupsAccount(context.Background(), 1, 5)
	if err != nil {
		t.Fatalf("ListGroupSubgroupsAccount: %v", err)
	}
	if len(groups) != 1 || groups[0].Title != "Account Subgroup" {
		t.Errorf("unexpected result: %v", groups)
	}
}

func TestOutcomesService_CreateGroupSubgroupAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_groups/5/subgroups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OutcomeGroup{ID: 20, Title: "New Account Subgroup"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	group, err := svc.CreateGroupSubgroupAccount(context.Background(), 1, 5, &OutcomeGroupParams{Title: "New Account Subgroup"})
	if err != nil {
		t.Fatalf("CreateGroupSubgroupAccount: %v", err)
	}
	if group.ID != 20 {
		t.Errorf("expected ID 20, got %d", group.ID)
	}
}

func TestOutcomesService_UpdateGroupCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/outcome_groups/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OutcomeGroup{ID: 5, Title: "Updated Course Group"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	group, err := svc.UpdateGroupCourse(context.Background(), 10, 5, &OutcomeGroupParams{Title: "Updated Course Group"})
	if err != nil {
		t.Fatalf("UpdateGroupCourse: %v", err)
	}
	if group.Title != "Updated Course Group" {
		t.Errorf("expected 'Updated Course Group', got %s", group.Title)
	}
}

func TestOutcomesService_ImportGroupCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/outcome_groups/5/import" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OutcomeGroup{ID: 5})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	group, err := svc.ImportGroupCourse(context.Background(), 10, 5, &OutcomeGroupImportParams{SourceOutcomeGroupID: 7})
	if err != nil {
		t.Fatalf("ImportGroupCourse: %v", err)
	}
	if group.ID != 5 {
		t.Errorf("expected ID 5, got %d", group.ID)
	}
}

func TestOutcomesService_ListGroupSubgroupsCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/outcome_groups/5/subgroups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]OutcomeGroup{{ID: 30, Title: "Course Subgroup"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	groups, err := svc.ListGroupSubgroupsCourse(context.Background(), 10, 5)
	if err != nil {
		t.Fatalf("ListGroupSubgroupsCourse: %v", err)
	}
	if len(groups) != 1 || groups[0].Title != "Course Subgroup" {
		t.Errorf("unexpected result: %v", groups)
	}
}

func TestOutcomesService_CreateGroupSubgroupCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/outcome_groups/5/subgroups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OutcomeGroup{ID: 40, Title: "New Course Subgroup"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	group, err := svc.CreateGroupSubgroupCourse(context.Background(), 10, 5, &OutcomeGroupParams{Title: "New Course Subgroup"})
	if err != nil {
		t.Fatalf("CreateGroupSubgroupCourse: %v", err)
	}
	if group.ID != 40 {
		t.Errorf("expected ID 40, got %d", group.ID)
	}
}
