package api

import (
	"encoding/json"
	"testing"
)

// TestSubmission_DecodesQuizSubmissionHistory pins the teacher-side shape
// observed live from
// GET /api/v1/courses/:course_id/assignments/:assignment_id/submissions/:user_id?include[]=submission_history
// for a classic quiz: one history entry per attempt, each with attempt,
// score and submission_data[{question_id, answer_id, correct, points}].
// answer_id arrives as a number on recent attempts and as a numeric string
// on older ones; both must survive decoding.
func TestSubmission_DecodesQuizSubmissionHistory(t *testing.T) {
	body := `{
		"id": 9001, "user_id": 12, "assignment_id": 77, "attempt": 2, "score": 9, "workflow_state": "graded",
		"submitted_at": null, "graded_at": null,
		"submission_history": [
			{"attempt": 1, "score": 5, "workflow_state": "graded", "submitted_at": null,
			 "submission_data": [
				{"question_id": 10482606, "answer_id": "1128", "correct": false, "points": 0, "text": ""}
			 ]},
			{"attempt": 2, "score": 9, "workflow_state": "graded",
			 "submission_data": [
				{"question_id": 10482606, "answer_id": 5357, "correct": true, "points": 1},
				{"question_id": 10482607, "correct": "partial", "points": 0.5}
			 ]}
		]
	}`

	var sub Submission
	if err := json.Unmarshal([]byte(body), &sub); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if sub.Attempt != 2 || sub.Score != 9 || len(sub.SubmissionHistory) != 2 {
		t.Fatalf("top level = attempt %d score %g history %d", sub.Attempt, sub.Score, len(sub.SubmissionHistory))
	}

	a1 := sub.SubmissionHistory[0]
	if a1.Attempt != 1 || a1.Score != 5 || len(a1.SubmissionData) != 1 {
		t.Errorf("attempt 1 = %+v", a1)
	}
	if d := a1.SubmissionData[0]; d.QuestionID != 10482606 || string(d.AnswerID) != `"1128"` || d.Points != 0 || string(d.Correct) != "false" {
		t.Errorf("attempt 1 data = %+v", d)
	}

	a2 := sub.SubmissionHistory[1]
	if len(a2.SubmissionData) != 2 {
		t.Fatalf("attempt 2 data = %+v", a2.SubmissionData)
	}
	if d := a2.SubmissionData[0]; string(d.AnswerID) != "5357" || d.Points != 1 || string(d.Correct) != "true" {
		t.Errorf("attempt 2 data[0] = %+v", d)
	}
	if d := a2.SubmissionData[1]; len(d.AnswerID) != 0 || d.Points != 0.5 || string(d.Correct) != `"partial"` {
		t.Errorf("attempt 2 data[1] (no answer_id, partial) = %+v", d)
	}

	// Quiz.AssignmentID is how the regrade reaches this endpoint.
	var quiz Quiz
	if err := json.Unmarshal([]byte(`{"id": 10, "title": "Q", "assignment_id": 77}`), &quiz); err != nil {
		t.Fatalf("Unmarshal quiz: %v", err)
	}
	if quiz.AssignmentID != 77 {
		t.Errorf("Quiz.AssignmentID = %d, want 77", quiz.AssignmentID)
	}
}
