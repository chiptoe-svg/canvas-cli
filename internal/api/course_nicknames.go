package api

import (
	"context"
	"fmt"
)

// CourseNicknamesService handles course nickname API calls
type CourseNicknamesService struct {
	client *Client
}

// NewCourseNicknamesService creates a new course nicknames service
func NewCourseNicknamesService(client *Client) *CourseNicknamesService {
	return &CourseNicknamesService{client: client}
}

// CourseNickname represents a nickname for a course
type CourseNickname struct {
	CourseID int64  `json:"course_id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

// SetCourseNicknameParams holds parameters for setting a course nickname
type SetCourseNicknameParams struct {
	Nickname string `json:"nickname"`
}

// List retrieves all course nicknames for the current user
func (s *CourseNicknamesService) List(ctx context.Context) ([]CourseNickname, error) {
	path := "/api/v1/users/self/course_nicknames"
	var result []CourseNickname
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing course nicknames: %w", err)
	}
	return result, nil
}

// Get retrieves a nickname for a specific course
func (s *CourseNicknamesService) Get(ctx context.Context, courseID int64) (*CourseNickname, error) {
	path := fmt.Sprintf("/api/v1/users/self/course_nicknames/%d", courseID)
	var result CourseNickname
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting course nickname for course %d: %w", courseID, err)
	}
	return &result, nil
}

// Set creates or updates a nickname for a course
func (s *CourseNicknamesService) Set(ctx context.Context, courseID int64, params SetCourseNicknameParams) (*CourseNickname, error) {
	path := fmt.Sprintf("/api/v1/users/self/course_nicknames/%d", courseID)
	var result CourseNickname
	if err := s.client.PutJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("setting course nickname for course %d: %w", courseID, err)
	}
	return &result, nil
}

// Delete removes the nickname for a specific course
func (s *CourseNicknamesService) Delete(ctx context.Context, courseID int64) error {
	path := fmt.Sprintf("/api/v1/users/self/course_nicknames/%d", courseID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// DeleteAll removes all course nicknames for the current user
func (s *CourseNicknamesService) DeleteAll(ctx context.Context) error {
	path := "/api/v1/users/self/course_nicknames"
	_, err := s.client.Delete(ctx, path)
	return err
}
