package options

// CompletionOptions contains options for the shell-completion command.
// The completion command currently exposes no flags; the struct is here for
// consistency with the project's options pattern.
type CompletionOptions struct{}

// Validate validates the options.
func (o *CompletionOptions) Validate() error { return nil }
