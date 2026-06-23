package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsCanvasDomain_SameHost(t *testing.T) {
	tests := []struct {
		redirectURL string
		baseURL     string
		want        bool
	}{
		{
			redirectURL: "https://canvas.example.com/files/1/download",
			baseURL:     "https://canvas.example.com",
			want:        true,
		},
		{
			redirectURL: "https://s3.amazonaws.com/bucket/file.pdf",
			baseURL:     "https://canvas.example.com",
			want:        false,
		},
		{
			redirectURL: "https://canvas.example.com/api/v1/files",
			baseURL:     "https://canvas.example.com/something",
			want:        true,
		},
		{
			redirectURL: "https://other.example.com/files/1",
			baseURL:     "https://canvas.example.com",
			want:        false,
		},
	}

	for _, tt := range tests {
		got := isCanvasDomain(tt.redirectURL, tt.baseURL)
		if got != tt.want {
			t.Errorf("isCanvasDomain(%q, %q) = %v, want %v",
				tt.redirectURL, tt.baseURL, got, tt.want)
		}
	}
}

func TestIsCanvasDomain_InvalidURL(t *testing.T) {
	// An unparseable redirect URL returns false.
	got := isCanvasDomain("://not-valid", "https://canvas.example.com")
	if got != false {
		t.Error("expected false for invalid redirect URL")
	}

	// An unparseable base URL returns false.
	got = isCanvasDomain("https://canvas.example.com/files/1", "://not-valid")
	if got != false {
		t.Error("expected false for invalid base URL")
	}
}

// TestFilesService_UpdateCourseDateDetails_NullClears asserts that nil date
// pointer fields marshal as JSON null (not omitted), so a PUT with nil DueAt
// explicitly clears the date on Canvas. Previously the struct had omitempty on
// the date fields, making it impossible to send null to clear a date.
func TestFilesService_UpdateCourseDateDetails_NullClears(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/files/42/date_details" || r.Method != http.MethodPut {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		var body map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		// due_at must be present and explicitly null (not omitted).
		raw, ok := body["due_at"]
		if !ok {
			t.Error("expected due_at key in body (should be null, not omitted)")
		} else if string(raw) != "null" {
			t.Errorf("expected due_at=null, got %s", raw)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FileDateDetails{})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewFilesService(client)
	// nil DueAt should marshal as null, explicitly clearing the date.
	_, err := svc.UpdateCourseDateDetails(context.Background(), 1, 42, &FileDateDetails{
		DueAt: nil, // intentional null-clear
	})
	if err != nil {
		t.Fatalf("UpdateCourseDateDetails: %v", err)
	}
}
