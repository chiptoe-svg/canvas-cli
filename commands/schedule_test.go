package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// scheduleCanvas is a fake Canvas for `schedule`: one course with classic
// quizzes (some graded, one practice), their assignments, plain
// assignments, an id that exists as both a quiz and an assignment, and
// mutable dates so a PUT is visible on the following GET.
//
// Fixture clock: Thursday 2026-09-03 10:00 America/New_York.
type scheduleCanvas struct {
	*httptest.Server
	mu       sync.Mutex
	requests []string // "METHOD URI BODY"
	quizzes  map[int64]*fakeSchedItem
	assigns  map[int64]*fakeSchedItem
	// staleQuiz makes the read-back of that quiz report the pre-write due date.
	staleQuiz int64
	// failAssignment makes PUTs to that assignment answer 400 (not retried).
	failAssignment int64
}

type fakeSchedItem struct {
	id                   int64
	title                string
	quizID, assignmentID int64
	quizAssignment       bool
	overrides            bool
	unlock, due, lock    *string
}

var scheduleFixtureNow = time.Date(2026, 9, 3, 14, 0, 0, 0, time.UTC)

func str(s string) *string { return &s }

func newScheduleCanvas(t *testing.T) *scheduleCanvas {
	t.Helper()
	sc := &scheduleCanvas{
		quizzes: map[int64]*fakeSchedItem{
			456: {id: 456, title: "Attendance 9/9", assignmentID: 789, unlock: str("2026-09-09T19:00:00Z")},
			457: {id: 457, title: "Attendance 9/16", assignmentID: 790},
			458: {id: 458, title: "Practice quiz"},
			460: {id: 460, title: "Reading quiz", assignmentID: 792, due: str("2026-09-05T03:59:00Z"), lock: str("2026-09-05T03:59:00Z")},
			5:   {id: 5, title: "Five quiz", assignmentID: 795},
		},
		assigns: map[int64]*fakeSchedItem{
			789: {id: 789, title: "Attendance 9/9", quizID: 456, quizAssignment: true, unlock: str("2026-09-09T19:00:00Z")},
			790: {id: 790, title: "Attendance 9/16", quizID: 457, quizAssignment: true},
			792: {id: 792, title: "Reading quiz", quizID: 460, quizAssignment: true, due: str("2026-09-05T03:59:00Z"), lock: str("2026-09-05T03:59:00Z")},
			795: {id: 795, title: "Five quiz", quizID: 5, quizAssignment: true},
			800: {id: 800, title: "Course lineup"},
			801: {id: 801, title: "Attendance essay", overrides: true, due: str("2026-09-20T03:59:00Z")},
			802: {id: 802, title: "Locked essay", unlock: str("2026-09-10T04:00:00Z"), due: str("2026-09-12T03:59:00Z"), lock: str("2026-09-12T03:59:00Z")},
			5:   {id: 5, title: "Five assignment"},
		},
	}
	sc.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			_, _ = w.Write([]byte(`[]`))
			return
		}
		body, _ := io.ReadAll(r.Body)
		sc.mu.Lock()
		sc.requests = append(sc.requests, strings.TrimSpace(r.Method+" "+r.URL.RequestURI()+" "+string(body)))
		sc.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")

		path := strings.TrimPrefix(r.URL.Path, "/api/v1/courses/1/")
		switch {
		case path == "quizzes" && r.Method == "GET":
			fmt.Fprint(w, sc.list(sc.quizzes, sc.quizJSON))
		case path == "assignments" && r.Method == "GET":
			fmt.Fprint(w, sc.list(sc.assigns, sc.assignJSON))
		case strings.HasPrefix(path, "quizzes/"):
			id, _ := strconv.ParseInt(strings.TrimPrefix(path, "quizzes/"), 10, 64)
			sc.mu.Lock()
			it := sc.quizzes[id]
			sc.mu.Unlock()
			if it == nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"errors":[{"message":"The specified resource does not exist."}]}`)
				return
			}
			if r.Method == "PUT" {
				sc.apply(it, body, "quiz")
				fmt.Fprint(w, sc.quizJSON(it))
				return
			}
			if id == sc.staleQuiz {
				stale := *it
				stale.due = str("2026-01-01T00:00:00Z")
				fmt.Fprint(w, sc.quizJSON(&stale))
				return
			}
			fmt.Fprint(w, sc.quizJSON(it))
		case strings.HasPrefix(path, "assignments/"):
			id, _ := strconv.ParseInt(strings.TrimPrefix(path, "assignments/"), 10, 64)
			sc.mu.Lock()
			it := sc.assigns[id]
			sc.mu.Unlock()
			if it == nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"errors":[{"message":"The specified resource does not exist."}]}`)
				return
			}
			if r.Method == "PUT" {
				if id == sc.failAssignment {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprint(w, `{"errors":[{"message":"boom"}]}`)
					return
				}
				sc.apply(it, body, "assignment")
			}
			fmt.Fprint(w, sc.assignJSON(it))
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.RequestURI())
			http.NotFound(w, r)
		}
	}))
	return sc
}

func (sc *scheduleCanvas) list(m map[int64]*fakeSchedItem, render func(*fakeSchedItem) string) string {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	var parts []string
	for _, it := range m {
		parts = append(parts, render(it))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func jsonTime(p *string) string {
	if p == nil {
		return "null"
	}
	return strconv.Quote(*p)
}

func (sc *scheduleCanvas) quizJSON(it *fakeSchedItem) string {
	return fmt.Sprintf(`{"id":%d,"title":%q,"assignment_id":%d,"quiz_type":"assignment","unlock_at":%s,"due_at":%s,"lock_at":%s}`,
		it.id, it.title, it.assignmentID, jsonTime(it.unlock), jsonTime(it.due), jsonTime(it.lock))
}

func (sc *scheduleCanvas) assignJSON(it *fakeSchedItem) string {
	types := `["online_upload"]`
	if it.quizAssignment {
		types = `["online_quiz"]`
	}
	return fmt.Sprintf(`{"id":%d,"name":%q,"quiz_id":%d,"is_quiz_assignment":%t,"submission_types":%s,"has_overrides":%t,"published":true,"unlock_at":%s,"due_at":%s,"lock_at":%s}`,
		it.id, it.title, it.quizID, it.quizAssignment, types, it.overrides, jsonTime(it.unlock), jsonTime(it.due), jsonTime(it.lock))
}

// apply mutates the item with the PUT body's date fields (null clears).
func (sc *scheduleCanvas) apply(it *fakeSchedItem, body []byte, wrapper string) {
	var req map[string]map[string]*string
	_ = json.Unmarshal(body, &req)
	sc.mu.Lock()
	defer sc.mu.Unlock()
	fields := req[wrapper]
	if v, ok := fields["unlock_at"]; ok {
		it.unlock = v
	}
	if v, ok := fields["due_at"]; ok {
		it.due = v
	}
	if v, ok := fields["lock_at"]; ok {
		it.lock = v
	}
	// The mirrored assignment/quiz keeps the same dates, as Canvas does.
	if wrapper == "quiz" && it.assignmentID != 0 {
		if a := sc.assigns[it.assignmentID]; a != nil {
			a.unlock, a.due, a.lock = it.unlock, it.due, it.lock
		}
	}
}

func (sc *scheduleCanvas) puts() []string {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	var out []string
	for _, r := range sc.requests {
		if strings.HasPrefix(r, "PUT ") {
			out = append(out, r)
		}
	}
	return out
}

// runScheduleCmd executes `schedule` against sc.
func runScheduleCmd(t *testing.T, sc *scheduleCanvas, format string, dry bool, stdin string, args ...string) (string, string, error) {
	t.Helper()
	t.Setenv("CANVAS_URL", sc.URL)
	t.Setenv("CANVAS_TOKEN", "test-token")
	t.Setenv("CANVAS_REQUESTS_PER_SEC", "1000")

	oldFormat, oldNoCache, oldNow, oldDry, oldQuiet := outputFormat, noCache, scheduleNow, dryRun, quiet
	outputFormat, noCache, dryRun, quiet = format, true, dry, false
	scheduleNow = func() time.Time { return scheduleFixtureNow }
	t.Cleanup(func() {
		outputFormat, noCache, scheduleNow, dryRun, quiet = oldFormat, oldNoCache, oldNow, oldDry, oldQuiet
	})

	cmd := newScheduleCmd()
	// rootCmd silences usage/error echo; the bare subcommand must too so
	// structured output stays parseable.
	cmd.SilenceUsage, cmd.SilenceErrors = true, true
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetIn(strings.NewReader(stdin))
	cmd.SetArgs(append([]string{"--course-id", "1", "--timezone", "America/New_York"}, args...))
	err := cmd.ExecuteContext(context.Background())
	return out.String(), errOut.String(), err
}

func decodeSchedule(t *testing.T, raw string) scheduleResult {
	t.Helper()
	var res scheduleResult
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		t.Fatalf("not json: %v\n%s", err, raw)
	}
	return res
}

func itemKeys(res scheduleResult) []string {
	var keys []string
	for _, it := range res.Items {
		keys = append(keys, fmt.Sprintf("%s:%d", it.Kind, it.ID))
	}
	return keys
}

func TestSchedule_DryRunByMatchQuiz(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()

	out, _, err := runScheduleCmd(t, sc, "table", true, "", "--match", "attendance", "--type", "quiz",
		"--date", "2026-09-09", "--available", "4pm", "--due", "4:50pm", "--closed", "4:50pm")
	if err != nil {
		t.Fatalf("dry run failed: %v\n%s", err, out)
	}
	if len(sc.puts()) != 0 {
		t.Errorf("dry run sent writes: %v", sc.puts())
	}
	for _, want := range []string{
		"DRY RUN: no changes were made.",
		"available: Wed 2026-09-09 4:00 PM EDT (2026-09-09T20:00:00Z)",
		"due:       Wed 2026-09-09 4:50 PM EDT (2026-09-09T20:50:00Z)",
		"Attendance 9/16",
		"Wed 2026-09-09 3:00 PM EDT", // the old available of quiz 456
		`PUT /api/v1/courses/1/quizzes/456 {"quiz":{"due_at":"2026-09-09T20:50:00Z","lock_at":"2026-09-09T20:50:00Z","unlock_at":"2026-09-09T20:00:00Z"}}`,
		`PUT /api/v1/courses/1/quizzes/457 {"quiz":{"due_at":"2026-09-09T20:50:00Z","lock_at":"2026-09-09T20:50:00Z","unlock_at":"2026-09-09T20:00:00Z"}}`,
		"Plan: 2 matched, 2 would change, 0 written.",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Attendance essay") {
		t.Errorf("--type quiz included a plain assignment:\n%s", out)
	}

	raw, _, err := runScheduleCmd(t, sc, "json", true, "", "--match", "attendance", "--type", "quiz",
		"--date", "2026-09-09", "--available", "4pm", "--due", "4:50pm", "--closed", "4:50pm")
	if err != nil {
		t.Fatal(err)
	}
	res := decodeSchedule(t, raw)
	if !res.DryRun || res.Timezone != "America/New_York" || res.Summary.Matched != 2 || res.Summary.Changed != 2 || res.Summary.Written != 0 {
		t.Errorf("summary = %+v dry=%v tz=%s", res.Summary, res.DryRun, res.Timezone)
	}
	if got := strings.Join(itemKeys(res), " "); got != "quiz:457 quiz:456" {
		t.Errorf("items = %s", got)
	}
	it := res.Items[1]
	if it.AssignmentID != 789 || it.Before.UnlockAt == nil || it.Before.DueAt != nil || it.After.DueAt == nil ||
		it.After.DueAt.Format(time.RFC3339) != "2026-09-09T20:50:00Z" || it.Verified != "-" {
		t.Errorf("item 456 = %+v", it)
	}
}

func TestSchedule_TypeAllNeverTwice(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()

	raw, _, err := runScheduleCmd(t, sc, "json", true, "", "--match", "attendance", "--date", "2026-09-09", "--due", "4:50pm")
	if err != nil {
		t.Fatal(err)
	}
	res := decodeSchedule(t, raw)
	// Quizzes first (by title), then plain assignments; 789/790 are their quizzes.
	if got := strings.Join(itemKeys(res), " "); got != "assignment:801 quiz:457 quiz:456" {
		t.Errorf("items = %s", got)
	}
	for _, it := range res.Items {
		if it.Kind == "assignment" && !it.HasOverrides {
			t.Errorf("assignment 801 should carry has_overrides")
		}
	}

	raw, _, err = runScheduleCmd(t, sc, "json", true, "", "--match", "attendance", "--type", "assignment", "--date", "2026-09-09", "--due", "4:50pm")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(itemKeys(decodeSchedule(t, raw)), " "); got != "assignment:801" {
		t.Errorf("--type assignment items = %s", got)
	}

	// A regex, anchored and case-sensitive: only the two attendance quizzes.
	raw, _, err = runScheduleCmd(t, sc, "json", true, "", "--match", "/^Attendance 9/", "--date", "2026-09-09", "--due", "4:50pm")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(itemKeys(decodeSchedule(t, raw)), " "); got != "quiz:457 quiz:456" {
		t.Errorf("regex items = %s", got)
	}
}

func TestSchedule_ByIDQuizWriteAndReadBack(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()

	out, _, err := runScheduleCmd(t, sc, "table", false, "", "--id", "456", "--type", "quiz",
		"--date", "2026-09-09", "--available", "4:00pm", "--due", "4:50pm", "--closed", "4:50pm")
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	puts := sc.puts()
	if len(puts) != 1 || !strings.HasPrefix(puts[0], "PUT /api/v1/courses/1/quizzes/456 ") ||
		!strings.Contains(puts[0], `"unlock_at":"2026-09-09T20:00:00Z"`) || !strings.Contains(puts[0], `"due_at":"2026-09-09T20:50:00Z"`) || !strings.Contains(puts[0], `"lock_at":"2026-09-09T20:50:00Z"`) {
		t.Errorf("writes = %v", puts)
	}
	for _, want := range []string{
		"Result for course 1 (times in America/New_York):",
		"Wed 2026-09-09 3:00 PM EDT", // before
		"Wed 2026-09-09 4:50 PM EDT",
		"2026-09-09T20:50:00Z",
		"yes",
		"Done: 1 matched, 1 changed, 1 written, 1 verified, 0 mismatched, 0 failed.",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	// The fake mirrored the quiz dates onto its assignment.
	if got := *sc.assigns[789].due; got != "2026-09-09T20:50:00Z" {
		t.Errorf("assignment 789 due = %s", got)
	}
}

func TestSchedule_QuizAssignmentIsRoutedThroughQuiz(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()

	raw, errOut, err := runScheduleCmd(t, sc, "json", false, "", "--id", "789", "--type", "assignment", "--date", "2026-09-09", "--due", "4:50pm")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(errOut, `note: assignment 789 "Attendance 9/9" is quiz 456; updating the quiz`) {
		t.Errorf("stderr = %q", errOut)
	}
	puts := sc.puts()
	if len(puts) != 1 || !strings.HasPrefix(puts[0], "PUT /api/v1/courses/1/quizzes/456 ") {
		t.Errorf("writes = %v", puts)
	}
	res := decodeSchedule(t, raw)
	if len(res.Items) != 1 || res.Items[0].Kind != "quiz" || res.Items[0].ID != 456 || res.Items[0].AssignmentID != 789 || res.Items[0].Verified != "yes" {
		t.Errorf("items = %+v", res.Items)
	}
}

func TestSchedule_ByIDAmbiguousAndMissing(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()

	_, _, err := runScheduleCmd(t, sc, "table", false, "", "--id", "5", "--due", "4pm")
	if err == nil || !strings.Contains(err.Error(), "ambiguous") || !strings.Contains(err.Error(), "--type quiz") {
		t.Errorf("ambiguous id: %v", err)
	}
	if len(sc.puts()) != 0 {
		t.Errorf("writes = %v", sc.puts())
	}

	raw, _, err := runScheduleCmd(t, sc, "json", true, "", "--id", "5", "--type", "assignment", "--due", "4pm")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(itemKeys(decodeSchedule(t, raw)), " "); got != "assignment:5" {
		t.Errorf("--type assignment resolved %s", got)
	}

	_, _, err = runScheduleCmd(t, sc, "table", false, "", "--id", "999", "--due", "4pm")
	if err == nil || !strings.Contains(err.Error(), "no quiz or assignment with ID 999") {
		t.Errorf("missing id: %v", err)
	}
	_, _, err = runScheduleCmd(t, sc, "table", false, "", "--match", "nothing here", "--due", "4pm")
	if err == nil || !strings.Contains(err.Error(), `matches "nothing here"`) {
		t.Errorf("no match: %v", err)
	}
}

func TestSchedule_RefusesOrderViolation(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()

	// 802 closes 2026-09-11 11:59 PM EDT; a due date after that is refused.
	_, _, err := runScheduleCmd(t, sc, "table", false, "", "--id", "802", "--type", "assignment", "--due", "2026-09-13 4pm", "--force")
	if err == nil || !strings.Contains(err.Error(), "refusing the plan") || !strings.Contains(err.Error(), "due Sun 2026-09-13 4:00 PM EDT is after closed Fri 2026-09-11 11:59 PM EDT") {
		t.Errorf("order violation: %v", err)
	}
	if len(sc.puts()) != 0 {
		t.Errorf("a refused plan was written: %v", sc.puts())
	}
	// Available after due, on a fresh item.
	_, _, err = runScheduleCmd(t, sc, "table", false, "", "--id", "800", "--type", "assignment", "--available", "2026-09-13 4pm", "--due", "2026-09-13 3pm")
	if err == nil || !strings.Contains(err.Error(), "available Sun 2026-09-13 4:00 PM EDT is after due") {
		t.Errorf("available after due: %v", err)
	}
	// In a bulk plan one bad item refuses the whole plan.
	_, _, err = runScheduleCmd(t, sc, "table", false, "", "--match", "essay", "--due", "2026-09-13 4pm", "--force")
	if err == nil || !strings.Contains(err.Error(), `assignment 802 "Locked essay"`) {
		t.Errorf("bulk refusal: %v", err)
	}
	if len(sc.puts()) != 0 {
		t.Errorf("a refused bulk plan was written: %v", sc.puts())
	}
}

func TestSchedule_ClearAndMergeKeepsOtherDates(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()

	out, _, err := runScheduleCmd(t, sc, "table", false, "", "--id", "802", "--type", "assignment", "--clear", "closed", "--due", "2026-09-13 4pm")
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	puts := sc.puts()
	if len(puts) != 1 || puts[0] != `PUT /api/v1/courses/1/assignments/802 {"assignment":{"due_at":"2026-09-13T20:00:00Z","lock_at":null}}` {
		t.Errorf("writes = %v", puts)
	}
	if sc.assigns[802].unlock == nil || *sc.assigns[802].unlock != "2026-09-10T04:00:00Z" {
		t.Errorf("unlock_at was disturbed: %v", sc.assigns[802].unlock)
	}
	if !strings.Contains(out, "closed:    (clear)") || !strings.Contains(out, "Done: 1 matched, 1 changed, 1 written, 1 verified, 0 mismatched, 0 failed.") {
		t.Errorf("output:\n%s", out)
	}
}

func TestSchedule_DateOnlyDueMeansEndOfDay(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()

	_, _, err := runScheduleCmd(t, sc, "table", false, "", "--id", "800", "--type", "assignment", "--available", "9/9/26", "--due", "9/9/26", "--closed", "9/9/26")
	if err != nil {
		t.Fatal(err)
	}
	puts := sc.puts()
	if len(puts) != 1 || puts[0] != `PUT /api/v1/courses/1/assignments/800 {"assignment":{"due_at":"2026-09-10T03:59:59Z","lock_at":"2026-09-10T03:59:59Z","unlock_at":"2026-09-09T04:00:00Z"}}` {
		t.Errorf("writes = %v", puts)
	}
}

func TestSchedule_ConfirmPromptForMatch(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()

	out, _, err := runScheduleCmd(t, sc, "table", false, "n\n", "--match", "course lineup", "--due", "9/9/26")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "About to update 1 of 1 matched items in course 1. Continue? [y/N]:") || !strings.Contains(out, "Cancelled.") {
		t.Errorf("prompt output:\n%s", out)
	}
	if !strings.Contains(out, "PUT /api/v1/courses/1/assignments/800") {
		t.Errorf("the plan was not shown before the prompt:\n%s", out)
	}
	if len(sc.puts()) != 0 {
		t.Errorf("declined prompt wrote: %v", sc.puts())
	}

	out, _, err = runScheduleCmd(t, sc, "table", false, "y\n", "--match", "course lineup", "--due", "9/9/26")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if puts := sc.puts(); len(puts) != 1 || puts[0] != `PUT /api/v1/courses/1/assignments/800 {"assignment":{"due_at":"2026-09-10T03:59:59Z"}}` {
		t.Errorf("writes = %v", puts)
	}

	// --force skips the prompt; --id never prompts.
	sc2 := newScheduleCanvas(t)
	defer sc2.Close()
	out, _, err = runScheduleCmd(t, sc2, "table", false, "", "--match", "course lineup", "--due", "9/9/26", "--force")
	if err != nil || strings.Contains(out, "Continue?") || len(sc2.puts()) != 1 {
		t.Errorf("--force: err=%v puts=%v\n%s", err, sc2.puts(), out)
	}
}

func TestSchedule_MismatchExitsNonZero(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()
	sc.staleQuiz = 457

	out, _, err := runScheduleCmd(t, sc, "table", false, "", "--id", "457", "--type", "quiz", "--date", "2026-09-09", "--due", "4:50pm")
	if err == nil || !strings.Contains(err.Error(), "1 of 1 items did not read back as requested") {
		t.Errorf("err = %v", err)
	}
	if !strings.Contains(out, `quiz 457 "Attendance 9/16": due read back Wed 2025-12-31 7:00 PM EST (2026-01-01T00:00:00Z), requested Wed 2026-09-09 4:50 PM EDT (2026-09-09T20:50:00Z)`) || !strings.Contains(out, "0 verified, 1 mismatched, 0 failed") {
		t.Errorf("output:\n%s", out)
	}
	if len(sc.puts()) != 1 {
		t.Errorf("writes = %v", sc.puts())
	}
}

func TestSchedule_WriteFailureContinuesAndFails(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()
	sc.failAssignment = 801

	raw, _, err := runScheduleCmd(t, sc, "json", false, "", "--match", "attendance", "--date", "2026-09-09", "--due", "4:50pm", "--force")
	if err == nil || !strings.Contains(err.Error(), "1 write failures") {
		t.Errorf("err = %v", err)
	}
	res := decodeSchedule(t, raw)
	if res.Summary.Written != 2 || res.Summary.Verified != 2 || res.Summary.Failed != 1 || res.Summary.Mismatched != 0 {
		t.Errorf("summary = %+v", res.Summary)
	}
	for _, it := range res.Items {
		if it.ID == 801 && (it.Written || it.Verified != "no" || !strings.Contains(it.Error, "write failed")) {
			t.Errorf("failed item = %+v", it)
		}
		if it.Kind == "quiz" && it.Verified != "yes" {
			t.Errorf("quiz %d not verified: %+v", it.ID, it)
		}
	}
	if len(sc.puts()) != 3 {
		t.Errorf("writes = %v", sc.puts())
	}
}

func TestSchedule_NothingToChange(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()

	out, _, err := runScheduleCmd(t, sc, "table", false, "", "--id", "456", "--type", "quiz", "--available", "2026-09-09 3pm")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Nothing to change") || !strings.Contains(out, "(unchanged)") || len(sc.puts()) != 0 {
		t.Errorf("output:\n%s\nputs=%v", out, sc.puts())
	}
}

func TestSchedule_TimeInputErrors(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()

	_, _, err := runScheduleCmd(t, sc, "table", false, "", "--id", "456", "--due", "4:50")
	if err == nil || !strings.Contains(err.Error(), "invalid --due") || !strings.Contains(err.Error(), "ambiguous") {
		t.Errorf("ambiguous time: %v", err)
	}
	_, _, err = runScheduleCmd(t, sc, "table", false, "", "--id", "456", "--date", "2026-09-09 4pm", "--due", "4:50pm")
	if err == nil || !strings.Contains(err.Error(), "invalid --date") {
		t.Errorf("date with time: %v", err)
	}
	_, _, err = runScheduleCmd(t, sc, "table", false, "", "--id", "456", "--due", "4:50pm", "--timezone", "Mars/Olympus")
	if err == nil || !strings.Contains(err.Error(), "unknown time zone") {
		t.Errorf("bad zone: %v", err)
	}
	_, _, err = runScheduleCmd(t, sc, "csv", false, "", "--id", "456", "--due", "4:50pm")
	if err == nil || !strings.Contains(err.Error(), "unsupported output format") {
		t.Errorf("csv: %v", err)
	}
	if len(sc.puts()) != 0 {
		t.Errorf("writes = %v", sc.puts())
	}
}

func TestSchedule_TimezoneFromConfig(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()
	withConfigTimezone(t, "Europe/Berlin")
	// Env auth is set by runScheduleCmd after this, so the config only
	// supplies the zone.
	oldFormat, oldNoCache, oldNow := outputFormat, noCache, scheduleNow
	outputFormat, noCache = "table", true
	scheduleNow = func() time.Time { return scheduleFixtureNow }
	t.Cleanup(func() { outputFormat, noCache, scheduleNow = oldFormat, oldNoCache, oldNow })
	t.Setenv("CANVAS_URL", sc.URL)
	t.Setenv("CANVAS_TOKEN", "test-token")
	t.Setenv("CANVAS_REQUESTS_PER_SEC", "1000")

	cmd := newScheduleCmd()
	cmd.SilenceUsage, cmd.SilenceErrors = true, true
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"--course-id", "1", "--id", "800", "--type", "assignment", "--due", "2026-09-09 4:50pm"})
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("%v\n%s", err, out.String())
	}
	if puts := sc.puts(); len(puts) != 1 || puts[0] != `PUT /api/v1/courses/1/assignments/800 {"assignment":{"due_at":"2026-09-09T14:50:00Z"}}` {
		t.Errorf("writes = %v", puts)
	}
	if !strings.Contains(out.String(), "times in Europe/Berlin") {
		t.Errorf("output:\n%s", out.String())
	}
}

func TestSchedule_Validation(t *testing.T) {
	sc := newScheduleCanvas(t)
	defer sc.Close()
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--due", "4pm"}, "one of --id or --match"},
		{[]string{"--id", "1", "--match", "x", "--due", "4pm"}, "mutually exclusive"},
		{[]string{"--id", "1"}, "nothing to change"},
		{[]string{"--id", "1", "--due", "4pm", "--clear", "due"}, "mutually exclusive"},
		{[]string{"--id", "1", "--due", "4pm", "--type", "page"}, "unknown --type"},
	} {
		_, _, err := runScheduleCmd(t, sc, "table", false, "", c.args...)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%v: err = %v, want %q", c.args, err, c.want)
		}
	}
}

func TestSchedule_UnitHelpers(t *testing.T) {
	a := time.Date(2026, 9, 9, 20, 0, 0, 0, time.UTC)
	b := a.In(time.FixedZone("X", 3600))
	if !sameInstant(&a, &b) || sameInstant(&a, nil) || !sameInstant(nil, nil) {
		t.Error("sameInstant")
	}
	if nonZeroTime(time.Time{}) != nil || nonZeroTime(a) == nil {
		t.Error("nonZeroTime")
	}
	if copyTime(nil) != nil || !copyTime(&b).Equal(a) || copyTime(&b).Location() != time.UTC {
		t.Error("copyTime")
	}
	if formatLocalOrNone(nil, time.UTC) != "—" || formatUTCOrNone(nil) != "—" || describeOrNone(nil, time.UTC) != "—" {
		t.Error("none rendering")
	}
	if scheduleKindWord("quiz") != "quiz" || scheduleKindWord("assignment") != "assignment" || scheduleKindWord("all") != "quiz or assignment" {
		t.Error("scheduleKindWord")
	}
}
