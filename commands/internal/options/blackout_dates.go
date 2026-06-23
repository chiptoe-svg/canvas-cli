package options

import "fmt"

// BlackoutDatesListOptions holds options for listing blackout dates.
type BlackoutDatesListOptions struct {
	CourseID int64
}

// Validate validates the options.
func (o *BlackoutDatesListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	return nil
}

// BlackoutDatesGetOptions holds options for getting a blackout date.
type BlackoutDatesGetOptions struct {
	CourseID int64
	ID       int64
}

// Validate validates the options.
func (o *BlackoutDatesGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.ID <= 0 {
		return fmt.Errorf("blackout-date-id is required")
	}
	return nil
}

// BlackoutDatesCreateOptions holds options for creating a blackout date.
type BlackoutDatesCreateOptions struct {
	CourseID   int64
	StartDate  string
	EndDate    string
	EventTitle string
}

// Validate validates the options.
func (o *BlackoutDatesCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.StartDate == "" {
		return fmt.Errorf("--start-date is required")
	}
	if o.EndDate == "" {
		return fmt.Errorf("--end-date is required")
	}
	if o.EventTitle == "" {
		return fmt.Errorf("--title is required")
	}
	return nil
}

// BlackoutDatesUpdateOptions holds options for updating a blackout date.
type BlackoutDatesUpdateOptions struct {
	CourseID   int64
	ID         int64
	StartDate  string
	EndDate    string
	EventTitle string
}

// Validate validates the options.
func (o *BlackoutDatesUpdateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.ID <= 0 {
		return fmt.Errorf("blackout-date-id is required")
	}
	return nil
}

// BlackoutDatesDeleteOptions holds options for deleting a blackout date.
type BlackoutDatesDeleteOptions struct {
	CourseID int64
	ID       int64
}

// Validate validates the options.
func (o *BlackoutDatesDeleteOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.ID <= 0 {
		return fmt.Errorf("blackout-date-id is required")
	}
	return nil
}
