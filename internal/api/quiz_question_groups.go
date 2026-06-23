package api

import (
	"context"
	"fmt"
)

// QuizQuestionGroupsService handles quiz question group-related API calls.
type QuizQuestionGroupsService struct {
	client *Client
}

// NewQuizQuestionGroupsService creates a new quiz question groups service.
func NewQuizQuestionGroupsService(client *Client) *QuizQuestionGroupsService {
	return &QuizQuestionGroupsService{client: client}
}

// QuizQuestionGroup represents a Canvas quiz question group.
type QuizQuestionGroup struct {
	ID             int64   `json:"id"`
	QuizID         int64   `json:"quiz_id"`
	Name           string  `json:"name"`
	PickCount      int     `json:"pick_count"`
	QuestionPoints float64 `json:"question_points"`
	Position       int     `json:"position"`
}

// QuizQuestionGroupResponse wraps the Canvas API envelope for question groups.
type QuizQuestionGroupResponse struct {
	QuizGroups []QuizQuestionGroup `json:"quiz_groups"`
}

// Get retrieves a single quiz question group.
// Canvas API: GET /api/v1/courses/:course_id/quizzes/:quiz_id/groups/:id
func (s *QuizQuestionGroupsService) Get(ctx context.Context, courseID, quizID, groupID int64) (*QuizQuestionGroup, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/groups/%d", courseID, quizID, groupID)

	var resp QuizQuestionGroupResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	if len(resp.QuizGroups) == 0 {
		return nil, fmt.Errorf("quiz question group not found")
	}

	return &resp.QuizGroups[0], nil
}

// CreateQuizQuestionGroupParams holds parameters for creating a quiz question group.
type CreateQuizQuestionGroupParams struct {
	Name           string  `json:"name,omitempty"`
	PickCount      int     `json:"pick_count,omitempty"`
	QuestionPoints float64 `json:"question_points,omitempty"`
}

// Create creates a new quiz question group.
// Canvas API: POST /api/v1/courses/:course_id/quizzes/:quiz_id/groups
func (s *QuizQuestionGroupsService) Create(ctx context.Context, courseID, quizID int64, params *CreateQuizQuestionGroupParams) (*QuizQuestionGroup, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/groups", courseID, quizID)

	body := map[string]interface{}{
		"quiz_groups": []map[string]interface{}{
			{},
		},
	}

	groups, ok := body["quiz_groups"].([]map[string]interface{})
	if !ok || len(groups) == 0 {
		return nil, fmt.Errorf("internal error: invalid group data structure")
	}

	groupData := groups[0]

	if params.Name != "" {
		groupData["name"] = params.Name
	}

	if params.PickCount > 0 {
		groupData["pick_count"] = params.PickCount
	}

	if params.QuestionPoints > 0 {
		groupData["question_points"] = params.QuestionPoints
	}

	var resp QuizQuestionGroupResponse
	if err := s.client.PostJSON(ctx, path, body, &resp); err != nil {
		return nil, err
	}

	if len(resp.QuizGroups) == 0 {
		return nil, fmt.Errorf("no group returned after create")
	}

	return &resp.QuizGroups[0], nil
}

// UpdateQuizQuestionGroupParams holds parameters for updating a quiz question group.
type UpdateQuizQuestionGroupParams struct {
	Name           *string  `json:"name,omitempty"`
	PickCount      *int     `json:"pick_count,omitempty"`
	QuestionPoints *float64 `json:"question_points,omitempty"`
}

// Update updates an existing quiz question group.
// Canvas API: PUT /api/v1/courses/:course_id/quizzes/:quiz_id/groups/:id
func (s *QuizQuestionGroupsService) Update(ctx context.Context, courseID, quizID, groupID int64, params *UpdateQuizQuestionGroupParams) (*QuizQuestionGroup, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/groups/%d", courseID, quizID, groupID)

	body := map[string]interface{}{
		"quiz_groups": []map[string]interface{}{
			{},
		},
	}

	groups, ok := body["quiz_groups"].([]map[string]interface{})
	if !ok || len(groups) == 0 {
		return nil, fmt.Errorf("internal error: invalid group data structure")
	}

	groupData := groups[0]

	if params.Name != nil {
		groupData["name"] = *params.Name
	}

	if params.PickCount != nil {
		groupData["pick_count"] = *params.PickCount
	}

	if params.QuestionPoints != nil {
		groupData["question_points"] = *params.QuestionPoints
	}

	var resp QuizQuestionGroupResponse
	if err := s.client.PutJSON(ctx, path, body, &resp); err != nil {
		return nil, err
	}

	if len(resp.QuizGroups) == 0 {
		return nil, fmt.Errorf("no group returned after update")
	}

	return &resp.QuizGroups[0], nil
}

// Delete deletes a quiz question group.
// Canvas API: DELETE /api/v1/courses/:course_id/quizzes/:quiz_id/groups/:id
func (s *QuizQuestionGroupsService) Delete(ctx context.Context, courseID, quizID, groupID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/groups/%d", courseID, quizID, groupID)

	_, err := s.client.Delete(ctx, path)
	return err
}

// ReorderItems reorders the items in a quiz question group.
// Canvas API: POST /api/v1/courses/:course_id/quizzes/:quiz_id/groups/:id/reorder
func (s *QuizQuestionGroupsService) ReorderItems(ctx context.Context, courseID, quizID, groupID int64, order []ReorderItem) error {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/groups/%d/reorder", courseID, quizID, groupID)

	body := map[string]interface{}{
		"order": order,
	}

	return s.client.PostJSON(ctx, path, body, nil)
}
