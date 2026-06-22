package api

import (
	"context"
	"fmt"
)

// QuizIPFiltersService handles quiz IP filter-related API calls.
type QuizIPFiltersService struct {
	client *Client
}

// NewQuizIPFiltersService creates a new quiz IP filters service.
func NewQuizIPFiltersService(client *Client) *QuizIPFiltersService {
	return &QuizIPFiltersService{client: client}
}

// QuizIPFilter represents a Canvas quiz IP filter.
type QuizIPFilter struct {
	Name    string `json:"name"`
	Account string `json:"account,omitempty"`
	Filter  string `json:"filter"`
}

// QuizIPFiltersResponse wraps the Canvas API response for IP filters.
type QuizIPFiltersResponse struct {
	QuizIPFilters []QuizIPFilter `json:"quiz_ip_filters"`
}

// List retrieves all IP filters available for a course's quizzes.
// Canvas API: GET /api/v1/courses/:course_id/quizzes/ip_filters
func (s *QuizIPFiltersService) List(ctx context.Context, courseID int64) ([]QuizIPFilter, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/ip_filters", courseID)

	var resp QuizIPFiltersResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.QuizIPFilters, nil
}
