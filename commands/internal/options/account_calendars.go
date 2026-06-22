package options

// AccountCalendarsListOptions holds options for listing account calendars.
type AccountCalendarsListOptions struct {
	AccountID int64
	Filter    string
	PerPage   int
}

// Validate validates the options.
func (o *AccountCalendarsListOptions) Validate() error {
	return nil
}

// AccountCalendarGetOptions holds options for getting an account calendar.
type AccountCalendarGetOptions struct {
	AccountID int64
}

// Validate validates the options.
func (o *AccountCalendarGetOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// AccountCalendarUpdateOptions holds options for updating an account calendar.
type AccountCalendarUpdateOptions struct {
	AccountID     int64
	Visible       bool
	AutoSubscribe bool
	VisibleSet    bool
	AutoSubSet    bool
}

// Validate validates the options.
func (o *AccountCalendarUpdateOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}
