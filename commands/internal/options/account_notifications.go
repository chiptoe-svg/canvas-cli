package options

// AccountNotificationsListOptions contains options for listing account notifications
type AccountNotificationsListOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *AccountNotificationsListOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// AccountNotificationsGetOptions contains options for getting an account notification
type AccountNotificationsGetOptions struct {
	AccountID      int64
	NotificationID int64
}

// Validate validates the options
func (o *AccountNotificationsGetOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("notification-id", o.NotificationID)
}

// AccountNotificationsCreateOptions contains options for creating an account notification
type AccountNotificationsCreateOptions struct {
	AccountID int64
	Subject   string
	Message   string
	StartAt   string
	EndAt     string
	Icon      string
	Roles     []string
}

// Validate validates the options
func (o *AccountNotificationsCreateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	if err := ValidateRequired("subject", o.Subject); err != nil {
		return err
	}
	if err := ValidateRequired("message", o.Message); err != nil {
		return err
	}
	if err := ValidateRequired("start-at", o.StartAt); err != nil {
		return err
	}
	return ValidateRequired("end-at", o.EndAt)
}

// AccountNotificationsUpdateOptions contains options for updating an account notification
type AccountNotificationsUpdateOptions struct {
	AccountID      int64
	NotificationID int64
	Subject        string
	Message        string
	StartAt        string
	EndAt          string
	Icon           string
	Roles          []string
}

// Validate validates the options
func (o *AccountNotificationsUpdateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("notification-id", o.NotificationID)
}

// AccountNotificationsDeleteOptions contains options for deleting an account notification
type AccountNotificationsDeleteOptions struct {
	AccountID      int64
	NotificationID int64
}

// Validate validates the options
func (o *AccountNotificationsDeleteOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("notification-id", o.NotificationID)
}
