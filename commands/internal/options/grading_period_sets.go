package options

// GradingPeriodSetsListOptions holds options for listing grading period sets.
type GradingPeriodSetsListOptions struct {
	AccountID int64
}

// Validate validates the options.
func (o *GradingPeriodSetsListOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// GradingPeriodSetsCreateOptions holds options for creating a grading period set.
type GradingPeriodSetsCreateOptions struct {
	AccountID              int64
	Title                  string
	WeightedGradingPeriods bool
}

// Validate validates the options.
func (o *GradingPeriodSetsCreateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("title", o.Title)
}

// GradingPeriodSetsUpdateOptions holds options for updating a grading period set.
type GradingPeriodSetsUpdateOptions struct {
	AccountID int64
	SetID     int64
	Title     string
}

// Validate validates the options.
func (o *GradingPeriodSetsUpdateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("id", o.SetID)
}

// GradingPeriodSetsDeleteOptions holds options for deleting a grading period set.
type GradingPeriodSetsDeleteOptions struct {
	AccountID int64
	SetID     int64
}

// Validate validates the options.
func (o *GradingPeriodSetsDeleteOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("id", o.SetID)
}

// GradingPeriodSetsListPeriodsOptions holds options for listing grading periods.
type GradingPeriodSetsListPeriodsOptions struct {
	AccountID int64
}

// Validate validates the options.
func (o *GradingPeriodSetsListPeriodsOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// GradingPeriodSetsDeletePeriodOptions holds options for deleting a grading period.
type GradingPeriodSetsDeletePeriodOptions struct {
	AccountID int64
	PeriodID  int64
}

// Validate validates the options.
func (o *GradingPeriodSetsDeletePeriodOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("id", o.PeriodID)
}
