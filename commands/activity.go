package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/jjuanrivvera/canvas-cli/internal/activity"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/batch"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

// activityCmd manages the local activity log.
var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "Inspect the local CLI activity log",
	Long: `Inspect and manage the local activity log: one JSON line per canvas
invocation (command, redacted arguments, every HTTP request with its status
and outcome, the Canvas objects touched, exit code and duration).

The log is off by default. Turn it on persistently:

  canvas activity configure --enable

or in the config file:

  activity_log:
    enabled: true
    path: /custom/activity.jsonl   # optional, default <config dir>/activity.jsonl
    writes_only: true              # optional, default true: reads and dry-runs leave no entry
    capture_bodies: false          # optional, record what every write sent and received
    required: false                # optional, audited mode: refuse writes when the log is unusable
    max_size_mb: 10                # optional, archive before exceeding

or with the environment: CANVAS_ACTIVITY_LOG=<path> enables the log and sets
its path; CANVAS_ACTIVITY_WRITES_ONLY, CANVAS_ACTIVITY_CAPTURE_BODIES and
CANVAS_ACTIVITY_REQUIRED (1/true or 0/false) override the matching keys. The
environment takes precedence over the config file.

The log records writes; reads and dry-runs leave no entry unless writes_only
is turned off. A write whose response was lost is always logged with
outcome "unknown" and verification_required, because Canvas may have applied
it. With capture_bodies on, every write's entry carries the full payload sent
(for example the text of a comment, a message or an announcement) and the
response: it is intended for operator audit logs and stores student-directed
text.

Tokens and other secrets never reach the log: secret-looking flag values and
JSON keys, Canvas-shaped tokens and Bearer credentials are redacted, query
strings are dropped from request paths, and CANVAS_TOKEN is never recorded.
The log directory is created owner-only (0700) and the file 0600; looser
permissions on either are tightened with a warning, never loosened.`,
}

func init() {
	rootCmd.AddCommand(activityCmd)
	activityCmd.AddCommand(newActivityListCmd())
	activityCmd.AddCommand(newActivityArchiveCmd())
	activityCmd.AddCommand(newActivityClearCmd())
	activityCmd.AddCommand(newActivityPathCmd())
	activityCmd.AddCommand(newActivityConfigureCmd())
}

// resolveActivityConfig reads the effective activity-log configuration from
// the config file (the same one the rest of the CLI uses) and the
// environment.
func resolveActivityConfig() activity.Config {
	configDir, err := config.GetConfigDir()
	if err != nil {
		configDir = "."
	}
	var settings activity.Settings
	if cfg, err := config.Load(); err == nil && cfg.ActivityLog != nil {
		al := cfg.ActivityLog
		settings = activity.Settings{
			Enabled:       al.Enabled,
			Path:          al.Path,
			WritesOnly:    al.WritesOnly,
			CaptureBodies: al.CaptureBodies,
			Required:      al.Required,
			MaxSizeMB:     al.MaxSizeMB,
		}
	}
	return activity.Resolve(configDir, settings, activity.Env{
		Path:          os.Getenv(activity.EnvVar),
		WritesOnly:    os.Getenv(activity.EnvWritesOnly),
		CaptureBodies: os.Getenv(activity.EnvCaptureBodies),
		Required:      os.Getenv(activity.EnvRequired),
	})
}

// activityEnvVars are the environment variables that override the config file.
var activityEnvVars = []string{activity.EnvVar, activity.EnvWritesOnly, activity.EnvCaptureBodies, activity.EnvRequired}

// runWithActivityLog executes the root command and, afterwards, appends the
// invocation to the activity log when it is enabled. Logging can never fail
// the command: problems go to stderr and the command's own error is returned.
// In audited mode (required) the log is checked before the first write to
// Canvas, and that write is refused when the log cannot be written.
func runWithActivityLog(execute func() (*cobra.Command, error)) error {
	start := time.Now()
	recorder := activity.Default()
	recorder.Reset()
	api.RequestObserver = func(o api.ObservedRequest) {
		recorder.Observe(activity.Observation{
			Method: o.Method, Path: o.Path, Status: o.Status, DryRun: o.DryRun,
			RequestBody: o.RequestBody, ResponseBody: o.ResponseBody,
		})
	}
	api.RequestGate = newActivityGate(resolveActivityConfig)

	cmd, err := execute()

	logActivity(cmd, err, start, recorder, os.Args[1:])
	return err
}

// newActivityGate returns the api.RequestGate for audited mode: before the
// first write of the process it prepares the log (directory 0700, file
// 0600, owned by the current user, writable) and refuses the write when
// that fails. The check runs once; reads are never gated (the client does
// not consult the gate for them).
func newActivityGate(resolve func() activity.Config) func(method, path string) error {
	var once sync.Once
	var result error
	return func(method, path string) error {
		once.Do(func() {
			cfg := resolve()
			if !cfg.Required {
				return
			}
			notes, err := activity.Prepare(cfg)
			warnActivityNotes(notes)
			if err != nil {
				result = fmt.Errorf("audit log unavailable: %v; refusing to write to Canvas", err)
			}
		})
		return result
	}
}

func warnActivityNotes(notes []string) {
	for _, n := range notes {
		fmt.Fprintf(os.Stderr, "warning: activity log: %s\n", n)
	}
}

// recordActivityInput attaches the parsed input of a file- or stdin-driven
// write command (CSV rows, a JSON document) to the activity entry under
// details.input, redacted, so the audit shows what was asked for, not only
// what was sent.
func recordActivityInput(v interface{}) {
	if v == nil {
		return
	}
	activity.Default().SetDetail("input", activity.Redact(v))
}

// recordActivityInputJSON is recordActivityInput for a raw JSON document.
func recordActivityInputJSON(data []byte) {
	recordActivityInput(activity.DecodeBody(data))
}

// gradeRecordsInput renders bulk-grade CSV rows for details.input.
func gradeRecordsInput(grades []batch.GradeRecord) interface{} {
	rows := make([]interface{}, 0, len(grades))
	for _, g := range grades {
		row := map[string]interface{}{"row": g.Row, "user_id": g.UserID, "assignment_id": g.AssignmentID, "grade": g.Grade}
		if g.Comment != "" {
			row["comment"] = g.Comment
		}
		rows = append(rows, row)
	}
	return rows
}

// activitySkippedGroups are commands whose invocations are not logged: the
// log's own management commands, and shell helpers that never touch Canvas.
var activitySkippedGroups = map[string]bool{"activity": true, "help": true, "completion": true, "__complete": true, "__completeNoDesc": true}

func logActivity(cmd *cobra.Command, runErr error, start time.Time, recorder *activity.Recorder, argv []string) {
	if cmd == nil {
		return
	}
	command := commandPathWithoutRoot(cmd)
	if command == "" || activitySkippedGroups[strings.Fields(command)[0]] {
		return
	}

	cfg := resolveActivityConfig()
	if !cfg.Enabled {
		return
	}

	requests := recorder.Requests()
	if cfg.CaptureBodies {
		requests = recorder.RequestsWithBodies()
	}
	touchedFromResponses(command, recorder)
	entry := activity.Entry{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Version:    version,
		Command:    command,
		Args:       activity.RedactArgs(argv),
		DryRun:     dryRun,
		ExitCode:   0,
		DurationMs: time.Since(start).Milliseconds(),
		Requests:   requests,
		Touched:    touchedFromCommand(cmd, recorder),
		Details:    recorder.Details(),
	}
	entry.VerificationRequired = entry.HasUnknownWrites()
	if runErr != nil {
		entry.ExitCode = 1 // main exits 1 on any error
	}

	if !entry.ShouldLog(cfg.WritesOnly) {
		return
	}

	res, err := activity.Append(cfg, entry)
	warnActivityNotes(res.Notes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: activity log not written: %v\n", err)
		return
	}
	if res.Archived != "" && verbose {
		fmt.Fprintf(os.Stderr, "activity log archived to %s\n", res.Archived)
	}
}

func commandPathWithoutRoot(cmd *cobra.Command) string {
	path := cmd.CommandPath()
	rootName := cmd.Root().Name()
	if path == rootName {
		return ""
	}
	return strings.TrimPrefix(path, rootName+" ")
}

// usePlaceholderRe matches "<question-id>"-style positional placeholders in
// a command's Use line.
var usePlaceholderRe = regexp.MustCompile(`<([a-z][a-z-]*)-id>`)

// touchedFromCommand derives the touched objects from the executed command:
// every changed "--<type>-id" flag and "--question", plus positional
// arguments whose placeholder in the Use line is "<type-id>".
func touchedFromCommand(cmd *cobra.Command, recorder *activity.Recorder) []activity.Touched {
	flags := cmd.Flags()
	flags.Visit(func(f *pflag.Flag) {
		var typ string
		switch {
		case f.Name == "question":
			typ = "question"
		case strings.HasSuffix(f.Name, "-id") && f.Name != "as-user":
			typ = strings.TrimSuffix(f.Name, "-id")
		default:
			return
		}
		if id, err := strconv.ParseInt(f.Value.String(), 10, 64); err == nil {
			recorder.Touch(typ, id)
		}
	})

	// Positional placeholders, e.g. "get <question-id>" or "regrade <quiz-id>".
	if matches := usePlaceholderRe.FindAllStringSubmatch(cmd.Use, -1); len(matches) > 0 {
		args := cmd.Flags().Args()
		for i, m := range matches {
			if i >= len(args) {
				break
			}
			if id, err := strconv.ParseInt(args[i], 10, 64); err == nil {
				recorder.Touch(m[1], id)
			}
		}
	}
	return recorder.Touched()
}

// touchedFromResponses adds, for the commands whose write responses carry
// the id of what they created or changed, those ids to the recorder. Only
// successful writes are read; the response bytes are already in memory.
func touchedFromResponses(command string, recorder *activity.Recorder) {
	for _, raw := range recorder.Raw() {
		if raw.IsRead() || raw.Outcome != activity.OutcomeOK || len(raw.ResponseBody) == 0 {
			continue
		}
		resp := activity.DecodeResponse(raw.ResponseBody)
		switch command {
		case "submissions grade", "submissions add-comment":
			obj := jsonObject(resp)
			recorder.Touch("submission", jsonID(obj, "id"))
			// The comment just posted is the one whose text matches the request.
			if text, _ := jsonObject(jsonObject(activity.DecodeBody(raw.RequestBody))["comment"])["text_comment"].(string); text != "" {
				for _, c := range jsonList(obj["submission_comments"]) {
					if cm := jsonObject(c); cm["comment"] == text {
						recorder.Touch("submission-comment", jsonID(cm, "id"))
					}
				}
			}
		case "conversations create":
			for _, c := range jsonList(resp) {
				touchConversation(recorder, jsonObject(c))
			}
		case "conversations reply":
			touchConversation(recorder, jsonObject(resp))
		case "announcements create", "announcements update":
			recorder.Touch("announcement", jsonID(jsonObject(resp), "id"))
		case "discussions post", "discussions reply":
			recorder.Touch("discussion-entry", jsonID(jsonObject(resp), "id"))
		case "pages create", "pages update":
			recorder.Touch("page", jsonID(jsonObject(resp), "page_id"))
		}
	}
}

func touchConversation(recorder *activity.Recorder, conv map[string]interface{}) {
	recorder.Touch("conversation", jsonID(conv, "id"))
	for _, m := range jsonList(conv["messages"]) {
		recorder.Touch("conversation-message", jsonID(jsonObject(m), "id"))
	}
}

func jsonObject(v interface{}) map[string]interface{} {
	m, _ := v.(map[string]interface{})
	return m
}

func jsonList(v interface{}) []interface{} {
	l, _ := v.([]interface{})
	return l
}

// jsonID reads an integer id from a decoded JSON object (numbers are
// json.Number after DecodeResponse, float64 after a plain Unmarshal).
func jsonID(obj map[string]interface{}, key string) int64 {
	switch n := obj[key].(type) {
	case json.Number:
		id, _ := n.Int64()
		return id
	case float64:
		return int64(n)
	}
	return 0
}

// ---- subcommands ----

// activityEntryRow is the compact table view of an entry.
type activityEntryRow struct {
	Timestamp  string `json:"ts"`
	Command    string `json:"command"`
	ExitCode   int    `json:"exit_code"`
	DurationMs int64  `json:"duration_ms"`
	Requests   int    `json:"requests"`
	Writes     int    `json:"writes"`
	Unknown    int    `json:"unknown"`
	DryRun     bool   `json:"dry_run"`
	Touched    string `json:"touched"`
}

func activityRow(e activity.Entry) activityEntryRow {
	touched := make([]string, 0, len(e.Touched))
	for _, t := range e.Touched {
		touched = append(touched, t.Type+":"+strconv.FormatInt(t.ID, 10))
	}
	unknown := 0
	for _, r := range e.Requests {
		if r.Outcome == activity.OutcomeUnknown {
			unknown++
		}
	}
	return activityEntryRow{
		Timestamp:  e.Timestamp,
		Command:    e.Command,
		ExitCode:   e.ExitCode,
		DurationMs: e.DurationMs,
		Requests:   len(e.Requests),
		Writes:     e.WriteCount(),
		Unknown:    unknown,
		DryRun:     e.DryRun,
		Touched:    strings.Join(touched, " "),
	}
}

func newActivityListCmd() *cobra.Command {
	var since, commandPrefix string
	var writesOnly bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List logged invocations",
		Long: `List logged invocations, oldest first. The table shows one line per
invocation; -o json prints the full entries including every request and the
touched objects.

Examples:
  canvas activity list --since 24h
  canvas activity list --writes --command quizzes
  canvas activity list --since 7d -o json | jq '.[] | select(.exit_code != 0)'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sinceDur, err := activity.ParseSince(since)
			if err != nil {
				return err
			}
			cfg := resolveActivityConfig()
			entries, skipped, err := activity.Read(cfg.Path)
			if err != nil {
				return fmt.Errorf("failed to read activity log %s: %w", cfg.Path, err)
			}
			if skipped > 0 {
				fmt.Fprintf(os.Stderr, "warning: %d malformed line(s) skipped in %s\n", skipped, cfg.Path)
			}
			entries = activity.Filter(entries, time.Now(), sinceDur, writesOnly, commandPrefix)
			if entries == nil {
				entries = []activity.Entry{}
			}
			empty := "No activity logged" + activityEnabledHint(cfg)
			if isStructuredOutput() {
				return formatEmptyOrOutput(entries, empty)
			}
			rows := make([]activityEntryRow, 0, len(entries))
			for _, e := range entries {
				rows = append(rows, activityRow(e))
			}
			return formatEmptyOrOutput(rows, empty)
		},
	}

	cmd.Flags().StringVar(&since, "since", "", "Only entries newer than this (e.g. 7d, 24h, 30m)")
	cmd.Flags().BoolVar(&writesOnly, "writes", false, "Only invocations that made a non-GET request")
	cmd.Flags().StringVar(&commandPrefix, "command", "", "Only commands starting with this prefix (e.g. \"quizzes\")")

	return cmd
}

func activityEnabledHint(cfg activity.Config) string {
	if cfg.Enabled {
		return ""
	}
	return " (logging is disabled; see 'canvas activity configure --help')"
}

func newActivityArchiveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "archive",
		Short: "Rename the current log to activity-<UTC timestamp>.jsonl",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := resolveActivityConfig()
			name, err := activity.Archive(cfg.Path)
			if err != nil {
				return fmt.Errorf("failed to archive activity log: %w", err)
			}
			if name == "" {
				printInfo("No activity log at %s\n", cfg.Path)
				return nil
			}
			printInfo("Archived %s to %s\n", cfg.Path, name)
			return nil
		},
	}
}

func newActivityClearCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Truncate the current log (requires --force)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				return fmt.Errorf("refusing to clear the activity log without --force (use 'canvas activity archive' to keep a copy)")
			}
			cfg := resolveActivityConfig()
			if err := activity.Clear(cfg.Path); err != nil {
				return fmt.Errorf("failed to clear activity log: %w", err)
			}
			printInfo("Cleared %s\n", cfg.Path)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Confirm truncating the log")
	return cmd
}

// activitySettingsMap is the structured view of the effective configuration.
func activitySettingsMap(cfg activity.Config) map[string]interface{} {
	return map[string]interface{}{
		"path":           cfg.Path,
		"enabled":        cfg.Enabled,
		"writes_only":    cfg.WritesOnly,
		"capture_bodies": cfg.CaptureBodies,
		"required":       cfg.Required,
		"max_size_mb":    cfg.MaxSizeBytes / (1024 * 1024),
	}
}

// printActivitySettings prints the effective settings, first the bare path
// so scripts can take the first line.
func printActivitySettings(cfg activity.Config) {
	fmt.Println(cfg.Path)
	fmt.Printf("enabled: %t\n", cfg.Enabled)
	fmt.Printf("writes_only: %t\n", cfg.WritesOnly)
	fmt.Printf("capture_bodies: %t\n", cfg.CaptureBodies)
	fmt.Printf("required: %t\n", cfg.Required)
	fmt.Printf("max_size_mb: %d\n", cfg.MaxSizeBytes/(1024*1024))
	var set []string
	for _, name := range activityEnvVars {
		if os.Getenv(name) != "" {
			set = append(set, name)
		}
	}
	if len(set) > 0 {
		fmt.Printf("(environment overrides in effect: %s)\n", strings.Join(set, ", "))
	}
}

func newActivityPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the resolved log path and the effective settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := resolveActivityConfig()
			return formatOutput(activitySettingsMap(cfg), func() { printActivitySettings(cfg) })
		},
	}
}

func newActivityConfigureCmd() *cobra.Command {
	var enable, disable bool
	var path string
	var maxSizeMB int

	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Write activity log settings to the config file",
		Long: `Write the activity_log settings into the CLI config file (the one
'canvas config' uses) and print the effective settings. Only the flags given
are changed; the others keep their current value.

Examples:
  canvas activity configure --enable
  canvas activity configure --enable --capture-bodies --required
  canvas activity configure --writes-only=false          # log reads and dry-runs too
  canvas activity configure --path /var/log/canvas/activity.jsonl --max-size-mb 50
  canvas activity configure --disable

Environment variables (CANVAS_ACTIVITY_LOG, CANVAS_ACTIVITY_WRITES_ONLY,
CANVAS_ACTIVITY_CAPTURE_BODIES, CANVAS_ACTIVITY_REQUIRED) still take
precedence over what is written here.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			if maxSizeMB < 0 {
				return fmt.Errorf("--max-size-mb must not be negative")
			}
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if cfg.ActivityLog == nil {
				cfg.ActivityLog = &config.ActivityLogSettings{}
			}
			al := cfg.ActivityLog
			if enable {
				al.Enabled = true
			}
			if disable {
				al.Enabled = false
			}
			if flags.Changed("path") {
				al.Path = path
			}
			if flags.Changed("writes-only") {
				v, _ := flags.GetBool("writes-only")
				al.WritesOnly = &v
			}
			if flags.Changed("capture-bodies") {
				al.CaptureBodies, _ = flags.GetBool("capture-bodies")
			}
			if flags.Changed("required") {
				al.Required, _ = flags.GetBool("required")
			}
			if flags.Changed("max-size-mb") {
				al.MaxSizeMB = maxSizeMB
			}
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			effective := resolveActivityConfig()
			configPath, _ := config.GetConfigPath()
			return formatOutput(activitySettingsMap(effective), func() {
				printInfo("Activity log settings saved to %s\n", configPath)
				printActivitySettings(effective)
			})
		},
	}

	cmd.Flags().BoolVar(&enable, "enable", false, "Turn the log on")
	cmd.Flags().BoolVar(&disable, "disable", false, "Turn the log off")
	cmd.Flags().StringVar(&path, "path", "", "Log file path (empty resets to <config dir>/activity.jsonl)")
	cmd.Flags().Bool("writes-only", true, "Log only invocations that wrote to Canvas; --writes-only=false logs reads and dry-runs too")
	cmd.Flags().Bool("capture-bodies", false, "Record the full payload sent and the response of every write")
	cmd.Flags().Bool("required", false, "Audited mode: refuse writes to Canvas when the log cannot be written (implies --enable)")
	cmd.Flags().IntVar(&maxSizeMB, "max-size-mb", 0, "Archive the log before it exceeds this many MiB")
	cmd.MarkFlagsMutuallyExclusive("enable", "disable")

	return cmd
}
