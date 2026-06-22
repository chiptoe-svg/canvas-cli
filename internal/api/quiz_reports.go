package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// QuizReportsService handles quiz report-related API calls.
type QuizReportsService struct {
	client *Client
}

// NewQuizReportsService creates a new quiz reports service.
func NewQuizReportsService(client *Client) *QuizReportsService {
	return &QuizReportsService{client: client}
}

// QuizReport represents a Canvas quiz report.
type QuizReport struct {
	ID           int64       `json:"id"`
	CourseID     int64       `json:"course_id"`
	QuizID       int64       `json:"quiz_id"`
	ReportType   string      `json:"report_type"`
	ReadableType string      `json:"readable_type"`
	IncludesAll  bool        `json:"includes_all_versions"`
	Anonymous    bool        `json:"anonymous"`
	GeneratesAs  string      `json:"generatable"`
	CreatedAt    string      `json:"created_at,omitempty"`
	UpdatedAt    string      `json:"updated_at,omitempty"`
	FileURL      string      `json:"file_url,omitempty"`
	Progress     interface{} `json:"progress,omitempty"`
	ProgressURL  string      `json:"progress_url,omitempty"`
}

// ListQuizReportsOptions holds options for listing quiz reports.
type ListQuizReportsOptions struct {
	IncludesAllVersions bool
}

// List retrieves all reports for a quiz.
// Canvas API: GET /api/v1/courses/:course_id/quizzes/:quiz_id/reports
func (s *QuizReportsService) List(ctx context.Context, courseID, quizID int64, opts *ListQuizReportsOptions) ([]QuizReport, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/reports", courseID, quizID)

	if opts != nil && opts.IncludesAllVersions {
		query := url.Values{}
		query.Set("includes_all_versions", strconv.FormatBool(opts.IncludesAllVersions))
		path += "?" + query.Encode()
	}

	var reports []QuizReport
	if err := s.client.GetAllPages(ctx, path, &reports); err != nil {
		return nil, err
	}

	return reports, nil
}

// Get retrieves a single quiz report.
// Canvas API: GET /api/v1/courses/:course_id/quizzes/:quiz_id/reports/:id
func (s *QuizReportsService) Get(ctx context.Context, courseID, quizID, reportID int64, includeProgress bool) (*QuizReport, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/reports/%d", courseID, quizID, reportID)

	if includeProgress {
		path += "?includes_all_versions=true"
	}

	var report QuizReport
	if err := s.client.GetJSON(ctx, path, &report); err != nil {
		return nil, err
	}

	return &report, nil
}

// CreateQuizReportParams holds parameters for creating a quiz report.
type CreateQuizReportParams struct {
	ReportType          string `json:"report_type"`
	IncludesAllVersions bool   `json:"includes_all_versions,omitempty"`
}

// Create creates (or retrieves) a quiz report.
// Canvas API: POST /api/v1/courses/:course_id/quizzes/:quiz_id/reports
func (s *QuizReportsService) Create(ctx context.Context, courseID, quizID int64, params *CreateQuizReportParams) (*QuizReport, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/reports", courseID, quizID)

	body := map[string]interface{}{
		"quiz_report": map[string]interface{}{
			"report_type": params.ReportType,
		},
	}

	if params.IncludesAllVersions {
		body["includes_all_versions"] = true
	}

	var report QuizReport
	if err := s.client.PostJSON(ctx, path, body, &report); err != nil {
		return nil, err
	}

	return &report, nil
}

// Delete deletes a quiz report.
// Canvas API: DELETE /api/v1/courses/:course_id/quizzes/:quiz_id/reports/:id
func (s *QuizReportsService) Delete(ctx context.Context, courseID, quizID, reportID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/reports/%d", courseID, quizID, reportID)

	_, err := s.client.Delete(ctx, path)
	return err
}
