# Security Policy

## Supported Versions

Only the most recent audited release is supported. There is no back-porting:
a fix ships as the next `v1.13.0+audited.N` tag, and updating means re-running
the installer.

| Version | Supported |
| ------- | --------- |
| The latest `v1.13.0+audited.N` release | :white_check_mark: |
| Any earlier `+audited.N` | :x: |
| Any build without `+audited` (upstream or package-manager builds) | :x: |

`canvas version` prints what you are running; it must read
`canvas-cli 1.13.0+audited.N`. The README explains how to verify the release
signature and reproduce the build.

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability, please report it responsibly.

### How to Report

1. **Do NOT** open a public GitHub issue for security vulnerabilities
2. Email security concerns to the maintainers via GitHub's private vulnerability reporting
3. Include as much detail as possible:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### What to Expect

- **Acknowledgment**: We will acknowledge receipt within 48 hours
- **Assessment**: We will assess the vulnerability and determine its severity
- **Fix Timeline**: Critical vulnerabilities will be addressed within 7 days
- **Disclosure**: We will coordinate with you on public disclosure timing

## Security Best Practices

### Token Storage

Canvas CLI stores authentication tokens securely:

- **macOS**: Keychain (preferred)
- **Linux**: Secret Service API or encrypted file
- **Windows**: Windows Credential Manager or encrypted file

### Configuration Security

- Never commit `.canvas-cli.yaml` or any file containing tokens
- The CLI automatically adds sensitive files to `.gitignore`
- Use environment variables (`CANVAS_TOKEN`) for CI/CD environments

### API Security

- All API communication uses HTTPS
- Tokens are never logged. They are redacted in `--dry-run` output unless you explicitly opt in with the `--show-token` flag
- Rate limiting prevents accidental API abuse

## Security Scanning

This project uses automated security tools:

- **gosec**: Static analysis for security issues
- **govulncheck**: Dependency vulnerability scanning
- **Dependabot**: Automated dependency updates

## Dependencies

We regularly update dependencies to patch security vulnerabilities. Run `go mod tidy` to ensure you have the latest versions.
