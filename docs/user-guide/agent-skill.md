# Agent Skill

Canvas CLI ships an **agent skill** — a `SKILL.md` (plus condensed reference
files) that teaches AI coding agents (Claude Code, Cursor, Codex, Gemini CLI,
Windsurf, GitHub Copilot, and more) how and when to drive the `canvas` CLI.
Once installed, you can just ask your agent things like *"grade everyone in
grades.csv for course 123"* or *"create next week's assignments"* and it knows
the commands, flags, and safety rules.

## Install the skill

The skill is bundled in the CLI binary. Installing it from the CLI guarantees
the skill matches the version you are running:

```bash
canvas skills install                 # current project (./.claude/skills)
canvas skills install --global        # user-level (~/.claude/skills)
canvas skills install --agent cursor --global
canvas skills path                    # show where it would install
canvas skills print                   # print the SKILL.md
```

Supported `--agent` values: `claude`, `cursor`, `windsurf`, `codex`, `gemini`,
`copilot`, `opencode`. Use `--dir` to write to any custom skills directory.

## Prerequisites

The skill **wraps** the `canvas` binary — it doesn't bundle it. So:

1. Install the CLI:
   `curl -fsSL https://raw.githubusercontent.com/chiptoe-svg/canvas-cli/release/audited/install.sh | sh`
2. Authenticate: `canvas auth login` or `canvas auth token set`
   (or set `CANVAS_URL`/`CANVAS_TOKEN`).
3. Verify: `canvas doctor`.

## Skill vs MCP

- **Skill** (this page) — markdown instructions for shell-capable agents. The
  agent runs `canvas …` commands directly. Lowest friction; best for CLI use.
- **[MCP server](mcp.md)** — `canvas mcp start` exposes every command as a
  structured tool, for clients that prefer tool schemas (Claude Desktop, etc.).

Use whichever your agent supports; they can coexist.

## What the agent learns

The skill teaches the auth → discover → act → verify workflow, the golden
rules (preview writes with `--dry-run`, parse with `-o json`, resolve IDs
live, bound large lists with `--limit`, confirm destructive actions), the
resource/action map, output filtering (`--filter`/`--columns`/`--sort`),
grading workflows (single and bulk CSV), course-content publishing,
cross-instance sync, the raw `canvas api` escape hatch, and error handling.

Condensed references ship alongside it:

- `references/canvas-commands.md` — command/flag cheatsheet
- `references/auth-and-config.md` — auth methods, env vars, instances, context
- `references/output-and-filtering.md` — formats, jq patterns, script hygiene
- `references/grading-workflows.md` — proposal and read-back discipline for grading
