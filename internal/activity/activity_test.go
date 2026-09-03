package activity

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	if (Entry{Requests: []Request{{"GET", "/x", 200}}}).HasWrites() {
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

func TestResolve(t *testing.T) {
	cfg := Resolve("/cfg", Settings{}, "")
	if cfg.Enabled || cfg.Path != filepath.Join("/cfg", DefaultFileName) || cfg.MaxSizeBytes != DefaultMaxSizeMB*1024*1024 || cfg.WritesOnly {
		t.Errorf("defaults = %+v", cfg)
	}
	cfg = Resolve("/cfg", Settings{Enabled: true, Path: "/custom.jsonl", WritesOnly: true, MaxSizeMB: 2}, "")
	if !cfg.Enabled || cfg.Path != "/custom.jsonl" || !cfg.WritesOnly || cfg.MaxSizeBytes != 2*1024*1024 {
		t.Errorf("config file = %+v", cfg)
	}
	// env wins and enables on its own
	cfg = Resolve("/cfg", Settings{Enabled: false, Path: "/custom.jsonl"}, "/env.jsonl")
	if !cfg.Enabled || cfg.Path != "/env.jsonl" {
		t.Errorf("env = %+v", cfg)
	}
}

func sampleEntry(ts string, cmd string, reqs ...Request) Entry {
	return Entry{Timestamp: ts, Version: "test", Command: cmd, Args: []string{cmd}, Requests: reqs, Touched: []Touched{}}
}

func TestAppendReadAndRotation(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{Enabled: true, Path: filepath.Join(dir, "sub", "activity.jsonl"), MaxSizeBytes: 300}

	if _, err := Append(cfg, sampleEntry("2026-09-03T10:00:00Z", "courses list", Request{"GET", "/api/v1/courses", 200})); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if _, err := Append(cfg, sampleEntry("2026-09-03T11:00:00Z", "quizzes regrade", Request{"PUT", "/api/v1/courses/1/quizzes/2/questions/3", 200})); err != nil {
		t.Fatalf("Append: %v", err)
	}
	info, err := os.Stat(cfg.Path)
	if err != nil {
		t.Fatalf("log not created: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("log perm = %o, want 600", info.Mode().Perm())
	}

	entries, skipped, err := Read(cfg.Path)
	if err != nil || skipped != 0 || len(entries) != 2 {
		t.Fatalf("Read = %d entries, %d skipped, err %v", len(entries), skipped, err)
	}
	if entries[1].Command != "quizzes regrade" || !entries[1].HasWrites() {
		t.Errorf("entries[1] = %+v", entries[1])
	}

	// the file is now over 300 bytes: the next append archives it first
	archived, err := Append(cfg, sampleEntry("2026-09-03T12:00:00Z", "courses get"))
	if err != nil {
		t.Fatalf("Append with rotation: %v", err)
	}
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
	if _, ok := m["details"]; ok {
		t.Errorf("details must be omitted when empty: %s", raw)
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
		sampleEntry("2026-08-01T12:00:00Z", "courses list", Request{"GET", "/a", 200}),
		sampleEntry("2026-09-03T10:00:00Z", "quizzes regrade", Request{"GET", "/a", 200}, Request{"PUT", "/b", 200}),
		sampleEntry("2026-09-03T11:30:00Z", "quizzes questions update", Request{"PUT", "/c", 0}), // dry run: no write reached Canvas
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
