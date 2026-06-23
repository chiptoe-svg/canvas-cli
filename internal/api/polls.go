package api

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Poll represents a Canvas poll
type Poll struct {
	ID           int64          `json:"id"`
	Question     string         `json:"question"`
	Description  string         `json:"description,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UserID       int64          `json:"user_id"`
	TotalResults map[string]int `json:"total_results,omitempty"`
}

// PollChoice represents a single answer option within a poll
type PollChoice struct {
	ID        int64     `json:"id"`
	PollID    int64     `json:"poll_id"`
	IsCorrect bool      `json:"is_correct"`
	Text      string    `json:"text"`
	Position  int       `json:"position,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PollSession represents an instance of a poll opened for a course
type PollSession struct {
	ID               int64            `json:"id"`
	PollID           int64            `json:"poll_id"`
	CourseID         int64            `json:"course_id"`
	CourseSectionID  int64            `json:"course_section_id,omitempty"`
	IsPublished      bool             `json:"is_published"`
	HasPublicResults bool             `json:"has_public_results"`
	CreatedAt        time.Time        `json:"created_at"`
	Results          map[string]int   `json:"results,omitempty"`
	PollSubmissions  []PollSubmission `json:"poll_submissions,omitempty"`
}

// PollSubmission represents a student's vote in a poll session
type PollSubmission struct {
	ID           int64     `json:"id"`
	PollChoiceID int64     `json:"poll_choice_id"`
	UserID       int64     `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// pollsResponse wraps the Canvas API envelope for polls endpoints.
type pollsResponse struct {
	Polls []Poll `json:"polls"`
}

// pollChoicesResponse wraps the Canvas API envelope for poll_choices endpoints.
type pollChoicesResponse struct {
	PollChoices []PollChoice `json:"poll_choices"`
}

// pollSessionsResponse wraps the Canvas API envelope for poll_sessions endpoints.
type pollSessionsResponse struct {
	PollSessions []PollSession `json:"poll_sessions"`
}

// pollSubmissionsResponse wraps the Canvas API envelope for poll_submissions endpoints.
type pollSubmissionsResponse struct {
	PollSubmissions []PollSubmission `json:"poll_submissions"`
}

// PollsService handles poll-related API calls
type PollsService struct {
	client *Client
}

// NewPollsService creates a new polls service
func NewPollsService(client *Client) *PollsService {
	return &PollsService{client: client}
}

// ---- Poll CRUD ----

// ListPolls retrieves a paginated list of polls for the current user
func (s *PollsService) ListPolls(ctx context.Context) ([]Poll, error) {
	path := "/api/v1/polls"

	var resp pollsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.Polls, nil
}

// GetPoll retrieves a single poll by ID
func (s *PollsService) GetPoll(ctx context.Context, pollID int64) (*Poll, error) {
	path := fmt.Sprintf("/api/v1/polls/%d", pollID)

	var resp pollsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	if len(resp.Polls) == 0 {
		return nil, fmt.Errorf("poll %d not found", pollID)
	}

	return &resp.Polls[0], nil
}

// CreatePollParams holds parameters for creating a poll
type CreatePollParams struct {
	Question    string
	Description string
}

// CreatePoll creates a new poll
func (s *PollsService) CreatePoll(ctx context.Context, params *CreatePollParams) (*Poll, error) {
	path := "/api/v1/polls"

	pollData := map[string]interface{}{
		"question": params.Question,
	}

	if params.Description != "" {
		pollData["description"] = params.Description
	}

	body := map[string]interface{}{
		"polls": []interface{}{pollData},
	}

	var resp pollsResponse
	if err := s.client.PostJSON(ctx, path, body, &resp); err != nil {
		return nil, err
	}

	if len(resp.Polls) == 0 {
		return nil, fmt.Errorf("unexpected empty response from poll creation")
	}

	return &resp.Polls[0], nil
}

// UpdatePollParams holds parameters for updating a poll
type UpdatePollParams struct {
	Question    string
	Description string
}

// UpdatePoll updates an existing poll
func (s *PollsService) UpdatePoll(ctx context.Context, pollID int64, params *UpdatePollParams) (*Poll, error) {
	path := fmt.Sprintf("/api/v1/polls/%d", pollID)

	pollData := map[string]interface{}{}

	if params.Question != "" {
		pollData["question"] = params.Question
	}

	if params.Description != "" {
		pollData["description"] = params.Description
	}

	body := map[string]interface{}{
		"polls": []interface{}{pollData},
	}

	var resp pollsResponse
	if err := s.client.PutJSON(ctx, path, body, &resp); err != nil {
		return nil, err
	}

	if len(resp.Polls) == 0 {
		return nil, fmt.Errorf("unexpected empty response from poll update")
	}

	return &resp.Polls[0], nil
}

// DeletePoll removes a poll and all associated data
func (s *PollsService) DeletePoll(ctx context.Context, pollID int64) error {
	path := fmt.Sprintf("/api/v1/polls/%d", pollID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// ---- Poll Choice CRUD ----

// ListPollChoices retrieves poll choices for a given poll
func (s *PollsService) ListPollChoices(ctx context.Context, pollID int64) ([]PollChoice, error) {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_choices", pollID)

	var resp pollChoicesResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.PollChoices, nil
}

// GetPollChoice retrieves a single poll choice
func (s *PollsService) GetPollChoice(ctx context.Context, pollID, choiceID int64) (*PollChoice, error) {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_choices/%d", pollID, choiceID)

	var resp pollChoicesResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	if len(resp.PollChoices) == 0 {
		return nil, fmt.Errorf("poll choice %d not found", choiceID)
	}

	return &resp.PollChoices[0], nil
}

// CreatePollChoiceParams holds parameters for creating a poll choice
type CreatePollChoiceParams struct {
	Text      string
	IsCorrect bool
	Position  int
}

// CreatePollChoice creates a new choice for a poll
func (s *PollsService) CreatePollChoice(ctx context.Context, pollID int64, params *CreatePollChoiceParams) (*PollChoice, error) {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_choices", pollID)

	choiceData := map[string]interface{}{
		"text": params.Text,
	}

	if params.IsCorrect {
		choiceData["is_correct"] = true
	}

	if params.Position > 0 {
		choiceData["position"] = params.Position
	}

	body := map[string]interface{}{
		"poll_choices": []interface{}{choiceData},
	}

	var resp pollChoicesResponse
	if err := s.client.PostJSON(ctx, path, body, &resp); err != nil {
		return nil, err
	}

	if len(resp.PollChoices) == 0 {
		return nil, fmt.Errorf("unexpected empty response from poll choice creation")
	}

	return &resp.PollChoices[0], nil
}

// UpdatePollChoiceParams holds parameters for updating a poll choice
type UpdatePollChoiceParams struct {
	Text      string
	IsCorrect *bool
	Position  int
}

// UpdatePollChoice updates an existing poll choice
func (s *PollsService) UpdatePollChoice(ctx context.Context, pollID, choiceID int64, params *UpdatePollChoiceParams) (*PollChoice, error) {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_choices/%d", pollID, choiceID)

	choiceData := map[string]interface{}{}

	if params.Text != "" {
		choiceData["text"] = params.Text
	}

	if params.IsCorrect != nil {
		choiceData["is_correct"] = *params.IsCorrect
	}

	if params.Position > 0 {
		choiceData["position"] = params.Position
	}

	body := map[string]interface{}{
		"poll_choices": []interface{}{choiceData},
	}

	var resp pollChoicesResponse
	if err := s.client.PutJSON(ctx, path, body, &resp); err != nil {
		return nil, err
	}

	if len(resp.PollChoices) == 0 {
		return nil, fmt.Errorf("unexpected empty response from poll choice update")
	}

	return &resp.PollChoices[0], nil
}

// DeletePollChoice removes a poll choice
func (s *PollsService) DeletePollChoice(ctx context.Context, pollID, choiceID int64) error {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_choices/%d", pollID, choiceID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// ---- Poll Session CRUD ----

// ListPollSessions retrieves poll sessions for a given poll
func (s *PollsService) ListPollSessions(ctx context.Context, pollID int64) ([]PollSession, error) {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_sessions", pollID)

	var resp pollSessionsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.PollSessions, nil
}

// GetPollSession retrieves a single poll session
func (s *PollsService) GetPollSession(ctx context.Context, pollID, sessionID int64) (*PollSession, error) {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_sessions/%d", pollID, sessionID)

	var resp pollSessionsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	if len(resp.PollSessions) == 0 {
		return nil, fmt.Errorf("poll session %d not found", sessionID)
	}

	return &resp.PollSessions[0], nil
}

// CreatePollSessionParams holds parameters for creating a poll session
type CreatePollSessionParams struct {
	CourseID         int64
	CourseSectionID  int64
	HasPublicResults bool
}

// CreatePollSession creates a new poll session
func (s *PollsService) CreatePollSession(ctx context.Context, pollID int64, params *CreatePollSessionParams) (*PollSession, error) {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_sessions", pollID)

	sessionData := map[string]interface{}{
		"course_id": params.CourseID,
	}

	if params.CourseSectionID > 0 {
		sessionData["course_section_id"] = params.CourseSectionID
	}

	if params.HasPublicResults {
		sessionData["has_public_results"] = true
	}

	body := map[string]interface{}{
		"poll_sessions": []interface{}{sessionData},
	}

	var resp pollSessionsResponse
	if err := s.client.PostJSON(ctx, path, body, &resp); err != nil {
		return nil, err
	}

	if len(resp.PollSessions) == 0 {
		return nil, fmt.Errorf("unexpected empty response from poll session creation")
	}

	return &resp.PollSessions[0], nil
}

// UpdatePollSessionParams holds parameters for updating a poll session
type UpdatePollSessionParams struct {
	CourseID         int64
	CourseSectionID  int64
	HasPublicResults *bool
}

// UpdatePollSession updates an existing poll session
func (s *PollsService) UpdatePollSession(ctx context.Context, pollID, sessionID int64, params *UpdatePollSessionParams) (*PollSession, error) {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_sessions/%d", pollID, sessionID)

	sessionData := map[string]interface{}{}

	if params.CourseID > 0 {
		sessionData["course_id"] = params.CourseID
	}

	if params.CourseSectionID > 0 {
		sessionData["course_section_id"] = params.CourseSectionID
	}

	if params.HasPublicResults != nil {
		sessionData["has_public_results"] = *params.HasPublicResults
	}

	body := map[string]interface{}{
		"poll_sessions": []interface{}{sessionData},
	}

	var resp pollSessionsResponse
	if err := s.client.PutJSON(ctx, path, body, &resp); err != nil {
		return nil, err
	}

	if len(resp.PollSessions) == 0 {
		return nil, fmt.Errorf("unexpected empty response from poll session update")
	}

	return &resp.PollSessions[0], nil
}

// DeletePollSession removes a poll session
func (s *PollsService) DeletePollSession(ctx context.Context, pollID, sessionID int64) error {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_sessions/%d", pollID, sessionID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// OpenPollSession opens a poll session for student participation
func (s *PollsService) OpenPollSession(ctx context.Context, pollID, sessionID int64) (*PollSession, error) {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_sessions/%d/open", pollID, sessionID)

	var resp pollSessionsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	if len(resp.PollSessions) == 0 {
		return nil, fmt.Errorf("unexpected empty response from poll session open")
	}

	return &resp.PollSessions[0], nil
}

// ClosePollSession closes a poll session to stop accepting submissions
func (s *PollsService) ClosePollSession(ctx context.Context, pollID, sessionID int64) (*PollSession, error) {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_sessions/%d/close", pollID, sessionID)

	var resp pollSessionsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	if len(resp.PollSessions) == 0 {
		return nil, fmt.Errorf("unexpected empty response from poll session close")
	}

	return &resp.PollSessions[0], nil
}

// ListOpenedPollSessions returns all currently open poll sessions across all polls
func (s *PollsService) ListOpenedPollSessions(ctx context.Context) ([]PollSession, error) {
	path := "/api/v1/poll_sessions/opened"

	var resp pollSessionsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.PollSessions, nil
}

// ListClosedPollSessions returns all closed poll sessions across all polls
func (s *PollsService) ListClosedPollSessions(ctx context.Context) ([]PollSession, error) {
	path := "/api/v1/poll_sessions/closed"

	var resp pollSessionsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.PollSessions, nil
}

// ---- Poll Submission ----

// GetPollSubmission retrieves a single poll submission
func (s *PollsService) GetPollSubmission(ctx context.Context, pollID, sessionID, submissionID int64) (*PollSubmission, error) {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_sessions/%d/poll_submissions/%d", pollID, sessionID, submissionID)

	var resp pollSubmissionsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	if len(resp.PollSubmissions) == 0 {
		return nil, fmt.Errorf("poll submission %d not found", submissionID)
	}

	return &resp.PollSubmissions[0], nil
}

// CreatePollSubmissionParams holds parameters for creating a poll submission (student vote)
type CreatePollSubmissionParams struct {
	PollChoiceID int64
}

// CreatePollSubmission creates a new poll submission (student vote)
func (s *PollsService) CreatePollSubmission(ctx context.Context, pollID, sessionID int64, params *CreatePollSubmissionParams) (*PollSubmission, error) {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_sessions/%d/poll_submissions", pollID, sessionID)

	body := map[string]interface{}{
		"poll_submissions": []interface{}{
			map[string]interface{}{
				"poll_choice_id": params.PollChoiceID,
			},
		},
	}

	var resp pollSubmissionsResponse
	if err := s.client.PostJSON(ctx, path, body, &resp); err != nil {
		return nil, err
	}

	if len(resp.PollSubmissions) == 0 {
		return nil, fmt.Errorf("unexpected empty response from poll submission creation")
	}

	return &resp.PollSubmissions[0], nil
}

// ListPollSubmissionsOptions holds query options for listing poll submissions
type ListPollSubmissionsOptions struct {
	Page    int
	PerPage int
}

// buildPollSubmissionsQuery builds query parameters for poll submissions
func buildPollSubmissionsQuery(opts *ListPollSubmissionsOptions) url.Values {
	q := url.Values{}

	if opts == nil {
		return q
	}

	if opts.Page > 0 {
		q.Set("page", fmt.Sprintf("%d", opts.Page))
	}

	if opts.PerPage > 0 {
		q.Set("per_page", fmt.Sprintf("%d", opts.PerPage))
	}

	return q
}

// ListPollSubmissions lists all submissions for a poll session (instructor-level access)
func (s *PollsService) ListPollSubmissions(ctx context.Context, pollID, sessionID int64, opts *ListPollSubmissionsOptions) ([]PollSubmission, error) {
	path := fmt.Sprintf("/api/v1/polls/%d/poll_sessions/%d/poll_submissions", pollID, sessionID)

	q := buildPollSubmissionsQuery(opts)
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	var resp pollSubmissionsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.PollSubmissions, nil
}
