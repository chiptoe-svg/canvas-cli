package api

import (
	"context"
	"fmt"
	"net/url"
)

// UserProfile represents a user's profile data
type UserProfile struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	ShortName       string `json:"short_name"`
	SortableName    string `json:"sortable_name"`
	Title           string `json:"title"`
	Bio             string `json:"bio"`
	PrimaryEmail    string `json:"primary_email"`
	LoginID         string `json:"login_id"`
	SISUserID       string `json:"sis_user_id"`
	LTIUserID       string `json:"lti_user_id"`
	AvatarURL       string `json:"avatar_url"`
	TimeZone        string `json:"time_zone"`
	Locale          string `json:"locale"`
	EffectiveLocale string `json:"effective_locale"`
	Pronouns        string `json:"pronouns"`
}

// UserAvatar represents an available avatar option for a user
type UserAvatar struct {
	Type        string `json:"type"`
	URL         string `json:"url"`
	Token       string `json:"token"`
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	ContentType string `json:"content_type"`
	FileName    string `json:"filename"`
	Size        int64  `json:"size"`
}

// UserSettings represents configurable user-level settings
type UserSettings struct {
	ManualMarkAsRead           bool `json:"manual_mark_as_read"`
	CollapseGlobalNav          bool `json:"collapse_global_nav"`
	HidesDashcardColorOverlays bool `json:"hide_dashcard_color_overlays"`
}

// UpdateUserSettingsParams holds settings to update for a user
type UpdateUserSettingsParams struct {
	ManualMarkAsRead           *bool `json:"manual_mark_as_read,omitempty"`
	CollapseGlobalNav          *bool `json:"collapse_global_nav,omitempty"`
	HidesDashcardColorOverlays *bool `json:"hide_dashcard_color_overlays,omitempty"`
}

// PageView represents a single entry in a user's page view history
type PageView struct {
	ID                 string  `json:"id"`
	AppName            string  `json:"app_name"`
	URL                string  `json:"url"`
	ContextType        string  `json:"context_type"`
	AssetType          string  `json:"asset_type"`
	Controller         string  `json:"controller"`
	Action             string  `json:"action"`
	Contributed        bool    `json:"contributed"`
	InteractionSeconds float64 `json:"interaction_seconds"`
	CreatedAt          string  `json:"created_at"`
	UserAgent          string  `json:"user_agent"`
	Participated       bool    `json:"participated"`
	HTTPMethod         string  `json:"http_method"`
	RemoteIP           string  `json:"remote_ip"`
}

// ActivityStreamItem represents an item in a user's activity stream
type ActivityStreamItem struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Type      string `json:"type"`
	ReadState bool   `json:"read_state"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	HTMLUrl   string `json:"html_url"`
	CourseID  int64  `json:"course_id"`
	GroupID   int64  `json:"group_id"`
}

// ActivityStreamSummary represents an aggregate count of activity stream items by type
type ActivityStreamSummary struct {
	Type                 string `json:"type"`
	Count                int    `json:"count"`
	UnreadCount          int    `json:"unread_count"`
	NotificationCategory string `json:"notification_category"`
}

// TodoItem represents an item on the current user's to-do list
type TodoItem struct {
	Type              string      `json:"type"`
	Assignment        interface{} `json:"assignment"`
	ContextType       string      `json:"context_type"`
	CourseID          int64       `json:"course_id"`
	GroupID           int64       `json:"group_id"`
	HTMLUrl           string      `json:"html_url"`
	Ignore            string      `json:"ignore"`
	IgnorePermanently string      `json:"ignore_permanently"`
}

// TodoItemCount holds the count of outstanding to-do items for a user
type TodoItemCount struct {
	NeedsGradingCount            int `json:"needs_grading_count"`
	AssignmentsNeedingSubmitting int `json:"assignments_needing_submitting"`
}

// UserLogin represents a login pseudonym for a user
type UserLogin struct {
	ID                       int64  `json:"id"`
	UserID                   int64  `json:"user_id"`
	AccountID                int64  `json:"account_id"`
	WorkflowState            string `json:"workflow_state"`
	UniqueID                 string `json:"unique_id"`
	CreatedAt                string `json:"created_at"`
	SISUserID                string `json:"sis_user_id"`
	IntegrationID            string `json:"integration_id"`
	SISImportID              int64  `json:"sis_import_id"`
	AuthenticationProviderID int64  `json:"authentication_provider_id"`
}

// UserColor represents a single color customization keyed by asset string
type UserColor struct {
	HexCode string `json:"hexcode"`
}

// UserColors holds all custom color settings for a user
type UserColors struct {
	CustomColors map[string]string `json:"custom_colors"`
}

// DashboardPositions holds the dashboard card position overrides for a user
type DashboardPositions struct {
	DashboardPositions map[string]int `json:"dashboard_positions"`
}

// Tab represents a navigation tab for a user context
type Tab struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Type     string `json:"type"`
	Hidden   bool   `json:"hidden"`
	Position int    `json:"position"`
	HTMLUrl  string `json:"html_url"`
}

// MissingSubmission represents an assignment with a missing submission for a user
type MissingSubmission struct {
	ID              int64    `json:"id"`
	Name            string   `json:"name"`
	PointsPossible  float64  `json:"points_possible"`
	DueAt           string   `json:"due_at"`
	HTMLURL         string   `json:"html_url"`
	CourseID        int64    `json:"course_id"`
	SubmissionTypes []string `json:"submission_types"`
}

// GetProfile retrieves a user's profile
func (s *UsersService) GetProfile(ctx context.Context, userID int64) (*UserProfile, error) {
	path := fmt.Sprintf("/api/v1/users/%d/profile", userID)
	var result UserProfile
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting profile for user %d: %w", userID, err)
	}
	return &result, nil
}

// GetAvatars retrieves available avatar options for a user
func (s *UsersService) GetAvatars(ctx context.Context, userID int64) ([]UserAvatar, error) {
	path := fmt.Sprintf("/api/v1/users/%d/avatars", userID)
	var result []UserAvatar
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting avatars for user %d: %w", userID, err)
	}
	return result, nil
}

// GetSettings retrieves settings for a user
func (s *UsersService) GetSettings(ctx context.Context, userID int64) (*UserSettings, error) {
	path := fmt.Sprintf("/api/v1/users/%d/settings", userID)
	var result UserSettings
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting settings for user %d: %w", userID, err)
	}
	return &result, nil
}

// UpdateSettings updates settings for a user
func (s *UsersService) UpdateSettings(ctx context.Context, id int64, params UpdateUserSettingsParams) (*UserSettings, error) {
	path := fmt.Sprintf("/api/v1/users/%d/settings", id)
	var result UserSettings
	if err := s.client.PutJSON(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("updating settings for user %d: %w", id, err)
	}
	return &result, nil
}

// GetHistory retrieves a user's recently visited items
func (s *UsersService) GetHistory(ctx context.Context, userID int64) ([]interface{}, error) {
	path := fmt.Sprintf("/api/v1/users/%d/history", userID)
	var result []interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting history for user %d: %w", userID, err)
	}
	return result, nil
}

// GetMissingSubmissions retrieves assignments with missing submissions for a user
func (s *UsersService) GetMissingSubmissions(ctx context.Context, userID int64) ([]MissingSubmission, error) {
	path := fmt.Sprintf("/api/v1/users/%d/missing_submissions", userID)
	var result []MissingSubmission
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting missing submissions for user %d: %w", userID, err)
	}
	return result, nil
}

// GetPageViews retrieves the page view history for a user
func (s *UsersService) GetPageViews(ctx context.Context, userID int64) ([]PageView, error) {
	path := fmt.Sprintf("/api/v1/users/%d/page_views", userID)
	var result []PageView
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting page views for user %d: %w", userID, err)
	}
	return result, nil
}

// GetTabs retrieves the navigation tabs available to a user
func (s *UsersService) GetTabs(ctx context.Context, userID int64) ([]Tab, error) {
	path := fmt.Sprintf("/api/v1/users/%d/tabs", userID)
	var result []Tab
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting tabs for user %d: %w", userID, err)
	}
	return result, nil
}

// ListLogins retrieves login pseudonyms for a user
func (s *UsersService) ListLogins(ctx context.Context, userID int64) ([]UserLogin, error) {
	path := fmt.Sprintf("/api/v1/users/%d/logins", userID)
	var result []UserLogin
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing logins for user %d: %w", userID, err)
	}
	return result, nil
}

// DeleteLogin deletes a login pseudonym for a user
func (s *UsersService) DeleteLogin(ctx context.Context, userID, loginID int64) (*UserLogin, error) {
	path := fmt.Sprintf("/api/v1/users/%d/logins/%d", userID, loginID)
	var result UserLogin
	if err := s.client.DeleteJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("deleting login %d for user %d: %w", loginID, userID, err)
	}
	return &result, nil
}

// ListUserCourses retrieves courses for a user
func (s *UsersService) ListUserCourses(ctx context.Context, userID int64) ([]Course, error) {
	path := fmt.Sprintf("/api/v1/users/%d/courses", userID)
	var result []Course
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing courses for user %d: %w", userID, err)
	}
	return result, nil
}

// GetColors retrieves all custom color settings for a user
func (s *UsersService) GetColors(ctx context.Context, id int64) (*UserColors, error) {
	path := fmt.Sprintf("/api/v1/users/%d/colors", id)
	var result UserColors
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting colors for user %d: %w", id, err)
	}
	return &result, nil
}

// GetColor retrieves a specific color customization for a user by asset string
func (s *UsersService) GetColor(ctx context.Context, id int64, assetString string) (*UserColor, error) {
	path := fmt.Sprintf("/api/v1/users/%d/colors/%s", id, url.PathEscape(assetString))
	var result UserColor
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting color %s for user %d: %w", assetString, id, err)
	}
	return &result, nil
}

// SetColor sets a color customization for a user by asset string
func (s *UsersService) SetColor(ctx context.Context, id int64, assetString, hexCode string) (*UserColor, error) {
	path := fmt.Sprintf("/api/v1/users/%d/colors/%s", id, url.PathEscape(assetString))
	body := map[string]string{"hexcode": hexCode}
	var result UserColor
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("setting color %s for user %d: %w", assetString, id, err)
	}
	return &result, nil
}

// GetDashboardPositions retrieves dashboard card position overrides for a user
func (s *UsersService) GetDashboardPositions(ctx context.Context, id int64) (*DashboardPositions, error) {
	path := fmt.Sprintf("/api/v1/users/%d/dashboard_positions", id)
	var result DashboardPositions
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting dashboard positions for user %d: %w", id, err)
	}
	return &result, nil
}

// SetDashboardPositions sets dashboard card positions for a user
func (s *UsersService) SetDashboardPositions(ctx context.Context, id int64, positions map[string]int) (*DashboardPositions, error) {
	path := fmt.Sprintf("/api/v1/users/%d/dashboard_positions", id)
	body := DashboardPositions{DashboardPositions: positions}
	var result DashboardPositions
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("setting dashboard positions for user %d: %w", id, err)
	}
	return &result, nil
}

// MergeInto merges a user into a destination user within the same account
func (s *UsersService) MergeInto(ctx context.Context, id, destinationUserID int64) (*User, error) {
	path := fmt.Sprintf("/api/v1/users/%d/merge_into/%d", id, destinationUserID)
	var result User
	if err := s.client.PutJSON(ctx, path, nil, &result); err != nil {
		return nil, fmt.Errorf("merging user %d into %d: %w", id, destinationUserID, err)
	}
	return &result, nil
}

// MergeIntoAccount merges a user into a destination user in a specific account
func (s *UsersService) MergeIntoAccount(ctx context.Context, id, destinationAccountID, destinationUserID int64) (*User, error) {
	path := fmt.Sprintf("/api/v1/users/%d/merge_into/accounts/%d/users/%d", id, destinationAccountID, destinationUserID)
	var result User
	if err := s.client.PutJSON(ctx, path, nil, &result); err != nil {
		return nil, fmt.Errorf("merging user %d into account %d user %d: %w", id, destinationAccountID, destinationUserID, err)
	}
	return &result, nil
}

// Split splits a merged user back into the source users
func (s *UsersService) Split(ctx context.Context, id int64) ([]User, error) {
	path := fmt.Sprintf("/api/v1/users/%d/split", id)
	var result []User
	if err := s.client.PostJSON(ctx, path, nil, &result); err != nil {
		return nil, fmt.Errorf("splitting user %d: %w", id, err)
	}
	return result, nil
}

// GetActivityStream retrieves the activity stream for the current user
func (s *UsersService) GetActivityStream(ctx context.Context) ([]ActivityStreamItem, error) {
	path := "/api/v1/users/self/activity_stream"
	var result []ActivityStreamItem
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting activity stream: %w", err)
	}
	return result, nil
}

// DeleteActivityStream deletes all activity stream items for the current user
func (s *UsersService) DeleteActivityStream(ctx context.Context) error {
	path := "/api/v1/users/self/activity_stream"
	_, err := s.client.Delete(ctx, path)
	return err
}

// DeleteActivityStreamItem deletes a specific activity stream item
func (s *UsersService) DeleteActivityStreamItem(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/api/v1/users/self/activity_stream/%d", id)
	_, err := s.client.Delete(ctx, path)
	return err
}

// GetActivityStreamSummary retrieves aggregate counts for the current user's activity stream
func (s *UsersService) GetActivityStreamSummary(ctx context.Context) ([]ActivityStreamSummary, error) {
	path := "/api/v1/users/self/activity_stream/summary"
	var result []ActivityStreamSummary
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting activity stream summary: %w", err)
	}
	return result, nil
}

// GetTodo retrieves outstanding to-do items for the current user
func (s *UsersService) GetTodo(ctx context.Context) ([]TodoItem, error) {
	path := "/api/v1/users/self/todo"
	var result []TodoItem
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting todo items: %w", err)
	}
	return result, nil
}

// GetTodoItemCount retrieves the count of outstanding to-do items for the current user
func (s *UsersService) GetTodoItemCount(ctx context.Context) (*TodoItemCount, error) {
	path := "/api/v1/users/self/todo_item_count"
	var result TodoItemCount
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting todo item count: %w", err)
	}
	return &result, nil
}

// GetUpcomingEvents retrieves upcoming calendar events for the current user
func (s *UsersService) GetUpcomingEvents(ctx context.Context) ([]interface{}, error) {
	path := "/api/v1/users/self/upcoming_events"
	var result []interface{}
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting upcoming events: %w", err)
	}
	return result, nil
}
