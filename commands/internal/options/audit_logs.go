package options

// AuditLogsListOptions holds common options for listing audit log events.
type AuditLogsListOptions struct {
	// Context selection — exactly one should be set
	AccountID    int64
	LoginID      int64
	UserID       int64
	CourseID     int64
	AssignmentID int64
	GraderID     int64
	StudentID    int64

	// Audit type selection
	AuditType string // authentication, course, grade_change

	// Filters
	StartTime string
	EndTime   string
	PerPage   int
}

// Validate validates the options.
func (o *AuditLogsListOptions) Validate() error {
	return nil
}
