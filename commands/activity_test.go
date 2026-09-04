package commands

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
	"github.com/jjuanrivvera/canvas-cli/internal/activity"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

// useTempActivityLog points the activity log at a temp file via the
// environment variable (which also enables it) and returns the path.
func useTempActivityLog(t *testing.T) string {
	t.Helper()
	// in its own 0700 directory: a bare t.TempDir() is 0755 and would be
	// tightened (with a warning) on the first write
	path := filepath.Join(t.TempDir(), "log", "activity.jsonl")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv(activity.EnvVar, path)
	return path
}

// fakeExecuted builds a command that looks like an executed
// "quizzes regrade <quiz-id>" invocation with parsed flags.
func fakeExecuted(t *testing.T, args ...string) *cobra.Command {
	t.Helper()
	root := &cobra.Command{Use: "canvas"}
	quizzes := &cobra.Command{Use: "quizzes"}
	regrade := &cobra.Command{Use: "regrade <quiz-id>", RunE: func(*cobra.Command, []string) error { return nil }}
	regrade.Flags().Int64("course-id", 0, "")
	regrade.Flags().Int64("question", 0, "")
	regrade.Flags().Int64("correct-answer-id", 0, "")
	regrade.Flags().String("token", "", "")
	quizzes.AddCommand(regrade)
	root.AddCommand(quizzes)
	root.SetArgs(append([]string{"quizzes", "regrade"}, args...))
	cmd, err := root.ExecuteContextC(context.Background())
	if err != nil {
		t.Fatalf("execute fake tree: %v", err)
	}
	return cmd
}

func TestLogActivity_WritesEntry(t *testing.T) {
	path := useTempActivityLog(t)
	rec := activity.NewRecorder()
	rec.Record("GET", "/api/v1/courses/1/quizzes/10/questions/789?include[]=x", 200)
	rec.Record("PUT", "/api/v1/courses/1/quizzes/10/questions/789", 200)
	rec.SetDetail("regrade", map[string]interface{}{"changed": 2, "mismatched": 0})

	cmd := fakeExecuted(t, "10", "--course-id", "1", "--question", "789", "--correct-answer-id", "1002", "--token", "7~AbCdEfGhIjKlMnOpQrStUv")
	argv := []string{"quizzes", "regrade", "10", "--course-id", "1", "--question", "789", "--correct-answer-id", "1002", "--token", "7~AbCdEfGhIjKlMnOpQrStUv"}
	logActivity(cmd, errors.New("boom"), time.Now().Add(-1500*time.Millisecond), rec, argv)

	entries, _, err := activity.Read(path)
	if err != nil || len(entries) != 1 {
		t.Fatalf("Read = %d entries, %v", len(entries), err)
	}
	e := entries[0]
	if e.Command != "quizzes regrade" || e.ExitCode != 1 || e.DurationMs < 1500 || e.DurationMs > 10000 {
		t.Errorf("entry = %+v", e)
	}
	if _, err := time.Parse(time.RFC3339, e.Timestamp); err != nil || !strings.HasSuffix(e.Timestamp, "Z") {
		t.Errorf("ts = %q, want RFC3339 UTC", e.Timestamp)
	}
	if strings.Join(e.Args, " ") != "quizzes regrade 10 --course-id 1 --question 789 --correct-answer-id 1002 --token "+activity.Redacted {
		t.Errorf("args = %q", e.Args)
	}
	if len(e.Requests) != 2 || strings.Contains(e.Requests[0].Path, "?") || e.Requests[1].Method != "PUT" {
		t.Errorf("requests = %+v", e.Requests)
	}
	want := "correct-answer:1002 course:1 question:789 quiz:10"
	var got []string
	for _, x := range e.Touched {
		got = append(got, x.Type+":"+strconv.FormatInt(x.ID, 10))
	}
	if strings.Join(got, " ") != want {
		t.Errorf("touched = %v, want %s", got, want)
	}
	// details (per-student score tables) ride only with capture_bodies
	if e.Details != nil {
		t.Errorf("details must not be written without capture_bodies: %+v", e.Details)
	}
	raw, _ := os.ReadFile(path)
	if strings.Contains(string(raw), "7~AbCd") {
		t.Errorf("token leaked into the log: %s", raw)
	}
}

// useTempHome points the config dir at a fresh home (no config file) and
// clears the config cache around the test.
func useTempHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	// os.UserHomeDir reads USERPROFILE on Windows, so setting HOME alone let
	// these tests write into the RUNNER's real profile and then assert against
	// the temp path they never used.
	t.Setenv("USERPROFILE", dir)
	for _, name := range activityEnvVars {
		t.Setenv(name, "")
	}
	config.ResetCache()
	t.Cleanup(config.ResetCache)
	return dir
}

func TestLogActivity_DisabledByDefault(t *testing.T) {
	dir := useTempHome(t) // config dir under the temp home: nothing should be written there either
	cmd := fakeExecuted(t, "10", "--course-id", "1", "--question", "789", "--correct-answer-id", "1002")
	logActivity(cmd, nil, time.Now(), activity.NewRecorder(), []string{"quizzes", "regrade", "10"})
	if _, err := os.Stat(filepath.Join(dir, ".canvas-cli", activity.DefaultFileName)); !os.IsNotExist(err) {
		t.Errorf("log written although disabled")
	}
}

func TestLogActivity_SkipsOwnCommands(t *testing.T) {
	path := useTempActivityLog(t)
	root := &cobra.Command{Use: "canvas"}
	act := &cobra.Command{Use: "activity"}
	list := &cobra.Command{Use: "list", RunE: func(*cobra.Command, []string) error { return nil }}
	act.AddCommand(list)
	root.AddCommand(act)
	root.SetArgs([]string{"activity", "list"})
	cmd, _ := root.ExecuteContextC(context.Background())
	logActivity(cmd, nil, time.Now(), activity.NewRecorder(), []string{"activity", "list"})
	logActivity(nil, nil, time.Now(), activity.NewRecorder(), nil)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("activity's own commands must not be logged")
	}
}

func TestLogActivity_WriteFailureDoesNotFailCommand(t *testing.T) {
	// a path inside a file cannot be created
	blocker := filepath.Join(t.TempDir(), "file")
	_ = os.WriteFile(blocker, []byte("x"), 0o600)
	t.Setenv(activity.EnvVar, filepath.Join(blocker, "activity.jsonl"))
	cmd := fakeExecuted(t, "10", "--course-id", "1", "--question", "789", "--correct-answer-id", "1002")
	// must not panic or return anything; the warning goes to stderr
	logActivity(cmd, nil, time.Now(), activity.NewRecorder(), []string{"quizzes", "regrade", "10"})
}

func TestRunWithActivityLog_ReturnsCommandError(t *testing.T) {
	useTempActivityLog(t)
	want := errors.New("command failed")
	err := runWithActivityLog(func() (*cobra.Command, error) { return nil, want })
	if !errors.Is(err, want) {
		t.Errorf("runWithActivityLog returned %v, want the command's error", err)
	}
}

func TestActivityListCmd(t *testing.T) {
	path := useTempActivityLog(t)
	cfg := activity.Config{Path: path, MaxSizeBytes: 1 << 20}
	old := time.Now().UTC().Add(-48 * time.Hour).Format(time.RFC3339)
	recent := time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339)
	mustAppend := func(e activity.Entry) {
		if _, err := activity.Append(cfg, e); err != nil {
			t.Fatal(err)
		}
	}
	mustAppend(activity.Entry{Timestamp: old, Command: "courses list", Requests: []activity.Request{{Method: "GET", Path: "/api/v1/courses", Status: 200}}})
	mustAppend(activity.Entry{Timestamp: recent, Command: "quizzes regrade", Requests: []activity.Request{{Method: "PUT", Path: "/api/v1/courses/1/quizzes/2/questions/3", Status: 200}}})
	mustAppend(activity.Entry{Timestamp: recent, Command: "quizzes questions get", Requests: []activity.Request{{Method: "GET", Path: "/api/v1/courses/1/quizzes/2/questions/3", Status: 200}}})

	tests := []cmdtest.CommandTestCase{
		{
			Name: "all entries",
			Args: []string{},
			ValidateOutput: func(t *testing.T, output string) {
				for _, want := range []string{"courses list", "quizzes regrade", "quizzes questions get"} {
					if !strings.Contains(output, want) {
						t.Errorf("expected %q in output:\n%s", want, output)
					}
				}
			},
		},
		{
			Name: "since 24h",
			Args: []string{"--since", "24h"},
			ValidateOutput: func(t *testing.T, output string) {
				if strings.Contains(output, "courses list") || !strings.Contains(output, "quizzes regrade") {
					t.Errorf("since filter wrong:\n%s", output)
				}
			},
		},
		{
			Name: "writes only",
			Args: []string{"--writes"},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "quizzes regrade") || strings.Contains(output, "questions get") || strings.Contains(output, "courses list") {
					t.Errorf("writes filter wrong:\n%s", output)
				}
			},
		},
		{
			Name: "command prefix",
			Args: []string{"--command", "quizzes questions"},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "questions get") || strings.Contains(output, "regrade") {
					t.Errorf("prefix filter wrong:\n%s", output)
				}
			},
		},
		{
			Name: "json output",
			Args: []string{"--writes", "json"},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, `"command": "quizzes regrade"`) && !strings.Contains(output, `"command":"quizzes regrade"`) {
					t.Errorf("json output wrong:\n%s", output)
				}
			},
		},
		{
			Name:        "bad since",
			Args:        []string{"--since", "yesterday"},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			// -o is a persistent root flag, unavailable on the bare
			// subcommand: set the global instead and drop the marker arg.
			if len(tc.Args) > 0 && tc.Args[len(tc.Args)-1] == "json" {
				setOutputFormat(t, "json")
				tc.Args = tc.Args[:len(tc.Args)-1]
			}
			cmdtest.RunCommandTest(t, newActivityListCmd(), tc)
		})
	}

	// empty / disabled hint
	t.Setenv(activity.EnvVar, filepath.Join(t.TempDir(), "none.jsonl"))
	cmdtest.RunCommandTest(t, newActivityListCmd(), cmdtest.CommandTestCase{Name: "empty", ExpectOutput: "No activity logged"})
}

func TestActivityClearCmd_RequiresForce(t *testing.T) {
	path := useTempActivityLog(t)
	_ = os.WriteFile(path, []byte("{}\n"), 0o600)

	cmdtest.RunCommandTest(t, newActivityClearCmd(), cmdtest.CommandTestCase{Name: "without force", ExpectError: true})
	if info, _ := os.Stat(path); info == nil || info.Size() == 0 {
		t.Fatal("clear without --force must not touch the log")
	}

	cmdtest.RunCommandTest(t, newActivityClearCmd(), cmdtest.CommandTestCase{Name: "with force", Args: []string{"--force"}, ExpectOutput: "Cleared"})
	if info, _ := os.Stat(path); info == nil || info.Size() != 0 {
		t.Fatal("clear --force must truncate the log")
	}
}

func TestActivityArchiveAndPathCmds(t *testing.T) {
	path := useTempActivityLog(t)
	_ = os.WriteFile(path, []byte("{}\n"), 0o600)

	cmdtest.RunCommandTest(t, newActivityArchiveCmd(), cmdtest.CommandTestCase{Name: "archive", ExpectOutput: "Archived"})
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("archive must move the log")
	}
	matches, _ := filepath.Glob(filepath.Join(filepath.Dir(path), "activity-*.jsonl"))
	if len(matches) != 1 {
		t.Errorf("archive file missing: %v", matches)
	}
	cmdtest.RunCommandTest(t, newActivityArchiveCmd(), cmdtest.CommandTestCase{Name: "archive again (nothing)", ExpectOutput: "No activity log"})

	cmdtest.RunCommandTest(t, newActivityPathCmd(), cmdtest.CommandTestCase{Name: "path", ValidateOutput: func(t *testing.T, output string) {
		if !strings.Contains(output, path) || !strings.Contains(output, "enabled: true") {
			t.Errorf("path output = %s", output)
		}
	}})
	t.Setenv(activity.EnvVar, "")
	cmdtest.RunCommandTest(t, newActivityPathCmd(), cmdtest.CommandTestCase{Name: "path disabled", ExpectOutput: "enabled: false"})
}
