package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// QuizSubmissionQuestionsService handles quiz submission question-related API calls.
type QuizSubmissionQuestionsService struct {
	client *Client
}

// NewQuizSubmissionQuestionsService creates a new quiz submission questions service.
func NewQuizSubmissionQuestionsService(client *Client) *QuizSubmissionQuestionsService {
	return &QuizSubmissionQuestionsService{client: client}
}

// QuizSubmissionQuestion represents a question in the context of a quiz submission.
//
// Answer is kept raw because its shape depends on the question type
// (https://canvas.instructure.com/doc/api/quiz_submission_questions.html,
// Appendix "Question Answer Formats"): for multiple_choice_question and
// true_false_question it is the selected answer's ID, which Canvas returns
// as a number or as a numeric string.
type QuizSubmissionQuestion struct {
	ID              int64           `json:"id"`
	FlaggedAt       interface{}     `json:"flagged,omitempty"`
	QuestionText    string          `json:"question_text,omitempty"`
	QuestionType    string          `json:"question_type,omitempty"`
	Answer          json.RawMessage `json:"answer,omitempty"`
	Answers         interface{}     `json:"answers,omitempty"`
	PointsPossible  float64         `json:"points_possible,omitempty"`
	CorrectComments string          `json:"correct_comments,omitempty"`
	Position        int             `json:"position,omitempty"`
}

// QuizSubmissionQuestionsResponse wraps the Canvas API envelope.
type QuizSubmissionQuestionsResponse struct {
	QuizSubmissionQuestions []QuizSubmissionQuestion `json:"quiz_submission_questions"`
}

// ListQuizSubmissionQuestionsOptions holds options for listing submission questions.
type ListQuizSubmissionQuestionsOptions struct {
	Include []string // e.g. "quiz_question"
}

// List retrieves all questions for a quiz submission.
// Canvas API: GET /api/v1/quiz_submissions/:quiz_submission_id/questions
func (s *QuizSubmissionQuestionsService) List(ctx context.Context, submissionID int64, opts *ListQuizSubmissionQuestionsOptions) ([]QuizSubmissionQuestion, error) {
	path := fmt.Sprintf("/api/v1/quiz_submissions/%d/questions", submissionID)

	if opts != nil && len(opts.Include) > 0 {
		query := url.Values{}
		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var resp QuizSubmissionQuestionsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.QuizSubmissionQuestions, nil
}

// QuizSubmissionAnswerParams holds a single question's answer in an answer batch.
type QuizSubmissionAnswerParams struct {
	ID     int64       `json:"id"`
	Answer interface{} `json:"answer"`
}

// AnswerQuizSubmissionQuestionsParams holds parameters for answering submission questions.
type AnswerQuizSubmissionQuestionsParams struct {
	Attempt         int                          `json:"attempt"`
	ValidationToken string                       `json:"validation_token"`
	AccessCode      string                       `json:"access_code,omitempty"`
	QuizQuestions   []QuizSubmissionAnswerParams `json:"quiz_questions"`
}

// Answer submits answers for questions in a quiz submission.
// Canvas API: POST /api/v1/quiz_submissions/:quiz_submission_id/questions
func (s *QuizSubmissionQuestionsService) Answer(ctx context.Context, submissionID int64, params *AnswerQuizSubmissionQuestionsParams) ([]QuizSubmissionQuestion, error) {
	path := fmt.Sprintf("/api/v1/quiz_submissions/%d/questions", submissionID)

	var resp QuizSubmissionQuestionsResponse
	if err := s.client.PostJSON(ctx, path, params, &resp); err != nil {
		return nil, err
	}

	return resp.QuizSubmissionQuestions, nil
}

// FlagQuizSubmissionQuestion flags a question in a quiz submission for review.
// Canvas API: PUT /api/v1/quiz_submissions/:quiz_submission_id/questions/:id/flag
func (s *QuizSubmissionQuestionsService) Flag(ctx context.Context, submissionID, questionID int64, attempt int, validationToken string) (*QuizSubmissionQuestion, error) {
	path := fmt.Sprintf("/api/v1/quiz_submissions/%d/questions/%d/flag", submissionID, questionID)

	body := map[string]interface{}{
		"attempt":          attempt,
		"validation_token": validationToken,
	}

	var resp QuizSubmissionQuestionsResponse
	if err := s.client.PutJSON(ctx, path, body, &resp); err != nil {
		return nil, err
	}

	if len(resp.QuizSubmissionQuestions) == 0 {
		return nil, fmt.Errorf("unexpected empty response from flag operation")
	}

	return &resp.QuizSubmissionQuestions[0], nil
}

// UnflagQuizSubmissionQuestion removes a flag from a question in a quiz submission.
// Canvas API: PUT /api/v1/quiz_submissions/:quiz_submission_id/questions/:id/unflag
func (s *QuizSubmissionQuestionsService) Unflag(ctx context.Context, submissionID, questionID int64, attempt int, validationToken string) (*QuizSubmissionQuestion, error) {
	path := fmt.Sprintf("/api/v1/quiz_submissions/%d/questions/%d/unflag", submissionID, questionID)

	body := map[string]interface{}{
		"attempt":          attempt,
		"validation_token": validationToken,
	}

	var resp QuizSubmissionQuestionsResponse
	if err := s.client.PutJSON(ctx, path, body, &resp); err != nil {
		return nil, err
	}

	if len(resp.QuizSubmissionQuestions) == 0 {
		return nil, fmt.Errorf("unexpected empty response from unflag operation")
	}

	return &resp.QuizSubmissionQuestions[0], nil
}
