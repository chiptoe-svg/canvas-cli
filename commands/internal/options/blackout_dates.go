package options

// BlackoutDatesListOptions contains options for listing blackout dates
type BlackoutDatesListOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *BlackoutDatesListOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// BlackoutDatesGetOptions contains options for getting a blackout date
type BlackoutDatesGetOptions struct {
	AccountID int64
	ID        int64
}

// Validate validates the options
func (o *BlackoutDatesGetOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("id", o.ID)
}

// BlackoutDatesCreateOptions contains options for creating a blackout date
type BlackoutDatesCreateOptions struct {
	AccountID  int64
	StartDate  string
	EndDate    string
	EventTitle string
}

// Validate validates the options
func (o *BlackoutDatesCreateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	if err := ValidateRequired("start-date", o.StartDate); err != nil {
		return err
	}
	return ValidateRequired("end-date", o.EndDate)
}

// BlackoutDatesUpdateOptions contains options for updating a blackout date
type BlackoutDatesUpdateOptions struct {
	AccountID  int64
	ID         int64
	StartDate  string
	EndDate    string
	EventTitle string
}

// Validate validates the options
func (o *BlackoutDatesUpdateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("id", o.ID)
}

// BlackoutDatesDeleteOptions contains options for deleting a blackout date
type BlackoutDatesDeleteOptions struct {
	AccountID int64
	ID        int64
}

// Validate validates the options
func (o *BlackoutDatesDeleteOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("id", o.ID)
}
