package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuizIPFiltersService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/ip_filters" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizIPFiltersResponse{
			QuizIPFilters: []QuizIPFilter{
				{Name: "Campus Network", Filter: "192.168.0.0/24"},
				{Name: "Library", Filter: "10.0.0.0/8"},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizIPFiltersService(client)

	filters, err := svc.List(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Errorf("expected 2 IP filters, got %d", len(filters))
	}
	if filters[0].Name != "Campus Network" {
		t.Errorf("expected 'Campus Network', got %q", filters[0].Name)
	}
	if filters[0].Filter != "192.168.0.0/24" {
		t.Errorf("expected '192.168.0.0/24', got %q", filters[0].Filter)
	}
}
