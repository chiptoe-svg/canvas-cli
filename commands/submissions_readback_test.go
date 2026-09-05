package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

// gradeServer is a fake Canvas whose submissions remember what was PUT.
// With reflect=false a PUT answers as if it worked but changes nothing,
// so the read-back disagrees with the request.
type gradeServer struct {
	*httptest.Server
	mu      sync.Mutex
	reflect bool
	subs    map[string]*api.Submission // "<assignment>/<user>"
	nextID  int64
	puts    []string
	gets    []string
}

func newGradeServer(t *testing.T) *gradeServer {
	t.Helper()
	gs := &gradeServer{
		reflect: true,
		nextID:  100,
		subs: map[string]*api.Submission{
			"100/10": {ID: 1, AssignmentID: 100, UserID: 10, Score: 88, Grade: "88", EnteredScore: 88, EnteredGrade: "88", WorkflowState: "graded",
				User:               &api.User{ID: 10, Name: "Ada Lovelace"},
				SubmissionComments: []api.SubmissionComment{{ID: 50, AuthorName: "Teacher", Comment: "first draft"}}},
			"100/11": {ID: 2, AssignmentID: 100, UserID: 11, WorkflowState: "submitted", User: &api.User{ID: 11, Name: "Bob"}},
			"100/12": {ID: 3, AssignmentID: 100, UserID: 12, Score: 70, Grade: "70", EnteredScore: 70, EnteredGrade: "70", WorkflowState: "graded", User: &api.User{ID: 12, Name: "Cy"}},
		},
	}
	gs.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gs.mu.Lock()
		defer gs.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path
		switch {
		case path == "/api/v1/courses/1":
			_, _ = w.Write([]byte(`{"id": 1, "name": "Course"}`))
		case path == "/api/v1/accounts":
			_, _ = w.Write([]byte(`[]`))
		case strings.HasPrefix(path, "/api/v1/courses/1/assignments/"):
			parts := strings.Split(strings.TrimPrefix(path, "/api/v1/courses/1/assignments/"), "/")
			if len(parts) != 3 || parts[1] != "submissions" {
				http.NotFound(w, r)
				return
			}
			key := parts[0] + "/" + parts[2]
			sub, ok := gs.subs[key]
			if !ok {
				http.NotFound(w, r)
				return
			}
			view := *sub
			if r.Method == http.MethodPut {
				raw, _ := io.ReadAll(r.Body)
				gs.puts = append(gs.puts, key+" "+string(raw))
				var body struct {
					Submission struct {
						PostedGrade string `json:"posted_grade"`
						Excuse      bool   `json:"excuse"`
					} `json:"submission"`
					Comment struct {
						TextComment string `json:"text_comment"`
					} `json:"comment"`
				}
				_ = json.Unmarshal(raw, &body)
				applied := *sub
				applied.SubmissionComments = append([]api.SubmissionComment(nil), sub.SubmissionComments...)
				if g := body.Submission.PostedGrade; g != "" {
					applied.Grade, applied.EnteredGrade = g, g
					if f, err := strconv.ParseFloat(g, 64); err == nil {
						applied.Score, applied.EnteredScore = f, f
					}
					applied.WorkflowState = "graded"
				}
				if body.Submission.Excuse {
					applied.ExcusedTLN = true
				}
				if text := body.Comment.TextComment; text != "" {
					gs.nextID++
					applied.SubmissionComments = append(applied.SubmissionComments, api.SubmissionComment{ID: gs.nextID, AuthorName: "Teacher", AuthorID: 7, Comment: text})
				}
				if gs.reflect {
					*sub = applied
				}
				// the PUT response always claims success
				_ = json.NewEncoder(w).Encode(applied)
				return
			}
			gs.gets = append(gs.gets, key+"?"+r.URL.RawQuery)
			if r.URL.Query().Get("include[]") != "submission_comments" {
				view.SubmissionComments = nil
			}
			_ = json.NewEncoder(w).Encode(view)
		default:
			t.Logf("unhandled %s %s", r.Method, path)
			http.NotFound(w, r)
		}
	}))
	return gs
}

// useGradeServer points getAPIClient at gs (env auth, cache enabled).
func useGradeServer(t *testing.T, gs *gradeServer) {
	t.Helper()
	t.Setenv("CANVAS_URL", gs.URL)
	t.Setenv("CANVAS_TOKEN", "test-token")
	t.Setenv("CANVAS_REQUESTS_PER_SEC", "1000")
}

func runGradeCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := newSubmissionsGradeCmd()
	cmd.SetArgs(args)
	var err error
	out := captureStdout(func() { err = cmd.ExecuteContext(context.Background()) })
	return out, err
}

func TestSubmissionsGrade_ReadBackVerified(t *testing.T) {
	gs := newGradeServer(t)
	defer gs.Close()
	useGradeServer(t, gs)

	out, err := runGradeCmd(t, "--course-id", "1", "--assignment-id", "100", "--user-id", "10", "--score", "95", "--comment", "Great work, see the rubric")
	if err != nil {
		t.Fatalf("grade: %v\n%s", err, out)
	}
	for _, want := range []string{"Successfully graded submission for Ada Lovelace (user 10, assignment 100)", "grade: 88 → 95", `comment: #101 by Teacher — "Great work, see the rubric"`, "verified: yes"} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if len(gs.puts) != 1 || !strings.Contains(gs.puts[0], `"posted_grade":"95.00"`) {
		t.Errorf("puts = %v", gs.puts)
	}
	// one read before and one after, both live and with comments
	if len(gs.gets) != 2 {
		t.Fatalf("gets = %v, want a read before and a read after the PUT", gs.gets)
	}
	for _, g := range gs.gets {
		if !strings.Contains(g, "include%5B%5D=submission_comments") {
			t.Errorf("read without comments: %s", g)
		}
	}
}

func TestSubmissionsGrade_ReadBackMismatchFails(t *testing.T) {
	gs := newGradeServer(t)
	defer gs.Close()
	gs.reflect = false
	useGradeServer(t, gs)

	out, err := runGradeCmd(t, "--course-id", "1", "--assignment-id", "100", "--user-id", "10", "--score", "95")
	if err == nil {
		t.Fatalf("expected a mismatch error:\n%s", out)
	}
	if !strings.Contains(err.Error(), "did not read back as requested") || !strings.Contains(err.Error(), "score read back 88, requested 95") {
		t.Errorf("error = %v", err)
	}
	if !strings.Contains(out, "grade: 88 → 88") || !strings.Contains(out, "verified: no — score read back 88, requested 95") {
		t.Errorf("output = %s", out)
	}

	// a comment that does not come back is a mismatch too
	out, err = runGradeCmd(t, "--course-id", "1", "--assignment-id", "100", "--user-id", "12", "--comment", "Please revise")
	if err == nil || !strings.Contains(err.Error(), "comment not found on read-back") {
		t.Errorf("comment mismatch: err=%v\n%s", err, out)
	}

	// an excuse that does not stick
	_, err = runGradeCmd(t, "--course-id", "1", "--assignment-id", "100", "--user-id", "11", "--excuse")
	if err == nil || !strings.Contains(err.Error(), "not excused on read-back") {
		t.Errorf("excuse mismatch: err=%v", err)
	}
}

func TestSubmissionsGrade_JSONIncludesBeforeAfterVerified(t *testing.T) {
	gs := newGradeServer(t)
	defer gs.Close()
	useGradeServer(t, gs)
	setOutputFormat(t, "json")

	out, err := runGradeCmd(t, "--course-id", "1", "--assignment-id", "100", "--user-id", "11", "--posted-grade", "A")
	if err != nil {
		t.Fatalf("grade: %v\n%s", err, out)
	}
	var got struct {
		Before    *api.Submission `json:"before"`
		After     *api.Submission `json:"after"`
		Requested gradeRequest    `json:"requested"`
		Verified  bool            `json:"verified"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not json: %v\n%s", err, out)
	}
	if got.Before == nil || got.Before.Grade != "" || got.After == nil || got.After.Grade != "A" || !got.Verified || got.Requested.PostedGrade != "A" {
		t.Errorf("json = %+v", got)
	}

	gs.reflect = false
	out, err = runGradeCmd(t, "--course-id", "1", "--assignment-id", "100", "--user-id", "12", "--posted-grade", "B")
	if err == nil {
		t.Fatal("mismatch must fail")
	}
	var bad struct {
		Verified   bool     `json:"verified"`
		Mismatches []string `json:"mismatches"`
	}
	_ = json.Unmarshal([]byte(out), &bad)
	if bad.Verified || len(bad.Mismatches) != 1 || !strings.Contains(bad.Mismatches[0], `grade read back "70", requested "B"`) {
		t.Errorf("json on mismatch = %s", out)
	}
}

func TestSubmissionsAddComment_ReadBack(t *testing.T) {
	gs := newGradeServer(t)
	defer gs.Close()
	useGradeServer(t, gs)

	cmd := newSubmissionsAddCommentCmd()
	cmd.SetArgs([]string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10", "--text", "first draft"}) // same text as an existing comment
	var err error
	out := captureStdout(func() { err = cmd.ExecuteContext(context.Background()) })
	if err != nil {
		t.Fatalf("add-comment: %v\n%s", err, out)
	}
	// the new comment (#101), not the pre-existing #50 with the same text
	if !strings.Contains(out, "Comment added successfully to submission for user 10") || !strings.Contains(out, `comment: #101 by Teacher — "first draft"`) || !strings.Contains(out, "verified: yes") {
		t.Errorf("output = %s", out)
	}

	gs.reflect = false
	cmd = newSubmissionsAddCommentCmd()
	cmd.SetArgs([]string{"--course-id", "1", "--assignment-id", "100", "--user-id", "11", "--text", "Missing citation"})
	out = captureStdout(func() { err = cmd.ExecuteContext(context.Background()) })
	if err == nil || !strings.Contains(err.Error(), "comment not found on read-back") || !strings.Contains(out, "verified: no") {
		t.Errorf("mismatch: err=%v\n%s", err, out)
	}
}

func TestSubmissionsBulkGrade_ReadBackSummary(t *testing.T) {
	gs := newGradeServer(t)
	defer gs.Close()
	useGradeServer(t, gs)
	csvPath := filepath.Join(t.TempDir(), "grades.csv")
	_ = os.WriteFile(csvPath, []byte("user_id,assignment_id,grade,comment\n10,100,95,Great\n11,100,80,\n12,100,72,\n99,100,50,\n"), 0o600)

	run := func() (string, error) {
		cmd := newSubmissionsBulkGradeCmd()
		cmd.SetArgs([]string{"--course-id", "1", "--csv-file", csvPath})
		var err error
		out := captureStdout(func() { err = cmd.ExecuteContext(context.Background()) })
		return out, err
	}

	// user 99 does not exist: an error row; the other three read back
	out, err := run()
	if err == nil || !strings.Contains(err.Error(), "1 errors") {
		t.Errorf("expected the missing user to fail the run: %v", err)
	}
	if !strings.Contains(out, "3 graded, 3 verified, 0 mismatched") {
		t.Errorf("summary missing:\n%s", out)
	}
	if !strings.Contains(out, "88 → 95, verified") || !strings.Contains(out, "ungraded → 80, verified") {
		t.Errorf("per-row lines missing:\n%s", out)
	}

	// now the server stops applying grades: every row of a new CSV mismatches
	gs.mu.Lock()
	gs.reflect = false
	gs.mu.Unlock()
	csv2 := filepath.Join(t.TempDir(), "grades2.csv")
	_ = os.WriteFile(csv2, []byte("user_id,assignment_id,grade\n10,100,60\n11,100,61\n12,100,62\n"), 0o600)
	run2 := func() (string, error) {
		cmd := newSubmissionsBulkGradeCmd()
		cmd.SetArgs([]string{"--course-id", "1", "--csv-file", csv2})
		var err error
		out := captureStdout(func() { err = cmd.ExecuteContext(context.Background()) })
		return out, err
	}
	out, err = run2()
	if err == nil || !strings.Contains(err.Error(), "3 of 3 graded submissions did not read back as requested") {
		t.Errorf("expected a mismatch error: %v", err)
	}
	for _, want := range []string{"3 graded, 0 verified, 3 mismatched", "95 → 95, NOT verified: score read back 95, requested 60", "80 → 80, NOT verified: score read back 80, requested 61"} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}

	// json: summary counters and rows
	setOutputFormat(t, "json")
	out, err = run2()
	if err == nil {
		t.Error("json output must still exit non-zero on mismatch")
	}
	var got struct {
		Total, Success, Errors, Verified, Mismatched int
		Rows                                         []bulkGradeRow `json:"rows"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not json: %v\n%s", err, out)
	}
	if got.Total != 3 || got.Success != 3 || got.Errors != 0 || got.Verified != 0 || got.Mismatched != 3 || len(got.Rows) != 3 {
		t.Errorf("json summary = %+v", got)
	}
	for _, row := range got.Rows {
		if row.Verified || row.Error != "" || len(row.Mismatches) != 1 {
			t.Errorf("row = %+v", row)
		}
		if row.UserID == 11 && (row.Before != "80" || row.After != "80" || row.Requested != "61") {
			t.Errorf("row 11 = %+v", row)
		}
	}
	setOutputFormat(t, "table")

	// an error row is counted separately from mismatches
	csv3 := filepath.Join(t.TempDir(), "grades3.csv")
	_ = os.WriteFile(csv3, []byte("user_id,assignment_id,grade\n10,100,40\n99,100,50\n"), 0o600)
	cmd := newSubmissionsBulkGradeCmd()
	cmd.SetArgs([]string{"--course-id", "1", "--csv-file", csv3})
	out = captureStdout(func() { err = cmd.ExecuteContext(context.Background()) })
	if err == nil || !strings.Contains(err.Error(), "1 errors and 1 submissions that did not read back") {
		t.Errorf("errors + mismatches: %v\n%s", err, out)
	}
	if !strings.Contains(out, "1 graded, 0 verified, 1 mismatched") || !strings.Contains(out, "Errors: 1") {
		t.Errorf("summary:\n%s", out)
	}
}

// The read-back must bypass the response cache: getAPIClient enables it,
// and a cached "before" read would be returned as the "after" read.
func TestSubmissionsGrade_ReadBackIsLive(t *testing.T) {
	gs := newGradeServer(t)
	defer gs.Close()
	useGradeServer(t, gs)
	t.Setenv("HOME", t.TempDir()) // fresh cache dir
	out, err := runGradeCmd(t, "--course-id", "1", "--assignment-id", "100", "--user-id", "10", "--score", "91")
	if err != nil {
		t.Fatalf("grade: %v\n%s", err, out)
	}
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if len(gs.gets) != 2 {
		t.Errorf("the read-back must hit the server (got %d reads: %v)", len(gs.gets), gs.gets)
	}
}

func TestVerifyGradeReadBack(t *testing.T) {
	after := &api.Submission{Score: 9.5, EnteredScore: 10, Grade: "9.5", EnteredGrade: "10", ExcusedTLN: false,
		SubmissionComments: []api.SubmissionComment{{ID: 1, Comment: "hi"}, {ID: 2, Comment: "hi"}}}
	before := &api.Submission{SubmissionComments: []api.SubmissionComment{{ID: 1, Comment: "hi"}}}

	// numeric grade compares with the entered score (before late policy)
	if c, m := verifyGradeReadBack(before, after, gradeRequest{PostedGrade: "10.00", Comment: "hi"}); len(m) != 0 || c == nil || c.ID != 2 {
		t.Errorf("entered score / new comment: comment=%v mismatches=%v", c, m)
	}
	if _, m := verifyGradeReadBack(before, after, gradeRequest{PostedGrade: "9.5"}); len(m) != 1 {
		t.Errorf("late-policy score must not satisfy the entered score: %v", m)
	}
	// letter grade, case-insensitive, entered grade first
	if _, m := verifyGradeReadBack(nil, &api.Submission{Grade: "a", EnteredGrade: "A"}, gradeRequest{PostedGrade: "a"}); len(m) != 0 {
		t.Errorf("letter grade: %v", m)
	}
	if _, m := verifyGradeReadBack(nil, after, gradeRequest{Excuse: true}); len(m) != 1 || !strings.Contains(m[0], "not excused") {
		t.Errorf("excuse: %v", m)
	}
	if _, m := verifyGradeReadBack(nil, after, gradeRequest{PostedGrade: "7", Excuse: true, Comment: "nope"}); len(m) != 3 {
		t.Errorf("every difference is reported: %v", m)
	}
	// tolerance: Canvas rounds to two decimals
	if _, m := verifyGradeReadBack(nil, &api.Submission{Score: 33.33}, gradeRequest{PostedGrade: fmt.Sprint(100.0 / 3)}); len(m) != 0 {
		t.Errorf("rounding tolerance: %v", m)
	}
}

// A rubric-only grade has no posted grade, excuse, or comment to compare, so
// without checking the rubric rows the read-back would verify an empty
// request. Every requested criterion must read back as asked.
func TestVerifyGradeReadBack_Rubric(t *testing.T) {
	eight, four := 8.0, 4.0
	req := gradeRequest{Rubric: map[string]rubricCriterionRequest{"_1": {Points: &eight}, "_2": {Points: &four, Rating: "r_low"}}}

	// the false-verified case: Canvas returned no rubric at all
	if _, m := verifyGradeReadBack(nil, &api.Submission{Score: 12}, req); len(m) != 2 {
		t.Errorf("rubric absent on read-back must fail for both criteria: %v", m)
	}
	ok := &api.Submission{Rubric: api.RubricAssessmentResult{"_1": {Points: 8}, "_2": {Points: 4, RatingID: "r_low"}}}
	if _, m := verifyGradeReadBack(nil, ok, req); len(m) != 0 {
		t.Errorf("matching rubric: %v", m)
	}
	wrong := &api.Submission{Rubric: api.RubricAssessmentResult{"_1": {Points: 8}, "_2": {Points: 3, RatingID: "r_mid"}}}
	if _, m := verifyGradeReadBack(nil, wrong, req); len(m) != 2 || !strings.Contains(m[0]+m[1], "3 points") {
		t.Errorf("points and rating differences must both be reported: %v", m)
	}
	// rounding tolerance applies to rubric points too
	third := 100.0 / 3
	if _, m := verifyGradeReadBack(nil, &api.Submission{Rubric: api.RubricAssessmentResult{"_1": {Points: 33.33}}}, gradeRequest{Rubric: map[string]rubricCriterionRequest{"_1": {Points: &third}}}); len(m) != 0 {
		t.Errorf("rubric rounding tolerance: %v", m)
	}
}

// gradeAndReadBack must not post a comment that is already on the
// submission; the grade still goes through and the row still verifies.
func TestGradeAndReadBack_DoesNotRepostExistingComment(t *testing.T) {
	var putBodies []map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v1/accounts":
			fmt.Fprint(w, `[]`)
		case r.Method == http.MethodPut:
			var body map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&body)
			putBodies = append(putBodies, body)
			fmt.Fprint(w, `{"id": 1, "user_id": 7, "score": 9, "grade": "9"}`)
		default:
			fmt.Fprint(w, `{"id": 1, "user_id": 7, "score": 9, "grade": "9", "entered_score": 9,
				"submission_comments": [{"id": 55, "comment": "Nice work", "author_id": 3}]}`)
		}
	}))
	defer server.Close()
	client, err := api.NewClient(api.ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatal(err)
	}
	svc := api.NewSubmissionsService(client)

	rb, err := gradeAndReadBack(context.Background(), svc, 1, 2, 7, &api.GradeSubmissionParams{
		PostedGrade: "9", Comment: &api.SubmissionCommentParams{TextComment: "Nice work"},
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	if !rb.CommentExisted || rb.Comment == nil || rb.Comment.ID != 55 {
		t.Errorf("existing comment not recognised: existed=%v comment=%+v", rb.CommentExisted, rb.Comment)
	}
	if !rb.Verified {
		t.Errorf("row should verify, mismatches: %v", rb.Mismatches)
	}
	if len(putBodies) != 1 {
		t.Fatalf("expected exactly one PUT (the grade), got %d", len(putBodies))
	}
	if _, has := putBodies[0]["comment"]; has {
		t.Errorf("comment must not be re-posted, body: %v", putBodies[0])
	}

	// comment-only request against an existing comment: nothing to write at all
	putBodies = nil
	rb, err = gradeAndReadBack(context.Background(), svc, 1, 2, 7, &api.GradeSubmissionParams{
		Comment: &api.SubmissionCommentParams{TextComment: "Nice work"},
	}, true)
	if err != nil || !rb.Verified || !rb.CommentExisted || len(putBodies) != 0 {
		t.Errorf("comment-only re-run must write nothing: err=%v verified=%v existed=%v puts=%d", err, rb.Verified, rb.CommentExisted, len(putBodies))
	}

	// without the flag (a single grade / add-comment) the comment is posted as asked
	putBodies = nil
	rb, err = gradeAndReadBack(context.Background(), svc, 1, 2, 7, &api.GradeSubmissionParams{
		PostedGrade: "9", Comment: &api.SubmissionCommentParams{TextComment: "Nice work"},
	}, false)
	if err != nil || rb.CommentExisted || len(putBodies) != 1 {
		t.Fatalf("explicit comment must be posted: err=%v existed=%v puts=%d", err, rb.CommentExisted, len(putBodies))
	}
	if _, has := putBodies[0]["comment"]; !has {
		t.Errorf("comment missing from explicit write: %v", putBodies[0])
	}
}
