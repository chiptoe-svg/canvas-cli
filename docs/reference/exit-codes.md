# Exit Codes

Canvas CLI uses two exit codes.

## Exit Code Reference

| Code | Name | Description |
|------|------|-------------|
| `0` | Success | Command completed successfully |
| `1` | Error | Any failure: authentication error, API error, invalid arguments, resource not found, network failure, etc. |

There are no specialized exit codes (no sysexits-style `64`-`78` codes). Every
error path exits with `1` and writes a descriptive message to **stderr**.

## Distinguishing Error Types

Since the exit code only tells you *that* a command failed, inspect stderr to
find out *why*:

```bash
#!/bin/bash
output=$(canvas courses list -o json 2>err.log > courses.json)

if [ $? -ne 0 ]; then
  if grep -qi "auth" err.log; then
    echo "Authentication required. Run: canvas auth login"
  elif grep -qi "not found" err.log; then
    echo "Resource does not exist"
  else
    echo "Command failed:"
    cat err.log
  fi
  exit 1
fi
```

Error messages include suggestions and, where relevant, links to the Canvas
API documentation.

## Using Exit Codes in Scripts

### Basic Check

```bash
#!/bin/bash
if canvas courses list -o json > courses.json; then
  echo "Success!"
else
  echo "Failed - see stderr output above"
  exit 1
fi
```

### Retry on Failure

The CLI already retries transient API failures internally (exponential
backoff, 3 attempts). If you still want a script-level retry:

```bash
#!/bin/bash
MAX_RETRIES=3
RETRY_DELAY=10

for i in $(seq 1 $MAX_RETRIES); do
  if canvas courses list -o json > courses.json; then
    break
  fi
  if [ "$i" -lt "$MAX_RETRIES" ]; then
    echo "Attempt $i failed, retrying in ${RETRY_DELAY}s..."
    sleep $RETRY_DELAY
  else
    echo "Failed after $MAX_RETRIES attempts"
    exit 1
  fi
done
```

## CI/CD Integration

### GitHub Actions

```yaml
- name: Fetch Canvas Data
  run: canvas courses list -o json > courses.json
  # Any failure exits 1 and fails the step
```

### Jenkins Pipeline

```groovy
pipeline {
  stages {
    stage('Fetch Data') {
      steps {
        script {
          def exitCode = sh(
            script: 'canvas courses list -o json > courses.json',
            returnStatus: true
          )
          if (exitCode != 0) {
            error('Canvas CLI failed - check stderr in the build log')
          }
        }
      }
    }
  }
}
```

## Best Practices

!!! tip "Check exit codes, parse stderr"
    The exit code tells you success or failure; stderr tells you the reason.
    Capture both in production scripts.

!!! tip "Use `--output json` for machine-readable results"
    Successful output goes to stdout as parseable JSON; diagnostics go to
    stderr, so the two streams never mix.

!!! warning "Don't Ignore Errors"
    Use `set -e` in bash scripts to exit on any error:
    ```bash
    #!/bin/bash
    set -e
    canvas courses list -o json > courses.json
    # Script stops here if command fails
    ```
