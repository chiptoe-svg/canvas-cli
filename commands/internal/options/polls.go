package options

// PollListOptions contains options for listing polls
type PollListOptions struct{}

// Validate validates the options
func (o *PollListOptions) Validate() error {
	return nil
}

// PollGetOptions contains options for getting a poll
type PollGetOptions struct {
	PollID int64
}

// Validate validates the options
func (o *PollGetOptions) Validate() error {
	return ValidateRequired("poll-id", o.PollID)
}

// PollCreateOptions contains options for creating a poll
type PollCreateOptions struct {
	Question    string
	Description string
}

// Validate validates the options
func (o *PollCreateOptions) Validate() error {
	return ValidateRequired("question", o.Question)
}

// PollUpdateOptions contains options for updating a poll
type PollUpdateOptions struct {
	PollID      int64
	Question    string
	Description string
}

// Validate validates the options
func (o *PollUpdateOptions) Validate() error {
	return ValidateRequired("poll-id", o.PollID)
}

// PollDeleteOptions contains options for deleting a poll
type PollDeleteOptions struct {
	PollID int64
	Force  bool
}

// Validate validates the options
func (o *PollDeleteOptions) Validate() error {
	return ValidateRequired("poll-id", o.PollID)
}

// PollChoiceListOptions contains options for listing poll choices
type PollChoiceListOptions struct {
	PollID int64
}

// Validate validates the options
func (o *PollChoiceListOptions) Validate() error {
	return ValidateRequired("poll-id", o.PollID)
}

// PollChoiceGetOptions contains options for getting a poll choice
type PollChoiceGetOptions struct {
	PollID   int64
	ChoiceID int64
}

// Validate validates the options
func (o *PollChoiceGetOptions) Validate() error {
	if err := ValidateRequired("poll-id", o.PollID); err != nil {
		return err
	}
	return ValidateRequired("choice-id", o.ChoiceID)
}

// PollChoiceCreateOptions contains options for creating a poll choice
type PollChoiceCreateOptions struct {
	PollID    int64
	Text      string
	IsCorrect bool
	Position  int
}

// Validate validates the options
func (o *PollChoiceCreateOptions) Validate() error {
	if err := ValidateRequired("poll-id", o.PollID); err != nil {
		return err
	}
	return ValidateRequired("text", o.Text)
}

// PollChoiceUpdateOptions contains options for updating a poll choice
type PollChoiceUpdateOptions struct {
	PollID    int64
	ChoiceID  int64
	Text      string
	IsCorrect *bool
	Position  int
}

// Validate validates the options
func (o *PollChoiceUpdateOptions) Validate() error {
	if err := ValidateRequired("poll-id", o.PollID); err != nil {
		return err
	}
	return ValidateRequired("choice-id", o.ChoiceID)
}

// PollChoiceDeleteOptions contains options for deleting a poll choice
type PollChoiceDeleteOptions struct {
	PollID   int64
	ChoiceID int64
	Force    bool
}

// Validate validates the options
func (o *PollChoiceDeleteOptions) Validate() error {
	if err := ValidateRequired("poll-id", o.PollID); err != nil {
		return err
	}
	return ValidateRequired("choice-id", o.ChoiceID)
}

// PollSessionListOptions contains options for listing poll sessions
type PollSessionListOptions struct {
	PollID int64
}

// Validate validates the options
func (o *PollSessionListOptions) Validate() error {
	return ValidateRequired("poll-id", o.PollID)
}

// PollSessionGetOptions contains options for getting a poll session
type PollSessionGetOptions struct {
	PollID    int64
	SessionID int64
}

// Validate validates the options
func (o *PollSessionGetOptions) Validate() error {
	if err := ValidateRequired("poll-id", o.PollID); err != nil {
		return err
	}
	return ValidateRequired("session-id", o.SessionID)
}

// PollSessionCreateOptions contains options for creating a poll session
type PollSessionCreateOptions struct {
	PollID           int64
	CourseID         int64
	CourseSectionID  int64
	HasPublicResults bool
}

// Validate validates the options
func (o *PollSessionCreateOptions) Validate() error {
	if err := ValidateRequired("poll-id", o.PollID); err != nil {
		return err
	}
	return ValidateRequired("course-id", o.CourseID)
}

// PollSessionUpdateOptions contains options for updating a poll session
type PollSessionUpdateOptions struct {
	PollID           int64
	SessionID        int64
	CourseID         int64
	CourseSectionID  int64
	HasPublicResults *bool
}

// Validate validates the options
func (o *PollSessionUpdateOptions) Validate() error {
	if err := ValidateRequired("poll-id", o.PollID); err != nil {
		return err
	}
	return ValidateRequired("session-id", o.SessionID)
}

// PollSessionDeleteOptions contains options for deleting a poll session
type PollSessionDeleteOptions struct {
	PollID    int64
	SessionID int64
	Force     bool
}

// Validate validates the options
func (o *PollSessionDeleteOptions) Validate() error {
	if err := ValidateRequired("poll-id", o.PollID); err != nil {
		return err
	}
	return ValidateRequired("session-id", o.SessionID)
}

// PollSessionOpenOptions contains options for opening a poll session
type PollSessionOpenOptions struct {
	PollID    int64
	SessionID int64
}

// Validate validates the options
func (o *PollSessionOpenOptions) Validate() error {
	if err := ValidateRequired("poll-id", o.PollID); err != nil {
		return err
	}
	return ValidateRequired("session-id", o.SessionID)
}

// PollSessionCloseOptions contains options for closing a poll session
type PollSessionCloseOptions struct {
	PollID    int64
	SessionID int64
}

// Validate validates the options
func (o *PollSessionCloseOptions) Validate() error {
	if err := ValidateRequired("poll-id", o.PollID); err != nil {
		return err
	}
	return ValidateRequired("session-id", o.SessionID)
}

// PollSubmissionGetOptions contains options for getting a poll submission
type PollSubmissionGetOptions struct {
	PollID       int64
	SessionID    int64
	SubmissionID int64
}

// Validate validates the options
func (o *PollSubmissionGetOptions) Validate() error {
	if err := ValidateRequired("poll-id", o.PollID); err != nil {
		return err
	}
	if err := ValidateRequired("session-id", o.SessionID); err != nil {
		return err
	}
	return ValidateRequired("submission-id", o.SubmissionID)
}

// PollSubmissionCreateOptions contains options for creating a poll submission
type PollSubmissionCreateOptions struct {
	PollID       int64
	SessionID    int64
	PollChoiceID int64
}

// Validate validates the options
func (o *PollSubmissionCreateOptions) Validate() error {
	if err := ValidateRequired("poll-id", o.PollID); err != nil {
		return err
	}
	if err := ValidateRequired("session-id", o.SessionID); err != nil {
		return err
	}
	return ValidateRequired("choice-id", o.PollChoiceID)
}
