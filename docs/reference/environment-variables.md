# Environment Variables

Canvas CLI can be configured using environment variables.

## Available Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `CANVAS_URL` | Canvas instance URL (env auth mode) | From config |
| `CANVAS_TOKEN` | API access token | From config/keyring |
| `CANVAS_REQUESTS_PER_SEC` | Request rate limit for env auth mode | `5.0` |
| `CANVAS_CLI_MACHINE_ID` | Override the machine ID used to derive the encryption key for file-based token storage | Auto-detected |

Output format and caching are controlled by the `--output` and `--no-cache`
flags (or the config file), not by environment variables.

## Usage

### Session Override

Set for the current shell session:

```bash
export CANVAS_URL=https://canvas.example.com
export CANVAS_TOKEN=your-api-token
canvas courses list
```

### Permanent Configuration

Add to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
# Canvas CLI configuration
export CANVAS_URL=https://canvas.example.com
export CANVAS_TOKEN=your-api-token
```

## Variable Details

### CANVAS_URL

The Canvas LMS instance URL.

```bash
export CANVAS_URL=https://canvas.instructure.com
```

!!! note "Priority"
    When both `CANVAS_URL` and `CANVAS_TOKEN` are set, CLI commands use env authentication directly.
    In that mode, `--instance` and `default_instance` are not used.

### CANVAS_TOKEN

Canvas API access token. Generate from Canvas Account Settings.

```bash
export CANVAS_TOKEN=7~AbCdEfGhIjKlMnOpQrStUvWxYz123456789
```

### CANVAS_REQUESTS_PER_SEC

Request rate limit used when env auth mode is active.

```bash
export CANVAS_REQUESTS_PER_SEC=10.0
```

!!! warning "Security"
    Avoid setting tokens in shared environments. Consider using the config file or keyring instead.

### CANVAS_CLI_MACHINE_ID

Advanced. When tokens are stored in the encrypted file fallback (no system
keyring available), the encryption key is derived from a machine identifier.
Set this to pin that identifier — for example in containers where the
auto-detected machine ID changes between runs.

```bash
export CANVAS_CLI_MACHINE_ID=my-stable-id
```

## Precedence Order

Configuration is applied in this order (highest precedence first):

1. If both `CANVAS_URL` and `CANVAS_TOKEN` are set, env auth mode is used
2. Otherwise, command-line flags (`--instance`, `--output`, etc.)
3. Otherwise, configuration file (`~/.canvas-cli/config.yaml`)
4. Built-in defaults

!!! note "MCP clients"
    This same precedence applies to MCP tool calls, because MCP executes the same command logic.
    Each MCP client runs its own server process and therefore may have different environment variables.

## Example Configurations

### Development Environment

```bash
# Use sandbox instance
export CANVAS_URL=https://canvas-sandbox.example.com
export CANVAS_TOKEN=sandbox-token
```

### CI/CD Pipeline

```yaml
# GitHub Actions example
env:
  CANVAS_URL: ${{ secrets.CANVAS_URL }}
  CANVAS_TOKEN: ${{ secrets.CANVAS_TOKEN }}
```

For JSON output in scripts, pass the flag: `canvas courses list -o json`.

### Multi-Instance Setup

```bash
# Function to switch instances
canvas-prod() {
  export CANVAS_URL=https://canvas.example.com
  export CANVAS_TOKEN=$CANVAS_PROD_TOKEN
}

canvas-sandbox() {
  export CANVAS_URL=https://canvas-sandbox.example.com
  export CANVAS_TOKEN=$CANVAS_SANDBOX_TOKEN
}
```

## Troubleshooting

### Variable Not Working

1. Verify the variable is set:
   ```bash
   echo $CANVAS_URL
   ```

2. Check for typos in variable names

3. Ensure you've sourced your profile:
   ```bash
   source ~/.bashrc
   ```

### Token Security

If your token is exposed:

1. Immediately regenerate it in Canvas Account Settings
2. Update your configuration
3. Consider using the config file with restrictive permissions:
   ```bash
   chmod 600 ~/.canvas-cli/config.yaml
   ```
