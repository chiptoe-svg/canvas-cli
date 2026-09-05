package options

// APIOptions contains options for the read-only "api get" command
type APIOptions struct {
	Query       []string
	Headers     []string
	Paginate    bool
	RawOutput   bool
	ShowHeaders bool
}

// Validate validates the options. No required fields beyond the positional args.
func (o *APIOptions) Validate() error {
	return nil
}
