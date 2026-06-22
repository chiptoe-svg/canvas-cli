package options

import "fmt"

// foldersValidateContext validates that exactly one of courseID, groupID, or userID is set.
func foldersValidateContext(courseID, groupID, userID int64) error {
	positive := 0
	if courseID > 0 {
		positive++
	}
	if groupID > 0 {
		positive++
	}
	if userID > 0 {
		positive++
	}
	if positive == 0 {
		return fmt.Errorf("one of --course-id, --group-id, or --user-id is required")
	}
	if positive > 1 {
		return fmt.Errorf("specify only one of --course-id, --group-id, or --user-id")
	}
	return nil
}

// FoldersListOptions contains options for listing folders
type FoldersListOptions struct {
	CourseID int64
	GroupID  int64
	UserID   int64
	FolderID int64 // list sub-folders of a specific folder
}

// Validate validates the options
func (o *FoldersListOptions) Validate() error {
	if o.FolderID > 0 {
		// --folder-id takes precedence; no context required
		return nil
	}
	return foldersValidateContext(o.CourseID, o.GroupID, o.UserID)
}

// FoldersGetOptions contains options for getting a folder
type FoldersGetOptions struct {
	FolderID int64
}

// Validate validates the options
func (o *FoldersGetOptions) Validate() error {
	return ValidateRequired("folder-id", o.FolderID)
}

// FoldersResolvePathOptions contains options for resolving a folder path
type FoldersResolvePathOptions struct {
	CourseID int64
	GroupID  int64
	UserID   int64
	Path     string
}

// Validate validates the options
func (o *FoldersResolvePathOptions) Validate() error {
	return foldersValidateContext(o.CourseID, o.GroupID, o.UserID)
}

// FoldersCreateOptions contains options for creating a folder
type FoldersCreateOptions struct {
	CourseID       int64
	GroupID        int64
	UserID         int64
	ParentFolderID int64
	Name           string
	LockAt         string
	UnlockAt       string
	Locked         bool
	Hidden         bool
	Position       int
}

// Validate validates the options
func (o *FoldersCreateOptions) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("name is required")
	}
	if o.ParentFolderID > 0 {
		// Creating inside a known folder; no context required
		return nil
	}
	return foldersValidateContext(o.CourseID, o.GroupID, o.UserID)
}

// FoldersUpdateOptions contains options for updating a folder
type FoldersUpdateOptions struct {
	FolderID       int64
	Name           string
	ParentFolderID int64
	LockAt         string
	UnlockAt       string
	Locked         bool
	Hidden         bool
	Position       int
	// Track which fields were set
	NameSet           bool
	ParentFolderIDSet bool
	LockAtSet         bool
	UnlockAtSet       bool
	LockedSet         bool
	HiddenSet         bool
	PositionSet       bool
}

// Validate validates the options
func (o *FoldersUpdateOptions) Validate() error {
	return ValidateRequired("folder-id", o.FolderID)
}

// FoldersDeleteOptions contains options for deleting a folder
type FoldersDeleteOptions struct {
	FolderID int64
	Force    bool
}

// Validate validates the options
func (o *FoldersDeleteOptions) Validate() error {
	return ValidateRequired("folder-id", o.FolderID)
}

// FoldersMediaOptions contains options for getting the media upload folder
type FoldersMediaOptions struct {
	CourseID int64
	GroupID  int64
}

// Validate validates the options
func (o *FoldersMediaOptions) Validate() error {
	if o.CourseID <= 0 && o.GroupID <= 0 {
		return fmt.Errorf("one of --course-id or --group-id is required")
	}
	if o.CourseID > 0 && o.GroupID > 0 {
		return fmt.Errorf("specify only one of --course-id or --group-id")
	}
	return nil
}

// FoldersCopyOptions contains options for copying a folder
type FoldersCopyOptions struct {
	DestFolderID   int64
	SourceFolderID int64
}

// Validate validates the options
func (o *FoldersCopyOptions) Validate() error {
	if err := ValidateRequired("dest-folder-id", o.DestFolderID); err != nil {
		return err
	}
	return ValidateRequired("source-folder-id", o.SourceFolderID)
}
