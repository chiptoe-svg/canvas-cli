package options

// CoursesListOptions encapsulates all flags for courses list command
type CoursesListOptions struct {
	EnrollmentType  string
	EnrollmentState string
	Include         []string
	State           []string

	// Pagination
	PerPage int
}

// Validate performs option validation
func (o *CoursesListOptions) Validate() error {
	// No required fields for listing courses
	return nil
}

// CoursesGetOptions encapsulates all flags for courses get command
type CoursesGetOptions struct {
	CourseID int64
	Include  []string
}

// Validate performs option validation
func (o *CoursesGetOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return nil
}

// CoursesUpdateOptions encapsulates all flags for courses update command
type CoursesUpdateOptions struct {
	CourseID    int64
	Name        string
	CourseCode  string
	StartAt     string
	EndAt       string
	License     string
	IsPublic    *bool // Pointer to differentiate between not set and false
	DefaultView string
}

// Validate performs option validation
func (o *CoursesUpdateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return nil
}
