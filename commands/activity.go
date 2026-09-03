package commands

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/jjuanrivvera/canvas-cli/internal/activity"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

// activityCmd manages the local activity log.
var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "Inspect the local CLI activity log",
	Long: `Inspect and manage the local activity log: one JSON line per canvas
invocation (command, redacted arguments, every HTTP request with its status,
the Canvas objects touched, exit code and duration).

The log is off by default. Enable it in the config file:

  activity_log:
    enabled: true
    path: /custom/activity.jsonl   # optional, default <config dir>/activity.jsonl
    writes_only: false             # optional, log only invocations that wrote
    max_size_mb: 10                # optional, archive before exceeding

or with the environment variable CANVAS_ACTIVITY_LOG=<path>, which enables
the log and sets its path (it takes precedence over the config file).

Tokens and other secrets never reach the log: secret-looking flag values,
Canvas-shaped tokens and Bearer credentials are redacted, query strings are
dropped from request paths, and CANVAS_TOKEN is never recorded.`,
}

func init() {
	rootCmd.AddCommand(activityCmd)
	activityCmd.AddCommand(newActivityListCmd())
	activityCmd.AddCommand(newActivityArchiveCmd())
	activityCmd.AddCommand(newActivityClearCmd())
	activityCmd.AddCommand(newActivityPathCmd())
}

// resolveActivityConfig reads the effective activity-log configuration from
// the loaded config file (viper) and the environment.
func resolveActivityConfig() activity.Config {
	configDir, err := config.GetConfigDir()
	if err != nil {
		configDir = "."
	}
	return activity.Resolve(configDir, activity.Settings{
		Enabled:    viper.GetBool("activity_log.enabled"),
		Path:       viper.GetString("activity_log.path"),
		WritesOnly: viper.GetBool("activity_log.writes_only"),
		MaxSizeMB:  viper.GetInt("activity_log.max_size_mb"),
	}, os.Getenv(activity.EnvVar))
}

// runWithActivityLog executes the root command and, afterwards, appends the
// invocation to the activity log when it is enabled. Logging can never fail
// the command: problems go to stderr and the command's own error is returned.
func runWithActivityLog(execute func() (*cobra.Command, error)) error {
	start := time.Now()
	recorder := activity.Default()
	recorder.Reset()
	api.RequestObserver = recorder.Record

	cmd, err := execute()

	logActivity(cmd, err, start, recorder, os.Args[1:])
	return err
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

	entry := activity.Entry{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Version:    version,
		Command:    command,
		Args:       activity.RedactArgs(argv),
		DryRun:     dryRun,
		ExitCode:   0,
		DurationMs: time.Since(start).Milliseconds(),
		Requests:   recorder.Requests(),
		Touched:    touchedFromCommand(cmd, recorder),
		Details:    recorder.Details(),
	}
	if runErr != nil {
		entry.ExitCode = 1 // main exits 1 on any error
	}

	if cfg.WritesOnly && !entry.HasWrites() {
		return
	}

	archived, err := activity.Append(cfg, entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: activity log not written: %v\n", err)
		return
	}
	if archived != "" && verbose {
		fmt.Fprintf(os.Stderr, "activity log archived to %s\n", archived)
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

// ---- subcommands ----

// activityEntryRow is the compact table view of an entry.
type activityEntryRow struct {
	Timestamp  string `json:"ts"`
	Command    string `json:"command"`
	ExitCode   int    `json:"exit_code"`
	DurationMs int64  `json:"duration_ms"`
	Requests   int    `json:"requests"`
	Writes     int    `json:"writes"`
	DryRun     bool   `json:"dry_run"`
	Touched    string `json:"touched"`
}

func activityRow(e activity.Entry) activityEntryRow {
	touched := make([]string, 0, len(e.Touched))
	for _, t := range e.Touched {
		touched = append(touched, t.Type+":"+strconv.FormatInt(t.ID, 10))
	}
	return activityEntryRow{
		Timestamp:  e.Timestamp,
		Command:    e.Command,
		ExitCode:   e.ExitCode,
		DurationMs: e.DurationMs,
		Requests:   len(e.Requests),
		Writes:     e.WriteCount(),
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
	return " (logging is disabled; see 'canvas activity --help')"
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

func newActivityPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the resolved log path and whether logging is enabled",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := resolveActivityConfig()
			return formatOutput(map[string]interface{}{
				"path":        cfg.Path,
				"enabled":     cfg.Enabled,
				"writes_only": cfg.WritesOnly,
				"max_size_mb": cfg.MaxSizeBytes / (1024 * 1024),
			}, func() {
				fmt.Println(cfg.Path)
				if cfg.Enabled {
					fmt.Println("enabled: true")
				} else {
					fmt.Println("enabled: false")
				}
			})
		},
	}
}
