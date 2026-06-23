package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHistoryService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/1/history" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		entries := []HistoryEntry{
			{AssetCode: "assignment_42", AssetName: "Midterm Essay", ContextType: "Course"},
			{AssetCode: "wiki_page_7", AssetName: "Syllabus", ContextType: "Course"},
		}
		if err := json.NewEncoder(w).Encode(entries); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewHistoryService(client)

	entries, err := svc.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].AssetCode != "assignment_42" {
		t.Errorf("unexpected asset code: %s", entries[0].AssetCode)
	}
}
