package api

import (
	"context"
	"fmt"
)

// GradedSubmission represents a graded submission returned by the user graded_submissions endpoint.
type GradedSubmission struct {
	ID            int64   `json:"id"`
	AssignmentID  int64   `json:"assignment_id"`
	CourseID      int64   `json:"course_id"`
	UserID        int64   `json:"user_id"`
	Score         float64 `json:"score"`
	Grade         string  `json:"grade"`
	SubmittedAt   string  `json:"submitted_at"`
	GradedAt      string  `json:"graded_at"`
	WorkflowState string  `json:"workflow_state"`
}

// TemporaryEnrollmentStatus holds temporary enrollment status information for a user.
type TemporaryEnrollmentStatus struct {
	IsProvider  bool `json:"is_provider"`
	IsRecipient bool `json:"is_recipient"`
	CanProvide  bool `json:"can_provide"`
}

// ContentExportParams holds parameters for creating a content export.
type ContentExportParams struct {
	ExportType string `json:"export_type"`
	SkipNotes  bool   `json:"skip_notifications,omitempty"`
}

// PageViewQuery represents a page-view async query result.
type PageViewQuery struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// PageViewQueryResults holds the results of a completed page view query.
type PageViewQueryResults struct {
	Results []PageView `json:"results"`
}

// GetGradedSubmissions retrieves graded submissions for a user.
// GET /api/v1/users/:id/graded_submissions
func (s *UsersService) GetGradedSubmissions(ctx context.Context, id int64) ([]GradedSubmission, error) {
	path := fmt.Sprintf("/api/v1/users/%d/graded_submissions", id)
	var result []GradedSubmission
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting graded submissions for user %d: %w", id, err)
	}
	return result, nil
}

// DeleteSessions ends all sessions for a user (admin action).
// DELETE /api/v1/users/:id/sessions
func (s *UsersService) DeleteSessions(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/api/v1/users/%d/sessions", id)
	_, err := s.client.Delete(ctx, path)
	return err
}

// DeleteMobileSessions ends mobile sessions for a user.
// DELETE /api/v1/users/:id/mobile_sessions
func (s *UsersService) DeleteMobileSessions(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/api/v1/users/%d/mobile_sessions", id)
	_, err := s.client.Delete(ctx, path)
	return err
}

// DeleteAllMobileSessions ends all mobile sessions (self endpoint).
// DELETE /api/v1/users/mobile_sessions
func (s *UsersService) DeleteAllMobileSessions(ctx context.Context) error {
	path := "/api/v1/users/mobile_sessions"
	_, err := s.client.Delete(ctx, path)
	return err
}

// GetActivityStreamAll retrieves the global activity stream for the current user.
// GET /api/v1/users/activity_stream
func (s *UsersService) GetActivityStreamAll(ctx context.Context) ([]ActivityStreamItem, error) {
	path := "/api/v1/users/activity_stream"
	var result []ActivityStreamItem
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting activity stream: %w", err)
	}
	return result, nil
}

// ListEportfolios retrieves ePortfolios for a user.
// GET /api/v1/users/:user_id/eportfolios
func (s *UsersService) ListEportfolios(ctx context.Context, userID int64) ([]Eportfolio, error) {
	path := fmt.Sprintf("/api/v1/users/%d/eportfolios", userID)
	var result []Eportfolio
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing eportfolios for user %d: %w", userID, err)
	}
	return result, nil
}

// GetTemporaryEnrollmentStatus retrieves the temporary enrollment status for a user.
// GET /api/v1/users/:user_id/temporary_enrollment_status
func (s *UsersService) GetTemporaryEnrollmentStatus(ctx context.Context, userID int64) (*TemporaryEnrollmentStatus, error) {
	path := fmt.Sprintf("/api/v1/users/%d/temporary_enrollment_status", userID)
	var result TemporaryEnrollmentStatus
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting temporary enrollment status for user %d: %w", userID, err)
	}
	return &result, nil
}

// GetPlannerItems retrieves planner items for a user.
// GET /api/v1/users/:user_id/planner/items
func (s *UsersService) GetPlannerItems(ctx context.Context, userID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/users/%d/planner/items", userID)
	var result []interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting planner items for user %d: %w", userID, err)
	}
	return result, nil
}

// ResetPassword initiates a password reset for a user.
// POST /api/v1/users/reset_password
func (s *UsersService) ResetPassword(ctx context.Context, email string) error {
	path := "/api/v1/users/reset_password"
	body := map[string]string{"pseudonym[unique_id]": email}
	var result interface{}
	return s.client.PostJSON(ctx, path, body, &result)
}

// SetFilesUIVersionPreference sets the Files UI version preference for a user.
// PUT /api/v1/users/:id/files_ui_version_preference
func (s *UsersService) SetFilesUIVersionPreference(ctx context.Context, id int64, version string) error {
	path := fmt.Sprintf("/api/v1/users/%d/files_ui_version_preference", id)
	body := map[string]string{"file_ui_version_preference": version}
	var result interface{}
	return s.client.PutJSON(ctx, path, body, &result)
}

// SetTextEditorPreference sets the text editor preference for a user.
// PUT /api/v1/users/:id/text_editor_preference
func (s *UsersService) SetTextEditorPreference(ctx context.Context, id int64, editor string) error {
	path := fmt.Sprintf("/api/v1/users/%d/text_editor_preference", id)
	body := map[string]string{"text_editor": editor}
	var result interface{}
	return s.client.PutJSON(ctx, path, body, &result)
}

// CreatePageViewQuery creates an async query for user page views.
// POST /api/v1/users/:user_id/page_views/query
func (s *UsersService) CreatePageViewQuery(ctx context.Context, userID int64, params map[string]interface{}) (*PageViewQuery, error) {
	path := fmt.Sprintf("/api/v1/users/%d/page_views/query", userID)
	var result PageViewQuery
	if err := s.client.PostJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("creating page view query for user %d: %w", userID, err)
	}
	return &result, nil
}

// GetPageViewQueryStatus returns the status of a page view query.
// GET /api/v1/users/:user_id/page_views/query/:query_id
func (s *UsersService) GetPageViewQueryStatus(ctx context.Context, userID int64, queryID string) (*PageViewQuery, error) {
	path := fmt.Sprintf("/api/v1/users/%d/page_views/query/%s", userID, queryID)
	var result PageViewQuery
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting page view query %s for user %d: %w", queryID, userID, err)
	}
	return &result, nil
}

// GetPageViewQueryResults returns the results of a completed page view query.
// GET /api/v1/users/:user_id/page_views/query/:query_id/results
func (s *UsersService) GetPageViewQueryResults(ctx context.Context, userID int64, queryID string) (*PageViewQueryResults, error) {
	path := fmt.Sprintf("/api/v1/users/%d/page_views/query/%s/results", userID, queryID)
	var result PageViewQueryResults
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting page view query results %s for user %d: %w", queryID, userID, err)
	}
	return &result, nil
}

// CreateSelfPageViewQuery creates an async page view query for self (no user_id).
// POST /api/v1/users/page_views/query
func (s *UsersService) CreateSelfPageViewQuery(ctx context.Context, params map[string]interface{}) (*PageViewQuery, error) {
	path := "/api/v1/users/page_views/query"
	var result PageViewQuery
	if err := s.client.PostJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("creating self page view query: %w", err)
	}
	return &result, nil
}

// GetSelfPageViewQueryStatus returns the status of a self page view query.
// GET /api/v1/users/page_views/query/:query_id
func (s *UsersService) GetSelfPageViewQueryStatus(ctx context.Context, queryID string) (*PageViewQuery, error) {
	path := fmt.Sprintf("/api/v1/users/page_views/query/%s", queryID)
	var result PageViewQuery
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting self page view query %s: %w", queryID, err)
	}
	return &result, nil
}

// GetSelfPageViewQueryResults returns the results of a completed self page view query.
// GET /api/v1/users/page_views/query/:query_id/results
func (s *UsersService) GetSelfPageViewQueryResults(ctx context.Context, queryID string) (*PageViewQueryResults, error) {
	path := fmt.Sprintf("/api/v1/users/page_views/query/%s/results", queryID)
	var result PageViewQueryResults
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting self page view query results %s: %w", queryID, err)
	}
	return &result, nil
}

// ListContentExports retrieves content exports for a user.
// GET /api/v1/users/:user_id/content_exports
func (s *UsersService) ListContentExports(ctx context.Context, userID int64) ([]ContentExport, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_exports", userID)
	var result []ContentExport
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing content exports for user %d: %w", userID, err)
	}
	return result, nil
}

// CreateContentExport creates a new content export for a user.
// POST /api/v1/users/:user_id/content_exports
func (s *UsersService) CreateContentExport(ctx context.Context, userID int64, params ContentExportParams) (*ContentExport, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_exports", userID)
	var result ContentExport
	if err := s.client.PostJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("creating content export for user %d: %w", userID, err)
	}
	return &result, nil
}

// GetContentExport retrieves a single content export for a user.
// GET /api/v1/users/:user_id/content_exports/:id
func (s *UsersService) GetContentExport(ctx context.Context, userID, exportID int64) (*ContentExport, error) {
	path := fmt.Sprintf("/api/v1/users/%d/content_exports/%d", userID, exportID)
	var result ContentExport
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting content export %d for user %d: %w", exportID, userID, err)
	}
	return &result, nil
}
