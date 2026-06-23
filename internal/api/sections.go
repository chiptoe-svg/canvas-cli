package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// SectionsService handles section-related API calls
type SectionsService struct {
	client *Client
}

// NewSectionsService creates a new sections service
func NewSectionsService(client *Client) *SectionsService {
	return &SectionsService{client: client}
}

// Section represents a Canvas course section
type Section struct {
	ID                                int64      `json:"id"`
	Name                              string     `json:"name"`
	SISSectionID                      string     `json:"sis_section_id,omitempty"`
	IntegrationID                     string     `json:"integration_id,omitempty"`
	SISImportID                       int64      `json:"sis_import_id,omitempty"`
	CourseID                          int64      `json:"course_id"`
	SISCourseID                       string     `json:"sis_course_id,omitempty"`
	StartAt                           *time.Time `json:"start_at,omitempty"`
	EndAt                             *time.Time `json:"end_at,omitempty"`
	RestrictEnrollmentsToSectionDates bool       `json:"restrict_enrollments_to_section_dates"`
	NonXlistCourseID                  *int64     `json:"nonxlist_course_id,omitempty"`
	TotalStudents                     int        `json:"total_students,omitempty"`
	CreatedAt                         *time.Time `json:"created_at,omitempty"`
}

// ListSectionsOptions holds options for listing sections
type ListSectionsOptions struct {
	Include []string // students, total_students, passback_status, permissions
	Page    int
	PerPage int
}

// ListCourse retrieves sections for a course
func (s *SectionsService) ListCourse(ctx context.Context, courseID int64, opts *ListSectionsOptions) ([]Section, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/sections", courseID)

	if opts != nil {
		query := url.Values{}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var sections []Section
	if err := s.client.GetAllPages(ctx, path, &sections); err != nil {
		return nil, err
	}

	return sections, nil
}

// Get retrieves a single section by ID
func (s *SectionsService) Get(ctx context.Context, sectionID int64, include []string) (*Section, error) {
	path := fmt.Sprintf("/api/v1/sections/%d", sectionID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var section Section
	if err := s.client.GetJSON(ctx, path, &section); err != nil {
		return nil, err
	}

	return &section, nil
}

// CreateSectionParams holds parameters for creating a section
type CreateSectionParams struct {
	Name                              string
	SISSectionID                      string
	IntegrationID                     string
	StartAt                           string
	EndAt                             string
	RestrictEnrollmentsToSectionDates bool
	EnableSISReactivation             bool
}

// Create creates a new section in a course
func (s *SectionsService) Create(ctx context.Context, courseID int64, params *CreateSectionParams) (*Section, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/sections", courseID)

	body := map[string]interface{}{
		"course_section": make(map[string]interface{}),
	}

	sectionData, ok := body["course_section"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid section data structure")
	}

	if params.Name != "" {
		sectionData["name"] = params.Name
	}

	if params.SISSectionID != "" {
		sectionData["sis_section_id"] = params.SISSectionID
	}

	if params.IntegrationID != "" {
		sectionData["integration_id"] = params.IntegrationID
	}

	if params.StartAt != "" {
		sectionData["start_at"] = params.StartAt
	}

	if params.EndAt != "" {
		sectionData["end_at"] = params.EndAt
	}

	if params.RestrictEnrollmentsToSectionDates {
		sectionData["restrict_enrollments_to_section_dates"] = true
	}

	if params.EnableSISReactivation {
		sectionData["enable_sis_reactivation"] = true
	}

	var section Section
	if err := s.client.PostJSON(ctx, path, body, &section); err != nil {
		return nil, err
	}

	return &section, nil
}

// UpdateSectionParams holds parameters for updating a section
type UpdateSectionParams struct {
	Name                              *string
	SISSectionID                      *string
	IntegrationID                     *string
	StartAt                           *string
	EndAt                             *string
	RestrictEnrollmentsToSectionDates *bool
	OverrideSISStickiness             bool
}

// Update updates an existing section
func (s *SectionsService) Update(ctx context.Context, sectionID int64, params *UpdateSectionParams) (*Section, error) {
	path := fmt.Sprintf("/api/v1/sections/%d", sectionID)

	body := map[string]interface{}{
		"course_section": make(map[string]interface{}),
	}

	sectionData, ok := body["course_section"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid section data structure")
	}

	if params.Name != nil {
		sectionData["name"] = *params.Name
	}

	if params.SISSectionID != nil {
		sectionData["sis_section_id"] = *params.SISSectionID
	}

	if params.IntegrationID != nil {
		sectionData["integration_id"] = *params.IntegrationID
	}

	if params.StartAt != nil {
		sectionData["start_at"] = *params.StartAt
	}

	if params.EndAt != nil {
		sectionData["end_at"] = *params.EndAt
	}

	if params.RestrictEnrollmentsToSectionDates != nil {
		sectionData["restrict_enrollments_to_section_dates"] = *params.RestrictEnrollmentsToSectionDates
	}

	if params.OverrideSISStickiness {
		body["override_sis_stickiness"] = true
	}

	var section Section
	if err := s.client.PutJSON(ctx, path, body, &section); err != nil {
		return nil, err
	}

	return &section, nil
}

// Delete deletes a section
func (s *SectionsService) Delete(ctx context.Context, sectionID int64) (*Section, error) {
	path := fmt.Sprintf("/api/v1/sections/%d", sectionID)

	var section Section
	if err := s.client.DeleteJSON(ctx, path, &section); err != nil {
		return nil, err
	}

	return &section, nil
}

// Crosslist moves a section to a different course
func (s *SectionsService) Crosslist(ctx context.Context, sectionID, newCourseID int64, overrideSISStickiness bool) (*Section, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/crosslist/%d", sectionID, newCourseID)

	if overrideSISStickiness {
		path += "?override_sis_stickiness=true"
	}

	var section Section
	if err := s.client.PostJSON(ctx, path, nil, &section); err != nil {
		return nil, err
	}

	return &section, nil
}

// Uncrosslist returns a crosslisted section to its original course
func (s *SectionsService) Uncrosslist(ctx context.Context, sectionID int64, overrideSISStickiness bool) (*Section, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/crosslist", sectionID)

	if overrideSISStickiness {
		path += "?override_sis_stickiness=true"
	}

	var section Section
	if err := s.client.DeleteJSON(ctx, path, &section); err != nil {
		return nil, err
	}

	return &section, nil
}

// GetCourse retrieves a section within a course context
// Canvas path: GET /api/v1/courses/:course_id/sections/:id
func (s *SectionsService) GetCourse(ctx context.Context, courseID, sectionID int64, include []string) (*Section, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/sections/%d", courseID, sectionID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var section Section
	if err := s.client.GetJSON(ctx, path, &section); err != nil {
		return nil, err
	}

	return &section, nil
}

// ListUsers retrieves users enrolled in a section
// Canvas path: GET /api/v1/sections/:id/users
func (s *SectionsService) ListUsers(ctx context.Context, sectionID int64, include []string) ([]User, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/users", sectionID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var users []User
	if err := s.client.GetAllPages(ctx, path, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// CreateEnrollment enrolls a user in a section
// Canvas path: POST /api/v1/sections/:section_id/enrollments
func (s *SectionsService) CreateEnrollment(ctx context.Context, sectionID int64, params map[string]interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/enrollments", sectionID)

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, params, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// ListSubmissions retrieves submissions for an assignment in a section
// Canvas path: GET /api/v1/sections/:section_id/assignments/:assignment_id/submissions
func (s *SectionsService) ListSubmissions(ctx context.Context, sectionID, assignmentID int64, opts *ListSubmissionsOptions) ([]Submission, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/assignments/%d/submissions", sectionID, assignmentID)

	if opts != nil {
		query := url.Values{}
		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}
		if opts.Grouped {
			query.Add("grouped", "true")
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var submissions []Submission
	if err := s.client.GetAllPages(ctx, path, &submissions); err != nil {
		return nil, err
	}

	return submissions, nil
}

// GetSubmission retrieves a single submission in a section
// Canvas path: GET /api/v1/sections/:section_id/assignments/:assignment_id/submissions/:user_id
func (s *SectionsService) GetSubmission(ctx context.Context, sectionID, assignmentID, userID int64, include []string) (*Submission, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/assignments/%d/submissions/%d", sectionID, assignmentID, userID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var submission Submission
	if err := s.client.GetJSON(ctx, path, &submission); err != nil {
		return nil, err
	}

	return &submission, nil
}

// ListSubmissionsForStudents retrieves submissions for multiple students in a section
// Canvas path: GET /api/v1/sections/:section_id/students/submissions
func (s *SectionsService) ListSubmissionsForStudents(ctx context.Context, sectionID int64, studentIDs []int64, opts *ListSubmissionsOptions) ([]Submission, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/students/submissions", sectionID)

	query := url.Values{}
	for _, id := range studentIDs {
		query.Add("student_ids[]", strconv.FormatInt(id, 10))
	}
	if opts != nil {
		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}
	}
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var submissions []Submission
	if err := s.client.GetAllPages(ctx, path, &submissions); err != nil {
		return nil, err
	}

	return submissions, nil
}

// GradeSubmission grades a submission in a section
// Canvas path: PUT /api/v1/sections/:section_id/assignments/:assignment_id/submissions/:user_id
func (s *SectionsService) GradeSubmission(ctx context.Context, sectionID, assignmentID, userID int64, params *GradeSubmissionParams) (*Submission, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/assignments/%d/submissions/%d", sectionID, assignmentID, userID)

	// Build the same body as SubmissionsService.Grade to ensure full field coverage
	// (SecondsLateOverride, RubricAssessment, all comment fields).
	body := buildGradeSubmissionBody(params)

	var submission Submission
	if err := s.client.PutJSON(ctx, path, body, &submission); err != nil {
		return nil, err
	}

	return &submission, nil
}

// GetAssignmentOverride retrieves assignment override for a section
// Canvas path: GET /api/v1/sections/:course_section_id/assignments/:assignment_id/override
func (s *SectionsService) GetAssignmentOverride(ctx context.Context, sectionID, assignmentID int64) (*AssignmentOverride, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/assignments/%d/override", sectionID, assignmentID)

	var override AssignmentOverride
	if err := s.client.GetJSON(ctx, path, &override); err != nil {
		return nil, err
	}

	return &override, nil
}

// GetSubmissionSummary returns the submission summary for an assignment in a section
// Canvas path: GET /api/v1/sections/:section_id/assignments/:assignment_id/submission_summary
func (s *SectionsService) GetSubmissionSummary(ctx context.Context, sectionID, assignmentID int64) (*SubmissionSummary, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/assignments/%d/submission_summary", sectionID, assignmentID)

	var summary SubmissionSummary
	if err := s.client.GetJSON(ctx, path, &summary); err != nil {
		return nil, err
	}

	return &summary, nil
}

// SubmitAssignment submits an assignment in a section context
// Canvas path: POST /api/v1/sections/:section_id/assignments/:assignment_id/submissions
func (s *SectionsService) SubmitAssignment(ctx context.Context, sectionID, assignmentID int64, params map[string]interface{}) (*Submission, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/assignments/%d/submissions", sectionID, assignmentID)

	var submission Submission
	if err := s.client.PostJSON(ctx, path, params, &submission); err != nil {
		return nil, err
	}

	return &submission, nil
}

// MarkSubmissionAsRead marks a submission as read in a section
// Canvas path: PUT /api/v1/sections/:section_id/assignments/:assignment_id/submissions/:user_id/read
func (s *SectionsService) MarkSubmissionAsRead(ctx context.Context, sectionID, assignmentID, userID int64) error {
	path := fmt.Sprintf("/api/v1/sections/%d/assignments/%d/submissions/%d/read", sectionID, assignmentID, userID)

	return s.client.PutJSON(ctx, path, nil, nil)
}

// MarkSubmissionAsUnread marks a submission as unread in a section
// Canvas path: DELETE /api/v1/sections/:section_id/assignments/:assignment_id/submissions/:user_id/read
func (s *SectionsService) MarkSubmissionAsUnread(ctx context.Context, sectionID, assignmentID, userID int64) error {
	path := fmt.Sprintf("/api/v1/sections/%d/assignments/%d/submissions/%d/read", sectionID, assignmentID, userID)

	_, err := s.client.Delete(ctx, path)
	return err
}

// BulkMarkRead marks all submissions as read in a section
// Canvas path: PUT /api/v1/sections/:section_id/submissions/bulk_mark_read
func (s *SectionsService) BulkMarkRead(ctx context.Context, sectionID int64) error {
	path := fmt.Sprintf("/api/v1/sections/%d/submissions/bulk_mark_read", sectionID)

	return s.client.PutJSON(ctx, path, nil, nil)
}

// UpdateGrades updates grades for multiple submissions in a section
// Canvas path: POST /api/v1/sections/:section_id/assignments/:assignment_id/submissions/update_grades
func (s *SectionsService) UpdateGrades(ctx context.Context, sectionID, assignmentID int64, params map[string]interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/assignments/%d/submissions/update_grades", sectionID, assignmentID)

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, params, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateSubmissionsGrades bulk updates submission grades in a section
// Canvas path: POST /api/v1/sections/:section_id/submissions/update_grades
func (s *SectionsService) UpdateSubmissionsGrades(ctx context.Context, sectionID int64, params map[string]interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/submissions/update_grades", sectionID)

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, params, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// ClearSubmissionUnread clears unread state for a user's submissions in a section
// Canvas path: PUT /api/v1/sections/:section_id/submissions/:user_id/clear_unread
func (s *SectionsService) ClearSubmissionUnread(ctx context.Context, sectionID, userID int64) error {
	path := fmt.Sprintf("/api/v1/sections/%d/submissions/%d/clear_unread", sectionID, userID)

	return s.client.PutJSON(ctx, path, nil, nil)
}

// ListPeerReviews retrieves peer reviews for an assignment in a section
// Canvas path: GET /api/v1/sections/:section_id/assignments/:assignment_id/peer_reviews
func (s *SectionsService) ListPeerReviews(ctx context.Context, sectionID, assignmentID int64) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/assignments/%d/peer_reviews", sectionID, assignmentID)

	var result []map[string]interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, err
	}

	return result, nil
}
