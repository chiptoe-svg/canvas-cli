package commands

import (
	"strconv"

	"github.com/njayp/ophis"
	"github.com/spf13/cobra"
)

// canvasLocalReadPaths is a hand-maintained allowlist of CLI paths under
// canvasLocalGroups that only read state. Those groups never touch the Canvas
// API, so classifyCanvasCommand returns canvasClassLocal for them and the
// verb-based read allowlist does not apply — "config account" and "context set"
// both end in tokens that look inert but write to local config.
//
// The value is the command's openWorldHint: true when it reaches the network.
//
// Misses are safe: an unlisted command emits no annotations and is therefore
// dropped by read-only MCP clients, never silently exposed as a read.
var canvasLocalReadPaths = map[string]bool{
	"activity list":    false,
	"activity path":    false,
	"alias list":       false,
	"auth status":      false,
	"cache stats":      false,
	"config list":      false,
	"config show":      false,
	"context show":     false,
	"doctor":           true, // probes Canvas connectivity
	"skills path":      false,
	"skills print":     false,
	"telemetry show":   false,
	"telemetry status": false,
	"update check":     true, // queries the GitHub releases API
	"update status":    false,
	"version":          false,
	"webhook events":   false, // prints the static event-type catalog
}

// applyMCPAnnotations stamps MCP tool annotations onto every command in the
// tree, derived from classifyCanvasCommand so the hints can never contradict
// what "canvas agent guard" gates.
//
// This exists because MCP clients running in read-only mode allow a tool only
// when annotations.readOnlyHint is strictly true. A missing annotation counts
// as a write, not as unknown, so without this every tool is filtered out and
// the whole server is dropped. See issue #58.
//
// Must run before the "mcp" subcommand executes — ophis builds its tool list
// from cmd.Annotations at run time, in registerTools. Execute calls it.
//
// Deliberately unannotated (and therefore invisible to read-only clients):
//   - "canvas api", which can issue any HTTP verb
//   - every local command not in canvasLocalReadPaths
func applyMCPAnnotations(root *cobra.Command) {
	walkCanvasCommands(root, func(sub *cobra.Command) {
		class, gc := classifyCanvasCommand(root, sub)
		switch class {
		case canvasClassRead:
			setMCPAnnotations(sub, true, false, true)
		case canvasClassWrite:
			// destructiveHint is left unset: MCP clients default it to true,
			// which is the conservative reading for an update that overwrites
			// fields. Only the irreversible bucket asserts it outright.
			setMCPAnnotations(sub, false, false, true)
		case canvasClassIrreversible:
			setMCPAnnotations(sub, false, true, true)
		case canvasClassLocal:
			if openWorld, ok := canvasLocalReadPaths[gc.cli]; ok {
				setMCPAnnotations(sub, true, false, openWorld)
			}
		case canvasClassSkip:
			// Never annotated — see the doc comment.
		}
	})
}

// setMCPAnnotations merges MCP hints into cmd.Annotations, preserving any
// unrelated keys cobra or the command itself already set.
//
// destructiveHint is only ever written when true; omitting it lets the client
// apply the spec default (true when readOnlyHint is false), which fails safe.
func setMCPAnnotations(cmd *cobra.Command, readOnly, destructive, openWorld bool) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[ophis.AnnotationReadOnly] = strconv.FormatBool(readOnly)
	cmd.Annotations[ophis.AnnotationOpenWorld] = strconv.FormatBool(openWorld)
	if destructive {
		cmd.Annotations[ophis.AnnotationDestructive] = "true"
	}
}
