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
| **Hard-block** | `delete`, `remove`, `conclude`, `crosslist`, `uncrosslist`, `deactivate`, `unpublish`, `unlink`, `clear`, `abort`, `reset`, `void`, `cancel`, `close`, `merge`, `split` | The PreToolUse hook (or permission rule) denies the tool call immediately, before execution. The agent sees an error and cannot proceed. |
| **Ask / approval** | `create`, `update`, `publish`, `grade`, `upload`, `add`, `sync`, … | The agent must receive explicit human approval before the tool call executes. |
| **Read — allowed** | `list`, `get`, `show`, `search`, `me`, `quota`, `revisions`, `entries`, `results`, … | No gate — the agent can read freely. |

Pass `--all-writes` to fold the approval tier into the hard-block set, leaving only reads ungated.

**Special case — bulk-destructive creates:** Some `create` commands that can destroy large amounts of data in a single call are hard-blocked regardless of their verb. `sis-imports create` initiates a full institutional data import that can conclude or delete many enrollments and courses at once and cannot be rolled back. It is always in the hard-block tier.

## How classification works

The guard classifies every leaf command in the live command tree by its verb (the last path segment, e.g. `delete` in `canvas courses delete`). Compound command names like `delete-comment` are split on `-` and any token matching an irreversible verb triggers the hard-block. Commands under local/utility top-level groups (`auth`, `config`, `cache`, `agent`, `telemetry`, etc.) are excluded entirely — they do not call the Canvas API and must never be gated.

**Classification is fail-safe.** A command is only left _allowed_ if its verb is on an explicit read allowlist (`get`, `list`, `show`, `me`, `search`, …). Irreversible verbs are hard-blocked; **everything else — ordinary writes _and_ any verb the guard doesn't recognize — defaults to "require approval."** So a future command, or a non-obvious mutating verb like `bind`, `copy`, or `assign-members`, is gated by default rather than slipping through as a read. Because the guard derives this from the live command tree, regenerating it after an upgrade automatically covers commands added since.

**Hard-blocked verbs include `merge`, `cancel`, `close`, and `split`** — these are irreversible and in the hard-block tier, not the approval tier.

After upgrading `canvas`, regenerate the guard config to capture any new commands:

```bash
canvas agent guard --host claude-code --write
```

## The `canvas api` escape hatch

The `canvas api <METHOD> <PATH>` command bypasses verb classification because it can issue any HTTP method against any endpoint. The guard partially handles this on the Bash surface:

- Bash patterns `Bash(canvas api DELETE:*)`, `Bash(canvas api PUT:*)`, `Bash(canvas api POST:*)`, and `Bash(canvas api PATCH:*)` are added to the deny list.
- The hook script also matches `canvas api (DELETE|PUT|POST|PATCH)` at the method position, so a `GET` whose URL path happens to contain `posts` or `delete` is NOT blocked.

However, these patterns are best-effort. Variable indirection (`M=DELETE; canvas api $M /api/v1/...`) and shell aliases can bypass them. For a hard guarantee on the `api` command, use MCP-only mode or a read-only sandbox.

## Obfuscation caveats

The Bash hook defeats trivial quote/backslash obfuscation (`canvas courses de""lete 1`, `canvas courses delete\ 1`) by stripping those characters before pattern matching, and matches path-invoked binaries (`./bin/canvas courses delete 1`, `/usr/local/bin/canvas courses delete 1`) as well as the bare `canvas` name. It does **not** defeat:

- **Variable indirection**: `v=delete; canvas courses $v 1`
- **Shell aliases**: `alias rm='canvas courses delete'`
- **Process substitution or eval**: `eval "canvas $(cat malicious.txt)"`

The MCP-tool branch has no equivalent weakness: structured tool names (`mcp__canvas__canvas_courses_delete`) cannot be obfuscated by definition.

## Strongest guarantee: MCP-only operation

If your agent host supports it, configure Claude (or another host) to use only the MCP server and disable the Bash tool entirely. The MCP tool names are deterministic and the hook's MCP branch is an exact set-membership check with no obfuscation surface. Combined with `--all-writes`, this is the safest configuration.

## Per-host config

### Claude Code

`canvas agent guard --host claude-code --write` writes two files under the project root (located by walking up to the nearest `.git` directory):

- `.claude/settings.json` — `permissions.deny` (one exact rule per irreversible command path, e.g. `Bash(canvas courses delete:*)`, plus `canvas api DELETE/PUT/POST/PATCH`) and `permissions.ask` (one exact rule per write command), plus `hooks.PreToolUse` wiring.
- `.claude/hooks/canvas-guard.sh` — the hook script that intercepts `Bash` and `mcp__canvas__` tool calls.

The permission rules use exact prefix patterns — `Bash(canvas courses delete:*)` matches only that specific subcommand and cannot false-match a command that has "delete" in an argument. MCP rules are exact tool names like `mcp__canvas__canvas_courses_delete`, not regex patterns.

The settings file is belt-and-suspenders: the hook catches calls that slip past the permission rules, and vice versa.

!!! warning "Existing settings.json"
    `--write` never overwrites an existing `.claude/settings.json`. If one already exists, the generated JSON is printed so you can merge the `permissions` and `hooks` sections manually.

### Codex

`canvas agent guard --host codex` recommends:

```toml
sandbox_mode    = "read-only"
approval_policy = "untrusted"
```

**Important:** Codex sandboxes restrict filesystem access and command execution, but **not network access**. This means Codex cannot hard-block Canvas API calls the way a Claude Code `PreToolUse` hook can — it can only gate them via `approval_policy`. Setting `approval_policy = "untrusted"` causes Codex to pause and request human approval before any shell command or MCP tool call, but a human could still mistakenly approve a destructive action. For a hard block on Canvas destructive actions, use `--host claude-code`.

### OpenCode

`canvas agent guard --host opencode --write` writes `opencode.json` with exact per-command `permission.bash` patterns (`deny`/`ask`) and exact MCP tool names.

## Example output (claude-code)

```json
{
  "permissions": {
    "deny": [
      "Bash(canvas courses delete:*)",
      "Bash(canvas sections crosslist:*)",
      "Bash(canvas enrollments conclude:*)",
      "Bash(canvas api DELETE:*)",
      "Bash(canvas api PUT:*)",
      "Bash(canvas api POST:*)",
      "Bash(canvas api PATCH:*)",
      "mcp__canvas__canvas_courses_delete",
      "mcp__canvas__canvas_sections_crosslist"
    ],
    "ask": [
      "Bash(canvas courses create:*)",
      "Bash(canvas assignments grade:*)",
      "mcp__canvas__canvas_courses_create",
      "mcp__canvas__canvas_assignments_grade"
    ]
  },
  "hooks": {
    "PreToolUse": [
      {"matcher": "Bash", "hooks": [{"type": "command", "command": "${CLAUDE_PROJECT_DIR}/.claude/hooks/canvas-guard.sh"}]},
      {"matcher": "mcp__canvas__", "hooks": [{"type": "command", "command": "${CLAUDE_PROJECT_DIR}/.claude/hooks/canvas-guard.sh"}]}
    ]
  }
}
```
