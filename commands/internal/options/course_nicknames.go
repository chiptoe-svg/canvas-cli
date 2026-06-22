package options

import "fmt"

// CourseNicknamesListOptions contains options for listing course nicknames
type CourseNicknamesListOptions struct{}

// Validate validates the options
func (o *CourseNicknamesListOptions) Validate() error { return nil }

// CourseNicknamesGetOptions contains options for getting a course nickname
type CourseNicknamesGetOptions struct {
	CourseID int64
}

// Validate validates the options
func (o *CourseNicknamesGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	return nil
}

// CourseNicknamesSetOptions contains options for setting a course nickname
type CourseNicknamesSetOptions struct {
	CourseID int64
	Nickname string
}

// Validate validates the options
func (o *CourseNicknamesSetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.Nickname == "" {
		return fmt.Errorf("nickname is required")
	}
	return nil
}

// CourseNicknamesDeleteOptions contains options for deleting a course nickname
type CourseNicknamesDeleteOptions struct {
	CourseID int64
}

// Validate validates the options
func (o *CourseNicknamesDeleteOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	return nil
}
