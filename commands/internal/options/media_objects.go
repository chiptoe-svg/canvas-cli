package options

// MediaObjectsListOptions holds options for listing media objects.
type MediaObjectsListOptions struct {
	CourseID int64
	GroupID  int64
	Sort     string
	Order    string
	PerPage  int
}

// Validate validates the options.
func (o *MediaObjectsListOptions) Validate() error {
	return nil
}

// MediaObjectUpdateOptions holds options for updating a media object.
type MediaObjectUpdateOptions struct {
	MediaID string
	Title   string
}

// Validate validates the options.
func (o *MediaObjectUpdateOptions) Validate() error {
	if err := ValidateRequired("media-id", o.MediaID); err != nil {
		return err
	}
	return ValidateRequired("title", o.Title)
}

// MediaTracksListOptions holds options for listing media tracks.
type MediaTracksListOptions struct {
	MediaID      string
	AttachmentID int64
}

// Validate validates the options.
func (o *MediaTracksListOptions) Validate() error {
	return nil
}

// MediaAttachmentsListOptions holds options for listing media attachments.
type MediaAttachmentsListOptions struct {
	CourseID int64
	GroupID  int64
	Sort     string
	Order    string
	PerPage  int
}

// Validate validates the options.
func (o *MediaAttachmentsListOptions) Validate() error {
	return nil
}

// MediaAttachmentUpdateOptions holds options for updating a media attachment.
type MediaAttachmentUpdateOptions struct {
	AttachmentID int64
	Title        string
}

// Validate validates the options.
func (o *MediaAttachmentUpdateOptions) Validate() error {
	if err := ValidateRequired("attachment-id", o.AttachmentID); err != nil {
		return err
	}
	return ValidateRequired("title", o.Title)
}
