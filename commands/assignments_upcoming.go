package commands

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/localtime"
	"github.com/jjuanrivvera/canvas-cli/internal/output"
)

// assignmentsUpcomingNow is the report's clock. Tests pin it.
var assignmentsUpcomingNow = time.Now

func init() {
	assignmentsCmd.AddCommand(newAssignmentsUpcomingCmd())
}

func newAssignmentsUpcomingCmd() *cobra.Command {
	opts := &options.AssignmentsUpcomingOptions{}

	cmd := &cobra.Command{
		Use:   "upcoming",
		Short: "What is due within a window or by a date, per course, in local time (read-only)",
		Long: `List the assignments (and quizzes and graded discussions, through their
assignments) due after now and up to a limit, course by course, sorted by
due date, with the due time in your local zone, points, submission type
and published state. Nothing is written.

The limit is --within (36h, 10d, 2w — measured from now) or --by (a local
date or date/time; a date alone means the end of that day, so
--by "this sunday" covers all of Sunday). Undated items are left out unless
--include-undated lists them in a separate section. Unpublished items are
left out unless --published-only=false.

Examples:
  # "Is anything due this Sunday in GC 1010?"
  canvas assignments upcoming --course-id 123 --by "this sunday"

  # "Which assignments are due in the next 10 days in GC 4800?"
  canvas assignments upcoming --course-id 456 --within 10d

  # Every course you teach or TA this term, for chat
  canvas assignments upcoming --all-active --within 2w -o markdown`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			if globalLimit > 0 {
				return fmt.Errorf("--limit is not supported by assignments upcoming (it would truncate the assignment list)")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runAssignmentsUpcoming(cmd.Context(), client, opts, cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}

	cmd.Flags().Int64SliceVar(&opts.CourseIDs, "course-id", nil, "Course ID (repeatable)")
	cmd.Flags().BoolVar(&opts.AllActive, "all-active", false, "Every course you teach or TA in the current term")
	cmd.Flags().StringVar(&opts.Within, "within", "", "Window from now: <n>h, <n>d or <n>w (e.g. 10d, 2w)")
	cmd.Flags().StringVar(&opts.By, "by", "", "Limit as a local date/time (2026-09-09, 9/9/26, this sunday, tomorrow 5pm); a date alone means the end of that day")
	cmd.Flags().BoolVar(&opts.IncludeUndated, "include-undated", false, "Also list items with no due date, in a separate section")
	cmd.Flags().BoolVar(&opts.PublishedOnly, "published-only", true, "Only published items")
	addTimezoneFlag(cmd, &opts.Timezone)

	return cmd
}

type upcomingCourseInfo struct {
	ID         int64  `json:"id" yaml:"id"`
	Name       string `json:"name" yaml:"name"`
	CourseCode string `json:"course_code" yaml:"course_code"`
	Term       string `json:"term,omitempty" yaml:"term,omitempty"`
}

type upcomingItem struct {
	ID              int64      `json:"id" yaml:"id"`
	Name            string     `json:"name" yaml:"name"`
	Type            string     `json:"type" yaml:"type"` // assignment | quiz | discussion
	QuizID          int64      `json:"quiz_id,omitempty" yaml:"quiz_id,omitempty"`
	DueAt           *time.Time `json:"due_at" yaml:"due_at"`
	DueLocal        string     `json:"due_local,omitempty" yaml:"due_local,omitempty"`
	PointsPossible  float64    `json:"points_possible" yaml:"points_possible"`
	SubmissionTypes []string   `json:"submission_types" yaml:"submission_types"`
	Published       bool       `json:"published" yaml:"published"`
	HTMLURL         string     `json:"html_url,omitempty" yaml:"html_url,omitempty"`
}

type upcomingCourseReport struct {
	Course      upcomingCourseInfo `json:"course" yaml:"course"`
	Assignments []upcomingItem     `json:"assignments" yaml:"assignments"`
	Undated     []upcomingItem     `json:"undated,omitempty" yaml:"undated,omitempty"`
	UndatedSeen int                `json:"undated_count" yaml:"undated_count"`
	Count       int                `json:"count" yaml:"count"`
}

type upcomingReport struct {
	Now      time.Time              `json:"now" yaml:"now"`
	Limit    time.Time              `json:"limit" yaml:"limit"`
	Timezone string                 `json:"timezone" yaml:"timezone"`
	Courses  []upcomingCourseReport `json:"courses" yaml:"courses"`
	loc      *time.Location
}

func runAssignmentsUpcoming(ctx context.Context, client *api.Client, opts *options.AssignmentsUpcomingOptions, out, errOut io.Writer) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "assignments.upcoming", map[string]interface{}{
		"course_ids": opts.CourseIDs,
		"all_active": opts.AllActive,
		"within":     opts.Within,
		"by":         opts.By,
	})

	tzName := resolveTimezone(opts.Timezone)
	loc, err := localtime.Location(tzName)
	if err != nil {
		return err
	}
	now := assignmentsUpcomingNow().UTC()
	limit, err := upcomingLimit(opts, tzName, now)
	if err != nil {
		return err
	}

	progress := func(format string, args ...interface{}) {
		if verbose {
			fmt.Fprintf(errOut, format+"\n", args...)
		}
	}
	courses, err := upcomingResolveCourses(ctx, client, opts.CourseIDs, opts.AllActive, now, progress)
	if err != nil {
		logger.LogCommandError(ctx, "assignments.upcoming", err, nil)
		return err
	}

	report := &upcomingReport{Now: now, Limit: limit, Timezone: loc.String(), Courses: []upcomingCourseReport{}, loc: loc}
	assignments := api.NewAssignmentsService(client)
	total := 0
	for i := range courses {
		c := &courses[i]
		progress("course %d: listing assignments", c.ID)
		// Canvas API: GET /api/v1/courses/:course_id/assignments — every
		// assignment, quiz-backed and discussion-backed ones included;
		// due_at is the base due date (overrides are not consulted).
		// https://canvas.instructure.com/doc/api/assignments.html#method.assignments_api.index
		list, err := assignments.List(ctx, c.ID, nil)
		if err != nil {
			logger.LogCommandError(ctx, "assignments.upcoming", err, map[string]interface{}{"course_id": c.ID})
			return fmt.Errorf("failed to list assignments for course %d: %w", c.ID, err)
		}
		cr := buildUpcomingCourseReport(c, list, opts, now, limit, loc)
		total += cr.Count
		report.Courses = append(report.Courses, cr)
	}

	logger.LogCommandComplete(ctx, "assignments.upcoming", total)
	return renderUpcomingReport(out, report, opts)
}

// upcomingLimit resolves --within / --by to the end instant of the window.
func upcomingLimit(opts *options.AssignmentsUpcomingOptions, tzName string, now time.Time) (time.Time, error) {
	if opts.Within != "" {
		d, err := options.ParseWithin(opts.Within)
		if err != nil {
			return time.Time{}, err
		}
		return now.Add(d), nil
	}
	p, err := localtime.Parse(opts.By, localtime.Options{Timezone: tzName, Now: now})
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid --by: %w", err)
	}
	if !p.HasTime {
		p = p.EndOfDay()
	}
	if !p.Time.After(now) {
		return time.Time{}, fmt.Errorf("--by %s is not after now (%s)", p.Describe(), localtime.Describe(now, p.Location))
	}
	return p.Time, nil
}

// upcomingResolveCourses returns the courses to report on, with term
// included: the explicit --course-id list, or every course the caller
// teaches or TAs in the current term.
//
// This is the same resolution `submissions missing` performs
// (resolveMissingCourses); both branches carry a copy so each merges on
// its own. Fold them into one helper once both are in.
func upcomingResolveCourses(ctx context.Context, client *api.Client, courseIDs []int64, allActive bool, now time.Time, progress func(string, ...interface{})) ([]api.Course, error) {
	svc := api.NewCoursesService(client)

	if !allActive {
		var courses []api.Course
		for _, id := range courseIDs {
			progress("course %d: fetching course", id)
			// Canvas API: GET /api/v1/courses/:id?include[]=term
			c, err := svc.Get(ctx, id, []string{"term"})
			if err != nil {
				if api.IsNotFoundError(err) {
					return nil, fmt.Errorf("course with ID %d not found. Use 'canvas courses list' to see available courses", id)
				}
				return nil, fmt.Errorf("failed to fetch course %d: %w", id, err)
			}
			courses = append(courses, *c)
		}
		return courses, nil
	}

	// Canvas API (List your active courses,
	// https://canvas.instructure.com/doc/api/courses.html#method.courses.index):
	// enrollment_type takes one value, so teacher and ta are two reads;
	// enrollment_state=active drops concluded enrollments; include[]=term
	// embeds the term for the header and the current-term filter.
	seen := map[int64]bool{}
	var courses []api.Course
	for _, et := range []string{"teacher", "ta"} {
		progress("listing active %s courses", et)
		list, err := svc.List(ctx, &api.ListCoursesOptions{EnrollmentType: et, EnrollmentState: "active", Include: []string{"term"}})
		if err != nil {
			return nil, fmt.Errorf("failed to list %s courses: %w", et, err)
		}
		for _, c := range list {
			if seen[c.ID] || !upcomingCourseInTerm(&c, now) {
				continue
			}
			seen[c.ID] = true
			courses = append(courses, c)
		}
	}
	sort.Slice(courses, func(i, j int) bool {
		if courses[i].CourseCode != courses[j].CourseCode {
			return courses[i].CourseCode < courses[j].CourseCode
		}
		return courses[i].ID < courses[j].ID
	})
	return courses, nil
}

// upcomingCourseInTerm keeps courses whose term (when it has dates)
// contains now and that are not concluded.
func upcomingCourseInTerm(c *api.Course, now time.Time) bool {
	if c.WorkflowState == "completed" || c.WorkflowState == "deleted" {
		return false
	}
	if c.Term == nil {
		return true
	}
	if !c.Term.StartAt.IsZero() && c.Term.StartAt.After(now) {
		return false
	}
	if !c.Term.EndAt.IsZero() && c.Term.EndAt.Before(now) {
		return false
	}
	return true
}

// upcomingKind classifies an assignment by what backs it.
func upcomingKind(a *api.Assignment) string {
	if a.IsQuizAssignment || a.QuizID > 0 {
		return "quiz"
	}
	for _, st := range a.SubmissionTypes {
		switch st {
		case "online_quiz":
			return "quiz"
		case "discussion_topic":
			return "discussion"
		}
	}
	return "assignment"
}

func buildUpcomingCourseReport(c *api.Course, list []api.Assignment, opts *options.AssignmentsUpcomingOptions, now, limit time.Time, loc *time.Location) upcomingCourseReport {
	cr := upcomingCourseReport{
		Course:      upcomingCourseInfo{ID: c.ID, Name: c.Name, CourseCode: c.CourseCode, Term: upcomingTermName(c)},
		Assignments: []upcomingItem{},
	}
	for i := range list {
		a := &list[i]
		if opts.PublishedOnly && !a.Published {
			continue
		}
		item := upcomingItem{
			ID:              a.ID,
			Name:            a.Name,
			Type:            upcomingKind(a),
			QuizID:          a.QuizID,
			PointsPossible:  a.PointsPossible,
			SubmissionTypes: a.SubmissionTypes,
			Published:       a.Published,
			HTMLURL:         a.HTMLURL,
		}
		if a.DueAt.IsZero() {
			cr.UndatedSeen++
			if opts.IncludeUndated {
				cr.Undated = append(cr.Undated, item)
			}
			continue
		}
		due := a.DueAt.UTC()
		if !due.After(now) || due.After(limit) {
			continue
		}
		item.DueAt = &due
		item.DueLocal = localtime.FormatLocal(due, loc)
		cr.Assignments = append(cr.Assignments, item)
	}
	sort.SliceStable(cr.Assignments, func(i, j int) bool {
		if !cr.Assignments[i].DueAt.Equal(*cr.Assignments[j].DueAt) {
			return cr.Assignments[i].DueAt.Before(*cr.Assignments[j].DueAt)
		}
		return cr.Assignments[i].Name < cr.Assignments[j].Name
	})
	sort.SliceStable(cr.Undated, func(i, j int) bool { return cr.Undated[i].Name < cr.Undated[j].Name })
	cr.Count = len(cr.Assignments)
	return cr
}

func upcomingTermName(c *api.Course) string {
	if c.Term == nil {
		return ""
	}
	return c.Term.Name
}

func upcomingCourseHeader(c *upcomingCourseInfo) string {
	h := c.CourseCode
	if h == "" {
		h = fmt.Sprintf("course %d", c.ID)
	}
	if c.Name != "" && c.Name != c.CourseCode {
		h += " — " + c.Name
	}
	if c.Term != "" {
		h += " (" + c.Term + ")"
	}
	return h
}

func upcomingSummary(cr *upcomingCourseReport, report *upcomingReport, opts *options.AssignmentsUpcomingOptions) string {
	s := fmt.Sprintf("%d due by %s", cr.Count, localtime.FormatLocal(report.Limit, report.loc))
	if cr.UndatedSeen > 0 && !opts.IncludeUndated {
		s += fmt.Sprintf(" · %d undated not shown (--include-undated)", cr.UndatedSeen)
	}
	return s
}

func upcomingPoints(p float64) string {
	return strconv.FormatFloat(p, 'f', -1, 64)
}

func upcomingSubmission(types []string) string {
	if len(types) == 0 {
		return "-"
	}
	return strings.Join(types, ",")
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func renderUpcomingReport(out io.Writer, report *upcomingReport, opts *options.AssignmentsUpcomingOptions) error {
	switch output.FormatType(outputFormat) {
	case output.FormatJSON, output.FormatYAML:
		return output.WriteWithOptions(out, report, output.FormatType(outputFormat), verbose)
	case output.FormatCSV:
		return renderUpcomingCSV(out, report)
	case "markdown", "md":
		return renderUpcomingMarkdown(out, report, opts)
	case output.FormatTable:
		return renderUpcomingTable(out, report, opts)
	default:
		return fmt.Errorf("unsupported output format %q for assignments upcoming (table, markdown, json, yaml, csv)", outputFormat)
	}
}

func renderUpcomingTable(out io.Writer, report *upcomingReport, opts *options.AssignmentsUpcomingOptions) error {
	fmt.Fprintf(out, "Due after %s, up to %s (times in %s)\n",
		localtime.FormatLocal(report.Now, report.loc), localtime.FormatLocal(report.Limit, report.loc), report.Timezone)
	for i := range report.Courses {
		cr := &report.Courses[i]
		fmt.Fprintln(out)
		fmt.Fprintln(out, upcomingCourseHeader(&cr.Course))
		fmt.Fprintln(out, upcomingSummary(cr, report, opts))
		if cr.Count > 0 {
			w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "  DUE\tTITLE\tTYPE\tID\tPOINTS\tSUBMISSION\tPUBLISHED")
			for _, it := range cr.Assignments {
				fmt.Fprintf(w, "  %s\t%s\t%s\t%d\t%s\t%s\t%s\n", it.DueLocal, it.Name, it.Type, it.ID, upcomingPoints(it.PointsPossible), upcomingSubmission(it.SubmissionTypes), yesNo(it.Published))
			}
			w.Flush()
		}
		if opts.IncludeUndated && len(cr.Undated) > 0 {
			fmt.Fprintf(out, "Undated (%d)\n", len(cr.Undated))
			w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "  TITLE\tTYPE\tID\tPOINTS\tSUBMISSION\tPUBLISHED")
			for _, it := range cr.Undated {
				fmt.Fprintf(w, "  %s\t%s\t%d\t%s\t%s\t%s\n", it.Name, it.Type, it.ID, upcomingPoints(it.PointsPossible), upcomingSubmission(it.SubmissionTypes), yesNo(it.Published))
			}
			w.Flush()
		}
	}
	return nil
}

func renderUpcomingMarkdown(out io.Writer, report *upcomingReport, opts *options.AssignmentsUpcomingOptions) error {
	fmt.Fprintf(out, "Due by **%s** (%s)\n", localtime.FormatLocal(report.Limit, report.loc), report.Timezone)
	for i := range report.Courses {
		cr := &report.Courses[i]
		fmt.Fprintf(out, "\n### %s\n", upcomingCourseHeader(&cr.Course))
		fmt.Fprintf(out, "_%s_\n", upcomingSummary(cr, report, opts))
		for _, it := range cr.Assignments {
			fmt.Fprintf(out, "- **%s** — %s (%s, %s pts, %s%s)\n", it.DueLocal, it.Name, it.Type, upcomingPoints(it.PointsPossible), upcomingSubmission(it.SubmissionTypes), upcomingUnpublishedTag(it))
		}
		if opts.IncludeUndated && len(cr.Undated) > 0 {
			fmt.Fprintf(out, "\n**Undated (%d)**\n", len(cr.Undated))
			for _, it := range cr.Undated {
				fmt.Fprintf(out, "- %s (%s, %s pts, %s%s)\n", it.Name, it.Type, upcomingPoints(it.PointsPossible), upcomingSubmission(it.SubmissionTypes), upcomingUnpublishedTag(it))
			}
		}
	}
	return nil
}

func upcomingUnpublishedTag(it upcomingItem) string {
	if it.Published {
		return ""
	}
	return ", unpublished"
}

func renderUpcomingCSV(out io.Writer, report *upcomingReport) error {
	w := csv.NewWriter(out)
	if err := w.Write([]string{"course_code", "course_id", "assignment_id", "name", "type", "due_local", "due_at", "points_possible", "submission_types", "published"}); err != nil {
		return err
	}
	for i := range report.Courses {
		cr := &report.Courses[i]
		rows := append([]upcomingItem{}, cr.Assignments...)
		rows = append(rows, cr.Undated...)
		for _, it := range rows {
			dueAt := ""
			if it.DueAt != nil {
				dueAt = localtime.FormatUTC(*it.DueAt)
			}
			if err := w.Write([]string{cr.Course.CourseCode, fmt.Sprint(cr.Course.ID), fmt.Sprint(it.ID), it.Name, it.Type, it.DueLocal, dueAt,
				upcomingPoints(it.PointsPossible), strings.Join(it.SubmissionTypes, ";"), yesNo(it.Published)}); err != nil {
				return err
			}
		}
	}
	w.Flush()
	return w.Error()
}
