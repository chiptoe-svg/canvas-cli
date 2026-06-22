package options

// CollaborationsListOptions holds options for listing collaborations.
type CollaborationsListOptions struct {
	CourseID        int64
	GroupID         int64
	CollaborationID int64
	PerPage         int
}

// Validate validates the options.
func (o *CollaborationsListOptions) Validate() error {
	return nil
}

// CollaborationsMembersOptions holds options for listing collaboration members.
type CollaborationsMembersOptions struct {
	CollaborationID int64
	PerPage         int
}

// Validate validates the options.
func (o *CollaborationsMembersOptions) Validate() error {
	return ValidateRequired("collaboration-id", o.CollaborationID)
}
