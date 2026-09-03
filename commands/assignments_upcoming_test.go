package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// upcomingCanvas is a fake Canvas for `assignments upcoming`: two courses
// in the current term and one in an ended term, and assignments on either
// side of the window, undated, unpublished, quiz- and discussion-backed.
//
// Fixture clock: Thursday 2026-09-03 10:00 America/New_York (14:00Z).
type upcomingCanvas struct {
	*httptest.Server
	mu       sync.Mutex
	requests []string
}

var upcomingFixtureNow = time.Date(2026, 9, 3, 14, 0, 0, 0, time.UTC)

const (
	upcomingCourse1 = `{"id":1,"name":"Intro to GC","course_code":"GC1010","workflow_state":"available","term":{"id":9,"name":"Fall 2026","start_at":"2026-08-15T00:00:00Z","end_at":"2026-12-20T00:00:00Z"}}`
	upcomingCourse2 = `{"id":2,"name":"Senior Studio","course_code":"GC4800","workflow_state":"available","term":{"id":9,"name":"Fall 2026","start_at":"2026-08-15T00:00:00Z","end_at":"2026-12-20T00:00:00Z"}}`
	upcomingCourse3 = `{"id":3,"name":"Old","course_code":"OLD100","workflow_state":"available","term":{"id":8,"name":"Spring 2026","start_at":"2026-01-10T00:00:00Z","end_at":"2026-05-10T00:00:00Z"}}`
)

var upcomingAssignments = map[int64]string{
	1: `[
{"id":101,"name":"Reading quiz","published":true,"due_at":"2026-09-07T03:59:00Z","points_possible":10,"submission_types":["online_quiz"],"is_quiz_assignment":true,"quiz_id":77,"html_url":"https://x/1"},
{"id":102,"name":"Essay 1","published":true,"due_at":"2026-09-12T03:59:00Z","points_possible":50,"submission_types":["online_upload","online_text_entry"]},
{"id":103,"name":"Past lab","published":true,"due_at":"2026-09-01T03:59:00Z","points_possible":5,"submission_types":["online_upload"]},
{"id":104,"name":"Undated reading","published":true,"due_at":null,"points_possible":0,"submission_types":["online_text_entry"]},
{"id":105,"name":"Draft essay","published":false,"due_at":"2026-09-05T03:59:00Z","points_possible":10,"submission_types":["online_upload"]},
{"id":106,"name":"Far project","published":true,"due_at":"2026-10-01T03:59:00Z","points_possible":100,"submission_types":["online_upload"]},
{"id":107,"name":"Week 2 discussion","published":true,"due_at":"2026-09-05T16:00:00Z","points_possible":2.5,"submission_types":["discussion_topic"]},
{"id":108,"name":"Attendance today","published":true,"due_at":"2026-09-03T13:00:00Z","points_possible":1,"submission_types":["online_url"]}
]`,
	2: `[
{"id":201,"name":"Studio critique","published":true,"due_at":"2026-09-10T03:59:00Z","points_possible":20,"submission_types":["on_paper"]},
{"id":202,"name":"Portfolio","published":true,"due_at":null,"points_possible":100,"submission_types":["online_upload"]}
]`,
	3: `[{"id":301,"name":"Should not be read","published":true,"due_at":"2026-09-05T03:59:00Z","points_possible":1,"submission_types":["online_upload"]}]`,
}

func newUpcomingCanvas(t *testing.T) *upcomingCanvas {
	t.Helper()
	uc := &upcomingCanvas{}
	uc.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			_, _ = w.Write([]byte(`[]`))
			return
		}
		uc.mu.Lock()
		uc.requests = append(uc.requests, r.URL.RequestURI())
		uc.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		q := r.URL.Query()
		switch r.URL.Path {
		case "/api/v1/courses":
			switch q.Get("enrollment_type") {
			case "teacher":
				fmt.Fprintf(w, "[%s,%s]", upcomingCourse1, upcomingCourse2)
			case "ta":
				fmt.Fprintf(w, "[%s,%s]", upcomingCourse2, upcomingCourse3)
			default:
				fmt.Fprint(w, "[]")
			}
		case "/api/v1/courses/1":
			fmt.Fprint(w, upcomingCourse1)
		case "/api/v1/courses/2":
			fmt.Fprint(w, upcomingCourse2)
		case "/api/v1/courses/1/assignments":
			fmt.Fprint(w, upcomingAssignments[1])
		case "/api/v1/courses/2/assignments":
			fmt.Fprint(w, upcomingAssignments[2])
		case "/api/v1/courses/3/assignments":
			t.Errorf("course 3 (ended term) was read")
			fmt.Fprint(w, upcomingAssignments[3])
		case "/api/v1/courses/404":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"errors":[{"message":"The specified resource does not exist."}]}`)
		default:
			t.Errorf("unexpected request %s", r.URL.RequestURI())
			http.NotFound(w, r)
		}
	}))
	return uc
}

func (uc *upcomingCanvas) count(substr string) int {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	n := 0
	for _, u := range uc.requests {
		if strings.Contains(u, substr) {
			n++
		}
	}
	return n
}

func runUpcoming(t *testing.T, uc *upcomingCanvas, format string, args ...string) (string, string, error) {
	t.Helper()
	t.Setenv("CANVAS_URL", uc.URL)
	t.Setenv("CANVAS_TOKEN", "test-token")
	t.Setenv("CANVAS_REQUESTS_PER_SEC", "1000")

	oldFormat, oldNoCache, oldNow, oldLimit := outputFormat, noCache, assignmentsUpcomingNow, globalLimit
	outputFormat, noCache, globalLimit = format, true, 0
	assignmentsUpcomingNow = func() time.Time { return upcomingFixtureNow }
	t.Cleanup(func() {
		outputFormat, noCache, assignmentsUpcomingNow, globalLimit = oldFormat, oldNoCache, oldNow, oldLimit
	})

	cmd := newAssignmentsUpcomingCmd()
	cmd.SilenceUsage, cmd.SilenceErrors = true, true
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs(append([]string{"--timezone", "America/New_York"}, args...))
	err := cmd.ExecuteContext(context.Background())
	return out.String(), errOut.String(), err
}

func decodeUpcoming(t *testing.T, raw string) upcomingReport {
	t.Helper()
	var rep upcomingReport
	if err := json.Unmarshal([]byte(raw), &rep); err != nil {
		t.Fatalf("not json: %v\n%s", err, raw)
	}
	return rep
}

func upcomingIDs(cr upcomingCourseReport) []int64 {
	var ids []int64
	for _, a := range cr.Assignments {
		ids = append(ids, a.ID)
	}
	return ids
}

func TestAssignmentsUpcoming_WithinWindow(t *testing.T) {
	uc := newUpcomingCanvas(t)
	defer uc.Close()

	raw, _, err := runUpcoming(t, uc, "json", "--course-id", "1", "--within", "10d")
	if err != nil {
		t.Fatal(err)
	}
	rep := decodeUpcoming(t, raw)
	if rep.Timezone != "America/New_York" || !rep.Now.Equal(upcomingFixtureNow) || !rep.Limit.Equal(upcomingFixtureNow.Add(10*24*time.Hour)) {
		t.Errorf("header = %+v", rep)
	}
	if len(rep.Courses) != 1 {
		t.Fatalf("courses = %d", len(rep.Courses))
	}
	cr := rep.Courses[0]
	// Past (103), undated (104), unpublished (105), beyond the window (106)
	// and already due today (108, 9 AM) are out; sorted by due date.
	if got := fmt.Sprint(upcomingIDs(cr)); got != "[107 101 102]" {
		t.Errorf("ids = %s", got)
	}
	if cr.Count != 3 || cr.UndatedSeen != 1 || len(cr.Undated) != 0 {
		t.Errorf("counts = %+v", cr)
	}
	if cr.Course.CourseCode != "GC1010" || cr.Course.Term != "Fall 2026" {
		t.Errorf("course = %+v", cr.Course)
	}
	quiz := cr.Assignments[1]
	if quiz.Type != "quiz" || quiz.QuizID != 77 || quiz.DueLocal != "Sun 2026-09-06 11:59 PM EDT" || quiz.DueAt.Format(time.RFC3339) != "2026-09-07T03:59:00Z" || quiz.HTMLURL != "https://x/1" {
		t.Errorf("quiz item = %+v", quiz)
	}
	if cr.Assignments[0].Type != "discussion" || cr.Assignments[2].Type != "assignment" {
		t.Errorf("kinds = %s / %s", cr.Assignments[0].Type, cr.Assignments[2].Type)
	}
	// One course GET, one assignments read, nothing else.
	if uc.count("/api/v1/courses/1?") != 1 || uc.count("/courses/1/assignments") != 1 || uc.count("enrollment_type") != 0 {
		t.Errorf("requests = %v", uc.requests)
	}
}

func TestAssignmentsUpcoming_ByThisSunday(t *testing.T) {
	uc := newUpcomingCanvas(t)
	defer uc.Close()

	raw, _, err := runUpcoming(t, uc, "json", "--course-id", "1", "--by", "this sunday")
	if err != nil {
		t.Fatal(err)
	}
	rep := decodeUpcoming(t, raw)
	// Sunday 2026-09-06 23:59:59 EDT — the quiz due at 11:59 PM is in, Essay 1 (Friday next) is out.
	if rep.Limit.Format(time.RFC3339) != "2026-09-07T03:59:59Z" {
		t.Errorf("limit = %s", rep.Limit.Format(time.RFC3339))
	}
	if got := fmt.Sprint(upcomingIDs(rep.Courses[0])); got != "[107 101]" {
		t.Errorf("ids = %s", got)
	}

	// A date/time limit is used as given.
	raw, _, err = runUpcoming(t, uc, "json", "--course-id", "1", "--by", "2026-09-05 11am")
	if err != nil {
		t.Fatal(err)
	}
	rep = decodeUpcoming(t, raw)
	if rep.Limit.Format(time.RFC3339) != "2026-09-05T15:00:00Z" || fmt.Sprint(upcomingIDs(rep.Courses[0])) != "[]" {
		t.Errorf("limit = %s ids = %v", rep.Limit.Format(time.RFC3339), upcomingIDs(rep.Courses[0]))
	}

	// A limit in the past is refused.
	_, _, err = runUpcoming(t, uc, "json", "--course-id", "1", "--by", "yesterday")
	if err == nil || !strings.Contains(err.Error(), "is not after now") {
		t.Errorf("past limit: %v", err)
	}
	_, _, err = runUpcoming(t, uc, "json", "--course-id", "1", "--by", "4:50")
	if err == nil || !strings.Contains(err.Error(), "invalid --by") {
		t.Errorf("ambiguous limit: %v", err)
	}
}

func TestAssignmentsUpcoming_UndatedAndUnpublished(t *testing.T) {
	uc := newUpcomingCanvas(t)
	defer uc.Close()

	raw, _, err := runUpcoming(t, uc, "json", "--course-id", "1", "--within", "10d", "--include-undated", "--published-only=false")
	if err != nil {
		t.Fatal(err)
	}
	cr := decodeUpcoming(t, raw).Courses[0]
	if got := fmt.Sprint(upcomingIDs(cr)); got != "[105 107 101 102]" {
		t.Errorf("ids = %s", got)
	}
	if len(cr.Undated) != 1 || cr.Undated[0].ID != 104 || cr.UndatedSeen != 1 {
		t.Errorf("undated = %+v", cr.Undated)
	}

	out, _, err := runUpcoming(t, uc, "table", "--course-id", "1", "--within", "10d", "--include-undated", "--published-only=false")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Undated (1)", "Undated reading", "Draft essay", "no"} {
		if !strings.Contains(out, want) {
			t.Errorf("table lacks %q:\n%s", want, out)
		}
	}
}

func TestAssignmentsUpcoming_TableAndMarkdown(t *testing.T) {
	uc := newUpcomingCanvas(t)
	defer uc.Close()

	out, _, err := runUpcoming(t, uc, "table", "--course-id", "1", "--course-id", "2", "--within", "10d")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Due after Thu 2026-09-03 10:00 AM EDT, up to Sun 2026-09-13 10:00 AM EDT (times in America/New_York)",
		"GC1010 — Intro to GC (Fall 2026)",
		"3 due by Sun 2026-09-13 10:00 AM EDT · 1 undated not shown (--include-undated)",
		"DUE", "TITLE", "TYPE", "POINTS", "SUBMISSION", "PUBLISHED",
		"Sat 2026-09-05 12:00 PM EDT  Week 2 discussion  discussion  107  2.5     discussion_topic",
		"Sun 2026-09-06 11:59 PM EDT  Reading quiz       quiz        101  10      online_quiz",
		"Fri 2026-09-11 11:59 PM EDT  Essay 1            assignment  102  50      online_upload,online_text_entry  yes",
		"GC4800 — Senior Studio (Fall 2026)",
		"1 due by Sun 2026-09-13 10:00 AM EDT · 1 undated not shown",
		"Studio critique",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("table lacks %q:\n%s", want, out)
		}
	}
	for _, banned := range []string{"Past lab", "Far project", "Draft essay", "Undated reading", "Attendance today", "Portfolio"} {
		if strings.Contains(out, banned) {
			t.Errorf("table shows %q:\n%s", banned, out)
		}
	}

	md, _, err := runUpcoming(t, uc, "markdown", "--course-id", "1", "--within", "10d")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Due by **Sun 2026-09-13 10:00 AM EDT** (America/New_York)",
		"### GC1010 — Intro to GC (Fall 2026)",
		"_3 due by Sun 2026-09-13 10:00 AM EDT · 1 undated not shown (--include-undated)_",
		"- **Sun 2026-09-06 11:59 PM EDT** — Reading quiz (quiz, 10 pts, online_quiz)",
		"- **Fri 2026-09-11 11:59 PM EDT** — Essay 1 (assignment, 50 pts, online_upload,online_text_entry)",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("markdown lacks %q:\n%s", want, md)
		}
	}
}

func TestAssignmentsUpcoming_CSVAndYAML(t *testing.T) {
	uc := newUpcomingCanvas(t)
	defer uc.Close()

	out, _, err := runUpcoming(t, uc, "csv", "--course-id", "1", "--within", "10d", "--include-undated")
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if lines[0] != "course_code,course_id,assignment_id,name,type,due_local,due_at,points_possible,submission_types,published" || len(lines) != 5 {
		t.Errorf("csv:\n%s", out)
	}
	if lines[2] != "GC1010,1,101,Reading quiz,quiz,Sun 2026-09-06 11:59 PM EDT,2026-09-07T03:59:00Z,10,online_quiz,yes" {
		t.Errorf("csv row = %s", lines[2])
	}
	if lines[4] != "GC1010,1,104,Undated reading,assignment,,,0,online_text_entry,yes" {
		t.Errorf("csv undated row = %s", lines[4])
	}

	y, _, err := runUpcoming(t, uc, "yaml", "--course-id", "1", "--within", "10d")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(y, "course_code: GC1010") || !strings.Contains(y, "due_local: Sun 2026-09-06 11:59 PM EDT") {
		t.Errorf("yaml:\n%s", y)
	}

	_, _, err = runUpcoming(t, uc, "xml", "--course-id", "1", "--within", "10d")
	if err == nil || !strings.Contains(err.Error(), "unsupported output format") {
		t.Errorf("xml: %v", err)
	}
}

func TestAssignmentsUpcoming_AllActive(t *testing.T) {
	uc := newUpcomingCanvas(t)
	defer uc.Close()

	raw, _, err := runUpcoming(t, uc, "json", "--all-active", "--within", "2w")
	if err != nil {
		t.Fatal(err)
	}
	rep := decodeUpcoming(t, raw)
	if len(rep.Courses) != 2 || rep.Courses[0].Course.ID != 1 || rep.Courses[1].Course.ID != 2 {
		t.Errorf("courses = %+v", rep.Courses)
	}
	if uc.count("enrollment_type=teacher") != 1 || uc.count("enrollment_type=ta") != 1 || uc.count("/api/v1/courses/1?") != 0 || uc.count("/courses/2/assignments") != 1 {
		t.Errorf("requests = %v", uc.requests)
	}
}

func TestAssignmentsUpcoming_Errors(t *testing.T) {
	uc := newUpcomingCanvas(t)
	defer uc.Close()

	_, _, err := runUpcoming(t, uc, "table", "--course-id", "404", "--within", "1d")
	if err == nil || !strings.Contains(err.Error(), "course with ID 404 not found") {
		t.Errorf("404: %v", err)
	}
	_, _, err = runUpcoming(t, uc, "table", "--within", "1d")
	if err == nil || !strings.Contains(err.Error(), "one of --course-id or --all-active") {
		t.Errorf("no scope: %v", err)
	}
	_, _, err = runUpcoming(t, uc, "table", "--course-id", "1")
	if err == nil || !strings.Contains(err.Error(), "one of --within or --by") {
		t.Errorf("no window: %v", err)
	}
	_, _, err = runUpcoming(t, uc, "table", "--course-id", "1", "--within", "1d", "--by", "tomorrow")
	if err == nil || !strings.Contains(err.Error(), "mutually exclusive") {
		t.Errorf("both windows: %v", err)
	}

	oldLimit := globalLimit
	globalLimit = 5
	t.Cleanup(func() { globalLimit = oldLimit })
	cmd := newAssignmentsUpcomingCmd()
	cmd.SilenceUsage, cmd.SilenceErrors = true, true
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetArgs([]string{"--course-id", "1", "--within", "1d"})
	if err := cmd.ExecuteContext(context.Background()); err == nil || !strings.Contains(err.Error(), "--limit is not supported") {
		t.Errorf("--limit: %v", err)
	}
}

func TestAssignmentsUpcoming_UnitHelpers(t *testing.T) {
	if upcomingPoints(10) != "10" || upcomingPoints(2.5) != "2.5" || upcomingPoints(0) != "0" || upcomingPoints(1.25) != "1.25" {
		t.Error("upcomingPoints")
	}
	if upcomingSubmission(nil) != "-" || upcomingSubmission([]string{"a", "b"}) != "a,b" {
		t.Error("upcomingSubmission")
	}
	if upcomingCourseHeader(&upcomingCourseInfo{ID: 7}) != "course 7" || upcomingCourseHeader(&upcomingCourseInfo{ID: 7, CourseCode: "X", Name: "X"}) != "X" {
		t.Error("upcomingCourseHeader")
	}
	if upcomingTermName(&api.Course{}) != "" || upcomingTermName(&api.Course{Term: &api.Term{Name: "Fall"}}) != "Fall" {
		t.Error("upcomingTermName")
	}
	now := upcomingFixtureNow
	if upcomingCourseInTerm(&api.Course{WorkflowState: "completed"}, now) {
		t.Error("completed course kept")
	}
	if !upcomingCourseInTerm(&api.Course{}, now) {
		t.Error("course without term dropped")
	}
	if upcomingCourseInTerm(&api.Course{Term: &api.Term{EndAt: now.Add(-time.Hour)}}, now) {
		t.Error("ended term kept")
	}
	if upcomingCourseInTerm(&api.Course{Term: &api.Term{StartAt: now.Add(time.Hour)}}, now) {
		t.Error("future term kept")
	}
	if upcomingKind(&api.Assignment{SubmissionTypes: []string{"online_upload"}}) != "assignment" ||
		upcomingKind(&api.Assignment{QuizID: 3}) != "quiz" ||
		upcomingKind(&api.Assignment{SubmissionTypes: []string{"discussion_topic"}}) != "discussion" {
		t.Error("upcomingKind")
	}
}
