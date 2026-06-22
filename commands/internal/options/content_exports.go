package options

import "fmt"

// ContentExportsListOptions holds options for listing content exports.
type ContentExportsListOptions struct {
	CourseID int64
}

// Validate validates the options.
func (o *ContentExportsListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	return nil
}

// ContentExportsGetOptions holds options for getting a content export.
type ContentExportsGetOptions struct {
	CourseID int64
	ID       int64
}

// Validate validates the options.
func (o *ContentExportsGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.ID <= 0 {
		return fmt.Errorf("export-id is required")
	}
	return nil
}

// ContentExportsCreateOptions holds options for creating a content export.
type ContentExportsCreateOptions struct {
	CourseID          int64
	ExportType        string
	SkipNotifications bool
}

// Validate validates the options.
func (o *ContentExportsCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.ExportType == "" {
		return fmt.Errorf("--export-type is required")
	}
	return nil
}

// EpubExportsGetOptions holds options for getting an epub export.
type EpubExportsGetOptions struct {
	CourseID int64
	ID       int64
}

// Validate validates the options.
func (o *EpubExportsGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("--course-id is required")
	}
	if o.ID <= 0 {
		return fmt.Errorf("epub-id is required")
	}
	return nil
}
