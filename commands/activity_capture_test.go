package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spf13/cobra"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
	"github.com/jjuanrivvera/canvas-cli/internal/activity"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

func captureStderr(fn func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	fn()
	w.Close()
	os.Stderr = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// newTestClient returns a client against server with the activity observer
// installed on it, exactly as runWithActivityLog does.
func newTestClient(t *testing.T, server *httptest.Server, rec *activity.Recorder) *api.Client {
	t.Helper()
	client, err := api.NewClient(api.ClientConfig{BaseURL: server.URL, Token: "test-token", RequestsPerSec: 100, UserAgent: "test", RetryInitialBackoff: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	api.RequestObserver = func(o api.ObservedRequest) {
		rec.Observe(activity.Observation{Method: o.Method, Path: o.Path, Status: o.Status, DryRun: o.DryRun, RequestBody: o.RequestBody, ResponseBody: o.ResponseBody})
	}
	t.Cleanup(func() { api.RequestObserver = nil; api.RequestGate = nil })
	return client
}

func readEntries(t *testing.T, path string) []activity.Entry {
	t.Helper()
	entries, skipped, err := activity.Read(path)
	if err != nil || skipped != 0 {
		t.Fatalf("Read = %d skipped, %v", skipped, err)
	}
	return entries
}

// ---- writes-only default and capture ----

func TestLogActivity_WritesOnlyByDefault(t *testing.T) {
	useTempHome(t)
	path := useTempActivityLog(t)
	cmd := fakeExecuted(t, "10", "--course-id", "1")

	// reads only: no entry
	rec := activity.NewRecorder()
	rec.Record("GET", "/api/v1/courses/1", 200)
	logActivity(cmd, nil, time.Now(), rec, []string{"quizzes", "regrade", "10"})
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("a read-only invocation must leave no entry by default")
	}

	// a dry-run write: still no entry
	rec = activity.NewRecorder()
	rec.Observe(activity.Observation{Method: "PUT", Path: "/api/v1/x", DryRun: true})
	logActivity(cmd, nil, time.Now(), rec, []string{"quizzes", "regrade", "10"})
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("a dry-run invocation must leave no entry by default")
	}

	// writes_only turned off by the environment: reads are logged
	t.Setenv(activity.EnvWritesOnly, "0")
	rec = activity.NewRecorder()
	rec.Record("GET", "/api/v1/courses/1", 200)
	logActivity(cmd, nil, time.Now(), rec, []string{"quizzes", "regrade", "10"})
	if got := readEntries(t, path); len(got) != 1 {
		t.Fatalf("with CANVAS_ACTIVITY_WRITES_ONLY=0 the read must be logged, got %d entries", len(got))
	}
}

func TestLogActivity_CaptureBodies(t *testing.T) {
	useTempHome(t)
	path := useTempActivityLog(t)
	cmd := fakeExecuted(t, "10", "--course-id", "1")
	argv := []string{"quizzes", "regrade", "10", "--course-id", "1"}
	observe := func(rec *activity.Recorder) {
		rec.Observe(activity.Observation{Method: "GET", Path: "/api/v1/courses/1", Status: 200, ResponseBody: []byte(`{"id":1,"name":"read"}`)})
		rec.Observe(activity.Observation{Method: "PUT", Path: "/api/v1/courses/1/x", Status: 200,
			RequestBody:  []byte(`{"comment":{"text_comment":"Great job"},"access_code":"1234","token":"7~AbCdEfGhIjKlMnOpQrStUv"}`),
			ResponseBody: []byte(`{"id":7,"api_key":"k","note":"Bearer zzz"}`)})
	}

	// off: no bodies, but outcomes
	rec := activity.NewRecorder()
	observe(rec)
	logActivity(cmd, nil, time.Now(), rec, argv)
	raw, _ := os.ReadFile(path)
	if strings.Contains(string(raw), "Great job") || strings.Contains(string(raw), `"body"`) || strings.Contains(string(raw), `"response"`) {
		t.Errorf("bodies logged although capture is off: %s", raw)
	}
	if !strings.Contains(string(raw), `"outcome":"ok"`) {
		t.Errorf("outcome missing: %s", raw)
	}

	// on: bodies for the write, nothing for the read, secrets redacted
	t.Setenv(activity.EnvCaptureBodies, "1")
	rec = activity.NewRecorder()
	observe(rec)
	logActivity(cmd, nil, time.Now(), rec, argv)
	entries := readEntries(t, path)
	if len(entries) != 2 {
		t.Fatalf("entries = %d", len(entries))
	}
	e := entries[1]
	if e.Requests[0].Body != nil || e.Requests[0].Response != nil {
		t.Errorf("GET must carry no body/response: %+v", e.Requests[0])
	}
	line, _ := json.Marshal(e.Requests[1])
	for _, want := range []string{`"text_comment":"Great job"`, `"access_code":"[REDACTED]"`, `"token":"[REDACTED]"`, `"api_key":"[REDACTED]"`, `"note":"Bearer [REDACTED]"`, `"id":7`} {
		if !strings.Contains(string(line), want) {
			t.Errorf("PUT entry lacks %s:\n%s", want, line)
		}
	}
	raw, _ = os.ReadFile(path)
	for _, leak := range []string{"7~AbCd", `"1234"`, `"k"`, "zzz", `"name":"read"`} {
		if strings.Contains(string(raw), leak) {
			t.Errorf("leaked %s into the log: %s", leak, raw)
		}
	}
}

func TestTouchedFromResponses(t *testing.T) {
	cases := []struct {
		command string
		method  string
		body    string
		resp    string
		want    string
	}{
		{"submissions grade", "PUT", `{"submission":{"posted_grade":"9"},"comment":{"text_comment":"new one"}}`,
			`{"id":42,"submission_comments":[{"id":5,"comment":"old"},{"id":6,"comment":"new one"}]}`, "submission:42 submission-comment:6"},
		{"submissions add-comment", "PUT", `{"comment":{"text_comment":"hi"}}`, `{"id":42,"submission_comments":[{"id":9,"comment":"hi"}]}`, "submission:42 submission-comment:9"},
		{"conversations create", "POST", `{"body":"x"}`, `[{"id":100,"messages":[{"id":1000}]},{"id":101,"messages":[{"id":1001}]}]`, "conversation:100 conversation:101 conversation-message:1000 conversation-message:1001"},
		{"conversations reply", "POST", `{"body":"x"}`, `{"id":100,"messages":[{"id":1002}]}`, "conversation:100 conversation-message:1002"},
		{"announcements create", "POST", `{}`, `{"id":55}`, "announcement:55"},
		{"announcements update", "PUT", `{}`, `{"id":55}`, "announcement:55"},
		{"discussions post", "POST", `{}`, `{"id":77}`, "discussion-entry:77"},
		{"discussions reply", "POST", `{}`, `{"id":78}`, "discussion-entry:78"},
		{"pages create", "POST", `{}`, `{"page_id":31,"url":"intro"}`, "page:31"},
		{"pages update", "PUT", `{}`, `{"page_id":31,"url":"intro"}`, "page:31"},
		{"courses update", "PUT", `{}`, `{"id":1}`, ""}, // not enriched
	}
	for _, tc := range cases {
		t.Run(tc.command, func(t *testing.T) {
			rec := activity.NewRecorder()
			rec.Observe(activity.Observation{Method: "GET", Path: "/r", Status: 200, ResponseBody: []byte(`{"id":999}`)})
			rec.Observe(activity.Observation{Method: tc.method, Path: "/w", Status: 200, RequestBody: []byte(tc.body), ResponseBody: []byte(tc.resp)})
			// a rejected write contributes nothing
			rec.Observe(activity.Observation{Method: tc.method, Path: "/w", Status: 422, RequestBody: []byte(tc.body), ResponseBody: []byte(`{"id":888,"page_id":888}`)})
			touchedFromResponses(tc.command, rec)
			var got []string
			for _, x := range rec.Touched() {
				got = append(got, x.Type+":"+strconv.FormatInt(x.ID, 10))
			}
			if strings.Join(got, " ") != tc.want {
				t.Errorf("touched = %q, want %q", strings.Join(got, " "), tc.want)
			}
		})
	}
}

// ---- unknown outcome ----

func TestLogActivity_UnknownOutcomeIsAlwaysLogged(t *testing.T) {
	useTempHome(t)
	path := useTempActivityLog(t)

	// the server hangs up before answering every PUT: the request was
	// received (and may have been applied) but no status ever arrives
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			hj, ok := w.(http.Hijacker)
			if !ok {
				t.Fatal("no hijacker")
			}
			conn, _, _ := hj.Hijack()
			_ = conn.Close()
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":1}`))
	}))
	defer server.Close()

	rec := activity.NewRecorder()
	client := newTestClient(t, server, rec)
	ctx := context.Background()
	var out map[string]interface{}
	if err := client.GetJSON(ctx, "/api/v1/courses/1", &out); err != nil {
		t.Fatalf("GET: %v", err)
	}
	err := client.PutJSON(ctx, "/api/v1/courses/1/assignments/2/submissions/3", map[string]interface{}{"submission": map[string]string{"posted_grade": "9"}}, &out)
	if err == nil {
		t.Fatal("the PUT must fail: the server hung up")
	}

	cmd := fakeExecuted(t, "10", "--course-id", "1")
	logActivity(cmd, err, time.Now(), rec, []string{"quizzes", "regrade", "10"})

	entries := readEntries(t, path)
	if len(entries) != 1 {
		t.Fatalf("an attempted write with a lost response must be logged under writes_only, got %d entries", len(entries))
	}
	e := entries[0]
	if !e.VerificationRequired || e.ExitCode != 1 {
		t.Errorf("entry = %+v", e)
	}
	var put *activity.Request
	for i := range e.Requests {
		if e.Requests[i].Method == "PUT" {
			put = &e.Requests[i]
		}
	}
	if put == nil || put.Status != 0 || put.Outcome != activity.OutcomeUnknown {
		t.Errorf("PUT request = %+v, want status 0 outcome unknown", put)
	}
	raw, _ := os.ReadFile(path)
	if !strings.Contains(string(raw), `"outcome":"unknown"`) || !strings.Contains(string(raw), `"verification_required":true`) {
		t.Errorf("line = %s", raw)
	}
	// the list view counts it
	if row := activityRow(e); row.Unknown != 1 || row.Writes != 0 {
		t.Errorf("row = %+v", row)
	}
}

// ---- audited mode ----

func TestActivityGate_Required(t *testing.T) {
	useTempHome(t)
	good := filepath.Join(t.TempDir(), "audit", "activity.jsonl")
	blocker := filepath.Join(t.TempDir(), "file")
	_ = os.WriteFile(blocker, []byte("x"), 0o600)
	bad := filepath.Join(blocker, "activity.jsonl")

	// not required: an unusable log never blocks
	gate := newActivityGate(func() activity.Config { return activity.Config{Enabled: true, Path: bad} })
	if err := gate("PUT", "/x"); err != nil {
		t.Errorf("without required the gate must pass: %v", err)
	}

	// required + unusable: refused, with the reason
	gate = newActivityGate(func() activity.Config { return activity.Config{Enabled: true, Required: true, Path: bad} })
	err := gate("PUT", "/x")
	if err == nil || !strings.Contains(err.Error(), "audit log unavailable") || !strings.Contains(err.Error(), "refusing to write to Canvas") {
		t.Errorf("gate error = %v", err)
	}
	if err2 := gate("POST", "/y"); err2 == nil || err2.Error() != err.Error() {
		t.Errorf("the verdict must be stable for the process: %v", err2)
	}

	// required + usable: passes and prepares dir 0700 / file 0600
	gate = newActivityGate(func() activity.Config { return activity.Config{Enabled: true, Required: true, Path: good} })
	if err := gate("PUT", "/x"); err != nil {
		t.Fatalf("usable log refused: %v", err)
	}
	if info, err := os.Stat(good); err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("preflight must create the log 0600: %v %v", info, err)
	}
	if info, _ := os.Stat(filepath.Dir(good)); info.Mode().Perm() != 0o700 {
		t.Errorf("preflight must create the dir 0700, got %o", info.Mode().Perm())
	}
}

func TestActivityGate_RefusesWriteBeforeCanvas(t *testing.T) {
	useTempHome(t)
	var puts, gets int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case !strings.Contains(r.URL.Path, "/courses/1"): // version probe etc.
		case r.Method == http.MethodGet:
			atomic.AddInt32(&gets, 1)
		default:
			atomic.AddInt32(&puts, 1)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":1}`))
	}))
	defer server.Close()

	blocker := filepath.Join(t.TempDir(), "file")
	_ = os.WriteFile(blocker, []byte("x"), 0o600)
	t.Setenv(activity.EnvVar, filepath.Join(blocker, "activity.jsonl"))
	t.Setenv(activity.EnvRequired, "1")

	rec := activity.NewRecorder()
	client := newTestClient(t, server, rec)
	api.RequestGate = newActivityGate(resolveActivityConfig)

	ctx := context.Background()
	var out map[string]interface{}
	if err := client.GetJSON(ctx, "/api/v1/courses/1", &out); err != nil {
		t.Fatalf("reads are never gated: %v", err)
	}
	err := client.PutJSON(ctx, "/api/v1/courses/1/assignments/2/submissions/3", map[string]string{"x": "y"}, &out)
	if err == nil || !strings.Contains(err.Error(), "refusing to write to Canvas") {
		t.Fatalf("PUT must be refused, got %v", err)
	}
	if atomic.LoadInt32(&puts) != 0 || atomic.LoadInt32(&gets) != 1 {
		t.Errorf("Canvas saw %d writes and %d reads; want 0 and 1", puts, gets)
	}
	for _, r := range rec.Requests() {
		if !r.IsRead() {
			t.Errorf("a refused write must not be recorded as sent: %+v", r)
		}
	}

	// without required, the same unusable log only warns at exit
	t.Setenv(activity.EnvRequired, "")
	api.RequestGate = newActivityGate(resolveActivityConfig)
	if err := client.PutJSON(ctx, "/api/v1/courses/1/assignments/2/submissions/3", map[string]string{"x": "y"}, &out); err != nil {
		t.Fatalf("without required the write must go through: %v", err)
	}
	if atomic.LoadInt32(&puts) != 1 {
		t.Errorf("Canvas saw %d writes, want 1", puts)
	}
	stderr := captureStderr(func() {
		logActivity(fakeExecuted(t, "10", "--course-id", "1"), nil, time.Now(), rec, []string{"quizzes", "regrade", "10"})
	})
	if !strings.Contains(stderr, "activity log not written") {
		t.Errorf("stderr = %q", stderr)
	}
}

// ---- permissions ----

func TestLogActivity_TightensLoosePermissions(t *testing.T) {
	useTempHome(t)
	path := useTempActivityLog(t)
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	_ = os.Chmod(path, 0o644)

	rec := activity.NewRecorder()
	rec.Record("PUT", "/api/v1/x", 200)
	cmd := fakeExecuted(t, "10", "--course-id", "1")
	stderr := captureStderr(func() { logActivity(cmd, nil, time.Now(), rec, []string{"quizzes", "regrade", "10"}) })
	if info, _ := os.Stat(path); info.Mode().Perm() != 0o600 {
		t.Errorf("log perm = %o, want 600", info.Mode().Perm())
	}
	if !strings.Contains(stderr, "tightened permissions on "+path+" from 0644 to 0600") {
		t.Errorf("stderr = %q", stderr)
	}
	if len(readEntries(t, path)) != 1 {
		t.Error("entry not written")
	}
	// second run: nothing to say
	stderr = captureStderr(func() { logActivity(cmd, nil, time.Now(), rec, []string{"quizzes", "regrade", "10"}) })
	if stderr != "" {
		t.Errorf("second run stderr = %q", stderr)
	}
}

// ---- bulk-grade capture ----

func TestBulkGrade_CapturesEveryWriteBody(t *testing.T) {
	useTempHome(t)
	path := useTempActivityLog(t)
	t.Setenv(activity.EnvCaptureBodies, "true")

	csvPath := filepath.Join(t.TempDir(), "grades.csv")
	// every row asks for the grade the static mock reports back, so the
	// test also holds once the command verifies its writes by read-back
	_ = os.WriteFile(csvPath, []byte("user_id,assignment_id,grade,comment\n10,100,95,Great work\n11,100,95,\n12,100,95,See rubric token 7~AbCdEfGhIjKlMnOpQrStUv\n"), 0o600)

	// the command attaches details.input to the process-wide recorder
	rec := activity.Default()
	rec.Reset()
	api.RequestObserver = func(o api.ObservedRequest) {
		rec.Observe(activity.Observation{Method: o.Method, Path: o.Path, Status: o.Status, DryRun: o.DryRun, RequestBody: o.RequestBody, ResponseBody: o.ResponseBody})
	}
	t.Cleanup(func() { api.RequestObserver = nil; rec.Reset() })

	leaf := newSubmissionsBulkGradeCmd()
	sub := &cobra.Command{Use: "submissions"}
	sub.AddCommand(leaf)
	root := &cobra.Command{Use: "canvas"}
	root.AddCommand(sub)
	argv := []string{"submissions", "bulk-grade", "--course-id", "1", "--csv-file", csvPath}
	cmdtest.RunCommandTest(t, root, cmdtest.CommandTestCase{
		Name: "bulk grade",
		Args: argv,
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments/100": cmdtest.NewMockResponse(`{"id":1,"assignment_id":100,"user_id":10,"score":95,"grade":"95",
				"submission_comments":[{"id":5,"author_name":"Teacher","comment":"Great work"},{"id":6,"author_name":"Teacher","comment":"See rubric token 7~AbCdEfGhIjKlMnOpQrStUv"}]}`),
		},
	})
	logActivity(leaf, nil, time.Now(), rec, argv)

	entries := readEntries(t, path)
	if len(entries) != 1 {
		t.Fatalf("entries = %d", len(entries))
	}
	e := entries[0]
	var bodies []string
	for _, r := range e.Requests {
		if r.Method == "PUT" {
			b, _ := json.Marshal(r.Body)
			bodies = append(bodies, string(b))
			if r.Outcome != activity.OutcomeOK || r.Response == nil {
				t.Errorf("PUT %s: outcome %q response %v", r.Path, r.Outcome, r.Response)
			}
		}
	}
	if len(bodies) != 3 {
		t.Fatalf("want 3 logged write bodies, got %d: %v", len(bodies), bodies)
	}
	for i := range bodies {
		if !strings.Contains(bodies[i], `"posted_grade":"95"`) {
			t.Errorf("body %d = %s, want the row's grade", i, bodies[i])
		}
	}
	if !strings.Contains(bodies[0], `"text_comment":"Great work"`) || strings.Contains(bodies[1], "text_comment") {
		t.Errorf("comments: %v", bodies)
	}
	if !strings.Contains(bodies[2], "[REDACTED]") || strings.Contains(bodies[2], "7~AbCd") {
		t.Errorf("token in a comment must be redacted: %s", bodies[2])
	}

	// details.input holds the parsed CSV rows, redacted
	input, _ := e.Details["input"].([]interface{})
	if len(input) != 3 {
		t.Fatalf("details.input = %v", e.Details["input"])
	}
	row0, _ := input[0].(map[string]interface{})
	if row0["grade"] != "95" || row0["comment"] != "Great work" || row0["user_id"] != float64(10) || row0["assignment_id"] != float64(100) {
		t.Errorf("row 0 = %v", row0)
	}
	raw, _ := os.ReadFile(path)
	if strings.Contains(string(raw), "7~AbCd") {
		t.Errorf("token leaked: %s", raw)
	}
	if got := readEntries(t, path)[0].Touched; len(got) == 0 || got[0].Type != "course" {
		t.Errorf("touched = %v", got)
	}
}

// ---- configure ----

func TestActivityConfigureCmd_RoundTrip(t *testing.T) {
	home := useTempHome(t)
	logPath := filepath.Join(home, "logs", "activity.jsonl")

	cmdtest.RunCommandTest(t, newActivityConfigureCmd(), cmdtest.CommandTestCase{
		Name: "enable",
		Args: []string{"--enable", "--capture-bodies", "--required", "--path", logPath, "--max-size-mb", "5"},
		ValidateOutput: func(t *testing.T, output string) {
			for _, want := range []string{"Activity log settings saved to " + filepath.Join(home, ".canvas-cli", "config.yaml"), logPath, "enabled: true", "writes_only: true", "capture_bodies: true", "required: true", "max_size_mb: 5"} {
				if !strings.Contains(output, want) {
					t.Errorf("output lacks %q:\n%s", want, output)
				}
			}
		},
	})
	cfg, err := config.Reload()
	if err != nil || cfg.ActivityLog == nil {
		t.Fatalf("config after configure: %+v, %v", cfg, err)
	}
	al := cfg.ActivityLog
	if !al.Enabled || !al.CaptureBodies || !al.Required || al.Path != logPath || al.MaxSizeMB != 5 || al.WritesOnly != nil {
		t.Errorf("activity_log = %+v", al)
	}
	if info, _ := os.Stat(filepath.Join(home, ".canvas-cli", "config.yaml")); info == nil || info.Mode().Perm() != 0o600 {
		t.Errorf("config file perm = %v", info)
	}

	// change some keys, keep the rest
	cmdtest.RunCommandTest(t, newActivityConfigureCmd(), cmdtest.CommandTestCase{
		Name: "update",
		Args: []string{"--writes-only=false", "--disable"},
		ValidateOutput: func(t *testing.T, output string) {
			// required still implies enabled at resolution time; the file says disabled
			if !strings.Contains(output, "writes_only: false") || !strings.Contains(output, "capture_bodies: true") {
				t.Errorf("output = %s", output)
			}
		},
	})
	cfg, _ = config.Reload()
	al = cfg.ActivityLog
	if al.Enabled || al.WritesOnly == nil || *al.WritesOnly || !al.CaptureBodies || !al.Required || al.Path != logPath || al.MaxSizeMB != 5 {
		t.Errorf("activity_log after update = %+v", al)
	}
	raw, _ := os.ReadFile(filepath.Join(home, ".canvas-cli", "config.yaml"))
	if !strings.Contains(string(raw), "writes_only: false") || !strings.Contains(string(raw), "activity_log:") {
		t.Errorf("config file = %s", raw)
	}

	// an explicit false survives another Save from elsewhere in the CLI
	// (the alias/context/instance commands all go through Config.Save)
	other, _ := config.Load()
	_ = other.SetAlias("ll", "courses list")
	if err := other.Save(); err != nil {
		t.Fatal(err)
	}
	cfg, _ = config.Reload()
	if cfg.ActivityLog == nil || cfg.ActivityLog.WritesOnly == nil || *cfg.ActivityLog.WritesOnly || cfg.ActivityLog.Path != logPath {
		t.Errorf("activity_log lost by another Save: %+v", cfg.ActivityLog)
	}

	// path reports the effective settings, json included
	cmdtest.RunCommandTest(t, newActivityPathCmd(), cmdtest.CommandTestCase{Name: "path", ValidateOutput: func(t *testing.T, output string) {
		if !strings.HasPrefix(output, logPath+"\n") || !strings.Contains(output, "required: true") || !strings.Contains(output, "capture_bodies: true") {
			t.Errorf("path output = %s", output)
		}
	}})
	setOutputFormat(t, "json")
	cmdtest.RunCommandTest(t, newActivityPathCmd(), cmdtest.CommandTestCase{Name: "path json", ValidateOutput: func(t *testing.T, output string) {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(output), &m); err != nil {
			t.Fatalf("not json: %v\n%s", err, output)
		}
		if m["path"] != logPath || m["required"] != true || m["writes_only"] != false || m["capture_bodies"] != true || m["max_size_mb"] != float64(5) {
			t.Errorf("json = %v", m)
		}
	}})
	setOutputFormat(t, "table")

	// mutually exclusive switches and a bad size are rejected
	cmdtest.RunCommandTest(t, newActivityConfigureCmd(), cmdtest.CommandTestCase{Name: "enable+disable", Args: []string{"--enable", "--disable"}, ExpectError: true})
	cmdtest.RunCommandTest(t, newActivityConfigureCmd(), cmdtest.CommandTestCase{Name: "negative size", Args: []string{"--max-size-mb", "-1"}, ExpectError: true})
}

func TestActivityConfigure_ThenLogs(t *testing.T) {
	home := useTempHome(t)
	cmdtest.RunCommandTest(t, newActivityConfigureCmd(), cmdtest.CommandTestCase{Name: "enable", Args: []string{"--enable"}})
	rec := activity.NewRecorder()
	rec.Record("PUT", "/api/v1/x", 200)
	logActivity(fakeExecuted(t, "10", "--course-id", "1"), nil, time.Now(), rec, []string{"quizzes", "regrade", "10"})
	got := readEntries(t, filepath.Join(home, ".canvas-cli", activity.DefaultFileName))
	if len(got) != 1 {
		t.Fatalf("after configure --enable the default log must receive the entry, got %d", len(got))
	}
}

func TestLogActivity_EnrichesTouchedFromResponse(t *testing.T) {
	useTempHome(t)
	path := useTempActivityLog(t)

	root := &cobra.Command{Use: "canvas"}
	sub := &cobra.Command{Use: "submissions"}
	grade := &cobra.Command{Use: "grade", RunE: func(*cobra.Command, []string) error { return nil }}
	grade.Flags().Int64("course-id", 0, "")
	grade.Flags().Int64("assignment-id", 0, "")
	grade.Flags().Int64("user-id", 0, "")
	grade.Flags().String("comment", "", "")
	sub.AddCommand(grade)
	root.AddCommand(sub)
	argv := []string{"submissions", "grade", "--course-id", "1", "--assignment-id", "100", "--user-id", "10", "--comment", "Well done"}
	root.SetArgs(argv)
	cmd, err := root.ExecuteContextC(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	rec := activity.NewRecorder()
	rec.Observe(activity.Observation{Method: "GET", Path: "/api/v1/courses/1", Status: 200, ResponseBody: []byte(`{"id":1}`)})
	rec.Observe(activity.Observation{Method: "PUT", Path: "/api/v1/courses/1/assignments/100/submissions/10", Status: 200,
		RequestBody:  []byte(`{"submission":{"posted_grade":"9"},"comment":{"text_comment":"Well done"}}`),
		ResponseBody: []byte(`{"id":4242,"submission_comments":[{"id":77,"comment":"Well done"}]}`)})
	logActivity(cmd, nil, time.Now(), rec, argv)

	entries := readEntries(t, path)
	if len(entries) != 1 {
		t.Fatalf("entries = %d", len(entries))
	}
	var got []string
	for _, x := range entries[0].Touched {
		got = append(got, x.Type+":"+strconv.FormatInt(x.ID, 10))
	}
	if want := "assignment:100 course:1 submission:4242 submission-comment:77 user:10"; strings.Join(got, " ") != want {
		t.Errorf("touched = %q, want %q", strings.Join(got, " "), want)
	}
}
