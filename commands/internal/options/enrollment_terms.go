package options

// EnrollmentTermsListOptions holds options for listing enrollment terms.
type EnrollmentTermsListOptions struct {
	AccountID int64
}

// Validate validates the options.
func (o *EnrollmentTermsListOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// EnrollmentTermsGetOptions holds options for getting a single enrollment term.
type EnrollmentTermsGetOptions struct {
	AccountID int64
	TermID    int64
}

// Validate validates the options.
func (o *EnrollmentTermsGetOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("term-id", o.TermID)
}

// EnrollmentTermsCreateOptions holds options for creating an enrollment term.
type EnrollmentTermsCreateOptions struct {
	AccountID int64
	Name      string
	StartAt   string
	EndAt     string
	SISTermID string
}

// Validate validates the options.
func (o *EnrollmentTermsCreateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("name", o.Name)
}

// EnrollmentTermsUpdateOptions holds options for updating an enrollment term.
type EnrollmentTermsUpdateOptions struct {
	AccountID int64
	TermID    int64
	Name      string
	StartAt   string
	EndAt     string
}

// Validate validates the options.
func (o *EnrollmentTermsUpdateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("term-id", o.TermID)
}

// EnrollmentTermsDeleteOptions holds options for deleting an enrollment term.
type EnrollmentTermsDeleteOptions struct {
	AccountID int64
	TermID    int64
}

// Validate validates the options.
func (o *EnrollmentTermsDeleteOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("term-id", o.TermID)
}
