package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOutcomesService_GetRollupsCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/outcome_rollups" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"rollups": []map[string]interface{}{{"name": "Alice", "scores": []interface{}{}}},
		})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	result, err := svc.GetRollupsCourse(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("GetRollupsCourse failed: %v", err)
	}
	if len(result.Rollups) != 1 {
		t.Errorf("expected 1 rollup, got %d", len(result.Rollups))
	}
}

func TestOutcomesService_GetProficiencyCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/outcome_proficiency" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ratings": []map[string]interface{}{{"description": "Mastery", "points": 3.0, "mastery": true, "color": "0000FF"}},
		})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	result, err := svc.GetProficiencyCourse(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetProficiencyCourse failed: %v", err)
	}
	if len(result.Ratings) != 1 {
		t.Errorf("expected 1 rating, got %d", len(result.Ratings))
	}
}

func TestOutcomesService_UpdateProficiencyCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/outcome_proficiency" || r.Method != http.MethodPost {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ratings": []map[string]interface{}{{"description": "Mastery", "points": 3.0}},
		})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	result, err := svc.UpdateProficiencyCourse(context.Background(), 1, &OutcomeProficiencyParams{
		Ratings: []ProficiencyRating{{Description: "Mastery", Points: 3.0, Mastery: true}},
	})
	if err != nil {
		t.Fatalf("UpdateProficiencyCourse failed: %v", err)
	}
	if len(result.Ratings) != 1 {
		t.Errorf("expected 1 rating, got %d", len(result.Ratings))
	}
}

func TestOutcomesService_ListGroupLinksCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/outcome_group_links" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"url": "/api/v1/courses/1/outcome_group_links/1"}})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	links, err := svc.ListGroupLinksCourse(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListGroupLinksCourse failed: %v", err)
	}
	if len(links) != 1 {
		t.Errorf("expected 1 link, got %d", len(links))
	}
}

func TestOutcomesService_GetGlobalRootGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/global/root_outcome_group" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 1, "title": "Root"})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	group, err := svc.GetGlobalRootGroup(context.Background())
	if err != nil {
		t.Fatalf("GetGlobalRootGroup failed: %v", err)
	}
	if group.ID != 1 {
		t.Errorf("expected ID 1, got %d", group.ID)
	}
}

func TestOutcomesService_GetGlobalGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/global/outcome_groups/5" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 5, "title": "Global Group"})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	group, err := svc.GetGlobalGroup(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetGlobalGroup failed: %v", err)
	}
	if group.ID != 5 {
		t.Errorf("expected ID 5, got %d", group.ID)
	}
}

func TestOutcomesService_UpdateGlobalGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/global/outcome_groups/5" || r.Method != http.MethodPut {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 5, "title": "Updated"})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	group, err := svc.UpdateGlobalGroup(context.Background(), 5, &OutcomeGroupParams{Title: "Updated"})
	if err != nil {
		t.Fatalf("UpdateGlobalGroup failed: %v", err)
	}
	if group.Title != "Updated" {
		t.Errorf("expected title Updated, got %s", group.Title)
	}
}

func TestOutcomesService_DeleteGlobalGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/global/outcome_groups/5" || r.Method != http.MethodDelete {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	err := svc.DeleteGlobalGroup(context.Background(), 5)
	if err != nil {
		t.Fatalf("DeleteGlobalGroup failed: %v", err)
	}
}

func TestOutcomesService_GetRootGroupCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/root_outcome_group" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 10, "title": "Course Root"})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	group, err := svc.GetRootGroupCourse(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetRootGroupCourse failed: %v", err)
	}
	if group.ID != 10 {
		t.Errorf("expected ID 10, got %d", group.ID)
	}
}

func TestOutcomesService_GetRootGroupAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/root_outcome_group" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 20, "title": "Account Root"})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	group, err := svc.GetRootGroupAccount(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetRootGroupAccount failed: %v", err)
	}
	if group.ID != 20 {
		t.Errorf("expected ID 20, got %d", group.ID)
	}
}

func TestOutcomesService_ListGroupLinksAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_group_links" || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{{"url": "/accounts/1/links/1"}})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	links, err := svc.ListGroupLinksAccount(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListGroupLinksAccount failed: %v", err)
	}
	if len(links) != 1 {
		t.Errorf("expected 1 link, got %d", len(links))
	}
}

func TestOutcomesService_DeleteGroupCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/outcome_groups/5" || r.Method != http.MethodDelete {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	err := svc.DeleteGroupCourse(context.Background(), 1, 5)
	if err != nil {
		t.Fatalf("DeleteGroupCourse failed: %v", err)
	}
}

func TestOutcomesService_DeleteGroupAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_groups/5" || r.Method != http.MethodDelete {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewOutcomesService(client)
	err := svc.DeleteGroupAccount(context.Background(), 1, 5)
	if err != nil {
		t.Fatalf("DeleteGroupAccount failed: %v", err)
	}
}
