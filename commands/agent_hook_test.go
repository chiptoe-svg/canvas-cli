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

	// Generate the hook from the real command tree so the blocked_cmds /
	// blocked_tools arrays are fully populated.
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

	// helper: build a PreToolUse JSON payload for an MCP tool call.
	mcpPayload := func(toolName string) string {
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
			name:        "bash_courses_delete_denied",
			payload:     bashPayload("canvas courses delete 5"),
			wantDenied:  true,
			description: "direct canvas courses delete must be denied",
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
			name:        "mcp_courses_delete_denied",
			payload:     mcpPayload("mcp__canvas__canvas_courses_delete"),
			wantDenied:  true,
			description: "MCP canvas_courses_delete tool must be denied",
		},
		{
			name:        "mcp_courses_list_allowed",
			payload:     mcpPayload("mcp__canvas__canvas_courses_list"),
			wantDenied:  false,
			description: "MCP canvas_courses_list tool must NOT be denied",
		},
		{
			name: "bash_obfuscated_delete_denied",
			// canvas courses de""lete 5 — quotes stripped by deobfuscate become
			// "canvas courses delete 5"
			payload:     bashPayload(`canvas courses de""lete 5`),
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

	// Create an empty directory to shadow jq so it's "unavailable".
	emptyPath := filepath.Join(tmpDir, "empty")
	if err := os.Mkdir(emptyPath, 0o750); err != nil {
		t.Fatalf("mkdir emptyPath: %v", err)
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

	// Build a PATH that keeps system tools (grep, sed, tr, etc.) but omits
	// jq specifically by placing a wrapper dir that has no jq first.
	// We keep the original PATH for all the POSIX utilities the hook uses.
	origPath := os.Getenv("PATH")
	// emptyPath prepended — jq not present there, but grep/sed/tr found later.
	noJqPath := emptyPath + ":" + origPath

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
			name:       "nojq_courses_delete_denied",
			payload:    bashPayload("canvas courses delete 5"),
			wantDenied: true,
		},
		{
			name:       "nojq_cat_courses_delete_allowed",
			payload:    bashPayload("cat courses_delete.go"),
			wantDenied: false,
		},
		{
			name:       "nojq_obfuscated_delete_denied",
			payload:    bashPayload(`canvas courses de""lete 5`),
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
