package options

import "fmt"

// GradingPeriodsListOptions holds options for listing grading periods.
type GradingPeriodsListOptions struct {
	CourseID int64
}

// Validate validates the options.
func (o *GradingPeriodsListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	return nil
}

// GradingPeriodsGetOptions holds options for getting a grading period.
type GradingPeriodsGetOptions struct {
	CourseID int64
	ID       int64
}

// Validate validates the options.
func (o *GradingPeriodsGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.ID <= 0 {
		return fmt.Errorf("grading-period-id is required")
	}
	return nil
}

// GradingPeriodsUpdateOptions holds options for updating a grading period.
type GradingPeriodsUpdateOptions struct {
	CourseID  int64
	ID        int64
	Title     string
	StartDate string
	EndDate   string
	CloseDate string
	Weight    float64
}

// Validate validates the options.
func (o *GradingPeriodsUpdateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.ID <= 0 {
		return fmt.Errorf("grading-period-id is required")
	}
	return nil
}

// GradingPeriodsDeleteOptions holds options for deleting a grading period.
type GradingPeriodsDeleteOptions struct {
	CourseID int64
	ID       int64
	Force    bool
}

// Validate validates the options.
func (o *GradingPeriodsDeleteOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.ID <= 0 {
		return fmt.Errorf("grading-period-id is required")
	}
	return nil
}

// GradingStandardsListOptions holds options for listing grading standards.
type GradingStandardsListOptions struct {
	CourseID int64
}

// Validate validates the options.
func (o *GradingStandardsListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	return nil
}

// GradingStandardsGetOptions holds options for getting a grading standard.
type GradingStandardsGetOptions struct {
	CourseID   int64
	StandardID int64
}

// Validate validates the options.
func (o *GradingStandardsGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.StandardID <= 0 {
		return fmt.Errorf("standard-id is required")
	}
	return nil
}

// GradingStandardsCreateOptions holds options for creating a grading standard.
type GradingStandardsCreateOptions struct {
	CourseID int64
	Title    string
}

// Validate validates the options.
func (o *GradingStandardsCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.Title == "" {
		return fmt.Errorf("--title is required")
	}
	return nil
}

// GradingStandardsDeleteOptions holds options for deleting a grading standard.
type GradingStandardsDeleteOptions struct {
	CourseID   int64
	StandardID int64
}

// Validate validates the options.
func (o *GradingStandardsDeleteOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.StandardID <= 0 {
		return fmt.Errorf("standard-id is required")
	}
	return nil
}
