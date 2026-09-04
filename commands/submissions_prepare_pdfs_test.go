package commands

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
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
	var manifest submissionPDFManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatal(err)
	}
	if len(manifest.Entries) != 1 || manifest.Entries[0].Status != "prepared" {
		t.Fatalf("unexpected manifest: %#v", manifest)
	}
	entry := manifest.Entries[0]
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
