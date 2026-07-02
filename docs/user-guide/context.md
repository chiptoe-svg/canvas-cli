# Context Management

Set default values for common flags like `--course-id` so you don't have to type them repeatedly.

## Overview

When working with a specific course, you typically run many commands with the same `--course-id`. Context management lets you set this once and have supported commands apply it automatically:

```bash
# Without context - repetitive
canvas assignments list --course-id 12345
canvas assignments get 67890 --course-id 12345

# With context - clean and simple
canvas context set course 12345
canvas assignments list
canvas assignments get 67890
```

!!! note "Which commands use context"
    Context is currently applied by `canvas assignments list` and
    `canvas assignments get` (course context). Other commands still require
    their flags explicitly. The other context types (`assignment`, `user`,
    `account`) are stored for future use but not yet consumed by any command.

## Setting Context

Use `canvas context set` to set default values:

```bash
canvas context set <type> <id>
```

### Available Context Types

| Type | Related flag | Example | Used by |
|------|--------------|---------|---------|
| `course` | `--course-id` | `canvas context set course 12345` | `assignments list`, `assignments get` |
| `assignment` | `--assignment-id` | `canvas context set assignment 67890` | (stored, not yet consumed) |
| `user` | `--user-id` | `canvas context set user 111` | (stored, not yet consumed) |
| `account` | `--account-id` | `canvas context set account 1` | (stored, not yet consumed) |

## Viewing Context

See your current context settings:

```bash
canvas context show
```

Output:
```
Current context:
  course_id:     12345
  assignment_id: 67890
```

## Clearing Context

### Clear All Context

```bash
canvas context clear
```

### Clear Specific Value

```bash
canvas context clear course
canvas context clear assignment
```

## How Context Works

1. When you run a supported command, Canvas CLI checks if the flag is provided
2. If the flag is missing, it falls back to the stored context value
3. Explicit flags always override context values

### Example Flow

```bash
# Set context
canvas context set course 12345

# These are equivalent:
canvas assignments list                    # Uses context (course 12345)
canvas assignments list --course-id 12345  # Explicit flag

# Override context with a different value
canvas assignments list --course-id 99999  # Uses 99999, ignores context
```

## Storage

Context is stored in your configuration file (`~/.canvas-cli/config.yaml`):

```yaml
context:
  course_id: 12345
  assignment_id: 67890
```

## Workflow Examples

### Course Review Workflow

```bash
# Working on a specific course
canvas context set course 12345

# Context fills in the course ID
canvas assignments list
canvas assignments get 67890

# Commands without context support still take explicit flags
canvas modules list --course-id 12345
canvas users list --course-id 12345 --enrollment-type student

# Switch to another course
canvas context set course 54321
```

### Combined with Aliases

Context and [aliases](aliases.md) work great together:

```bash
# Set course context
canvas context set course 12345

# Create an alias that benefits from context
canvas alias set hw "assignments list"

# Use it - context provides the course ID
canvas hw
```

## Tips

- Set context at the start of a work session
- Clear context when switching tasks to avoid confusion
- Use `canvas context show` to verify your current context
- Explicit flags always take precedence over context
