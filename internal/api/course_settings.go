package api

import (
	"context"
	"fmt"
)

// CourseSettingsService handles settings, todo, tabs, permissions and other
// course-scoped utility endpoints that don't warrant their own file.
type CourseSettingsService struct {
	client *Client
}

// NewCourseSettingsService creates a new course settings service.
func NewCourseSettingsService(client *Client) *CourseSettingsService {
	return &CourseSettingsService{client: client}
}

// CourseSettings represents Canvas course settings.
type CourseSettings struct {
	AllowStudentDiscussionTopics  bool  `json:"allow_student_discussion_topics,omitempty"`
	AllowStudentForumAttachments  bool  `json:"allow_student_forum_attachments,omitempty"`
	AllowStudentDiscussionEditing bool  `json:"allow_student_discussion_editing,omitempty"`
	GradingStandardEnabled        bool  `json:"grading_standard_enabled,omitempty"`
	GradingStandardID             int64 `json:"grading_standard_id,omitempty"`
	AllowStudentOrganizedGroups   bool  `json:"allow_student_organized_groups,omitempty"`
	HideFinalGrades               bool  `json:"hide_final_grades,omitempty"`
	HideDistributionGraphs        bool  `json:"hide_distribution_graphs,omitempty"`
	HideSectionsOnCourseUsersPage bool  `json:"hide_sections_on_course_users_page,omitempty"`
	LockAllAnnouncements          bool  `json:"lock_all_announcements,omitempty"`
	UsagePrightToggles            bool  `json:"usage_rights_required,omitempty"`
	RestrictStudentPastView       bool  `json:"restrict_student_past_view,omitempty"`
	RestrictStudentFutureView     bool  `json:"restrict_student_future_view,omitempty"`
	ShowAnnouncementsOnHomePage   bool  `json:"show_announcements_on_home_page,omitempty"`
	HomePageAnnouncementLimit     int   `json:"home_page_announcement_limit,omitempty"`
}

// CourseTodo represents a to-do item for a course.
type CourseTodo struct {
	Type       string      `json:"type"`
	Assignment interface{} `json:"assignment,omitempty"`
	Quiz       interface{} `json:"quiz,omitempty"`
	Ignore     string      `json:"ignore,omitempty"`
	IgnoreURL  string      `json:"ignore_permanently_url,omitempty"`
}

// CourseTab represents a navigation tab in a Canvas course.
type CourseTab struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Type     string `json:"type"`
	Hidden   bool   `json:"hidden,omitempty"`
	Position int    `json:"position,omitempty"`
	URL      string `json:"url,omitempty"`
	FullURL  string `json:"full_url,omitempty"`
}

// EffectiveDueDates is a map from assignment_id to user due-date info.
// Map shape: { "assignment_id": { "user_id": { due_at, grading_type, ... } } }
type EffectiveDueDates map[string]map[string]interface{}

// LatePolicy represents Canvas late/missing submission penalty settings.
type LatePolicy struct {
	ID                                  int64   `json:"id,omitempty"`
	CourseID                            int64   `json:"course_id,omitempty"`
	LateSubmissionDeductionEnabled      bool    `json:"late_submission_deduction_enabled"`
	LateSubmissionDeduction             float64 `json:"late_submission_deduction,omitempty"`
	LateSubmissionInterval              string  `json:"late_submission_interval,omitempty"`
	LateSubmissionMinimumPercent        float64 `json:"late_submission_minimum_percent,omitempty"`
	LateSubmissionMinimumPercentEnabled bool    `json:"late_submission_minimum_percent_enabled,omitempty"`
	MissingSubmissionDeductionEnabled   bool    `json:"missing_submission_deduction_enabled,omitempty"`
	MissingSubmissionDeduction          float64 `json:"missing_submission_deduction,omitempty"`
}

// latePolicyEnvelope wraps the Canvas {"late_policy": {...}} envelope.
type latePolicyEnvelope struct {
	LatePolicy LatePolicy `json:"late_policy"`
}

// GetSettings retrieves settings for a course.
// GET /api/v1/courses/:course_id/settings
func (s *CourseSettingsService) GetSettings(ctx context.Context, courseID int64) (*CourseSettings, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/settings", courseID)
	var out CourseSettings
	if err := s.client.GetJSON(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("getting settings for course %d: %w", courseID, err)
	}
	return &out, nil
}

// UpdateSettings updates settings for a course.
// PUT /api/v1/courses/:course_id/settings
func (s *CourseSettingsService) UpdateSettings(ctx context.Context, courseID int64, settings CourseSettings) (*CourseSettings, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/settings", courseID)
	var out CourseSettings
	if err := s.client.PutJSON(ctx, path, settings, &out); err != nil {
		return nil, fmt.Errorf("updating settings for course %d: %w", courseID, err)
	}
	return &out, nil
}

// GetTodo lists to-do items for a course.
// GET /api/v1/courses/:course_id/todo
func (s *CourseSettingsService) GetTodo(ctx context.Context, courseID int64) ([]CourseTodo, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/todo", courseID)
	var out []CourseTodo
	if err := s.client.GetAllPages(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("getting todo for course %d: %w", courseID, err)
	}
	return out, nil
}

// ListTabs lists navigation tabs for a course.
// GET /api/v1/courses/:course_id/tabs
func (s *CourseSettingsService) ListTabs(ctx context.Context, courseID int64) ([]CourseTab, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/tabs", courseID)
	var out []CourseTab
	if err := s.client.GetAllPages(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("listing tabs for course %d: %w", courseID, err)
	}
	return out, nil
}

// UpdateTab updates a navigation tab for a course.
// PUT /api/v1/courses/:course_id/tabs/:tab_id
func (s *CourseSettingsService) UpdateTab(ctx context.Context, courseID int64, tabID string, tab CourseTab) (*CourseTab, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/tabs/%s", courseID, tabID)
	var out CourseTab
	if err := s.client.PutJSON(ctx, path, tab, &out); err != nil {
		return nil, fmt.Errorf("updating tab %s for course %d: %w", tabID, courseID, err)
	}
	return &out, nil
}

// GetPermissions retrieves permission settings for the current user in a course.
// GET /api/v1/courses/:course_id/permissions
func (s *CourseSettingsService) GetPermissions(ctx context.Context, courseID int64) (map[string]bool, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/permissions", courseID)
	var out map[string]bool
	if err := s.client.GetJSON(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("getting permissions for course %d: %w", courseID, err)
	}
	return out, nil
}

// GetEffectiveDueDates retrieves effective due dates for all assignments in a course.
// GET /api/v1/courses/:course_id/effective_due_dates
func (s *CourseSettingsService) GetEffectiveDueDates(ctx context.Context, courseID int64) (EffectiveDueDates, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/effective_due_dates", courseID)
	var out EffectiveDueDates
	if err := s.client.GetJSON(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("getting effective due dates for course %d: %w", courseID, err)
	}
	return out, nil
}

// GetLatePolicy retrieves the late policy for a course.
// GET /api/v1/courses/:id/late_policy
func (s *CourseSettingsService) GetLatePolicy(ctx context.Context, courseID int64) (*LatePolicy, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/late_policy", courseID)
	var envelope latePolicyEnvelope
	if err := s.client.GetJSON(ctx, path, &envelope); err != nil {
		return nil, fmt.Errorf("getting late policy for course %d: %w", courseID, err)
	}
	return &envelope.LatePolicy, nil
}

// CreateLatePolicy creates a late policy for a course.
// POST /api/v1/courses/:id/late_policy
func (s *CourseSettingsService) CreateLatePolicy(ctx context.Context, courseID int64, policy LatePolicy) (*LatePolicy, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/late_policy", courseID)
	body := latePolicyEnvelope{LatePolicy: policy}
	var envelope latePolicyEnvelope
	if err := s.client.PostJSON(ctx, path, body, &envelope); err != nil {
		return nil, fmt.Errorf("creating late policy for course %d: %w", courseID, err)
	}
	return &envelope.LatePolicy, nil
}

// UpdateLatePolicy patches the late policy for a course.
// PATCH /api/v1/courses/:id/late_policy
func (s *CourseSettingsService) UpdateLatePolicy(ctx context.Context, courseID int64, policy LatePolicy) (*LatePolicy, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/late_policy", courseID)
	body := latePolicyEnvelope{LatePolicy: policy}
	var envelope latePolicyEnvelope
	if err := s.client.PatchJSON(ctx, path, body, &envelope); err != nil {
		return nil, fmt.Errorf("updating late policy for course %d: %w", courseID, err)
	}
	return &envelope.LatePolicy, nil
}

// GetRecentStudents retrieves recently-enrolled students for a course.
// GET /api/v1/courses/:course_id/recent_students
func (s *CourseSettingsService) GetRecentStudents(ctx context.Context, courseID int64) ([]User, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/recent_students", courseID)
	var out []User
	if err := s.client.GetAllPages(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("getting recent students for course %d: %w", courseID, err)
	}
	return out, nil
}
