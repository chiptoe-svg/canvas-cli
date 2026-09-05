# Command surface

Every command group in the faculty build and its subcommands, from
`canvas <group> --help`. This is the whole tool: if something is not here, the
binary does not have it. Flags are not listed — read them from
`canvas <group> <sub> --help` rather than guessing.

## Global flags (any command)

| Flag | Meaning |
|---|---|
| `-o, --output table\|json\|yaml\|csv` | Output format (default table) |
| `--filter TEXT` | Case-insensitive substring match across all fields |
| `--columns a,b,c` | Select columns to display |
| `--sort field` / `--sort -field` | Sort ascending / descending |
| `--limit N` | Cap list results (0 = unlimited) |
| `--instance NAME` | Use a specific configured Canvas instance |
| `--dry-run` | Print the request as a curl command, send nothing |
| `--show-token` | Do not redact auth in `--dry-run` output |
| `--no-cache` | Bypass the response cache |
| `--quiet` | Data and errors only (for scripts) |
| `-v, --verbose` | Debug logging to stderr |

## Teaching and content

| Group | Subcommands |
|---|---|
| `courses` | `get` `list` `update`, `settings {get\|late-policy\|permissions\|tabs\|todo\|recent-students\|effective-due-dates}` |
| `course-features` | `list` `get-flag` `set-flag` `delete-flag` `enabled` |
| `assignments` | `create` `delete` `get` `list` `update` `upcoming` |
| `assignment-groups` | `create` `delete` `get` `list` `update` |
| `schedule` | (no subcommands) available/due/closed times, by `--id` or `--match`, local time |
| `overrides` | `create` `delete` `get` `list` `update` |
| `modules` | `create` `delete` `get` `list` `update` `publish` `unpublish` `relock`, `items {…}` |
| `pages` | `create` `delete` `duplicate` `front` `get` `list` `revert` `revisions` `update` |
| `quizzes` | `create` `delete` `get` `list` `update` `questions` `groups` `submissions` `reports` `statistics` `regrade` `extensions` `ip-filters` `assignment-overrides` |
| `discussions` | `create` `delete` `get` `list` `update` `entries` `entry-list` `post` `reply` `replies` `update-entry` `delete-entry` `rate-entry` `reorder` `duplicate` `subscribe` `unsubscribe` `view` `mark-read` `mark-unread` `mark-all-read` `mark-entry-read` `mark-entry-unread` `mark-entries-read` `mark-entries-unread` |
| `announcements` | `create` `delete` `get` `list` `update` |
| `rubrics` | `create` `delete` `get` `list` `update` `associate` |
| `rubric-associations` | `assess` `update` `delete` `delete-assessment` |
| `outcomes` | `create` `get` `list` `update` `groups` `link` `unlink` `results` `alignments` |
| `peer-reviews` | `create` `delete` `list` |
| `collaborations` | `list` `members` |

## Grading

| Group | Subcommands |
|---|---|
| `submissions` | `list` `get` `grade` `bulk-grade` `excuse` `missing` `download` `comments` `add-comment` `delete-comment` |
| `grades` | `columns` `feed` `history` `history-day` `history-submissions` |
| `grading-periods` | `list` `get` `update` `delete` |
| `grading-standards` | `create` `delete` `get` `list` |
| `course-extensions` | `quiz` `assignment` |

## People

| Group | Subcommands |
|---|---|
| `users` | `list` `get` `me` `profile` `search` `todo` `upcoming-events` `activity-stream` `missing-submissions` |
| `enrollments` | `create` `get` `list` `accept` `reject` `conclude` `reactivate` |
| `sections` | `create` `delete` `get` `list` `update` `crosslist` `uncrosslist` |
| `groups` | `create` `create-standalone` `delete` `get` `list` `update` `members` `memberships` `users` `invite` `categories` `permissions` `tabs` `activity-stream` `assignment-override` `collaborations` `conferences` `content-exports` `external-feeds` |
| `conversations` | `create` `get` `list` `reply` `delete` `add-recipients` `archive` `unarchive` `star` `unstar` `mark-read` `mark-all-read` `unread-count` |
| `appointment-groups` | `create` `delete` `get` `list` `update` `groups` `users` `next` |
| `analytics` | `activity` `assignments` `students` `user` |

## Files, calendar, content transfer

| Group | Subcommands |
|---|---|
| `files` | `list` `get` `upload` `download` `copy` `delete` `quota` `licenses` `set-usage-rights` `remove-usage-rights` `reset-verifier` |
| `folders` | `create` `list` `get` `update` `delete` `copy` `media` `resolve-path` |
| `calendar` | `create` `delete` `get` `list` `update` `reserve` |
| `content-exports` | `create` `get` `list` `epub-create` `epub-get` |
| `content-migrations` | `create` `get` `list` `content` `issues` `migrators` |
| `content-shares` | `list-sent` `list-received` `get` `delete` |

## Tool itself

| Group | Subcommands |
|---|---|
| `auth` | `login` `logout` `status` `token` |
| `config` | `add` `list` `show` `use` `remove` |
| `context` | `set` `show` `clear` |
| `doctor` | (no subcommands) install, config, auth and connectivity diagnostics |
| `api` | `get` only — read-only raw GET, response wrapped under `.body` |
| `activity` | `list` `archive` `clear` `configure` `path` |
| `agent` | `guard` — generate permission rules and hooks for an AI agent host |
| `alias` | `set` `list` `delete` |
| `cache` | `stats` `clear` |
| `skills` | `install` `path` `print` |
| `update` | `check` `status` `enable` `disable` |
| `completion` | shell completion for bash, zsh, fish, powershell |
| `version` | (no subcommands) prints the build; the audited one ends `+audited.<n>` |

## Scoping

Course-scoped commands take `--course-id`. Submission commands add
`--assignment-id` and `--user-id`. `canvas context set course 12345` fills
`--course-id` for `assignments list` and `assignments get` only — everywhere
else pass it explicitly, and run `canvas context show` before acting, because a
stale context silently targets the wrong course.

There is no account administration in this build: no `accounts`, `admins`,
`roles`, `sis-imports`, `enrollment-terms`, or `developer-keys`, and no way to
create a Canvas user. No command takes an account-scoped flag: `courses list`
returns the courses you are enrolled in, and every other list that could once
be account-scoped (`users`, `rubrics`, `groups`, `outcomes`) now requires
`--course-id`.
