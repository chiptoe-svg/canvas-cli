package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// CoursesService handles course-related API calls
type CoursesService struct {
	client *Client
}

// NewCoursesService creates a new courses service
func NewCoursesService(client *Client) *CoursesService {
	return &CoursesService{client: client}
}

// ListCoursesOptions holds options for listing courses
type ListCoursesOptions struct {
	EnrollmentType  string   // student, teacher, ta, observer, designer
	EnrollmentState string   // active, invited_or_pending, completed
	Include         []string // needs_grading_count, syllabus_body, total_scores, term, etc.
	State           []string // unpublished, available, completed, deleted
	Page            int
	PerPage         int
}

// List retrieves all courses for the current user
func (s *CoursesService) List(ctx context.Context, opts *ListCoursesOptions) ([]Course, error) {
	path := "/api/v1/courses"

	if opts != nil {
		query := url.Values{}

		if opts.EnrollmentType != "" {
			query.Add("enrollment_type", opts.EnrollmentType)
		}

		if opts.EnrollmentState != "" {
			query.Add("enrollment_state", opts.EnrollmentState)
		}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		for _, state := range opts.State {
			query.Add("state[]", state)
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

	courses, err := GetAllPagesGeneric[Course](s.client, ctx, path)
	if err != nil {
		return nil, err
	}

	return NormalizeCourses(courses), nil
}

// Get retrieves a single course by ID
func (s *CoursesService) Get(ctx context.Context, courseID int64, include []string) (*Course, error) {
	path := fmt.Sprintf("/api/v1/courses/%d", courseID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var course Course
	if err := s.client.GetJSON(ctx, path, &course); err != nil {
		return nil, err
	}

	return NormalizeCourse(&course), nil
}

// CreateCourseParams holds parameters for creating a course
type CreateCourseParams struct {
	AccountID                        int64
	Name                             string
	CourseCode                       string
	StartAt                          string
	EndAt                            string
	License                          string
	IsPublic                         bool
	IsPublicToAuthUsers              bool
	PublicSyllabus                   bool
	PublicSyllabusToAuth             bool
	PublicDescription                string
	AllowStudentWikiEdits            bool
	AllowWikiComments                bool
	AllowStudentForumAttachments     bool
	OpenEnrollment                   bool
	SelfEnrollment                   bool
	RestrictEnrollmentsToCourseDates bool
	TermID                           int64
	SISCourseID                      string
	IntegrationID                    string
	HideFinalGrades                  bool
	ApplyAssignmentGroupWeights      bool
	TimeZone                         string
	Offer                            bool
	EnrollMe                         bool
	DefaultView                      string
	SyllabusBody                     string
	GradingStandardID                int64
	CourseFormat                     string
}

// addCourseStringField adds a string field to courseData if non-empty
func addCourseStringField(courseData map[string]interface{}, key, value string) {
	if value != "" {
		courseData[key] = value
	}
}

// addCourseBoolField adds a bool field to courseData if true
func addCourseBoolField(courseData map[string]interface{}, key string, value bool) {
	if value {
		courseData[key] = true
	}
}

// addCourseBoolPtrField adds a bool pointer field to courseData if not nil
func addCourseBoolPtrField(courseData map[string]interface{}, key string, value *bool) {
	if value != nil {
		courseData[key] = *value
	}
}

// addCourseInt64Field adds an int64 field to courseData if positive
func addCourseInt64Field(courseData map[string]interface{}, key string, value int64) {
	if value > 0 {
		courseData[key] = value
	}
}

// Create creates a new course
func (s *CoursesService) Create(ctx context.Context, params *CreateCourseParams) (*Course, error) {
	if params.AccountID == 0 {
		return nil, fmt.Errorf("account_id is required")
	}

	path := fmt.Sprintf("/api/v1/accounts/%d/courses", params.AccountID)

	body := map[string]interface{}{
		"course": make(map[string]interface{}),
	}

	courseData := body["course"].(map[string]interface{})

	// String fields
	addCourseStringField(courseData, "name", params.Name)
	addCourseStringField(courseData, "course_code", params.CourseCode)
	addCourseStringField(courseData, "start_at", params.StartAt)
	addCourseStringField(courseData, "end_at", params.EndAt)
	addCourseStringField(courseData, "license", params.License)
	addCourseStringField(courseData, "public_description", params.PublicDescription)
	addCourseStringField(courseData, "sis_course_id", params.SISCourseID)
	addCourseStringField(courseData, "integration_id", params.IntegrationID)
	addCourseStringField(courseData, "time_zone", params.TimeZone)
	addCourseStringField(courseData, "default_view", params.DefaultView)
	addCourseStringField(courseData, "syllabus_body", params.SyllabusBody)
	addCourseStringField(courseData, "course_format", params.CourseFormat)

	// Boolean fields
	addCourseBoolField(courseData, "is_public", params.IsPublic)
	addCourseBoolField(courseData, "is_public_to_auth_users", params.IsPublicToAuthUsers)
	addCourseBoolField(courseData, "public_syllabus", params.PublicSyllabus)
	addCourseBoolField(courseData, "public_syllabus_to_auth", params.PublicSyllabusToAuth)
	addCourseBoolField(courseData, "allow_student_wiki_edits", params.AllowStudentWikiEdits)
	addCourseBoolField(courseData, "allow_wiki_comments", params.AllowWikiComments)
	addCourseBoolField(courseData, "allow_student_forum_attachments", params.AllowStudentForumAttachments)
	addCourseBoolField(courseData, "open_enrollment", params.OpenEnrollment)
	addCourseBoolField(courseData, "self_enrollment", params.SelfEnrollment)
	addCourseBoolField(courseData, "restrict_enrollments_to_course_dates", params.RestrictEnrollmentsToCourseDates)
	addCourseBoolField(courseData, "hide_final_grades", params.HideFinalGrades)
	addCourseBoolField(courseData, "apply_assignment_group_weights", params.ApplyAssignmentGroupWeights)
	addCourseBoolField(courseData, "offer", params.Offer)
	addCourseBoolField(courseData, "enroll_me", params.EnrollMe)

	// Integer fields
	addCourseInt64Field(courseData, "term_id", params.TermID)
	addCourseInt64Field(courseData, "grading_standard_id", params.GradingStandardID)

	var course Course
	if err := s.client.PostJSON(ctx, path, body, &course); err != nil {
		return nil, err
	}

	if course.ID == 0 {
		return nil, fmt.Errorf("course creation failed: Canvas returned no course ID (possible API error or maintenance mode)")
	}

	return NormalizeCourse(&course), nil
}

// UpdateCourseParams holds parameters for updating a course
type UpdateCourseParams struct {
	Name                             string
	CourseCode                       string
	StartAt                          string
	EndAt                            string
	License                          string
	IsPublic                         *bool
	IsPublicToAuthUsers              *bool
	PublicSyllabus                   *bool
	PublicSyllabusToAuth             *bool
	PublicDescription                string
	AllowStudentWikiEdits            *bool
	AllowWikiComments                *bool
	AllowStudentForumAttachments     *bool
	OpenEnrollment                   *bool
	SelfEnrollment                   *bool
	RestrictEnrollmentsToCourseDates *bool
	HideFinalGrades                  *bool
	ApplyAssignmentGroupWeights      *bool
	TimeZone                         string
	DefaultView                      string
	SyllabusBody                     string
	GradingStandardID                int64
	CourseFormat                     string
	ImageID                          int64
	ImageURL                         string
	RemoveImage                      bool
}

// Update updates an existing course
func (s *CoursesService) Update(ctx context.Context, courseID int64, params *UpdateCourseParams) (*Course, error) {
	path := fmt.Sprintf("/api/v1/courses/%d", courseID)

	body := map[string]interface{}{
		"course": make(map[string]interface{}),
	}

	courseData := body["course"].(map[string]interface{})

	// String fields
	addCourseStringField(courseData, "name", params.Name)
	addCourseStringField(courseData, "course_code", params.CourseCode)
	addCourseStringField(courseData, "start_at", params.StartAt)
	addCourseStringField(courseData, "end_at", params.EndAt)
	addCourseStringField(courseData, "license", params.License)
	addCourseStringField(courseData, "public_description", params.PublicDescription)
	addCourseStringField(courseData, "time_zone", params.TimeZone)
	addCourseStringField(courseData, "default_view", params.DefaultView)
	addCourseStringField(courseData, "syllabus_body", params.SyllabusBody)
	addCourseStringField(courseData, "course_format", params.CourseFormat)
	addCourseStringField(courseData, "image_url", params.ImageURL)

	// Boolean pointer fields (allow explicit false)
	addCourseBoolPtrField(courseData, "is_public", params.IsPublic)
	addCourseBoolPtrField(courseData, "is_public_to_auth_users", params.IsPublicToAuthUsers)
	addCourseBoolPtrField(courseData, "public_syllabus", params.PublicSyllabus)
	addCourseBoolPtrField(courseData, "public_syllabus_to_auth", params.PublicSyllabusToAuth)
	addCourseBoolPtrField(courseData, "allow_student_wiki_edits", params.AllowStudentWikiEdits)
	addCourseBoolPtrField(courseData, "allow_wiki_comments", params.AllowWikiComments)
	addCourseBoolPtrField(courseData, "allow_student_forum_attachments", params.AllowStudentForumAttachments)
	addCourseBoolPtrField(courseData, "open_enrollment", params.OpenEnrollment)
	addCourseBoolPtrField(courseData, "self_enrollment", params.SelfEnrollment)
	addCourseBoolPtrField(courseData, "restrict_enrollments_to_course_dates", params.RestrictEnrollmentsToCourseDates)
	addCourseBoolPtrField(courseData, "hide_final_grades", params.HideFinalGrades)
	addCourseBoolPtrField(courseData, "apply_assignment_group_weights", params.ApplyAssignmentGroupWeights)

	// Boolean field
	addCourseBoolField(courseData, "remove_image", params.RemoveImage)

	// Integer fields
	addCourseInt64Field(courseData, "grading_standard_id", params.GradingStandardID)
	addCourseInt64Field(courseData, "image_id", params.ImageID)

	var course Course
	if err := s.client.PutJSON(ctx, path, body, &course); err != nil {
		return nil, err
	}

	if course.ID == 0 {
		return nil, fmt.Errorf("course update failed: Canvas returned no course ID (possible API error or maintenance mode)")
	}

	return NormalizeCourse(&course), nil
}

// Delete deletes a course (sets to deleted state)
func (s *CoursesService) Delete(ctx context.Context, courseID int64, event string) error {
	path := fmt.Sprintf("/api/v1/courses/%d", courseID)

	if event != "" {
		path += "?event=" + url.QueryEscape(event)
	}

	_, err := s.client.Delete(ctx, path)
	return err
}

// GetActivityStream retrieves the activity stream for a course
// Canvas path: GET /api/v1/courses/:course_id/activity_stream
func (s *CoursesService) GetActivityStream(ctx context.Context, courseID int64) ([]ActivityStreamItem, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/activity_stream", courseID)

	var items []ActivityStreamItem
	if err := s.client.GetAllPages(ctx, path, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetActivityStreamSummary retrieves a summary of the activity stream for a course
// Canvas path: GET /api/v1/courses/:course_id/activity_stream/summary
func (s *CoursesService) GetActivityStreamSummary(ctx context.Context, courseID int64) ([]ActivityStreamSummary, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/activity_stream/summary", courseID)

	var items []ActivityStreamSummary
	if err := s.client.GetJSON(ctx, path, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetStudents retrieves all students enrolled in a course
// Canvas path: GET /api/v1/courses/:course_id/students
func (s *CoursesService) GetStudents(ctx context.Context, courseID int64) ([]User, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/students", courseID)

	var users []User
	if err := s.client.GetAllPages(ctx, path, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// ListCourseUsersOptions holds options for listing users in a course
type ListCourseUsersOptions struct {
	SearchTerm      string
	Include         []string // enrollments, locked, avatar_url, test_student, bio, custom_links, current_grading_period_scores, uuid
	Sort            string
	EnrollmentType  []string // student, teacher, ta, observer, designer
	EnrollmentRole  []string
	EnrollmentState []string
	Page            int
	PerPage         int
}

// GetUser retrieves a single user in a course
// Canvas path: GET /api/v1/courses/:course_id/users/:id
func (s *CoursesService) GetUser(ctx context.Context, courseID, userID int64, include []string) (*User, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/users/%d", courseID, userID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var user User
	if err := s.client.GetJSON(ctx, path, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// ListUsers retrieves all users enrolled in a course
// Canvas path: GET /api/v1/courses/:course_id/users
// Already handled by ListEnrollments in enrollment, but exposed here for convenience
func (s *CoursesService) ListUsers(ctx context.Context, courseID int64, opts *ListCourseUsersOptions) ([]User, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/users", courseID)

	if opts != nil {
		query := url.Values{}
		if opts.SearchTerm != "" {
			query.Add("search_term", opts.SearchTerm)
		}
		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}
		if opts.Sort != "" {
			query.Add("sort", opts.Sort)
		}
		for _, et := range opts.EnrollmentType {
			query.Add("enrollment_type[]", et)
		}
		for _, er := range opts.EnrollmentRole {
			query.Add("enrollment_role[]", er)
		}
		for _, es := range opts.EnrollmentState {
			query.Add("enrollment_state[]", es)
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

	var users []User
	if err := s.client.GetAllPages(ctx, path, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserProgress retrieves a user's progress in a course
// Canvas path: GET /api/v1/courses/:course_id/users/:user_id/progress
func (s *CoursesService) GetUserProgress(ctx context.Context, courseID, userID int64) (*CourseProgress, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/users/%d/progress", courseID, userID)

	var progress CourseProgress
	if err := s.client.GetJSON(ctx, path, &progress); err != nil {
		return nil, err
	}

	return &progress, nil
}

// SearchUsers searches for users in a course
// Canvas path: GET /api/v1/courses/:course_id/search_users
func (s *CoursesService) SearchUsers(ctx context.Context, courseID int64, searchTerm string) ([]User, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/search_users", courseID)

	if searchTerm != "" {
		path += "?search_term=" + url.QueryEscape(searchTerm)
	}

	var users []User
	if err := s.client.GetAllPages(ctx, path, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetStudentViewStudent retrieves the test student for a course
// Canvas path: GET /api/v1/courses/:course_id/student_view_student
func (s *CoursesService) GetStudentViewStudent(ctx context.Context, courseID int64) (*User, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/student_view_student", courseID)

	var user User
	if err := s.client.GetJSON(ctx, path, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// ResetContent resets a course back to default blank state
// Canvas path: POST /api/v1/courses/:course_id/reset_content
func (s *CoursesService) ResetContent(ctx context.Context, courseID int64) (*Course, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/reset_content", courseID)

	var course Course
	if err := s.client.PostJSON(ctx, path, nil, &course); err != nil {
		return nil, err
	}

	return &course, nil
}

// CourseProgress represents a user's progress in a course
type CourseProgress struct {
	RequirementCount          int     `json:"requirement_count"`
	RequirementCompletedCount int     `json:"requirement_completed_count"`
	NextRequirementURL        string  `json:"next_requirement_url,omitempty"`
	CompletedAt               *string `json:"completed_at,omitempty"`
}

// ContentShareUser represents a user available for content sharing
type ContentShareUser struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_image_url,omitempty"`
}

// GetContentShareUsers retrieves users available for content sharing in a course
// Canvas path: GET /api/v1/courses/:course_id/content_share_users
func (s *CoursesService) GetContentShareUsers(ctx context.Context, courseID int64, searchTerm string) ([]ContentShareUser, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/content_share_users", courseID)
	if searchTerm != "" {
		path += "?search_term=" + url.QueryEscape(searchTerm)
	}

	var users []ContentShareUser
	if err := s.client.GetAllPages(ctx, path, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetPotentialCollaborators retrieves potential collaborators for a course
// Canvas path: GET /api/v1/courses/:course_id/potential_collaborators
func (s *CoursesService) GetPotentialCollaborators(ctx context.Context, courseID int64) ([]User, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/potential_collaborators", courseID)

	var users []User
	if err := s.client.GetAllPages(ctx, path, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetCSPSettings retrieves CSP settings for a course
// Canvas path: GET /api/v1/courses/:course_id/csp_settings
func (s *CoursesService) GetCSPSettings(ctx context.Context, courseID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/csp_settings", courseID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateCSPSettings updates CSP settings for a course
// Canvas path: PUT /api/v1/courses/:course_id/csp_settings
func (s *CoursesService) UpdateCSPSettings(ctx context.Context, courseID int64, params map[string]interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/csp_settings", courseID)

	var result map[string]interface{}
	if err := s.client.PutJSON(ctx, path, params, &result); err != nil {
		return nil, err
	}

	return result, nil
}
