# Canvas CLI — auth, instances, and context

Loaded on demand by the `canvas-cli` skill. Authoritative docs:
https://jjuanrivvera.github.io/canvas-cli/getting-started/authentication/

## Three auth methods

| Method | Setup | Best for |
|---|---|---|
| Environment variables | `export CANVAS_URL=… CANVAS_TOKEN=…` | CI/CD, containers, one-off scripts |
| API token | `canvas auth token set <name> --url URL --token 7~…` | Personal use, simplest persistent setup |
| OAuth 2.0 + PKCE | `canvas auth login --instance URL` | Most secure; tokens in OS keychain |

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
| `CANVAS_OUTPUT` | Default output format | `table` |
| `CANVAS_NO_CACHE` | Disable response caching | `false` |

CI example (GitHub Actions):

```yaml
- name: Run Canvas CLI
  env:
    CANVAS_URL: ${{ secrets.CANVAS_URL }}
    CANVAS_TOKEN: ${{ secrets.CANVAS_TOKEN }}
  run: canvas courses list -o json
```

## Multiple instances

```bash
canvas config add production --url https://canvas.instructure.com
canvas config add staging --url https://staging.canvas.example.com
canvas config list                  # see all + which is default
canvas config use staging           # switch default
canvas courses list --instance production   # one-off override
canvas auth login --instance production     # authenticate each separately
```

Mixing auth types per instance is fine (OAuth for prod, token for sandbox);
`canvas auth status` shows the auth type and state of every instance.

Config lives at `~/.canvas-cli/config.yaml`:

```yaml
default_instance: production
instances:
  production:
    url: https://canvas.instructure.com
context:
  course_id: 12345
```

## Context (default IDs)

Avoid repeating `--course-id` etc. during a working session:

```bash
canvas context set course 12345      # fills --course-id when omitted
canvas context set assignment 67890  # fills --assignment-id
canvas context set user 111          # fills --user-id
canvas context set account 1         # fills --account-id
canvas context show
canvas context clear [course|assignment|user|account]
```

Explicit flags always beat context. Before acting on the user's behalf, run
`canvas context show` — a stale context can silently target the wrong course.

## Headless / remote auth

OAuth on a machine without a browser:

```bash
canvas auth login --instance https://myschool.instructure.com --mode oob
```

Prints an authorization URL, then prompts for the pasted code. For fully
non-interactive setups prefer env vars or `auth token set`.

## Verification

```bash
canvas auth status     # per-instance auth type + state
canvas doctor          # full diagnostics: binary, config, auth, connectivity
canvas users me        # confirms who the API sees you as
```
