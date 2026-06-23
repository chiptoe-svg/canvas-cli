package options

// BrandGetOptions holds options for retrieving brand variables.
type BrandGetOptions struct {
	AccountID int64
	CourseID  int64
}

// Validate validates the options.
func (o *BrandGetOptions) Validate() error {
	return nil
}
