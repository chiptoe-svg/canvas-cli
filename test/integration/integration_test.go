//go:build integration

// Package integration holds binary-level integration tests for canvas-cli.
// These tests compile the binary once in TestMain and then exercise it as a
// black box, using per-test isolated HOME directories and a lightweight
// httptest server to stand in for Canvas.
//
// Run with:
//
//	go test -tags integration -v -timeout 5m ./test/integration/
package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
)

// binaryPath holds the compiled canvas binary built once in TestMain.
var binaryPath string

// TestMain compiles the binary into a temp dir before any test runs.
func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "canvas-integration-bin-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmp)

	binName := "canvas"
	if runtime.GOOS == "windows" {
		binName = "canvas.exe"
	}
	binaryPath = filepath.Join(tmp, binName)

	// Resolve repo root from this file's location: test/integration → ../..
	_, file, _, _ := runtime.Caller(0)
	repoRoot := filepath.Join(filepath.Dir(file), "..", "..")

	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/canvas")
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build binary: %v\n%s\n", err, out)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

// runCanvas executes the compiled canvas binary with the given environment and
// arguments. Each call is fully isolated: callers supply their own HOME dir.
// Returns stdout, stderr, and the process exit code.
func runCanvas(t *testing.T, env []string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()

	cmd := exec.Command(binaryPath, args...)
	cmd.Env = env

	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	stdout = outBuf.String()
	stderr = errBuf.String()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return stdout, stderr, exitCode
}

// testEnv returns a minimal environment slice for an isolated run.
// homeDir is the per-test HOME/USERPROFILE directory (config/cache isolation).
// Extra key=value pairs may be appended via extras.
func testEnv(homeDir string, extras ...string) []string {
	env := []string{
		"HOME=" + homeDir,
		"USERPROFILE=" + homeDir, // Windows: os.UserHomeDir() reads %USERPROFILE%
		"CANVAS_CLI_MACHINE_ID=test-machine-id-integration",
		// PATH must be forwarded so exec.Command can find system tools (e.g. go
		// itself when the binary delegates subprocesses, or the shell on Windows).
		"PATH=" + os.Getenv("PATH"),
	}
	// Forward GOPATH/GOROOT so any internal go invocations can find the stdlib.
	if v := os.Getenv("GOPATH"); v != "" {
		env = append(env, "GOPATH="+v)
	}
	if v := os.Getenv("GOROOT"); v != "" {
		env = append(env, "GOROOT="+v)
	}
	env = append(env, extras...)
	return env
}

// mockServer returns a started httptest.Server that handles a small set of
// Canvas API endpoints with canned JSON responses. It is the caller's
// responsibility to call server.Close().
func mockServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

// authEnv returns a testEnv with CANVAS_URL and CANVAS_TOKEN pre-set.
func authEnv(homeDir, serverURL string) []string {
	return testEnv(homeDir,
		"CANVAS_URL="+serverURL,
		"CANVAS_TOKEN=test-token",
	)
}

// ─── Tests ────────────────────────────────────────────────────────────────────

// Case 1: canvas version → exit 0, output contains version string.
func TestVersion(t *testing.T) {
	home := t.TempDir()
	stdout, stderr, code := runCanvas(t, testEnv(home), "version")

	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if !strings.Contains(stdout, "canvas-cli") {
		t.Errorf("expected 'canvas-cli' in stdout, got: %q", stdout)
	}
}

// Case 2: canvas --help → exit 0, lists core commands.
func TestHelp(t *testing.T) {
	home := t.TempDir()
	stdout, stderr, code := runCanvas(t, testEnv(home), "--help")

	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	for _, cmd := range []string{"courses", "assignments", "auth"} {
		if !strings.Contains(stdout, cmd) {
			t.Errorf("expected %q in help output; got: %s", cmd, stdout)
		}
	}
}

// Case 3: unknown command → exit 1, stderr mentions "unknown".
func TestUnknownCommand(t *testing.T) {
	home := t.TempDir()
	_, stderr, code := runCanvas(t, testEnv(home), "this-command-does-not-exist")

	if code == 0 {
		t.Fatal("expected non-zero exit for unknown command")
	}
	combined := strings.ToLower(stderr)
	if !strings.Contains(combined, "unknown") {
		t.Errorf("expected 'unknown' in stderr, got: %q", stderr)
	}
}

// Case 4: courses list with mock server → exit 0, table output contains course name.
func TestCoursesListTableOutput(t *testing.T) {
	srv := mockServer(t, newCoursesHandler())
	home := t.TempDir()

	stdout, stderr, code := runCanvas(t, authEnv(home, srv.URL), "courses", "list")

	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr=%q stdout=%q", code, stderr, stdout)
	}
	if !strings.Contains(stdout, "Introduction to Go") {
		t.Errorf("expected course name in table output; got: %q", stdout)
	}
}

// Case 5: courses list -o json → stdout is a JSON array containing the course.
func TestCoursesListJSONOutput(t *testing.T) {
	srv := mockServer(t, newCoursesHandler())
	home := t.TempDir()

	stdout, stderr, code := runCanvas(t, authEnv(home, srv.URL), "courses", "list", "-o", "json")

	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr=%q", code, stderr)
	}

	var courses []map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &courses); err != nil {
		t.Fatalf("stdout is not valid JSON array: %v\nstdout: %q", err, stdout)
	}
	if len(courses) == 0 {
		t.Fatal("expected at least one course in JSON output")
	}
	name, _ := courses[0]["name"].(string)
	if name != "Introduction to Go" {
		t.Errorf("expected course name 'Introduction to Go', got: %q", name)
	}
}

// Case 6: courses list -o csv → header row present.
func TestCoursesListCSVOutput(t *testing.T) {
	srv := mockServer(t, newCoursesHandler())
	home := t.TempDir()

	stdout, stderr, code := runCanvas(t, authEnv(home, srv.URL), "courses", "list", "-o", "csv")

	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr=%q", code, stderr)
	}
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	if len(lines) < 1 {
		t.Fatal("expected at least a header row in CSV output")
	}
	header := strings.ToLower(lines[0])
	// CSV header must contain at least one known field
	if !strings.Contains(header, "id") && !strings.Contains(header, "name") {
		t.Errorf("CSV header row doesn't look like a header: %q", lines[0])
	}
}

// Case 7: missing required flag → exit 1, error mentions course-id.
func TestAssignmentsListMissingCourseID(t *testing.T) {
	home := t.TempDir()
	// No CANVAS_URL/CANVAS_TOKEN needed — validation fires before any HTTP call.
	stdout, stderr, code := runCanvas(t, testEnv(home), "assignments", "list")

	if code == 0 {
		t.Fatal("expected non-zero exit when --course-id is missing")
	}
	combined := stdout + stderr
	if !strings.Contains(combined, "course-id") {
		t.Errorf("expected 'course-id' mentioned in output; got stdout=%q stderr=%q", stdout, stderr)
	}
}

// Case 8: API 401 → exit 1, output suggests auth problem.
func TestCoursesListUnauthorized(t *testing.T) {
	srv := mockServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"errors":[{"message":"Invalid access token"}]}`))
	}))
	home := t.TempDir()

	stdout, stderr, code := runCanvas(t, authEnv(home, srv.URL), "courses", "list")

	if code == 0 {
		t.Fatal("expected non-zero exit on 401 from server")
	}
	combined := strings.ToLower(stdout + stderr)
	// The CLI should surface something indicating auth/token failure.
	if !strings.Contains(combined, "401") &&
		!strings.Contains(combined, "unauthorized") &&
		!strings.Contains(combined, "token") &&
		!strings.Contains(combined, "auth") {
		t.Errorf("expected auth-related text in output; stdout=%q stderr=%q", stdout, stderr)
	}
}

// Case 9: alias set then invoke → alias expansion works end-to-end.
func TestAliasSetAndExpand(t *testing.T) {
	srv := mockServer(t, newCoursesHandler())
	home := t.TempDir()
	env := authEnv(home, srv.URL)

	// Step 1: set alias "ll" → "courses list"
	_, stderr, code := runCanvas(t, env, "alias", "set", "ll", "courses list")
	if code != 0 {
		t.Fatalf("alias set failed (exit %d): %s", code, stderr)
	}

	// Step 2: invoke "ll" — the main() expandAliases() should substitute "courses list"
	stdout, stderr, code := runCanvas(t, env, "ll")
	if code != 0 {
		t.Fatalf("alias invocation failed (exit %d): stderr=%q stdout=%q", code, stderr, stdout)
	}
	if !strings.Contains(stdout, "Introduction to Go") {
		t.Errorf("expected alias-expanded output to contain course name; got: %q", stdout)
	}
}

// Case 10: context set course → assignments list without --course-id uses context.
// The mock server asserts that the course path was actually hit.
func TestContextCourseUsedByAssignmentsList(t *testing.T) {
	const courseID = "789"
	var coursePathHit atomic.Bool

	srv := mockServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/accounts":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[]`))
		case r.URL.Path == "/api/v1/courses/"+courseID:
			// Course validation
			coursePathHit.Store(true)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"id":789,"name":"Test Course"}`))
		case r.URL.Path == "/api/v1/courses/"+courseID+"/assignments":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[{"id":1,"name":"Homework 1","course_id":789}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	home := t.TempDir()
	env := authEnv(home, srv.URL)

	// Set context
	_, stderr, code := runCanvas(t, env, "context", "set", "course", courseID)
	if code != 0 {
		t.Fatalf("context set failed (exit %d): %s", code, stderr)
	}

	// List assignments without --course-id; context should supply it
	stdout, stderr, code := runCanvas(t, env, "assignments", "list")
	if code != 0 {
		t.Fatalf("assignments list failed (exit %d): stderr=%q stdout=%q", code, stderr, stdout)
	}
	if !coursePathHit.Load() {
		t.Error("expected the mock server to receive a request for course 789 (course validation)")
	}
	if !strings.Contains(stdout, "Homework 1") {
		t.Errorf("expected assignment name in output; got: %q", stdout)
	}
}

// Case 11: --dry-run prints curl preview with token REDACTED; mock receives NO request.
func TestDryRunRedactsToken(t *testing.T) {
	var serverHit atomic.Bool

	srv := mockServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverHit.Store(true)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	}))
	home := t.TempDir()
	env := authEnv(home, srv.URL)

	stdout, stderr, code := runCanvas(t, env, "courses", "list", "--dry-run")

	if code != 0 {
		t.Fatalf("expected exit 0 for --dry-run, got %d; stderr=%q", code, stderr)
	}
	if serverHit.Load() {
		t.Error("--dry-run should not send any request to the mock server")
	}
	if !strings.Contains(stdout, "REDACTED") {
		t.Errorf("expected '[REDACTED]' token in dry-run output; got: %q", stdout)
	}
	if strings.Contains(stdout, "test-token") {
		t.Errorf("real token must not appear in dry-run output; got: %q", stdout)
	}
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// newCoursesHandler returns an http.Handler that serves a single hardcoded
// course on GET /api/v1/courses plus a stub for /api/v1/accounts (used by the
// client for Canvas version detection).
func newCoursesHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	})

	mux.HandleFunc("/api/v1/courses", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":42,"name":"Introduction to Go","course_code":"CS101","workflow_state":"available"}]`))
	})

	return mux
}
