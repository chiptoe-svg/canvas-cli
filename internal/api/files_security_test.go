package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSanitizeDestPath verifies the internal path validation helper.
func TestSanitizeDestPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"normal filename", "file.pdf", false},
		{"path with directory", "subdir/file.pdf", false},
		{"absolute path", "/tmp/file.pdf", false},
		{"empty string", "", true},
		{"bare dot", ".", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sanitizeDestPath(tt.path)
			if tt.wantErr && err == nil {
				t.Errorf("sanitizeDestPath(%q): expected error, got nil", tt.path)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("sanitizeDestPath(%q): unexpected error: %v", tt.path, err)
			}
		})
	}
}

// TestFilesService_Download_PathTraversal ensures that Download rejects
// an empty destPath (which would be the only unsafe value it receives after
// the command layer has already sanitised the server-supplied filename).
func TestFilesService_Download_PathTraversal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		// Return a file record with a server-controlled filename.
		file := Attachment{
			ID:          1,
			Filename:    "../../etc/passwd",
			DisplayName: "evil",
			URL:         "http://example.com/download",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(file)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	// Download with empty destPath should be rejected immediately.
	ctx := context.Background()
	svc := NewFilesService(client)
	if err := svc.Download(ctx, 1, ""); err == nil {
		t.Error("Download with empty destPath: expected error, got nil")
	}
}

// TestFilesService_Download_WritesToExpectedPath confirms that a normal
// download (non-malicious filename) writes the file where specified.
func TestFilesService_Download_WritesToExpectedPath(t *testing.T) {
	content := []byte("hello world")

	// Track request count to distinguish file-info from download calls.
	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		requestCount++

		if strings.HasPrefix(r.URL.Path, "/api/v1/files/") {
			// Return file metadata — URL points back to our test server.
			file := map[string]interface{}{
				"id":           1,
				"filename":     "report.pdf",
				"display_name": "report.pdf",
				"url":          "http://" + r.Host + "/download",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(file)
			return
		}

		if r.URL.Path == "/download" {
			w.WriteHeader(http.StatusOK)
			w.Write(content)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	tmpDir := t.TempDir()
	dest := filepath.Join(tmpDir, "report.pdf")

	ctx := context.Background()
	svc := NewFilesService(client)
	if err := svc.Download(ctx, 1, dest); err != nil {
		t.Fatalf("Download: %v", err)
	}

	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != string(content) {
		t.Errorf("expected %q, got %q", content, data)
	}
}
