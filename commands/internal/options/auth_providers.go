package options

// AuthProvidersListOptions contains options for listing authentication providers
type AuthProvidersListOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *AuthProvidersListOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// AuthProvidersGetOptions contains options for getting an authentication provider
type AuthProvidersGetOptions struct {
	AccountID  int64
	ProviderID int64
}

// Validate validates the options
func (o *AuthProvidersGetOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("provider-id", o.ProviderID)
}

// AuthProvidersCreateOptions contains options for creating an authentication provider
type AuthProvidersCreateOptions struct {
	AccountID      int64
	AuthType       string
	ClientID       string
	ClientSecret   string
	LoginAttribute string
}

// Validate validates the options
func (o *AuthProvidersCreateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("auth-type", o.AuthType)
}

// AuthProvidersDeleteOptions contains options for deleting an authentication provider
type AuthProvidersDeleteOptions struct {
	AccountID  int64
	ProviderID int64
}

// Validate validates the options
func (o *AuthProvidersDeleteOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("provider-id", o.ProviderID)
}

// AuthProvidersRestoreOptions contains options for restoring an authentication provider
type AuthProvidersRestoreOptions struct {
	AccountID  int64
	ProviderID int64
}

// Validate validates the options
func (o *AuthProvidersRestoreOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("provider-id", o.ProviderID)
}

// AuthProvidersForcePasswordResetOptions contains options for force password reset
type AuthProvidersForcePasswordResetOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *AuthProvidersForcePasswordResetOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// AuthProvidersSSOSettingsOptions contains options for getting SSO settings
type AuthProvidersSSOSettingsOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *AuthProvidersSSOSettingsOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}
