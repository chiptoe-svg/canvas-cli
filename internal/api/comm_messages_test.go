package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommMessagesService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/comm_messages" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		msgs := []CommMessage{
			{ID: 1, Subject: "Welcome", WorkflowState: "sent"},
			{ID: 2, Subject: "Assignment due", WorkflowState: "sent"},
		}
		if err := json.NewEncoder(w).Encode(msgs); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCommMessagesService(client)

	msgs, err := svc.List(context.Background(), &ListCommMessagesOptions{UserID: 5})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(msgs) != 2 {
		t.Errorf("expected 2 messages, got %d", len(msgs))
	}
	if msgs[0].Subject != "Welcome" {
		t.Errorf("unexpected subject: %s", msgs[0].Subject)
	}
}
