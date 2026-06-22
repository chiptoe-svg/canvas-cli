package options

// DeveloperKeysListOptions holds options for listing developer keys.
type DeveloperKeysListOptions struct {
	AccountID int64
}

// Validate validates the options.
func (o *DeveloperKeysListOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// DeveloperKeysCreateOptions holds options for creating a developer key.
type DeveloperKeysCreateOptions struct {
	AccountID   int64
	Name        string
	Email       string
	RedirectURI string
	Notes       string
}

// Validate validates the options.
func (o *DeveloperKeysCreateOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// DeveloperKeysBindOptions holds options for binding a developer key.
type DeveloperKeysBindOptions struct {
	AccountID      int64
	DeveloperKeyID int64
	WorkflowState  string
}

// Validate validates the options.
func (o *DeveloperKeysBindOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	if err := ValidateRequired("developer-key-id", o.DeveloperKeyID); err != nil {
		return err
	}
	return ValidateRequired("workflow-state", o.WorkflowState)
}
