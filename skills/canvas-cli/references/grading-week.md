# Grading week

Reading student work, scoring it, leaving feedback. This is judgment applied to
someone's education record: the discipline matters as much as the commands.

## Find the pile, read before you score

```bash
canvas users todo                                          # what Canvas says needs grading
canvas assignments list --course-id 123 --filter "Essay"   # assignment ids
canvas assignments get 456 --course-id 123 -o json
canvas submissions list --course-id 123 --assignment-id 456 \
  --workflow-state submitted --include submission_comments,user -o json
canvas submissions get --course-id 123 --assignment-id 456 --user-id 789 \
  --include submission_history,submission_comments,rubric_assessment -o json
```

`canvas users todo` carries the instructor's to-do items with needs-grading
counts; take ids from there, since a copied course keeps none of last term's. Re-read the rubric every time — a copied course keeps the old
point values and the criterion ids change. Pull the existing comments so you do
not repeat feedback already left, and read the **history**, not just the latest
attempt: students routinely upload a multi-part assignment across separate
attempts, so the current attempt is often missing files that arrived earlier.

## Download the files

```bash
canvas submissions download --course-id 123 --assignment-id 456 \
  --destination ./assignment-456-submissions --no-cache
```

Nothing is written to Canvas, but student work lands on the local disk. Get the
instructor's approval for the destination, keep it private, and upload that
folder nowhere. Files land under `user-<Canvas user id>/` with the attachment
id, so duplicate filenames cannot collide, and every attempt is fetched. The
`submission-download-manifest.json` lists every submission, marking text-entry
and URL ones as having no file. Re-running skips what is there, so an
interrupted run resumes safely; `--overwrite` replaces local student files and
needs approval. Report the downloaded/skipped/failed counts and the manifest
path — a run with any failure is not complete. Where layout or
composition is part of the grade read the PDF visually: text extraction drops
exactly the evidence a design critique rests on; for prose-only work convert
once locally (`pdftotext`, never an online converter).

## Propose, then post

Show the instructor the proposed total, the per-criterion breakdown and the
exact comment text; wait for a yes. Then write score and rubric in one update:

```bash
canvas submissions grade --course-id 123 --assignment-id 456 --user-id 789 \
  --score 14 --rubric _1234=4 --rubric _5678=0 \
  --rubric-comment _5678="no third review" --comment "…" --dry-run
# then the same line without --dry-run, once the instructor says yes
```

Criterion ids come from the assignment's own rubric. Pass one `--rubric
<criterion-id>=<points>` per criterion **including the ones scored zero** — an
omitted criterion keeps its old value instead of going to zero — and
`--rubric-comment` only alongside a matching `--rubric`. A rubric alone is a
complete grade; Canvas totals it. `--posted-grade` replaces `--score` for
letter or pass/fail, and `--excuse` excuses in the same call. The command then
reads the submission back and prints `grade: 88 → 95`, the posted
comment's id, author and first 80 characters, and `verified: yes`. That
read-back is what done means; `verified: no — <reason>` exits non-zero, so
report the line and never claim the work was graded. `-o json` carries
`{before, after, requested, comment, verified, mismatches}`. Feedback with no
score change is `canvas submissions add-comment … --text "…"`, read back the
same way, and `canvas submissions comments` lists what is already there. Write
feedback that names the student's actual evidence — the choice they made, its
effect, one concrete next step — and lead with what works.

## Bulk grading from a CSV

```bash
canvas submissions bulk-grade --course-id 123 --csv-file grades.csv --dry-run
canvas submissions bulk-grade --course-id 123 --csv-file grades.csv
```

Columns: `user_id,assignment_id,score,comment`. The first line is always read
as a header and skipped — a CSV whose first row is real data loses that
student. The score column is passed through as the posted grade, so a letter or
`Pass` works; the comment column is optional. Rows are applied one at a time
and each is read back: the line shows `<before> → <after>` and whether it
verified, and the run ends with `N graded, N verified, N mismatched`. Any
mismatch or error exits non-zero — report that summary line. Re-running the
same CSV is safe: rows already at the target score verify unchanged and an
identical comment is not posted twice. Show the instructor the whole `--dry-run`
first; it is the only cheap place to catch a wrong `user_id`.

## Excuse a student

```bash
canvas submissions excuse --course-id 123 --student "Ada Lovelace" --assignment "Quiz 3" --dry-run
canvas submissions excuse --course-id 123 --student 789 --assignment 456 --unexcuse
```

`--student` takes an id, exact name, sortable name ("Lovelace, Ada"), login or
SIS id, or a substring matching exactly one active student; `--assignment`
takes an assignment id, quiz id, name, or a unique substring (a quiz resolves
to the assignment Canvas grades it under). Zero or several matches are refused
with the candidates listed — pick one and re-run, never guess. It prints what
it resolved, then `excused: not excused → excused` and `verified: yes`. Use
`--force` (skips the prompt) only once the instructor has approved the names.

## What changed, and what you cannot export

```bash
canvas grades history --course-id 123 --start-date 2026-03-01 --end-date 2026-03-31
canvas activity list --writes
```

`canvas grades history` is Canvas's own gradebook history: who changed what,
when. `canvas activity list` is the local log of what this CLI did, when the
instructor has enabled it. Do not rebuild a gradebook from `submissions list` —
Canvas's *Export Entire Gradebook* applies calculated scores, section
membership, assignment visibility and the Points Possible row, and no
instructor-level API command reproduces it. Ask the instructor to export it and
hand you the CSV; reformat that file preserving rows, columns, values and
ordering, and say plainly that the export step is theirs.
