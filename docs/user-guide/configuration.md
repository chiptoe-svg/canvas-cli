# Configuration

Canvas CLI can be configured to work with multiple Canvas instances and customized for your workflow.

## Configuration File

Canvas CLI stores configuration in `~/.canvas-cli/config.yaml`.

```yaml
default_instance: production
instances:
  production:
    name: production
    url: https://canvas.example.com
    # OAuth credentials (optional - for auto-refresh)
    client_id: "your-client-id"
    client_secret: "your-client-secret"
  sandbox:
    name: sandbox
    url: https://canvas-sandbox.example.com
    # API token (alternative to OAuth)
    token: "7~your-api-token-here"
settings:
  # IANA zone for commands that take or show local times (default: $TZ, then the system zone)
  timezone: America/New_York
```

### Local times

Commands that accept a wall-clock time (`--available`, `--due`, `--closed`,
`--by`, `--date`, …) read it in the local zone and send Canvas the UTC
instant. The zone is, in order: the command's `--timezone <IANA>` flag,
`settings.timezone` in the config file, the `TZ` environment variable, the
system zone. Accepted forms:

```
times:  4pm, 4:00pm, 4:00 PM, 16:50, noon, midnight
dates:  2026-09-09, 9/9/26, 9/9/2026, today, tomorrow, yesterday, sunday, this sunday, next monday
both:   2026-09-09 4:50pm, tomorrow 9am, this sunday 11:59pm
exact:  2026-09-09T16:50 (local), 2026-09-09T20:50:00Z, 2026-09-09T16:50:00-04:00
```

Dates are month/day/year. `this sunday` is the first Sunday on or after
today; `next sunday` is the one a week later. `4:50` without am/pm, a time
that does not exist on a spring-forward night, and a time that occurs twice
on a fall-back night are refused rather than guessed. Every resolved time is
printed in both the local zone and UTC.

!!! note "Authentication Methods"
    - **API Token**: Stored directly in config file (set via `canvas auth token set`)
    - **OAuth tokens**: Stored securely in your system keychain (set via `canvas auth login`)

    You can mix both methods - use OAuth for production and API tokens for testing/sandbox environments.

## Managing Instances

### Add an Instance

```bash
# Add a new instance
canvas config add production --url https://canvas.example.com

# Add with description
canvas config add staging --url https://canvas-staging.example.com --description "Staging environment"
```

After adding, authenticate with OAuth:
```bash
canvas auth login --instance production
```

Or set an API token:
```bash
canvas auth token set production --token 7~your-token-here
```

### List Instances

```bash
canvas config list
```

### Switch Default Instance

```bash
canvas config use sandbox
```

### Remove an Instance

```bash
canvas config remove sandbox
```

### Set a Default Account

Set (or auto-detect) the default account ID used when a command needs an
account ID and none is given:

```bash
# Set manually
canvas config account production 1

# Auto-detect from the API
canvas config account production
```

### Show Current Configuration

```bash
canvas config show
```

## Environment Variables

Canvas CLI supports environment variables for configuration (useful for CI/CD):

| Variable | Description |
|----------|-------------|
| `CANVAS_URL` | Canvas instance URL |
| `CANVAS_TOKEN` | API access token |
| `CANVAS_REQUESTS_PER_SEC` | Rate limit (default: 5.0) |

Example:

```bash
export CANVAS_URL=https://canvas.example.com
export CANVAS_TOKEN=your-api-token
canvas courses list
```

!!! tip "Priority"
    Environment variables take precedence over the config file.

See the [Environment Variables reference](../reference/environment-variables.md) for the full list and details.

## Command-Line Overrides

You can override configuration with command-line flags:

```bash
# Override instance
canvas courses list --instance https://other-canvas.example.com

# Override output format
canvas courses list --output json

# Disable caching
canvas courses list --no-cache
```

## Multiple Instances

Canvas CLI supports working with multiple Canvas instances. This is useful for:

- Development vs. production environments
- Multiple institutions
- Testing and staging

### Switching Instances

```bash
# Use a specific instance for one command
canvas courses list --instance sandbox

# Switch default instance
canvas config use sandbox
```

### Syncing Between Instances

```bash
# Sync course 123 from production to course 456 on sandbox
canvas sync course production 123 sandbox 456
```

See the [Course Sync Tutorial](../tutorials/course-sync.md) for more details.
