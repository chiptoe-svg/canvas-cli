# Agent Safety

Canvas CLI can be driven by AI agents (Claude Code, Codex, OpenCode, etc.) through the MCP server or the Bash tool. Agents that can issue destructive commands — `courses delete`, `sections crosslist`, `enrollments conclude` — present a risk of irreversible data loss if they are confused by bad instructions or adversarial prompts.

`canvas agent guard` generates a set of permission rules and a PreToolUse hook script tailored to your agent host, derived from the live command tree so the list is always complete.

## Quick start

```bash
# Print the config for review (default: no files are written)
canvas agent guard --host claude-code

# Write config and hook script into the current project
canvas agent guard --host claude-code --write

# Also block write operations (create, update, grade, …) — strictest mode
canvas agent guard --host claude-code --all-writes --write
```

Supported hosts: `claude-code`, `codex`, `opencode`.

## The hard-block vs ask model

By default the guard generates two tiers of protection:

| Tier | Verbs | What happens |
|------|-------|--------------|
| **Hard-block** | `delete`, `remove`, `conclude`, `crosslist`, `uncrosslist`, `deactivate`, `unpublish`, `unlink`, `clear`, `abort`, `reset`, `void` | The PreToolUse hook (or permission rule) denies the tool call immediately, before execution. The agent sees an error and cannot proceed. |
| **Ask / approval** | `create`, `update`, `publish`, `grade`, `upload`, `add`, `sync`, … | The agent must receive explicit human approval before the tool call executes. |
| **Read — allowed** | `list`, `get`, `show`, `search`, `me`, `status`, `quota`, `revisions`, `entries`, `results`, … | No gate — the agent can read freely. |

Pass `--all-writes` to fold the approval tier into the hard-block set, leaving only reads ungated.

## How classification works

The guard classifies every leaf command in the live command tree by its verb (the last path segment, e.g. `delete` in `canvas courses delete`). Compound command names like `delete-comment` are split on `-` and any token matching an irreversible verb triggers the hard-block. Commands under local/utility top-level groups (`auth`, `config`, `cache`, `agent`, `telemetry`, etc.) are excluded entirely — they do not call the Canvas API and must never be gated.

**Classification is fail-safe.** A command is only left _allowed_ if its verb is on an explicit read allowlist (`get`, `list`, `show`, `me`, `search`, …). Irreversible verbs are hard-blocked; **everything else — ordinary writes _and_ any verb the guard doesn't recognize — defaults to "require approval."** So a future command, or a non-obvious mutating verb like `merge`, `cancel`, `bind`, `copy`, or `assign-members`, is gated by default rather than slipping through as a read. Because the guard derives this from the live command tree, regenerating it after an upgrade automatically covers commands added since.

After upgrading `canvas`, regenerate the guard config to capture any new commands:

```bash
canvas agent guard --host claude-code --write
```

## The `canvas api` escape hatch

The `canvas api <METHOD> <PATH>` command bypasses verb classification because it can issue any HTTP method against any endpoint. The guard partially handles this on the Bash surface:

- Bash patterns `Bash(canvas api DELETE:*)`, `Bash(canvas api PUT:*)`, `Bash(canvas api POST:*)` are added to the deny list.
- The hook script also scans for `canvas api ... (DELETE|PUT|POST)` patterns with obfuscation stripping.

However, these patterns are best-effort. Variable indirection (`M=DELETE; canvas api $M /api/v1/...`) and shell aliases can bypass them. For a hard guarantee on the `api` command, use MCP-only mode or a read-only sandbox.

## Obfuscation caveats

The Bash hook defeats trivial quote/backslash obfuscation (`canvas courses de""lete 1`, `canvas courses delete\ 1`) by stripping those characters before pattern matching. It does **not** defeat:

- **Variable indirection**: `v=delete; canvas courses $v 1`
- **Shell aliases**: `alias rm='canvas courses delete'`
- **Process substitution or eval**: `eval "canvas $(cat malicious.txt)"`

The MCP-tool branch has no equivalent weakness: structured tool names (`canvas_courses_delete`) cannot be obfuscated by definition.

## Strongest guarantee: MCP-only operation

If your agent host supports it, configure Claude (or another host) to use only the MCP server and disable the Bash tool entirely. The MCP tool names are deterministic (`canvas_<resource>_<verb>`) and the hook's MCP branch is a hard block with no obfuscation surface. Combined with `--all-writes`, this is the safest configuration.

## Per-host config

### Claude Code

`canvas agent guard --host claude-code --write` writes two files:

- `.claude/settings.json` — `permissions.deny` (irreversible verbs + `canvas api DELETE/PUT/POST`) and `permissions.ask` (writes), plus `hooks.PreToolUse` wiring.
- `.claude/hooks/canvas-guard.sh` — the hook script that intercepts `Bash` and `mcp__*canvas*` tool calls.

The settings file is belt-and-suspenders: the hook catches calls that slip past the permission rules, and vice versa.

!!! warning "Existing settings.json"
    `--write` never overwrites an existing `.claude/settings.json`. If one already exists, the generated JSON is printed so you can merge the `permissions` and `hooks` sections manually.

### Codex

`canvas agent guard --host codex` recommends:

```toml
sandbox_mode    = "read-only"
approval_policy = "untrusted"
```

Codex has no per-command deny hook; the read-only sandbox is the hard block.

### OpenCode

`canvas agent guard --host opencode --write` writes `opencode.json` with `permission.bash` patterns (`deny`/`ask`) and MCP tool patterns.

## Example output (claude-code)

```json
{
  "permissions": {
    "deny": [
      "Bash(canvas * delete:*)",
      "Bash(canvas * remove:*)",
      "Bash(canvas * conclude:*)",
      "Bash(canvas api DELETE:*)",
      "Bash(canvas api PUT:*)",
      "Bash(canvas api POST:*)",
      "mcp__.*canvas.*_(delete|remove|conclude|...)"
    ],
    "ask": [
      "Bash(canvas * create:*)",
      "Bash(canvas * update:*)",
      "Bash(canvas * grade:*)",
      "mcp__.*canvas.*_(create|update|grade|...)"
    ]
  },
  "hooks": {
    "PreToolUse": [
      {"matcher": "Bash", "hooks": [{"type": "command", "command": "${CLAUDE_PROJECT_DIR}/.claude/hooks/canvas-guard.sh"}]},
      {"matcher": "mcp__.*canvas.*", "hooks": [{"type": "command", "command": "${CLAUDE_PROJECT_DIR}/.claude/hooks/canvas-guard.sh"}]}
    ]
  }
}
```
