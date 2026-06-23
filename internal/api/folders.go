package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Folder represents a Canvas file folder
type Folder struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	FullName       string    `json:"full_name"`
	ContextID      int64     `json:"context_id"`
	ContextType    string    `json:"context_type"`
	ParentFolderID int64     `json:"parent_folder_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	LockAt         *string   `json:"lock_at,omitempty"`
	UnlockAt       *string   `json:"unlock_at,omitempty"`
	Position       int       `json:"position"`
	Locked         bool      `json:"locked"`
	LockedForUser  bool      `json:"locked_for_user"`
	Hidden         bool      `json:"hidden"`
	HiddenForUser  bool      `json:"hidden_for_user"`
	ForSubmissions bool      `json:"for_submissions"`
	FilesCount     int       `json:"files_count"`
	FoldersCount   int       `json:"folders_count"`
	FilesURL       string    `json:"files_url"`
	FoldersURL     string    `json:"folders_url"`
}

// FoldersService handles folder-related API calls
type FoldersService struct {
	client *Client
}

// NewFoldersService creates a new folders service
func NewFoldersService(client *Client) *FoldersService {
	return &FoldersService{client: client}
}

// foldersContextPath returns the context segment for course, group, or user.
// Exactly one ID must be positive.
func foldersContextPath(courseID, groupID, userID int64) (string, error) {
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
		return "", fmt.Errorf("one of course-id, group-id, or user-id is required")
	}
	if positive > 1 {
		return "", fmt.Errorf("specify only one of course-id, group-id, or user-id")
	}
	switch {
	case courseID > 0:
		return fmt.Sprintf("courses/%d", courseID), nil
	case groupID > 0:
		return fmt.Sprintf("groups/%d", groupID), nil
	default:
		return fmt.Sprintf("users/%d", userID), nil
	}
}

// ListFoldersOptions holds options for listing folders
type ListFoldersOptions struct {
	Page    int
	PerPage int
}

// ListContextFolders lists all folders for a course, group, or user context.
func (s *FoldersService) ListContextFolders(ctx context.Context, courseID, groupID, userID int64, opts *ListFoldersOptions) ([]Folder, error) {
	ctxSeg, err := foldersContextPath(courseID, groupID, userID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/%s/folders", ctxSeg)

	if opts != nil {
		query := url.Values{}
		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}
		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var folders []Folder
	if err := s.client.GetAllPages(ctx, path, &folders); err != nil {
		return nil, err
	}

	return folders, nil
}

// ListFolderSubFolders lists sub-folders within a specific folder.
func (s *FoldersService) ListFolderSubFolders(ctx context.Context, folderID int64, opts *ListFoldersOptions) ([]Folder, error) {
	if err := ValidatePositiveID(folderID, "folder_id"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/folders/%d/folders", folderID)

	if opts != nil {
		query := url.Values{}
		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}
		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var folders []Folder
	if err := s.client.GetAllPages(ctx, path, &folders); err != nil {
		return nil, err
	}

	return folders, nil
}

// GetFolder returns details for a specific folder by ID.
func (s *FoldersService) GetFolder(ctx context.Context, folderID int64) (*Folder, error) {
	if err := ValidatePositiveID(folderID, "folder_id"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/folders/%d", folderID)

	var folder Folder
	if err := s.client.GetJSON(ctx, path, &folder); err != nil {
		return nil, err
	}

	return &folder, nil
}

// ResolvePath resolves a full folder path and returns the folder hierarchy.
// fullPath may be empty (resolves to root) or a slash-separated path like "subfolder/child".
func (s *FoldersService) ResolvePath(ctx context.Context, courseID, groupID, userID int64, fullPath string) ([]Folder, error) {
	ctxSeg, err := foldersContextPath(courseID, groupID, userID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/%s/folders/by_path", ctxSeg)
	if fullPath != "" {
		path += "/" + fullPath
	}

	var folders []Folder
	if err := s.client.GetJSON(ctx, path, &folders); err != nil {
		return nil, err
	}

	return folders, nil
}

// CreateFolderParams holds parameters for creating a folder
type CreateFolderParams struct {
	Name           string
	ParentFolderID int64
	LockAt         string
	UnlockAt       string
	Locked         bool
	Hidden         bool
	Position       int
}

// CreateContextFolder creates a folder within a course, group, or user context.
func (s *FoldersService) CreateContextFolder(ctx context.Context, courseID, groupID, userID int64, params *CreateFolderParams) (*Folder, error) {
	if params == nil {
		return nil, ErrNilParams
	}
	if err := ValidateNonEmpty(params.Name, "name"); err != nil {
		return nil, err
	}

	ctxSeg, err := foldersContextPath(courseID, groupID, userID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/%s/folders", ctxSeg)
	body := buildFolderBody(params)

	var folder Folder
	if err := s.client.PostJSON(ctx, path, body, &folder); err != nil {
		return nil, err
	}

	return &folder, nil
}

// CreateSubFolder creates a folder inside a parent folder.
func (s *FoldersService) CreateSubFolder(ctx context.Context, parentFolderID int64, params *CreateFolderParams) (*Folder, error) {
	if err := ValidatePositiveID(parentFolderID, "folder_id"); err != nil {
		return nil, err
	}
	if params == nil {
		return nil, ErrNilParams
	}
	if err := ValidateNonEmpty(params.Name, "name"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/folders/%d/folders", parentFolderID)
	body := buildFolderBody(params)

	var folder Folder
	if err := s.client.PostJSON(ctx, path, body, &folder); err != nil {
		return nil, err
	}

	return &folder, nil
}

// buildFolderBody converts CreateFolderParams to a request body map.
func buildFolderBody(params *CreateFolderParams) map[string]interface{} {
	body := map[string]interface{}{
		"name": params.Name,
	}
	if params.ParentFolderID > 0 {
		body["parent_folder_id"] = params.ParentFolderID
	}
	if params.LockAt != "" {
		body["lock_at"] = params.LockAt
	}
	if params.UnlockAt != "" {
		body["unlock_at"] = params.UnlockAt
	}
	if params.Locked {
		body["locked"] = true
	}
	if params.Hidden {
		body["hidden"] = true
	}
	if params.Position > 0 {
		body["position"] = params.Position
	}
	return body
}

// UpdateFolderParams holds parameters for updating a folder
type UpdateFolderParams struct {
	Name           *string
	ParentFolderID *int64
	LockAt         *string
	UnlockAt       *string
	Locked         *bool
	Hidden         *bool
	Position       *int
}

// UpdateFolder updates an existing folder.
func (s *FoldersService) UpdateFolder(ctx context.Context, folderID int64, params *UpdateFolderParams) (*Folder, error) {
	if err := ValidatePositiveID(folderID, "folder_id"); err != nil {
		return nil, err
	}
	if params == nil {
		return nil, ErrNilParams
	}

	path := fmt.Sprintf("/api/v1/folders/%d", folderID)

	body := make(map[string]interface{})
	if params.Name != nil {
		body["name"] = *params.Name
	}
	if params.ParentFolderID != nil {
		body["parent_folder_id"] = *params.ParentFolderID
	}
	if params.LockAt != nil {
		body["lock_at"] = *params.LockAt
	}
	if params.UnlockAt != nil {
		body["unlock_at"] = *params.UnlockAt
	}
	if params.Locked != nil {
		body["locked"] = *params.Locked
	}
	if params.Hidden != nil {
		body["hidden"] = *params.Hidden
	}
	if params.Position != nil {
		body["position"] = *params.Position
	}

	var folder Folder
	if err := s.client.PutJSON(ctx, path, body, &folder); err != nil {
		return nil, err
	}

	return &folder, nil
}

// DeleteFolder deletes a folder. Set force=true to delete non-empty folders.
func (s *FoldersService) DeleteFolder(ctx context.Context, folderID int64, force bool) error {
	if err := ValidatePositiveID(folderID, "folder_id"); err != nil {
		return err
	}

	path := fmt.Sprintf("/api/v1/folders/%d", folderID)
	if force {
		path += "?force=true"
	}

	_, err := s.client.Delete(ctx, path)
	return err
}

// GetMediaFolder returns the designated upload folder for a course or group.
func (s *FoldersService) GetMediaFolder(ctx context.Context, courseID, groupID int64) (*Folder, error) {
	positive := 0
	if courseID > 0 {
		positive++
	}
	if groupID > 0 {
		positive++
	}
	if positive == 0 {
		return nil, fmt.Errorf("course-id or group-id is required")
	}
	if positive > 1 {
		return nil, fmt.Errorf("specify only one of course-id or group-id")
	}

	var path string
	if courseID > 0 {
		path = fmt.Sprintf("/api/v1/courses/%d/folders/media", courseID)
	} else {
		path = fmt.Sprintf("/api/v1/groups/%d/folders/media", groupID)
	}

	var folder Folder
	if err := s.client.GetJSON(ctx, path, &folder); err != nil {
		return nil, err
	}

	return &folder, nil
}

// CopyFolderParams holds parameters for copying a folder
type CopyFolderParams struct {
	SourceFolderID int64
}

// CopyFolder copies a folder into a destination folder.
func (s *FoldersService) CopyFolder(ctx context.Context, destFolderID int64, params *CopyFolderParams) (*Folder, error) {
	if err := ValidatePositiveID(destFolderID, "dest_folder_id"); err != nil {
		return nil, err
	}
	if params == nil {
		return nil, ErrNilParams
	}
	if err := ValidatePositiveID(params.SourceFolderID, "source_folder_id"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/folders/%d/copy_folder", destFolderID)
	body := map[string]interface{}{
		"source_folder_id": params.SourceFolderID,
	}

	var folder Folder
	if err := s.client.PostJSON(ctx, path, body, &folder); err != nil {
		return nil, err
	}

	return &folder, nil
}

// GetGroupFolder retrieves a folder by ID in a group context
// Canvas path: GET /api/v1/groups/:group_id/folders/:id
func (s *FoldersService) GetGroupFolder(ctx context.Context, groupID, folderID int64) (*Folder, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/folders/%d", groupID, folderID)

	var folder Folder
	if err := s.client.GetJSON(ctx, path, &folder); err != nil {
		return nil, err
	}

	return &folder, nil
}

// ListFolderFiles retrieves files in a folder
// Canvas path: GET /api/v1/folders/:id/files
func (s *FoldersService) ListFolderFiles(ctx context.Context, folderID int64) ([]Attachment, error) {
	path := fmt.Sprintf("/api/v1/folders/%d/files", folderID)

	var files []Attachment
	if err := s.client.GetAllPages(ctx, path, &files); err != nil {
		return nil, err
	}

	return files, nil
}
