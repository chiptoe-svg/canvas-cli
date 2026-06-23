package options

// AccountCalendarsListOptions contains options for listing all account calendars
type AccountCalendarsListOptions struct {
	Search string
}

// Validate validates the options
func (o *AccountCalendarsListOptions) Validate() error {
	// No required fields
	return nil
}

// AccountCalendarsGetOptions contains options for getting an account calendar
type AccountCalendarsGetOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *AccountCalendarsGetOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// AccountCalendarsUpdateOptions contains options for updating an account calendar
type AccountCalendarsUpdateOptions struct {
	AccountID     int64
	Visible       bool
	VisibleSet    bool // tracks whether --visible flag was provided
	AutoSubscribe bool
	AutoSubSet    bool // tracks whether --auto-subscribe flag was provided
}

// Validate validates the options
func (o *AccountCalendarsUpdateOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// AccountCalendarsListForAccountOptions contains options for listing calendars for an account
type AccountCalendarsListForAccountOptions struct {
	AccountID int64
	Search    string
}

// Validate validates the options
func (o *AccountCalendarsListForAccountOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}
