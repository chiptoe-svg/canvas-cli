package options

// APIOptions contains options for the raw API command
type APIOptions struct {
	Data        string
	DataFile    string
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
