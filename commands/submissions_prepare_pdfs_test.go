package commands

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/pdfprep"
)

type preparePDFRunner struct{}

func (preparePDFRunner) Run(_ context.Context, name string, args ...string) (pdfprep.CommandResult, error) {
	switch name {
	case "pdfinfo":
		return pdfprep.CommandResult{Stdout: "Pages:          1\n"}, nil
	case "pdftotext":
		return pdfprep.CommandResult{Stdout: "student notes"}, nil
	case "pdffonts":
		return pdfprep.CommandResult{}, nil
	case "pdfimages":
		if len(args) == 3 && args[0] == "-all" {
			return pdfprep.CommandResult{}, os.WriteFile(args[2]+"-000.jpg", []byte("original"), 0o600)
		}
		return pdfprep.CommandResult{Stdout: "   0   1 image\n"}, nil
	}
	return pdfprep.CommandResult{}, nil
}

func TestRunSubmissionsPreparePDFsWritesLocalManifest(t *testing.T) {
	folder := t.TempDir()
	output := filepath.Join(t.TempDir(), "review")
	studentFolder := filepath.Join(folder, "user-42")
	if err := os.Mkdir(studentFolder, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(studentFolder, "notes.pdf"), []byte("not a real PDF; runner is faked"), 0o600); err != nil {
		t.Fatal(err)
	}
	opts := &options.SubmissionsPreparePDFsOptions{Folder: folder, Output: output}
	if err := runSubmissionsPreparePDFs(context.Background(), opts, pdfprep.Preparer{Runner: preparePDFRunner{}}); err != nil {
		t.Fatalf("runSubmissionsPreparePDFs: %v", err)
	}
	manifestBytes, err := os.ReadFile(filepath.Join(output, submissionPDFManifestName))
	if err != nil {
		t.Fatal(err)
	}
	// JSON Lines: the header record, then one record per PDF.
	lines := strings.Split(strings.TrimSpace(string(manifestBytes)), "\n")
	if len(lines) != 2 {
		t.Fatalf("want a header line plus one entry line, got %d:\n%s", len(lines), manifestBytes)
	}
	var header submissionPDFManifest
	if err := json.Unmarshal([]byte(lines[0]), &header); err != nil {
		t.Fatalf("header line: %v", err)
	}
	var entry submissionPDFRecord
	if err := json.Unmarshal([]byte(lines[1]), &entry); err != nil {
		t.Fatalf("entry line: %v", err)
	}
	if entry.Status != "prepared" {
		t.Fatalf("unexpected entry status: %#v", entry)
	}
	if entry.SourceRelativePath != "user-42/notes.pdf" || entry.PageImages.Source != pdfprep.EmbeddedImagePages {
		t.Fatalf("unexpected entry: %#v", entry)
	}
	if len(entry.PageImages.Paths) != 1 {
		t.Fatalf("page images: %#v", entry.PageImages)
	}
	if _, err := os.Stat(filepath.Join(output, entry.PageImages.Paths[0])); err != nil {
		t.Fatalf("prepared page image missing: %v", err)
	}
}

func TestSubmissionsPreparePDFsOptionsValidate(t *testing.T) {
	if err := (&options.SubmissionsPreparePDFsOptions{Folder: "submissions", Output: "review"}).Validate(); err != nil {
		t.Fatalf("valid options: %v", err)
	}
	if err := (&options.SubmissionsPreparePDFsOptions{Output: "review"}).Validate(); err == nil {
		t.Fatal("missing folder must fail")
	}
	if err := (&options.SubmissionsPreparePDFsOptions{Folder: "submissions"}).Validate(); err == nil {
		t.Fatal("missing output must fail")
	}
}

// The manifest is named .jsonl, so it must BE JSON Lines: one parsable record
// per line, header first. It used to be a single indented document, which no
// line-oriented reader could consume.
func TestWriteSubmissionPDFManifest_IsJSONLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, submissionPDFManifestName)
	manifest := submissionPDFManifest{
		SchemaVersion: 1,
		SourceFolder:  "/tmp/work",
		Entries: []submissionPDFRecord{
			{SourceRelativePath: "user-42/a.pdf", Status: "prepared"},
			{SourceRelativePath: "user-43/b.pdf", Status: "prepared"},
		},
	}
	if err := writeSubmissionPDFManifest(path, manifest); err != nil {
		t.Fatalf("writeSubmissionPDFManifest: %v", err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	if len(lines) != 3 {
		t.Fatalf("want a header line plus one line per entry (3), got %d:\n%s", len(lines), raw)
	}
	for i, line := range lines {
		var record map[string]interface{}
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("line %d is not a JSON record: %v\n%s", i, err, line)
		}
		if i == 0 {
			if _, present := record["entries"]; present {
				t.Error("header line must not carry an entries key")
			}
			continue
		}
		if record["source_relative_path"] == nil {
			t.Errorf("line %d is missing its source path: %s", i, line)
		}
	}
}

// The Poppler check belongs to the real command, not the injected run
// function: a CI box without Poppler must still exercise the logic with a fake
// runner. This pins the message so the check is not quietly lost either.
func TestRequirePopplerTools_NamesTheFix(t *testing.T) {
	err := requirePopplerTools()
	if err == nil {
		return // Poppler is installed here; nothing to assert about the message
	}
	for _, want := range []string{"pdftotext", "brew install poppler"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("missing-tools error should mention %q: %v", want, err)
		}
	}
}
