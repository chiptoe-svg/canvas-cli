// Package activity implements the local CLI activity log: one JSON line per
// invocation, written at exit, recording which command ran, which HTTP
// requests it made and which Canvas objects it touched. It is off by
// default and never contains credentials.
package activity

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// EnvVar names the environment variable that enables the log and sets its
// path in one go. It takes precedence over the config file.
const EnvVar = "CANVAS_ACTIVITY_LOG"

// DefaultFileName is the log file name inside the canvas-cli config dir.
const DefaultFileName = "activity.jsonl"

// DefaultMaxSizeMB is the size at which the log is archived before writing.
const DefaultMaxSizeMB = 10

// Redacted replaces every secret-looking value in logged arguments.
const Redacted = "[REDACTED]"

// Request is one HTTP request made during the invocation.
type Request struct {
	Method string `json:"method"`
	Path   string `json:"path"`   // URL path only, never the query string
	Status int    `json:"status"` // 0 when no response was received (dry run, transport error)
}

// Touched is a Canvas object the invocation addressed, as far as the flags
// and arguments make cheaply known.
type Touched struct {
	Type string `json:"type"`
	ID   int64  `json:"id"`
}

// Entry is one logged invocation.
type Entry struct {
	Timestamp  string                 `json:"ts"`
	Version    string                 `json:"version"`
	Command    string                 `json:"command"`
	Args       []string               `json:"args"`
	DryRun     bool                   `json:"dry_run"`
	ExitCode   int                    `json:"exit_code"`
	DurationMs int64                  `json:"duration_ms"`
	Requests   []Request              `json:"requests"`
	Touched    []Touched              `json:"touched"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// HasWrites reports whether the entry made at least one non-GET request
// that reached Canvas.
func (e Entry) HasWrites() bool {
	for _, r := range e.Requests {
		if r.Method != "GET" && r.Method != "HEAD" && r.Status > 0 {
			return true
		}
	}
	return false
}

// WriteCount counts the non-GET requests that reached Canvas.
func (e Entry) WriteCount() int {
	n := 0
	for _, r := range e.Requests {
		if r.Method != "GET" && r.Method != "HEAD" && r.Status > 0 {
			n++
		}
	}
	return n
}

// Recorder accumulates what one process does. It is safe for concurrent
// use (batch commands issue requests from several goroutines).
type Recorder struct {
	mu       sync.Mutex
	requests []Request
	touched  []Touched
	details  map[string]interface{}
}

// NewRecorder returns an empty recorder.
func NewRecorder() *Recorder {
	return &Recorder{details: map[string]interface{}{}}
}

var defaultRecorder = NewRecorder()

// Default is the process-wide recorder the CLI writes at exit.
func Default() *Recorder { return defaultRecorder }

// Record adds one HTTP request. path may carry a query string; it is
// dropped so tokens or search terms never reach the log.
func (r *Recorder) Record(method, path string, status int) {
	if i := strings.IndexByte(path, '?'); i >= 0 {
		path = path[:i]
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests = append(r.requests, Request{Method: method, Path: path, Status: status})
}

// Touch notes a Canvas object the command addressed. Duplicates are dropped.
func (r *Recorder) Touch(typ string, id int64) {
	if id <= 0 || typ == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, t := range r.touched {
		if t.Type == typ && t.ID == id {
			return
		}
	}
	r.touched = append(r.touched, Touched{Type: typ, ID: id})
}

// SetDetail attaches a command-specific structure (for example the regrade
// verification table) under key.
func (r *Recorder) SetDetail(key string, v interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.details[key] = v
}

// Requests returns a copy of the recorded requests.
func (r *Recorder) Requests() []Request {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]Request(nil), r.requests...)
}

// Touched returns a copy of the touched objects, sorted for stable output.
func (r *Recorder) Touched() []Touched {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := append([]Touched(nil), r.touched...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Type != out[j].Type {
			return out[i].Type < out[j].Type
		}
		return out[i].ID < out[j].ID
	})
	return out
}

// Details returns a copy of the attached details (nil when empty).
func (r *Recorder) Details() map[string]interface{} {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.details) == 0 {
		return nil
	}
	out := make(map[string]interface{}, len(r.details))
	for k, v := range r.details {
		out[k] = v
	}
	return out
}

// Reset clears the recorder (tests).
func (r *Recorder) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests, r.touched, r.details = nil, nil, map[string]interface{}{}
}

// ---- redaction ----

// secretFlags are flag names whose value must never be logged. Matched as a
// substring of the flag name, so --token, --api-token, --client-secret,
// --access-code and --password are all covered.
var secretFlags = []string{"token", "secret", "password", "passwd", "authorization", "access-code", "api-key", "apikey"}

var (
	// canvasTokenRe matches a Canvas API token (e.g. "7~AbC...") anywhere in
	// an argument, so a token pasted as a bare value or inside a header is
	// still caught.
	canvasTokenRe = regexp.MustCompile(`\d+~[A-Za-z0-9]{16,}`)
	// bearerRe matches "Bearer <value>" inside an argument such as a raw
	// -H "Authorization: Bearer ..." header for "canvas api".
	bearerRe = regexp.MustCompile(`(?i)(bearer\s+)\S+`)
)

func isSecretFlag(name string) bool {
	name = strings.ToLower(strings.TrimLeft(name, "-"))
	for _, s := range secretFlags {
		if strings.Contains(name, s) {
			return true
		}
	}
	return false
}

// RedactArgs returns a copy of argv (without the binary) with secret values
// replaced: the value of any token/secret/password/authorization flag in
// both "--flag value" and "--flag=value" form, Canvas-shaped tokens, and
// Bearer credentials embedded in header arguments.
func RedactArgs(args []string) []string {
	out := make([]string, len(args))
	redactNext := false
	for i, a := range args {
		switch {
		case redactNext:
			out[i] = Redacted
			redactNext = false
		case strings.HasPrefix(a, "-") && strings.Contains(a, "="):
			name, _, _ := strings.Cut(a, "=")
			if isSecretFlag(name) {
				out[i] = name + "=" + Redacted
			} else {
				out[i] = redactValue(a)
			}
		case strings.HasPrefix(a, "-") && isSecretFlag(a):
			out[i] = a
			redactNext = true
		default:
			out[i] = redactValue(a)
		}
	}
	return out
}

func redactValue(s string) string {
	s = canvasTokenRe.ReplaceAllString(s, Redacted)
	s = bearerRe.ReplaceAllString(s, "${1}"+Redacted)
	return s
}

// ---- configuration ----

// Config is the resolved activity-log configuration.
type Config struct {
	Enabled      bool
	Path         string
	WritesOnly   bool
	MaxSizeBytes int64
}

// Settings mirrors the config-file block:
//
//	activity_log:
//	  enabled: true
//	  path: /custom/activity.jsonl   # optional
//	  writes_only: false             # optional
//	  max_size_mb: 10                # optional
type Settings struct {
	Enabled    bool
	Path       string
	WritesOnly bool
	MaxSizeMB  int
}

// Resolve combines the config-file settings, the environment variable and
// the config directory into the effective configuration. The environment
// variable wins: when set, logging is enabled and its value is the path.
func Resolve(configDir string, s Settings, envValue string) Config {
	cfg := Config{
		Enabled:      s.Enabled,
		Path:         s.Path,
		WritesOnly:   s.WritesOnly,
		MaxSizeBytes: int64(s.MaxSizeMB) * 1024 * 1024,
	}
	if envValue != "" {
		cfg.Enabled = true
		cfg.Path = envValue
	}
	if cfg.Path == "" {
		cfg.Path = filepath.Join(configDir, DefaultFileName)
	}
	if cfg.MaxSizeBytes <= 0 {
		cfg.MaxSizeBytes = DefaultMaxSizeMB * 1024 * 1024
	}
	return cfg
}

// ---- writing ----

// Append writes one entry as a JSON line, archiving the file first when it
// exceeds cfg.MaxSizeBytes. The archive name is returned when one was made.
func Append(cfg Config, e Entry) (archived string, err error) {
	if info, statErr := os.Stat(cfg.Path); statErr == nil && info.Size() > cfg.MaxSizeBytes {
		archived, err = Archive(cfg.Path)
		if err != nil {
			return "", err
		}
	}

	if e.Requests == nil {
		e.Requests = []Request{}
	}
	if e.Touched == nil {
		e.Touched = []Touched{}
	}
	if e.Args == nil {
		e.Args = []string{}
	}
	line, err := json.Marshal(e)
	if err != nil {
		return archived, fmt.Errorf("encode activity entry: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(cfg.Path), 0o700); err != nil {
		return archived, err
	}
	f, err := os.OpenFile(cfg.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600) // #nosec G304 -- operator-configured log path
	if err != nil {
		return archived, err
	}
	defer f.Close()
	if _, err := f.Write(append(line, '\n')); err != nil {
		return archived, err
	}
	return archived, nil
}

// ArchiveName returns the archive file name for path at time t.
func ArchiveName(path string, t time.Time) string {
	dir := filepath.Dir(path)
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return filepath.Join(dir, fmt.Sprintf("%s-%s.jsonl", base, t.UTC().Format("20060102T150405Z")))
}

// Archive renames the current log next to itself with a UTC timestamp and
// returns the new name. A missing log is not an error (returns "").
func Archive(path string) (string, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	name := ArchiveName(path, time.Now())
	for i := 1; ; i++ {
		if _, err := os.Stat(name); os.IsNotExist(err) {
			break
		}
		name = strings.TrimSuffix(ArchiveName(path, time.Now()), ".jsonl") + "-" + strconv.Itoa(i) + ".jsonl"
	}
	if err := os.Rename(path, name); err != nil {
		return "", err
	}
	return name, nil
}

// Clear truncates the log in place. A missing log is not an error.
func Clear(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.Truncate(path, 0)
}

// ---- reading ----

// Read loads every entry of the log, oldest first. A missing log yields no
// entries. Malformed lines are skipped and counted in the second result.
func Read(path string) ([]Entry, int, error) {
	f, err := os.Open(path) // #nosec G304 -- operator-configured log path
	if err != nil {
		if os.IsNotExist(err) {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	defer f.Close()

	var entries []Entry
	skipped := 0
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			skipped++
			continue
		}
		entries = append(entries, e)
	}
	return entries, skipped, scanner.Err()
}

// Filter selects entries newer than now-since (0 = all), optionally only
// those with writes, optionally only commands starting with prefix.
func Filter(entries []Entry, now time.Time, since time.Duration, writesOnly bool, commandPrefix string) []Entry {
	var out []Entry
	cutoff := now.Add(-since)
	for _, e := range entries {
		if since > 0 {
			ts, err := time.Parse(time.RFC3339, e.Timestamp)
			if err != nil || ts.Before(cutoff) {
				continue
			}
		}
		if writesOnly && !e.HasWrites() {
			continue
		}
		if commandPrefix != "" && !strings.HasPrefix(e.Command, commandPrefix) {
			continue
		}
		out = append(out, e)
	}
	return out
}

// ParseSince parses "7d", "24h", "30m", "90s" (days are 24h).
func ParseSince(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	if strings.HasSuffix(s, "d") {
		days, err := strconv.Atoi(strings.TrimSuffix(s, "d"))
		if err != nil || days < 0 {
			return 0, fmt.Errorf("invalid duration %q: expected e.g. 7d, 24h, 30m", s)
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil || d < 0 {
		return 0, fmt.Errorf("invalid duration %q: expected e.g. 7d, 24h, 30m", s)
	}
	return d, nil
}
