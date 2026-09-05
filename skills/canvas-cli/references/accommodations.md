# Accommodations

Extra time on quizzes, extra attempts, and later due dates for named students.
These are usually disability-services accommodations: apply exactly what the
instructor states, do not ask why a student has one, and do not record the
reason anywhere.

## Resolve the ids first, and read them back

Every command here takes numeric ids. Get them, then say the names out loud to
the instructor before writing anything — a wrong `--user-id` gives the wrong
student extra time, and the wrong `--quiz-id` changes the wrong exam.

```bash
canvas quizzes list --course-id 123
canvas quizzes list --course-id 123 --search "midterm"
canvas assignments list --course-id 123 --filter "Essay"
canvas users list --course-id 123 --enrollment-type student \
  --columns id,name,sortable_name
```

`canvas quizzes list` is the whole quiz surface of the course; run it first
when the instructor says "all her quizzes", so the list you act on is Canvas's
and not your memory of the syllabus.

## Extra time or attempts on one quiz

```bash
canvas quizzes extensions create --course-id 123 --quiz-id 456 --user-id 789 \
  --extra-time 30 --dry-run
canvas quizzes extensions create --course-id 123 --quiz-id 456 --user-id 789 \
  --extra-time 30
```

`--extra-time` is minutes added to that student's time limit on that quiz.
`--extra-attempts` grants extra tries. `--manually-unlocked` lets the student
in past the lock date. `--extend-from-now <minutes>` extends an attempt that is
running right now — it is the emergency lever for a student mid-quiz, and it
is measured from the moment the command runs, so do not queue it up in advance.

This is one quiz at a time. For a student who gets time-and-a-half on
everything, loop over the quiz ids from `canvas quizzes list` and dry-run the
whole set before applying any of it:

```bash
canvas quizzes list --course-id 123 -o json | jq -r '.[] | "\(.id)\t\(.title)"'
canvas quizzes extensions create --course-id 123 --quiz-id 456 --user-id 789 --extra-time 30 --dry-run
canvas quizzes extensions create --course-id 123 --quiz-id 457 --user-id 789 --extra-time 30 --dry-run
```

`canvas course-extensions quiz --course-id 123 --user-id 789 --extra-time 30`
is the course-wide form of the same thing, and
`canvas course-extensions assignment --course-id 123 --assignment-id 456
--user-id 789 --extra-attempts 2` grants extra submission attempts on one
assignment. Both take `--dry-run`.

## A later due date for named students

Quiz extensions do not move a due date. That is an assignment override, and it
is the right tool for "Ada gets until Friday".

```bash
canvas overrides list --course-id 123 --assignment-id 456
canvas overrides create --course-id 123 --assignment-id 456 \
  --student-ids "789" --title "Extended deadline" \
  --due-at "2026-03-20T23:59:00Z" --dry-run
canvas overrides create --course-id 123 --assignment-id 456 \
  --student-ids "789" --title "Extended deadline" --due-at "2026-03-20T23:59:00Z"
```

Read the existing overrides first: a student can belong to only one override
per assignment, so adding a second for the same student moves them off the
first rather than stacking. `--student-ids` is comma-separated (add every
student who shares the new date to one override); `--section-id` and
`--group-id` are
the other two audiences and you must pass exactly one of the three.
`--unlock-at` and `--lock-at` move the available and close times the same way.

Override times are **ISO 8601, not local wall-clock** — unlike `canvas
schedule`, which reads local time. Convert the instructor's time to UTC and
show them both values before writing, or ask which zone they mean. Getting this
wrong by a time zone is the most common way an accommodation lands on the wrong
day.

A student with an override no longer appears under the base due date in
`canvas assignments upcoming` — that command reports base dates only.

## Confirm it landed

None of the commands on this page print a `verified:` read-back, so read the
object back yourself and show the instructor the result:

```bash
canvas overrides list --course-id 123 --assignment-id 456 -o json \
  | jq '.[] | {id, title, student_ids, due_at}'
canvas courses settings effective-due-dates --course-id 123 -o json
canvas quizzes get 456 --course-id 123 -o json
```

`canvas courses settings effective-due-dates` is the check that matters for a
date change: it is the date Canvas will actually enforce per student, overrides
included. For a quiz extension, confirm with the instructor and, if they want
certainty, ask them to open the quiz's Moderate This Quiz page — the extension
shows there per student.

Report what you read back, not what you sent. If the read-back does not show
the accommodation, say so plainly and do not re-run the write blindly.
