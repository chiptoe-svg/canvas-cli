package options

// AppointmentGroupListOptions contains options for listing appointment groups
type AppointmentGroupListOptions struct {
	Scope                   string
	ContextCodes            []string
	IncludePastAppointments bool
	Include                 []string
}

// Validate validates the options
func (o *AppointmentGroupListOptions) Validate() error {
	return nil
}

// AppointmentGroupGetOptions contains options for getting an appointment group
type AppointmentGroupGetOptions struct {
	GroupID int64
	Include []string
}

// Validate validates the options
func (o *AppointmentGroupGetOptions) Validate() error {
	return ValidateRequired("group-id", o.GroupID)
}

// AppointmentGroupCreateOptions contains options for creating an appointment group
type AppointmentGroupCreateOptions struct {
	ContextCodes                  []string
	SubContextCodes               []string
	Title                         string
	Description                   string
	LocationName                  string
	LocationAddress               string
	Publish                       bool
	ParticipantsPerAppointment    int
	MinAppointmentsPerParticipant int
	MaxAppointmentsPerParticipant int
	ParticipantVisibility         string
	AllowObserverSignup           bool
}

// Validate validates the options
func (o *AppointmentGroupCreateOptions) Validate() error {
	if len(o.ContextCodes) == 0 {
		return ValidateRequired("context-codes", "")
	}
	return ValidateRequired("title", o.Title)
}

// AppointmentGroupUpdateOptions contains options for updating an appointment group
type AppointmentGroupUpdateOptions struct {
	GroupID                       int64
	ContextCodes                  []string
	SubContextCodes               []string
	Title                         string
	Description                   string
	LocationName                  string
	LocationAddress               string
	Publish                       bool
	ParticipantsPerAppointment    int
	MinAppointmentsPerParticipant int
	MaxAppointmentsPerParticipant int
	ParticipantVisibility         string
	AllowObserverSignup           *bool
}

// Validate validates the options
func (o *AppointmentGroupUpdateOptions) Validate() error {
	return ValidateRequired("group-id", o.GroupID)
}

// AppointmentGroupDeleteOptions contains options for deleting an appointment group
type AppointmentGroupDeleteOptions struct {
	GroupID      int64
	CancelReason string
	Force        bool
}

// Validate validates the options
func (o *AppointmentGroupDeleteOptions) Validate() error {
	return ValidateRequired("group-id", o.GroupID)
}

// AppointmentGroupUsersOptions contains options for listing appointment group users
type AppointmentGroupUsersOptions struct {
	GroupID            int64
	RegistrationStatus string
}

// Validate validates the options
func (o *AppointmentGroupUsersOptions) Validate() error {
	return ValidateRequired("group-id", o.GroupID)
}

// AppointmentGroupGroupsOptions contains options for listing appointment group student groups
type AppointmentGroupGroupsOptions struct {
	GroupID            int64
	RegistrationStatus string
}

// Validate validates the options
func (o *AppointmentGroupGroupsOptions) Validate() error {
	return ValidateRequired("group-id", o.GroupID)
}

// AppointmentGroupNextOptions contains options for the next-appointment endpoint
type AppointmentGroupNextOptions struct {
	AppointmentGroupIDs []int64
}

// Validate validates the options
func (o *AppointmentGroupNextOptions) Validate() error {
	return nil
}
