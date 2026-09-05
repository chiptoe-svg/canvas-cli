# Canvas CLI — auth, instances, and context

Loaded on demand by the `canvas-cli` skill. Authoritative docs:
https://jjuanrivvera.github.io/canvas-cli/getting-started/authentication/

## Three auth methods

| Method | Setup | Best for |
|---|---|---|
| Environment variables | `export CANVAS_URL=… CANVAS_TOKEN=…` | One-off scripts, a shared machine |
| API token | `canvas auth token set <name> --url URL --token 7~…` | Simplest persistent setup |
| OAuth 2.0 + PKCE | `canvas auth login <url>`, then `--instance <name>` | Most secure; tokens in the OS keychain |

Precedence (highest first):

1. `CANVAS_URL` + `CANVAS_TOKEN` env vars — when both are set, env-auth mode is
   used and `--instance`/`default_instance` are ignored.
2. API token stored in config (`canvas auth token set`).
3. OAuth tokens in the system keychain (`canvas auth login`).

Tokens are generated in Canvas under **Account → Settings → Approved
Integrations → + New Access Token** (shown only once).

## Environment variables

| Variable | Meaning | Default |
|---|---|---|
| `CANVAS_URL` | Canvas instance URL (enables env-auth with token) | from config |
| `CANVAS_TOKEN` | API access token | from config/keyring |
| `CANVAS_REQUESTS_PER_SEC` | Rate limit in env-auth mode | `5.0` |

Output format and caching are flag-controlled (`-o`, `--no-cache`) — there are
no `CANVAS_OUTPUT`/`CANVAS_NO_CACHE` env vars.

A token is a live credential for the instructor's whole Canvas account. Never
print one, never write one into a file the agent produces, and never pass
`--show-token` unless the instructor asked to see the request verbatim.

## Multiple instances

An instructor who teaches at two schools, or has a sandbox alongside the live
instance, keeps both in the config:

```bash
canvas config add myschool --url https://myschool.instructure.com
canvas config add sandbox --url https://myschool.test.instructure.com
canvas config list                        # all instances + which is default
canvas config use myschool                # switch the default
canvas courses list --instance sandbox    # one-off override
canvas auth login --instance sandbox      # each is authenticated separately
```

Mixing auth types per instance is fine (OAuth for the live one, a token for the
sandbox); `canvas auth status` shows the auth type and state of every instance.
Before any write, check which instance is active — grading against the sandbox
looks like it worked and changes nothing the students see, and grading the live
instance when the sandbox was meant is worse.

Config lives at `~/.canvas-cli/config.yaml`:

```yaml
default_instance: myschool
instances:
  myschool:
    url: https://myschool.instructure.com
settings:
  timezone: America/New_York
context:
  course_id: 12345
```

## Context (default IDs)

Store default IDs for a working session:

```bash
canvas context set course 12345      # fills --course-id for assignments list/get
canvas context set assignment 67890  # stored, not yet consumed by commands
canvas context set user 111          # stored, not yet consumed by commands
canvas context show
canvas context clear [course|assignment|user]
```

Only `assignments list`/`assignments get` read the course context today; pass
explicit flags everywhere else. Explicit flags always beat context. Before
acting on the instructor's behalf, run `canvas context show` — a stale context
from last week's course silently targets the wrong students.

## Headless / remote auth

OAuth on a machine without a browser:

```bash
canvas auth login https://myschool.instructure.com --mode oob
```

Prints an authorization URL, then prompts for the pasted code. For fully
non-interactive setups prefer env vars or `auth token set`.

## Verification

```bash
canvas auth status     # per-instance auth type + state
canvas doctor          # full diagnostics: binary, config, auth, connectivity
canvas users me        # confirms who the API sees you as
```
