package options

// AccountExtToolsFavoritesOptions contains options for managing account external tool favorites
type AccountExtToolsFavoritesOptions struct {
	AccountID int64
	ToolID    int64
}

// Validate validates the options
func (o *AccountExtToolsFavoritesOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("tool-id", o.ToolID)
}
