# Mid-term check

Halfway through a term the instructor wants to know who is falling behind, who
has gone quiet, what lands next, and how to reach the students who need a
nudge. Everything up to the last section is read-only.

## Who is missing work

```bash
canvas submissions missing --course-id 123
canvas submissions missing --all-active -o markdown
canvas submissions missing --course-id 123 --zero-only
canvas submissions missing --course-id 123 --assignment-match "/^Quiz/"
canvas submissions missing --course-id 123 --due-before 2026-03-01 --min-missing 2
```

Read-only: two paginated reads per course, the active roster and the
course-wide submissions grid. Nothing is written.

It keeps two populations apart and you should too. **Zero submissions** —
students with nothing turned in for any in-scope assignment — is an advising
problem, often a student who never started the course. **Missing one or more**
specific assignments is a different conversation. `--zero-only` shows just the
first group.

By default it counts published assignments, quizzes and graded discussions that
are already due. Late-but-submitted work is not missing, and excused work is
not missing. Undated assignments need `--include-undated`; `--types
assignment,quiz,discussion` narrows the kinds, `--exclude-zero-points` drops
ungraded busywork, `--include-inactive` adds inactive and completed
enrollments. `-o markdown` is the shape to paste into chat; `-o json` gives
`courses[].students[]` with `user_id`, `name`, `missing` (assignment ids),
`missing_names`, `missing_count`, `late` and `zero_submissions`.

## Who has gone quiet

```bash
canvas analytics students --course-id 123 --sort score
canvas analytics students --course-id 123 --sort participations
canvas analytics activity --course-id 123
```

`analytics students` is one row per student: current score, page views,
participations. Sorted by `score` or `participations` ascending, the students
at risk are at the top. `analytics activity` is the course-wide participation
and page-view series over time — useful for "did anyone look at the module I
posted", not for judging an individual.

Read these as signals, not verdicts. Low page views can mean a student works
from a downloaded PDF. Do not report a student as disengaged; report what the
numbers say and let the instructor decide.

## What lands next

```bash
canvas assignments upcoming --course-id 123 --within 10d
canvas assignments upcoming --course-id 123 --by "this sunday"
canvas assignments upcoming --all-active --within 2w -o markdown
```

Read-only, per course, sorted by due date, with the due time in the
instructor's local zone plus points, submission type and published state. The
window is `--within 36h|10d|2w` (from now) or `--by <local date>`; a date alone
covers the whole of that day. Undated items only with `--include-undated`;
unpublished only with `--published-only=false`. Report the per-course summary
line and the rows. These are the base due dates — per-student and per-section
overrides are not consulted, so a student with an extension will still appear
under the original date.

## Message the students who need a nudge

This is the one write in this workflow, and it is irreversible: a Canvas
conversation lands in every recipient's inbox and their email.

```bash
# 1. Get the ids from the missing report, and read them back as names
canvas submissions missing --course-id 123 --min-missing 2 -o json \
  | jq '[.courses[].students[] | {user_id, name, missing_count, missing_names}]'

# 2. Preview: the exact recipients, subject, and body
canvas conversations create --recipients 789,790,791 \
  --context-code course_123 --subject "Checking in on GC 1010" \
  --body "…" --dry-run

# 3. Post only after the instructor approves the text and the list
canvas conversations create --recipients 789,790,791 \
  --context-code course_123 --subject "Checking in on GC 1010" --body "…"
```

Before the real run, tell the instructor the recipient **count** and read back
the names behind the ids — a stale id list is how a message reaches a student
who dropped the course. `--recipients` is comma-separated ids.

Without `--group`, each recipient gets a separate one-to-one conversation and
cannot see the others; with `--group` it is one thread they all share and every
name is visible to every student. Default to separate conversations for
anything about a student's own standing, and say which mode you are using when
you propose the message. `--context-code course_123` files the thread under the
course.

Write the message as the instructor, not about the student: name the specific
assignments, say what the student can still do, and give one next step. Never
paste the missing-work table into a message that goes to more than one student.
