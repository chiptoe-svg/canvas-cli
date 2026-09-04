package commands

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/activity"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// Canvas encodes correctness of a multiple-choice / true-false answer as its
// weight: 100 for the correct answer, 0 for every other one
// (https://canvas.instructure.com/doc/api/quiz_questions.html, Answer object:
// "answer_weight ... Incorrect answers should be 0, correct answers should be 100").
const (
	quizAnswerWeightCorrect = 100
	quizAnswerWeightWrong   = 0
)

// Where a submission's answers were read from (quizRegradeRow.AnswerSource).
const (
	regradeSourceHistory       = "submission_history"        // assignment submission, include[]=submission_history
	regradeSourceQuizQuestions = "quiz_submission_questions" // GET /quiz_submissions/:id/questions (student-side fallback)
	regradeSourceNone          = "none"                      // no answer record found for the question
)

// quizRegradePerPage is the page size used when listing quiz submissions.
// A variable so tests can force multi-page listings with small fixtures.
var quizRegradePerPage = 100

// quizRegradeScoreTolerance bounds float noise when comparing Canvas scores.
const quizRegradeScoreTolerance = 1e-6

// Canvas applies a score update asynchronously enough that an immediate
// read-back can still return the previous value (observed live on an older
// attempt). A mismatching read is therefore retried a few times before it
// is reported. Variables so tests can shorten the delay.
var (
	quizRegradeReadBackAttempts = 3
	quizRegradeReadBackDelay    = time.Second
)

func init() {
	quizzesCmd.AddCommand(newQuizzesRegradeCmd())
}

func newQuizzesRegradeCmd() *cobra.Command {
	opts := &options.QuizzesRegradeOptions{}

	cmd := &cobra.Command{
		Use:   "regrade <quiz-id>",
		Short: "Change a question's correct answer and rescore existing submissions",
		Long: `Regrade one multiple-choice or true/false question of a classic quiz.

The command:
  1. fetches the question and rewrites its answers so --correct-answer-id
     has weight 100 and every other answer weight 0, then saves it;
  2. lists the quiz submissions (workflow_state=complete by default,
     --attempts all to include pending_review too);
  3. reads each student's answers from the assignment submission's
     submission_history (the per-attempt record graders can see), works
     out the new per-question score (points_possible if the student picked
     the new correct answer, else 0) and the resulting change to that
     attempt's score. By default only the latest attempt is rescored;
     --attempts all rescores every attempt in the history;
  4. writes the new question score to each affected attempt;
  5. reads every affected attempt back and prints a before/after table
     with a verified column. The command exits non-zero if any read-back
     does not match the expected score.

Only multiple_choice_question and true_false_question are supported. Other
question types are refused before anything is written.

--dry-run prints the full plan (the answer-key change and every attempt
that would change, with old and expected scores) and writes nothing.

Examples:
  canvas quizzes regrade 456 --course-id 123 --question 789 --correct-answer-id 1002 --dry-run
  canvas quizzes regrade 456 --course-id 123 --question 789 --correct-answer-id 1002 --force
  canvas quizzes regrade 456 --course-id 123 --question 789 --correct-answer-id 1002 --attempts all`,
		Args: ExactArgsWithUsage(1, "quiz-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			quizID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid quiz ID: %s", args[0])
			}
			opts.QuizID = quizID
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
			// Every read here must be live: the plan is derived from current
			// scores and the read-back retries re-fetch the same URL until
			// Canvas has applied the update. A cached GET would defeat both.
			client.SetCacheEnabled(false)

			return runQuizzesRegrade(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuestionID, "question", 0, "Quiz question ID to regrade (required)")
	cmd.Flags().Int64Var(&opts.CorrectAnswerID, "correct-answer-id", 0, "Answer ID that becomes the only correct answer (required)")
	cmd.Flags().StringVar(&opts.Attempts, "attempts", "completed", "Which attempts to rescore: completed (latest attempt of complete submissions) or all (every attempt, pending_review included)")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")
	mustMarkRequired(cmd, "course-id", "question", "correct-answer-id")

	return cmd
}

// quizRegradeRow is one attempt's before/after line.
type quizRegradeRow struct {
	SubmissionID     int64   `json:"submission_id"` // quiz submission ID
	UserID           int64   `json:"user_id"`
	Attempt          int     `json:"attempt"`
	WorkflowState    string  `json:"workflow_state"`
	AnswerSource     string  `json:"answer_source"`
	SelectedAnswerID int64   `json:"selected_answer_id"` // 0 when unanswered
	OldQuestionScore float64 `json:"old_question_score"`
	NewQuestionScore float64 `json:"new_question_score"`
	OldScore         float64 `json:"old_score"`
	ExpectedScore    float64 `json:"expected_score"`
	NewScore         float64 `json:"new_score"`
	Changed          bool    `json:"changed"`
	Verified         string  `json:"verified"`          // yes | no | - (not written)
	ReadBackRetries  int     `json:"read_back_retries"` // extra reads needed before the score matched (or gave up)
}

// quizRegradeSummary counts the outcome.
type quizRegradeSummary struct {
	Considered int `json:"considered"`
	Changed    int `json:"changed"`
	Verified   int `json:"verified"`
	Mismatched int `json:"mismatched"`
}

// quizRegradeResult is the structured output of the command.
type quizRegradeResult struct {
	CourseID        int64              `json:"course_id"`
	QuizID          int64              `json:"quiz_id"`
	AssignmentID    int64              `json:"assignment_id"`
	QuestionID      int64              `json:"question_id"`
	QuestionType    string             `json:"question_type"`
	PointsPossible  float64            `json:"points_possible"`
	OldCorrectIDs   []int64            `json:"old_correct_answer_ids"`
	CorrectAnswerID int64              `json:"correct_answer_id"`
	DryRun          bool               `json:"dry_run"`
	Submissions     []quizRegradeRow   `json:"submissions"`
	Summary         quizRegradeSummary `json:"summary"`
}

func runQuizzesRegrade(ctx context.Context, client *api.Client, opts *options.QuizzesRegradeOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.regrade", map[string]interface{}{
		"course_id":         opts.CourseID,
		"quiz_id":           opts.QuizID,
		"question_id":       opts.QuestionID,
		"correct_answer_id": opts.CorrectAnswerID,
		"attempts":          opts.Attempts,
		"dry_run":           opts.DryRun,
	})

	quizzes := api.NewQuizzesService(client)
	questions := api.NewQuizQuestionsService(client)
	quizSubmissions := api.NewQuizSubmissionsService(client)
	submissionQuestions := api.NewQuizSubmissionQuestionsService(client)
	assignmentSubmissions := api.NewSubmissionsService(client)

	// (a) the question as it is now
	question, err := questions.Get(ctx, opts.CourseID, opts.QuizID, opts.QuestionID)
	if err != nil {
		return fmt.Errorf("failed to get question %d: %w", opts.QuestionID, err)
	}

	// (b) the new answer key
	newAnswers, err := regradeAnswers(question, opts.CorrectAnswerID)
	if err != nil {
		return err
	}
	oldCorrect := correctAnswerIDs(question.Answers)

	// The per-attempt answer record lives on the assignment submission, so
	// the quiz's assignment is needed to reach it.
	quiz, err := quizzes.Get(ctx, opts.CourseID, opts.QuizID)
	if err != nil {
		return fmt.Errorf("failed to get quiz %d: %w", opts.QuizID, err)
	}
	if quiz.AssignmentID == 0 {
		return fmt.Errorf("quiz %d has no assignment_id; only graded (assignment) quizzes can be regraded", opts.QuizID)
	}

	// (d) candidate submissions
	all, err := listAllQuizSubmissions(ctx, quizSubmissions, opts.CourseID, opts.QuizID)
	if err != nil {
		return fmt.Errorf("failed to list quiz submissions: %w", err)
	}
	candidates := selectRegradeSubmissions(all, opts.Attempts)

	// (e) per-attempt plan
	result := &quizRegradeResult{
		CourseID:        opts.CourseID,
		QuizID:          opts.QuizID,
		AssignmentID:    quiz.AssignmentID,
		QuestionID:      question.ID,
		QuestionType:    question.QuestionType,
		PointsPossible:  question.PointsPossible,
		OldCorrectIDs:   oldCorrect,
		CorrectAnswerID: opts.CorrectAnswerID,
		DryRun:          opts.DryRun,
	}
	for _, sub := range candidates {
		asub, err := assignmentSubmissions.Get(ctx, opts.CourseID, quiz.AssignmentID, sub.UserID, []string{"submission_history"})
		if err != nil {
			return fmt.Errorf("failed to read submission history of user %d: %w", sub.UserID, err)
		}
		attempts := selectRegradeAttempts(sub, asub, opts.Attempts)
		if len(attempts) == 0 {
			// No graded attempt record on the assignment submission. Fall
			// back to the quiz-submission questions endpoint, which carries
			// the selected answer only when the caller can see it
			// (student-side); graders get {id, correct, flagged} there.
			answered, err := submissionQuestions.List(ctx, sub.ID, nil)
			if err != nil {
				return fmt.Errorf("failed to read answers of submission %d: %w", sub.ID, err)
			}
			attempts = []regradeAttempt{fallbackAttemptFromQuizQuestions(sub, question, oldCorrect, answered)}
		}
		for _, a := range attempts {
			result.Submissions = append(result.Submissions, planRegradeRow(sub, a, question, opts.CorrectAnswerID))
		}
	}
	result.Summary.Considered = len(result.Submissions)
	for _, row := range result.Submissions {
		if row.Changed {
			result.Summary.Changed++
		}
	}

	if opts.DryRun {
		return printRegradeResult(result)
	}

	if !opts.Force {
		ok, err := confirmRegrade(result)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// (c) save the answer key first; if this fails nothing else is touched
	updatedQuestion, err := questions.Update(ctx, opts.CourseID, opts.QuizID, opts.QuestionID, &api.UpdateQuizQuestionParams{
		Answers: &newAnswers,
	})
	if err != nil {
		return fmt.Errorf("failed to update question %d: %w", opts.QuestionID, err)
	}
	if got := correctAnswerIDs(updatedQuestion.Answers); len(got) != 1 || got[0] != opts.CorrectAnswerID {
		return fmt.Errorf("question %d saved, but Canvas reports correct answer(s) %v instead of %d", opts.QuestionID, got, opts.CorrectAnswerID)
	}

	// (e) write the new question score to each affected attempt
	for i := range result.Submissions {
		row := &result.Submissions[i]
		if !row.Changed {
			continue
		}
		attempt := row.Attempt
		score := row.NewQuestionScore
		_, err := quizSubmissions.Update(ctx, opts.CourseID, opts.QuizID, row.SubmissionID, &api.UpdateQuizSubmissionParams{
			Attempt: &attempt,
			Questions: map[int64]api.QuizSubmissionQuestionScore{
				question.ID: {Score: &score},
			},
		})
		if err != nil {
			row.Verified = "no"
			result.Summary.Mismatched++
			logger.LogCommandError(ctx, "quizzes.regrade", err, map[string]interface{}{
				"submission_id": row.SubmissionID,
				"attempt":       row.Attempt,
			})
			printInfo("submission %d attempt %d: update failed: %v\n", row.SubmissionID, row.Attempt, err)
			continue
		}

		// (f) read back and verify
		newScore, retries, err := readBackAttemptScore(ctx, quizSubmissions, opts.CourseID, opts.QuizID, row.SubmissionID, row.Attempt, row.ExpectedScore)
		row.ReadBackRetries = retries
		if err != nil {
			row.Verified = "no"
			result.Summary.Mismatched++
			printInfo("submission %d attempt %d: read-back failed: %v\n", row.SubmissionID, row.Attempt, err)
			continue
		}
		row.NewScore = newScore
		if scoresEqual(newScore, row.ExpectedScore) {
			row.Verified = "yes"
			result.Summary.Verified++
		} else {
			row.Verified = "no"
			result.Summary.Mismatched++
		}
	}

	// The verification table and summary go into the activity log entry.
	activity.Default().SetDetail("regrade", result)

	if err := printRegradeResult(result); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "quizzes.regrade", result.Summary.Changed)

	if result.Summary.Mismatched > 0 {
		return fmt.Errorf("%d of %d regraded attempts did not read back the expected score", result.Summary.Mismatched, result.Summary.Changed)
	}
	return nil
}

// regradeAnswers returns a copy of the question's answers with correctID at
// full weight and every other answer at zero. Only multiple-choice and
// true/false questions are supported: for those Canvas treats weight 100 as
// "the correct answer" and 0 as wrong, and the student's selection is a
// single answer ID.
func regradeAnswers(question *api.QuizQuestion, correctID int64) ([]api.QuizAnswer, error) {
	switch question.QuestionType {
	case "multiple_choice_question", "true_false_question":
	default:
		return nil, fmt.Errorf("regrade supports multiple_choice_question and true_false_question; question %d is %q", question.ID, question.QuestionType)
	}

	found := false
	answers := make([]api.QuizAnswer, len(question.Answers))
	for i, a := range question.Answers {
		answers[i] = a
		if a.ID == correctID {
			found = true
			answers[i].Weight = quizAnswerWeightCorrect
		} else {
			answers[i].Weight = quizAnswerWeightWrong
		}
	}
	if !found {
		ids := make([]string, 0, len(question.Answers))
		for _, a := range question.Answers {
			ids = append(ids, strconv.FormatInt(a.ID, 10))
		}
		return nil, fmt.Errorf("answer %d is not an answer of question %d (available: %s)", correctID, question.ID, strings.Join(ids, ", "))
	}
	return answers, nil
}

// correctAnswerIDs lists the answers Canvas currently treats as correct.
// A missing weight decodes as 0, i.e. wrong.
func correctAnswerIDs(answers []api.QuizAnswer) []int64 {
	var ids []int64
	for _, a := range answers {
		if a.Weight == quizAnswerWeightCorrect {
			ids = append(ids, a.ID)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

// listAllQuizSubmissions pages through GET .../quizzes/:quiz_id/submissions.
// The endpoint answers with a {"quiz_submissions": [...]} envelope, so the
// generic array paginator does not apply. It pages until an empty page: a
// page shorter than what was asked for is not the end, because a Canvas
// admin can cap per_page below the request, and stopping there would
// silently rescore only the first page. A page that adds no new submission
// also ends the walk, so a server that ignores the page parameter cannot
// loop forever.
func listAllQuizSubmissions(ctx context.Context, service *api.QuizSubmissionsService, courseID, quizID int64) ([]api.QuizSubmission, error) {
	var all []api.QuizSubmission
	seen := map[int64]bool{}
	for page := 1; ; page++ {
		batch, err := service.List(ctx, courseID, quizID, &api.ListQuizSubmissionsOptions{
			Page:    page,
			PerPage: quizRegradePerPage,
		})
		if err != nil {
			return nil, err
		}
		added := 0
		for _, s := range batch {
			if seen[s.ID] {
				continue
			}
			seen[s.ID] = true
			all = append(all, s)
			added++
		}
		if len(batch) == 0 || added == 0 {
			return all, nil
		}
	}
}

// selectRegradeSubmissions keeps the quiz submissions Canvas can rescore.
// Canvas has one quiz submission per user whose attempt is the latest one.
// "completed" keeps workflow_state=complete. "all" also keeps pending_review;
// untaken, settings_only and preview submissions have no scorable attempt
// (Canvas answers "invalid attempt") and are always skipped.
func selectRegradeSubmissions(subs []api.QuizSubmission, attempts string) []api.QuizSubmission {
	var out []api.QuizSubmission
	for _, s := range subs {
		switch s.WorkflowState {
		case "complete":
			out = append(out, s)
		case "pending_review":
			if attempts == "all" {
				out = append(out, s)
			}
		}
	}
	return out
}

// regradeAttempt is one attempt's answer record, whichever source it came from.
type regradeAttempt struct {
	Attempt  int
	OldScore float64 // the attempt's score before the regrade
	Source   string
	// Data is the grader-visible per-question record (submission_history).
	Data []api.QuizSubmissionData
	// Fallback values used when Data is unavailable (quiz_submission_questions).
	SelectedAnswerID int64
	OldQuestionScore float64
}

// selectRegradeAttempts picks the attempts to rescore from the assignment
// submission's history. Graders see a student's classic-quiz answers there:
// GET /api/v1/courses/:course_id/assignments/:assignment_id/submissions/:user_id?include[]=submission_history
// returns one entry per attempt with "attempt", "score" and "submission_data"
// ({question_id, answer_id, correct, points} per question). By default only
// the entry for the quiz submission's current (latest) attempt is used;
// "all" uses every attempt that has a submission_data record. Entries
// without submission_data (never graded, or not a quiz attempt) are skipped.
func selectRegradeAttempts(sub api.QuizSubmission, asub *api.Submission, attempts string) []regradeAttempt {
	entries := asub.SubmissionHistory
	if len(entries) == 0 && len(asub.SubmissionData) > 0 {
		// Some responses carry the current attempt's data on the top-level
		// object only.
		entries = []api.Submission{*asub}
	}

	byAttempt := map[int]api.Submission{}
	for _, e := range entries {
		if e.Attempt <= 0 || len(e.SubmissionData) == 0 {
			continue
		}
		byAttempt[e.Attempt] = e // later entries win for a repeated attempt
	}

	var out []regradeAttempt
	if attempts == "all" {
		for n, e := range byAttempt {
			out = append(out, regradeAttempt{Attempt: n, OldScore: e.Score, Source: regradeSourceHistory, Data: e.SubmissionData})
		}
		sort.Slice(out, func(i, j int) bool { return out[i].Attempt < out[j].Attempt })
		return out
	}
	if e, ok := byAttempt[sub.Attempt]; ok {
		out = append(out, regradeAttempt{Attempt: sub.Attempt, OldScore: e.Score, Source: regradeSourceHistory, Data: e.SubmissionData})
	}
	return out
}

// fallbackAttemptFromQuizQuestions builds the latest attempt's record from
// GET /api/v1/quiz_submissions/:id/questions
// (https://canvas.instructure.com/doc/api/quiz_submission_questions.html):
// each entry's "id" is the quiz question ID and, for multiple_choice /
// true_false questions, "answer" is the selected answer ID when the caller
// is allowed to see it. Without a points record the old question score is
// derived from the old answer key.
func fallbackAttemptFromQuizQuestions(sub api.QuizSubmission, question *api.QuizQuestion, oldCorrect []int64, answered []api.QuizSubmissionQuestion) regradeAttempt {
	a := regradeAttempt{Attempt: sub.Attempt, OldScore: sub.Score, Source: regradeSourceNone}
	for _, q := range answered {
		if q.ID != question.ID {
			continue
		}
		if id, ok := selectedAnswerID(q.Answer); ok {
			a.Source = regradeSourceQuizQuestions
			a.SelectedAnswerID = id
			for _, c := range oldCorrect {
				if c == id {
					a.OldQuestionScore = question.PointsPossible
				}
			}
		}
		break
	}
	return a
}

// planRegradeRow computes what the regrade does to one attempt.
//
// With a submission_history record, the old per-question score is the
// "points" Canvas awarded for that question (authoritative: it reflects any
// manual adjustment), and the selected answer is "answer_id". The new score
// is points_possible when that answer is the new correct one and 0
// otherwise; the attempt's expected score is old_score + (new - old). Canvas
// recomputes the attempt score from per-question scores when
// quiz_submissions[][questions][<qid>][score] is written, so a read-back
// that differs from expected_score is reported as a mismatch.
func planRegradeRow(sub api.QuizSubmission, a regradeAttempt, question *api.QuizQuestion, newCorrect int64) quizRegradeRow {
	row := quizRegradeRow{
		SubmissionID:     sub.ID,
		UserID:           sub.UserID,
		Attempt:          a.Attempt,
		WorkflowState:    sub.WorkflowState,
		AnswerSource:     a.Source,
		SelectedAnswerID: a.SelectedAnswerID,
		OldQuestionScore: a.OldQuestionScore,
		OldScore:         a.OldScore,
		Verified:         "-",
	}

	// Only answer_id and points are trusted from submission_data. Its
	// "correct" flag is NOT recomputed by Canvas when the answer key changes
	// (observed live: after a regrade an attempt still reported correct:true
	// with points 0), so correctness is always answer_id against the key.
	for _, d := range a.Data {
		if d.QuestionID != question.ID {
			continue
		}
		row.OldQuestionScore = d.Points
		if id, ok := selectedAnswerID(d.AnswerID); ok {
			row.SelectedAnswerID = id
		}
		break
	}

	if row.SelectedAnswerID != 0 && row.SelectedAnswerID == newCorrect {
		row.NewQuestionScore = question.PointsPossible
	}

	delta := row.NewQuestionScore - row.OldQuestionScore
	row.ExpectedScore = row.OldScore + delta
	row.NewScore = row.OldScore
	row.Changed = !scoresEqual(delta, 0)
	return row
}

// readBackAttemptScore re-reads an attempt's score after the update through
// GET .../quizzes/:quiz_id/submissions/:id?attempt=N, which reports that
// attempt's score for every attempt, latest included (the assignment
// submission's history entry lags behind the update and must not be used).
// A read that does not match expected is retried; the number of extra reads
// is returned so the row can show that propagation was needed.
func readBackAttemptScore(ctx context.Context, quizSubmissions *api.QuizSubmissionsService, courseID, quizID, submissionID int64, attempt int, expected float64) (score float64, retries int, err error) {
	for i := 0; i < quizRegradeReadBackAttempts; i++ {
		if i > 0 {
			retries++
			select {
			case <-ctx.Done():
				return score, retries, ctx.Err()
			case <-time.After(quizRegradeReadBackDelay):
			}
		}
		after, getErr := quizSubmissions.GetAttempt(ctx, courseID, quizID, submissionID, attempt)
		if getErr != nil {
			err = getErr
			continue
		}
		err = nil
		score = after.Score
		if scoresEqual(score, expected) {
			return score, retries, nil
		}
	}
	return score, retries, err
}

// selectedAnswerID decodes a selected answer ID as Canvas returns it: a JSON
// number or a numeric string; null, absent or empty means unanswered.
func selectedAnswerID(raw json.RawMessage) (int64, bool) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return 0, false
	}

	var asNumber json.Number
	if err := json.Unmarshal(raw, &asNumber); err == nil {
		if id, err := asNumber.Int64(); err == nil && id > 0 {
			return id, true
		}
		if f, err := asNumber.Float64(); err == nil && f > 0 && f == math.Trunc(f) {
			return int64(f), true
		}
		return 0, false
	}

	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		id, err := strconv.ParseInt(strings.TrimSpace(asString), 10, 64)
		if err == nil && id > 0 {
			return id, true
		}
	}
	return 0, false
}

func scoresEqual(a, b float64) bool {
	return math.Abs(a-b) < quizRegradeScoreTolerance
}

func confirmRegrade(result *quizRegradeResult) (bool, error) {
	fmt.Printf("About to set answer %d as the only correct answer of question %d and rescore %d of %d attempts.\n",
		result.CorrectAnswerID, result.QuestionID, result.Summary.Changed, result.Summary.Considered)
	fmt.Print("Continue? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

func printRegradeResult(result *quizRegradeResult) error {
	return formatOutput(result, func() {
		if result.DryRun {
			fmt.Println("DRY RUN: no changes were made.")
		}
		fmt.Printf("Question %d (%s, %g points): correct answer %v -> %d\n",
			result.QuestionID, result.QuestionType, result.PointsPossible, result.OldCorrectIDs, result.CorrectAnswerID)
		fmt.Println()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SUBMISSION\tUSER\tATTEMPT\tSTATE\tANSWER\tQ OLD\tQ NEW\tOLD SCORE\tEXPECTED\tNEW SCORE\tVERIFIED\tRETRIES")
		for _, r := range result.Submissions {
			answer := "-"
			if r.SelectedAnswerID != 0 {
				answer = strconv.FormatInt(r.SelectedAnswerID, 10)
			}
			newScore := "-"
			if r.Verified != "-" {
				newScore = fmt.Sprintf("%g", r.NewScore)
			}
			retries := "-"
			if r.Verified != "-" {
				retries = strconv.Itoa(r.ReadBackRetries)
			}
			fmt.Fprintf(w, "%d\t%d\t%d\t%s\t%s\t%g\t%g\t%g\t%g\t%s\t%s\t%s\n",
				r.SubmissionID, r.UserID, r.Attempt, r.WorkflowState, answer,
				r.OldQuestionScore, r.NewQuestionScore, r.OldScore, r.ExpectedScore, newScore, r.Verified, retries)
		}
		w.Flush()

		fmt.Println()
		if result.DryRun {
			fmt.Printf("Plan: %d attempts considered, %d would change, 0 written.\n",
				result.Summary.Considered, result.Summary.Changed)
			return
		}
		fmt.Printf("Done: %d attempts considered, %d changed, %d verified, %d mismatched.\n",
			result.Summary.Considered, result.Summary.Changed, result.Summary.Verified, result.Summary.Mismatched)
	})
}
