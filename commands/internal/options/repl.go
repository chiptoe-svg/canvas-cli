package options

// ReplOptions contains options for the repl/shell command.
// The repl command currently exposes no flags; the struct is here for consistency
// and to make future additions straightforward.
type ReplOptions struct{}

// Validate validates the options.
func (o *ReplOptions) Validate() error { return nil }
