package commands

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestCanvasHookScript_BashExecution exercises the generated hook script with
// real bash to verify the adversarial cases specified in the code review.
// Gated on runtime.GOOS != "windows" and bash being available, so it is safe
// to include in the regular test suite.
func TestCanvasHookScript_BashExecution(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("bash hook tests require a POSIX shell; skipping on windows")
	}
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found in PATH; skipping hook execution tests")
	}

	// Generate the hook from the real command tree so the blocked_cmds array
	// is fully populated.
	_, _, irreversible := classifyCanvasCommands(rootCmd)

	// Write the hook to a temp file.
	hookContent := canvasHookScript(irreversible)
	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "canvas-guard.sh")
	if err := os.WriteFile(hookFile, []byte(hookContent), 0o755); err != nil {
		t.Fatalf("write hook: %v", err)
	}

	// helper: build a PreToolUse JSON payload for a Bash tool call.
	bashPayload := func(command string) string {
		b, _ := json.Marshal(map[string]any{
			"tool_name": "Bash",
			"tool_input": map[string]any{
				"command": command,
			},
		})
		return string(b)
	}

	// helper: build a PreToolUse JSON payload for a non-Bash tool call.
	otherToolPayload := func(toolName string) string {
		b, _ := json.Marshal(map[string]any{
			"tool_name":  toolName,
			"tool_input": map[string]any{},
		})
		return string(b)
	}

	// runHook executes the hook with the given payload and returns the output.
	runHook := func(t *testing.T, payload string) string {
		t.Helper()
		cmd := exec.Command(bash, hookFile)
		cmd.Stdin = strings.NewReader(payload)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		// Exit code 0 is expected in all cases (the hook always exits 0;
		// the decision is in the JSON output, not the exit code).
		if err := cmd.Run(); err != nil {
			t.Logf("hook output: %s", out.String())
			t.Fatalf("hook script exited non-zero: %v", err)
		}
		return out.String()
	}

	// isDenied checks whether the hook output contains a deny decision.
	isDenied := func(output string) bool {
		return strings.Contains(output, `"permissionDecision":"deny"`)
	}

	cases := []struct {
		name        string
		payload     string
		wantDenied  bool
		description string
	}{
		{
			name:        "bash_assignments_delete_denied",
			payload:     bashPayload("canvas assignments delete 5"),
			wantDenied:  true,
			description: "direct canvas assignments delete must be denied",
		},
		{
			name:        "bash_pages_create_with_delete_title_allowed",
			payload:     bashPayload(`canvas pages create --title "How to remove a student"`),
			wantDenied:  false,
			description: "pages create with benign 'remove' in title arg must NOT be denied",
		},
		{
			name:        "bash_api_patch_denied",
			payload:     bashPayload("canvas api PATCH /api/v1/courses/1 -d x"),
			wantDenied:  true,
			description: "canvas api PATCH must be denied",
		},
		{
			name:        "bash_api_get_with_posts_in_path_allowed",
			payload:     bashPayload("canvas api GET /api/v1/courses/1/discussion_topics/2/posts"),
			wantDenied:  false,
			description: "canvas api GET whose PATH contains 'posts' must NOT be denied",
		},
		{
			name:        "bash_unrelated_file_allowed",
			payload:     bashPayload("cat courses_delete.go"),
			wantDenied:  false,
			description: "non-canvas Bash command mentioning 'delete' must NOT be denied",
		},
		{
			// Bash is the only surface this edition gates: with no MCP server
			// there are no canvas tool names, and the hook must let every
			// other tool through rather than guessing at one.
			name:        "non_bash_tool_allowed",
			payload:     otherToolPayload("Read"),
			wantDenied:  false,
			description: "a non-Bash tool call must NOT be denied",
		},
		{
			name: "bash_obfuscated_delete_denied",
			// canvas assignments de""lete 5 — quotes stripped by deobfuscate become
			// "canvas assignments delete 5"
			payload:     bashPayload(`canvas assignments de""lete 5`),
			wantDenied:  true,
			description: "obfuscated canvas courses de\"\"lete 5 must be denied (jq path)",
		},
		{
			name:        "bash_api_delete_denied",
			payload:     bashPayload("canvas api DELETE /api/v1/courses/99"),
			wantDenied:  true,
			description: "canvas api DELETE must be denied",
		},
		{
			name:        "bash_api_put_denied",
			payload:     bashPayload("canvas api PUT /api/v1/courses/1/assignments/2"),
			wantDenied:  true,
			description: "canvas api PUT must be denied",
		},
		{
			name:        "bash_api_post_denied",
			payload:     bashPayload("canvas api POST /api/v1/courses/1/assignments"),
			wantDenied:  true,
			description: "canvas api POST must be denied",
		},
		{
			name:        "bash_submissions_delete_comment_denied",
			payload:     bashPayload("canvas submissions delete-comment --course-id 1 --assignment-id 2 --submission-id 3 --id 4"),
			wantDenied:  true,
			description: "canvas submissions delete-comment must be denied (compound path with delete token)",
		},
		{
			name:        "bash_relative_path_binary_denied",
			payload:     bashPayload("./bin/canvas assignments delete 5"),
			wantDenied:  true,
			description: "path-invoked binary ./bin/canvas assignments delete must be denied",
		},
		{
			name:        "bash_absolute_path_binary_denied",
			payload:     bashPayload("/usr/local/bin/canvas assignments delete 5"),
			wantDenied:  true,
			description: "path-invoked binary /usr/local/bin/canvas assignments delete must be denied",
		},
		{
			name:        "bash_absolute_path_api_delete_denied",
			payload:     bashPayload("/usr/local/bin/canvas api DELETE /api/v1/courses/5"),
			wantDenied:  true,
			description: "path-invoked canvas api DELETE must be denied",
		},
		{
			name:        "bash_other_binary_named_like_canvas_allowed",
			payload:     bashPayload("mycanvas assignments delete 5"),
			wantDenied:  false,
			description: "a different binary whose name merely ends in 'canvas' must NOT be denied",
		},
		{
			name:        "bash_glued_separator_denied",
			payload:     bashPayload("canvas content-shares delete;true"),
			wantDenied:  true,
			description: "no-arg irreversible command with a shell separator glued to the verb must be denied",
		},
		{
			name:        "bash_glued_pipe_denied",
			payload:     bashPayload("canvas folders delete|cat"),
			wantDenied:  true,
			description: "irreversible command with a pipe glued to the verb must be denied",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			output := runHook(t, tc.payload)
			denied := isDenied(output)
			if denied != tc.wantDenied {
				action := "allowed"
				if tc.wantDenied {
					action = "denied"
				}
				t.Errorf("%s: want %s, got %s\noutput: %s",
					tc.description, action,
					map[bool]string{true: "denied", false: "allowed"}[denied],
					output)
			}
		})
	}
}

// TestCanvasHookScript_BashExecutionNoJq exercises the no-jq fallback path by
// temporarily making jq unavailable via PATH manipulation.
func TestCanvasHookScript_BashExecutionNoJq(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("bash hook tests require a POSIX shell; skipping on windows")
	}
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found in PATH; skipping hook execution tests")
	}

	_, _, irreversible := classifyCanvasCommands(rootCmd)
	hookContent := canvasHookScript(irreversible)
	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "canvas-guard.sh")
	if err := os.WriteFile(hookFile, []byte(hookContent), 0o755); err != nil {
		t.Fatalf("write hook: %v", err)
	}

	// Build a strict PATH containing ONLY the tools the hook needs, minus jq.
	// (Merely prepending an empty dir does NOT hide jq — it stays reachable
	// later in PATH and the no-jq branch never runs; that flaw previously
	// masked a fail-open bug in this branch.)
	strictPath := filepath.Join(tmpDir, "strictbin")
	if err := os.Mkdir(strictPath, 0o750); err != nil {
		t.Fatalf("mkdir strictPath: %v", err)
	}
	for _, tool := range []string{"cat", "tr", "grep", "sed", "printf"} {
		p, err := exec.LookPath(tool)
		if err != nil {
			t.Skipf("required tool %s not found: %v", tool, err)
		}
		if err := os.Symlink(p, filepath.Join(strictPath, tool)); err != nil {
			t.Fatalf("symlink %s: %v", tool, err)
		}
	}

	bashPayload := func(command string) string {
		b, _ := json.Marshal(map[string]any{
			"tool_name": "Bash",
			"tool_input": map[string]any{
				"command": command,
			},
		})
		return string(b)
	}

	noJqPath := strictPath

	runHookNoJq := func(t *testing.T, payload string) string {
		t.Helper()
		cmd := exec.Command(bash, hookFile)
		cmd.Stdin = strings.NewReader(payload)
		// Replace just PATH; keep all other env vars intact so the hook's
		// system utilities (tr, grep, sed, printf) remain available.
		env := make([]string, 0, len(os.Environ()))
		for _, e := range os.Environ() {
			if !strings.HasPrefix(e, "PATH=") {
				env = append(env, e)
			}
		}
		env = append(env, "PATH="+noJqPath)
		cmd.Env = env
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			t.Logf("hook output: %s", out.String())
			t.Fatalf("hook script exited non-zero: %v", err)
		}
		return out.String()
	}

	isDenied := func(output string) bool {
		return strings.Contains(output, `"permissionDecision":"deny"`)
	}

	cases := []struct {
		name       string
		payload    string
		wantDenied bool
	}{
		{
			name:       "nojq_assignments_delete_denied",
			payload:    bashPayload("canvas assignments delete 5"),
			wantDenied: true,
		},
		{
			name:       "nojq_cat_courses_delete_allowed",
			payload:    bashPayload("cat courses_delete.go"),
			wantDenied: false,
		},
		{
			name:       "nojq_obfuscated_delete_denied",
			payload:    bashPayload(`canvas assignments de""lete 5`),
			wantDenied: true,
		},
		{
			name:       "nojq_pages_create_allowed",
			payload:    bashPayload(`canvas pages create --title "How to delete a student"`),
			wantDenied: false,
		},
		{
			name:       "nojq_api_delete_denied",
			payload:    bashPayload("canvas api DELETE /api/v1/courses/5"),
			wantDenied: true,
		},
		{
			// Regression: the compact JSON payload glues the command to its
			// key ("command":"canvas ...). Without JSON-punctuation
			// flattening the anchor can never match and the branch is
			// silently fail-open.
			name:       "nojq_glued_json_key_denied",
			payload:    bashPayload("canvas enrollments conclude 1"),
			wantDenied: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			output := runHookNoJq(t, tc.payload)
			denied := isDenied(output)
			if denied != tc.wantDenied {
				t.Errorf("no-jq case %q: want denied=%v, got denied=%v\noutput: %s",
					tc.name, tc.wantDenied, denied, output)
			}
		})
	}
}
