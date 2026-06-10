# Installation Guide

This guide covers all the ways to install Canvas CLI on your system.

## System Requirements

- **Operating Systems**: macOS, Linux, Windows
- **Architecture**: amd64, arm64
- **Go Version** (if building from source): 1.25 or later

## Installation Methods

### Method 1: Homebrew (macOS/Linux) - Recommended

```bash
# Add the Canvas CLI tap
brew tap jjuanrivvera/canvas-cli

# Install Canvas CLI
brew install canvas-cli

# Verify installation
canvas version
```

#### Upgrading

```bash
brew upgrade canvas-cli
```

### Method 2: Using Go

If you have Go installed:

```bash
# Install latest version
go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest

# Install specific version
go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@v1.0.0

# Verify installation
canvas version
```

!!! note "MCP + go install"
    MCP support (`canvas mcp ...`) requires Go `1.25+` when installing with `go install`.
    Homebrew and release binaries are not affected by local Go version.

### Method 3: Download Binary (All Platforms)

1. Visit the [Releases page](https://github.com/jjuanrivvera/canvas-cli/releases)
2. Download the appropriate archive for your platform:
   - **macOS (Intel)**: `canvas-cli_darwin_x86_64.tar.gz`
   - **macOS (Apple Silicon)**: `canvas-cli_darwin_arm64.tar.gz`
   - **Linux (64-bit)**: `canvas-cli_linux_x86_64.tar.gz`
   - **Linux (ARM64)**: `canvas-cli_linux_arm64.tar.gz`
   - **Windows (64-bit)**: `canvas-cli_windows_x86_64.zip`

3. Extract the archive:
   ```bash
   # macOS/Linux
   tar -xzf canvas-cli_*.tar.gz

   # Windows - use your preferred extraction tool
   ```

4. Move the binary to your PATH:
   ```bash
   # macOS/Linux
   sudo mv canvas /usr/local/bin/

   # Windows - add to PATH or move to C:\Windows\System32\
   ```

5. Verify installation:
   ```bash
   canvas version
   ```

#### Verifying a downloaded release

Every release ships a `checksums.txt` signed keylessly with
[cosign](https://github.com/sigstore/cosign), and each archive includes a
Software Bill of Materials (`*.sbom.json`).

```bash
# 1. Verify the checksums file signature (requires cosign)
cosign verify-blob \
  --certificate checksums.txt.pem \
  --signature checksums.txt.sig \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  --certificate-identity-regexp 'https://github.com/jjuanrivvera/canvas-cli.*' \
  checksums.txt

# 2. Verify your archive against the checksums
shasum -a 256 --check --ignore-missing checksums.txt
```

!!! tip "Self-update verifies too"
    `canvas update` downloads the matching `checksums.txt` and refuses to
    install a binary it cannot verify.

### Method 4: Docker

Images are published to GHCR for every release: `:latest`, `:v1` (major), and
exact versions like `:1.9.0`.

```bash
# Run Canvas CLI in Docker
docker run --rm ghcr.io/jjuanrivvera/canvas-cli:latest version

# Authenticate via environment variables (recommended for containers)
docker run --rm -e CANVAS_URL -e CANVAS_TOKEN \
  ghcr.io/jjuanrivvera/canvas-cli:latest courses list

# Create an alias for easier use
alias canvas='docker run -i --rm -e CANVAS_URL -e CANVAS_TOKEN ghcr.io/jjuanrivvera/canvas-cli:latest'

# Now use as normal
canvas courses list
```

!!! note "Minimal image"
    The image is distroless and runs as a non-root user — there is no shell
    inside. Pass credentials with `CANVAS_URL`/`CANVAS_TOKEN` environment
    variables instead of mounting a config directory.

### Method 5: Build from Source

```bash
# Clone the repository
git clone https://github.com/jjuanrivvera/canvas-cli.git
cd canvas-cli

# Build
make build

# Install
make install

# Verify
canvas version
```

## Shell Completion

Enable tab completion for your shell:

### Bash

```bash
# Generate completion script
canvas completion bash > /etc/bash_completion.d/canvas

# Or for user-level installation
canvas completion bash > ~/.canvas-completion.bash
echo 'source ~/.canvas-completion.bash' >> ~/.bashrc
```

### Zsh

```bash
# Generate completion script
canvas completion zsh > "${fpath[1]}/_canvas"

# Reload completions
autoload -U compinit && compinit
```

### Fish

```bash
# Generate completion script
canvas completion fish > ~/.config/fish/completions/canvas.fish
```

### PowerShell

```powershell
# Generate completion script
canvas completion powershell | Out-String | Invoke-Expression

# Add to profile for persistence
canvas completion powershell >> $PROFILE
```

## Verify Installation

After installation, verify everything is working:

```bash
# Check version
canvas version

# Run diagnostics
canvas doctor

# Test authentication (will prompt for credentials)
canvas auth login --instance https://canvas.instructure.com
```

## Troubleshooting

### Command not found

If you get `command not found`, ensure the installation directory is in your PATH:

```bash
# Check PATH
echo $PATH

# Add to PATH (macOS/Linux)
export PATH="$PATH:/usr/local/bin"

# Make permanent by adding to ~/.bashrc or ~/.zshrc
echo 'export PATH="$PATH:/usr/local/bin"' >> ~/.bashrc
```

### Permission denied

If you get permission errors:

```bash
# Make binary executable
chmod +x /path/to/canvas

# Or use sudo for installation
sudo mv canvas /usr/local/bin/
```

### macOS Security Warning

On macOS, you may need to allow the app in System Preferences:

1. Try to run `canvas version`
2. Go to **System Preferences > Security & Privacy**
3. Click **"Allow Anyway"** for Canvas CLI
4. Run the command again

## Updating

### Self-Update (built in)

```bash
# Check for and install the latest release (checksum-verified)
canvas update
```

### Homebrew

```bash
brew upgrade canvas-cli
```

### Go

```bash
go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest
```

### Binary

Download and replace the binary with the latest version from the releases page.

### Docker

```bash
docker pull ghcr.io/jjuanrivvera/canvas-cli:latest
```

## Uninstalling

### Homebrew

```bash
brew uninstall canvas-cli
brew untap jjuanrivvera/canvas-cli
```

### Go

```bash
rm $(which canvas)
```

### Binary

```bash
# Find and remove the binary
sudo rm /usr/local/bin/canvas

# Remove configuration and cache
rm -rf ~/.canvas-cli
```

### Docker

```bash
docker rmi ghcr.io/jjuanrivvera/canvas-cli:latest
```

## Next Steps

After installation, continue with:
- [Authentication Guide](authentication.md) - Set up OAuth
- [Command Reference](../commands/index.md) - Learn available commands
- [Tutorials](../tutorials/index.md) - See common use cases
- [AI Agent Skill](../user-guide/agent-skill.md) - Let Claude Code, Cursor, and other agents drive the CLI
