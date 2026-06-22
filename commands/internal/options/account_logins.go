package options

// AccountLoginsListOptions holds options for listing account logins.
type AccountLoginsListOptions struct {
	AccountID int64
	UserID    int64
}

// Validate validates the options.
func (o *AccountLoginsListOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// AccountLoginsCreateOptions holds options for creating an account login.
type AccountLoginsCreateOptions struct {
	AccountID int64
	UserID    int64
	UniqueID  string
	Password  string
	SISUserID string
}

// Validate validates the options.
func (o *AccountLoginsCreateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	if err := ValidateRequired("user-id", o.UserID); err != nil {
		return err
	}
	return ValidateRequired("unique-id", o.UniqueID)
}

// AccountLoginsUpdateOptions holds options for updating an account login.
type AccountLoginsUpdateOptions struct {
	AccountID int64
	LoginID   int64
	UniqueID  string
	Password  string
}

// Validate validates the options.
func (o *AccountLoginsUpdateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("login-id", o.LoginID)
}
