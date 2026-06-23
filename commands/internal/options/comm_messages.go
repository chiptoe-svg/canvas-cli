package options

// CommMessagesListOptions holds options for listing communication messages.
type CommMessagesListOptions struct {
	UserID  int64
	PerPage int
}

// Validate validates the options.
func (o *CommMessagesListOptions) Validate() error {
	return nil
}
