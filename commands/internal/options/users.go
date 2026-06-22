package options

import "fmt"

// UsersListOptions contains options for listing users
type UsersListOptions struct {
	AccountID       int64
	CourseID        int64
	SearchTerm      string
	EnrollmentType  string
	EnrollmentState string
	Include         []string
}

// Validate validates the options
func (o *UsersListOptions) Validate() error {
	if o.AccountID > 0 && o.CourseID > 0 {
		return fmt.Errorf("can only specify one of --account-id or --course-id")
	}
	return nil
}

// UsersGetOptions contains options for getting a user
type UsersGetOptions struct {
	UserID  int64
	Include []string
}

// Validate validates the options
func (o *UsersGetOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// UsersMeOptions contains options for getting current user
type UsersMeOptions struct {
	// No options needed
}

// Validate validates the options
func (o *UsersMeOptions) Validate() error {
	return nil
}

// UsersSearchOptions contains options for searching users
type UsersSearchOptions struct {
	SearchTerm string
}

// Validate validates the options
func (o *UsersSearchOptions) Validate() error {
	if o.SearchTerm == "" {
		return fmt.Errorf("search-term is required")
	}
	return nil
}

// UsersCreateOptions contains options for creating a user
type UsersCreateOptions struct {
	AccountID        int64
	Name             string
	ShortName        string
	SortableName     string
	Email            string
	LoginID          string
	Password         string
	SISUserID        string
	TimeZone         string
	Locale           string
	SkipRegistration bool
	SkipConfirmation bool
	JSONFile         string
	Stdin            bool
}

// Validate validates the options
func (o *UsersCreateOptions) Validate() error {
	if o.AccountID <= 0 {
		return fmt.Errorf("account-id is required and must be greater than 0")
	}
	return nil
}

// UsersUpdateOptions contains options for updating a user
type UsersUpdateOptions struct {
	UserID       int64
	Name         string
	ShortName    string
	SortableName string
	Email        string
	TimeZone     string
	Locale       string
	JSONFile     string
	Stdin        bool
}

// Validate validates the options
func (o *UsersUpdateOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// UsersProfileOptions contains options for getting a user profile
type UsersProfileOptions struct {
	UserID int64
}

// Validate validates the options
func (o *UsersProfileOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// UsersSettingsOptions contains options for getting user settings
type UsersSettingsOptions struct {
	UserID int64
}

// Validate validates the options
func (o *UsersSettingsOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// UsersUpdateSettingsOptions contains options for updating user settings
type UsersUpdateSettingsOptions struct {
	UserID            int64
	ManualMarkAsRead  bool
	CollapseGlobalNav bool
}

// Validate validates the options
func (o *UsersUpdateSettingsOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// UsersPageViewsOptions contains options for getting page views
type UsersPageViewsOptions struct {
	UserID int64
}

// Validate validates the options
func (o *UsersPageViewsOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// UsersLoginsOptions contains options for listing logins
type UsersLoginsOptions struct {
	UserID int64
}

// Validate validates the options
func (o *UsersLoginsOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// UsersCoursesOptions contains options for listing user courses
type UsersCoursesOptions struct {
	UserID int64
}

// Validate validates the options
func (o *UsersCoursesOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// UsersMissingSubmissionsOptions contains options for getting missing submissions
type UsersMissingSubmissionsOptions struct {
	UserID int64
}

// Validate validates the options
func (o *UsersMissingSubmissionsOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// UsersActivityStreamOptions contains options for getting the activity stream
type UsersActivityStreamOptions struct{}

// Validate validates the options
func (o *UsersActivityStreamOptions) Validate() error { return nil }

// UsersTodoOptions contains options for getting todo items
type UsersTodoOptions struct{}

// Validate validates the options
func (o *UsersTodoOptions) Validate() error { return nil }

// UsersUpcomingEventsOptions contains options for getting upcoming events
type UsersUpcomingEventsOptions struct{}

// Validate validates the options
func (o *UsersUpcomingEventsOptions) Validate() error { return nil }

// UsersMergeOptions contains options for merging users
type UsersMergeOptions struct {
	UserID            int64
	DestinationUserID int64
}

// Validate validates the options
func (o *UsersMergeOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	if o.DestinationUserID <= 0 {
		return fmt.Errorf("destination-user-id is required and must be greater than 0")
	}
	return nil
}

// UsersSplitOptions contains options for splitting a user
type UsersSplitOptions struct {
	UserID int64
}

// Validate validates the options
func (o *UsersSplitOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}
