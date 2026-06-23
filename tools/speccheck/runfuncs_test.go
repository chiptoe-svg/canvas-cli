package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// repoRoot returns the path two directories up from the tools/speccheck package dir.
// Go tests run with cwd = the package directory.
func repoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// tools/speccheck → repo root
	return filepath.Join(cwd, "..", "..")
}

// TestLoadManifestFromRepoRoot calls loadManifest with the working directory set
// to the repo root so that the relative manifestPath constant resolves correctly.
func TestLoadManifestFromRepoRoot(t *testing.T) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	if err := os.Chdir(repoRoot(t)); err != nil {
		t.Fatalf("chdir to repo root: %v", err)
	}

	man, err := loadManifest()
	if err != nil {
		t.Fatalf("loadManifest: %v", err)
	}
	if len(man.Endpoints) == 0 {
		t.Error("expected non-empty endpoints from committed manifest")
	}
}

// TestRunCoverage calls runCoverage from the repo root so all relative paths resolve.
// It simply verifies the function returns no error (manifest exists, CLI paths are harvested).
func TestRunCoverage(t *testing.T) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	if err := os.Chdir(repoRoot(t)); err != nil {
		t.Fatalf("chdir to repo root: %v", err)
	}

	if err := runCoverage(); err != nil {
		t.Fatalf("runCoverage: %v", err)
	}
}

// TestWriteManifest_Integration tests the writeManifest function by temporarily
// redirecting the output to a writable temp path via a sub-process-free approach.
// Since manifestPath is a const, we exercise writeManifest's underlying writeJSON logic
// by calling it with a path we control (same code path).
func TestWriteManifest_Integration(t *testing.T) {
	dir := t.TempDir()

	// Patch the global manifestPath for this test by shadowing it.
	// We call writeJSON directly (same code writeManifest calls).
	outPath := filepath.Join(dir, "testdata", "spec", "canvas_endpoints.json")

	endpoints := []Endpoint{
		{Method: "GET", Path: "/api/v1/courses"},
	}
	man := Manifest{
		GeneratedNote: "integration test",
		Source:        "https://test.example.com",
		Endpoints:     endpoints,
	}

	if err := writeJSON(outPath, man); err != nil {
		t.Fatalf("writeJSON: %v", err)
	}

	raw, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}

	var got Manifest
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Source != "https://test.example.com" {
		t.Errorf("Source mismatch: %q", got.Source)
	}
}

// TestWriteModels_Integration exercises the writeModels-equivalent code path.
func TestWriteModels_Integration(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "models.json")

	models := map[string]*ModelDef{
		"Assignment": {
			Resource: "assignments",
			Fields: map[string]ModelField{
				"id":    {Type: "integer"},
				"title": {Type: "string"},
			},
		},
	}
	mm := ModelsManifest{
		GeneratedNote: "test",
		Source:        "test",
		Models:        models,
	}
	if err := writeJSON(outPath, mm); err != nil {
		t.Fatalf("writeJSON models: %v", err)
	}

	raw, _ := os.ReadFile(outPath)
	var got ModelsManifest
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := got.Models["Assignment"]; !ok {
		t.Error("expected 'Assignment' model in output")
	}
}

// TestLoadManifest_ErrorPaths verifies that loadManifest returns meaningful errors
// for invalid or missing files.
func TestLoadManifest_ErrorPaths(t *testing.T) {
	// Point at a temp dir that has no manifest. t.Chdir restores the working
	// directory at cleanup in the correct order relative to t.TempDir's
	// RemoveAll, which matters on Windows (a dir that is the process cwd cannot
	// be deleted).
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	_, err := loadManifest()
	if err == nil {
		t.Error("expected error when manifest file does not exist")
	}
}

// TestLoadManifest_EmptyEndpoints verifies that a manifest with 0 endpoints
// is rejected by loadManifest.
func TestLoadManifest_EmptyEndpoints(t *testing.T) {
	dir := t.TempDir()
	specDir := filepath.Join(dir, "testdata", "spec")
	if err := os.MkdirAll(specDir, 0o750); err != nil {
		t.Fatal(err)
	}
	empty := Manifest{
		GeneratedNote: "test",
		Source:        "test",
		Endpoints:     []Endpoint{},
	}
	data, _ := json.Marshal(empty)
	path := filepath.Join(specDir, "canvas_endpoints.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	_, err := loadManifest()
	if err == nil {
		t.Error("expected error for manifest with 0 endpoints")
	}
}

// TestLoadManifest_InvalidJSON verifies that corrupt JSON returns an error.
func TestLoadManifest_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	specDir := filepath.Join(dir, "testdata", "spec")
	if err := os.MkdirAll(specDir, 0o750); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(specDir, "canvas_endpoints.json")
	if err := os.WriteFile(path, []byte("{not valid json"), 0o600); err != nil {
		t.Fatal(err)
	}

	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	_, err := loadManifest()
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
