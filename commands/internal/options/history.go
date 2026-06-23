package options

// HistoryListOptions holds options for listing user history.
type HistoryListOptions struct {
	UserID int64
}

// Validate validates the options.
func (o *HistoryListOptions) Validate() error {
	return ValidateRequired("user-id", o.UserID)
}
