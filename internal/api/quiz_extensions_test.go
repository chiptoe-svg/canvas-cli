package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuizExtensionsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/extensions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		extensions, ok := body["quiz_extensions"].([]interface{})
		if !ok || len(extensions) == 0 {
			t.Errorf("expected quiz_extensions array in body")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizExtensionsResponse{
			QuizExtensions: []QuizExtension{
				{UserID: 42, QuizID: 2, ExtraTime: 30},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizExtensionsService(client)

	entries := []QuizExtensionEntry{
		{UserID: 42, ExtraTime: 30},
	}
	exts, err := svc.Create(context.Background(), 1, 2, entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exts) != 1 {
		t.Fatalf("expected 1 extension, got %d", len(exts))
	}
	if exts[0].UserID != 42 {
		t.Errorf("expected UserID 42, got %d", exts[0].UserID)
	}
	if exts[0].ExtraTime != 30 {
		t.Errorf("expected ExtraTime 30, got %d", exts[0].ExtraTime)
	}
}
