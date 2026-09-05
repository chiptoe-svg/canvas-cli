package commands

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/commands/internal/options"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
	"github.com/chiptoe-svg/canvas-cli/internal/localtime"
	"github.com/chiptoe-svg/canvas-cli/internal/output"
)

// scheduleNow is the clock behind today/tomorrow/this sunday. Tests pin it.
var scheduleNow = time.Now

func init() {
	rootCmd.AddCommand(newScheduleCmd())
}

func newScheduleCmd() *cobra.Command {
	opts := &options.ScheduleOptions{}

	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Set available/due/closed times on quizzes and assignments, in local time",
		Long: `Set the three availability times Canvas keeps on every quiz and
assignment — available from (unlock_at), due (due_at) and closed (lock_at)
— on one item (--id) or on every item whose title matches (--match), in
your local time zone.

Times are read in the local zone (see --timezone) and shown in both local
time and UTC. A time-only value such as "4:50pm" applies to each matched
item's OWN existing date for that field (falling back to its due date,
then available date, then closed date, and refusing an item with no date
at all) unless --date is given, which moves every matched item to that
one day instead; a date-only --due or --closed means 11:59 PM of that day
and a date-only --available means 12:00 AM.

Quizzes are updated through the quiz endpoint, plain assignments through
the assignment endpoint. A quiz-backed assignment is always updated through
its quiz so the two stay consistent, and no item is written twice. After
merging with the item's current values the order available ≤ due ≤ closed
must hold on every item; otherwise nothing is written.

Every item is read back after the write and its before/after printed; the
command exits non-zero if any read-back does not match. --dry-run prints
the full plan and the requests it would send and writes nothing.

Examples:
  # One quiz: available at 4:00, due and closed at 4:50, on its own date
  canvas schedule --course-id 123 --id 456 --type quiz --available 4:00pm --due 4:50pm --closed 4:50pm

  # Every attendance quiz gets these times, each on its own existing date
  canvas schedule --course-id 123 --match attendance --type quiz --available 4pm --due 4:50pm --closed 4:50pm --dry-run
  canvas schedule --course-id 123 --match attendance --type quiz --available 4pm --due 4:50pm --closed 4:50pm --force

  # Move every attendance quiz to one specific day instead
  canvas schedule --course-id 123 --match attendance --type quiz --date 2026-09-09 --available 4pm --due 4:50pm --closed 4:50pm --force

  # An assignment due on a date (11:59 PM), in a specific zone
  canvas schedule --course-id 123 --match "course lineup" --due 9/9/26 --timezone America/New_York

  # Remove the close date
  canvas schedule --course-id 123 --id 789 --type assignment --clear closed`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.DryRun = dryRun
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			// The plan needs real reads even under --dry-run; the command
			// implements its own preview and never writes in that mode.
			client.SetDryRun(false)
			// The plan merges with current values and the read-back must be
			// live, so no cached GET may be used.
			client.SetCacheEnabled(false)
			return runSchedule(cmd.Context(), client, opts, cmd.OutOrStdout(), cmd.ErrOrStderr(), cmd.InOrStdin())
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.ID, "id", 0, "Quiz or assignment ID (see --type)")
	cmd.Flags().StringVar(&opts.Match, "match", "", "Every quiz/assignment whose title contains this text (case-insensitive) or matches /regex/")
	cmd.Flags().StringVar(&opts.Type, "type", options.ScheduleTypeAll, "Item kind: quiz, assignment or all")
	cmd.Flags().StringVar(&opts.Available, "available", "", "Available-from time (unlock_at), local")
	cmd.Flags().StringVar(&opts.Due, "due", "", "Due time (due_at), local")
	cmd.Flags().StringVar(&opts.Closed, "closed", "", "Close time (lock_at), local")
	cmd.Flags().StringVar(&opts.Date, "date", "", "Move every matched item's time-only values to this day instead of each item's own date (2026-09-09, 9/9/26, today, this sunday)")
	cmd.Flags().StringSliceVar(&opts.Clear, "clear", nil, "Clear a date: available, due or closed (repeatable)")
	addTimezoneFlag(cmd, &opts.Timezone)
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip the confirmation prompt shown for --match")
	mustMarkRequired(cmd, "course-id")

	return cmd
}

// scheduleDates are the three Canvas timestamps; nil means not set.
type scheduleDates struct {
	UnlockAt *time.Time `json:"unlock_at" yaml:"unlock_at"`
	DueAt    *time.Time `json:"due_at" yaml:"due_at"`
	LockAt   *time.Time `json:"lock_at" yaml:"lock_at"`
}

// scheduleFields lists the fields in display order with their flag names.
var scheduleFields = []struct{ canvas, flag string }{
	{api.DateFieldUnlockAt, options.ScheduleFieldAvailable},
	{api.DateFieldDueAt, options.ScheduleFieldDue},
	{api.DateFieldLockAt, options.ScheduleFieldClosed},
}

func (d *scheduleDates) get(field string) *time.Time {
	switch field {
	case api.DateFieldUnlockAt:
		return d.UnlockAt
	case api.DateFieldDueAt:
		return d.DueAt
	default:
		return d.LockAt
	}
}

func (d *scheduleDates) set(field string, t *time.Time) {
	switch field {
	case api.DateFieldUnlockAt:
		d.UnlockAt = t
	case api.DateFieldDueAt:
		d.DueAt = t
	default:
		d.LockAt = t
	}
}

// scheduleItem is one quiz or assignment in the plan and, after the write,
// its evidence.
type scheduleItem struct {
	Kind         string         `json:"type" yaml:"type"` // quiz | assignment
	ID           int64          `json:"id" yaml:"id"`     // quiz id or assignment id
	AssignmentID int64          `json:"assignment_id,omitempty" yaml:"assignment_id,omitempty"`
	Title        string         `json:"title" yaml:"title"`
	HasOverrides bool           `json:"has_overrides,omitempty" yaml:"has_overrides,omitempty"`
	Before       scheduleDates  `json:"before" yaml:"before"`
	After        scheduleDates  `json:"after" yaml:"after"`
	ReadBack     *scheduleDates `json:"read_back,omitempty" yaml:"read_back,omitempty"`
	Changed      bool           `json:"changed" yaml:"changed"`
	Written      bool           `json:"written" yaml:"written"`
	Verified     string         `json:"verified" yaml:"verified"` // yes | no | -
	Mismatches   []string       `json:"mismatches,omitempty" yaml:"mismatches,omitempty"`
	Error        string         `json:"error,omitempty" yaml:"error,omitempty"`
}

func (it *scheduleItem) label() string {
	return fmt.Sprintf("%s %d %q", it.Kind, it.ID, it.Title)
}

// changes is the DateChanges to send: only fields that differ.
func (it *scheduleItem) changes() api.DateChanges {
	c := api.DateChanges{}
	for _, f := range scheduleFields {
		if !sameInstant(it.Before.get(f.canvas), it.After.get(f.canvas)) {
			c[f.canvas] = it.After.get(f.canvas)
		}
	}
	return c
}

func (it *scheduleItem) path(courseID int64) string {
	if it.Kind == options.ScheduleTypeQuiz {
		return fmt.Sprintf("/api/v1/courses/%d/quizzes/%d", courseID, it.ID)
	}
	return fmt.Sprintf("/api/v1/courses/%d/assignments/%d", courseID, it.ID)
}

func (it *scheduleItem) wrapper() string {
	if it.Kind == options.ScheduleTypeQuiz {
		return "quiz"
	}
	return "assignment"
}

type scheduleSummary struct {
	Matched    int `json:"matched" yaml:"matched"`
	Refused    int `json:"refused" yaml:"refused"`
	Changed    int `json:"changed" yaml:"changed"`
	Written    int `json:"written" yaml:"written"`
	Verified   int `json:"verified" yaml:"verified"`
	Mismatched int `json:"mismatched" yaml:"mismatched"`
	Failed     int `json:"failed" yaml:"failed"`
}

type scheduleResult struct {
	CourseID int64           `json:"course_id" yaml:"course_id"`
	Timezone string          `json:"timezone" yaml:"timezone"`
	DryRun   bool            `json:"dry_run" yaml:"dry_run"`
	Items    []scheduleItem  `json:"items" yaml:"items"`
	Summary  scheduleSummary `json:"summary" yaml:"summary"`
	loc      *time.Location  // for rendering
	// movedCount/movedDate report a uniform --date move, for the plan
	// header warning; movedCount is 0 unless --date was given.
	movedCount int
	movedDate  string
}

// scheduleRequest is what the flags asked for. Fields with an own date
// (from --date, or a value that names its own date) resolve to a concrete
// instant in set, applied to every item alike. A pure time-only value
// (e.g. "4:50pm") with no --date is deferred: perItemRaw holds the raw
// text, keyed by Canvas field name, and mergeScheduleDatesForItem combines
// it with each item's own existing date.
type scheduleRequest struct {
	set        scheduleDates
	clear      map[string]bool
	perItemRaw map[string]string
	// hasDate/dateCtx record a --date, for the plan's "moving N items" warning.
	hasDate bool
	dateCtx time.Time
}

func runSchedule(ctx context.Context, client *api.Client, opts *options.ScheduleOptions, out, errOut io.Writer, in io.Reader) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "schedule", map[string]interface{}{
		"course_id": opts.CourseID,
		"id":        opts.ID,
		"match":     opts.Match,
		"type":      opts.ItemType(),
		"dry_run":   opts.DryRun,
	})

	if output.FormatType(outputFormat) == output.FormatCSV {
		return fmt.Errorf("unsupported output format %q for schedule (table, json, yaml)", outputFormat)
	}
	tzName := resolveTimezone(opts.Timezone)
	loc, err := localtime.Location(tzName)
	if err != nil {
		return err
	}
	now := scheduleNow()

	req, err := resolveScheduleRequest(opts, tzName, now)
	if err != nil {
		return err
	}
	if !isStructuredOutput() {
		printScheduleRequest(out, req, loc)
	}

	items, err := resolveScheduleTargets(ctx, client, opts, errOut)
	if err != nil {
		logger.LogCommandError(ctx, "schedule", err, nil)
		return err
	}
	if len(items) == 0 {
		return fmt.Errorf("no %s in course %d matches %q", scheduleKindWord(opts.ItemType()), opts.CourseID, opts.Match)
	}

	result := &scheduleResult{CourseID: opts.CourseID, Timezone: loc.String(), DryRun: opts.DryRun, Items: items, loc: loc}
	result.Summary.Matched = len(items)
	var violations []string
	for i := range result.Items {
		it := &result.Items[i]
		it.Verified = "-"
		after, mergeErr := mergeScheduleDatesForItem(it.Before, req, tzName, now, it.Title)
		if mergeErr != nil {
			// This item has nothing to combine a time-only value with; it
			// is refused on its own and does not block the rest of the
			// plan (unlike an order violation, which does).
			it.Error = mergeErr.Error()
			it.Mismatches = []string{it.Error}
			it.Verified = "no"
			it.After = it.Before
			result.Summary.Refused++
			continue
		}
		it.After = after
		it.Changed = !sameDates(it.Before, it.After)
		if it.Changed {
			result.Summary.Changed++
		}
		if v := scheduleOrderViolation(it.After, loc); v != "" {
			violations = append(violations, fmt.Sprintf("%s: %s", it.label(), v))
		}
	}
	if len(violations) > 0 {
		return fmt.Errorf("refusing the plan: the order available ≤ due ≤ closed would not hold after the change; nothing was written\n  %s",
			strings.Join(violations, "\n  "))
	}
	if req.hasDate {
		result.movedCount = scheduleDateMoveCount(result.Items, req, loc)
		result.movedDate = req.dateCtx.Format("2006-01-02")
	}

	if opts.DryRun {
		if err := printScheduleResult(out, result); err != nil {
			return err
		}
		if result.Summary.Refused > 0 {
			return fmt.Errorf("%d of %d matched items could not be planned; see above", result.Summary.Refused, result.Summary.Matched)
		}
		return nil
	}
	if result.Summary.Changed == 0 {
		if err := printScheduleResult(out, result); err != nil {
			return err
		}
		if result.Summary.Refused > 0 {
			return fmt.Errorf("%d of %d matched items could not be planned; see above", result.Summary.Refused, result.Summary.Matched)
		}
		fmt.Fprintln(out, "Nothing to change: every matched item already has the requested times.")
		return nil
	}

	if opts.Match != "" && !opts.Force {
		if !isStructuredOutput() {
			printSchedulePlan(out, result)
		}
		ok, err := confirmSchedule(out, in, result)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintln(out, "Cancelled.")
			return nil
		}
	}

	writeScheduleItems(ctx, client, opts.CourseID, result)

	logger.LogCommandComplete(ctx, "schedule", result.Summary.Written)
	if err := printScheduleResult(out, result); err != nil {
		return err
	}
	if result.Summary.Mismatched > 0 || result.Summary.Failed > 0 {
		return fmt.Errorf("%d of %d items did not read back as requested (%d write failures)",
			result.Summary.Mismatched+result.Summary.Failed, result.Summary.Changed, result.Summary.Failed)
	}
	if result.Summary.Refused > 0 {
		return fmt.Errorf("%d of %d matched items could not be planned; see above", result.Summary.Refused, result.Summary.Matched)
	}
	return nil
}

// resolveScheduleRequest turns the time flags into instants where the
// value carries its own date (--date, or a value that names one, such as
// "2026-09-09 4:50pm" or "tomorrow"). A pure time-only value with no
// --date (e.g. "4:50pm") cannot be resolved yet — it depends on each
// matched item's own existing date — so it is deferred into perItemRaw for
// mergeScheduleDatesForItem to combine per item.
func resolveScheduleRequest(opts *options.ScheduleOptions, tzName string, now time.Time) (*scheduleRequest, error) {
	req := &scheduleRequest{perItemRaw: map[string]string{}}
	var err error
	if req.clear, err = opts.ClearSet(); err != nil {
		return nil, err
	}

	if opts.Date != "" {
		p, err := localtime.Parse(opts.Date, localtime.Options{Timezone: tzName, Now: now})
		if err != nil {
			return nil, fmt.Errorf("invalid --date: %w", err)
		}
		if !p.HasDate || p.HasTime {
			return nil, fmt.Errorf("invalid --date %q: give a calendar day only (2026-09-09, 9/9/26, today, this sunday)", opts.Date)
		}
		req.hasDate = true
		req.dateCtx = p.Local
	}

	for _, f := range []struct{ flag, canvas, value string }{
		{options.ScheduleFieldAvailable, api.DateFieldUnlockAt, opts.Available},
		{options.ScheduleFieldDue, api.DateFieldDueAt, opts.Due},
		{options.ScheduleFieldClosed, api.DateFieldLockAt, opts.Closed},
	} {
		if f.value == "" {
			continue
		}
		// Classify first: does the value name its own date, or is it a
		// pure time? Without --date, DateContext is left zero so Parse
		// defaults it to today — only p.HasDate is used from this call
		// when deferring; the instant itself is recomputed per item.
		p, err := localtime.Parse(f.value, localtime.Options{Timezone: tzName, Now: now, DateContext: req.dateCtx})
		if err != nil {
			return nil, fmt.Errorf("invalid --%s: %w", f.flag, err)
		}
		if !req.hasDate && !p.HasDate {
			req.perItemRaw[f.canvas] = f.value
			continue
		}
		// A date-only due/close means the end of that day, as in the Canvas UI.
		if !p.HasTime && f.canvas != api.DateFieldUnlockAt {
			p = p.EndOfDay()
		}
		t := p.Time
		req.set.set(f.canvas, &t)
	}
	return req, nil
}

func printScheduleRequest(out io.Writer, req *scheduleRequest, loc *time.Location) {
	for _, f := range scheduleFields {
		if t := req.set.get(f.canvas); t != nil {
			fmt.Fprintf(out, "%-10s %s\n", f.flag+":", localtime.Describe(*t, loc))
		} else if req.clear[f.flag] {
			fmt.Fprintf(out, "%-10s (clear)\n", f.flag+":")
		} else if raw, ok := req.perItemRaw[f.canvas]; ok {
			fmt.Fprintf(out, "%-10s %s (each item's own date)\n", f.flag+":", raw)
		}
	}
}

func scheduleKindWord(kind string) string {
	switch kind {
	case options.ScheduleTypeQuiz:
		return "quiz"
	case options.ScheduleTypeAssignment:
		return "assignment"
	default:
		return "quiz or assignment"
	}
}

// resolveScheduleTargets finds the items to change. Quiz-backed assignments
// are represented by their quiz, so each Canvas object appears once.
func resolveScheduleTargets(ctx context.Context, client *api.Client, opts *options.ScheduleOptions, errOut io.Writer) ([]scheduleItem, error) {
	quizzes := api.NewQuizzesService(client)
	assignments := api.NewAssignmentsService(client)
	kind := opts.ItemType()

	if opts.ID != 0 {
		return resolveScheduleByID(ctx, quizzes, assignments, opts.CourseID, opts.ID, kind, errOut)
	}

	match, err := options.CompileNameMatch("--match", opts.Match)
	if err != nil {
		return nil, err
	}
	var items []scheduleItem
	seenQuiz := map[int64]bool{}

	if kind != options.ScheduleTypeAssignment {
		// Canvas API: GET /api/v1/courses/:course_id/quizzes — every classic
		// quiz, including practice quizzes and surveys that have no
		// assignment. https://canvas.instructure.com/doc/api/quizzes.html#method.quizzes/quizzes_api.index
		list, err := quizzes.List(ctx, opts.CourseID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to list quizzes: %w", err)
		}
		for i := range list {
			q := &list[i]
			if !match(q.Title) {
				continue
			}
			seenQuiz[q.ID] = true
			items = append(items, itemFromQuiz(q))
		}
	}
	if kind != options.ScheduleTypeQuiz {
		// Canvas API: GET /api/v1/courses/:course_id/assignments — includes
		// the assignment behind every graded quiz (quiz_id set), which is
		// routed through the quiz. https://canvas.instructure.com/doc/api/assignments.html#method.assignments_api.index
		list, err := assignments.List(ctx, opts.CourseID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to list assignments: %w", err)
		}
		for i := range list {
			a := &list[i]
			if !match(a.Name) {
				continue
			}
			if isQuizAssignment(a) {
				if kind == options.ScheduleTypeAssignment {
					continue // --type assignment means plain assignments
				}
				if a.QuizID == 0 || seenQuiz[a.QuizID] {
					continue
				}
				// Matched through its assignment but not through the quiz
				// list (a title that differs, or a quiz the list omitted).
				seenQuiz[a.QuizID] = true
				items = append(items, itemFromQuizAssignment(a))
				continue
			}
			items = append(items, itemFromAssignment(a))
		}
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Kind != items[j].Kind {
			return items[i].Kind < items[j].Kind
		}
		if items[i].Title != items[j].Title {
			return items[i].Title < items[j].Title
		}
		return items[i].ID < items[j].ID
	})
	return items, nil
}

func resolveScheduleByID(ctx context.Context, quizzes *api.QuizzesService, assignments *api.AssignmentsService, courseID, id int64, kind string, errOut io.Writer) ([]scheduleItem, error) {
	var quizItem, assignmentItem *scheduleItem

	if kind != options.ScheduleTypeAssignment {
		// Canvas API: GET /api/v1/courses/:course_id/quizzes/:id
		q, err := quizzes.Get(ctx, courseID, id)
		if err != nil && !api.IsNotFoundError(err) {
			return nil, fmt.Errorf("failed to get quiz %d: %w", id, err)
		}
		if err == nil {
			it := itemFromQuiz(q)
			quizItem = &it
		}
	}
	if kind != options.ScheduleTypeQuiz {
		// Canvas API: GET /api/v1/courses/:course_id/assignments/:id
		a, err := assignments.Get(ctx, courseID, id, nil)
		if err != nil && !api.IsNotFoundError(err) {
			return nil, fmt.Errorf("failed to get assignment %d: %w", id, err)
		}
		if err == nil {
			var it scheduleItem
			if isQuizAssignment(a) {
				if a.QuizID == 0 {
					return nil, fmt.Errorf("assignment %d %q is a quiz assignment without a quiz_id; schedule it through the quiz", a.ID, a.Name)
				}
				it = itemFromQuizAssignment(a)
				if !quiet {
					fmt.Fprintf(errOut, "note: assignment %d %q is quiz %d; updating the quiz\n", a.ID, a.Name, a.QuizID)
				}
			} else {
				it = itemFromAssignment(a)
			}
			assignmentItem = &it
		}
	}

	switch {
	case quizItem != nil && assignmentItem != nil:
		if assignmentItem.Kind == options.ScheduleTypeQuiz && assignmentItem.ID == quizItem.ID {
			return []scheduleItem{*quizItem}, nil
		}
		return nil, fmt.Errorf("--id %d is ambiguous: quiz %d %q and assignment %d %q both exist; add --type quiz or --type assignment",
			id, quizItem.ID, quizItem.Title, assignmentItem.ID, assignmentItem.Title)
	case quizItem != nil:
		return []scheduleItem{*quizItem}, nil
	case assignmentItem != nil:
		return []scheduleItem{*assignmentItem}, nil
	}
	return nil, fmt.Errorf("no %s with ID %d in course %d. Use 'canvas quizzes list' or 'canvas assignments list' to see available items",
		scheduleKindWord(kind), id, courseID)
}

func isQuizAssignment(a *api.Assignment) bool {
	if a.IsQuizAssignment || a.QuizID > 0 {
		return true
	}
	for _, st := range a.SubmissionTypes {
		if st == "online_quiz" {
			return true
		}
	}
	return false
}

func itemFromQuiz(q *api.Quiz) scheduleItem {
	return scheduleItem{
		Kind:         options.ScheduleTypeQuiz,
		ID:           q.ID,
		AssignmentID: q.AssignmentID,
		Title:        q.Title,
		Before:       scheduleDates{UnlockAt: copyTime(q.UnlockAt), DueAt: copyTime(q.DueAt), LockAt: copyTime(q.LockAt)},
	}
}

func itemFromQuizAssignment(a *api.Assignment) scheduleItem {
	return scheduleItem{
		Kind:         options.ScheduleTypeQuiz,
		ID:           a.QuizID,
		AssignmentID: a.ID,
		Title:        a.Name,
		HasOverrides: a.HasOverrides,
		Before:       assignmentDates(a),
	}
}

func itemFromAssignment(a *api.Assignment) scheduleItem {
	return scheduleItem{
		Kind:         options.ScheduleTypeAssignment,
		ID:           a.ID,
		Title:        a.Name,
		HasOverrides: a.HasOverrides,
		Before:       assignmentDates(a),
	}
}

func assignmentDates(a *api.Assignment) scheduleDates {
	return scheduleDates{UnlockAt: nonZeroTime(a.UnlockAt), DueAt: nonZeroTime(a.DueAt), LockAt: nonZeroTime(a.LockAt)}
}

func quizDates(q *api.Quiz) scheduleDates {
	return scheduleDates{UnlockAt: copyTime(q.UnlockAt), DueAt: copyTime(q.DueAt), LockAt: copyTime(q.LockAt)}
}

func copyTime(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	u := t.UTC()
	return &u
}

func nonZeroTime(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	u := t.UTC()
	return &u
}

func sameInstant(a, b *time.Time) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Equal(*b)
}

func sameDates(a, b scheduleDates) bool {
	return sameInstant(a.UnlockAt, b.UnlockAt) && sameInstant(a.DueAt, b.DueAt) && sameInstant(a.LockAt, b.LockAt)
}

// fallbackDateSource picks the calendar day a per-item time-only value
// combines with: the field's own current value, else the item's due date,
// then its available date, then its closed date. nil means the item has
// no date at all to apply a time to.
func fallbackDateSource(before scheduleDates, field string) *time.Time {
	for _, f := range []string{field, api.DateFieldDueAt, api.DateFieldUnlockAt, api.DateFieldLockAt} {
		if t := before.get(f); t != nil {
			return t
		}
	}
	return nil
}

// mergeScheduleDatesForItem applies req on top of before: cleared fields
// become nil, concrete (--date or own-dated) values replace directly, and
// a deferred time-only value combines with the item's own existing date
// via fallbackDateSource — returning an error naming title when the item
// has no date at all for that field.
func mergeScheduleDatesForItem(before scheduleDates, req *scheduleRequest, tzName string, now time.Time, title string) (scheduleDates, error) {
	after := before
	for _, f := range scheduleFields {
		switch {
		case req.clear[f.flag]:
			after.set(f.canvas, nil)
		case req.set.get(f.canvas) != nil:
			after.set(f.canvas, copyTime(req.set.get(f.canvas)))
		default:
			raw, ok := req.perItemRaw[f.canvas]
			if !ok {
				continue
			}
			dateSrc := fallbackDateSource(before, f.canvas)
			if dateSrc == nil {
				return before, fmt.Errorf("%s: no existing date to apply a time to; pass --date", title)
			}
			p, err := localtime.Parse(raw, localtime.Options{Timezone: tzName, Now: now, DateContext: *dateSrc})
			if err != nil {
				return before, fmt.Errorf("invalid --%s: %w", f.flag, err)
			}
			t := p.Time
			after.set(f.canvas, &t)
		}
	}
	return after, nil
}

// scheduleDateMoveCount counts matched items where a --date move actually
// changes the calendar day of a field the item previously had set (a
// newly-created date on a previously-unset field is not a "move").
func scheduleDateMoveCount(items []scheduleItem, req *scheduleRequest, loc *time.Location) int {
	if !req.hasDate {
		return 0
	}
	count := 0
	for i := range items {
		it := &items[i]
		moved := false
		for _, f := range scheduleFields {
			if req.set.get(f.canvas) == nil {
				continue
			}
			before := it.Before.get(f.canvas)
			after := it.After.get(f.canvas)
			if before != nil && after != nil && !sameLocalDate(*before, *after, loc) {
				moved = true
				break
			}
		}
		if moved {
			count++
		}
	}
	return count
}

func sameLocalDate(a, b time.Time, loc *time.Location) bool {
	ay, am, ad := a.In(loc).Date()
	by, bm, bd := b.In(loc).Date()
	return ay == by && am == bm && ad == bd
}

// scheduleOrderViolation checks available ≤ due ≤ closed (and available ≤
// closed) among the dates that are set.
func scheduleOrderViolation(d scheduleDates, loc *time.Location) string {
	var problems []string
	if d.UnlockAt != nil && d.DueAt != nil && d.UnlockAt.After(*d.DueAt) {
		problems = append(problems, fmt.Sprintf("available %s is after due %s", localtime.FormatLocal(*d.UnlockAt, loc), localtime.FormatLocal(*d.DueAt, loc)))
	}
	if d.DueAt != nil && d.LockAt != nil && d.DueAt.After(*d.LockAt) {
		problems = append(problems, fmt.Sprintf("due %s is after closed %s", localtime.FormatLocal(*d.DueAt, loc), localtime.FormatLocal(*d.LockAt, loc)))
	}
	if d.UnlockAt != nil && d.LockAt != nil && d.UnlockAt.After(*d.LockAt) {
		problems = append(problems, fmt.Sprintf("available %s is after closed %s", localtime.FormatLocal(*d.UnlockAt, loc), localtime.FormatLocal(*d.LockAt, loc)))
	}
	return strings.Join(problems, "; ")
}

// writeScheduleItems performs the writes and read-backs, recording the
// evidence on each item. A failure on one item does not stop the others.
func writeScheduleItems(ctx context.Context, client *api.Client, courseID int64, result *scheduleResult) {
	quizzes := api.NewQuizzesService(client)
	assignments := api.NewAssignmentsService(client)

	for i := range result.Items {
		it := &result.Items[i]
		if !it.Changed {
			continue
		}
		changes := it.changes()
		var readBack scheduleDates
		var err error
		if it.Kind == options.ScheduleTypeQuiz {
			if _, err = quizzes.UpdateDates(ctx, courseID, it.ID, changes); err == nil {
				it.Written = true
				var q *api.Quiz
				if q, err = quizzes.Get(ctx, courseID, it.ID); err == nil {
					readBack = quizDates(q)
				}
			}
		} else {
			if _, err = assignments.UpdateDates(ctx, courseID, it.ID, changes); err == nil {
				it.Written = true
				var a *api.Assignment
				if a, err = assignments.Get(ctx, courseID, it.ID, nil); err == nil {
					readBack = assignmentDates(a)
				}
			}
		}
		if it.Written {
			result.Summary.Written++
		}
		if err != nil {
			if it.Written {
				it.Error = fmt.Sprintf("written, but the read-back failed: %v", err)
			} else {
				it.Error = fmt.Sprintf("write failed: %v", err)
			}
			it.Verified = "no"
			it.Mismatches = []string{it.Error}
			result.Summary.Failed++
			continue
		}
		it.ReadBack = &readBack
		for _, f := range scheduleFields {
			want, got := it.After.get(f.canvas), readBack.get(f.canvas)
			if !sameInstant(want, got) {
				it.Mismatches = append(it.Mismatches, fmt.Sprintf("%s read back %s, requested %s",
					f.flag, describeOrNone(got, result.loc), describeOrNone(want, result.loc)))
			}
		}
		if len(it.Mismatches) == 0 {
			it.Verified = "yes"
			result.Summary.Verified++
		} else {
			it.Verified = "no"
			result.Summary.Mismatched++
		}
	}
}

func describeOrNone(t *time.Time, loc *time.Location) string {
	if t == nil {
		return "—"
	}
	return localtime.Describe(*t, loc)
}

func formatLocalOrNone(t *time.Time, loc *time.Location) string {
	if t == nil {
		return "—"
	}
	return localtime.FormatLocal(*t, loc)
}

func formatUTCOrNone(t *time.Time) string {
	if t == nil {
		return "—"
	}
	return localtime.FormatUTC(*t)
}

func confirmSchedule(out io.Writer, in io.Reader, result *scheduleResult) (bool, error) {
	fmt.Fprintf(out, "About to update %d of %d matched items in course %d. Continue? [y/N]: ", result.Summary.Changed, result.Summary.Matched, result.CourseID)
	response, err := bufio.NewReader(in).ReadString('\n')
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("failed to read response: %w", err)
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

// printSchedulePlan prints the per-field old → new table and the requests
// that would be sent.
func printSchedulePlan(out io.Writer, result *scheduleResult) {
	loc := result.loc
	fmt.Fprintf(out, "Plan for course %d (times in %s):\n", result.CourseID, result.Timezone)
	if result.movedCount > 1 {
		fmt.Fprintf(out, "Moving %d items to %s\n", result.movedCount, result.movedDate)
	}
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE\tID\tTITLE\tFIELD\tOLD (local)\tNEW (local)\tNEW (UTC)")
	for i := range result.Items {
		it := &result.Items[i]
		title := it.Title
		if it.HasOverrides {
			title += " (has overrides; base dates only)"
		}
		first := true
		for _, f := range scheduleFields {
			before, after := it.Before.get(f.canvas), it.After.get(f.canvas)
			kind, id := it.Kind, fmt.Sprint(it.ID)
			if !first {
				kind, id, title = "", "", ""
			}
			first = false
			newLocal, newUTC := "(unchanged)", ""
			if !sameInstant(before, after) {
				newLocal, newUTC = formatLocalOrNone(after, loc), formatUTCOrNone(after)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", kind, id, title, f.flag, formatLocalOrNone(before, loc), newLocal, newUTC)
		}
	}
	w.Flush()
	for i := range result.Items {
		it := &result.Items[i]
		if it.Error != "" {
			fmt.Fprintf(out, "  %s\n", it.Error)
		}
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Requests:")
	for i := range result.Items {
		it := &result.Items[i]
		if !it.Changed {
			continue
		}
		body, _ := it.changes().Payload(it.wrapper())
		raw, _ := json.Marshal(body)
		fmt.Fprintf(out, "  PUT %s %s\n", it.path(result.CourseID), raw)
	}
	fmt.Fprintln(out)
}

// printScheduleResult prints the plan (dry-run or nothing to change) or the
// verification table, or the structured result.
func printScheduleResult(out io.Writer, result *scheduleResult) error {
	format := output.FormatType(outputFormat)
	switch format {
	case output.FormatJSON, output.FormatYAML:
		return output.WriteWithOptions(out, result, format, verbose)
	}

	if result.DryRun || (result.Summary.Written == 0 && result.Summary.Failed == 0) {
		if result.DryRun {
			fmt.Fprintln(out, "DRY RUN: no changes were made.")
		}
		printSchedulePlan(out, result)
		fmt.Fprintf(out, "Plan: %d matched, %d would change, 0 written.\n", result.Summary.Matched, result.Summary.Changed)
		return nil
	}

	loc := result.loc
	fmt.Fprintf(out, "Result for course %d (times in %s):\n", result.CourseID, result.Timezone)
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE\tID\tTITLE\tFIELD\tBEFORE (local)\tAFTER (read back)\tAFTER (UTC)\tVERIFIED")
	for i := range result.Items {
		it := &result.Items[i]
		if !it.Changed {
			continue
		}
		first := true
		for _, f := range scheduleFields {
			before := it.Before.get(f.canvas)
			var after *time.Time
			if it.ReadBack != nil {
				after = it.ReadBack.get(f.canvas)
			}
			kind, id, title, verified := it.Kind, fmt.Sprint(it.ID), it.Title, it.Verified
			if !first {
				kind, id, title, verified = "", "", "", ""
			}
			first = false
			afterLocal, afterUTC := "?", "?"
			if it.ReadBack != nil {
				afterLocal, afterUTC = formatLocalOrNone(after, loc), formatUTCOrNone(after)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", kind, id, title, f.flag, formatLocalOrNone(before, loc), afterLocal, afterUTC, verified)
		}
	}
	w.Flush()
	for i := range result.Items {
		it := &result.Items[i]
		if it.Verified == "no" {
			fmt.Fprintf(out, "  %s: %s\n", it.label(), strings.Join(it.Mismatches, "; "))
		}
	}
	fmt.Fprintln(out)
	fmt.Fprintf(out, "Done: %d matched, %d changed, %d written, %d verified, %d mismatched, %d failed.\n",
		result.Summary.Matched, result.Summary.Changed, result.Summary.Written, result.Summary.Verified, result.Summary.Mismatched, result.Summary.Failed)
	return nil
}
