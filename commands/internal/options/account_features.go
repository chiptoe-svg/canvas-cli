package options

// AccountFeaturesListOptions contains options for listing account features
type AccountFeaturesListOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *AccountFeaturesListOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// AccountFeaturesListEnabledOptions contains options for listing enabled account features
type AccountFeaturesListEnabledOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *AccountFeaturesListEnabledOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// AccountFeaturesGetFlagOptions contains options for getting an account feature flag
type AccountFeaturesGetFlagOptions struct {
	AccountID int64
	Feature   string
}

// Validate validates the options
func (o *AccountFeaturesGetFlagOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("feature", o.Feature)
}

// AccountFeaturesSetFlagOptions contains options for setting an account feature flag
type AccountFeaturesSetFlagOptions struct {
	AccountID int64
	Feature   string
	State     string
}

// Validate validates the options
func (o *AccountFeaturesSetFlagOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	if err := ValidateRequired("feature", o.Feature); err != nil {
		return err
	}
	return ValidateRequired("state", o.State)
}

// AccountFeaturesDeleteFlagOptions contains options for deleting an account feature flag
type AccountFeaturesDeleteFlagOptions struct {
	AccountID int64
	Feature   string
}

// Validate validates the options
func (o *AccountFeaturesDeleteFlagOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("feature", o.Feature)
}

// AccountFeaturesSettingsOptions contains options for getting account settings
type AccountFeaturesSettingsOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *AccountFeaturesSettingsOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// AccountFeaturesPermissionsOptions contains options for getting account permissions
type AccountFeaturesPermissionsOptions struct {
	AccountID   int64
	Permissions []string
}

// Validate validates the options
func (o *AccountFeaturesPermissionsOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}
