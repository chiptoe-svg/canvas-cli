package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorReportsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/error_reports" || r.Method != http.MethodPost {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		result := ErrorReportResult{Logged: true, ID: "abc123"}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewErrorReportsService(client)

	report := &ErrorReport{
		Subject:  "Test bug",
		Comments: "Something broke",
	}

	result, err := svc.Create(context.Background(), report)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if !result.Logged {
		t.Error("expected Logged=true")
	}
}
