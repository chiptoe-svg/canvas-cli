package api

import (
	"context"
	"fmt"
)

// CourseExtrasService handles miscellaneous course-scoped endpoints that cover
// quiz/assignment extensions and other utility actions.
type CourseExtrasService struct {
	client *Client
}

// NewCourseExtrasService creates a new course extras service.
func NewCourseExtrasService(client *Client) *CourseExtrasService {
	return &CourseExtrasService{client: client}
}

// QuizExtensionParams holds parameters for a quiz extension grant.
type QuizExtensionParams struct {
	UserID           int64 `json:"user_id"`
	ExtraAttempts    int   `json:"extra_attempts,omitempty"`
	ExtraTime        int   `json:"extra_time,omitempty"`
	ManuallyUnlocked bool  `json:"manually_unlocked,omitempty"`
	ExtendFromNow    int   `json:"extend_from_now,omitempty"`
}

// quizExtensionsBody wraps the Canvas quiz extensions request body.
type quizExtensionsBody struct {
	QuizExtensions []QuizExtensionParams `json:"quiz_extensions"`
}

// AssignmentExtension represents a time/attempt extension for an assignment.
type AssignmentExtension struct {
	AssignmentID  int64 `json:"assignment_id"`
	UserID        int64 `json:"user_id"`
	ExtraAttempts int   `json:"extra_attempts,omitempty"`
}

// assignmentExtensionsBody wraps the Canvas assignment extensions request body.
type assignmentExtensionsBody struct {
	AssignmentExtensions []AssignmentExtension `json:"assignment_extensions"`
}

// CreateQuizExtensions creates quiz extensions for students in a course.
// POST /api/v1/courses/:course_id/quiz_extensions
func (s *CourseExtrasService) CreateQuizExtensions(ctx context.Context, courseID int64, extensions []QuizExtensionParams) error {
	path := fmt.Sprintf("/api/v1/courses/%d/quiz_extensions", courseID)
	body := quizExtensionsBody{QuizExtensions: extensions}
	if err := s.client.PostJSON(ctx, path, body, nil); err != nil {
		return fmt.Errorf("creating quiz extensions for course %d: %w", courseID, err)
	}
	return nil
}

// CreateAssignmentExtensions creates assignment extensions for students.
// POST /api/v1/courses/:course_id/assignments/:assignment_id/extensions
func (s *CourseExtrasService) CreateAssignmentExtensions(ctx context.Context, courseID, assignmentID int64, extensions []AssignmentExtension) ([]AssignmentExtension, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/extensions", courseID, assignmentID)
	body := assignmentExtensionsBody{AssignmentExtensions: extensions}
	var out assignmentExtensionsBody
	if err := s.client.PostJSON(ctx, path, body, &out); err != nil {
		return nil, fmt.Errorf("creating assignment extensions for course %d assignment %d: %w", courseID, assignmentID, err)
	}
	return out.AssignmentExtensions, nil
}
