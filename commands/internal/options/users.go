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
