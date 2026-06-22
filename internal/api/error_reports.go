package api

import "context"

// ErrorReport represents the body sent to Canvas when filing an error report.
type ErrorReport struct {
	Subject               string `json:"subject,omitempty"`
	URL                   string `json:"url,omitempty"`
	Email                 string `json:"email,omitempty"`
	Comments              string `json:"comments,omitempty"`
	UserPerceivedSeverity string `json:"user_perceived_severity,omitempty"`
}

// ErrorReportResult is the response Canvas returns after creating an error report.
type ErrorReportResult struct {
	Logged bool   `json:"logged"`
	ID     string `json:"id,omitempty"`
}

// ErrorReportsService handles error-report submission.
type ErrorReportsService struct {
	client *Client
}

// NewErrorReportsService creates a new ErrorReportsService.
func NewErrorReportsService(client *Client) *ErrorReportsService {
	return &ErrorReportsService{client: client}
}

// Create submits an error report to Canvas.
func (s *ErrorReportsService) Create(ctx context.Context, report *ErrorReport) (*ErrorReportResult, error) {
	path := "/api/v1/error_reports"

	body := map[string]interface{}{
		"error": map[string]interface{}{
			"subject":                 report.Subject,
			"url":                     report.URL,
			"email":                   report.Email,
			"comments":                report.Comments,
			"user_perceived_severity": report.UserPerceivedSeverity,
		},
	}

	var result ErrorReportResult
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
