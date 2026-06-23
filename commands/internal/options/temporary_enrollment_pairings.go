package options

// TemporaryEnrollmentPairingsListOptions contains options for listing temporary enrollment pairings
type TemporaryEnrollmentPairingsListOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *TemporaryEnrollmentPairingsListOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// TemporaryEnrollmentPairingsGetOptions contains options for getting a temporary enrollment pairing
type TemporaryEnrollmentPairingsGetOptions struct {
	AccountID int64
	ID        int64
}

// Validate validates the options
func (o *TemporaryEnrollmentPairingsGetOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("id", o.ID)
}

// TemporaryEnrollmentPairingsCreateOptions contains options for creating a temporary enrollment pairing
type TemporaryEnrollmentPairingsCreateOptions struct {
	AccountID               int64
	StartingEnrollmentState string
	RoleID                  int64
}

// Validate validates the options
func (o *TemporaryEnrollmentPairingsCreateOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// TemporaryEnrollmentPairingsDeleteOptions contains options for deleting a temporary enrollment pairing
type TemporaryEnrollmentPairingsDeleteOptions struct {
	AccountID int64
	ID        int64
}

// Validate validates the options
func (o *TemporaryEnrollmentPairingsDeleteOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("id", o.ID)
}
