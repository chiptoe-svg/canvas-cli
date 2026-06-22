package api

import (
	"context"
	"fmt"
	"net/url"
)

// AuditLogEvent represents a single event in a Canvas audit log.
type AuditLogEvent struct {
	ID          string      `json:"id,omitempty"`
	CreatedAt   string      `json:"created_at,omitempty"`
	EventType   string      `json:"event_type,omitempty"`
	PseudonymID int64       `json:"pseudonym_id,omitempty"`
	UserID      int64       `json:"user_id,omitempty"`
	User        *User       `json:"user,omitempty"`
	Pseudonym   interface{} `json:"pseudonym,omitempty"`
	Event       interface{} `json:"event,omitempty"`
}

// AuditLogsService handles Canvas audit log API calls.
type AuditLogsService struct {
	client *Client
}

// NewAuditLogsService creates a new AuditLogsService.
func NewAuditLogsService(client *Client) *AuditLogsService {
	return &AuditLogsService{client: client}
}

// AuditLogOptions holds common query parameters for audit log requests.
type AuditLogOptions struct {
	StartTime string // ISO 8601
	EndTime   string // ISO 8601
	PerPage   int
}

func buildAuditLogQuery(opts *AuditLogOptions) string {
	if opts == nil {
		return ""
	}
	q := url.Values{}
	if opts.StartTime != "" {
		q.Set("start_time", opts.StartTime)
	}
	if opts.EndTime != "" {
		q.Set("end_time", opts.EndTime)
	}
	if opts.PerPage > 0 {
		q.Set("per_page", fmt.Sprintf("%d", opts.PerPage))
	}
	if len(q) > 0 {
		return "?" + q.Encode()
	}
	return ""
}

// ListAuthenticationForAccount retrieves authentication audit events for an account.
func (s *AuditLogsService) ListAuthenticationForAccount(ctx context.Context, accountID int64, opts *AuditLogOptions) ([]AuditLogEvent, error) {
	path := fmt.Sprintf("/api/v1/audit/authentication/accounts/%d", accountID) + buildAuditLogQuery(opts)

	var events []AuditLogEvent
	if err := s.client.GetAllPages(ctx, path, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// ListAuthenticationForLogin retrieves authentication audit events for a login.
func (s *AuditLogsService) ListAuthenticationForLogin(ctx context.Context, loginID int64, opts *AuditLogOptions) ([]AuditLogEvent, error) {
	path := fmt.Sprintf("/api/v1/audit/authentication/logins/%d", loginID) + buildAuditLogQuery(opts)

	var events []AuditLogEvent
	if err := s.client.GetAllPages(ctx, path, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// ListAuthenticationForUser retrieves authentication audit events for a user.
func (s *AuditLogsService) ListAuthenticationForUser(ctx context.Context, userID int64, opts *AuditLogOptions) ([]AuditLogEvent, error) {
	path := fmt.Sprintf("/api/v1/audit/authentication/users/%d", userID) + buildAuditLogQuery(opts)

	var events []AuditLogEvent
	if err := s.client.GetAllPages(ctx, path, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// ListCourseEventsForAccount retrieves course audit events for an account.
func (s *AuditLogsService) ListCourseEventsForAccount(ctx context.Context, accountID int64, opts *AuditLogOptions) ([]AuditLogEvent, error) {
	path := fmt.Sprintf("/api/v1/audit/course/accounts/%d", accountID) + buildAuditLogQuery(opts)

	var events []AuditLogEvent
	if err := s.client.GetAllPages(ctx, path, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// ListCourseEventsForCourse retrieves course audit events for a specific course.
func (s *AuditLogsService) ListCourseEventsForCourse(ctx context.Context, courseID int64, opts *AuditLogOptions) ([]AuditLogEvent, error) {
	path := fmt.Sprintf("/api/v1/audit/course/courses/%d", courseID) + buildAuditLogQuery(opts)

	var events []AuditLogEvent
	if err := s.client.GetAllPages(ctx, path, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// ListGradeChangeEvents retrieves grade-change audit events (global).
func (s *AuditLogsService) ListGradeChangeEvents(ctx context.Context, opts *AuditLogOptions) ([]AuditLogEvent, error) {
	path := "/api/v1/audit/grade_change" + buildAuditLogQuery(opts)

	var events []AuditLogEvent
	if err := s.client.GetAllPages(ctx, path, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// ListGradeChangeForAssignment retrieves grade-change events for an assignment.
func (s *AuditLogsService) ListGradeChangeForAssignment(ctx context.Context, assignmentID int64, opts *AuditLogOptions) ([]AuditLogEvent, error) {
	path := fmt.Sprintf("/api/v1/audit/grade_change/assignments/%d", assignmentID) + buildAuditLogQuery(opts)

	var events []AuditLogEvent
	if err := s.client.GetAllPages(ctx, path, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// ListGradeChangeForCourse retrieves grade-change events for a course.
func (s *AuditLogsService) ListGradeChangeForCourse(ctx context.Context, courseID int64, opts *AuditLogOptions) ([]AuditLogEvent, error) {
	path := fmt.Sprintf("/api/v1/audit/grade_change/courses/%d", courseID) + buildAuditLogQuery(opts)

	var events []AuditLogEvent
	if err := s.client.GetAllPages(ctx, path, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// ListGradeChangeForGrader retrieves grade-change events for a grader.
func (s *AuditLogsService) ListGradeChangeForGrader(ctx context.Context, graderID int64, opts *AuditLogOptions) ([]AuditLogEvent, error) {
	path := fmt.Sprintf("/api/v1/audit/grade_change/graders/%d", graderID) + buildAuditLogQuery(opts)

	var events []AuditLogEvent
	if err := s.client.GetAllPages(ctx, path, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// ListGradeChangeForStudent retrieves grade-change events for a student.
func (s *AuditLogsService) ListGradeChangeForStudent(ctx context.Context, studentID int64, opts *AuditLogOptions) ([]AuditLogEvent, error) {
	path := fmt.Sprintf("/api/v1/audit/grade_change/students/%d", studentID) + buildAuditLogQuery(opts)

	var events []AuditLogEvent
	if err := s.client.GetAllPages(ctx, path, &events); err != nil {
		return nil, err
	}

	return events, nil
}
