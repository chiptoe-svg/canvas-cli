package options

// ConferencesListOptions holds options for listing conferences.
type ConferencesListOptions struct {
	CourseID int64
	GroupID  int64
	State    string
	PerPage  int
}

// Validate validates the options.
func (o *ConferencesListOptions) Validate() error {
	return nil
}
