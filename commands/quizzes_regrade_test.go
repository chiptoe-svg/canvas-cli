package commands

import (
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

	"github.com/chiptoe-svg/canvas-cli/commands/internal/options"
	cmdtest "github.com/chiptoe-svg/canvas-cli/commands/internal/testing"
	"github.com/chiptoe-svg/canvas-cli/internal/activity"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

// ---- pure computation ----

func regradeFixtureQuestion() *api.QuizQuestion {
	return &api.QuizQuestion{
		ID:             789,
		QuizID:         10,
		QuestionType:   "multiple_choice_question",
		PointsPossible: 4,
		Answers: []api.QuizAnswer{
			{ID: 1001, Text: "old correct", Weight: 100},
			{ID: 1002, Text: "new correct"}, // weight absent in Canvas JSON == 0
			{ID: 1003, Text: "distractor", Weight: 0},
		},
	}
}

func TestRegradeAnswers(t *testing.T) {
	q := regradeFixtureQuestion()
	got, err := regradeAnswers(q, 1002)
	if err != nil {
		t.Fatalf("regradeAnswers: %v", err)
	}
	want := map[int64]float64{1001: 0, 1002: 100, 1003: 0}
	for _, a := range got {
		if a.Weight != want[a.ID] {
			t.Errorf("answer %d weight = %v, want %v", a.ID, a.Weight, want[a.ID])
		}
	}
	// the input must not be mutated
	if q.Answers[0].Weight != 100 || q.Answers[1].Weight != 0 {
		t.Errorf("regradeAnswers mutated its input: %+v", q.Answers)
	}

	if _, err := regradeAnswers(q, 9999); err == nil || !strings.Contains(err.Error(), "1001, 1002, 1003") {
		t.Errorf("expected an error listing available answers, got %v", err)
	}

	tf := &api.QuizQuestion{ID: 1, QuestionType: "true_false_question", Answers: []api.QuizAnswer{{ID: 1, Weight: 100}, {ID: 2}}}
	if _, err := regradeAnswers(tf, 2); err != nil {
		t.Errorf("true_false_question should be supported: %v", err)
	}

	for _, typ := range []string{"essay_question", "multiple_answers_question", "numerical_question", ""} {
		bad := &api.QuizQuestion{ID: 1, QuestionType: typ, Answers: []api.QuizAnswer{{ID: 1}}}
		if _, err := regradeAnswers(bad, 1); err == nil {
			t.Errorf("question type %q should be refused", typ)
		}
	}
}

func TestCorrectAnswerIDs(t *testing.T) {
	got := correctAnswerIDs(regradeFixtureQuestion().Answers)
	if len(got) != 1 || got[0] != 1001 {
		t.Errorf("correctAnswerIDs = %v, want [1001]", got)
	}
	if got := correctAnswerIDs([]api.QuizAnswer{{ID: 5}, {ID: 6, Weight: 50}}); len(got) != 0 {
		t.Errorf("no answer at weight 100 should yield none, got %v", got)
	}
}

func TestSelectedAnswerID(t *testing.T) {
	tests := []struct {
		raw    string
		wantID int64
		wantOK bool
	}{
		{`1002`, 1002, true},
		{`"1002"`, 1002, true},
		{`" 1002 "`, 1002, true},
		{`1002.0`, 1002, true},
		{`null`, 0, false},
		{``, 0, false},
		{`""`, 0, false},
		{`"abc"`, 0, false},
		{`0`, 0, false},
		{`-3`, 0, false},
		{`[1,2]`, 0, false},
		{`{"a":1}`, 0, false},
		{`1002.5`, 0, false},
	}
	for _, tt := range tests {
		id, ok := selectedAnswerID(json.RawMessage(tt.raw))
		if id != tt.wantID || ok != tt.wantOK {
			t.Errorf("selectedAnswerID(%s) = (%d, %v), want (%d, %v)", tt.raw, id, ok, tt.wantID, tt.wantOK)
		}
	}
}

func TestSelectRegradeSubmissions(t *testing.T) {
	subs := []api.QuizSubmission{
		{ID: 1, WorkflowState: "complete"},
		{ID: 2, WorkflowState: "pending_review"},
		{ID: 3, WorkflowState: "untaken"},
		{ID: 4, WorkflowState: "settings_only"},
		{ID: 5, WorkflowState: "preview"},
	}
	ids := func(s []api.QuizSubmission) []int64 {
		var out []int64
		for _, x := range s {
			out = append(out, x.ID)
		}
		return out
	}
	if got := ids(selectRegradeSubmissions(subs, "completed")); fmt.Sprint(got) != "[1]" {
		t.Errorf("completed = %v, want [1]", got)
	}
	if got := ids(selectRegradeSubmissions(subs, "")); fmt.Sprint(got) != "[1]" {
		t.Errorf("default = %v, want [1]", got)
	}
	if got := ids(selectRegradeSubmissions(subs, "all")); fmt.Sprint(got) != "[1 2]" {
		t.Errorf("all = %v, want [1 2]", got)
	}
}

// historyFixture mirrors the shape observed live from
// GET /courses/:id/assignments/:id/submissions/:uid?include[]=submission_history
// for a two-attempt classic quiz: answer_id numeric on the latest attempt,
// numeric string on the older one.
func historyFixture() *api.Submission {
	return &api.Submission{
		ID: 9001, UserID: 12, Attempt: 2, Score: 9,
		SubmissionHistory: []api.Submission{
			{Attempt: 1, Score: 5, SubmissionData: []api.QuizSubmissionData{
				{QuestionID: 1, AnswerID: json.RawMessage(`7`), Correct: json.RawMessage(`true`), Points: 1},
				{QuestionID: 789, AnswerID: json.RawMessage(`"1002"`), Correct: json.RawMessage(`false`), Points: 0},
			}},
			{Attempt: 2, Score: 9, SubmissionData: []api.QuizSubmissionData{
				{QuestionID: 1, AnswerID: json.RawMessage(`7`), Correct: json.RawMessage(`true`), Points: 1},
				{QuestionID: 789, AnswerID: json.RawMessage(`1001`), Correct: json.RawMessage(`true`), Points: 4},
			}},
			{Attempt: 3, Score: 0}, // in-progress / no data: never rescored
		},
	}
}

func TestSelectRegradeAttempts(t *testing.T) {
	sub := api.QuizSubmission{ID: 502, UserID: 12, Attempt: 2, WorkflowState: "complete", Score: 9}
	asub := historyFixture()

	latest := selectRegradeAttempts(sub, asub, "completed")
	if len(latest) != 1 || latest[0].Attempt != 2 || latest[0].OldScore != 9 || latest[0].Source != regradeSourceHistory {
		t.Errorf("completed = %+v, want the attempt-2 history entry", latest)
	}

	all := selectRegradeAttempts(sub, asub, "all")
	if len(all) != 2 || all[0].Attempt != 1 || all[1].Attempt != 2 || all[0].OldScore != 5 {
		t.Errorf("all = %+v, want attempts 1 and 2 (attempt 3 has no data)", all)
	}

	// no history at all, data on the top-level object only
	flat := &api.Submission{UserID: 12, Attempt: 2, Score: 9, SubmissionData: asub.SubmissionHistory[1].SubmissionData}
	if got := selectRegradeAttempts(sub, flat, "completed"); len(got) != 1 || got[0].Attempt != 2 {
		t.Errorf("top-level submission_data should be used, got %+v", got)
	}

	// nothing usable
	if got := selectRegradeAttempts(sub, &api.Submission{UserID: 12}, "completed"); len(got) != 0 {
		t.Errorf("no data should select nothing, got %+v", got)
	}
	// latest attempt has no data (not graded) -> nothing, even if older attempts do
	older := &api.Submission{UserID: 12, Attempt: 3, SubmissionHistory: asub.SubmissionHistory}
	if got := selectRegradeAttempts(api.QuizSubmission{Attempt: 3}, older, "completed"); len(got) != 0 {
		t.Errorf("latest attempt without data should select nothing, got %+v", got)
	}
}

func TestPlanRegradeRow_History(t *testing.T) {
	q := regradeFixtureQuestion()
	sub := api.QuizSubmission{ID: 502, UserID: 12, Attempt: 2, WorkflowState: "complete", Score: 9}
	attempts := selectRegradeAttempts(sub, historyFixture(), "all")

	rows := []quizRegradeRow{
		planRegradeRow(sub, attempts[0], q, 1002),
		planRegradeRow(sub, attempts[1], q, 1002),
	}
	// attempt 1 picked "1002" (string), 0 points before -> +4
	if r := rows[0]; r.Attempt != 1 || r.SelectedAnswerID != 1002 || r.OldQuestionScore != 0 || r.NewQuestionScore != 4 || r.ExpectedScore != 9 || !r.Changed {
		t.Errorf("attempt 1 row = %+v", r)
	}
	// attempt 2 picked 1001 (number), 4 points before -> -4
	if r := rows[1]; r.Attempt != 2 || r.SelectedAnswerID != 1001 || r.OldQuestionScore != 4 || r.NewQuestionScore != 0 || r.ExpectedScore != 5 || !r.Changed {
		t.Errorf("attempt 2 row = %+v", r)
	}

	// manual partial credit: points (2) beats the key-derived value (4)
	partial := regradeAttempt{Attempt: 1, OldScore: 6, Source: regradeSourceHistory, Data: []api.QuizSubmissionData{
		{QuestionID: 789, AnswerID: json.RawMessage(`1001`), Correct: json.RawMessage(`"partial"`), Points: 2},
	}}
	if r := planRegradeRow(api.QuizSubmission{ID: 508, Attempt: 1, Score: 6}, partial, q, 1002); r.OldQuestionScore != 2 || r.ExpectedScore != 4 || !r.Changed {
		t.Errorf("partial-credit row = %+v, want old 2 / expected 4", r)
	}

	// submission_data.correct is stale after a key change (observed live:
	// correct:true with points 0) and must never drive the computation:
	// correct:true on the now-wrong answer still scores 0, correct:false on
	// the now-right answer still scores full points.
	staleCorrect := regradeAttempt{Attempt: 1, OldScore: 9, Source: regradeSourceHistory, Data: []api.QuizSubmissionData{
		{QuestionID: 789, AnswerID: json.RawMessage(`1001`), Correct: json.RawMessage(`true`), Points: 4},
	}}
	if r := planRegradeRow(api.QuizSubmission{ID: 509, Attempt: 1, Score: 9}, staleCorrect, q, 1002); r.NewQuestionScore != 0 || r.ExpectedScore != 5 {
		t.Errorf("stale correct:true must not keep the points: %+v", r)
	}
	staleWrong := regradeAttempt{Attempt: 1, OldScore: 5, Source: regradeSourceHistory, Data: []api.QuizSubmissionData{
		{QuestionID: 789, AnswerID: json.RawMessage(`1002`), Correct: json.RawMessage(`false`), Points: 0},
	}}
	if r := planRegradeRow(api.QuizSubmission{ID: 509, Attempt: 1, Score: 5}, staleWrong, q, 1002); r.NewQuestionScore != 4 || r.ExpectedScore != 9 {
		t.Errorf("stale correct:false must not withhold the points: %+v", r)
	}

	// question not answered in this attempt
	none := regradeAttempt{Attempt: 1, OldScore: 2, Source: regradeSourceHistory, Data: []api.QuizSubmissionData{
		{QuestionID: 1, AnswerID: json.RawMessage(`7`), Points: 1},
	}}
	if r := planRegradeRow(api.QuizSubmission{ID: 506, Attempt: 1, Score: 2}, none, q, 1002); r.SelectedAnswerID != 0 || r.Changed || r.ExpectedScore != 2 {
		t.Errorf("unanswered row = %+v", r)
	}
}

func TestFallbackAttemptFromQuizQuestions(t *testing.T) {
	q := regradeFixtureQuestion()
	oldCorrect := correctAnswerIDs(q.Answers)
	sub := api.QuizSubmission{ID: 507, UserID: 17, Attempt: 1, Score: 4}

	// student-side shape: answer present
	a := fallbackAttemptFromQuizQuestions(sub, q, oldCorrect, []api.QuizSubmissionQuestion{
		{ID: 1, Answer: json.RawMessage(`7`)},
		{ID: 789, Answer: json.RawMessage(`1001`)},
	})
	if a.Source != regradeSourceQuizQuestions || a.SelectedAnswerID != 1001 || a.OldQuestionScore != 4 || a.OldScore != 4 {
		t.Errorf("fallback = %+v", a)
	}
	row := planRegradeRow(sub, a, q, 1002)
	if row.AnswerSource != regradeSourceQuizQuestions || row.ExpectedScore != 0 || !row.Changed {
		t.Errorf("fallback row = %+v", row)
	}

	// grader-side shape: {id, correct, flagged} without answer -> unknown
	b := fallbackAttemptFromQuizQuestions(sub, q, oldCorrect, []api.QuizSubmissionQuestion{{ID: 789}})
	if b.Source != regradeSourceNone || b.SelectedAnswerID != 0 {
		t.Errorf("grader-side fallback = %+v, want no answer", b)
	}
	if r := planRegradeRow(sub, b, q, 1002); r.Changed {
		t.Errorf("unknown answer must never change a score: %+v", r)
	}
}

// ---- end-to-end against a stateful fixture server ----

type historyEntry struct {
	attempt int
	score   float64
	data    []api.QuizSubmissionData // nil => no submission_data (never graded here)
}

type staleState struct {
	left  int
	score float64
}

type regradeServer struct {
	*httptest.Server
	mu           sync.Mutex
	questionType string
	weights      map[int64]float64 // answer id -> weight
	subs         map[int64]*api.QuizSubmission
	order        []int64
	history      map[int64][]*historyEntry // user id -> attempts
	studentView  map[int64]json.RawMessage // quiz submission id -> "answer" on /quiz_submissions/:id/questions
	ignoreScores bool                      // simulate Canvas not applying the update
	staleReads   int                       // reads after a PUT that still return the previous score
	stale        map[string]*staleState    // "<sub>/<attempt>" -> pending stale reads
	attemptGETs  []string
	questionPUTs []map[string]interface{}
	subPUTs      map[int64][]map[string]interface{}
	listPages    []string
	historyGETs  []string
	capPerPage   int // when > 0, the server never returns more than this per page whatever was asked
}

func data(q int64, answer string, points float64) api.QuizSubmissionData {
	correct := `false`
	if points > 0 {
		correct = `true`
	}
	return api.QuizSubmissionData{QuestionID: q, AnswerID: json.RawMessage(answer), Correct: json.RawMessage(correct), Points: points}
}

func newRegradeServer(t *testing.T) *regradeServer {
	t.Helper()
	other := data(1, `7`, 1) // an unrelated question answered correctly
	rs := &regradeServer{
		questionType: "multiple_choice_question",
		weights:      map[int64]float64{1001: 100, 1002: 0, 1003: 0},
		subs: map[int64]*api.QuizSubmission{
			501: {ID: 501, QuizID: 10, UserID: 11, Attempt: 1, WorkflowState: "complete", Score: 6},
			502: {ID: 502, QuizID: 10, UserID: 12, Attempt: 2, WorkflowState: "complete", Score: 9},
			503: {ID: 503, QuizID: 10, UserID: 13, Attempt: 1, WorkflowState: "complete", Score: 3},
			504: {ID: 504, QuizID: 10, UserID: 14, Attempt: 0, WorkflowState: "untaken", Score: 0},
			505: {ID: 505, QuizID: 10, UserID: 15, Attempt: 1, WorkflowState: "pending_review", Score: 7},
			506: {ID: 506, QuizID: 10, UserID: 16, Attempt: 1, WorkflowState: "complete", Score: 2},
			507: {ID: 507, QuizID: 10, UserID: 17, Attempt: 1, WorkflowState: "complete", Score: 4},
			508: {ID: 508, QuizID: 10, UserID: 18, Attempt: 1, WorkflowState: "complete", Score: 6},
		},
		order: []int64{501, 502, 503, 504, 505, 506, 507, 508},
		history: map[int64][]*historyEntry{
			11: {{1, 6, []api.QuizSubmissionData{other, data(789, `1002`, 0)}}},
			12: {{1, 5, []api.QuizSubmissionData{other, data(789, `"1002"`, 0)}}, {2, 9, []api.QuizSubmissionData{other, data(789, `"1001"`, 4)}}},
			13: {{1, 3, []api.QuizSubmissionData{other, data(789, `1003`, 0)}}},
			15: {{1, 7, []api.QuizSubmissionData{other, data(789, `1002`, 0)}}},
			16: {{1, 2, []api.QuizSubmissionData{other}}},                       // question 789 not answered
			17: {{1, 4, nil}},                                                   // no submission_data at all -> fallback
			18: {{1, 6, []api.QuizSubmissionData{other, data(789, `1001`, 2)}}}, // manual partial credit
		},
		studentView: map[int64]json.RawMessage{507: json.RawMessage(`1002`)},
		subPUTs:     map[int64][]map[string]interface{}{},
		stale:       map[string]*staleState{},
	}

	rs.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rs.mu.Lock()
		defer rs.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path

		switch {
		case path == "/api/v1/accounts":
			_, _ = w.Write([]byte(`[]`))

		case path == "/api/v1/courses/1/quizzes/10":
			_, _ = w.Write([]byte(`{"id": 10, "title": "Quiz", "quiz_type": "assignment", "assignment_id": 77}`))

		case path == "/api/v1/courses/1/quizzes/10/questions/789":
			if r.Method == http.MethodPut {
				var body map[string]interface{}
				raw, _ := io.ReadAll(r.Body)
				_ = json.Unmarshal(raw, &body)
				rs.questionPUTs = append(rs.questionPUTs, body)
				q, _ := body["question"].(map[string]interface{})
				answers, _ := q["answers"].([]interface{})
				for _, a := range answers {
					m, _ := a.(map[string]interface{})
					id, _ := m["id"].(float64)
					wgt, _ := m["weight"].(float64)
					rs.weights[int64(id)] = wgt
				}
			}
			_ = json.NewEncoder(w).Encode(rs.question())

		case path == "/api/v1/courses/1/quizzes/10/submissions":
			rs.listPages = append(rs.listPages, r.URL.RawQuery)
			page, _ := strconv.Atoi(r.URL.Query().Get("page"))
			perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
			if rs.capPerPage > 0 && perPage > rs.capPerPage {
				perPage = rs.capPerPage
			}
			start := (page - 1) * perPage
			out := []api.QuizSubmission{}
			for i := start; i < start+perPage && i < len(rs.order); i++ {
				out = append(out, *rs.subs[rs.order[i]])
			}
			_ = json.NewEncoder(w).Encode(api.QuizSubmissionsResponse{QuizSubmissions: out})

		case strings.HasPrefix(path, "/api/v1/courses/1/assignments/77/submissions/"):
			rs.historyGETs = append(rs.historyGETs, path+"?"+r.URL.RawQuery)
			uid, _ := strconv.ParseInt(strings.TrimPrefix(path, "/api/v1/courses/1/assignments/77/submissions/"), 10, 64)
			var sub *api.QuizSubmission
			for _, s := range rs.subs {
				if s.UserID == uid {
					sub = s
				}
			}
			if sub == nil {
				http.NotFound(w, r)
				return
			}
			asub := api.Submission{ID: 9000 + uid, UserID: uid, AssignmentID: 77, Attempt: sub.Attempt, Score: sub.Score}
			if r.URL.Query().Get("include[]") == "submission_history" {
				for _, h := range rs.history[uid] {
					asub.SubmissionHistory = append(asub.SubmissionHistory, api.Submission{
						UserID: uid, Attempt: h.attempt, Score: h.score, SubmissionData: h.data,
					})
				}
			}
			_ = json.NewEncoder(w).Encode(asub)

		case strings.HasPrefix(path, "/api/v1/quiz_submissions/") && strings.HasSuffix(path, "/questions"):
			id, _ := strconv.ParseInt(strings.TrimSuffix(strings.TrimPrefix(path, "/api/v1/quiz_submissions/"), "/questions"), 10, 64)
			// Grader-side shape by default: {id, correct, flagged}, no answer.
			qs := []api.QuizSubmissionQuestion{{ID: 1}, {ID: 789}}
			if ans, ok := rs.studentView[id]; ok {
				qs[1].Answer = ans
			}
			_ = json.NewEncoder(w).Encode(api.QuizSubmissionQuestionsResponse{QuizSubmissionQuestions: qs})

		case strings.HasPrefix(path, "/api/v1/courses/1/quizzes/10/submissions/"):
			id, _ := strconv.ParseInt(strings.TrimPrefix(path, "/api/v1/courses/1/quizzes/10/submissions/"), 10, 64)
			sub, ok := rs.subs[id]
			if !ok {
				http.NotFound(w, r)
				return
			}
			if r.Method == http.MethodPut {
				var body map[string]interface{}
				raw, _ := io.ReadAll(r.Body)
				_ = json.Unmarshal(raw, &body)
				rs.subPUTs[id] = append(rs.subPUTs[id], body)
				entries, _ := body["quiz_submissions"].([]interface{})
				entry, _ := entries[0].(map[string]interface{})
				attempt, _ := entry["attempt"].(float64)
				if rs.staleReads > 0 {
					rs.stale[fmt.Sprintf("%d/%d", id, int(attempt))] = &staleState{left: rs.staleReads, score: rs.attemptScore(sub, int(attempt))}
				}
				if !rs.ignoreScores {
					questions, _ := entry["questions"].(map[string]interface{})
					q789, _ := questions["789"].(map[string]interface{})
					if newScore, ok := q789["score"].(float64); ok {
						// Canvas recomputes the attempt's score from its per-question points.
						rs.applyScore(sub, int(attempt), newScore)
					}
				}
				_ = json.NewEncoder(w).Encode(api.QuizSubmissionsResponse{QuizSubmissions: []api.QuizSubmission{*sub}})
				return
			}
			// GET: ?attempt=N reports that attempt's score (observed live);
			// without it the object describes the latest attempt.
			view := *sub
			if a := r.URL.Query().Get("attempt"); a != "" {
				rs.attemptGETs = append(rs.attemptGETs, path+"?"+r.URL.RawQuery)
				n, _ := strconv.Atoi(a)
				view.Attempt = n
				view.Score = rs.attemptScore(sub, n)
				if st, ok := rs.stale[fmt.Sprintf("%d/%d", id, n)]; ok && st.left > 0 {
					st.left--
					view.Score = st.score
				}
			}
			_ = json.NewEncoder(w).Encode(api.QuizSubmissionsResponse{QuizSubmissions: []api.QuizSubmission{view}})

		default:
			t.Logf("unhandled path %s %s", r.Method, path)
			http.NotFound(w, r)
		}
	}))
	return rs
}

// attemptScore is the stored score of one attempt (history entry, or the
// quiz submission itself when there is no history record).
func (rs *regradeServer) attemptScore(sub *api.QuizSubmission, attempt int) float64 {
	for _, h := range rs.history[sub.UserID] {
		if h.attempt == attempt {
			return h.score
		}
	}
	return sub.Score
}

// applyScore updates the stored per-question points of one attempt and the
// attempt's score; the quiz submission score follows the latest attempt.
func (rs *regradeServer) applyScore(sub *api.QuizSubmission, attempt int, newScore float64) {
	for _, h := range rs.history[sub.UserID] {
		if h.attempt != attempt {
			continue
		}
		old := 0.0
		found := false
		for i := range h.data {
			if h.data[i].QuestionID == 789 {
				old = h.data[i].Points
				h.data[i].Points = newScore
				found = true
			}
		}
		if !found {
			h.data = append(h.data, data(789, `null`, newScore))
		}
		h.score = h.score - old + newScore
		if attempt == sub.Attempt {
			sub.Score = h.score
		}
	}
}

// withFastReadBack removes the propagation delay between read-back retries.
func withFastReadBack(t *testing.T) {
	t.Helper()
	orig := quizRegradeReadBackDelay
	quizRegradeReadBackDelay = 0
	t.Cleanup(func() { quizRegradeReadBackDelay = orig })
}

func (rs *regradeServer) question() api.QuizQuestion {
	return api.QuizQuestion{
		ID: 789, QuizID: 10, QuestionType: rs.questionType, PointsPossible: 4,
		Answers: []api.QuizAnswer{
			{ID: 1001, Text: "old correct", Weight: rs.weights[1001]},
			{ID: 1002, Text: "new correct", Weight: rs.weights[1002]},
			{ID: 1003, Text: "distractor", Weight: rs.weights[1003]},
		},
	}
}

func newRegradeClient(t *testing.T, rs *regradeServer) *api.Client {
	t.Helper()
	client, err := api.NewClient(api.ClientConfig{BaseURL: rs.URL, Token: "test-token", RequestsPerSec: 1000})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

func setRegradePerPage(t *testing.T, n int) {
	t.Helper()
	orig := quizRegradePerPage
	quizRegradePerPage = n
	t.Cleanup(func() { quizRegradePerPage = orig })
}

func regradeOpts(attempts string, dry bool) *options.QuizzesRegradeOptions {
	return &options.QuizzesRegradeOptions{
		CourseID: 1, QuizID: 10, QuestionID: 789, CorrectAnswerID: 1002,
		Attempts: attempts, Force: true, DryRun: dry,
	}
}

func TestRunQuizzesRegrade_EndToEnd(t *testing.T) {
	setRegradePerPage(t, 2)
	withFastReadBack(t)
	rs := newRegradeServer(t)
	defer rs.Close()
	rs.staleReads = 1 // the first read after each PUT still shows the old score

	var runErr error
	out := captureStdout(func() {
		runErr = runQuizzesRegrade(context.Background(), newRegradeClient(t, rs), regradeOpts("completed", false))
	})
	if runErr != nil {
		t.Fatalf("runQuizzesRegrade: %v\n%s", runErr, out)
	}

	rs.mu.Lock()
	defer rs.mu.Unlock()

	// pagination: 8 submissions at 2 per page -> four full pages, then the
	// empty page 5 ends the loop
	if len(rs.listPages) != 5 {
		t.Errorf("expected 5 list pages, got %d: %v", len(rs.listPages), rs.listPages)
	}
	if !strings.Contains(rs.listPages[0], "page=1") || !strings.Contains(rs.listPages[0], "per_page=2") {
		t.Errorf("first list query = %q, want page=1&per_page=2", rs.listPages[0])
	}

	// answers were read from the assignment submission history, once per
	// candidate user (6 complete), all with include[]=submission_history
	if len(rs.historyGETs) != 6 {
		t.Errorf("expected 6 history reads, got %d: %v", len(rs.historyGETs), rs.historyGETs)
	}
	for _, g := range rs.historyGETs {
		if !strings.Contains(g, "include%5B%5D=submission_history") {
			t.Errorf("history read without include[]=submission_history: %s", g)
		}
	}

	// (c) the answer key was written with explicit weights
	if len(rs.questionPUTs) != 1 {
		t.Fatalf("expected exactly one question PUT, got %d", len(rs.questionPUTs))
	}
	q, _ := rs.questionPUTs[0]["question"].(map[string]interface{})
	answers, _ := q["answers"].([]interface{})
	gotWeights := map[float64]interface{}{}
	for _, a := range answers {
		m, _ := a.(map[string]interface{})
		w, present := m["weight"]
		if !present {
			t.Errorf("answer %v was sent without an explicit weight", m["id"])
		}
		gotWeights[m["id"].(float64)] = w
	}
	if gotWeights[1001] != float64(0) || gotWeights[1002] != float64(100) || gotWeights[1003] != float64(0) {
		t.Errorf("weights sent = %v, want 1001:0 1002:100 1003:0", gotWeights)
	}
	if rs.weights[1002] != 100 || rs.weights[1001] != 0 {
		t.Errorf("server-side key after PUT = %v", rs.weights)
	}

	// (e) exactly the affected attempts were written, with the documented body
	if len(rs.subPUTs) != 4 {
		t.Fatalf("expected PUTs for 501, 502, 507, 508 only, got %v", keysOf(rs.subPUTs))
	}
	assertSubmissionPUT(t, rs.subPUTs[501], 1, 4)
	assertSubmissionPUT(t, rs.subPUTs[502], 2, 0) // latest attempt only
	assertSubmissionPUT(t, rs.subPUTs[507], 1, 4) // via the student-side fallback
	assertSubmissionPUT(t, rs.subPUTs[508], 1, 0) // partial credit 2 -> 0

	// (f) scores after the writes
	want := map[int64]float64{501: 10, 502: 5, 503: 3, 505: 7, 506: 2, 507: 8, 508: 4}
	for id, s := range want {
		if rs.subs[id].Score != s {
			t.Errorf("submission %d score = %g, want %g", id, rs.subs[id].Score, s)
		}
	}

	for _, want := range []string{
		"correct answer [1001] -> 1002",
		"Done: 6 attempts considered, 4 changed, 4 verified, 0 mismatched.",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
	// every read-back went through ?attempt=N and needed one retry (stale first read)
	for _, g := range rs.attemptGETs {
		if !strings.Contains(g, "attempt=") {
			t.Errorf("read-back without ?attempt: %s", g)
		}
	}
	if len(rs.attemptGETs) != 8 {
		t.Errorf("expected 8 attempt reads (4 changed x 2), got %d: %v", len(rs.attemptGETs), rs.attemptGETs)
	}
	assertTableRow(t, out, "501", "11", "1", "complete", "1002", "0", "4", "6", "10", "10", "yes", "1")
	assertTableRow(t, out, "502", "12", "2", "complete", "1001", "4", "0", "9", "5", "5", "yes", "1")
	assertTableRow(t, out, "503", "13", "1", "complete", "1003", "0", "0", "3", "3", "-", "-", "-")
	assertTableRow(t, out, "506", "16", "1", "complete", "-", "0", "0", "2", "2", "-", "-", "-")
	assertTableRow(t, out, "507", "17", "1", "complete", "1002", "0", "4", "4", "8", "8", "yes", "1")
	assertTableRow(t, out, "508", "18", "1", "complete", "1001", "2", "0", "6", "4", "4", "yes", "1")
	if strings.Contains(out, "\n504") || strings.Contains(out, "\n505") {
		t.Errorf("untaken / pending_review submissions must not appear by default:\n%s", out)
	}
}

func TestRunQuizzesRegrade_AttemptsAll(t *testing.T) {
	setRegradePerPage(t, 100)
	withFastReadBack(t)
	rs := newRegradeServer(t)
	defer rs.Close()

	var runErr error
	out := captureStdout(func() {
		runErr = runQuizzesRegrade(context.Background(), newRegradeClient(t, rs), regradeOpts("all", false))
	})
	if runErr != nil {
		t.Fatalf("runQuizzesRegrade: %v\n%s", runErr, out)
	}
	rs.mu.Lock()
	defer rs.mu.Unlock()

	// pending_review 505 is included, and user 12's older attempt 1 gets its own PUT
	if len(rs.subPUTs) != 5 || rs.subPUTs[505] == nil || len(rs.subPUTs[502]) != 2 {
		t.Errorf("expected PUTs for 501, 502 (x2), 505, 507, 508, got %v (502: %d)", keysOf(rs.subPUTs), len(rs.subPUTs[502]))
	}
	a1, a2 := rs.subPUTs[502][0], rs.subPUTs[502][1]
	if attemptOf(a1) != 1 || attemptOf(a2) != 2 {
		t.Errorf("502 PUT attempts = %d, %d; want 1 then 2", attemptOf(a1), attemptOf(a2))
	}
	// attempt 1 of user 12: 5 -> 9 in history; latest attempt 2: 9 -> 5
	if h := rs.history[12]; h[0].score != 9 || h[1].score != 5 || rs.subs[502].Score != 5 {
		t.Errorf("user 12 after regrade: attempt1=%g attempt2=%g quiz submission=%g", h[0].score, h[1].score, rs.subs[502].Score)
	}
	if rs.subs[505].Score != 11 {
		t.Errorf("505 score = %g, want 11", rs.subs[505].Score)
	}
	// the older attempt was read back through ?attempt=1 on the quiz submission
	found := false
	for _, g := range rs.attemptGETs {
		if strings.HasSuffix(g, "/submissions/502?attempt=1") {
			found = true
		}
	}
	if !found {
		t.Errorf("attempt 1 of 502 was not read back via ?attempt=1: %v", rs.attemptGETs)
	}
	if !strings.Contains(out, "Done: 8 attempts considered, 6 changed, 6 verified, 0 mismatched.") {
		t.Errorf("unexpected summary:\n%s", out)
	}
	assertTableRow(t, out, "502", "12", "1", "complete", "1002", "0", "4", "5", "9", "9", "yes", "0")
	assertTableRow(t, out, "505", "15", "1", "pending_review", "1002", "0", "4", "7", "11", "11", "yes", "0")
}

func TestRunQuizzesRegrade_DryRunWritesNothing(t *testing.T) {
	setRegradePerPage(t, 100)
	rs := newRegradeServer(t)
	defer rs.Close()

	var runErr error
	out := captureStdout(func() {
		runErr = runQuizzesRegrade(context.Background(), newRegradeClient(t, rs), regradeOpts("completed", true))
	})
	if runErr != nil {
		t.Fatalf("runQuizzesRegrade: %v", runErr)
	}
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if len(rs.questionPUTs) != 0 || len(rs.subPUTs) != 0 {
		t.Errorf("dry run must not write: question PUTs=%d submission PUTs=%v", len(rs.questionPUTs), keysOf(rs.subPUTs))
	}
	if rs.weights[1001] != 100 {
		t.Errorf("dry run changed the server-side key: %v", rs.weights)
	}
	for _, want := range []string{
		"DRY RUN: no changes were made.",
		"Plan: 6 attempts considered, 4 would change, 0 written.",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
	assertTableRow(t, out, "501", "11", "1", "complete", "1002", "0", "4", "6", "10", "-", "-", "-")
	assertTableRow(t, out, "502", "12", "2", "complete", "1001", "4", "0", "9", "5", "-", "-", "-")
	assertTableRow(t, out, "508", "18", "1", "complete", "1001", "2", "0", "6", "4", "-", "-", "-")
}

func TestRunQuizzesRegrade_ReadBackMismatchFails(t *testing.T) {
	setRegradePerPage(t, 100)
	withFastReadBack(t)
	rs := newRegradeServer(t)
	defer rs.Close()
	rs.ignoreScores = true

	var runErr error
	out := captureStdout(func() {
		runErr = runQuizzesRegrade(context.Background(), newRegradeClient(t, rs), regradeOpts("completed", false))
	})
	if runErr == nil {
		t.Fatalf("expected a non-nil error when read-back does not match:\n%s", out)
	}
	if !strings.Contains(runErr.Error(), "4 of 4 regraded attempts did not read back the expected score") {
		t.Errorf("unexpected error: %v", runErr)
	}
	if !strings.Contains(out, "Done: 6 attempts considered, 4 changed, 0 verified, 4 mismatched.") {
		t.Errorf("unexpected summary:\n%s", out)
	}
	assertTableRow(t, out, "501", "11", "1", "complete", "1002", "0", "4", "6", "10", "6", "no", "2")
}

// TestQuizzesRegradeCmd_MismatchReturnsError runs the cobra command itself
// (RunE, as main does) against the fixture server and checks a mismatch
// surfaces as an error — which main turns into a non-zero exit code.
func TestQuizzesRegradeCmd_MismatchReturnsError(t *testing.T) {
	setRegradePerPage(t, 100)
	withFastReadBack(t)
	rs := newRegradeServer(t)
	defer rs.Close()
	rs.ignoreScores = true
	t.Setenv("CANVAS_URL", rs.URL)
	t.Setenv("CANVAS_TOKEN", "test-token")
	t.Setenv("CANVAS_REQUESTS_PER_SEC", "1000")

	cmd := newQuizzesRegradeCmd()
	cmd.SetArgs([]string{"10", "--course-id", "1", "--question", "789", "--correct-answer-id", "1002", "--force"})
	var err error
	out := captureStdout(func() { err = cmd.ExecuteContext(context.Background()) })
	if err == nil {
		t.Fatalf("expected the command to fail on a read-back mismatch:\n%s", out)
	}
	if !strings.Contains(err.Error(), "did not read back the expected score") {
		t.Errorf("unexpected error: %v", err)
	}

	// and the happy path, same wiring, returns nil — with a stale first
	// read after each PUT, which the retry must see through. getAPIClient
	// enables response caching here, so this also pins that the read-back
	// bypasses the cache (a cached stale score would never converge).
	rs2 := newRegradeServer(t)
	defer rs2.Close()
	rs2.staleReads = 1
	t.Setenv("CANVAS_URL", rs2.URL)
	cmd = newQuizzesRegradeCmd()
	cmd.SetArgs([]string{"10", "--course-id", "1", "--question", "789", "--correct-answer-id", "1002", "--force"})
	out = captureStdout(func() { err = cmd.ExecuteContext(context.Background()) })
	if err != nil {
		t.Fatalf("unexpected error on the happy path: %v\n%s", err, out)
	}
}

func TestRunQuizzesRegrade_RefusesUnsupportedType(t *testing.T) {
	setRegradePerPage(t, 100)
	rs := newRegradeServer(t)
	defer rs.Close()
	rs.questionType = "essay_question"

	err := runQuizzesRegrade(context.Background(), newRegradeClient(t, rs), regradeOpts("completed", false))
	if err == nil || !strings.Contains(err.Error(), "essay_question") {
		t.Fatalf("expected an unsupported-type error, got %v", err)
	}
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if len(rs.questionPUTs) != 0 || len(rs.subPUTs) != 0 {
		t.Errorf("nothing may be written for an unsupported type")
	}
}

func TestRunQuizzesRegrade_RefusesUnknownAnswer(t *testing.T) {
	setRegradePerPage(t, 100)
	rs := newRegradeServer(t)
	defer rs.Close()

	opts := regradeOpts("completed", false)
	opts.CorrectAnswerID = 4242
	err := runQuizzesRegrade(context.Background(), newRegradeClient(t, rs), opts)
	if err == nil || !strings.Contains(err.Error(), "4242 is not an answer of question 789") {
		t.Fatalf("expected an unknown-answer error, got %v", err)
	}
}

func TestRunQuizzesRegrade_JSONOutput(t *testing.T) {
	setRegradePerPage(t, 100)
	withFastReadBack(t)
	setOutputFormat(t, "json")
	rs := newRegradeServer(t)
	defer rs.Close()

	var runErr error
	out := captureStdout(func() {
		runErr = runQuizzesRegrade(context.Background(), newRegradeClient(t, rs), regradeOpts("completed", false))
	})
	if runErr != nil {
		t.Fatalf("runQuizzesRegrade: %v", runErr)
	}
	var result quizRegradeResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, out)
	}
	if result.AssignmentID != 77 || result.Summary.Changed != 4 || result.Summary.Verified != 4 || result.Summary.Mismatched != 0 {
		t.Errorf("result = %+v", result)
	}
	if len(result.Submissions) != 6 || result.Submissions[0].SubmissionID != 501 || result.Submissions[0].NewScore != 10 || result.Submissions[0].AnswerSource != regradeSourceHistory {
		t.Errorf("submissions = %+v", result.Submissions)
	}
}

// TestQuizzesRegradeCmd exercises the cobra wiring: flag validation, and a
// no-op regrade (the mock key already has 1002 correct) under --dry-run,
// which must still perform real reads.
func TestQuizzesRegradeCmd(t *testing.T) {
	withGlobalDryRun(t, true)
	mocks := map[string]cmdtest.MockResponse{
		"/api/v1/courses/1/quizzes/10": cmdtest.NewMockResponse(`{"id": 10, "assignment_id": 77}`),
		"/api/v1/courses/1/quizzes/10/questions/789": cmdtest.NewMockResponse(`{
			"id": 789, "quiz_id": 10, "question_type": "true_false_question", "points_possible": 1,
			"answers": [{"id": 1001, "text": "True", "weight": 0}, {"id": 1002, "text": "False", "weight": 100}]
		}`),
		"/api/v1/courses/1/quizzes/10/submissions": cmdtest.NewMockResponse(`{"quiz_submissions": [
			{"id": 501, "quiz_id": 10, "user_id": 11, "attempt": 1, "workflow_state": "complete", "score": 1}
		]}`),
		"/api/v1/courses/1/assignments/77/submissions/11": cmdtest.NewMockResponse(`{
			"id": 9011, "user_id": 11, "attempt": 1, "score": 1,
			"submission_history": [{"attempt": 1, "score": 1, "submission_data": [
				{"question_id": 789, "answer_id": 1002, "correct": true, "points": 1}
			]}]
		}`),
	}
	tests := []cmdtest.CommandTestCase{
		{
			Name:          "dry run reads and prints the plan",
			Args:          []string{"10", "--course-id", "1", "--question", "789", "--correct-answer-id", "1002"},
			MockResponses: mocks,
			ExpectError:   false,
			ValidateOutput: func(t *testing.T, output string) {
				for _, want := range []string{
					"DRY RUN: no changes were made.",
					"Question 789 (true_false_question, 1 points): correct answer [1002] -> 1002",
					"Plan: 1 attempts considered, 0 would change, 0 written.",
				} {
					if !strings.Contains(output, want) {
						t.Errorf("expected %q in output:\n%s", want, output)
					}
				}
				assertTableRow(t, output, "501", "11", "1", "complete", "1002", "1", "1", "1", "1", "-", "-", "-")
				if strings.Contains(output, "curl") {
					t.Errorf("reads must be real under --dry-run, not curl echoes:\n%s", output)
				}
			},
		},
		{
			Name:        "missing quiz ID",
			Args:        []string{"--course-id", "1", "--question", "789", "--correct-answer-id", "1002"},
			ExpectError: true,
		},
		{
			Name:        "invalid quiz ID",
			Args:        []string{"abc", "--course-id", "1", "--question", "789", "--correct-answer-id", "1002"},
			ExpectError: true,
		},
		{
			Name:        "missing question",
			Args:        []string{"10", "--course-id", "1", "--correct-answer-id", "1002"},
			ExpectError: true,
		},
		{
			Name:        "missing correct answer",
			Args:        []string{"10", "--course-id", "1", "--question", "789"},
			ExpectError: true,
		},
		{
			Name:        "bad attempts value",
			Args:        []string{"10", "--course-id", "1", "--question", "789", "--correct-answer-id", "1002", "--attempts", "some"},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newQuizzesRegradeCmd(), tc)
		})
	}
}

// ---- helpers ----

func keysOf(m map[int64][]map[string]interface{}) []int64 {
	var out []int64
	for k := range m {
		out = append(out, k)
	}
	return out
}

func attemptOf(body map[string]interface{}) int {
	entries, _ := body["quiz_submissions"].([]interface{})
	if len(entries) == 0 {
		return -1
	}
	entry, _ := entries[0].(map[string]interface{})
	a, _ := entry["attempt"].(float64)
	return int(a)
}

func assertSubmissionPUT(t *testing.T, puts []map[string]interface{}, wantAttempt int, wantScore float64) {
	t.Helper()
	if len(puts) != 1 {
		t.Fatalf("expected exactly one PUT, got %d", len(puts))
	}
	entries, _ := puts[0]["quiz_submissions"].([]interface{})
	if len(entries) != 1 {
		t.Fatalf("quiz_submissions = %v", puts[0]["quiz_submissions"])
	}
	entry, _ := entries[0].(map[string]interface{})
	if entry["attempt"] != float64(wantAttempt) {
		t.Errorf("attempt = %v, want %d", entry["attempt"], wantAttempt)
	}
	questions, _ := entry["questions"].(map[string]interface{})
	q789, _ := questions["789"].(map[string]interface{})
	if q789["score"] != wantScore {
		t.Errorf("questions[789].score = %v, want %g (body: %v)", q789["score"], wantScore, entry)
	}
	if len(entry) != 2 {
		t.Errorf("body should carry only attempt and questions, got %v", entry)
	}
}

// assertTableRow finds the table line for (submission id, attempt) and
// checks its whitespace-separated columns.
func assertTableRow(t *testing.T, out string, cols ...string) {
	t.Helper()
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 || fields[0] != cols[0] || fields[2] != cols[2] {
			continue
		}
		if strings.Join(fields, " ") != strings.Join(cols, " ") {
			t.Errorf("row %s/%s = %q, want %q", cols[0], cols[2], strings.Join(fields, " "), strings.Join(cols, " "))
		}
		return
	}
	t.Errorf("no table row for submission %s attempt %s in:\n%s", cols[0], cols[2], out)
}

// TestRunQuizzesRegrade_ActivityLogEntry runs a regrade with the activity
// log enabled and checks the entry: the verification table and summary
// under details.regrade, and outcome "ok" on every PUT.
func TestRunQuizzesRegrade_ActivityLogEntry(t *testing.T) {
	useTempHome(t)
	logPath := useTempActivityLog(t)
	t.Setenv(activity.EnvCaptureBodies, "true") // details.regrade rides with the bodies
	withFastReadBack(t)
	rs := newRegradeServer(t)
	defer rs.Close()

	rec := activity.Default()
	rec.Reset()
	api.RequestObserver = func(o api.ObservedRequest) {
		rec.Observe(activity.Observation{Method: o.Method, Path: o.Path, Status: o.Status, DryRun: o.DryRun, RequestBody: o.RequestBody, ResponseBody: o.ResponseBody})
	}
	t.Cleanup(func() { api.RequestObserver = nil; rec.Reset() })

	var runErr error
	out := captureStdout(func() {
		runErr = runQuizzesRegrade(context.Background(), newRegradeClient(t, rs), regradeOpts("completed", false))
	})
	if runErr != nil {
		t.Fatalf("runQuizzesRegrade: %v\n%s", runErr, out)
	}
	argv := []string{"quizzes", "regrade", "10", "--course-id", "1", "--question", "789", "--correct-answer-id", "1002"}
	logActivity(fakeExecuted(t, argv[2:]...), nil, time.Now(), rec, argv)

	entries, skipped, err := activity.Read(logPath)
	if err != nil || skipped != 0 || len(entries) != 1 {
		t.Fatalf("activity log: %d entries, %d skipped, %v", len(entries), skipped, err)
	}
	e := entries[0]
	if e.Command != "quizzes regrade" || e.ExitCode != 0 || e.VerificationRequired {
		t.Errorf("entry = %+v", e)
	}

	puts := 0
	for _, r := range e.Requests {
		if r.Method == http.MethodPut {
			puts++
			if r.Outcome != activity.OutcomeOK || r.Status != 200 {
				t.Errorf("PUT %s: status %d outcome %q, want 200 ok", r.Path, r.Status, r.Outcome)
			}
		}
	}
	// one question PUT plus one per changed attempt, as the server saw them
	rs.mu.Lock()
	want := 1
	for _, p := range rs.subPUTs {
		want += len(p)
	}
	rs.mu.Unlock()
	if puts != want {
		t.Errorf("PUTs logged = %d, want %d", puts, want)
	}

	raw, _ := json.Marshal(e.Details["regrade"])
	var detail struct {
		QuestionID      int64 `json:"question_id"`
		CorrectAnswerID int64 `json:"correct_answer_id"`
		Summary         struct {
			Considered, Changed, Verified, Mismatched int
		} `json:"summary"`
		Submissions []struct {
			SubmissionID  int64   `json:"submission_id"`
			Attempt       int     `json:"attempt"`
			OldScore      float64 `json:"old_score"`
			ExpectedScore float64 `json:"expected_score"`
			NewScore      float64 `json:"new_score"`
			Verified      string  `json:"verified"`
		} `json:"submissions"`
	}
	if err := json.Unmarshal(raw, &detail); err != nil || detail.QuestionID != 789 || detail.CorrectAnswerID != 1002 {
		t.Fatalf("details.regrade = %s (%v)", raw, err)
	}
	if detail.Summary.Changed == 0 || detail.Summary.Verified != detail.Summary.Changed || detail.Summary.Mismatched != 0 {
		t.Errorf("summary = %+v", detail.Summary)
	}
	verified := 0
	for _, row := range detail.Submissions {
		if row.Verified == "yes" {
			verified++
			if row.NewScore != row.ExpectedScore {
				t.Errorf("row %+v: new_score must equal expected_score when verified", row)
			}
		}
	}
	if verified != detail.Summary.Verified {
		t.Errorf("verification table has %d verified rows, summary says %d", verified, detail.Summary.Verified)
	}
}

// A Canvas admin can cap per_page below what was asked. The first page then
// comes back short without being the last one; the walk must continue to
// the empty page instead of rescoring only the first page and reporting
// success.
func TestRunQuizzesRegrade_ServerCapsPageSize(t *testing.T) {
	setRegradePerPage(t, 100)
	withFastReadBack(t)
	rs := newRegradeServer(t)
	defer rs.Close()
	rs.capPerPage = 2

	var runErr error
	out := captureStdout(func() {
		runErr = runQuizzesRegrade(context.Background(), newRegradeClient(t, rs), regradeOpts("completed", false))
	})
	if runErr != nil {
		t.Fatalf("runQuizzesRegrade: %v\n%s", runErr, out)
	}
	rs.mu.Lock()
	defer rs.mu.Unlock()
	// 8 submissions at a capped 2 per page: four pages, then the empty fifth
	if len(rs.listPages) != 5 {
		t.Errorf("expected 5 list pages, got %d: %v", len(rs.listPages), rs.listPages)
	}
	// every complete submission was considered, not just the first page's
	if len(rs.historyGETs) != 6 {
		t.Errorf("expected all 6 complete submissions to be read, got %d: %v", len(rs.historyGETs), rs.historyGETs)
	}
}
