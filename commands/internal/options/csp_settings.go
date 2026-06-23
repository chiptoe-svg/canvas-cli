package options

// CSPSettingsGetOptions contains options for getting CSP settings
type CSPSettingsGetOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *CSPSettingsGetOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// CSPSettingsAddDomainOptions contains options for adding a domain to CSP allowlist
type CSPSettingsAddDomainOptions struct {
	AccountID int64
	Domain    string
}

// Validate validates the options
func (o *CSPSettingsAddDomainOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("domain", o.Domain)
}

// CSPSettingsRemoveDomainOptions contains options for removing a domain from CSP allowlist
type CSPSettingsRemoveDomainOptions struct {
	AccountID int64
	Domain    string
}

// Validate validates the options
func (o *CSPSettingsRemoveDomainOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("domain", o.Domain)
}

// CSPSettingsLockOptions contains options for locking CSP settings
type CSPSettingsLockOptions struct {
	AccountID int64
	Locked    bool
}

// Validate validates the options
func (o *CSPSettingsLockOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}
