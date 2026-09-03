---
name: canvas-cli
description: Manage Canvas LMS (https://www.instructure.com/canvas) from the terminal with the `canvas` CLI — courses, assignments, submissions and grading, modules, pages, quizzes, discussions, announcements, users, enrollments, sections, files, and analytics. Use this whenever the user wants to list or create assignments, grade submissions (single or bulk from CSV), manage course content, enroll users, upload or download course files, pull course or student analytics, sync content between Canvas instances, or script any Canvas LMS teaching/administration task.
version: 1.9.0
homepage: https://github.com/jjuanrivvera/canvas-cli
license: MIT
allowed-tools: Bash(canvas:*)
metadata: {"openclaw":{"category":"education","emoji":"🎓","requires":{"bins":["canvas"]},"install":[{"kind":"brew","formula":"jjuanrivvera/canvas-cli/canvas-cli","bins":["canvas"]},{"kind":"go","package":"github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest","bins":["canvas"]}]}}
---

# Canvas CLI

Drive [Canvas LMS](https://www.instructure.com/canvas) through the `canvas`
command-line tool. This skill teaches you how and when to use it.

## Prerequisites

- The `canvas` binary must be on `PATH`. Check with `canvas version`. If
  missing, install it: `brew tap jjuanrivvera/canvas-cli && brew install
  canvas-cli` or `go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest`.
- Credentials, one of:
  - **Environment variables** (best for CI/non-interactive): set `CANVAS_URL`
    and `CANVAS_TOKEN`. This takes priority over everything else.
  - **API token**: `canvas auth token set myschool --url
    https://myschool.instructure.com --token 7~...` (token from Canvas →
    Account → Settings → New Access Token).
  - **OAuth**: `canvas auth login --instance https://myschool.instructure.com`
    (opens a browser; add `--mode oob` on headless machines).
- Confirm with `canvas auth status` or `canvas doctor`.

Details on multi-instance setup and precedence:
`references/auth-and-config.md`.

## Golden rules (read before acting)

1. **Preview writes.** Every command accepts `--dry-run`, which prints the
   exact HTTP request as a curl command instead of executing it (tokens
   redacted). For any create/update/delete/grade, run once with `--dry-run`,
   show the user, then run for real. `submissions bulk-grade --dry-run`
   previews the whole batch.
2. **Parse with JSON.** Add `-o json` and pipe to `jq` when you need to read
   values; the default `table` output is for humans.
3. **Resolve IDs live.** Course, assignment, user, section, and module IDs are
   instance-specific — never guess them. Look them up
   (`canvas courses list`, `canvas assignments list --course-id N`,
   `canvas users search "name"`).
4. **Bound big lists.** Account-level lists can be huge; use `--limit N` and
   `--search`/`--filter` instead of dumping everything.
5. **Confirm destructive actions** (`delete`, `conclude`) with the user before
   running them, and prefer `--dry-run` first.
6. **Set context for a working session.** `canvas context set course 12345`
   makes `--course-id` implicit for `assignments list`/`assignments get` (other
   commands still need explicit flags); explicit flags always override it.
   Check what's active with `canvas context show`.
7. **Mind the instance.** With multiple configured instances, verify which one
   is active (`canvas config list`) before writing; switch with
   `canvas config use <name>` or per-command `--instance`.
8. **Prefer `canvas` over curl.** Never hand-roll `curl` against the Canvas
   API when this CLI is available: it handles auth, pagination, rate limiting,
   and retries for you. For endpoints without a dedicated command, use
   `canvas api` (see the raw API escape hatch below).

## Workflow: auth → discover → act → verify

```bash
canvas doctor                          # 1. verify install, auth, connectivity
canvas courses list -o json            # 2. find real IDs
canvas <resource> --help               # 3. discover actions & flags
canvas assignments create --course-id 123 --name "Quiz 1" --points 100 --dry-run   # 4. preview
canvas assignments create --course-id 123 --name "Quiz 1" --points 100             #    act
canvas assignments list --course-id 123 --filter "Quiz 1"                          # 5. verify
```

## Command map

`canvas <resource> <action>` — most resources support
`list|get|create|update|delete` plus resource-specific actions. Main resources:

| Area | Resources |
|---|---|
| Teaching | `courses`, `assignments`, `assignment-groups`, `modules`, `pages`, `quizzes` (incl. `reports`, `statistics`, `question-groups`, `ip-filters`), `discussions`, `announcements`, `rubrics`, `rubric-associations`, `outcomes`, `overrides`, `peer-reviews`, `polls` |
| Grading | `submissions` (`grade`, `bulk-grade`, `add-comment`), `grades`, `grading-periods`, `grading-standards`, `grading-period-sets`, `live-assessments` |
| People | `users`, `enrollments`, `sections`, `groups` (`memberships`, `categories`), `conversations`, `comm-channels`, `observees`, `appointment-groups` |
| Content & files | `files`, `folders`, `calendar`, `content-migrations`, `content-exports`, `content-shares`, `blueprint`, `course-pacing`, `blackout-dates`, `media`, `eportfolios` |
| Personal | `favorites`, `bookmarks`, `course-nicknames`, `planner`, `history` |
| Admin | `accounts`, `admins`, `roles`, `analytics`, `sis-imports`, `external-tools`, `auth-providers`, `csp-settings`, `account-notifications`, `account-reports`, `enrollment-terms`, `developer-keys`, `audit` |
| Utility | `api` (raw requests), `sync`, `context`, `alias`, `cache`, `doctor`, `mcp`, `repl`/`shell`, `webhook`, `jwts`, `progress` |

The CLI has ~93 command groups covering most of the Canvas REST API. The table
above is a guide, not the full list — **always discover the real surface with
`canvas --help` and `canvas <resource> --help`** rather than assuming a command
exists. Most resources that exist under a course also exist under a group or
user context via `--group-id`/`--user-id` (e.g. discussions, pages, files,
folders, content-migrations). A condensed cheatsheet ships in
`references/canvas-commands.md`.

```bash
# Courses
canvas courses list                              # your enrolled courses
canvas courses list --account-id 1 --search "Biology"   # admin: account courses
canvas courses get 123 -o json | jq '{id,name,course_code}'

# Assignments
canvas assignments list --course-id 123 --bucket upcoming
canvas assignments create --course-id 123 --name "Essay" --points 50 \
  --due-at "2026-08-01T23:59:00Z" --grading-type points
echo '{"name":"Quiz 1","points_possible":100}' | canvas assignments create --course-id 123 --stdin

# Users & enrollments
canvas users search "john doe"
canvas users list --course-id 123 --enrollment-type student
canvas enrollments create --course-id 123 --user-id 456 --type StudentEnrollment --state active

# Modules (content structure)
canvas modules create --course-id 123 --name "Week 1"
canvas modules items create --course-id 123 --module-id 9 --type Assignment --content-id 456
canvas modules publish --course-id 123 9
```

## Output: formats, filtering, columns, sorting

Global flags work on every command:

```bash
canvas courses list -o json | jq '.[].id'        # json | yaml | csv | table
canvas assignments list --course-id 123 --filter "exam"      # substring, all fields
canvas assignments list --course-id 123 --columns id,name,due_at,points_possible
canvas assignments list --course-id 123 --sort -due_at       # '-' prefix = descending
canvas users list --course-id 123 -o csv > roster.csv
canvas users list --account-id 1 --limit 100                 # cap result count
```

Use built-in `--filter` for simple matching and `-o json | jq` for anything
structural. Details: `references/output-and-filtering.md`.

## Local times

Commands that take a wall-clock time (`--due`, `--available`, `--closed`,
`--by`, `--date`) read it in the user's local zone and send Canvas the UTC
instant — never convert to UTC yourself. Pass what the user said:
`--due "4:50pm"`, `--due "9/9/26"`, `--due "this sunday 11:59pm"`,
`--by "next monday"`; `4pm`, `16:50`, `noon`, `today`, `tomorrow`,
`2026-09-09`, `9/9/2026` and RFC 3339 also work. The zone is `--timezone
<IANA>`, else `settings.timezone` in the config, else `$TZ`, else the
system zone — ask the user for their zone if none of those is set and the
resolved zone in the output looks wrong. Ambiguous input (`4:50` without
am/pm, a time in a DST gap) is refused with the accepted forms; fix the
input rather than guessing. Output always shows the resolved local time and
the UTC value — report both to the user.

## What is due soon (`assignments upcoming`)

`canvas assignments upcoming` is read-only: per course it lists the
assignments (quizzes and graded discussions included, through their
assignments) due after now and up to a limit, sorted by due date, with the
due time in local time, points, submission type and published state, and a
one-line summary per course. The limit is `--within 36h|10d|2w` or `--by
<local date>` (a date alone covers the whole day). Undated items only with
`--include-undated`; unpublished only with `--published-only=false`.

```bash
# "In GC 1010, is anything due this Sunday?"
canvas assignments upcoming --course-id 123 --by "this sunday"

# "Which assignments are due in the next 10 days in GC 4800?"
canvas assignments upcoming --course-id 456 --within 10d

# Every course you teach this term, formatted for chat
canvas assignments upcoming --all-active --within 2w -o markdown
```

Report the per-course summary line (`N due by <local time>`) and the rows;
`-o json` carries `{now, limit, timezone, courses[].assignments[]}` with
both `due_at` (UTC) and `due_local`. Due dates are the base dates —
section or student overrides are not consulted.

## Workflow: grade a submission

```bash
# Find who/what to grade
canvas assignments list --course-id 123 --filter "Essay"
canvas submissions list --course-id 123 --assignment-id 456 --workflow-state submitted

# Grade one submission (score, letter grade, or excuse)
canvas submissions grade --course-id 123 --assignment-id 456 --user-id 789 \
  --score 95 --comment "Great work"

# Verify
canvas submissions get --course-id 123 --assignment-id 456 --user-id 789 -o json
```

## Workflow: excuse a student (`submissions excuse`)

"Excuse <student> from <quiz/assignment>" is one command — no id hunting:

```bash
canvas submissions excuse --course-id 123 --student "Ada Lovelace" --assignment "Quiz 3"
canvas submissions excuse --course-id 123 --student "lovelace" --assignment "lineup" --dry-run   # preview
canvas submissions excuse --course-id 123 --student 789 --assignment 456 --force                 # ids, no prompt
canvas submissions excuse --course-id 123 --student "Ada Lovelace" --assignment "Quiz 3" --unexcuse
```

`--student` takes an id, the exact name, the sortable name ("Lovelace,
Ada"), a login/SIS id, or a substring that matches exactly one active
student; `--assignment` takes an assignment id, a quiz id, the exact name,
or a substring that matches exactly one item (a quiz resolves to the
assignment it is graded under). Zero or several matches are refused with
the candidates listed — pick one and re-run; never guess. The command
prints who and what it resolved, reads the submission back and prints
`excused: not excused → excused` and `verified: yes`; "done" means
`verified: yes`. A `verified: no` line means the read-back disagreed and
the command exited non-zero — report it. `-o json` carries
`{student, assignment, before_excused, after_excused, verified}`.

## Workflow: bulk grading from CSV

CSV columns: `user_id,assignment_id,score,comment`.

```bash
canvas submissions bulk-grade --course-id 123 --csv-file grades.csv --dry-run   # preview every change
canvas submissions bulk-grade --course-id 123 --csv-file grades.csv             # apply (concurrent)
canvas grades history --course-id 123                                      # audit afterwards
```

## Workflow: set quiz / assignment times (`schedule`)

`canvas schedule` sets the three availability times — available from
(`unlock_at`), due (`due_at`), closed (`lock_at`) — on quizzes **and**
assignments in one command, in local time, one item (`--id`) or every item
whose title matches (`--match <substring|/regex/>`, narrowed by `--type
quiz|assignment|all`). A quiz-backed assignment is always updated through
its quiz, and no item is written twice. Time-only values apply to each
matched item's own date; add `--date` to move them. A date-only
`--due`/`--closed` means 11:59 PM that day. It refuses a plan where
available ≤ due ≤ closed would not hold after merging with the current
values, reads every item back, prints before/after, and exits non-zero on
any mismatch.

```bash
# "make that quiz available at 4:00 and both due and closed at 4:50pm"
canvas schedule --course-id 123 --id 456 --type quiz --available 4:00pm --due 4:50pm --closed 4:50pm

# "make every attendance quiz have these times" — times apply to each
# quiz's own date; add --date to move them. Preview, then apply.
canvas schedule --course-id 123 --match attendance --type quiz --available 4pm --due 4:50pm --closed 4:50pm --dry-run
canvas schedule --course-id 123 --match attendance --type quiz --available 4pm --due 4:50pm --closed 4:50pm --force

# "make the course lineup assignment due 9/9/26"
canvas schedule --course-id 123 --match "course lineup" --due 9/9/26

# remove a close date
canvas schedule --course-id 123 --id 789 --type assignment --clear closed
```

Always `--dry-run` a `--match` first and show the user the plan (item,
type, each time old → new in local and UTC, and the exact `PUT` requests).
Without `--force`, `--match` prompts for confirmation. "Done" means the
result table shows `yes` under VERIFIED for every changed item and the
summary line ends in `0 mismatched, 0 failed`; otherwise report the
mismatch lines to the user. `-o json` carries `{items[].before, after,
read_back, verified, mismatches}, summary`.

## Workflow: publish course content

```bash
canvas pages create --course-id 123 --title "Syllabus" --body "<p>…</p>" --published
canvas announcements create --course-id 123 --title "Welcome" --message "Class starts Monday"
canvas quizzes create --course-id 123 --title "Midterm" --quiz-type assignment --time-limit 60
canvas discussions create --course-id 123 --title "Week 1 discussion"
canvas files upload syllabus.pdf --course-id 123
```

## Sync between instances

`canvas sync` copies content across configured instances (e.g. staging →
production). Both instances must exist in config and be authenticated.

```bash
canvas sync assignments prod 12345 staging 67890
canvas sync course prod 12345 staging 67890 --interactive   # interactive conflict resolution
```

## Raw API escape hatch

For endpoints without a dedicated command:

```bash
canvas api GET /api/v1/courses/123/todo
canvas api POST /api/v1/accounts/1/courses -d '{"course":{"name":"New Course"}}'
canvas api GET /api/v1/users -q "search_term=john" --paginate
```

`--dry-run` works here too — use it to show the user the exact request.

**Gotcha:** `canvas api` wraps the response in an envelope — the payload is
under `.body`, not at the top level:

```bash
# {"body": [...actual data...], "status_code": 200}
canvas api GET /api/v1/courses/123/tabs -o json | jq '.body'      # the data
canvas api GET /api/v1/courses/123/tabs -o json | jq '.body[0]'  # first item
```

Dedicated commands (`canvas modules list`, …) return the data directly,
without this wrapper.

## MCP server mode

The same binary is an MCP server exposing each command as a typed tool — use
it when a client wants structured tools instead of shell:

```bash
canvas mcp start            # STDIO MCP server
canvas mcp stream --port 8080   # HTTP
canvas mcp claude enable    # auto-configure Claude Desktop (also: cursor, vscode)
```

The skill (shell) and MCP modes can coexist; prefer the shell when you can run
commands directly.

## Errors & troubleshooting

- `canvas doctor` diagnoses install/auth/connectivity in one shot.
- `401` → re-auth (`canvas auth login` / check `CANVAS_TOKEN`); `403` → missing
  permission or masquerade (`--as-user`) not allowed; `404` → wrong ID or wrong
  instance.
- Rate limits are handled automatically (adaptive throttling + retries); for
  long batch jobs in env-auth mode you can tune `CANVAS_REQUESTS_PER_SEC`.
- Stale data? Responses are cached — add `--no-cache` or run
  `canvas cache clear`.
- Add `-v/--verbose` to see request logging; `--quiet` for clean script output.

## More

Full docs: https://jjuanrivvera.github.io/canvas-cli/ . Condensed references
ship alongside this skill in `references/canvas-commands.md`,
`references/auth-and-config.md`, and `references/output-and-filtering.md`.
