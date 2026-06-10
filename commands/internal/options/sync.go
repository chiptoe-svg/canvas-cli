package options

// SyncAssignmentsOptions contains options for the sync assignments command
type SyncAssignmentsOptions struct {
	Interactive bool
}

// Validate validates the options.
func (o *SyncAssignmentsOptions) Validate() error {
	return nil
}

// SyncCourseOptions contains options for the sync course command
type SyncCourseOptions struct {
	Interactive bool
}

// Validate validates the options.
func (o *SyncCourseOptions) Validate() error {
	return nil
}
