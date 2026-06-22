package options

// ProgressGetOptions holds options for getting a progress job.
type ProgressGetOptions struct {
	ID int64
}

// Validate validates the options.
func (o *ProgressGetOptions) Validate() error {
	return ValidateRequired("id", o.ID)
}

// ProgressCancelOptions holds options for cancelling a progress job.
type ProgressCancelOptions struct {
	ID    int64
	Force bool
}

// Validate validates the options.
func (o *ProgressCancelOptions) Validate() error {
	return ValidateRequired("id", o.ID)
}
