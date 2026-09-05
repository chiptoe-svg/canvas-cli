package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/options"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

// missingCanvas is a fake Canvas for `submissions missing`: two courses, a
// roster with an inactive student and a Test Student in each, a grid with
// every workflow state the report distinguishes, Link-header pagination on
// course 1's roster and grid, and a student with no grid rows at all.
//
// Fixture clock: now = 2026-03-15T12:00:00Z (see missingFixtureNow).
type missingCanvas struct {
	*httptest.Server
	mu       sync.Mutex
	requests []string // request URIs, excluding version detection
}

var missingFixtureNow = time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)

const (
	missingCourse1 = `{"id":1,"name":"Test Course","course_code":"TEST101","workflow_state":"available","term":{"id":9,"name":"Spring 2026","start_at":"2026-01-10T00:00:00Z","end_at":"2026-05-10T00:00:00Z"}}`
	missingCourse2 = `{"id":2,"name":"Biology","course_code":"BIO200","workflow_state":"available","term":{"id":9,"name":"Spring 2026","start_at":"2026-01-10T00:00:00Z","end_at":"2026-05-10T00:00:00Z"}}`
	// Course 3 belongs to a term that ended before the fixture clock.
	missingCourse3 = `{"id":3,"name":"Old Course","course_code":"OLD100","workflow_state":"available","term":{"id":8,"name":"Fall 2025","start_at":"2025-08-15T00:00:00Z","end_at":"2025-12-20T00:00:00Z"}}`
)

// Assignment objects as embedded by include[]=assignment.
var missingAssignments = map[int64]string{
	100: `{"id":100,"name":"Essay 1","published":true,"due_at":"2026-03-01T23:59:00Z","points_possible":10,"submission_types":["online_text_entry"],"is_quiz_assignment":false}`,
	101: `{"id":101,"name":"Quiz 1","published":true,"due_at":"2026-03-10T23:59:00Z","points_possible":5,"submission_types":["online_quiz"],"is_quiz_assignment":true,"quiz_id":77}`,
	102: `{"id":102,"name":"Draft essay","published":false,"due_at":"2026-03-01T23:59:00Z","points_possible":10,"submission_types":["online_upload"]}`,
	103: `{"id":103,"name":"Undated reading","published":true,"due_at":null,"points_possible":10,"submission_types":["online_text_entry"]}`,
	104: `{"id":104,"name":"Final project","published":true,"due_at":"2026-04-01T23:59:00Z","points_possible":100,"submission_types":["online_upload"]}`,
	105: `{"id":105,"name":"Attendance check","published":true,"due_at":"2026-03-01T23:59:00Z","points_possible":0,"submission_types":["online_url"]}`,
	106: `{"id":106,"name":"Discussion 1","published":true,"due_at":"2026-03-02T23:59:00Z","points_possible":10,"submission_types":["discussion_topic"]}`,
	200: `{"id":200,"name":"Lab 1","published":true,"due_at":"2026-02-20T23:59:00Z","points_possible":20,"submission_types":["online_upload"]}`,
	201: `{"id":201,"name":"Lab 2","published":true,"due_at":"2026-03-05T23:59:00Z","points_possible":20,"submission_types":["online_upload"]}`,
}

type missingRow struct {
	user, assignment int64
	state            string
	extra            string // additional JSON fields
}

// Course 1 grid. Students: 10 Ada (active), 11 Bob (active), 12 Cy (inactive),
// 13 Test Student, 14 Dee (active, NO rows at all).
var missingGrid1 = []missingRow{
	{10, 100, "graded", `"score":9,"submitted_at":"2026-02-28T10:00:00Z"`},
	{10, 101, "submitted", `"late":true,"submitted_at":"2026-03-11T10:00:00Z"`},
	{10, 102, "unsubmitted", ""},
	{10, 103, "unsubmitted", ""},
	{10, 104, "unsubmitted", ""},
	{10, 105, "unsubmitted", ""},
	{10, 106, "graded", `"excused":true`},
	{11, 100, "unsubmitted", ""},
	{11, 101, "graded", `"score":0,"missing":true`},
	{11, 102, "unsubmitted", ""},
	{11, 103, "unsubmitted", ""},
	{11, 104, "unsubmitted", ""},
	{11, 105, "unsubmitted", ""},
	// 11/106: no row at all
	{12, 100, "unsubmitted", ""},
	{12, 101, "unsubmitted", ""},
	{12, 105, "unsubmitted", ""},
	{12, 106, "unsubmitted", ""},
	{13, 100, "unsubmitted", ""},
	{13, 101, "unsubmitted", ""},
	{13, 105, "unsubmitted", ""},
	{13, 106, "unsubmitted", ""},
}

// Course 2 grid. Students: 20 Eve, 21 Frank (active), 22 Gus (inactive), 23 Test Student.
var missingGrid2 = []missingRow{
	{20, 200, "submitted", `"submitted_at":"2026-02-19T10:00:00Z"`},
	{20, 201, "submitted", `"submitted_at":"2026-03-04T10:00:00Z"`},
	{21, 200, "unsubmitted", ""},
	{21, 201, "submitted", `"submitted_at":"2026-03-04T10:00:00Z"`},
	{22, 200, "unsubmitted", ""},
	{22, 201, "unsubmitted", ""},
	{23, 200, "unsubmitted", ""},
	{23, 201, "unsubmitted", ""},
}

type missingUser struct {
	id             int64
	name, sortable string
	state          string // enrollment state
	enrollType     string
}

var missingUsers = map[int64][]missingUser{
	1: {
		{10, "Ada Lovelace", "Lovelace, Ada", "active", "StudentEnrollment"},
		{11, "Bob Builder", "Builder, Bob", "active", "StudentEnrollment"},
		{12, "Cy Inactive", "Inactive, Cy", "inactive", "StudentEnrollment"},
		{13, "Test Student", "Student, Test", "active", "StudentViewEnrollment"},
		{14, "Dee Ghost", "Ghost, Dee", "active", "StudentEnrollment"},
	},
	2: {
		{20, "Eve Adams", "Adams, Eve", "active", "StudentEnrollment"},
		{21, "Frank Zappa", "Zappa, Frank", "active", "StudentEnrollment"},
		{22, "Gus Gone", "Gone, Gus", "completed", "StudentEnrollment"},
		// Course 2's Test Student arrives with a plain StudentEnrollment, so it
		// can only be recognised by name.
		{23, "Test Student", "Student, Test", "active", "StudentEnrollment"},
	},
}

func (r missingRow) json() string {
	extra := ""
	if r.extra != "" {
		extra = "," + r.extra
	}
	return fmt.Sprintf(`{"id":%d,"user_id":%d,"assignment_id":%d,"workflow_state":%q,"assignment":%s%s}`,
		r.user*1000+r.assignment, r.user, r.assignment, r.state, missingAssignments[r.assignment], extra)
}

func (u missingUser) json() string {
	return fmt.Sprintf(`{"id":%d,"name":%q,"sortable_name":%q,"enrollments":[{"id":%d,"type":%q,"enrollment_state":%q,"course_id":1}]}`,
		u.id, u.name, u.sortable, u.id+500, u.enrollType, u.state)
}

func newMissingCanvas(t *testing.T) *missingCanvas {
	t.Helper()
	mc := &missingCanvas{}
	mc.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			_, _ = w.Write([]byte(`[]`))
			return
		}
		mc.mu.Lock()
		mc.requests = append(mc.requests, r.URL.RequestURI())
		mc.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		q := r.URL.Query()
		path := r.URL.Path

		switch {
		case path == "/api/v1/courses":
			switch q.Get("enrollment_type") {
			case "teacher":
				fmt.Fprintf(w, "[%s,%s]", missingCourse1, missingCourse2)
			case "ta":
				fmt.Fprintf(w, "[%s,%s]", missingCourse2, missingCourse3)
			default:
				fmt.Fprint(w, "[]")
			}
		case path == "/api/v1/courses/1":
			fmt.Fprint(w, missingCourse1)
		case path == "/api/v1/courses/2":
			fmt.Fprint(w, missingCourse2)
		case path == "/api/v1/courses/404":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"errors":[{"message":"The specified resource does not exist."}]}`)
		case strings.HasSuffix(path, "/users"):
			courseID, _ := strconv.ParseInt(strings.TrimSuffix(strings.TrimPrefix(path, "/api/v1/courses/"), "/users"), 10, 64)
			states := map[string]bool{}
			for _, s := range q["enrollment_state[]"] {
				states[s] = true
			}
			var users []string
			for _, u := range missingUsers[courseID] {
				if states[u.state] {
					users = append(users, u.json())
				}
			}
			// Course 1's roster is paginated: one user on page 1, the rest on page 2.
			if courseID == 1 && q.Get("page") != "2" {
				next := mc.URL + path + "?" + q.Encode() + "&page=2"
				w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="next"`, next))
				fmt.Fprintf(w, "[%s]", users[0])
				return
			}
			if courseID == 1 {
				users = users[1:]
			}
			fmt.Fprintf(w, "[%s]", strings.Join(users, ","))
		case strings.HasSuffix(path, "/students/submissions"):
			courseID, _ := strconv.ParseInt(strings.TrimSuffix(strings.TrimPrefix(path, "/api/v1/courses/"), "/students/submissions"), 10, 64)
			if q.Get("student_ids[]") != "all" {
				t.Errorf("grid request without student_ids[]=all: %s", r.URL.RequestURI())
			}
			grid := missingGrid1
			if courseID == 2 {
				grid = missingGrid2
			}
			only := map[int64]bool{}
			for _, id := range q["assignment_ids[]"] {
				n, _ := strconv.ParseInt(id, 10, 64)
				only[n] = true
			}
			var rows []string
			for _, row := range grid {
				if len(only) > 0 && !only[row.assignment] {
					continue
				}
				rows = append(rows, row.json())
			}
			// Course 1's grid is paginated: first 5 rows on page 1.
			if courseID == 1 && len(rows) > 5 && q.Get("page") != "2" {
				next := mc.URL + path + "?" + q.Encode() + "&page=2"
				w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="next"`, next))
				fmt.Fprintf(w, "[%s]", strings.Join(rows[:5], ","))
				return
			}
			if courseID == 1 && len(rows) > 5 {
				rows = rows[5:]
			}
			fmt.Fprintf(w, "[%s]", strings.Join(rows, ","))
		case strings.Contains(path, "/assignments/"):
			// Metadata fallback for an explicitly named assignment with no grid rows.
			id, _ := strconv.ParseInt(path[strings.LastIndex(path, "/")+1:], 10, 64)
			if body, ok := missingAssignments[id]; ok {
				fmt.Fprint(w, body)
				return
			}
			if id >= 5000 {
				fmt.Fprintf(w, `{"id":%d,"name":"Extra %d","published":true,"due_at":"2026-03-01T23:59:00Z","points_possible":1,"submission_types":["online_text_entry"]}`, id, id)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"errors":[{"message":"The specified resource does not exist."}]}`)
		default:
			t.Errorf("unexpected request %s", r.URL.RequestURI())
			http.NotFound(w, r)
		}
	}))
	return mc
}

// count returns how many recorded requests (across pages) hit the given
// path suffix for a course.
func (mc *missingCanvas) count(courseID int64, suffix string) int {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	n := 0
	prefix := fmt.Sprintf("/api/v1/courses/%d%s", courseID, suffix)
	for _, u := range mc.requests {
		if strings.HasPrefix(u, prefix+"?") || u == prefix {
			n++
		}
	}
	return n
}

func (mc *missingCanvas) find(substr string) []string {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	var out []string
	for _, u := range mc.requests {
		if strings.Contains(u, substr) {
			out = append(out, u)
		}
	}
	return out
}

// runMissing executes `submissions missing` against mc with the given output
// format, returning stdout and stderr.
func runMissing(t *testing.T, mc *missingCanvas, format string, args ...string) (string, string, error) {
	t.Helper()
	t.Setenv("CANVAS_URL", mc.URL)
	t.Setenv("CANVAS_TOKEN", "test-token")
	t.Setenv("CANVAS_REQUESTS_PER_SEC", "1000")

	oldFormat, oldNoCache, oldNow, oldVerbose := outputFormat, noCache, submissionsMissingNow, verbose
	outputFormat, noCache = format, true
	submissionsMissingNow = func() time.Time { return missingFixtureNow }
	t.Cleanup(func() {
		outputFormat, noCache, submissionsMissingNow, verbose = oldFormat, oldNoCache, oldNow, oldVerbose
	})

	cmd := newSubmissionsMissingCmd()
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs(args)
	err := cmd.ExecuteContext(context.Background())
	return out.String(), errOut.String(), err
}

func decodeMissing(t *testing.T, raw string) missingReport {
	t.Helper()
	var rep missingReport
	if err := json.Unmarshal([]byte(raw), &rep); err != nil {
		t.Fatalf("not json: %v\n%s", err, raw)
	}
	return rep
}

func assignmentIDs(cr missingCourseReport) []int64 {
	var ids []int64
	for _, a := range cr.Assignments {
		ids = append(ids, a.ID)
	}
	return ids
}

func findStudent(cr missingCourseReport, id int64) *missingStudentReport {
	for i := range cr.Students {
		if cr.Students[i].UserID == id {
			return &cr.Students[i]
		}
	}
	return nil
}

func equalIDs(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSubmissionsMissing_DefaultJSON(t *testing.T) {
	mc := newMissingCanvas(t)
	defer mc.Close()

	out, _, err := runMissing(t, mc, "json", "--course-id", "1", "--course-id", "2")
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	rep := decodeMissing(t, out)
	if len(rep.Courses) != 2 {
		t.Fatalf("courses = %d, want 2", len(rep.Courses))
	}
	c1, c2 := rep.Courses[0], rep.Courses[1]

	if c1.Course.CourseCode != "TEST101" || c1.Course.Term != "Spring 2026" {
		t.Errorf("course header = %+v", c1.Course)
	}
	// Default scope: published, due by now, dated, any points, all three kinds
	// — sorted by due date then name (100 and 105 share a due date).
	// Unpublished 102, undated 103, future 104 are out.
	if got := assignmentIDs(c1); !equalIDs(got, []int64{105, 100, 106, 101}) {
		t.Errorf("in-scope assignments = %v, want [105 100 106 101]", got)
	}
	if c1.Assignments[3].Type != "quiz" || c1.Assignments[2].Type != "discussion" || c1.Assignments[1].Type != "assignment" {
		t.Errorf("assignment kinds = %+v", c1.Assignments)
	}
	// Roster: Ada, Bob, Dee. Cy is inactive, 13 is the Test Student.
	if c1.StudentsConsidered != 3 {
		t.Errorf("students_considered = %d, want 3", c1.StudentsConsidered)
	}
	if c1.StudentsMissingAny != 3 || c1.StudentsZeroSubmissions != 2 {
		t.Errorf("missing_any=%d zero=%d, want 3 and 2", c1.StudentsMissingAny, c1.StudentsZeroSubmissions)
	}
	for _, id := range []int64{12, 13} {
		if findStudent(c1, id) != nil {
			t.Errorf("student %d must not be reported (inactive / Test Student)", id)
		}
	}
	// Ada: graded 100, late-but-submitted 101, excused 106 → only 105 missing.
	ada := findStudent(c1, 10)
	if ada == nil {
		t.Fatal("Ada missing from report")
	}
	if !equalIDs(ada.Missing, []int64{105}) || ada.MissingCount != 1 || ada.Late != 1 || ada.ZeroSubmissions {
		t.Errorf("Ada = %+v, want missing [105], late 1, not zero", *ada)
	}
	if len(ada.MissingNames) != 1 || ada.MissingNames[0] != "Attendance check" {
		t.Errorf("Ada missing_names = %v", ada.MissingNames)
	}
	// Bob: unsubmitted 100, missing:true 101, unsubmitted 105, no row 106 → zero.
	bob := findStudent(c1, 11)
	if bob == nil || !equalIDs(bob.Missing, []int64{105, 100, 106, 101}) || !bob.ZeroSubmissions {
		t.Errorf("Bob = %+v, want all four missing and zero_submissions", bob)
	}
	// Dee has no grid rows at all → zero submissions.
	dee := findStudent(c1, 14)
	if dee == nil || dee.MissingCount != 4 || !dee.ZeroSubmissions {
		t.Errorf("Dee = %+v, want 4 missing and zero_submissions", dee)
	}
	// Sorted by missing count desc, then sortable name: Builder, Ghost, Lovelace.
	if c1.Students[0].UserID != 11 || c1.Students[1].UserID != 14 || c1.Students[2].UserID != 10 {
		t.Errorf("order = %v", []int64{c1.Students[0].UserID, c1.Students[1].UserID, c1.Students[2].UserID})
	}

	// Course 2: Frank missing Lab 1 only; Eve clean and therefore absent.
	if c2.StudentsConsidered != 2 || c2.StudentsMissingAny != 1 || c2.StudentsZeroSubmissions != 0 {
		t.Errorf("course 2 summary = %+v", c2)
	}
	if findStudent(c2, 23) != nil {
		t.Error("course 2's Test Student (recognised by name only) must not be reported")
	}
	if len(c2.Students) != 1 || c2.Students[0].UserID != 21 || !equalIDs(c2.Students[0].Missing, []int64{200}) {
		t.Errorf("course 2 students = %+v", c2.Students)
	}

	// Exactly two data reads per course (plus pagination), no assignments list.
	if n := mc.count(1, "/users"); n != 2 {
		t.Errorf("course 1 roster requests = %d, want 2 (paginated)", n)
	}
	if n := mc.count(1, "/students/submissions"); n != 2 {
		t.Errorf("course 1 grid requests = %d, want 2 (paginated)", n)
	}
	if n := mc.count(2, "/users"); n != 1 {
		t.Errorf("course 2 roster requests = %d, want 1", n)
	}
	if n := mc.count(2, "/students/submissions"); n != 1 {
		t.Errorf("course 2 grid requests = %d, want 1", n)
	}
	if got := mc.find("/assignments"); len(got) != 0 {
		t.Errorf("no assignments read expected, got %v", got)
	}
	grid := mc.find("/students/submissions")[0]
	for _, want := range []string{"student_ids%5B%5D=all", "include%5B%5D=assignment", "per_page=100"} {
		if !strings.Contains(grid, want) {
			t.Errorf("grid request %q lacks %s", grid, want)
		}
	}
	if strings.Contains(grid, "assignment_ids") {
		t.Errorf("default grid must not filter assignment_ids: %s", grid)
	}
	roster := mc.find("/users")[0]
	for _, want := range []string{"enrollment_type%5B%5D=student", "enrollment_state%5B%5D=active", "include%5B%5D=enrollments"} {
		if !strings.Contains(roster, want) {
			t.Errorf("roster request %q lacks %s", roster, want)
		}
	}
	if strings.Contains(roster, "inactive") {
		t.Errorf("default roster must be active only: %s", roster)
	}
}

func TestSubmissionsMissing_TableOutput(t *testing.T) {
	mc := newMissingCanvas(t)
	defer mc.Close()

	out, _, err := runMissing(t, mc, "table", "--course-id", "1", "--course-id", "2")
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	for _, want := range []string{
		"TEST101 — Test Course (Spring 2026)",
		"3 active students · 4 assignments in scope · 3 students missing ≥1 · 2 with zero submissions",
		"Zero submissions — nothing turned in for any of the 4 in-scope assignments (2 students)",
		"Missing one or more — excludes the zero-submission students above (1 students)",
		"Builder, Bob",
		"Ghost, Dee",
		"Lovelace, Ada",
		"Attendance check",
		"BIO200 — Biology (Spring 2026)",
		"2 active students · 2 assignments in scope · 1 students missing ≥1 · 0 with zero submissions",
		"Zappa, Frank",
		"Lab 1",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("table lacks %q:\n%s", want, out)
		}
	}
	for _, banned := range []string{"Test Student", "Student, Test", "Inactive, Cy", "Adams, Eve", "Final project", "Draft essay", "Undated reading"} {
		if strings.Contains(out, banned) {
			t.Errorf("table must not contain %q:\n%s", banned, out)
		}
	}
	// The zero-submission section comes before the per-assignment section.
	if strings.Index(out, "Zero submissions") > strings.Index(out, "Missing one or more") {
		t.Errorf("zero-submission section must come first:\n%s", out)
	}
	// Ada belongs to the second section only.
	zeroSection := out[strings.Index(out, "Zero submissions"):strings.Index(out, "Missing one or more")]
	if strings.Contains(zeroSection, "Lovelace") {
		t.Errorf("Ada is not a zero-submission student:\n%s", zeroSection)
	}
}

func TestSubmissionsMissing_MarkdownOutput(t *testing.T) {
	mc := newMissingCanvas(t)
	defer mc.Close()

	out, _, err := runMissing(t, mc, "markdown", "--course-id", "1")
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	for _, want := range []string{
		"### TEST101 — Test Course (Spring 2026)",
		"**Zero submissions (2)**",
		"- Builder, Bob (id 11)",
		"- Ghost, Dee (id 14)",
		"**Missing one or more (1)**",
		"- Lovelace, Ada — 1 missing, 1 late: Attendance check",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("markdown lacks %q:\n%s", want, out)
		}
	}
	for _, line := range strings.Split(out, "\n") {
		if len(line) > 120 {
			t.Errorf("markdown line too long for chat (%d chars): %s", len(line), line)
		}
	}
}

func TestSubmissionsMissing_ZeroOnlyAndMinMissing(t *testing.T) {
	mc := newMissingCanvas(t)
	defer mc.Close()

	out, _, err := runMissing(t, mc, "json", "--course-id", "1", "--zero-only")
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	c1 := decodeMissing(t, out).Courses[0]
	if len(c1.Students) != 2 || findStudent(c1, 10) != nil || findStudent(c1, 11) == nil || findStudent(c1, 14) == nil {
		t.Errorf("--zero-only students = %+v, want only Bob and Dee", c1.Students)
	}
	// Summary counts still describe the whole course.
	if c1.StudentsMissingAny != 3 || c1.StudentsZeroSubmissions != 2 {
		t.Errorf("summary = missing_any %d zero %d", c1.StudentsMissingAny, c1.StudentsZeroSubmissions)
	}

	table, _, err := runMissing(t, mc, "table", "--course-id", "1", "--zero-only")
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if strings.Contains(table, "Missing one or more") || strings.Contains(table, "Lovelace") {
		t.Errorf("--zero-only table must show only the zero section:\n%s", table)
	}

	out, _, err = runMissing(t, mc, "json", "--course-id", "1", "--min-missing", "2")
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	c1 = decodeMissing(t, out).Courses[0]
	if len(c1.Students) != 2 || findStudent(c1, 10) != nil {
		t.Errorf("--min-missing 2 students = %+v, want Bob and Dee only", c1.Students)
	}
}

func TestSubmissionsMissing_IncludeInactive(t *testing.T) {
	mc := newMissingCanvas(t)
	defer mc.Close()

	out, _, err := runMissing(t, mc, "json", "--course-id", "1", "--include-inactive")
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	c1 := decodeMissing(t, out).Courses[0]
	if c1.StudentsConsidered != 4 || findStudent(c1, 12) == nil {
		t.Errorf("expected Cy with --include-inactive: considered=%d students=%+v", c1.StudentsConsidered, c1.Students)
	}
	if findStudent(c1, 13) != nil {
		t.Error("Test Student must stay excluded even with --include-inactive")
	}
	roster := mc.find("/users")[0]
	for _, want := range []string{"enrollment_state%5B%5D=active", "enrollment_state%5B%5D=inactive", "enrollment_state%5B%5D=completed"} {
		if !strings.Contains(roster, want) {
			t.Errorf("roster request %q lacks %s", roster, want)
		}
	}
}

func TestSubmissionsMissing_InclusionRules(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want []int64
	}{
		{"exclude zero points", []string{"--exclude-zero-points"}, []int64{100, 106, 101}},
		{"include undated", []string{"--include-undated"}, []int64{105, 100, 106, 101, 103}},
		{"include unpublished", []string{"--published-only=false"}, []int64{105, 102, 100, 106, 101}},
		{"due-before date cuts at local midnight", []string{"--due-before", "2026-03-05"}, []int64{105, 100, 106}},
		{"due-before future includes not-yet-due", []string{"--due-before", "2026-05-01T00:00:00Z"}, []int64{105, 100, 106, 101, 104}},
		{"due-after", []string{"--due-after", "2026-03-05"}, []int64{101}},
		{"types quiz", []string{"--types", "quiz"}, []int64{101}},
		{"types discussion", []string{"--types", "discussion"}, []int64{106}},
		{"types assignment", []string{"--types", "assignment"}, []int64{105, 100}},
		{"match substring", []string{"--assignment-match", "quiz"}, []int64{101}},
		{"match regex", []string{"--assignment-match", "/^(Essay|Discussion) \\d$/"}, []int64{100, 106}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mc := newMissingCanvas(t)
			defer mc.Close()
			out, _, err := runMissing(t, mc, "json", append([]string{"--course-id", "1"}, tc.args...)...)
			if err != nil {
				t.Fatalf("run: %v", err)
			}
			got := assignmentIDs(decodeMissing(t, out).Courses[0])
			if !equalIDs(got, tc.want) {
				t.Errorf("in-scope = %v, want %v", got, tc.want)
			}
			if n := mc.count(1, "/assignments"); n != 0 {
				t.Errorf("inclusion rules must use the embedded assignment, not an assignments read (%d)", n)
			}
		})
	}
}

func TestSubmissionsMissing_ExplicitAssignmentIDs(t *testing.T) {
	mc := newMissingCanvas(t)
	defer mc.Close()

	// 101 has rows; 200-series ids are course 2's; 5001 has no rows in the
	// grid, so its metadata comes from a single assignment GET.
	out, errOut, err := runMissing(t, mc, "json", "--course-id", "1", "--assignment-id", "101", "--assignment-id", "5001", "--assignment-id", "104")
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	c1 := decodeMissing(t, out).Courses[0]
	if got := assignmentIDs(c1); !equalIDs(got, []int64{5001, 101}) {
		t.Errorf("in-scope = %v, want [5001 101] (104 is not due yet)", got)
	}
	if !strings.Contains(errOut, "assignment 104") || !strings.Contains(errOut, "not due until") {
		t.Errorf("explicitly named but excluded assignment must be reported on stderr, got %q", errOut)
	}
	// Everybody is missing 5001 (no rows at all); Ada submitted 101 late.
	ada := findStudent(c1, 10)
	if ada == nil || !equalIDs(ada.Missing, []int64{5001}) || ada.Late != 1 {
		t.Errorf("Ada = %+v", ada)
	}
	grid := mc.find("/students/submissions")
	if len(grid) == 0 || !strings.Contains(grid[0], "assignment_ids%5B%5D=101") || !strings.Contains(grid[0], "assignment_ids%5B%5D=5001") {
		t.Errorf("grid should be filtered to the named assignments: %v", grid)
	}
	if got := mc.find("/assignments/5001"); len(got) != 1 {
		t.Errorf("expected exactly one metadata GET for 5001, got %v", got)
	}
	if got := mc.find("/assignments/101"); len(got) != 0 {
		t.Errorf("101 has grid rows; no metadata GET expected, got %v", got)
	}

	// An unknown assignment is an error, not a silent empty scope.
	if _, _, err := runMissing(t, mc, "json", "--course-id", "1", "--assignment-id", "4242"); err == nil || !strings.Contains(err.Error(), "assignment 4242 not found") {
		t.Errorf("expected not-found error, got %v", err)
	}
}

func TestSubmissionsMissing_ChunksManyAssignmentIDs(t *testing.T) {
	mc := newMissingCanvas(t)
	defer mc.Close()

	args := []string{"--course-id", "2"}
	for i := int64(0); i < 150; i++ {
		args = append(args, "--assignment-id", strconv.FormatInt(5000+i, 10))
	}
	out, _, err := runMissing(t, mc, "json", args...)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	c2 := decodeMissing(t, out).Courses[0]
	if len(c2.Assignments) != 150 {
		t.Errorf("in-scope = %d, want 150", len(c2.Assignments))
	}
	grid := mc.find("/students/submissions")
	if len(grid) < 2 {
		t.Fatalf("150 ids should be chunked into several grid requests, got %d", len(grid))
	}
	seen := map[string]bool{}
	for _, u := range grid {
		if len(u) > 2000 {
			t.Errorf("grid URL is %d chars, over 2000", len(u))
		}
		parsed, _ := url.Parse(u)
		for _, id := range parsed.Query()["assignment_ids[]"] {
			seen[id] = true
		}
	}
	if len(seen) != 150 {
		t.Errorf("chunks covered %d ids, want 150", len(seen))
	}
	for _, s := range c2.Students {
		if !s.ZeroSubmissions || s.MissingCount != 150 {
			t.Errorf("student %d = %+v, want all 150 missing", s.UserID, s)
		}
	}
}

func TestSubmissionsMissing_AllActive(t *testing.T) {
	mc := newMissingCanvas(t)
	defer mc.Close()

	out, _, err := runMissing(t, mc, "json", "--all-active")
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	rep := decodeMissing(t, out)
	// Course 2 appears in both the teacher and ta lists (deduplicated);
	// course 3's term ended before the fixture clock.
	if len(rep.Courses) != 2 || rep.Courses[0].Course.ID != 2 || rep.Courses[1].Course.ID != 1 {
		var ids []int64
		for _, c := range rep.Courses {
			ids = append(ids, c.Course.ID)
		}
		t.Fatalf("courses = %v, want [2 1] (BIO200 sorts before TEST101)", ids)
	}
	lists := mc.find("/api/v1/courses?")
	if len(lists) != 2 {
		t.Fatalf("course list requests = %v, want teacher + ta", lists)
	}
	for _, u := range lists {
		for _, want := range []string{"enrollment_state=active", "include%5B%5D=term"} {
			if !strings.Contains(u, want) {
				t.Errorf("course list %q lacks %s", u, want)
			}
		}
	}
	if !strings.Contains(lists[0], "enrollment_type=teacher") || !strings.Contains(lists[1], "enrollment_type=ta") {
		t.Errorf("course lists = %v", lists)
	}
	// No per-course GET: the list already carried the term. Two reads per course.
	for _, id := range []int64{1, 2} {
		if got := mc.find(fmt.Sprintf("/api/v1/courses/%d?", id)); len(got) != 0 {
			t.Errorf("unexpected course GET %v", got)
		}
		if n := mc.count(id, "/users"); n == 0 {
			t.Errorf("course %d: no roster read", id)
		}
		if n := mc.count(id, "/students/submissions"); n == 0 {
			t.Errorf("course %d: no grid read", id)
		}
	}
	if n := mc.count(1, "/users") + mc.count(1, "/students/submissions"); n != 4 {
		t.Errorf("course 1 data requests = %d, want 2 reads × 2 pages", n)
	}
	if n := mc.count(2, "/users") + mc.count(2, "/students/submissions"); n != 2 {
		t.Errorf("course 2 data requests = %d, want exactly 2", n)
	}
	if got := mc.find("/api/v1/courses/3"); len(got) != 0 {
		t.Errorf("concluded-term course must not be read: %v", got)
	}
}

func TestSubmissionsMissing_CSVAndYAML(t *testing.T) {
	mc := newMissingCanvas(t)
	defer mc.Close()

	out, _, err := runMissing(t, mc, "csv", "--course-id", "1")
	if err != nil {
		t.Fatalf("csv: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if lines[0] != "course_id,course_code,user_id,name,sortable_name,missing_count,late,zero_submissions,missing_ids,missing_names" || len(lines) != 4 {
		t.Errorf("csv = %q", out)
	}
	if !strings.Contains(out, `1,TEST101,10,Ada Lovelace,"Lovelace, Ada",1,1,false,105,Attendance check`) {
		t.Errorf("csv lacks Ada's row:\n%s", out)
	}

	out, _, err = runMissing(t, mc, "yaml", "--course-id", "1")
	if err != nil {
		t.Fatalf("yaml: %v", err)
	}
	if !strings.Contains(out, "zero_submissions: true") || !strings.Contains(out, "course_code: TEST101") {
		t.Errorf("yaml = %s", out)
	}

	if _, _, err := runMissing(t, mc, "xml", "--course-id", "1"); err == nil || !strings.Contains(err.Error(), "unsupported output format") {
		t.Errorf("expected unsupported format error, got %v", err)
	}
}

func TestSubmissionsMissing_VerboseProgressOnStderrOnly(t *testing.T) {
	mc := newMissingCanvas(t)
	defer mc.Close()

	oldVerbose := verbose
	verbose = true
	t.Cleanup(func() { verbose = oldVerbose })
	out, errOut, err := runMissing(t, mc, "json", "--course-id", "2")
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	decodeMissing(t, out) // stdout must still be clean JSON
	if !strings.Contains(errOut, "course 2: fetching students") || !strings.Contains(errOut, "fetching submissions grid") {
		t.Errorf("verbose progress missing from stderr: %q", errOut)
	}
}

func TestSubmissionsMissing_CacheOnAndOffAgree(t *testing.T) {
	mc := newMissingCanvas(t)
	defer mc.Close()

	withoutCache, _, err := runMissing(t, mc, "json", "--course-id", "2")
	if err != nil {
		t.Fatalf("no-cache run: %v", err)
	}
	// runMissing forces noCache=true; flip it for the second run.
	t.Setenv("CANVAS_URL", mc.URL)
	t.Setenv("CANVAS_TOKEN", "test-token")
	oldFormat, oldNoCache, oldNow := outputFormat, noCache, submissionsMissingNow
	outputFormat, noCache = "json", false
	submissionsMissingNow = func() time.Time { return missingFixtureNow }
	t.Cleanup(func() { outputFormat, noCache, submissionsMissingNow = oldFormat, oldNoCache, oldNow })
	cmd := newSubmissionsMissingCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--course-id", "2"})
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("cached run: %v", err)
	}
	if out.String() != withoutCache {
		t.Errorf("cached and uncached reports differ:\n%s\n---\n%s", out.String(), withoutCache)
	}
}

func TestSubmissionsMissing_Validation(t *testing.T) {
	mc := newMissingCanvas(t)
	defer mc.Close()

	for _, tc := range []struct {
		args []string
		want string
	}{
		{nil, "one of --course-id or --all-active"},
		{[]string{"--course-id", "1", "--all-active"}, "mutually exclusive"},
		{[]string{"--course-id", "1", "--assignment-id", "5", "--assignment-match", "x"}, "mutually exclusive"},
		{[]string{"--course-id", "1", "--types", "page"}, "unknown --types value"},
		{[]string{"--course-id", "1", "--due-before", "soon"}, "invalid --due-before"},
		{[]string{"--course-id", "404"}, "course with ID 404 not found"},
	} {
		_, _, err := runMissing(t, mc, "table", tc.args...)
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("args %v: err = %v, want containing %q", tc.args, err, tc.want)
		}
	}
	if len(mc.find("/users")) != 0 {
		t.Error("validation failures must not reach the roster read")
	}

	oldLimit := globalLimit
	globalLimit = 50
	t.Cleanup(func() { globalLimit = oldLimit })
	if _, _, err := runMissing(t, mc, "table", "--course-id", "1"); err == nil || !strings.Contains(err.Error(), "--limit is not supported") {
		t.Errorf("--limit must be refused, got %v", err)
	}
}

func TestSubmissionsMissing_UnitHelpers(t *testing.T) {
	if !submissionIsMissing(nil) {
		t.Error("no row must be missing")
	}
	if submissionIsMissing(&api.Submission{WorkflowState: "unsubmitted", ExcusedTLN: true}) {
		t.Error("excused must never be missing")
	}
	if !submissionIsMissing(&api.Submission{WorkflowState: "graded", Missing: true}) {
		t.Error("Canvas missing flag must count")
	}
	if submissionIsMissing(&api.Submission{WorkflowState: "submitted", Late: true}) {
		t.Error("late-but-submitted is not missing")
	}
	if !isTestStudent(&api.User{Name: "Student", Enrollments: []api.Enrollment{{Type: "StudentViewEnrollment"}}}) {
		t.Error("StudentViewEnrollment must be recognised")
	}
	if got := truncateMissingItems([]string{"Alpha", "Beta", "Gamma", "Delta"}, 14); got != "Alpha, … (+3 more)" {
		t.Errorf("truncate = %q", got)
	}
	if missingCourseHeader(&missingCourseInfo{Name: "Only name"}) != "Only name" {
		t.Error("header without code/term")
	}
	if got := missingAssignmentKind(&api.Assignment{SubmissionTypes: []string{"on_paper"}}); got != "" {
		t.Errorf("on_paper kind = %q, want not submittable", got)
	}
	if !courseInCurrentTerm(&api.Course{}, missingFixtureNow) || courseInCurrentTerm(&api.Course{WorkflowState: "completed"}, missingFixtureNow) {
		t.Error("courseInCurrentTerm")
	}
	var b bytes.Buffer
	if err := renderMissingTable(&b, &missingReport{}, &options.SubmissionsMissingOptions{}); err != nil || !strings.Contains(b.String(), "No courses in scope") {
		t.Errorf("empty table = %q, %v", b.String(), err)
	}
	b.Reset()
	if err := renderMissingMarkdown(&b, &missingReport{}, &options.SubmissionsMissingOptions{}); err != nil || !strings.Contains(b.String(), "No courses") {
		t.Errorf("empty markdown = %q, %v", b.String(), err)
	}
}
