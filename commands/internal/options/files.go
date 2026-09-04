package options

import "fmt"

// FilesListOptions contains options for listing files
type FilesListOptions struct {
	CourseID     int64
	GroupID      int64
	FolderID     int64
	UserID       int64
	ContentTypes []string
	SearchTerm   string
	Include      []string
	Sort         string
	Order        string
}

// Validate validates the options
func (o *FilesListOptions) Validate() error {
	contextsSpecified := 0
	if o.CourseID > 0 {
		contextsSpecified++
	}
	if o.GroupID > 0 {
		contextsSpecified++
	}
	if o.FolderID > 0 {
		contextsSpecified++
	}
	if o.UserID > 0 {
		contextsSpecified++
	}

	if contextsSpecified == 0 {
		return fmt.Errorf("must specify one of course-id, group-id, folder-id, or user-id")
	}
	if contextsSpecified > 1 {
		return fmt.Errorf("can only specify one of course-id, group-id, folder-id, or user-id")
	}

	return nil
}

// FilesGetOptions contains options for getting a file
type FilesGetOptions struct {
	FileID  int64
	Include []string
}

// Validate validates the options
func (o *FilesGetOptions) Validate() error {
	return ValidateRequired("file-id", o.FileID)
}

// FilesUploadOptions contains options for uploading a file
type FilesUploadOptions struct {
	FilePath       string
	CourseID       int64
	GroupID        int64
	FolderID       int64
	UserID         int64
	OnDuplicate    string
	ParentFolderID int64
	Hidden         bool
	Locked         bool
}

// Validate validates the options
func (o *FilesUploadOptions) Validate() error {
	if o.FilePath == "" {
		return fmt.Errorf("file-path is required")
	}

	contextsSpecified := 0
	if o.CourseID > 0 {
		contextsSpecified++
	}
	if o.GroupID > 0 {
		contextsSpecified++
	}
	if o.FolderID > 0 {
		contextsSpecified++
	}
	if o.UserID > 0 {
		contextsSpecified++
	}

	if contextsSpecified == 0 {
		return fmt.Errorf("must specify one of course-id, group-id, folder-id, or user-id")
	}
	if contextsSpecified > 1 {
		return fmt.Errorf("can only specify one of course-id, group-id, folder-id, or user-id")
	}

	return nil
}

// FilesDownloadOptions contains options for downloading a file
type FilesDownloadOptions struct {
	FileID      int64
	Destination string
	Overwrite   bool
}

// Validate validates the options
func (o *FilesDownloadOptions) Validate() error {
	return ValidateRequired("file-id", o.FileID)
}

// FilesDeleteOptions contains options for deleting a file
type FilesDeleteOptions struct {
	FileID int64
	Force  bool
}

// Validate validates the options
func (o *FilesDeleteOptions) Validate() error {
	return ValidateRequired("file-id", o.FileID)
}

// FilesQuotaOptions contains options for getting quota information
type FilesQuotaOptions struct {
	CourseID int64
	GroupID  int64
	UserID   int64
}

// Validate validates the options
func (o *FilesQuotaOptions) Validate() error {
	positive := 0
	if o.CourseID > 0 {
		positive++
	}
	if o.GroupID > 0 {
		positive++
	}
	if o.UserID > 0 {
		positive++
	}
	if positive == 0 {
		return fmt.Errorf("must specify one of course-id, group-id, or user-id")
	}
	if positive > 1 {
		return fmt.Errorf("can only specify one of course-id, group-id, or user-id")
	}
	return nil
}

// FilesResetVerifierOptions contains options for resetting a file's link verifier
type FilesResetVerifierOptions struct {
	FileID int64
}

// Validate validates the options
func (o *FilesResetVerifierOptions) Validate() error {
	return ValidateRequired("file-id", o.FileID)
}

// FilesCopyOptions contains options for copying a file into a folder
type FilesCopyOptions struct {
	DestFolderID int64
	SourceFileID int64
	OnDuplicate  string
}

// Validate validates the options
func (o *FilesCopyOptions) Validate() error {
	if err := ValidateRequired("dest-folder-id", o.DestFolderID); err != nil {
		return err
	}
	return ValidateRequired("source-file-id", o.SourceFileID)
}

// FilesUsageRightsOptions contains options for setting/removing usage rights
type FilesUsageRightsOptions struct {
	CourseID         int64
	GroupID          int64
	UserID           int64
	FileIDs          []int64
	FolderIDs        []int64
	UseJustification string
	LegalCopyright   string
	License          string
	Publish          bool
}

// Validate validates the options
func (o *FilesUsageRightsOptions) Validate() error {
	positive := 0
	if o.CourseID > 0 {
		positive++
	}
	if o.GroupID > 0 {
		positive++
	}
	if o.UserID > 0 {
		positive++
	}
	if positive == 0 {
		return fmt.Errorf("must specify one of course-id, group-id, or user-id")
	}
	if positive > 1 {
		return fmt.Errorf("can only specify one of course-id, group-id, or user-id")
	}
	if len(o.FileIDs) == 0 && len(o.FolderIDs) == 0 {
		return fmt.Errorf("at least one --file-id or --folder-id is required")
	}
	if o.UseJustification == "" {
		return fmt.Errorf("use-justification is required")
	}
	return nil
}

// FilesRemoveUsageRightsOptions contains options for removing usage rights
type FilesRemoveUsageRightsOptions struct {
	CourseID  int64
	GroupID   int64
	UserID    int64
	FileIDs   []int64
	FolderIDs []int64
}

// Validate validates the options
func (o *FilesRemoveUsageRightsOptions) Validate() error {
	positive := 0
	if o.CourseID > 0 {
		positive++
	}
	if o.GroupID > 0 {
		positive++
	}
	if o.UserID > 0 {
		positive++
	}
	if positive == 0 {
		return fmt.Errorf("must specify one of course-id, group-id, or user-id")
	}
	if positive > 1 {
		return fmt.Errorf("can only specify one of course-id, group-id, or user-id")
	}
	if len(o.FileIDs) == 0 && len(o.FolderIDs) == 0 {
		return fmt.Errorf("at least one --file-id or --folder-id is required")
	}
	return nil
}

// FilesLicensesOptions contains options for listing content licenses
type FilesLicensesOptions struct {
	CourseID int64
	GroupID  int64
	UserID   int64
}

// Validate validates the options
func (o *FilesLicensesOptions) Validate() error {
	positive := 0
	if o.CourseID > 0 {
		positive++
	}
	if o.GroupID > 0 {
		positive++
	}
	if o.UserID > 0 {
		positive++
	}
	if positive == 0 {
		return fmt.Errorf("must specify one of course-id, group-id, or user-id")
	}
	if positive > 1 {
		return fmt.Errorf("can only specify one of course-id, group-id, or user-id")
	}
	return nil
}
