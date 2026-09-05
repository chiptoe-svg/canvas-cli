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

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/commands/internal/options"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
	"github.com/chiptoe-svg/canvas-cli/internal/output"
)

// submissionsMissingNow is the clock behind the default --due-before "now".
// Tests replace it so the fixture's past/future due dates are stable.
var submissionsMissingNow = time.Now

// missingItemsWidth caps the "missing items" column in table and markdown
// output; the JSON report always carries the full list.
const missingItemsWidth = 72

func newSubmissionsMissingCmd() *cobra.Command {
	opts := &options.SubmissionsMissingOptions{}

	cmd := &cobra.Command{
		Use:   "missing",
		Short: "Report students missing work (read-only)",
		Long: `Report which students are missing work in one or more courses.

Read-only: two paginated reads per course — the active student roster and the
course-wide submissions grid (students/submissions?student_ids[]=all with the
assignment embedded in every row). Nothing is written to Canvas.

Scope: --course-id (repeatable) or --all-active (every course where your own
enrollment is an active teacher or TA in a current term).

Assignments in scope, after the inclusion rules: --assignment-id (repeatable)
or --assignment-match (substring, or /regex/) on the assignment name. By default
only published, already-due assignments, quizzes and graded discussions count;
undated assignments are skipped unless --include-undated. An explicitly named
--assignment-id that a rule excludes is reported on stderr, not silently dropped.

"Missing" means no submission row, workflow_state unsubmitted, or Canvas'
missing flag — and not excused. Late-but-submitted work is not missing; it is
counted separately. Two populations are always kept apart: students with zero
submissions for every in-scope assignment, and students missing one or more
specific assignments. --zero-only shows just the first group.

Examples:
  canvas submissions missing --all-active
  canvas submissions missing --course-id 123 --course-id 456
  canvas submissions missing --course-id 123 --assignment-id 789
  canvas submissions missing --course-id 123 --assignment-match "/^Quiz/"
  canvas submissions missing --all-active --zero-only
  canvas submissions missing --course-id 123 --due-before 2026-03-01 --min-missing 2
  canvas submissions missing --course-id 123 -o markdown`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			// The global --limit caps every paginated read, which would silently
			// truncate the roster or the grid and misreport students as missing.
			if globalLimit > 0 {
				return fmt.Errorf("--limit is not supported by submissions missing: it would truncate the roster or submissions grid")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runSubmissionsMissing(cmd.Context(), client, opts, cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}

	cmd.Flags().Int64SliceVar(&opts.CourseIDs, "course-id", nil, "Course ID (repeatable)")
	cmd.Flags().BoolVar(&opts.AllActive, "all-active", false, "Every course where you are an active teacher or TA this term")
	cmd.Flags().Int64SliceVar(&opts.AssignmentIDs, "assignment-id", nil, "Only these assignment IDs (repeatable)")
	cmd.Flags().StringVar(&opts.AssignmentMatch, "assignment-match", "", "Only assignments whose name contains this text, or matches /regex/")
	cmd.Flags().BoolVar(&opts.PublishedOnly, "published-only", true, "Only published assignments")
	cmd.Flags().StringSliceVar(&opts.Types, "types", options.DefaultMissingTypes, "Assignment kinds to include: assignment, quiz, discussion")
	cmd.Flags().BoolVar(&opts.ExcludeZeroPoints, "exclude-zero-points", false, "Skip assignments worth 0 points")
	cmd.Flags().StringVar(&opts.DueBefore, "due-before", "now", "Only assignments due at or before this time (RFC 3339, YYYY-MM-DD, or \"now\")")
	cmd.Flags().StringVar(&opts.DueAfter, "due-after", "", "Only assignments due at or after this time (RFC 3339 or YYYY-MM-DD)")
	cmd.Flags().BoolVar(&opts.IncludeUndated, "include-undated", false, "Include assignments with no due date")
	cmd.Flags().BoolVar(&opts.IncludeInactive, "include-inactive", false, "Also consider inactive and completed student enrollments")
	cmd.Flags().BoolVar(&opts.ZeroOnly, "zero-only", false, "Only students with no submission for any in-scope assignment")
	cmd.Flags().IntVar(&opts.MinMissing, "min-missing", 0, "Only students missing at least this many assignments")

	return cmd
}

// --- report model (JSON/YAML shape) ---

type missingReport struct {
	Courses []missingCourseReport `json:"courses" yaml:"courses"`
}

type missingCourseInfo struct {
	ID         int64  `json:"id" yaml:"id"`
	Name       string `json:"name" yaml:"name"`
	CourseCode string `json:"course_code" yaml:"course_code"`
	Term       string `json:"term" yaml:"term"`
}

type missingAssignmentInfo struct {
	ID             int64      `json:"id" yaml:"id"`
	Name           string     `json:"name" yaml:"name"`
	Type           string     `json:"type" yaml:"type"`
	DueAt          *time.Time `json:"due_at" yaml:"due_at"`
	PointsPossible float64    `json:"points_possible" yaml:"points_possible"`
}

type missingStudentReport struct {
	UserID          int64    `json:"user_id" yaml:"user_id"`
	Name            string   `json:"name" yaml:"name"`
	SortableName    string   `json:"sortable_name" yaml:"sortable_name"`
	Missing         []int64  `json:"missing" yaml:"missing"`
	MissingNames    []string `json:"missing_names" yaml:"missing_names"`
	MissingCount    int      `json:"missing_count" yaml:"missing_count"`
	Late            int      `json:"late" yaml:"late"`
	ZeroSubmissions bool     `json:"zero_submissions" yaml:"zero_submissions"`
}

type missingCourseReport struct {
	Course                  missingCourseInfo       `json:"course" yaml:"course"`
	StudentsConsidered      int                     `json:"students_considered" yaml:"students_considered"`
	Assignments             []missingAssignmentInfo `json:"assignments" yaml:"assignments"`
	Students                []missingStudentReport  `json:"students" yaml:"students"`
	StudentsMissingAny      int                     `json:"students_missing_any" yaml:"students_missing_any"`
	StudentsZeroSubmissions int                     `json:"students_zero_submissions" yaml:"students_zero_submissions"`
}

// --- inclusion rules ---

type missingScope struct {
	publishedOnly     bool
	types             map[string]bool
	excludeZeroPoints bool
	dueBefore         time.Time
	dueAfter          time.Time // zero = unset
	includeUndated    bool
	match             func(string) bool
	explicit          map[int64]bool // non-empty when --assignment-id was given
}

func newMissingScope(opts *options.SubmissionsMissingOptions, now time.Time) (*missingScope, error) {
	types, err := opts.TypeSet()
	if err != nil {
		return nil, err
	}
	match, err := options.CompileAssignmentMatch(opts.AssignmentMatch)
	if err != nil {
		return nil, err
	}
	before, err := options.ParseDueBound(opts.DueBefore, now)
	if err != nil {
		return nil, fmt.Errorf("invalid --due-before: %w", err)
	}
	var after time.Time
	if opts.DueAfter != "" {
		if after, err = options.ParseDueBound(opts.DueAfter, now); err != nil {
			return nil, fmt.Errorf("invalid --due-after: %w", err)
		}
	}
	s := &missingScope{
		publishedOnly:     opts.PublishedOnly,
		types:             types,
		excludeZeroPoints: opts.ExcludeZeroPoints,
		dueBefore:         before,
		dueAfter:          after,
		includeUndated:    opts.IncludeUndated,
		match:             match,
		explicit:          map[int64]bool{},
	}
	for _, id := range opts.AssignmentIDs {
		s.explicit[id] = true
	}
	return s, nil
}

// exclude returns why a is out of scope, or "" when it is in scope.
func (s *missingScope) exclude(a *api.Assignment) string {
	if len(s.explicit) > 0 && !s.explicit[a.ID] {
		return "not in --assignment-id"
	}
	if s.publishedOnly && !a.Published {
		return "unpublished (see --published-only=false)"
	}
	kind := missingAssignmentKind(a)
	if kind == "" {
		return "not submittable online (submission type none/on_paper/not_graded)"
	}
	if !s.types[kind] {
		return "type " + kind + " not in --types"
	}
	if s.excludeZeroPoints && a.PointsPossible == 0 {
		return "worth 0 points (--exclude-zero-points)"
	}
	if !s.match(a.Name) {
		return "name does not match --assignment-match"
	}
	if a.DueAt.IsZero() {
		if s.includeUndated {
			return ""
		}
		return "no due date (see --include-undated)"
	}
	if a.DueAt.After(s.dueBefore) {
		return "not due until " + a.DueAt.Format(time.RFC3339) + " (see --due-before)"
	}
	if !s.dueAfter.IsZero() && a.DueAt.Before(s.dueAfter) {
		return "due before --due-after"
	}
	return ""
}

// missingAssignmentKind maps a Canvas assignment to one of the --types kinds
// from its submission_types / is_quiz_assignment (Assignments API:
// https://canvas.instructure.com/doc/api/assignments.html#Assignment). It
// returns "" for assignments students cannot submit online (submission types
// none, not_graded, on_paper only), which can never be "missing" a submission.
func missingAssignmentKind(a *api.Assignment) string {
	if a.IsQuizAssignment || a.QuizID > 0 {
		return options.MissingTypeQuiz
	}
	submittable := false
	for _, st := range a.SubmissionTypes {
		switch st {
		case "online_quiz":
			return options.MissingTypeQuiz
		case "discussion_topic":
			return options.MissingTypeDiscussion
		case "none", "not_graded", "on_paper":
		default:
			submittable = true
		}
	}
	if !submittable {
		return ""
	}
	return options.MissingTypeAssignment
}

// submissionIsMissing applies the report's definition of "missing" to one
// (student, assignment) cell of the grid. A nil cell is a student the grid
// returned no row for at all.
func submissionIsMissing(s *api.Submission) bool {
	if s == nil {
		return true
	}
	if s.ExcusedTLN {
		return false
	}
	if s.Missing {
		return true
	}
	return s.WorkflowState == "unsubmitted"
}

// isTestStudent recognizes the course's Test Student, which must never appear
// in the report whether it arrives by name or as a StudentViewEnrollment.
func isTestStudent(u *api.User) bool {
	if strings.EqualFold(strings.TrimSpace(u.Name), "Test Student") {
		return true
	}
	for _, e := range u.Enrollments {
		if e.Type == "StudentViewEnrollment" {
			return true
		}
	}
	return false
}

// --- run ---

func runSubmissionsMissing(ctx context.Context, client *api.Client, opts *options.SubmissionsMissingOptions, out, errOut io.Writer) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "submissions.missing", map[string]interface{}{
		"course_ids": opts.CourseIDs,
		"all_active": opts.AllActive,
		"zero_only":  opts.ZeroOnly,
	})

	now := submissionsMissingNow()
	scope, err := newMissingScope(opts, now)
	if err != nil {
		return err
	}
	progress := func(format string, args ...interface{}) {
		if verbose {
			fmt.Fprintf(errOut, format+"\n", args...)
		}
	}

	courses, err := resolveMissingCourses(ctx, client, opts, now, progress)
	if err != nil {
		logger.LogCommandError(ctx, "submissions.missing", err, nil)
		return err
	}

	report := missingReport{Courses: []missingCourseReport{}}
	reported := 0
	for i := range courses {
		cr, err := buildMissingCourseReport(ctx, client, &courses[i], opts, scope, progress, errOut)
		if err != nil {
			logger.LogCommandError(ctx, "submissions.missing", err, map[string]interface{}{"course_id": courses[i].ID})
			return err
		}
		reported += len(cr.Students)
		report.Courses = append(report.Courses, cr)
	}

	logger.LogCommandComplete(ctx, "submissions.missing", reported)
	return renderMissingReport(out, &report, opts)
}

// resolveMissingCourses returns the courses to report on, with term included.
func resolveMissingCourses(ctx context.Context, client *api.Client, opts *options.SubmissionsMissingOptions, now time.Time, progress func(string, ...interface{})) ([]api.Course, error) {
	svc := api.NewCoursesService(client)

	if !opts.AllActive {
		var courses []api.Course
		for _, id := range opts.CourseIDs {
			progress("course %d: fetching course", id)
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

	// Canvas docs (List your active courses,
	// https://canvas.instructure.com/doc/api/courses.html#method.courses.index):
	// enrollment_type takes one value, so teacher and ta are two reads;
	// enrollment_state=active drops concluded enrollments; include[]=term
	// embeds the term so the header and the current-term filter need no
	// further reads.
	seen := map[int64]bool{}
	var courses []api.Course
	for _, et := range []string{"teacher", "ta"} {
		progress("listing active %s courses", et)
		list, err := svc.List(ctx, &api.ListCoursesOptions{EnrollmentType: et, EnrollmentState: "active", Include: []string{"term"}})
		if err != nil {
			return nil, fmt.Errorf("failed to list %s courses: %w", et, err)
		}
		for _, c := range list {
			if seen[c.ID] || !courseInCurrentTerm(&c, now) {
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

// courseInCurrentTerm keeps courses whose term (when it has dates) contains
// now and that are not concluded.
func courseInCurrentTerm(c *api.Course, now time.Time) bool {
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

type missingCell struct{ user, assignment int64 }

// buildMissingCourseReport performs the two data reads for one course and
// derives the report.
func buildMissingCourseReport(ctx context.Context, client *api.Client, course *api.Course, opts *options.SubmissionsMissingOptions, scope *missingScope, progress func(string, ...interface{}), errOut io.Writer) (missingCourseReport, error) {
	cr := missingCourseReport{
		Course:      missingCourseInfo{ID: course.ID, Name: course.Name, CourseCode: course.CourseCode, Term: courseTermName(course)},
		Assignments: []missingAssignmentInfo{},
		Students:    []missingStudentReport{},
	}

	// Read 1: the roster. Canvas docs (List users in course,
	// https://canvas.instructure.com/doc/api/courses.html#method.courses.users):
	// enrollment_type[]=student, enrollment_state[] filters server-side;
	// include[]=enrollments lets the Test Student be recognised by its
	// StudentViewEnrollment even if it leaks through.
	states := []string{"active"}
	if opts.IncludeInactive {
		states = append(states, "inactive", "completed")
	}
	progress("course %d: fetching students", course.ID)
	users, err := api.NewCoursesService(client).ListUsers(ctx, course.ID, &api.ListCourseUsersOptions{
		EnrollmentType:  []string{"student"},
		EnrollmentState: states,
		Include:         []string{"enrollments"},
	})
	if err != nil {
		return cr, fmt.Errorf("course %d: failed to list students: %w", course.ID, err)
	}
	var roster []api.User
	seenUser := map[int64]bool{}
	for _, u := range users {
		if seenUser[u.ID] || isTestStudent(&u) {
			continue
		}
		seenUser[u.ID] = true
		roster = append(roster, u)
	}
	cr.StudentsConsidered = len(roster)

	// Read 2: the grid, one row per (student, assignment) Canvas has a
	// submission object for, with the assignment embedded (include[]=assignment,
	// Submissions API "List submissions for multiple assignments",
	// https://canvas.instructure.com/doc/api/submissions.html#method.submissions_api.for_students).
	// Canvas omits rows for students it holds no submission object for, so a
	// (student, in-scope assignment) pair with no row counts as missing below.
	progress("course %d: fetching submissions grid", course.ID)
	grid, err := api.NewSubmissionsService(client).ListForAllStudents(ctx, course.ID, opts.AssignmentIDs, &api.ListSubmissionsOptions{
		Include: []string{"assignment"},
		PerPage: 100,
	})
	if err != nil {
		return cr, fmt.Errorf("course %d: failed to list submissions: %w", course.ID, err)
	}

	// The in-scope assignment set is the union of embedded assignments seen
	// across the grid ...
	assignments := map[int64]*api.Assignment{}
	cells := map[missingCell]*api.Submission{}
	for i := range grid {
		s := &grid[i]
		if s.Assignment != nil {
			if _, ok := assignments[s.AssignmentID]; !ok {
				assignments[s.AssignmentID] = s.Assignment
			}
		}
		cells[missingCell{s.UserID, s.AssignmentID}] = s
	}
	// ... plus one metadata GET for an explicitly named assignment the grid
	// returned no rows for at all.
	for _, id := range opts.AssignmentIDs {
		if _, ok := assignments[id]; ok {
			continue
		}
		progress("course %d: assignment %d has no submission rows; fetching it", course.ID, id)
		a, err := api.NewAssignmentsService(client).Get(ctx, course.ID, id, nil)
		if err != nil {
			if api.IsNotFoundError(err) {
				return cr, fmt.Errorf("assignment %d not found in course %d", id, course.ID)
			}
			return cr, fmt.Errorf("course %d: failed to fetch assignment %d: %w", course.ID, id, err)
		}
		assignments[id] = a
	}

	var inScope []*api.Assignment
	for _, a := range assignments {
		reason := scope.exclude(a)
		if reason == "" {
			inScope = append(inScope, a)
			continue
		}
		if scope.explicit[a.ID] {
			fmt.Fprintf(errOut, "note: course %d: assignment %d %q excluded: %s\n", course.ID, a.ID, a.Name, reason)
		}
	}
	sort.Slice(inScope, func(i, j int) bool {
		a, b := inScope[i], inScope[j]
		if !a.DueAt.Equal(b.DueAt) {
			if a.DueAt.IsZero() != b.DueAt.IsZero() {
				return !a.DueAt.IsZero() // dated before undated
			}
			return a.DueAt.Before(b.DueAt)
		}
		if a.Name != b.Name {
			return a.Name < b.Name
		}
		return a.ID < b.ID
	})
	for _, a := range inScope {
		info := missingAssignmentInfo{ID: a.ID, Name: a.Name, Type: missingAssignmentKind(a), PointsPossible: a.PointsPossible}
		if !a.DueAt.IsZero() {
			due := a.DueAt
			info.DueAt = &due
		}
		cr.Assignments = append(cr.Assignments, info)
	}

	for _, u := range roster {
		sr := missingStudentReport{UserID: u.ID, Name: u.Name, SortableName: u.SortableName, Missing: []int64{}, MissingNames: []string{}}
		if sr.SortableName == "" {
			sr.SortableName = u.Name
		}
		for _, a := range inScope {
			cell := cells[missingCell{u.ID, a.ID}]
			if cell != nil && cell.Late {
				sr.Late++
			}
			if submissionIsMissing(cell) {
				sr.Missing = append(sr.Missing, a.ID)
				sr.MissingNames = append(sr.MissingNames, a.Name)
			}
		}
		sr.MissingCount = len(sr.Missing)
		if sr.MissingCount == 0 {
			continue
		}
		sr.ZeroSubmissions = sr.MissingCount == len(inScope)
		cr.StudentsMissingAny++
		if sr.ZeroSubmissions {
			cr.StudentsZeroSubmissions++
		}
		if sr.MissingCount < opts.MinMissing || (opts.ZeroOnly && !sr.ZeroSubmissions) {
			continue
		}
		cr.Students = append(cr.Students, sr)
	}
	sort.Slice(cr.Students, func(i, j int) bool {
		a, b := cr.Students[i], cr.Students[j]
		if a.MissingCount != b.MissingCount {
			return a.MissingCount > b.MissingCount
		}
		an, bn := strings.ToLower(a.SortableName), strings.ToLower(b.SortableName)
		if an != bn {
			return an < bn
		}
		return a.UserID < b.UserID
	})
	return cr, nil
}

func courseTermName(c *api.Course) string {
	if c.Term != nil && c.Term.Name != "" {
		return c.Term.Name
	}
	return ""
}

// --- rendering ---

func renderMissingReport(out io.Writer, report *missingReport, opts *options.SubmissionsMissingOptions) error {
	switch strings.ToLower(outputFormat) {
	case "table", "":
		return renderMissingTable(out, report, opts)
	case "markdown", "md":
		return renderMissingMarkdown(out, report, opts)
	case "csv":
		return renderMissingCSV(out, report)
	case "json", "yaml":
		return output.WriteWithOptions(out, report, output.FormatType(strings.ToLower(outputFormat)), verbose)
	default:
		return fmt.Errorf("unsupported output format %q for submissions missing (table, markdown, json, yaml, csv)", outputFormat)
	}
}

func missingCourseHeader(c *missingCourseInfo) string {
	h := c.Name
	if c.CourseCode != "" {
		h = c.CourseCode + " — " + c.Name
	}
	if c.Term != "" {
		h += " (" + c.Term + ")"
	}
	return h
}

func missingSummaryLine(cr *missingCourseReport, opts *options.SubmissionsMissingOptions) string {
	students := "active students"
	if opts.IncludeInactive {
		students = "students (incl. inactive)"
	}
	return fmt.Sprintf("%d %s · %d assignments in scope · %d students missing ≥1 · %d with zero submissions",
		cr.StudentsConsidered, students, len(cr.Assignments), cr.StudentsMissingAny, cr.StudentsZeroSubmissions)
}

// splitMissingPopulations separates the two populations the report never
// conflates: zero-submission students, and students missing some but not all.
func splitMissingPopulations(cr *missingCourseReport) (zero, partial []missingStudentReport) {
	for _, s := range cr.Students {
		if s.ZeroSubmissions {
			zero = append(zero, s)
		} else {
			partial = append(partial, s)
		}
	}
	return zero, partial
}

func truncateMissingItems(names []string, width int) string {
	joined := strings.Join(names, ", ")
	if len(joined) <= width {
		return joined
	}
	var b strings.Builder
	for i, n := range names {
		next := n
		if i > 0 {
			next = ", " + n
		}
		if b.Len()+len(next) > width-12 && i > 0 {
			fmt.Fprintf(&b, ", … (+%d more)", len(names)-i)
			return b.String()
		}
		b.WriteString(next)
	}
	return b.String()
}

func renderMissingTable(out io.Writer, report *missingReport, opts *options.SubmissionsMissingOptions) error {
	if len(report.Courses) == 0 {
		fmt.Fprintln(out, "No courses in scope")
		return nil
	}
	for i := range report.Courses {
		cr := &report.Courses[i]
		if i > 0 {
			fmt.Fprintln(out)
		}
		fmt.Fprintln(out, missingCourseHeader(&cr.Course))
		fmt.Fprintln(out, missingSummaryLine(cr, opts))
		zero, partial := splitMissingPopulations(cr)
		if len(zero) == 0 && len(partial) == 0 {
			fmt.Fprintln(out, "No students missing work")
			continue
		}
		if len(zero) > 0 || opts.ZeroOnly {
			fmt.Fprintf(out, "\nZero submissions — nothing turned in for any of the %d in-scope assignments (%d students)\n", len(cr.Assignments), len(zero))
			if err := writeMissingRows(out, zero, len(cr.Assignments)); err != nil {
				return err
			}
		}
		if !opts.ZeroOnly && len(partial) > 0 {
			fmt.Fprintf(out, "\nMissing one or more — excludes the zero-submission students above (%d students)\n", len(partial))
			if err := writeMissingRows(out, partial, len(cr.Assignments)); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeMissingRows(out io.Writer, students []missingStudentReport, inScope int) error {
	if len(students) == 0 {
		fmt.Fprintln(out, "  (none)")
		return nil
	}
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  STUDENT\tID\tMISSING\tLATE\tMISSING ITEMS")
	for _, s := range students {
		items := truncateMissingItems(s.MissingNames, missingItemsWidth)
		if s.ZeroSubmissions {
			items = fmt.Sprintf("all %d", inScope)
		}
		fmt.Fprintf(w, "  %s\t%d\t%d\t%d\t%s\n", s.SortableName, s.UserID, s.MissingCount, s.Late, items)
	}
	return w.Flush()
}

func renderMissingMarkdown(out io.Writer, report *missingReport, opts *options.SubmissionsMissingOptions) error {
	if len(report.Courses) == 0 {
		fmt.Fprintln(out, "_No courses in scope._")
		return nil
	}
	for i := range report.Courses {
		cr := &report.Courses[i]
		if i > 0 {
			fmt.Fprintln(out)
		}
		fmt.Fprintf(out, "### %s\n", missingCourseHeader(&cr.Course))
		fmt.Fprintln(out, missingSummaryLine(cr, opts))
		zero, partial := splitMissingPopulations(cr)
		if len(zero) == 0 && len(partial) == 0 {
			fmt.Fprintln(out, "_No students missing work._")
			continue
		}
		if len(zero) > 0 || opts.ZeroOnly {
			fmt.Fprintf(out, "\n**Zero submissions (%d)** — nothing turned in for any of the %d in-scope assignments\n", len(zero), len(cr.Assignments))
			if len(zero) == 0 {
				fmt.Fprintln(out, "- _none_")
			}
			for _, s := range zero {
				fmt.Fprintf(out, "- %s (id %d)\n", s.SortableName, s.UserID)
			}
		}
		if !opts.ZeroOnly && len(partial) > 0 {
			fmt.Fprintf(out, "\n**Missing one or more (%d)** — excludes the zero-submission students above\n", len(partial))
			for _, s := range partial {
				late := ""
				if s.Late > 0 {
					late = fmt.Sprintf(", %d late", s.Late)
				}
				fmt.Fprintf(out, "- %s — %d missing%s: %s\n", s.SortableName, s.MissingCount, late, truncateMissingItems(s.MissingNames, missingItemsWidth))
			}
		}
	}
	return nil
}

func renderMissingCSV(out io.Writer, report *missingReport) error {
	w := csv.NewWriter(out)
	if err := w.Write([]string{"course_id", "course_code", "user_id", "name", "sortable_name", "missing_count", "late", "zero_submissions", "missing_ids", "missing_names"}); err != nil {
		return err
	}
	for _, cr := range report.Courses {
		for _, s := range cr.Students {
			ids := make([]string, len(s.Missing))
			for i, id := range s.Missing {
				ids[i] = strconv.FormatInt(id, 10)
			}
			row := []string{
				strconv.FormatInt(cr.Course.ID, 10), cr.Course.CourseCode,
				strconv.FormatInt(s.UserID, 10), s.Name, s.SortableName,
				strconv.Itoa(s.MissingCount), strconv.Itoa(s.Late), strconv.FormatBool(s.ZeroSubmissions),
				strings.Join(ids, ";"), strings.Join(s.MissingNames, ";"),
			}
			if err := w.Write(row); err != nil {
				return err
			}
		}
	}
	w.Flush()
	return w.Error()
}
