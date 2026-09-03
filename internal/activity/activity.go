// Package activity implements the local CLI activity log: one JSON line per
// invocation, written at exit, recording which command ran, which HTTP
// requests it made and which Canvas objects it touched. It is off by
// default and never contains credentials.
package activity

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
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

// EnvWritesOnly ("1"/"true" or "0"/"false") overrides activity_log.writes_only.
const EnvWritesOnly = "CANVAS_ACTIVITY_WRITES_ONLY"

// EnvCaptureBodies ("1"/"true" or "0"/"false") overrides activity_log.capture_bodies.
const EnvCaptureBodies = "CANVAS_ACTIVITY_CAPTURE_BODIES"

// EnvRequired ("1"/"true") overrides activity_log.required: audited mode,
// where a write to Canvas is refused unless the log is writable first.
const EnvRequired = "CANVAS_ACTIVITY_REQUIRED"

// Outcomes of a non-GET request.
const (
	OutcomeOK       = "ok"       // Canvas answered 2xx/3xx
	OutcomeRejected = "rejected" // Canvas answered 4xx/5xx
	OutcomeUnknown  = "unknown"  // no response: Canvas may or may not have applied it
)

// MaxTextCapture bounds a captured non-JSON body: only its first 4 KiB are kept.
const MaxTextCapture = 4 * 1024

// DefaultFileName is the log file name inside the canvas-cli config dir.
const DefaultFileName = "activity.jsonl"

// DefaultMaxSizeMB is the size at which the log is archived before writing.
const DefaultMaxSizeMB = 10

// Redacted replaces every secret-looking value in logged arguments.
const Redacted = "[REDACTED]"

// Request is one HTTP request made during the invocation. Body and Response
// are only present when capture_bodies is on, and only for non-GET requests:
// the payload as sent and the parsed response, both redacted.
type Request struct {
	Method   string      `json:"method"`
	Path     string      `json:"path"`              // URL path only, never the query string
	Status   int         `json:"status"`            // 0 when no response was received (dry run, transport error)
	Outcome  string      `json:"outcome,omitempty"` // non-GET, sent: ok | rejected | unknown
	Body     interface{} `json:"body,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

// RawRequest is a recorded request with the bytes that went over the wire.
type RawRequest struct {
	Request
	RequestBody  []byte
	ResponseBody []byte
}

// IsRead reports whether the request could not have changed anything.
func (r Request) IsRead() bool { return r.Method == "GET" || r.Method == "HEAD" }

// Touched is a Canvas object the invocation addressed, as far as the flags
// and arguments make cheaply known.
type Touched struct {
	Type string `json:"type"`
	ID   int64  `json:"id"`
}

// Entry is one logged invocation. VerificationRequired is set when a write
// was sent but its outcome is unknown: Canvas may have applied it.
type Entry struct {
	Timestamp            string                 `json:"ts"`
	Version              string                 `json:"version"`
	Command              string                 `json:"command"`
	Args                 []string               `json:"args"`
	DryRun               bool                   `json:"dry_run"`
	ExitCode             int                    `json:"exit_code"`
	DurationMs           int64                  `json:"duration_ms"`
	Requests             []Request              `json:"requests"`
	Touched              []Touched              `json:"touched"`
	VerificationRequired bool                   `json:"verification_required,omitempty"`
	Details              map[string]interface{} `json:"details,omitempty"`
}

// HasUnknownWrites reports whether a non-GET request was sent without a
// response arriving (transport error, timeout).
func (e Entry) HasUnknownWrites() bool {
	for _, r := range e.Requests {
		if r.Outcome == OutcomeUnknown {
			return true
		}
	}
	return false
}

// ShouldLog decides whether the entry is written under writesOnly: entries
// with a write that reached Canvas, or one whose outcome is unknown, always
// are; with writesOnly off everything is.
func (e Entry) ShouldLog(writesOnly bool) bool {
	return !writesOnly || e.HasWrites() || e.HasUnknownWrites()
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
	requests []RawRequest
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
	r.RecordBodies(method, path, status, nil, nil)
}

// RecordBodies is Record with the request payload and the response body as
// sent and received.
func (r *Recorder) RecordBodies(method, path string, status int, requestBody, responseBody []byte) {
	r.Observe(Observation{Method: method, Path: path, Status: status, RequestBody: requestBody, ResponseBody: responseBody})
}

// Observation is everything the HTTP client reports about one request.
type Observation struct {
	Method       string
	Path         string
	Status       int  // 0 when no response was received
	DryRun       bool // the request was rendered, not sent
	RequestBody  []byte
	ResponseBody []byte
}

// Observe records one request. The bytes are kept in memory only; they
// reach the log solely through RequestsWithBodies, redacted. A body is
// never stored for a GET or HEAD. A sent non-GET request gets its Outcome:
// ok, rejected, or unknown when no response arrived.
func (r *Recorder) Observe(o Observation) {
	path := o.Path
	if i := strings.IndexByte(path, '?'); i >= 0 {
		path = path[:i]
	}
	req := RawRequest{Request: Request{Method: o.Method, Path: path, Status: o.Status}}
	if !req.IsRead() {
		req.RequestBody = append([]byte(nil), o.RequestBody...)
		req.ResponseBody = append([]byte(nil), o.ResponseBody...)
		switch {
		case o.DryRun:
		case o.Status == 0:
			req.Outcome = OutcomeUnknown
		case o.Status >= 400:
			req.Outcome = OutcomeRejected
		default:
			req.Outcome = OutcomeOK
		}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests = append(r.requests, req)
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

// Requests returns a copy of the recorded requests without their bodies.
func (r *Recorder) Requests() []Request {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Request, 0, len(r.requests))
	for _, req := range r.requests {
		out = append(out, req.Request)
	}
	return out
}

// RequestsWithBodies returns the recorded requests with, for every non-GET
// request, the payload sent (Body) and the response received (Response),
// decoded and redacted. Reads never carry either.
func (r *Recorder) RequestsWithBodies() []Request {
	raw := r.Raw()
	out := make([]Request, 0, len(raw))
	for _, req := range raw {
		e := req.Request
		if !e.IsRead() {
			e.Body = Redact(DecodeBody(req.RequestBody))
			e.Response = Redact(DecodeResponse(req.ResponseBody))
		}
		out = append(out, e)
	}
	return out
}

// Raw returns a copy of the recorded requests with their wire bytes.
func (r *Recorder) Raw() []RawRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]RawRequest(nil), r.requests...)
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

// secretKeys are JSON/form key fragments whose values are never logged.
// Keys are compared lower-cased with "-" folded to "_", so access_code,
// accessCode and access-code all match.
var secretKeys = []string{"token", "secret", "password", "passwd", "authorization", "access_code", "accesscode", "api_key", "apikey"}

func isSecretKey(name string) bool {
	name = strings.ToLower(strings.ReplaceAll(name, "-", "_"))
	for _, s := range secretKeys {
		if strings.Contains(name, s) {
			return true
		}
	}
	return false
}

// Redact returns v (a decoded JSON value) with the value of every secret
// key replaced by Redacted at any depth, and Canvas-shaped tokens and
// Bearer credentials inside any string replaced as well.
func Redact(v interface{}) interface{} {
	switch x := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(x))
		for k, val := range x {
			if isSecretKey(k) {
				out[k] = Redacted
			} else {
				out[k] = Redact(val)
			}
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(x))
		for i, val := range x {
			out[i] = Redact(val)
		}
		return out
	case string:
		return redactValue(x)
	default:
		return v
	}
}

// DecodeBody turns a request payload into a JSON value: JSON is parsed as
// is (numbers kept verbatim), a form-encoded body becomes an object, and
// anything else is {"text": <first 4 KiB>}. Empty yields nil.
func DecodeBody(b []byte) interface{} {
	if len(bytes.TrimSpace(b)) == 0 {
		return nil
	}
	if v, ok := decodeJSON(b); ok {
		return v
	}
	if form, ok := decodeForm(b); ok {
		return form
	}
	return textValue(b)
}

// DecodeResponse turns a response body into a JSON value, or
// {"text": <first 4 KiB>} when it is not JSON. Empty yields nil.
func DecodeResponse(b []byte) interface{} {
	if len(bytes.TrimSpace(b)) == 0 {
		return nil
	}
	if v, ok := decodeJSON(b); ok {
		return v
	}
	return textValue(b)
}

func decodeJSON(b []byte) (interface{}, bool) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var v interface{}
	if err := dec.Decode(&v); err != nil {
		return nil, false
	}
	if _, err := dec.Token(); err != io.EOF { // trailing garbage: not a JSON document
		return nil, false
	}
	return v, true
}

func decodeForm(b []byte) (map[string]interface{}, bool) {
	s := string(b)
	if strings.ContainsAny(s, " \n{}\"") || !strings.Contains(s, "=") {
		return nil, false
	}
	values, err := url.ParseQuery(s)
	if err != nil || len(values) == 0 {
		return nil, false
	}
	out := make(map[string]interface{}, len(values))
	for k, v := range values {
		if len(v) == 1 {
			out[k] = v[0]
		} else {
			items := make([]interface{}, len(v))
			for i := range v {
				items[i] = v[i]
			}
			out[k] = items
		}
	}
	return out, true
}

func textValue(b []byte) interface{} {
	if len(b) > MaxTextCapture {
		b = b[:MaxTextCapture]
	}
	return map[string]interface{}{"text": string(b)}
}

// ---- configuration ----

// Config is the resolved activity-log configuration.
type Config struct {
	Enabled       bool
	Path          string
	WritesOnly    bool
	CaptureBodies bool
	Required      bool
	MaxSizeBytes  int64
}

// Settings mirrors the config-file block:
//
//	activity_log:
//	  enabled: true
//	  path: /custom/activity.jsonl   # optional
//	  writes_only: true              # optional, default true
//	  capture_bodies: false          # optional
//	  required: false                # optional: refuse writes when the log is unusable
//	  max_size_mb: 10                # optional
//
// WritesOnly is nil when the key is absent (default: true).
type Settings struct {
	Enabled       bool
	Path          string
	WritesOnly    *bool
	CaptureBodies bool
	Required      bool
	MaxSizeMB     int
}

// Env carries the raw values of the activity environment variables
// (CANVAS_ACTIVITY_LOG, CANVAS_ACTIVITY_WRITES_ONLY,
// CANVAS_ACTIVITY_CAPTURE_BODIES, CANVAS_ACTIVITY_REQUIRED); empty means
// unset.
type Env struct {
	Path          string
	WritesOnly    string
	CaptureBodies string
	Required      string
}

// ParseBoolEnv reads a boolean environment value: "1"/"true"/"yes"/"on" are
// true, "0"/"false"/"no"/"off" are false; anything else (including empty)
// is not set.
func ParseBoolEnv(s string) (value, set bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "on":
		return true, true
	case "0", "false", "no", "off":
		return false, true
	}
	return false, false
}

// Resolve combines the config-file settings, the environment and the config
// directory into the effective configuration. The environment wins:
// CANVAS_ACTIVITY_LOG enables logging and sets the path, and the boolean
// variables override their config keys. writes_only defaults to true:
// reads and dry-runs leave no entry unless it is explicitly turned off.
// required implies enabled.
func Resolve(configDir string, s Settings, env Env) Config {
	cfg := Config{
		Enabled:       s.Enabled,
		Path:          s.Path,
		WritesOnly:    true,
		CaptureBodies: s.CaptureBodies,
		Required:      s.Required,
		MaxSizeBytes:  int64(s.MaxSizeMB) * 1024 * 1024,
	}
	if s.WritesOnly != nil {
		cfg.WritesOnly = *s.WritesOnly
	}
	if env.Path != "" {
		cfg.Enabled = true
		cfg.Path = env.Path
	}
	if v, ok := ParseBoolEnv(env.WritesOnly); ok {
		cfg.WritesOnly = v
	}
	if v, ok := ParseBoolEnv(env.CaptureBodies); ok {
		cfg.CaptureBodies = v
	}
	if v, ok := ParseBoolEnv(env.Required); ok {
		cfg.Required = v
	}
	if cfg.Required {
		cfg.Enabled = true
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

// AppendResult reports what Append did besides writing the line.
type AppendResult struct {
	Archived string   // archive file name when the log was rotated first
	Notes    []string // permission tightening performed, for a stderr warning
}

// Prepare makes sure the log can be written: the parent directory exists
// (created 0700 when missing), the file exists (created 0600) and is owned
// by the current user, and any looser permissions on either — when owned by
// the current user — are tightened to 0700/0600, never loosened. It returns
// a note for each tightening. This is the audited-mode preflight and the
// first step of every Append.
func Prepare(cfg Config) ([]string, error) {
	var notes []string
	dir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("create %s: %w", dir, err)
	}
	n, err := tighten(dir, 0o700, false)
	if err != nil {
		return nil, err
	}
	notes = append(notes, n...)

	f, err := os.OpenFile(cfg.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600) // #nosec G304 -- operator-configured log path
	if err != nil {
		return nil, fmt.Errorf("open %s for append: %w", cfg.Path, err)
	}
	_ = f.Close()
	n, err = tighten(cfg.Path, 0o600, true)
	if err != nil {
		return nil, err
	}
	return append(notes, n...), nil
}

// tighten chmods path down to max when it is owned by the current user and
// carries any permission bit outside max. Something owned by someone else
// is left alone — or, with requireOwner (the log file itself), refused: the
// log must be the operator's own.
func tighten(path string, max os.FileMode, requireOwner bool) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if owned, known := ownedByCurrentUser(info); known && !owned {
		if requireOwner {
			return nil, fmt.Errorf("%s is not owned by the current user", path)
		}
		return nil, nil
	}
	perm := info.Mode().Perm()
	if perm&^max == 0 {
		return nil, nil
	}
	if err := os.Chmod(path, perm&max); err != nil {
		return nil, fmt.Errorf("tighten permissions on %s: %w", path, err)
	}
	return []string{fmt.Sprintf("tightened permissions on %s from %04o to %04o", path, perm, perm&max)}, nil
}

// Append writes one entry as a JSON line, preparing the file first (see
// Prepare) and archiving it when it exceeds cfg.MaxSizeBytes.
func Append(cfg Config, e Entry) (AppendResult, error) {
	var res AppendResult
	notes, err := Prepare(cfg)
	if err != nil {
		return res, err
	}
	res.Notes = notes

	if info, statErr := os.Stat(cfg.Path); statErr == nil && info.Size() > cfg.MaxSizeBytes {
		res.Archived, err = Archive(cfg.Path)
		if err != nil {
			return res, err
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
		return res, fmt.Errorf("encode activity entry: %w", err)
	}

	f, err := os.OpenFile(cfg.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600) // #nosec G304 -- operator-configured log path
	if err != nil {
		return res, err
	}
	defer f.Close()
	if _, err := f.Write(append(line, '\n')); err != nil {
		return res, err
	}
	return res, nil
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
