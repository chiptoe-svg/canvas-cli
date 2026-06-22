package api

import (
	"context"
	"fmt"
)

// QuizExtensionsService handles quiz extension-related API calls.
type QuizExtensionsService struct {
	client *Client
}

// NewQuizExtensionsService creates a new quiz extensions service.
func NewQuizExtensionsService(client *Client) *QuizExtensionsService {
	return &QuizExtensionsService{client: client}
}

// QuizExtension represents a Canvas quiz extension for a student.
type QuizExtension struct {
	UserID           int64  `json:"user_id"`
	QuizID           int64  `json:"quiz_id"`
	ExtraAttempts    int    `json:"extra_attempts,omitempty"`
	ExtraTime        int    `json:"extra_time,omitempty"`
	ManuallyUnlocked bool   `json:"manually_unlocked,omitempty"`
	EndAt            string `json:"end_at,omitempty"`
}

// QuizExtensionEntry is a single extension entry sent in the API body.
type QuizExtensionEntry struct {
	UserID           int64 `json:"user_id"`
	ExtraAttempts    int   `json:"extra_attempts,omitempty"`
	ExtraTime        int   `json:"extra_time,omitempty"`
	ManuallyUnlocked bool  `json:"manually_unlocked,omitempty"`
	ExtendFromNow    int   `json:"extend_from_now,omitempty"`
}

// QuizExtensionsResponse wraps the API response envelope.
type QuizExtensionsResponse struct {
	QuizExtensions []QuizExtension `json:"quiz_extensions"`
}

// Create sets quiz extensions for one or more students.
// Canvas API: POST /api/v1/courses/:course_id/quizzes/:quiz_id/extensions
func (s *QuizExtensionsService) Create(ctx context.Context, courseID, quizID int64, extensions []QuizExtensionEntry) ([]QuizExtension, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/extensions", courseID, quizID)

	body := map[string]interface{}{
		"quiz_extensions": extensions,
	}

	var resp QuizExtensionsResponse
	if err := s.client.PostJSON(ctx, path, body, &resp); err != nil {
		return nil, err
	}

	return resp.QuizExtensions, nil
}
