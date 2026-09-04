package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// excuseCanvas is a fake Canvas for `submissions excuse`: a roster with
// look-alike names, assignments including a quiz-backed one, and a
// submission store that a PUT mutates so the read-back sees it.
type excuseCanvas struct {
	*httptest.Server
	mu       sync.Mutex
	requests []string // "METHOD URI BODY"
	excused  map[string]bool
	stale    bool // read-backs ignore the write
}

func newExcuseCanvas(t *testing.T) *excuseCanvas {
	t.Helper()
	ec := &excuseCanvas{excused: map[string]bool{}}
	ec.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			_, _ = w.Write([]byte(`[]`))
			return
		}
		body, _ := io.ReadAll(r.Body)
		ec.mu.Lock()
		ec.requests = append(ec.requests, strings.TrimSpace(r.Method+" "+r.URL.RequestURI()+" "+string(body)))
		ec.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v1/courses/1/users":
			fmt.Fprint(w, `[
{"id":10,"name":"Ada Lovelace","sortable_name":"Lovelace, Ada","short_name":"Ada","login_id":"ada@example.edu"},
{"id":11,"name":"Ann Smith","sortable_name":"Smith, Ann","login_id":"asmith"},
{"id":12,"name":"Anne Smithers","sortable_name":"Smithers, Anne","login_id":"asmithers"},
{"id":13,"name":"Test Student","sortable_name":"Student, Test"}]`)
		case r.URL.Path == "/api/v1/courses/1/assignments":
			fmt.Fprint(w, `[
{"id":456,"name":"Quiz 3","quiz_id":77,"is_quiz_assignment":true,"submission_types":["online_quiz"]},
{"id":457,"name":"Essay 1","submission_types":["online_upload"]},
{"id":458,"name":"Essay 10","submission_types":["online_upload"]},
{"id":459,"name":"Course lineup","submission_types":["online_text_entry"]}]`)
		case strings.HasPrefix(r.URL.Path, "/api/v1/courses/1/assignments/") && strings.Contains(r.URL.Path, "/submissions/"):
			key := strings.TrimPrefix(r.URL.Path, "/api/v1/courses/1/assignments/")
			if r.Method == "PUT" && !ec.stale {
				var req map[string]map[string]bool
				_ = json.Unmarshal(body, &req)
				ec.mu.Lock()
				ec.excused[key] = req["submission"]["excuse"]
				ec.mu.Unlock()
			}
			ec.mu.Lock()
			ex := ec.excused[key]
			ec.mu.Unlock()
			parts := strings.Split(key, "/submissions/")
			fmt.Fprintf(w, `{"id":1,"assignment_id":%s,"user_id":%s,"excused":%t,"workflow_state":"unsubmitted"}`, parts[0], parts[1], ex)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.RequestURI())
			http.NotFound(w, r)
		}
	}))
	return ec
}

func (ec *excuseCanvas) puts() []string {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	var out []string
	for _, r := range ec.requests {
		if strings.HasPrefix(r, "PUT ") {
			out = append(out, r)
		}
	}
	return out
}

func runExcuse(t *testing.T, ec *excuseCanvas, format string, dry bool, stdin string, args ...string) (string, error) {
	t.Helper()
	t.Setenv("CANVAS_URL", ec.URL)
	t.Setenv("CANVAS_TOKEN", "test-token")
	t.Setenv("CANVAS_REQUESTS_PER_SEC", "1000")

	oldFormat, oldNoCache, oldDry := outputFormat, noCache, dryRun
	outputFormat, noCache, dryRun = format, true, dry
	t.Cleanup(func() { outputFormat, noCache, dryRun = oldFormat, oldNoCache, oldDry })

	cmd := newSubmissionsExcuseCmd()
	cmd.SilenceUsage, cmd.SilenceErrors = true, true
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(io.Discard)
	cmd.SetIn(strings.NewReader(stdin))
	cmd.SetArgs(append([]string{"--course-id", "1"}, args...))
	err := cmd.ExecuteContext(context.Background())
	return out.String(), err
}

const excusePut = `PUT /api/v1/courses/1/assignments/456/submissions/10 {"submission":{"excuse":true}}`

func TestSubmissionsExcuse_ByNameWithReadBack(t *testing.T) {
	ec := newExcuseCanvas(t)
	defer ec.Close()

	out, err := runExcuse(t, ec, "table", false, "", "--student", "ada lovelace", "--assignment", "quiz 3", "--force")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if puts := ec.puts(); len(puts) != 1 || puts[0] != excusePut {
		t.Errorf("writes = %v", puts)
	}
	for _, want := range []string{
		"student:    10 Ada Lovelace (Lovelace, Ada) <ada@example.edu>",
		"assignment: 456 Quiz 3 [quiz 77]",
		"excused:    not excused → excused",
		"verified:   yes",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
}

func TestSubmissionsExcuse_ResolutionForms(t *testing.T) {
	for _, c := range []struct{ student, assignment string }{
		{"10", "456"},
		{"10", "77"}, // quiz id → its assignment
		{"Lovelace, Ada", "QUIZ"},
		{"ada@example.edu", "quiz 3"},
	} {
		ec := newExcuseCanvas(t)
		out, err := runExcuse(t, ec, "table", false, "", "--student", c.student, "--assignment", c.assignment, "--force")
		if err != nil {
			t.Errorf("%v: %v\n%s", c, err, out)
		}
		if puts := ec.puts(); len(puts) != 1 || puts[0] != excusePut {
			t.Errorf("%v: writes = %v", c, puts)
		}
		ec.Close()
	}
}

func TestSubmissionsExcuse_RefusesAmbiguousOrUnknown(t *testing.T) {
	ec := newExcuseCanvas(t)
	defer ec.Close()

	_, err := runExcuse(t, ec, "table", false, "", "--student", "smith", "--assignment", "quiz 3", "--force")
	if err == nil || !strings.Contains(err.Error(), `--student: no student is named exactly "smith"`) ||
		!strings.Contains(err.Error(), "11 Ann Smith (Smith, Ann) <asmith>") || !strings.Contains(err.Error(), "12 Anne Smithers") {
		t.Errorf("ambiguous student: %v", err)
	}
	_, err = runExcuse(t, ec, "table", false, "", "--student", "Ann Smith", "--assignment", "essay", "--force")
	if err == nil || !strings.Contains(err.Error(), `--assignment: "essay" matches 2 assignments`) || !strings.Contains(err.Error(), "458 Essay 10") {
		t.Errorf("ambiguous assignment: %v", err)
	}
	// part of a name is not enough for a write, even when only one student matches it
	_, err = runExcuse(t, ec, "table", false, "", "--student", "lovel", "--assignment", "quiz 3", "--force")
	if err == nil || !strings.Contains(err.Error(), `--student: no student is named exactly "lovel"`) || !strings.Contains(err.Error(), "10 Ada Lovelace") {
		t.Errorf("partial student name: %v", err)
	}
	if puts := ec.puts(); len(puts) != 0 {
		t.Errorf("partial name must not write: %v", puts)
	}
	_, err = runExcuse(t, ec, "table", false, "", "--student", "nobody", "--assignment", "quiz 3", "--force")
	if err == nil || !strings.Contains(err.Error(), `no student matches "nobody" (3 searched)`) {
		t.Errorf("unknown student: %v", err)
	}
	_, err = runExcuse(t, ec, "table", false, "", "--student", "Test Student", "--assignment", "quiz 3", "--force")
	if err == nil || !strings.Contains(err.Error(), "no student matches") {
		t.Errorf("test student must not resolve: %v", err)
	}
	_, err = runExcuse(t, ec, "table", false, "", "--student", "10", "--assignment", "999", "--force")
	if err == nil || !strings.Contains(err.Error(), `no assignment matches "999"`) {
		t.Errorf("unknown assignment: %v", err)
	}
	if len(ec.puts()) != 0 {
		t.Errorf("a refused resolution wrote: %v", ec.puts())
	}
}

func TestSubmissionsExcuse_UnexcuseAndAlready(t *testing.T) {
	ec := newExcuseCanvas(t)
	defer ec.Close()
	ec.excused["456/submissions/10"] = true

	out, err := runExcuse(t, ec, "table", false, "", "--student", "10", "--assignment", "456", "--force")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "excused:    excused (already excused; nothing to do)") || len(ec.puts()) != 0 {
		t.Errorf("already excused:\n%s\nputs=%v", out, ec.puts())
	}

	out, err = runExcuse(t, ec, "table", false, "", "--student", "10", "--assignment", "456", "--unexcuse", "--force")
	if err != nil {
		t.Fatal(err)
	}
	if puts := ec.puts(); len(puts) != 1 || puts[0] != `PUT /api/v1/courses/1/assignments/456/submissions/10 {"submission":{"excuse":false}}` {
		t.Errorf("writes = %v", puts)
	}
	if !strings.Contains(out, "excused:    excused → not excused") || !strings.Contains(out, "verified:   yes") {
		t.Errorf("output:\n%s", out)
	}
}

func TestSubmissionsExcuse_DryRunAndPrompt(t *testing.T) {
	ec := newExcuseCanvas(t)
	defer ec.Close()

	out, err := runExcuse(t, ec, "table", true, "", "--student", "ada", "--assignment", "lineup")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"assignment: 459 Course lineup",
		"excused:    not excused → excused (planned)",
		`request:    PUT /api/v1/courses/1/assignments/459/submissions/10 {"submission":{"excuse":true}}`,
		"DRY RUN: no changes were made.",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("dry run lacks %q:\n%s", want, out)
		}
	}
	if len(ec.puts()) != 0 {
		t.Errorf("dry run wrote: %v", ec.puts())
	}

	out, err = runExcuse(t, ec, "table", false, "n\n", "--student", "ada", "--assignment", "lineup")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `Excuse Ada Lovelace from "Course lineup" (assignment 459)? [y/N]:`) || !strings.Contains(out, "Cancelled.") || len(ec.puts()) != 0 {
		t.Errorf("declined prompt:\n%s\nputs=%v", out, ec.puts())
	}
	out, err = runExcuse(t, ec, "table", false, "", "--student", "ada", "--assignment", "lineup")
	if err != nil || len(ec.puts()) != 0 || !strings.Contains(out, "Cancelled.") {
		t.Errorf("empty stdin must not write: %v %v", err, ec.puts())
	}
	_, err = runExcuse(t, ec, "table", false, "yes\n", "--student", "ada", "--assignment", "lineup")
	if err != nil || len(ec.puts()) != 1 {
		t.Errorf("accepted prompt: %v %v", err, ec.puts())
	}
}

func TestSubmissionsExcuse_MismatchExitsNonZero(t *testing.T) {
	ec := newExcuseCanvas(t)
	defer ec.Close()
	ec.stale = true

	out, err := runExcuse(t, ec, "table", false, "", "--student", "10", "--assignment", "456", "--force")
	if err == nil || !strings.Contains(err.Error(), "did not read back as requested: excused read back not excused, requested excused") {
		t.Errorf("err = %v", err)
	}
	if !strings.Contains(out, "excused:    not excused → not excused") || !strings.Contains(out, "verified:   no — excused read back not excused, requested excused") {
		t.Errorf("output:\n%s", out)
	}
	if len(ec.puts()) != 1 {
		t.Errorf("writes = %v", ec.puts())
	}
}

func TestSubmissionsExcuse_JSON(t *testing.T) {
	ec := newExcuseCanvas(t)
	defer ec.Close()

	raw, err := runExcuse(t, ec, "json", false, "", "--student", "10", "--assignment", "77", "--force")
	if err != nil {
		t.Fatal(err)
	}
	var res excuseResult
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		t.Fatalf("not json: %v\n%s", err, raw)
	}
	if res.Student.ID != 10 || res.Assignment.ID != 456 || res.Assignment.QuizID != 77 || !res.Requested ||
		res.Before == nil || *res.Before || res.After == nil || !*res.After || !res.Written || !res.Verified || res.DryRun {
		t.Errorf("result = %+v", res)
	}
	if strings.Contains(raw, "student:") {
		t.Errorf("json output carries prose:\n%s", raw)
	}

	_, err = runExcuse(t, ec, "csv", false, "", "--student", "10", "--assignment", "77", "--force")
	if err == nil || !strings.Contains(err.Error(), "unsupported output format") {
		t.Errorf("csv: %v", err)
	}
}

func TestSubmissionsExcuse_Validation(t *testing.T) {
	ec := newExcuseCanvas(t)
	defer ec.Close()
	if _, err := runExcuse(t, ec, "table", false, "", "--student", " ", "--assignment", "x"); err == nil || !strings.Contains(err.Error(), "--student is required") {
		t.Errorf("blank student: %v", err)
	}
	if _, err := runExcuse(t, ec, "table", false, "", "--student", "x"); err == nil {
		t.Error("missing --assignment accepted")
	}
	if excusedState(nil) != "?" || excusedWord(true) != "excused" {
		t.Error("helpers")
	}
}
