package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// canvasIrreversibleVerbs are the Canvas operations that cannot be undone.
// "canvas agent guard" hard-blocks these by default; ordinary writes (create,
// update, publish, grade, …) only require approval.
var canvasIrreversibleVerbs = map[string]bool{
	"delete":      true,
	"remove":      true,
	"conclude":    true,
	"reset":       true,
	"abort":       true,
	"crosslist":   true,
	"uncrosslist": true,
	"deactivate":  true,
	"unpublish":   true,
	"unlink":      true,
	"clear":       true,
	"void":        true,
}

// canvasWriteVerbs are Canvas operations that mutate data but can be recovered
// from (creates can be deleted, updates can be reverted, etc.). These surface
// as approval-required rather than hard-blocked.
var canvasWriteVerbs = map[string]bool{
	"create":              true,
	"update":              true,
	"publish":             true,
	"grade":               true,
	"bulk-grade":          true,
	"upload":              true,
	"add":                 true,
	"set":                 true,
	"move":                true,
	"duplicate":           true,
	"reply":               true,
	"post":                true,
	"send":                true,
	"enroll":              true,
	"accept":              true,
	"reject":              true,
	"reactivate":          true,
	"star":                true,
	"unstar":              true,
	"archive":             true,
	"unarchive":           true,
	"subscribe":           true,
	"unsubscribe":         true,
	"relock":              true,
	"revert":              true,
	"associate":           true,
	"sync":                true,
	"restore":             true,
	"regenerate-secret":   true,
	"switch-experience":   true,
	"switch-role":         true,
	"import":              true,
	"batch-update":        true,
	"mark-read":           true,
	"mark-all-read":       true,
	"complete":            true,
	"dismiss":             true,
	"done":                true,
	"conclude-enrollment": true,
	"reserve":             true,
	"crosslist-section":   true,
	"emit":                true,
	"link":                true,
	"enable":              true,
	"disable":             true,
}

// canvasLocalGroups are top-level command groups that never call the Canvas
// API — they perform local operations only and must never be gated.
var canvasLocalGroups = map[string]bool{
	"auth":       true,
	"config":     true,
	"context":    true,
	"alias":      true,
	"cache":      true,
	"doctor":     true,
	"completion": true,
	"version":    true,
	"mcp":        true,
	"skills":     true,
	"repl":       true,
	"shell":      true,
	"update":     true,
	"telemetry":  true,
	"webhook":    true,
	"agent":      true,
	"help":       true,
}

// canvasGuardCmd represents one classified Canvas operation.
type canvasGuardCmd struct {
	cli  string // CLI path without root, e.g. "courses delete"
	tool string // MCP tool name, e.g. "canvas_courses_delete"
	verb string // leaf command name, e.g. "delete"
}

func init() {
	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Helpers for running canvas under an AI agent",
		Long:  "Helpers for running canvas under an AI agent (Claude Code, Codex, OpenCode, …).",
	}
	agentCmd.AddCommand(newAgentGuardCmd())
	rootCmd.AddCommand(agentCmd)
}

func newAgentGuardCmd() *cobra.Command {
	var host string
	var allWrites bool
	var write bool

	cmd := &cobra.Command{
		Use:   "guard",
		Short: "Generate agent-safety config that blocks destructive canvas operations",
		Long: `guard generates the permission rules and hooks that stop an AI agent from
running destructive canvas operations, derived from the live command tree so the
list is always complete.

By default it hard-blocks the irreversible actions (delete, remove, conclude,
crosslist, uncrosslist, deactivate, unpublish, unlink, clear, abort, reset)
and makes ordinary writes (create, update, publish, grade, ...) require
approval; read operations stay allowed. Pass --all-writes to block writes too.

IMPORTANT: the "canvas api" escape hatch can issue any HTTP verb. The guard
blocks "canvas api DELETE/PUT/POST" patterns on the Bash surface but cannot
enumerate arbitrary path arguments. For a hard guarantee, run the agent MCP-only
(no Bash tool) or inside a read-only sandbox.

Output is printed for review by default; pass --write to install it. See the
Agent Safety guide: https://jjuanrivvera.github.io/canvas-cli/user-guide/agent-safety/`,
		Args: cobra.NoArgs,
		Example: "  canvas agent guard --host claude-code\n" +
			"  canvas agent guard --host codex\n" +
			"  canvas agent guard --host opencode --all-writes\n" +
			"  canvas agent guard --host claude-code --write",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, writes, irreversible := classifyCanvasCommands(rootCmd)
			g := canvasGuardPlan{
				irreversible: irreversible,
				writes:       writes,
				allWrites:    allWrites,
			}
			switch host {
			case "claude-code", "claude":
				return emitCanvasClaudeCode(cmd, g, write)
			case "codex":
				return emitCanvasCodex(cmd, g, write)
			case "opencode":
				return emitCanvasOpenCode(cmd, g, write)
			case "":
				return fmt.Errorf("--host is required (claude-code, codex, or opencode)")
			default:
				return fmt.Errorf("unknown host %q (use claude-code, codex, or opencode)", host)
			}
		},
	}
	cmd.Flags().StringVar(&host, "host", "", "Target agent host: claude-code, codex, opencode")
	cmd.Flags().BoolVar(&allWrites, "all-writes", false, "Also block create/update/grade (default: those require approval)")
	cmd.Flags().BoolVar(&write, "write", false, "Write the config/hook files instead of printing them")
	_ = cmd.RegisterFlagCompletionFunc("host", func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		vals := []string{"claude-code", "codex", "opencode"}
		out := make([]string, 0, len(vals))
		for _, v := range vals {
			if strings.HasPrefix(v, toComplete) {
				out = append(out, v)
			}
		}
		return out, cobra.ShellCompDirectiveNoFileComp
	})
	return cmd
}

// canvasGuardPlan is the classified set of operations the generators convert
// into host-specific config.
type canvasGuardPlan struct {
	irreversible []canvasGuardCmd
	writes       []canvasGuardCmd
	allWrites    bool
}

// blocked returns the operations to hard-block (irreversible, plus writes when
// --all-writes is set).
func (g canvasGuardPlan) blocked() []canvasGuardCmd {
	if g.allWrites {
		return append(append([]canvasGuardCmd{}, g.irreversible...), g.writes...)
	}
	return g.irreversible
}

// asked returns the operations that require approval (writes, unless
// --all-writes folds them into the hard-block set).
func (g canvasGuardPlan) asked() []canvasGuardCmd {
	if g.allWrites {
		return nil
	}
	return g.writes
}

// blockedVerbs returns the distinct verb strings from the hard-block set.
func (g canvasGuardPlan) blockedVerbs() []string {
	return distinctCanvasVerbs(g.blocked())
}

// isCanvasIrreversibleVerb reports whether a command name contains an
// irreversible verb, handling compound names like "bulk-grade" by splitting on
// "-" (matching alegra's isIrreversibleVerb approach).
func isCanvasIrreversibleVerb(name string) bool {
	for _, tok := range strings.Split(name, "-") {
		if canvasIrreversibleVerbs[tok] {
			return true
		}
	}
	return false
}

// isCanvasWriteVerb reports whether a command name is a write operation.
func isCanvasWriteVerb(name string) bool {
	// Full name first (handles "bulk-grade", "mark-read", etc.)
	if canvasWriteVerbs[name] {
		return true
	}
	// Compound check: any token that is a write verb
	for _, tok := range strings.Split(name, "-") {
		if canvasWriteVerbs[tok] {
			return true
		}
	}
	return false
}

// topLevelGroup returns the top-level command group name for a given command by
// walking up the parent chain.
func topLevelGroup(c *cobra.Command) string {
	// Walk up to find the child of root
	cur := c
	for cur.Parent() != nil && cur.Parent().Parent() != nil {
		cur = cur.Parent()
	}
	if cur.Parent() == nil {
		// c itself is root
		return ""
	}
	return cur.Name()
}

// classifyCanvasCommands walks the command tree and buckets every runnable
// leaf command into read, write, or irreversible. Commands under local/utility
// top-level groups (auth, config, cache, agent, etc.) are excluded entirely —
// they never call the Canvas API.
//
// Unlike alegra-cli, canvas-cli does not set openWorldHint/readOnlyHint cobra
// annotations, so classification is purely by verb (leaf command name).
func classifyCanvasCommands(root *cobra.Command) (read, writes, irreversible []canvasGuardCmd) {
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		for _, sub := range c.Commands() {
			// Recurse first so we always visit leaves.
			walk(sub)

			if !sub.Runnable() || sub.Hidden || sub.Name() == "help" {
				continue
			}

			group := topLevelGroup(sub)

			// Skip the "api" raw-escape command — it is documented separately
			// in the hook script header and cannot be safely classified by verb.
			if group == "api" || sub.Name() == "api" {
				continue
			}

			// Skip local/utility groups that never hit the Canvas API.
			if canvasLocalGroups[group] {
				continue
			}

			gc := canvasGuardCmd{
				cli:  strings.TrimPrefix(sub.CommandPath(), root.Name()+" "),
				tool: strings.ReplaceAll(sub.CommandPath(), " ", "_"),
				verb: sub.Name(),
			}

			switch {
			case isCanvasIrreversibleVerb(sub.Name()):
				irreversible = append(irreversible, gc)
			case isCanvasWriteVerb(sub.Name()):
				writes = append(writes, gc)
			default:
				read = append(read, gc)
			}
		}
	}
	walk(root)
	sortCanvasGuard(read)
	sortCanvasGuard(writes)
	sortCanvasGuard(irreversible)
	return read, writes, irreversible
}

func sortCanvasGuard(cs []canvasGuardCmd) {
	sort.Slice(cs, func(i, j int) bool { return cs[i].tool < cs[j].tool })
}

func distinctCanvasVerbs(cs []canvasGuardCmd) []string {
	seen := map[string]bool{}
	var out []string
	for _, c := range cs {
		if !seen[c.verb] {
			seen[c.verb] = true
			out = append(out, c.verb)
		}
	}
	sort.Strings(out)
	return out
}

// canvasVerbGroup renders verbs as a regex alternation, e.g. "(delete|remove|void)".
func canvasVerbGroup(verbs []string) string {
	return "(" + strings.Join(verbs, "|") + ")"
}

// canvasWriteOrPrint either writes content to path (creating parent dirs) when
// write is set and the file does not already exist, or prints it to the
// command's output with a header. It never overwrites an existing file.
func canvasWriteOrPrint(cmd *cobra.Command, write bool, path, content string, perm os.FileMode) error {
	out := cmd.OutOrStdout()
	if !write {
		fmt.Fprintf(out, "# ----- %s -----\n%s\n", path, content)
		return nil
	}
	if _, err := os.Stat(path); err == nil {
		fmt.Fprintf(out, "# %s already exists — review and merge manually:\n%s\n", path, content)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(content), perm); err != nil { // #nosec G306 -- hook script needs 0o755
		return err
	}
	fmt.Fprintf(out, "wrote %s\n", path)
	return nil
}
