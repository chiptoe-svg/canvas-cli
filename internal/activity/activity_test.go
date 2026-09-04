package activity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRecorder(t *testing.T) {
	r := NewRecorder()
	r.Record("GET", "/api/v1/courses/1?access_token=7~abc&page=2", 200)
	r.Record("PUT", "/api/v1/courses/1/quizzes/2/questions/3", 200)
	r.Record("PUT", "/api/v1/courses/1/quizzes/2/submissions/4", 0) // dry run / no response
	r.Touch("course", 1)
	r.Touch("course", 1) // duplicate
	r.Touch("quiz", 2)
	r.Touch("", 9)     // ignored
	r.Touch("user", 0) // ignored
	r.SetDetail("regrade", map[string]int{"changed": 2})

	reqs := r.Requests()
	if len(reqs) != 3 {
		t.Fatalf("requests = %v", reqs)
	}
	if reqs[0].Path != "/api/v1/courses/1" || strings.Contains(reqs[0].Path, "?") {
		t.Errorf("query string must be dropped, got %q", reqs[0].Path)
	}
	if got := r.Touched(); len(got) != 2 || got[0] != (Touched{"course", 1}) || got[1] != (Touched{"quiz", 2}) {
		t.Errorf("touched = %v", got)
	}
	if d := r.Details(); d["regrade"] == nil {
		t.Errorf("details = %v", d)
	}

	e := Entry{Requests: reqs}
	if !e.HasWrites() || e.WriteCount() != 1 {
		t.Errorf("HasWrites/WriteCount: only the PUT with a status counts, got %v / %d", e.HasWrites(), e.WriteCount())
	}
	if (Entry{Requests: []Request{{Method: "GET", Path: "/x", Status: 200}}}).HasWrites() {
		t.Error("a GET-only entry has no writes")
	}

	r.Reset()
	if len(r.Requests()) != 0 || len(r.Touched()) != 0 || r.Details() != nil {
		t.Error("Reset did not clear the recorder")
	}
}

func TestRedactArgs(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"token flag with space", []string{"auth", "token", "set", "x", "--token", "7~AbCdEfGhIjKlMnOpQrStUv"}, []string{"auth", "token", "set", "x", "--token", Redacted}},
		{"token flag with equals", []string{"--token=7~AbCdEfGhIjKlMnOpQrStUv"}, []string{"--token=" + Redacted}},
		{"client secret", []string{"--client-secret", "hunter2", "--url", "https://x"}, []string{"--client-secret", Redacted, "--url", "https://x"}},
		{"password and access code", []string{"--password=pw", "--access-code", "1234"}, []string{"--password=" + Redacted, "--access-code", Redacted}},
		{"bare canvas token value", []string{"api", "GET", "/api/v1/users?access_token=7~AbCdEfGhIjKlMnOpQrStUv"}, []string{"api", "GET", "/api/v1/users?access_token=" + Redacted}},
		{"authorization header", []string{"api", "GET", "/x", "-H", "Authorization: Bearer abc.def"}, []string{"api", "GET", "/x", "-H", "Authorization: Bearer " + Redacted}},
		{"show-token flag is a switch, not a value", []string{"--dry-run", "--show-token", "courses", "list"}, []string{"--dry-run", "--show-token", Redacted, "list"}},
		{"ordinary ids untouched", []string{"quizzes", "regrade", "456", "--course-id", "123", "--question", "789"}, []string{"quizzes", "regrade", "456", "--course-id", "123", "--question", "789"}},
		// student-directed text and student names are not for the log; ids and scores are
		{"comment text redacted, score kept", []string{"submissions", "grade", "--user-id", "10", "--score", "95", "--comment", "Great work, Ada"}, []string{"submissions", "grade", "--user-id", "10", "--score", "95", "--comment", Redacted}},
		{"student name in equals form", []string{"submissions", "excuse", "--student=Ada Lovelace", "--assignment", "456"}, []string{"submissions", "excuse", "--student=" + Redacted, "--assignment", "456"}},
		{"rubric comment and message", []string{"--rubric-comment", "_1=vague", "--message", "see me"}, []string{"--rubric-comment", Redacted, "--message", Redacted}},
		{"empty", nil, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RedactArgs(tt.in)
			if strings.Join(got, "\x00") != strings.Join(tt.want, "\x00") {
				t.Errorf("RedactArgs(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}

	// nothing secret-shaped survives in any position
	for _, a := range RedactArgs([]string{"--token", "7~AbCdEfGhIjKlMnOpQrStUv", "7~AbCdEfGhIjKlMnOpQrStUv", "Bearer 7~AbCdEfGhIjKlMnOpQrStUv"}) {
		if strings.Contains(a, "7~") || strings.Contains(a, "AbCdEf") {
			t.Errorf("token leaked: %q", a)
		}
	}
}

func boolPtr(b bool) *bool { return &b }

func TestResolve(t *testing.T) {
	cfg := Resolve("/cfg", Settings{}, Env{})
	if cfg.Enabled || cfg.Path != filepath.Join("/cfg", DefaultFileName) || cfg.MaxSizeBytes != DefaultMaxSizeMB*1024*1024 || cfg.CaptureBodies || cfg.Required {
		t.Errorf("defaults = %+v", cfg)
	}
	if !cfg.WritesOnly {
		t.Error("writes_only must default to true")
	}
	cfg = Resolve("/cfg", Settings{Enabled: true, Path: "/custom.jsonl", WritesOnly: boolPtr(false), CaptureBodies: true, MaxSizeMB: 2}, Env{})
	if !cfg.Enabled || cfg.Path != "/custom.jsonl" || cfg.WritesOnly || !cfg.CaptureBodies || cfg.MaxSizeBytes != 2*1024*1024 {
		t.Errorf("config file = %+v", cfg)
	}
	// env wins and enables on its own
	cfg = Resolve("/cfg", Settings{Enabled: false, Path: "/custom.jsonl"}, Env{Path: "/env.jsonl"})
	if !cfg.Enabled || cfg.Path != "/env.jsonl" || !cfg.WritesOnly {
		t.Errorf("env = %+v", cfg)
	}
	// boolean env overrides beat the config file in both directions
	cfg = Resolve("/cfg", Settings{WritesOnly: boolPtr(true), CaptureBodies: false}, Env{WritesOnly: "0", CaptureBodies: "1"})
	if cfg.WritesOnly || !cfg.CaptureBodies {
		t.Errorf("env bools = %+v", cfg)
	}
	cfg = Resolve("/cfg", Settings{WritesOnly: boolPtr(false), CaptureBodies: true}, Env{WritesOnly: "true", CaptureBodies: "false"})
	if !cfg.WritesOnly || cfg.CaptureBodies {
		t.Errorf("env bools (reverse) = %+v", cfg)
	}
	// an unparsable env value changes nothing
	cfg = Resolve("/cfg", Settings{WritesOnly: boolPtr(false)}, Env{WritesOnly: "maybe"})
	if cfg.WritesOnly {
		t.Errorf("garbage env must be ignored: %+v", cfg)
	}
	// required implies enabled, from the file or the environment
	if cfg := Resolve("/cfg", Settings{Required: true}, Env{}); !cfg.Required || !cfg.Enabled {
		t.Errorf("required (file) = %+v", cfg)
	}
	if cfg := Resolve("/cfg", Settings{}, Env{Required: "1"}); !cfg.Required || !cfg.Enabled {
		t.Errorf("required (env) = %+v", cfg)
	}
	if cfg := Resolve("/cfg", Settings{Required: true}, Env{Required: "0"}); cfg.Required || cfg.Enabled {
		t.Errorf("required overridden off = %+v", cfg)
	}
}

func TestParseBoolEnv(t *testing.T) {
	for in, want := range map[string]bool{"1": true, "true": true, "TRUE": true, " yes ": true, "on": true, "0": false, "false": false, "no": false, "off": false} {
		got, set := ParseBoolEnv(in)
		if !set || got != want {
			t.Errorf("ParseBoolEnv(%q) = %v,%v; want %v,true", in, got, set, want)
		}
	}
	for _, in := range []string{"", "maybe", "2"} {
		if _, set := ParseBoolEnv(in); set {
			t.Errorf("ParseBoolEnv(%q) must not be set", in)
		}
	}
}

func TestObserveOutcomes(t *testing.T) {
	r := NewRecorder()
	r.Observe(Observation{Method: "GET", Path: "/a", Status: 200, ResponseBody: []byte(`{"x":1}`)})
	r.Observe(Observation{Method: "PUT", Path: "/b", Status: 200, RequestBody: []byte(`{"submission":{"posted_grade":"9"}}`), ResponseBody: []byte(`{"id":1}`)})
	r.Observe(Observation{Method: "POST", Path: "/c", Status: 422, RequestBody: []byte(`{}`)})
	r.Observe(Observation{Method: "PUT", Path: "/d", Status: 0, RequestBody: []byte(`{}`)})
	r.Observe(Observation{Method: "PUT", Path: "/e", Status: 0, DryRun: true, RequestBody: []byte(`{}`)})

	reqs := r.Requests()
	want := []string{"", OutcomeOK, OutcomeRejected, OutcomeUnknown, ""}
	for i, w := range want {
		if reqs[i].Outcome != w {
			t.Errorf("request %d (%s %s) outcome = %q, want %q", i, reqs[i].Method, reqs[i].Path, reqs[i].Outcome, w)
		}
	}
	raw := r.Raw()
	if raw[0].ResponseBody != nil || raw[0].RequestBody != nil {
		t.Error("a GET must never keep a body")
	}
	if string(raw[1].RequestBody) == "" || string(raw[1].ResponseBody) != `{"id":1}` {
		t.Errorf("PUT bodies not kept: %+v", raw[1])
	}

	e := Entry{Requests: reqs}
	if !e.HasWrites() || !e.HasUnknownWrites() {
		t.Errorf("HasWrites=%v HasUnknownWrites=%v", e.HasWrites(), e.HasUnknownWrites())
	}

	// ShouldLog: reads-only and dry-run-only entries are dropped under
	// writes_only; an unknown-outcome write is always kept.
	readOnly := Entry{Requests: []Request{{Method: "GET", Path: "/a", Status: 200}}}
	dryOnly := Entry{Requests: []Request{{Method: "PUT", Path: "/e", Status: 0}}, DryRun: true}
	unknownOnly := Entry{Requests: []Request{{Method: "PUT", Path: "/d", Status: 0, Outcome: OutcomeUnknown}}}
	rejectedOnly := Entry{Requests: []Request{{Method: "PUT", Path: "/d", Status: 422, Outcome: OutcomeRejected}}}
	for name, tc := range map[string]struct {
		e    Entry
		want bool
	}{"read": {readOnly, false}, "dry": {dryOnly, false}, "unknown": {unknownOnly, true}, "rejected": {rejectedOnly, true}, "ok": {e, true}} {
		if got := tc.e.ShouldLog(true); got != tc.want {
			t.Errorf("%s: ShouldLog(writesOnly) = %v, want %v", name, got, tc.want)
		}
		if !tc.e.ShouldLog(false) {
			t.Errorf("%s: everything is logged with writes_only off", name)
		}
	}
}

func TestRequestsWithBodies(t *testing.T) {
	r := NewRecorder()
	r.Observe(Observation{Method: "GET", Path: "/a", Status: 200, ResponseBody: []byte(`{"secret":"x"}`)})
	r.Observe(Observation{Method: "PUT", Path: "/b", Status: 200,
		RequestBody:  []byte(`{"comment":{"text_comment":"Nice work, see 7~AbCdEfGhIjKlMnOpQrStUv"},"quiz":{"access_code":"1234","points":1.50,"id":12345678901234567}}`),
		ResponseBody: []byte(`{"id":1,"submission_comments":[{"id":5,"comment":"Nice work"}],"user":{"login_id":"u","integration_id":null,"api_key":"k"}}`)})
	r.Observe(Observation{Method: "POST", Path: "/c", Status: 200, RequestBody: []byte("a=1&b=2&b=3&password=pw"), ResponseBody: []byte("<html>not json</html>")})
	r.Observe(Observation{Method: "DELETE", Path: "/d", Status: 200})

	reqs := r.RequestsWithBodies()
	if reqs[0].Body != nil || reqs[0].Response != nil {
		t.Errorf("GET must carry nothing: %+v", reqs[0])
	}
	raw, _ := json.Marshal(reqs[1])
	line := string(raw)
	for _, want := range []string{`"text_comment":"Nice work, see [REDACTED]"`, `"access_code":"[REDACTED]"`, `"points":1.50`, `"id":12345678901234567`, `"api_key":"[REDACTED]"`, `"integration_id":null`, `"comment":"Nice work"`} {
		if !strings.Contains(line, want) {
			t.Errorf("PUT entry lacks %s:\n%s", want, line)
		}
	}
	if strings.Contains(line, "7~AbCd") || strings.Contains(line, "1234\"") || strings.Contains(line, `"k"`) {
		t.Errorf("secret leaked: %s", line)
	}
	form, _ := reqs[2].Body.(map[string]interface{})
	if form["a"] != "1" || form["password"] != Redacted {
		t.Errorf("form body = %+v", reqs[2].Body)
	}
	if b, _ := form["b"].([]interface{}); len(b) != 2 {
		t.Errorf("repeated form key = %+v", form["b"])
	}
	if text, _ := reqs[2].Response.(map[string]interface{}); text["text"] != "<html>not json</html>" {
		t.Errorf("non-JSON response = %+v", reqs[2].Response)
	}
	if reqs[3].Body != nil || reqs[3].Response != nil {
		t.Errorf("empty bodies must be omitted: %+v", reqs[3])
	}

	// a non-JSON body is cut at MaxTextCapture
	big := strings.Repeat("x", MaxTextCapture+100)
	if text, _ := DecodeResponse([]byte(big)).(map[string]interface{}); len(text["text"].(string)) != MaxTextCapture {
		t.Errorf("text capture not bounded: %d", len(text["text"].(string)))
	}
	// trailing garbage is not JSON
	if m, _ := DecodeBody([]byte(`{"a":1} trailing`)).(map[string]interface{}); m["text"] != `{"a":1} trailing` {
		t.Errorf("a JSON prefix with trailing text must be kept as text, got %v", m)
	}
}

func TestRedact(t *testing.T) {
	in := map[string]interface{}{
		"Authorization": "Bearer abc",
		"nested":        map[string]interface{}{"api-key": "k", "accessCode": "c", "clientSecret": "s", "PASSWORD": "p", "ok": "Bearer xyz"},
		"list":          []interface{}{map[string]interface{}{"token": "t"}, "7~AbCdEfGhIjKlMnOpQrStUv", json.Number("1")},
		"n":             json.Number("2"),
		"b":             true,
	}
	out, _ := json.Marshal(Redact(in))
	s := string(out)
	for _, leak := range []string{"abc", `"k"`, `"c"`, `"s"`, `"p"`, "xyz", `"t"`, "7~AbCd"} {
		if strings.Contains(s, leak) {
			t.Errorf("leaked %s in %s", leak, s)
		}
	}
	for _, keep := range []string{`"n":2`, `"b":true`, `"ok":"Bearer [REDACTED]"`, `1]`} {
		if !strings.Contains(s, keep) {
			t.Errorf("lost %s in %s", keep, s)
		}
	}
	// the input is not mutated
	if in["Authorization"] != "Bearer abc" {
		t.Error("Redact mutated its input")
	}
}

func sampleEntry(ts string, cmd string, reqs ...Request) Entry {
	return Entry{Timestamp: ts, Version: "test", Command: cmd, Args: []string{cmd}, Requests: reqs, Touched: []Touched{}}
}

func TestAppendReadAndRotation(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{Enabled: true, Path: filepath.Join(dir, "sub", "activity.jsonl"), MaxSizeBytes: 300}

	if _, err := Append(cfg, sampleEntry("2026-09-03T10:00:00Z", "courses list", Request{Method: "GET", Path: "/api/v1/courses", Status: 200})); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if _, err := Append(cfg, sampleEntry("2026-09-03T11:00:00Z", "quizzes regrade", Request{Method: "PUT", Path: "/api/v1/courses/1/quizzes/2/questions/3", Status: 200, Outcome: OutcomeOK})); err != nil {
		t.Fatalf("Append: %v", err)
	}
	info, err := os.Stat(cfg.Path)
	if err != nil {
		t.Fatalf("log not created: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("log perm = %o, want 600", info.Mode().Perm())
	}
	if dinfo, _ := os.Stat(filepath.Dir(cfg.Path)); dinfo == nil || dinfo.Mode().Perm() != 0o700 {
		t.Errorf("created log dir perm = %v, want 700", dinfo)
	}

	entries, skipped, err := Read(cfg.Path)
	if err != nil || skipped != 0 || len(entries) != 2 {
		t.Fatalf("Read = %d entries, %d skipped, err %v", len(entries), skipped, err)
	}
	if entries[1].Command != "quizzes regrade" || !entries[1].HasWrites() {
		t.Errorf("entries[1] = %+v", entries[1])
	}

	// the file is now over 300 bytes: the next append archives it first
	res, err := Append(cfg, sampleEntry("2026-09-03T12:00:00Z", "courses get"))
	if err != nil {
		t.Fatalf("Append with rotation: %v", err)
	}
	archived := res.Archived
	if archived == "" || !strings.HasPrefix(filepath.Base(archived), "activity-") || !strings.HasSuffix(archived, ".jsonl") || filepath.Dir(archived) != filepath.Dir(cfg.Path) {
		t.Errorf("archive name = %q", archived)
	}
	old, _, _ := Read(archived)
	current, _, _ := Read(cfg.Path)
	if len(old) != 2 || len(current) != 1 || current[0].Command != "courses get" {
		t.Errorf("after rotation: archive has %d, current has %d", len(old), len(current))
	}

	// a malformed line is skipped, not fatal
	f, _ := os.OpenFile(cfg.Path, os.O_APPEND|os.O_WRONLY, 0o600)
	_, _ = f.WriteString("{not json\n")
	f.Close()
	entries, skipped, err = Read(cfg.Path)
	if err != nil || skipped != 1 || len(entries) != 1 {
		t.Errorf("Read with malformed line = %d entries, %d skipped, err %v", len(entries), skipped, err)
	}

	// JSON shape of a line
	raw, _ := os.ReadFile(cfg.Path)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(strings.SplitN(string(raw), "\n", 2)[0]), &m); err != nil {
		t.Fatalf("line is not JSON: %v", err)
	}
	for _, k := range []string{"ts", "version", "command", "args", "dry_run", "exit_code", "duration_ms", "requests", "touched"} {
		if _, ok := m[k]; !ok {
			t.Errorf("line lacks %q: %s", k, raw)
		}
	}
	for _, k := range []string{"details", "verification_required"} {
		if _, ok := m[k]; ok {
			t.Errorf("%s must be omitted when empty/false: %s", k, raw)
		}
	}
}

func TestPrepareTightensPermissions(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "logs")
	path := filepath.Join(dir, "activity.jsonl")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// the umask may already have narrowed these; make the looseness explicit
	_ = os.Chmod(dir, 0o755)
	_ = os.Chmod(path, 0o644)

	notes, err := Prepare(Config{Path: path})
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if len(notes) != 2 || !strings.Contains(notes[0], "0755 to 0700") || !strings.Contains(notes[1], "0644 to 0600") {
		t.Errorf("notes = %q", notes)
	}
	if info, _ := os.Stat(dir); info.Mode().Perm() != 0o700 {
		t.Errorf("dir perm = %o, want 700", info.Mode().Perm())
	}
	if info, _ := os.Stat(path); info.Mode().Perm() != 0o600 {
		t.Errorf("file perm = %o, want 600", info.Mode().Perm())
	}
	// the content is untouched and a second run is silent
	if raw, _ := os.ReadFile(path); string(raw) != "{}\n" {
		t.Errorf("content changed: %q", raw)
	}
	if notes, err := Prepare(Config{Path: path}); err != nil || len(notes) != 0 {
		t.Errorf("second Prepare = %q, %v", notes, err)
	}

	// a read-only file is reported as unusable and never loosened
	_ = os.Chmod(path, 0o400)
	if _, err := Prepare(Config{Path: path}); err == nil || !strings.Contains(err.Error(), "for append") {
		t.Errorf("Prepare on 0400 must fail as not writable, got %v", err)
	}
	if info, _ := os.Stat(path); info.Mode().Perm() != 0o400 {
		t.Errorf("0400 must not be loosened, got %o", info.Mode().Perm())
	}
	_ = os.Chmod(path, 0o600)

	// Append performs the same tightening and reports it
	_ = os.Chmod(path, 0o666)
	res, err := Append(Config{Path: path, MaxSizeBytes: 1 << 20}, sampleEntry("2026-09-03T10:00:00Z", "courses list"))
	if err != nil || len(res.Notes) != 1 || !strings.Contains(res.Notes[0], "0666 to 0600") {
		t.Errorf("Append notes = %q, %v", res.Notes, err)
	}

	// an unusable path is an error, not a panic
	blocker := filepath.Join(t.TempDir(), "file")
	_ = os.WriteFile(blocker, []byte("x"), 0o600)
	if _, err := Prepare(Config{Path: filepath.Join(blocker, "activity.jsonl")}); err == nil {
		t.Error("Prepare under a file must fail")
	}
}

func TestArchiveAndClear(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "activity.jsonl")

	// missing log: both are no-ops
	if name, err := Archive(path); err != nil || name != "" {
		t.Errorf("Archive(missing) = %q, %v", name, err)
	}
	if err := Clear(path); err != nil {
		t.Errorf("Clear(missing) = %v", err)
	}

	_ = os.WriteFile(path, []byte("{}\n"), 0o600)
	name, err := Archive(path)
	if err != nil || name == "" {
		t.Fatalf("Archive = %q, %v", name, err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("archive must move the current log away")
	}
	if !strings.HasPrefix(filepath.Base(name), "activity-2") { // activity-<YYYY...>Z.jsonl
		t.Errorf("archive name = %q", name)
	}

	// a second archive within the same second gets a distinct name
	_ = os.WriteFile(path, []byte("{}\n"), 0o600)
	name2, err := Archive(path)
	if err != nil || name2 == "" || name2 == name {
		t.Errorf("second archive = %q (first %q), %v", name2, name, err)
	}

	_ = os.WriteFile(path, []byte("{}\n{}\n"), 0o600)
	if err := Clear(path); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if info, _ := os.Stat(path); info.Size() != 0 {
		t.Errorf("Clear must truncate in place, size = %d", info.Size())
	}
}

func TestFilterAndParseSince(t *testing.T) {
	now := time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC)
	entries := []Entry{
		sampleEntry("2026-08-01T12:00:00Z", "courses list", Request{Method: "GET", Path: "/a", Status: 200}),
		sampleEntry("2026-09-03T10:00:00Z", "quizzes regrade", Request{Method: "GET", Path: "/a", Status: 200}, Request{Method: "PUT", Path: "/b", Status: 200, Outcome: OutcomeOK}),
		sampleEntry("2026-09-03T11:30:00Z", "quizzes questions update", Request{Method: "PUT", Path: "/c", Status: 0}), // dry run: no write reached Canvas
		sampleEntry("bad-timestamp", "courses get"),
	}

	if got := Filter(entries, now, 0, false, ""); len(got) != 4 {
		t.Errorf("no filter = %d, want 4", len(got))
	}
	if got := Filter(entries, now, 24*time.Hour, false, ""); len(got) != 2 {
		t.Errorf("since 24h = %d, want 2 (old and unparsable dropped)", len(got))
	}
	if got := Filter(entries, now, 0, true, ""); len(got) != 1 || got[0].Command != "quizzes regrade" {
		t.Errorf("writes only = %v", got)
	}
	if got := Filter(entries, now, 0, false, "quizzes"); len(got) != 2 {
		t.Errorf("prefix quizzes = %d, want 2", len(got))
	}
	if got := Filter(entries, now, 3*time.Hour, true, "quizzes"); len(got) != 1 {
		t.Errorf("combined = %d, want 1", len(got))
	}

	for in, want := range map[string]time.Duration{"": 0, "7d": 7 * 24 * time.Hour, "24h": 24 * time.Hour, "30m": 30 * time.Minute} {
		if got, err := ParseSince(in); err != nil || got != want {
			t.Errorf("ParseSince(%q) = %v, %v; want %v", in, got, err, want)
		}
	}
	for _, in := range []string{"x", "-1d", "7days", "1w"} {
		if _, err := ParseSince(in); err == nil {
			t.Errorf("ParseSince(%q) should fail", in)
		}
	}
}

// A log path straight under $HOME must not get the home directory chmodded
// to 0700; the file is still tightened and the operator is told to move it.
func TestPrepareLeavesHomeDirectoryPermissionsAlone(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file mode bits are not enforced on Windows")
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	_ = os.Chmod(home, 0o755)
	path := filepath.Join(home, "activity.jsonl")

	notes, err := Prepare(Config{Path: path})
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	info, err := os.Stat(home)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o755 {
		t.Errorf("home directory mode changed to %04o", got)
	}
	if len(notes) == 0 || !strings.Contains(notes[0], "home directory") {
		t.Errorf("expected a note about the home directory, got %q", notes)
	}
	if info, err := os.Stat(path); err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("log file should still be 0600, got %v %v", info, err)
	}
}
