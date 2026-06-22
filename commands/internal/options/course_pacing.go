package options

import "fmt"

// CoursePacingGetOptions holds options for getting a course pace.
type CoursePacingGetOptions struct {
	CourseID int64
	PaceID   int64
}

// Validate validates the options.
func (o *CoursePacingGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.PaceID <= 0 {
		return fmt.Errorf("pace-id is required")
	}
	return nil
}

// CoursePacingCreateOptions holds options for creating a course pace.
type CoursePacingCreateOptions struct {
	CourseID        int64
	ExcludeWeekends bool
	HardEndDates    bool
	EndDate         string
}

// Validate validates the options.
func (o *CoursePacingCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	return nil
}

// CoursePacingUpdateOptions holds options for updating a course pace.
type CoursePacingUpdateOptions struct {
	CourseID        int64
	PaceID          int64
	ExcludeWeekends bool
	HardEndDates    bool
	EndDate         string
}

// Validate validates the options.
func (o *CoursePacingUpdateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.PaceID <= 0 {
		return fmt.Errorf("pace-id is required")
	}
	return nil
}

// CoursePacingDeleteOptions holds options for deleting a course pace.
type CoursePacingDeleteOptions struct {
	CourseID int64
	PaceID   int64
}

// Validate validates the options.
func (o *CoursePacingDeleteOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.PaceID <= 0 {
		return fmt.Errorf("pace-id is required")
	}
	return nil
}
