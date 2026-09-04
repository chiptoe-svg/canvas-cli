package commands

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

func TestSubmissionsDownloadOptionsValidate(t *testing.T) {
	tests := []struct {
		name    string
		opts    options.SubmissionsDownloadOptions
		wantErr bool
	}{
		{"valid", options.SubmissionsDownloadOptions{CourseID: 1, AssignmentID: 2, Destination: "downloads"}, false},
		{"missing course", options.SubmissionsDownloadOptions{AssignmentID: 2, Destination: "downloads"}, true},
		{"missing assignment", options.SubmissionsDownloadOptions{CourseID: 1, Destination: "downloads"}, true},
		{"missing destination", options.SubmissionsDownloadOptions{CourseID: 1, AssignmentID: 2}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.opts.Validate(); (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunSubmissionsDownload_ManifestAndResume(t *testing.T) {
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/courses/1":
			_, _ = w.Write([]byte(`{"id":1,"name":"Course"}`))
		case "/api/v1/courses/1/assignments/2/submissions":
			_, _ = w.Write([]byte(`[
                {"id":100,"user_id":42,"workflow_state":"submitted","submission_type":"online_upload","attachments":[{"id":10,"filename":"../essay.pdf"}]},
                {"id":101,"user_id":43,"workflow_state":"submitted","submission_type":"online_text_entry","body":"private student text"}
            ]`))
		case "/api/v1/files/10":
			_, _ = w.Write([]byte(`{"id":10,"filename":"essay.pdf","url":"` + serverURL + `/download/10"}`))
		case "/download/10":
			_, _ = w.Write([]byte("submission content"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client, err := api.NewClient(api.ClientConfig{BaseURL: server.URL, Token: "test-token", RequestsPerSec: 100})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	destination := t.TempDir()
	opts := &options.SubmissionsDownloadOptions{CourseID: 1, AssignmentID: 2, Destination: destination}

	if err := runSubmissionsDownload(context.Background(), client, opts); err != nil {
		t.Fatalf("runSubmissionsDownload() error = %v", err)
	}
	filePath := filepath.Join(destination, "user-42", "attachment-10-essay.pdf")
	contents, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", filePath, err)
	}
	if got, want := string(contents), "submission content"; got != want {
		t.Errorf("downloaded content = %q, want %q", got, want)
	}

	manifestBytes, err := os.ReadFile(filepath.Join(destination, submissionDownloadManifestName))
	if err != nil {
		t.Fatalf("ReadFile(manifest): %v", err)
	}
	if strings.Contains(string(manifestBytes), "private student text") {
		t.Fatal("manifest must not store submission body text")
	}
	var manifest submissionDownloadManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatalf("Unmarshal(manifest): %v", err)
	}
	if len(manifest.Entries) != 2 || manifest.Entries[0].Status != "downloaded" || manifest.Entries[1].Status != "no_attachment" {
		t.Fatalf("unexpected manifest entries: %#v", manifest.Entries)
	}

	if err := runSubmissionsDownload(context.Background(), client, opts); err != nil {
		t.Fatalf("resume run error = %v", err)
	}
	manifestBytes, err = os.ReadFile(filepath.Join(destination, submissionDownloadManifestName))
	if err != nil {
		t.Fatalf("ReadFile(resume manifest): %v", err)
	}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatalf("Unmarshal(resume manifest): %v", err)
	}
	if manifest.Entries[0].Status != "skipped_existing" {
		t.Fatalf("resume status = %q, want skipped_existing", manifest.Entries[0].Status)
	}
}

func TestSafeSubmissionAttachmentFilename(t *testing.T) {
	if got, want := safeSubmissionAttachmentFilename(api.Attachment{ID: 7, Filename: "../../report.pdf"}), "report.pdf"; got != want {
		t.Errorf("safe filename = %q, want %q", got, want)
	}
	if got, want := safeSubmissionAttachmentFilename(api.Attachment{ID: 8, Filename: ".."}), "file-8"; got != want {
		t.Errorf("fallback filename = %q, want %q", got, want)
	}
}

// A student who uploaded the parts of an assignment across separate attempts
// must get EVERY file, not just the ones on the latest attempt. Canvas returns
// the earlier files only under submission_history, and it repeats a
// carried-forward file in later attempts — so this pins both halves: nothing
// from an earlier attempt is lost, and nothing is downloaded twice.
func TestRunSubmissionsDownload_CollectsEveryAttempt(t *testing.T) {
	var serverURL string
	var listQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/courses/1":
			_, _ = w.Write([]byte(`{"id":1,"name":"Course"}`))
		case "/api/v1/courses/1/assignments/2/submissions":
			listQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[{
                "id":100,"user_id":42,"attempt":2,"workflow_state":"submitted","submission_type":"online_upload",
                "attachments":[{"id":11,"filename":"reviews.pdf"},{"id":10,"filename":"notes.pdf"}],
                "submission_history":[
                  {"id":100,"user_id":42,"attempt":1,"attachments":[{"id":10,"filename":"notes.pdf"}]},
                  {"id":100,"user_id":42,"attempt":2,"attachments":[{"id":11,"filename":"reviews.pdf"}]}
                ]
            }]`))
		case "/api/v1/files/10":
			_, _ = w.Write([]byte(`{"id":10,"filename":"notes.pdf","url":"` + serverURL + `/download/10"}`))
		case "/api/v1/files/11":
			_, _ = w.Write([]byte(`{"id":11,"filename":"reviews.pdf","url":"` + serverURL + `/download/11"}`))
		case "/download/10":
			_, _ = w.Write([]byte("attempt one notes"))
		case "/download/11":
			_, _ = w.Write([]byte("attempt two reviews"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client, err := api.NewClient(api.ClientConfig{BaseURL: server.URL, Token: "test-token", RequestsPerSec: 100})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	destination := t.TempDir()
	opts := &options.SubmissionsDownloadOptions{CourseID: 1, AssignmentID: 2, Destination: destination}
	if err := runSubmissionsDownload(context.Background(), client, opts); err != nil {
		t.Fatalf("runSubmissionsDownload() error = %v", err)
	}

	if !strings.Contains(listQuery, "submission_history") {
		t.Fatalf("submissions were listed without the history: query = %q", listQuery)
	}

	for name, want := range map[string]string{
		"attachment-10-notes.pdf":   "attempt one notes",
		"attachment-11-reviews.pdf": "attempt two reviews",
	} {
		got, err := os.ReadFile(filepath.Join(destination, "user-42", name))
		if err != nil {
			t.Fatalf("ReadFile(%s): %v", name, err)
		}
		if string(got) != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}

	manifestBytes, err := os.ReadFile(filepath.Join(destination, submissionDownloadManifestName))
	if err != nil {
		t.Fatalf("ReadFile(manifest): %v", err)
	}
	var manifest submissionDownloadManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatalf("Unmarshal(manifest): %v", err)
	}
	if len(manifest.Entries) != 2 {
		t.Fatalf("want one entry per distinct file, got %d: %#v", len(manifest.Entries), manifest.Entries)
	}
	attempts := map[int64]int{}
	for _, entry := range manifest.Entries {
		if entry.Status != "downloaded" {
			t.Errorf("attachment %d status = %q, want downloaded", entry.AttachmentID, entry.Status)
		}
		attempts[entry.AttachmentID] = entry.Attempt
	}
	if attempts[10] != 1 || attempts[11] != 2 {
		t.Errorf("manifest attempts = %v, want attachment 10 on attempt 1 and 11 on attempt 2", attempts)
	}
}
