package options

// AccountAnalyticsTermOptions contains options for term-scoped analytics
type AccountAnalyticsTermOptions struct {
	AccountID int64
	TermID    int64
}

// Validate validates the options
func (o *AccountAnalyticsTermOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("term-id", o.TermID)
}

// AccountAnalyticsCompletedOptions contains options for completed-course analytics
type AccountAnalyticsCompletedOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *AccountAnalyticsCompletedOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}
