# MCP Integration

Use Canvas CLI as a Model Context Protocol (MCP) server so AI clients can call Canvas tools directly.

Canvas CLI exposes each command as an MCP tool (for example `canvas_courses_list`, `canvas_assignments_get`, `canvas_version`).

## What You Get

- One MCP server with the same behavior as the normal CLI
- 530+ typed tools generated from command flags and args
- Support for local stdio clients and HTTP streaming mode
- Built-in editor helpers for Claude Desktop, Cursor, and VS Code

## Install Path Notes

Use an absolute binary path in MCP client configs.

Common paths:

- Homebrew (Apple Silicon macOS): `/opt/homebrew/bin/canvas`
- Homebrew (Intel macOS / many Linux setups): `/usr/local/bin/canvas`
- Manual binary install: wherever you placed `canvas` in `PATH`

Quick check:

```bash
which canvas
canvas version
```

## Requirements

- Canvas CLI `v1.10.0+` recommended
- If installed with `go install`, use Go `1.25+`
- Valid Canvas auth using either:
  - `canvas auth login` (OAuth), or
  - `canvas auth token set`, or
  - `CANVAS_URL` + `CANVAS_TOKEN` environment variables

## Quick Start

```bash
# 1) Check MCP commands
canvas mcp --help

# 2) Export tool schemas and verify discovery
canvas mcp tools

# 3) Configure one client (examples below), then restart the client
```

After setup, ask your AI client to run `canvas_version`.

## Transport Modes

### STDIO (recommended for local clients)

```bash
canvas mcp start
```

Use this for Claude Code, Cursor, VS Code, OpenCode, and Codex local integration.

### HTTP Stream

```bash
canvas mcp stream --host 127.0.0.1 --port 8080
```

Use this only when a client requires HTTP transport.

!!! warning "Security"
    `canvas mcp stream` does not add authentication by itself. Keep it on localhost, or place it behind an authenticated proxy/tunnel.

## Client Setup

### Claude Desktop

Auto-configure with built-in helper:

```bash
canvas mcp claude enable
canvas mcp claude list
canvas mcp claude disable   # remove the server entry
```

Optional overrides:

- `--server-name`
- `--config-path`
- `--env KEY=value`

### Claude Code (CLI)

Claude Code CLI uses its own MCP config. Add Canvas manually:

```bash
claude mcp add canvas /opt/homebrew/bin/canvas mcp start
claude mcp list
```

### Cursor

User-level setup:

```bash
canvas mcp cursor enable
canvas mcp cursor list
canvas mcp cursor disable   # remove the server entry
```

Workspace setup (recommended per project):

```bash
canvas mcp cursor enable --workspace
canvas mcp cursor list --workspace
```

This writes `.cursor/mcp.json` in your project.

### VS Code

User-level setup:

```bash
canvas mcp vscode enable
canvas mcp vscode disable   # remove the server entry
```

Workspace setup:

```bash
canvas mcp vscode enable --workspace
canvas mcp vscode list --workspace
```

This writes `.vscode/mcp.json` in your project.

### OpenCode

Use OpenCode MCP management:

```bash
opencode mcp add
opencode mcp list
```

When prompted, configure a local/stdio server with:

- command: `/opt/homebrew/bin/canvas`
- args: `mcp start`

You can also define it directly under the `mcp` section in OpenCode config.

### Codex CLI

Add Canvas as a stdio MCP server:

```bash
codex mcp add canvas -- /opt/homebrew/bin/canvas mcp start
codex mcp list
```

## Authentication and Environment Behavior

MCP uses the same auth resolution path as normal CLI commands.

Each AI client starts its own `canvas mcp start` process, so the server sees that specific client's environment variables.

Auth/input precedence is:

1. `CANVAS_URL` + `CANVAS_TOKEN` (if both are present)
2. `flags.instance` (same as CLI `--instance`)
3. `default_instance` from `~/.canvas-cli/config.yaml`

So yes:

- You can select instance per tool call with `flags.instance`.
- Without instance flag, default instance is used.
- If both `CANVAS_URL` and `CANVAS_TOKEN` exist in that process environment, they override instance/default.

### Practical Example

If Claude Code has `CANVAS_URL` and `CANVAS_TOKEN` exported, but Cursor does not:

- Claude Code MCP calls will use env auth mode.
- Cursor MCP calls will use `flags.instance` or `default_instance`.

This is expected because each client owns its own MCP process environment.

## Passing Flags in MCP Tool Calls

Tool input follows generated schema. Example shape:

```json
{
  "args": ["123"],
  "flags": {
    "instance": "staging",
    "output": "json",
    "quiet": true
  }
}
```

## Sensitive Flags

For safety, inherited flags below are excluded from MCP exposure:

- `--show-token`
- `--config`

## Tool Annotations and Read-Only Clients

Every tool carries MCP annotations describing what it does, so clients can
decide what to auto-approve:

| Annotation | Meaning |
|---|---|
| `readOnlyHint: true` | Reads only — `list`, `get`, `show`, … |
| `readOnlyHint: false` | Writes — `create`, `update`, `grade`, … |
| `destructiveHint: true` | Cannot be undone — `delete`, `conclude`, `crosslist`, … |

The hints are derived from the same classification `canvas agent guard` uses, so
a tool can never advertise itself as read-only while the guard gates it as a
write. See [Agent Safety](agent-safety.md).

Some clients run their MCP session in read-only mode and only allow a tool when
`readOnlyHint` is strictly `true`. A tool with no annotation counts as a write,
not as unknown. In those clients the write and destructive tools are filtered
out and the read tools remain usable.

Two kinds of tool are deliberately left unannotated, so read-only clients drop
them entirely:

- **`canvas api`** — the raw escape hatch can issue any HTTP verb, so no
  read-only claim can be made about it.
- **Local commands that mutate state** — `auth login`, `config set`,
  `config account`, `cache clear`, `skills install`, `update`, and similar.

To inspect what your client will see:

```bash
canvas mcp tools   # writes mcp-tools.json in the current directory
jq '[.[] | select(.annotations.readOnlyHint == true)] | length' mcp-tools.json
```

## `--instance` in MCP Calls

Use the same flag name under MCP `flags` as in CLI:

```json
{
  "flags": {
    "instance": "production"
  }
}
```

This behaves like:

```bash
canvas ... --instance production
```

## Verification Checklist

```bash
# Export and inspect available tools
canvas mcp tools

# Confirm your client sees configured MCP server
# (use client-specific list command)

# In client chat, run a harmless call first:
# "Use tool canvas_version and show the output"
```

## Troubleshooting

### Tool server not appearing in client

- Restart the client after adding MCP server
- Confirm command path is absolute (for example `/opt/homebrew/bin/canvas`)
- Run `canvas mcp tools` to verify server-side tool discovery

### Server dropped by a read-only client

The client reports that the server has no available tools after filtering, and
drops it. This means it is enforcing read-only mode and found no tool with
`readOnlyHint: true`. Check what your build emits:

```bash
canvas mcp tools
jq '[.[] | select(.annotations.readOnlyHint == true)] | length' mcp-tools.json
```

A count of `0` means the build predates tool annotations — upgrade. Anything
above zero means the server will survive filtering.

### Using wrong Canvas environment

- Check if client process has `CANVAS_URL` and `CANVAS_TOKEN` set
- If yes, those override `flags.instance` / default instance logic
- If no, pass `flags.instance` explicitly or update `default_instance`

### Unauthorized / invalid grant

- Re-authenticate: `canvas auth login --instance <name>`
- Or refresh token setup: `canvas auth token set <name> --token <token>`

### HTTP stream cannot connect

- Confirm host/port values
- Keep `--host 127.0.0.1` for local-only usage
- If remote access is needed, add your own authenticated gateway layer
