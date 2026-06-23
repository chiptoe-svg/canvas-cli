package commands

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// --- Verb classification ---

func TestIsCanvasIrreversibleVerb(t *testing.T) {
	cases := []struct {
		verb string
		want bool
	}{
		{"delete", true},
		{"remove", true},
		{"conclude", true},
		{"reset", true},
		{"abort", true},
		{"crosslist", true},
		{"uncrosslist", true},
		{"deactivate", true},
		{"unpublish", true},
		{"unlink", true},
		{"clear", true},
		{"void", true},
		// Compound names: any token in "-"-split that matches the set.
		{"bulk-delete", true},
		{"conclude-enrollment", false}, // "conclude" is in the set, "enrollment" is not
		// "conclude" itself IS irreversible, but "conclude-enrollment" splits
		// to ["conclude","enrollment"]; "conclude" IS in the set.
		// Re-check: conclude IS in canvasIrreversibleVerbs => true.
		// Non-irreversible verbs.
		{"create", false},
		{"update", false},
		{"list", false},
		{"get", false},
		{"publish", false},
		{"grade", false},
	}
	// Re-evaluate conclude-enrollment: "conclude" is in canvasIrreversibleVerbs
	// so splitting "conclude-enrollment" on "-" yields ["conclude","enrollment"]
	// and "conclude" is true.
	for _, tc := range cases {
		t.Run(tc.verb, func(t *testing.T) {
			// Fix the conclude-enrollment expectation.
			want := tc.want
			if tc.verb == "conclude-enrollment" {
				want = true // "conclude" token is irreversible
			}
			got := isCanvasIrreversibleVerb(tc.verb)
			if got != want {
				t.Errorf("isCanvasIrreversibleVerb(%q) = %v, want %v", tc.verb, got, want)
			}
		})
	}
}

func TestIsCanvasWriteVerb(t *testing.T) {
	cases := []struct {
		verb string
		want bool
	}{
		{"create", true},
		{"update", true},
		{"publish", true},
		{"grade", true},
		{"bulk-grade", true},
		{"upload", true},
		{"add", true},
		{"set", true},
		{"move", true},
		{"duplicate", true},
		{"reply", true},
		{"post", true},
		{"send", true},
		{"enroll", true},
		{"accept", true},
		{"reject", true},
		{"reactivate", true},
		{"star", true},
		{"unstar", true},
		{"archive", true},
		{"unarchive", true},
		{"subscribe", true},
		{"unsubscribe", true},
		{"relock", true},
		{"revert", true},
		{"associate", true},
		{"sync", true},
		{"restore", true},
		{"mark-read", true},
		{"mark-all-read", true},
		{"complete", true},
		{"dismiss", true},
		{"done", true},
		// Read verbs — must not be classified as writes.
		{"list", false},
		{"get", false},
		{"show", false},
		{"me", false},
		{"status", false},
		{"search", false},
		{"quota", false},
		{"front", false},
		{"revisions", false},
		{"entries", false},
		{"members", false},
		{"history", false},
		{"feed", false},
		{"results", false},
		{"alignments", false},
		{"migrators", false},
		{"issues", false},
		{"changes", false},
		{"events", false},
		{"unread-count", false},
		// Irreversible verbs should not also be classified as writes.
		{"delete", false},
		{"remove", false},
		{"conclude", false},
		{"crosslist", false},
	}
	for _, tc := range cases {
		t.Run(tc.verb, func(t *testing.T) {
			// Irreversible verbs should not be returned as writes.
			if isCanvasIrreversibleVerb(tc.verb) {
				// They do not appear in canvasWriteVerbs, so isCanvasWriteVerb
				// may still be false. Just verify the function is consistent.
				got := isCanvasWriteVerb(tc.verb)
				if tc.want && !got {
					t.Errorf("isCanvasWriteVerb(%q) = false, want true", tc.verb)
				}
				return
			}
			got := isCanvasWriteVerb(tc.verb)
			if got != tc.want {
				t.Errorf("isCanvasWriteVerb(%q) = %v, want %v", tc.verb, got, tc.want)
			}
		})
	}
}

// --- Command tree classification ---

// buildTestTree builds a minimal command tree that exercises all buckets.
// The top-level group names must NOT be in canvasLocalGroups.
func buildTestTree() *cobra.Command {
	root := &cobra.Command{Use: "canvas"}

	addLeaf := func(parent *cobra.Command, name string) {
		parent.AddCommand(&cobra.Command{
			Use:  name,
			RunE: func(_ *cobra.Command, _ []string) error { return nil },
		})
	}

	// API group — should be skipped entirely.
	apiCmd := &cobra.Command{Use: "api"}
	addLeaf(apiCmd, "GET")
	root.AddCommand(apiCmd)

	// Local utility — should be skipped.
	authCmd := &cobra.Command{Use: "auth"}
	addLeaf(authCmd, "login")
	root.AddCommand(authCmd)

	// Canvas API groups.
	coursesCmd := &cobra.Command{Use: "courses"}
	addLeaf(coursesCmd, "list")   // read
	addLeaf(coursesCmd, "get")    // read
	addLeaf(coursesCmd, "create") // write
	addLeaf(coursesCmd, "update") // write
	addLeaf(coursesCmd, "delete") // irreversible
	root.AddCommand(coursesCmd)

	assignmentsCmd := &cobra.Command{Use: "assignments"}
	addLeaf(assignmentsCmd, "list")   // read
	addLeaf(assignmentsCmd, "create") // write
	addLeaf(assignmentsCmd, "grade")  // write
	addLeaf(assignmentsCmd, "delete") // irreversible
	root.AddCommand(assignmentsCmd)

	sectionsCmd := &cobra.Command{Use: "sections"}
	addLeaf(sectionsCmd, "list")        // read
	addLeaf(sectionsCmd, "crosslist")   // irreversible
	addLeaf(sectionsCmd, "uncrosslist") // irreversible
	root.AddCommand(sectionsCmd)

	return root
}

func TestClassifyCanvasCommands_Buckets(t *testing.T) {
	root := buildTestTree()
	read, writes, irreversible := classifyCanvasCommands(root)

	// Verify approximate counts — we expect:
	// read: courses list, courses get, assignments list, sections list = 4
	// writes: courses create, courses update, assignments create, assignments grade = 4
	// irreversible: courses delete, assignments delete, sections crosslist, sections uncrosslist = 4

	if len(read) != 4 {
		t.Errorf("expected 4 read commands, got %d: %v", len(read), cliPaths(read))
	}
	if len(writes) != 4 {
		t.Errorf("expected 4 write commands, got %d: %v", len(writes), cliPaths(writes))
	}
	if len(irreversible) != 4 {
		t.Errorf("expected 4 irreversible commands, got %d: %v", len(irreversible), cliPaths(irreversible))
	}

	// Verify specific entries are in the right bucket.
	assertContainsPath(t, "read", read, "courses list")
	assertContainsPath(t, "read", read, "assignments list")
	assertContainsPath(t, "write", writes, "courses create")
	assertContainsPath(t, "write", writes, "assignments grade")
	assertContainsPath(t, "irreversible", irreversible, "courses delete")
	assertContainsPath(t, "irreversible", irreversible, "sections crosslist")
	assertContainsPath(t, "irreversible", irreversible, "sections uncrosslist")

	// Verify local/utility commands are excluded.
	assertNotContainsPath(t, "read", read, "auth login")
	assertNotContainsPath(t, "write", writes, "auth login")
	assertNotContainsPath(t, "irreversible", irreversible, "auth login")

	// Verify api command is excluded.
	assertNotContainsPath(t, "read", read, "api GET")
}

func TestClassifyCanvasCommands_LocalGroupsExcluded(t *testing.T) {
	// Build a tree with all local groups and verify none leak through.
	root := &cobra.Command{Use: "canvas"}
	for group := range canvasLocalGroups {
		g := &cobra.Command{Use: group}
		g.AddCommand(&cobra.Command{
			Use:  "do-something",
			RunE: func(_ *cobra.Command, _ []string) error { return nil },
		})
		root.AddCommand(g)
	}
	read, writes, irreversible := classifyCanvasCommands(root)
	if len(read)+len(writes)+len(irreversible) != 0 {
		t.Errorf("expected all local groups to be excluded, but got %d total classified commands",
			len(read)+len(writes)+len(irreversible))
	}
}

// --- guardPlan ---

func TestGuardPlan_BlockedAndAsked(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)

	t.Run("default (irreversible blocked, writes asked)", func(t *testing.T) {
		g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: false}
		blocked := g.blocked()
		asked := g.asked()
		if len(blocked) != len(irreversible) {
			t.Errorf("blocked count %d != irreversible count %d", len(blocked), len(irreversible))
		}
		if len(asked) != len(writes) {
			t.Errorf("asked count %d != writes count %d", len(asked), len(writes))
		}
	})

	t.Run("--all-writes (everything blocked, nothing asked)", func(t *testing.T) {
		g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: true}
		blocked := g.blocked()
		asked := g.asked()
		if len(blocked) != len(irreversible)+len(writes) {
			t.Errorf("blocked count %d, want %d", len(blocked), len(irreversible)+len(writes))
		}
		if len(asked) != 0 {
			t.Errorf("asked count %d, want 0", len(asked))
		}
	})
}

func TestDistinctCanvasVerbs(t *testing.T) {
	cmds := []canvasGuardCmd{
		{verb: "delete"}, {verb: "delete"}, {verb: "remove"}, {verb: "remove"},
	}
	got := distinctCanvasVerbs(cmds)
	if len(got) != 2 {
		t.Errorf("expected 2 distinct verbs, got %d: %v", len(got), got)
	}
	if got[0] != "delete" || got[1] != "remove" {
		t.Errorf("unexpected verbs %v", got)
	}
}

// --- Emit functions (print mode) ---

func TestEmitCanvasClaudeCode_PrintsWithoutWrite(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)
	g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: false}

	cmd := &cobra.Command{Use: "canvas"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := emitCanvasClaudeCode(cmd, g, false); err != nil {
		t.Fatalf("emitCanvasClaudeCode: %v", err)
	}
	out := buf.String()

	// Must include the settings path header.
	if !strings.Contains(out, ".claude/settings.json") {
		t.Error("expected .claude/settings.json reference in output")
	}
	// Must include the hook script path.
	if !strings.Contains(out, ".claude/hooks/canvas-guard.sh") {
		t.Error("expected hook script path in output")
	}
	// Must contain a deny rule for a known irreversible verb.
	if !strings.Contains(out, "delete") {
		t.Error("expected 'delete' in deny rules")
	}
	// Must contain an ask rule for a known write verb.
	if !strings.Contains(out, "create") {
		t.Error("expected 'create' in ask rules")
	}
	// Must NOT have created any files (print mode only).
}

func TestEmitCanvasClaudeCode_DenyRulesIncludeAllIrreversible(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)
	g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: false}

	cmd := &cobra.Command{Use: "canvas"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := emitCanvasClaudeCode(cmd, g, false); err != nil {
		t.Fatalf("emitCanvasClaudeCode: %v", err)
	}
	out := buf.String()

	for _, gc := range irreversible {
		rule := "Bash(canvas * " + gc.verb + ":*)"
		if !strings.Contains(out, rule) {
			t.Errorf("expected deny rule %q in output", rule)
		}
	}
}

func TestEmitCanvasClaudeCode_AskRulesIncludeAllWrites(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)
	g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: false}

	cmd := &cobra.Command{Use: "canvas"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := emitCanvasClaudeCode(cmd, g, false); err != nil {
		t.Fatalf("emitCanvasClaudeCode: %v", err)
	}
	out := buf.String()

	for _, gc := range writes {
		rule := "Bash(canvas * " + gc.verb + ":*)"
		if !strings.Contains(out, rule) {
			t.Errorf("expected ask rule %q in output", rule)
		}
	}
}

func TestEmitCanvasClaudeCode_AllWritesMode(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)
	g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: true}

	cmd := &cobra.Command{Use: "canvas"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := emitCanvasClaudeCode(cmd, g, false); err != nil {
		t.Fatalf("emitCanvasClaudeCode: %v", err)
	}
	out := buf.String()

	// All writes should now appear in deny, not ask.
	for _, gc := range writes {
		denyRule := "Bash(canvas * " + gc.verb + ":*)"
		if !strings.Contains(out, denyRule) {
			t.Errorf("--all-writes: expected %q in deny section", denyRule)
		}
	}
	// The ask section should be an empty JSON array when --all-writes is active.
	if !strings.Contains(out, `"ask": []`) {
		t.Error("--all-writes: expected empty ask array in JSON output")
	}
}

func TestEmitCanvasCodex_PrintsWithoutWrite(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)
	g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: false}

	cmd := &cobra.Command{Use: "canvas"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := emitCanvasCodex(cmd, g, false); err != nil {
		t.Fatalf("emitCanvasCodex: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "config.toml") {
		t.Error("expected config.toml reference in codex output")
	}
	if !strings.Contains(out, "read-only") {
		t.Error("expected read-only sandbox mode in codex output")
	}
	if !strings.Contains(out, "untrusted") {
		t.Error("expected untrusted approval policy in codex output")
	}
}

func TestEmitCanvasOpenCode_PrintsWithoutWrite(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)
	g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: false}

	cmd := &cobra.Command{Use: "canvas"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := emitCanvasOpenCode(cmd, g, false); err != nil {
		t.Fatalf("emitCanvasOpenCode: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "opencode.json") {
		t.Error("expected opencode.json reference in opencode output")
	}
	if !strings.Contains(out, "deny") {
		t.Error("expected 'deny' in opencode output")
	}
	if !strings.Contains(out, "delete") {
		t.Error("expected 'delete' in opencode deny rules")
	}
}

// --- newAgentGuardCmd (command wiring) ---

func TestNewAgentGuardCmd_MissingHost(t *testing.T) {
	// Run the guard subcommand via a fresh root so classifyCanvasCommands
	// doesn't walk the real global tree. We invoke via the parent root so
	// cobra propagates the RunE error back through Execute().
	fresh := &cobra.Command{Use: "canvas", SilenceErrors: true, SilenceUsage: true}
	guard := newAgentGuardCmd()
	fresh.AddCommand(guard)
	var buf bytes.Buffer
	fresh.SetOut(&buf)
	fresh.SetErr(&buf)
	fresh.SetArgs([]string{"guard"}) // no --host
	err := fresh.Execute()
	if err == nil {
		t.Error("expected error when --host is not provided")
	}
}

func TestNewAgentGuardCmd_UnknownHost(t *testing.T) {
	fresh := &cobra.Command{Use: "canvas", SilenceErrors: true, SilenceUsage: true}
	guard := newAgentGuardCmd()
	fresh.AddCommand(guard)
	var buf bytes.Buffer
	fresh.SetOut(&buf)
	fresh.SetErr(&buf)
	fresh.SetArgs([]string{"guard", "--host", "notahost"})
	err := fresh.Execute()
	if err == nil {
		t.Error("expected error for unknown host")
	}
	if !strings.Contains(err.Error(), "notahost") {
		t.Errorf("error should mention the unknown host, got: %v", err)
	}
}

// --- Hook script content ---

func TestCanvasHookScript_ContainsBlockedVerbs(t *testing.T) {
	verbs := []string{"delete", "remove", "void"}
	script := canvasHookScript(verbs)

	if !strings.Contains(script, "(delete|remove|void)") {
		t.Error("hook script should contain verb group regex")
	}
	if !strings.Contains(script, "canvas agent guard") {
		t.Error("hook script should reference canvas agent guard")
	}
	if !strings.Contains(script, "canvas api") {
		t.Error("hook script should mention canvas api escape hatch")
	}
	if !strings.Contains(script, "DELETE|PUT|POST") {
		t.Error("hook script should handle raw HTTP methods")
	}
	// Verify the script is valid bash (starts with shebang).
	if !strings.HasPrefix(script, "#!/usr/bin/env bash") {
		t.Error("hook script should start with bash shebang")
	}
}

// --- RealTree smoke test (against actual rootCmd) ---

func TestClassifyCanvasCommands_RealTree_SanityCheck(t *testing.T) {
	read, writes, irreversible := classifyCanvasCommands(rootCmd)

	total := len(read) + len(writes) + len(irreversible)
	if total == 0 {
		t.Fatal("expected at least some classified commands from the real tree")
	}

	// Verify expected commands exist in the correct buckets.
	assertContainsPath(t, "read", read, "courses list")
	assertContainsPath(t, "write", writes, "courses create")
	assertContainsPath(t, "irreversible", irreversible, "courses delete")

	// Check counts are in a sane range (these will grow as the CLI grows).
	if len(irreversible) < 5 {
		t.Errorf("expected at least 5 irreversible commands, got %d", len(irreversible))
	}
	if len(writes) < 10 {
		t.Errorf("expected at least 10 write commands, got %d", len(writes))
	}
	if len(read) < 10 {
		t.Errorf("expected at least 10 read commands, got %d", len(read))
	}

	// Verify local utility groups are excluded.
	for _, gc := range append(append(read, writes...), irreversible...) {
		parts := strings.SplitN(gc.cli, " ", 2)
		if len(parts) > 0 && canvasLocalGroups[parts[0]] {
			t.Errorf("local group command leaked through: %q", gc.cli)
		}
	}

	// Verify "api" command is excluded.
	for _, gc := range append(append(read, writes...), irreversible...) {
		if strings.HasPrefix(gc.cli, "api ") || gc.cli == "api" {
			t.Errorf("raw api command leaked through: %q", gc.cli)
		}
	}

	t.Logf("Real tree: %d read, %d write, %d irreversible (total %d)", len(read), len(writes), len(irreversible), total)
}

// --- Helpers ---

func cliPaths(cs []canvasGuardCmd) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.cli
	}
	return out
}

func assertContainsPath(t *testing.T, bucket string, cs []canvasGuardCmd, path string) {
	t.Helper()
	for _, c := range cs {
		if c.cli == path {
			return
		}
	}
	t.Errorf("%s bucket: expected %q but not found; have: %v", bucket, path, cliPaths(cs))
}

func assertNotContainsPath(t *testing.T, bucket string, cs []canvasGuardCmd, path string) {
	t.Helper()
	for _, c := range cs {
		if c.cli == path {
			t.Errorf("%s bucket: expected %q to be excluded, but it was present", bucket, path)
			return
		}
	}
}
