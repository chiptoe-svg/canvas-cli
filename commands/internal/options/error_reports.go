package options

// ErrorReportCreateOptions holds options for submitting an error report.
type ErrorReportCreateOptions struct {
	Subject  string
	Comments string
	URL      string
	Email    string
	Severity string
}

// Validate validates the options.
func (o *ErrorReportCreateOptions) Validate() error {
	return ValidateRequired("subject", o.Subject)
}
