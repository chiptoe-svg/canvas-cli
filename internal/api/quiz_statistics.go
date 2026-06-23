package api

import (
	"context"
	"fmt"
)

// QuizStatisticsService handles quiz statistics-related API calls.
type QuizStatisticsService struct {
	client *Client
}

// NewQuizStatisticsService creates a new quiz statistics service.
func NewQuizStatisticsService(client *Client) *QuizStatisticsService {
	return &QuizStatisticsService{client: client}
}

// QuizStatistics represents statistics for a Canvas quiz.
type QuizStatistics struct {
	ID                    int64                     `json:"id"`
	QuizID                int64                     `json:"quiz_id"`
	Includes              map[string]interface{}    `json:"includes,omitempty"`
	MultipleAttemptsExist bool                      `json:"multiple_attempts_exist"`
	GeneratedAt           string                    `json:"generated_at,omitempty"`
	QuestionStatistics    []QuizQuestionStatistics  `json:"question_statistics,omitempty"`
	Links                 map[string]interface{}    `json:"links,omitempty"`
	SubmissionStatistics  *QuizSubmissionStatistics `json:"submission_statistics,omitempty"`
}

// QuizQuestionStatistics holds per-question statistics.
type QuizQuestionStatistics struct {
	ID             int64       `json:"id"`
	QuestionType   string      `json:"question_type"`
	QuestionText   string      `json:"question_text"`
	Position       int         `json:"position"`
	Responses      int         `json:"responses"`
	AnswerSets     interface{} `json:"answer_sets,omitempty"`
	Answers        interface{} `json:"answers,omitempty"`
	PointBiserials interface{} `json:"point_biserials,omitempty"`
}

// QuizSubmissionStatistics holds aggregate submission statistics.
type QuizSubmissionStatistics struct {
	UniqueCount           int     `json:"unique_count"`
	ScoreAverage          float64 `json:"score_average"`
	ScoreHigh             float64 `json:"score_high"`
	ScoreLow              float64 `json:"score_low"`
	ScoreStdev            float64 `json:"score_stdev"`
	CorrectCountAverage   float64 `json:"correct_count_average"`
	IncorrectCountAverage float64 `json:"incorrect_count_average"`
	Duration              float64 `json:"duration_average"`
}

// QuizStatisticsResponse wraps the Canvas API envelope for quiz statistics.
type QuizStatisticsResponse struct {
	QuizStatistics []QuizStatistics `json:"quiz_statistics"`
}

// List retrieves statistics for a quiz.
// Canvas API: GET /api/v1/courses/:course_id/quizzes/:quiz_id/statistics
func (s *QuizStatisticsService) List(ctx context.Context, courseID, quizID int64) ([]QuizStatistics, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/statistics", courseID, quizID)

	var resp QuizStatisticsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.QuizStatistics, nil
}
