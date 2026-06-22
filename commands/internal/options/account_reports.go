package options

// AccountReportsListOptions contains options for listing account reports
type AccountReportsListOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *AccountReportsListOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// AccountReportsRunsOptions contains options for listing report runs
type AccountReportsRunsOptions struct {
	AccountID  int64
	ReportName string
}

// Validate validates the options
func (o *AccountReportsRunsOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("report-name", o.ReportName)
}

// AccountReportsStartOptions contains options for starting a report run
type AccountReportsStartOptions struct {
	AccountID  int64
	ReportName string
}

// Validate validates the options
func (o *AccountReportsStartOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("report-name", o.ReportName)
}

// AccountReportsGetRunOptions contains options for getting a specific report run
type AccountReportsGetRunOptions struct {
	AccountID  int64
	ReportName string
	RunID      int64
}

// Validate validates the options
func (o *AccountReportsGetRunOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	if err := ValidateRequired("report-name", o.ReportName); err != nil {
		return err
	}
	return ValidateRequired("run-id", o.RunID)
}

// AccountReportsDeleteRunOptions contains options for deleting a report run
type AccountReportsDeleteRunOptions struct {
	AccountID  int64
	ReportName string
	RunID      int64
}

// Validate validates the options
func (o *AccountReportsDeleteRunOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	if err := ValidateRequired("report-name", o.ReportName); err != nil {
		return err
	}
	return ValidateRequired("run-id", o.RunID)
}

// AccountReportsAbortRunOptions contains options for aborting a report run
type AccountReportsAbortRunOptions struct {
	AccountID  int64
	ReportName string
	RunID      int64
}

// Validate validates the options
func (o *AccountReportsAbortRunOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	if err := ValidateRequired("report-name", o.ReportName); err != nil {
		return err
	}
	return ValidateRequired("run-id", o.RunID)
}
