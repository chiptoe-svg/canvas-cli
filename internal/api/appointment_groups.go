package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// AppointmentGroup represents a Canvas appointment group (office hours, sign-up slots, etc.)
type AppointmentGroup struct {
	ID                            int64             `json:"id"`
	Title                         string            `json:"title"`
	StartAt                       *time.Time        `json:"start_at,omitempty"`
	EndAt                         *time.Time        `json:"end_at,omitempty"`
	Description                   string            `json:"description,omitempty"`
	LocationName                  string            `json:"location_name,omitempty"`
	LocationAddress               string            `json:"location_address,omitempty"`
	ParticipantCount              int               `json:"participant_count,omitempty"`
	ReservedTimes                 []AppointmentSlot `json:"reserved_times,omitempty"`
	AllowObserverSignup           bool              `json:"allow_observer_signup"`
	ContextCodes                  []string          `json:"context_codes"`
	SubContextCodes               []string          `json:"sub_context_codes,omitempty"`
	WorkflowState                 string            `json:"workflow_state"`
	RequiringAction               bool              `json:"requiring_action"`
	AppointmentsCount             int               `json:"appointments_count"`
	Appointments                  []CalendarEvent   `json:"appointments,omitempty"`
	NewAppointments               []CalendarEvent   `json:"new_appointments,omitempty"`
	MaxAppointmentsPerParticipant *int              `json:"max_appointments_per_participant,omitempty"`
	MinAppointmentsPerParticipant *int              `json:"min_appointments_per_participant,omitempty"`
	ParticipantsPerAppointment    *int              `json:"participants_per_appointment,omitempty"`
	ParticipantVisibility         string            `json:"participant_visibility,omitempty"`
	ParticipantType               string            `json:"participant_type,omitempty"`
	URL                           string            `json:"url,omitempty"`
	HTMLURL                       string            `json:"html_url,omitempty"`
	CreatedAt                     time.Time         `json:"created_at"`
	UpdatedAt                     time.Time         `json:"updated_at"`
}

// AppointmentSlot represents a reserved time slot
type AppointmentSlot struct {
	ID      int64      `json:"id"`
	StartAt *time.Time `json:"start_at,omitempty"`
	EndAt   *time.Time `json:"end_at,omitempty"`
}

// AppointmentGroupsService handles appointment-group-related API calls
type AppointmentGroupsService struct {
	client *Client
}

// NewAppointmentGroupsService creates a new appointment groups service
func NewAppointmentGroupsService(client *Client) *AppointmentGroupsService {
	return &AppointmentGroupsService{client: client}
}

// ListAppointmentGroupsOptions holds query options for listing appointment groups
type ListAppointmentGroupsOptions struct {
	Scope                   string   // "reservable" or "manageable"
	ContextCodes            []string // filter by context codes
	IncludePastAppointments bool
	Include                 []string // appointments, child_events, participant_count, reserved_times, all_context_codes
}

// List retrieves appointment groups that can be reserved or managed by the current user
func (s *AppointmentGroupsService) List(ctx context.Context, opts *ListAppointmentGroupsOptions) ([]AppointmentGroup, error) {
	path := "/api/v1/appointment_groups"

	if opts != nil {
		query := url.Values{}

		if opts.Scope != "" {
			query.Set("scope", opts.Scope)
		}

		for _, code := range opts.ContextCodes {
			query.Add("context_codes[]", code)
		}

		if opts.IncludePastAppointments {
			query.Set("include_past_appointments", "1")
		}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var groups []AppointmentGroup
	if err := s.client.GetAllPages(ctx, path, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// Get retrieves a single appointment group
func (s *AppointmentGroupsService) Get(ctx context.Context, groupID int64, include []string) (*AppointmentGroup, error) {
	path := fmt.Sprintf("/api/v1/appointment_groups/%d", groupID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var group AppointmentGroup
	if err := s.client.GetJSON(ctx, path, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// CreateAppointmentGroupParams holds parameters for creating an appointment group
type CreateAppointmentGroupParams struct {
	// Required
	ContextCodes []string
	Title        string

	// Optional
	SubContextCodes               []string
	Description                   string
	LocationName                  string
	LocationAddress               string
	Publish                       bool
	ParticipantsPerAppointment    int
	MinAppointmentsPerParticipant int
	MaxAppointmentsPerParticipant int
	ParticipantVisibility         string // "private" or "protected"
	AllowObserverSignup           bool
	// NewAppointments is a list of [startAt, endAt] pairs (ISO 8601)
	NewAppointments [][2]string
}

// Create creates a new appointment group
func (s *AppointmentGroupsService) Create(ctx context.Context, params *CreateAppointmentGroupParams) (*AppointmentGroup, error) {
	path := "/api/v1/appointment_groups"

	ag := buildAppointmentGroupBody(params.ContextCodes, params.SubContextCodes, params.Title,
		params.Description, params.LocationName, params.LocationAddress,
		params.ParticipantsPerAppointment, params.MinAppointmentsPerParticipant,
		params.MaxAppointmentsPerParticipant, params.ParticipantVisibility)

	if params.Publish {
		ag["publish"] = true
	}

	if params.AllowObserverSignup {
		ag["allow_observer_signup"] = true
	}

	if len(params.NewAppointments) > 0 {
		newAppts := make(map[string]interface{})
		for i, pair := range params.NewAppointments {
			newAppts[strconv.Itoa(i)] = []string{pair[0], pair[1]}
		}
		ag["new_appointments"] = newAppts
	}

	body := map[string]interface{}{
		"appointment_group": ag,
	}

	var group AppointmentGroup
	if err := s.client.PostJSON(ctx, path, body, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// UpdateAppointmentGroupParams holds parameters for updating an appointment group
type UpdateAppointmentGroupParams struct {
	// Context codes are required on update per spec
	ContextCodes []string

	// Optional
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
	NewAppointments               [][2]string
}

// Update updates an existing appointment group
func (s *AppointmentGroupsService) Update(ctx context.Context, groupID int64, params *UpdateAppointmentGroupParams) (*AppointmentGroup, error) {
	path := fmt.Sprintf("/api/v1/appointment_groups/%d", groupID)

	ag := buildAppointmentGroupBody(params.ContextCodes, params.SubContextCodes, params.Title,
		params.Description, params.LocationName, params.LocationAddress,
		params.ParticipantsPerAppointment, params.MinAppointmentsPerParticipant,
		params.MaxAppointmentsPerParticipant, params.ParticipantVisibility)

	if params.Publish {
		ag["publish"] = true
	}

	if params.AllowObserverSignup != nil {
		ag["allow_observer_signup"] = *params.AllowObserverSignup
	}

	if len(params.NewAppointments) > 0 {
		newAppts := make(map[string]interface{})
		for i, pair := range params.NewAppointments {
			newAppts[strconv.Itoa(i)] = []string{pair[0], pair[1]}
		}
		ag["new_appointments"] = newAppts
	}

	body := map[string]interface{}{
		"appointment_group": ag,
	}

	var group AppointmentGroup
	if err := s.client.PutJSON(ctx, path, body, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// Delete removes an appointment group and all associated slots and reservations
func (s *AppointmentGroupsService) Delete(ctx context.Context, groupID int64, cancelReason string) (*AppointmentGroup, error) {
	path := fmt.Sprintf("/api/v1/appointment_groups/%d", groupID)

	if cancelReason != "" {
		query := url.Values{}
		query.Set("cancel_reason", cancelReason)
		path += "?" + query.Encode()
	}

	var group AppointmentGroup
	if err := s.client.DeleteJSON(ctx, path, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// ListUsers returns a paginated list of users participating in an appointment group
func (s *AppointmentGroupsService) ListUsers(ctx context.Context, groupID int64, registrationStatus string) ([]User, error) {
	path := fmt.Sprintf("/api/v1/appointment_groups/%d/users", groupID)

	if registrationStatus != "" {
		query := url.Values{}
		query.Set("registration_status", registrationStatus)
		path += "?" + query.Encode()
	}

	var users []User
	if err := s.client.GetAllPages(ctx, path, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// ListGroups returns a paginated list of student groups participating in an appointment group
func (s *AppointmentGroupsService) ListGroups(ctx context.Context, groupID int64, registrationStatus string) ([]Group, error) {
	path := fmt.Sprintf("/api/v1/appointment_groups/%d/groups", groupID)

	if registrationStatus != "" {
		query := url.Values{}
		query.Set("registration_status", registrationStatus)
		path += "?" + query.Encode()
	}

	var groups []Group
	if err := s.client.GetAllPages(ctx, path, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// NextAppointment returns the next appointment available to sign up for
func (s *AppointmentGroupsService) NextAppointment(ctx context.Context, appointmentGroupIDs []int64) ([]CalendarEvent, error) {
	path := "/api/v1/appointment_groups/next_appointment"

	if len(appointmentGroupIDs) > 0 {
		query := url.Values{}
		for _, id := range appointmentGroupIDs {
			query.Add("appointment_group_ids[]", strconv.FormatInt(id, 10))
		}
		path += "?" + query.Encode()
	}

	var events []CalendarEvent
	if err := s.client.GetJSON(ctx, path, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// buildAppointmentGroupBody constructs the shared appointment_group body map
func buildAppointmentGroupBody(contextCodes, subContextCodes []string, title, description,
	locationName, locationAddress string, participantsPerAppt, minAppts, maxAppts int,
	participantVisibility string) map[string]interface{} {

	ag := map[string]interface{}{}

	if len(contextCodes) > 0 {
		ag["context_codes"] = contextCodes
	}

	if len(subContextCodes) > 0 {
		ag["sub_context_codes"] = subContextCodes
	}

	if title != "" {
		ag["title"] = title
	}

	if description != "" {
		ag["description"] = description
	}

	if locationName != "" {
		ag["location_name"] = locationName
	}

	if locationAddress != "" {
		ag["location_address"] = locationAddress
	}

	if participantsPerAppt > 0 {
		ag["participants_per_appointment"] = participantsPerAppt
	}

	if minAppts > 0 {
		ag["min_appointments_per_participant"] = minAppts
	}

	if maxAppts > 0 {
		ag["max_appointments_per_participant"] = maxAppts
	}

	if participantVisibility != "" {
		ag["participant_visibility"] = participantVisibility
	}

	return ag
}
