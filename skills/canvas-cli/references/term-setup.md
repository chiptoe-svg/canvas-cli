# Term setup

Getting a course ready before students arrive: check what the shell holds, move
the dates, publish the content, staff the course, open the office-hour calendar.

## Start from what is there

```bash
canvas courses list --enrollment-type teacher
canvas courses get 123 --include syllabus_body,term -o json
canvas courses settings get --course-id 123
canvas courses settings late-policy --course-id 123
canvas courses settings tabs --course-id 123
```

A course copied from last term keeps the old dates, announcements and often the
old late policy. Read all four before proposing anything, and tell the
instructor what looks stale rather than fixing it silently. `canvas courses
settings effective-due-dates --course-id 123` shows the dates Canvas will apply
per student, overrides included.

## Move the dates

`canvas schedule` sets the three availability times Canvas keeps on quizzes and
assignments — available from (`unlock_at`), due (`due_at`), closed (`lock_at`)
— in the instructor's local zone. Pass the wall-clock time they said, never a
UTC conversion of your own. A date-only `--due`/`--closed` means 11:59 PM that
day; a time-only value applies to each matched item's own date unless `--date`
moves them all to one day.

```bash
# One item, by id — then the same line again without --dry-run
canvas schedule --course-id 123 --id 456 --type quiz \
  --available 4:00pm --due 4:50pm --closed 4:50pm --dry-run

# Every item whose title matches — ALWAYS dry-run this one first
canvas schedule --course-id 123 --match attendance --type quiz \
  --available 4pm --due 4:50pm --closed 4:50pm --dry-run
canvas schedule --course-id 123 --match attendance --type quiz \
  --available 4pm --due 4:50pm --closed 4:50pm --force

# Move one assignment to a date, or drop a close date
canvas schedule --course-id 123 --match "course lineup" --due 9/9/26 --dry-run
canvas schedule --course-id 123 --id 789 --type assignment --clear closed --dry-run
```

The `--dry-run` prints the whole plan: every matched item, each time old → new
in local and UTC, and the exact requests. Show it to the instructor — a
`--match` that swept in one extra quiz is only cheap to catch there. Without
`--force`, `--match` prompts anyway. Nothing is written unless available ≤ due
≤ closed holds on every item after merging with its current values. Every item
is read back; done means `yes` under VERIFIED on every row and a summary
ending `0 mismatched, 0 failed`.

## Publish the content

```bash
canvas modules list --course-id 123 --include items
canvas modules publish --course-id 123 9 --dry-run   # then again without it
canvas pages list --course-id 123
canvas pages create --course-id 123 --title "Syllabus" \
  --body "<p>…</p>" --published --dry-run
canvas pages update --course-id 123 syllabus --body "<p>…</p>" --dry-run
```

Publishing a module makes it visible to students immediately, so confirm the
module and its items before the real run. `canvas pages create --front-page`
sets the home page; page bodies are HTML, not markdown.

## Announce it

```bash
canvas announcements list --course-id 123
canvas announcements create --course-id 123 --title "Welcome to GC 1010" \
  --message "<p>Class starts Monday.</p>" --dry-run    # then without --dry-run
```

An announcement notifies the whole class the moment it posts: show the exact
title and message and get an explicit yes first. `--delayed-at
"2026-09-01T09:00:00Z"` stages one instead of posting now.

## Staff the course

```bash
canvas users search "jamie rivera"
canvas enrollments list --course-id 123 --type TaEnrollment
canvas enrollments create --course-id 123 --user-id 456 \
  --type TaEnrollment --state active --dry-run
canvas enrollments create --course-id 123 --user-id 456 --type TaEnrollment --state active
canvas enrollments create --course-id 123 --user-id 457 --type TeacherEnrollment --state active
```

This enrolls an account that already exists — the faculty build cannot create
Canvas users, so if the TA has none the registrar or a Canvas admin must make
one. `canvas users search` is Canvas's messageable-users search
(`/api/v1/search/recipients`), scoped to people the instructor already shares a
course or group with, so a new TA in none of their courses returns nothing.
Then ask the instructor for the TA's Canvas user id — the number in the URL of
that person's Canvas People page — or search a course they do share with
`canvas users list --course-id 123 --search "rivera"`. `enrollments create`
needs the numeric `--user-id` either way, and a wrong id gives a stranger
access to student records, so confirm the id and the name first. `--state
active` skips the invitation; `--section-id` limits it to one section.

## Office hours and interview slots

```bash
canvas appointment-groups create --context course_123 --title "Office Hours" \
  --description "15-minute slots" --participants-per-slot 1 --max-slots 2 \
  --location "Room 214" --dry-run          # then the same line without --dry-run
canvas appointment-groups list --scope manageable --context course_123 \
  --include appointments,participant_count
```

`--context` takes context codes (`course_123`) and is how the group is scoped
to a course; `--sub-context` narrows it to sections. Leave `--publish` off
until the slots exist: the individual time slots are added in the Canvas
calendar UI — the CLI creates the group, not its slots — so hand that step back
to the instructor, then publish. The `list` above confirms it: the group, its
appointments, and how many students have reserved.
