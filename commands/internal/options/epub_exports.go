package options

// EpubExportsListOptions holds options for listing ePub exports.
type EpubExportsListOptions struct{}

// Validate validates the options.
func (o *EpubExportsListOptions) Validate() error {
	return nil
}

// EpubExportCreateOptions holds options for creating an ePub export.
type EpubExportCreateOptions struct {
	CourseID int64
}

// Validate validates the options.
func (o *EpubExportCreateOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// EpubExportGetOptions holds options for getting an ePub export.
type EpubExportGetOptions struct {
	CourseID int64
	ID       int64
}

// Validate validates the options.
func (o *EpubExportGetOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("id", o.ID)
}
