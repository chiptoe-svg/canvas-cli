package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// FilesService handles file-related API calls
type FilesService struct {
	client *Client
}

// isCanvasDomain checks if the given URL belongs to the same Canvas instance
// This is used to prevent leaking the Authorization header to third-party storage providers
func isCanvasDomain(redirectURL, baseURL string) bool {
	redirectParsed, err := url.Parse(redirectURL)
	if err != nil {
		return false
	}
	baseParsed, err := url.Parse(baseURL)
	if err != nil {
		return false
	}
	return redirectParsed.Host == baseParsed.Host
}

// NewFilesService creates a new files service
func NewFilesService(client *Client) *FilesService {
	return &FilesService{client: client}
}

// ListFilesOptions holds options for listing files
type ListFilesOptions struct {
	ContentTypes []string // Filter by MIME type
	SearchTerm   string   // Search by file name
	Include      []string // Additional data to include (user)
	Only         []string // Filter by type (names, folders)
	Sort         string   // Sort by (name, size, created_at, updated_at, content_type)
	Order        string   // Order direction (asc, desc)
	Page         int
	PerPage      int
}

// ListCourseFiles retrieves files for a course
func (s *FilesService) ListCourseFiles(ctx context.Context, courseID int64, opts *ListFilesOptions) ([]Attachment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/files", courseID)
	return s.listFiles(ctx, path, opts)
}

// ListFolderFiles retrieves files in a folder
func (s *FilesService) ListFolderFiles(ctx context.Context, folderID int64, opts *ListFilesOptions) ([]Attachment, error) {
	path := fmt.Sprintf("/api/v1/folders/%d/files", folderID)
	return s.listFiles(ctx, path, opts)
}

// ListUserFiles retrieves files for a user
func (s *FilesService) ListUserFiles(ctx context.Context, userID int64, opts *ListFilesOptions) ([]Attachment, error) {
	path := fmt.Sprintf("/api/v1/users/%d/files", userID)
	return s.listFiles(ctx, path, opts)
}

// listFiles is a helper for listing files with options
func (s *FilesService) listFiles(ctx context.Context, basePath string, opts *ListFilesOptions) ([]Attachment, error) {
	path := basePath

	if opts != nil {
		query := url.Values{}

		for _, ct := range opts.ContentTypes {
			query.Add("content_types[]", ct)
		}

		if opts.SearchTerm != "" {
			query.Add("search_term", opts.SearchTerm)
		}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		for _, only := range opts.Only {
			query.Add("only[]", only)
		}

		if opts.Sort != "" {
			query.Add("sort", opts.Sort)
		}

		if opts.Order != "" {
			query.Add("order", opts.Order)
		}

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

	var files []Attachment
	if err := s.client.GetAllPages(ctx, path, &files); err != nil {
		return nil, err
	}

	return files, nil
}

// Get retrieves a single file by ID
func (s *FilesService) Get(ctx context.Context, fileID int64, include []string) (*Attachment, error) {
	path := fmt.Sprintf("/api/v1/files/%d", fileID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var file Attachment
	if err := s.client.GetJSON(ctx, path, &file); err != nil {
		return nil, err
	}

	return &file, nil
}

// UploadParams holds parameters for uploading a file
type UploadParams struct {
	Name           string // File name
	Size           int64  // File size in bytes (required for Canvas)
	ContentType    string // MIME type
	ParentFolderID int64  // Folder to upload to
	OnDuplicate    string // How to handle duplicates: overwrite, rename
	LockAt         string // ISO8601 date
	UnlockAt       string // ISO8601 date
	Locked         bool   // Lock the file
	Hidden         bool   // Hide from students
}

// UploadToCourse uploads a file to a course
func (s *FilesService) UploadToCourse(ctx context.Context, courseID int64, filePath string, params *UploadParams) (*Attachment, error) {
	uploadPath := fmt.Sprintf("/api/v1/courses/%d/files", courseID)
	return s.upload(ctx, uploadPath, filePath, params)
}

// UploadToFolder uploads a file to a specific folder
func (s *FilesService) UploadToFolder(ctx context.Context, folderID int64, filePath string, params *UploadParams) (*Attachment, error) {
	uploadPath := fmt.Sprintf("/api/v1/folders/%d/files", folderID)
	return s.upload(ctx, uploadPath, filePath, params)
}

// UploadToUser uploads a file to a user's files
func (s *FilesService) UploadToUser(ctx context.Context, userID int64, filePath string, params *UploadParams) (*Attachment, error) {
	uploadPath := fmt.Sprintf("/api/v1/users/%d/files", userID)
	return s.upload(ctx, uploadPath, filePath, params)
}

// UploadToGroup uploads a file to a Canvas group's files.
func (s *FilesService) UploadToGroup(ctx context.Context, groupID int64, filePath string, params *UploadParams) (*Attachment, error) {
	uploadPath := fmt.Sprintf("/api/v1/groups/%d/files", groupID)
	return s.upload(ctx, uploadPath, filePath, params)
}

// upload is a helper that handles the Canvas three-step upload process
// Step 1: Request upload parameters from Canvas
// Step 2: Upload file to the provided URL using multipart/form-data with upload_params
// Step 3: Confirm upload (follow redirect or parse response)
func (s *FilesService) upload(ctx context.Context, uploadPath, filePath string, params *UploadParams) (*Attachment, error) {
	// Open the file to upload — filePath comes from the --file flag, user-controlled by design.
	file, err := os.Open(filePath) // #nosec G304 -- filePath is provided by the user via --file flag
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Use provided name or default to filename
	fileName := params.Name
	if fileName == "" {
		fileName = filepath.Base(filePath)
	}

	// Build upload request parameters
	uploadBody := map[string]interface{}{
		"name": fileName,
		"size": fileInfo.Size(),
	}

	if params.ContentType != "" {
		uploadBody["content_type"] = params.ContentType
	}
	if params.ParentFolderID > 0 {
		uploadBody["parent_folder_id"] = params.ParentFolderID
	}
	if params.OnDuplicate != "" {
		uploadBody["on_duplicate"] = params.OnDuplicate
	}
	if params.LockAt != "" {
		uploadBody["lock_at"] = params.LockAt
	}
	if params.UnlockAt != "" {
		uploadBody["unlock_at"] = params.UnlockAt
	}
	if params.Locked {
		uploadBody["locked"] = true
	}
	if params.Hidden {
		uploadBody["hidden"] = true
	}

	// Step 1: Tell Canvas we want to upload
	var uploadResponse struct {
		UploadURL    string                 `json:"upload_url"`
		UploadParams map[string]interface{} `json:"upload_params"`
		FileParam    string                 `json:"file_param"`
	}

	if err := s.client.PostJSON(ctx, uploadPath, uploadBody, &uploadResponse); err != nil {
		return nil, fmt.Errorf("failed to initialize upload: %w", err)
	}

	// Step 2: Upload the actual file to the provided URL using multipart/form-data
	// Create a buffer to write our multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add all upload params from Canvas first (order matters for some storage providers)
	for key, value := range uploadResponse.UploadParams {
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case float64:
			strValue = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			strValue = strconv.FormatBool(v)
		default:
			strValue = fmt.Sprintf("%v", v)
		}
		if err := writer.WriteField(key, strValue); err != nil {
			return nil, fmt.Errorf("failed to write form field %s: %w", key, err)
		}
	}

	// Determine the file field name (usually "file" but Canvas tells us)
	fileFieldName := "file"
	if uploadResponse.FileParam != "" {
		fileFieldName = uploadResponse.FileParam
	}

	// Add the file itself
	part, err := writer.CreateFormFile(fileFieldName, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Reset file position and copy content
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close the writer to finalize the multipart message
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create the upload request
	req, err := http.NewRequestWithContext(ctx, "POST", uploadResponse.UploadURL, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a custom client that doesn't follow redirects automatically
	// We need to handle the redirect ourselves to get the file confirmation
	noRedirectClient := &http.Client{
		Timeout: 5 * time.Minute, // File uploads may take longer
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Execute the upload
	resp, err := noRedirectClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	// Handle different response scenarios
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		// Direct response with file info
		var uploadedFile Attachment
		if err := json.NewDecoder(resp.Body).Decode(&uploadedFile); err != nil {
			return nil, fmt.Errorf("failed to parse upload response: %w", err)
		}
		return &uploadedFile, nil

	case http.StatusFound, http.StatusSeeOther, http.StatusMovedPermanently, http.StatusTemporaryRedirect:
		// Step 3: Follow redirect to confirm upload
		location := resp.Header.Get("Location")
		if location == "" {
			return nil, fmt.Errorf("redirect response missing Location header")
		}

		// Make a GET request to the redirect URL to confirm the upload
		confirmReq, err := http.NewRequestWithContext(ctx, "GET", location, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create confirmation request: %w", err)
		}

		// Add authorization header only if the redirect is to Canvas domain.
		// This prevents leaking the bearer token to third-party storage providers.
		// Use getToken() so OAuth token-source users (who may have an empty
		// s.client.token field) send a valid, possibly auto-refreshed token.
		if isCanvasDomain(location, s.client.baseURL) {
			tok, err := s.client.getToken()
			if err != nil {
				return nil, fmt.Errorf("failed to get token for upload confirmation: %w", err)
			}
			confirmReq.Header.Set("Authorization", "Bearer "+tok)
		}

		confirmResp, err := s.client.httpClient.Do(confirmReq)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm upload: %w", err)
		}
		defer confirmResp.Body.Close()

		if confirmResp.StatusCode != http.StatusOK && confirmResp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(confirmResp.Body)
			return nil, fmt.Errorf("upload confirmation failed with status %d: %s", confirmResp.StatusCode, string(bodyBytes))
		}

		var uploadedFile Attachment
		if err := json.NewDecoder(confirmResp.Body).Decode(&uploadedFile); err != nil {
			return nil, fmt.Errorf("failed to parse confirmation response: %w", err)
		}
		return &uploadedFile, nil

	default:
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}
}

// sanitizeDestPath validates a destination path to prevent path traversal.
// It accepts any caller-supplied path as-is (the caller controls it), but
// rejects the empty string so os.Create does not create an unnamed file.
// The primary path-traversal defence lives in the command layer which
// sanitizes server-controlled filenames before they reach here.
func sanitizeDestPath(destPath string) error {
	if destPath == "" {
		return fmt.Errorf("destination path must not be empty")
	}
	// Reject a bare "." which would create or truncate the current directory entry.
	if filepath.Clean(destPath) == "." {
		return fmt.Errorf("unsafe destination path: %q", destPath)
	}
	return nil
}

// Download downloads a file to the specified destination
func (s *FilesService) Download(ctx context.Context, fileID int64, destPath string) error {
	if err := sanitizeDestPath(destPath); err != nil {
		return err
	}

	// Get file info first to get the download URL
	file, err := s.Get(ctx, fileID, nil)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	if file.URL == "" {
		return fmt.Errorf("file has no download URL")
	}

	// Create the destination file — destPath comes from the --destination flag, user-controlled by design.
	out, err := os.Create(destPath) // #nosec G304 -- destPath is provided by the user via --destination flag
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	// Download the file content
	req, err := http.NewRequestWithContext(ctx, "GET", file.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Copy the content to the destination file
	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	return nil
}

// UpdateParams holds parameters for updating a file
type UpdateParams struct {
	Name           string  // New file name
	ParentFolderID *int64  // Move to different folder
	LockAt         *string // ISO8601 date
	UnlockAt       *string // ISO8601 date
	Locked         *bool
	Hidden         *bool
}

// Update updates file metadata
func (s *FilesService) Update(ctx context.Context, fileID int64, params *UpdateParams) (*Attachment, error) {
	path := fmt.Sprintf("/api/v1/files/%d", fileID)

	body := make(map[string]interface{})

	if params.Name != "" {
		body["name"] = params.Name
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

	var updatedFile Attachment
	if err := s.client.PutJSON(ctx, path, body, &updatedFile); err != nil {
		return nil, err
	}

	return &updatedFile, nil
}

// Delete deletes a file
func (s *FilesService) Delete(ctx context.Context, fileID int64) error {
	path := fmt.Sprintf("/api/v1/files/%d", fileID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// GetQuota retrieves quota information for a course or user
func (s *FilesService) GetCourseQuota(ctx context.Context, courseID int64) (*QuotaInfo, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/files/quota", courseID)
	var quota QuotaInfo
	if err := s.client.GetJSON(ctx, path, &quota); err != nil {
		return nil, err
	}
	return &quota, nil
}

// GetUserQuota retrieves quota information for a user
func (s *FilesService) GetUserQuota(ctx context.Context, userID int64) (*QuotaInfo, error) {
	path := fmt.Sprintf("/api/v1/users/%d/files/quota", userID)
	var quota QuotaInfo
	if err := s.client.GetJSON(ctx, path, &quota); err != nil {
		return nil, err
	}
	return &quota, nil
}

// QuotaInfo represents storage quota information
type QuotaInfo struct {
	QuotaUsed int64 `json:"quota_used"`
	Quota     int64 `json:"quota"`
}

// ListGroupFiles retrieves files for a group
func (s *FilesService) ListGroupFiles(ctx context.Context, groupID int64, opts *ListFilesOptions) ([]Attachment, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/files", groupID)
	return s.listFiles(ctx, path, opts)
}

// GetGroupQuota retrieves quota information for a group
func (s *FilesService) GetGroupQuota(ctx context.Context, groupID int64) (*QuotaInfo, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/files/quota", groupID)
	var quota QuotaInfo
	if err := s.client.GetJSON(ctx, path, &quota); err != nil {
		return nil, err
	}
	return &quota, nil
}

// ResetVerifier resets the link verifier for a file.
// Any existing links using the previous verifier parameter will no longer work.
func (s *FilesService) ResetVerifier(ctx context.Context, fileID int64) (*Attachment, error) {
	if err := ValidatePositiveID(fileID, "file_id"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/files/%d/reset_verifier", fileID)

	var file Attachment
	if err := s.client.PostJSON(ctx, path, nil, &file); err != nil {
		return nil, err
	}

	return &file, nil
}

// CopyFileParams holds parameters for copying a file into a folder
type CopyFileParams struct {
	SourceFileID int64
	OnDuplicate  string // overwrite | rename
}

// CopyFile copies a file from elsewhere in Canvas into a destination folder.
func (s *FilesService) CopyFile(ctx context.Context, destFolderID int64, params *CopyFileParams) (*Attachment, error) {
	if err := ValidatePositiveID(destFolderID, "dest_folder_id"); err != nil {
		return nil, err
	}
	if params == nil {
		return nil, ErrNilParams
	}
	if err := ValidatePositiveID(params.SourceFileID, "source_file_id"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/folders/%d/copy_file", destFolderID)

	body := map[string]interface{}{
		"source_file_id": params.SourceFileID,
	}
	if params.OnDuplicate != "" {
		body["on_duplicate"] = params.OnDuplicate
	}

	var file Attachment
	if err := s.client.PostJSON(ctx, path, body, &file); err != nil {
		return nil, err
	}

	return &file, nil
}

// UsageRights represents copyright and license information for files
type UsageRights struct {
	UseJustification string `json:"use_justification"`
	LegalCopyright   string `json:"legal_copyright,omitempty"`
	License          string `json:"license,omitempty"`
	LicenseName      string `json:"license_name,omitempty"`
	Message          string `json:"message,omitempty"`
	Restricted       bool   `json:"restricted,omitempty"`
}

// SetUsageRightsParams holds parameters for setting usage rights on files
type SetUsageRightsParams struct {
	FileIDs          []int64
	FolderIDs        []int64
	Publish          bool
	UseJustification string // required: own_copyright|used_by_permission|fair_use|public_domain|creative_commons
	LegalCopyright   string
	License          string // required if UseJustification == "creative_commons"
}

// ContentLicense represents an available content license
type ContentLicense struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// filesUsageRightsContextPath returns the context path segment for usage rights.
// Exactly one of courseID, groupID, userID must be positive.
func filesUsageRightsContextPath(courseID, groupID, userID int64) (string, error) {
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

// SetUsageRights sets copyright and license information on files.
func (s *FilesService) SetUsageRights(ctx context.Context, courseID, groupID, userID int64, params *SetUsageRightsParams) (*UsageRights, error) {
	if params == nil {
		return nil, ErrNilParams
	}
	if err := ValidateNonEmpty(params.UseJustification, "use_justification"); err != nil {
		return nil, err
	}
	if len(params.FileIDs) == 0 && len(params.FolderIDs) == 0 {
		return nil, fmt.Errorf("at least one file-id or folder-id is required")
	}

	ctxSeg, err := filesUsageRightsContextPath(courseID, groupID, userID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/%s/usage_rights", ctxSeg)

	body := map[string]interface{}{
		"usage_rights[use_justification]": params.UseJustification,
	}
	if len(params.FileIDs) > 0 {
		body["file_ids"] = params.FileIDs
	}
	if len(params.FolderIDs) > 0 {
		body["folder_ids"] = params.FolderIDs
	}
	if params.Publish {
		body["publish"] = true
	}
	if params.LegalCopyright != "" {
		body["usage_rights[legal_copyright]"] = params.LegalCopyright
	}
	if params.License != "" {
		body["usage_rights[license]"] = params.License
	}

	var rights UsageRights
	if err := s.client.PutJSON(ctx, path, body, &rights); err != nil {
		return nil, err
	}

	return &rights, nil
}

// RemoveUsageRightsParams holds parameters for removing usage rights
type RemoveUsageRightsParams struct {
	FileIDs   []int64
	FolderIDs []int64
}

// RemoveUsageRights removes copyright and license information from files.
func (s *FilesService) RemoveUsageRights(ctx context.Context, courseID, groupID, userID int64, params *RemoveUsageRightsParams) error {
	if params == nil {
		return ErrNilParams
	}
	if len(params.FileIDs) == 0 && len(params.FolderIDs) == 0 {
		return fmt.Errorf("at least one file-id or folder-id is required")
	}

	ctxSeg, err := filesUsageRightsContextPath(courseID, groupID, userID)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/api/v1/%s/usage_rights", ctxSeg)

	query := url.Values{}
	for _, id := range params.FileIDs {
		query.Add("file_ids[]", strconv.FormatInt(id, 10))
	}
	for _, id := range params.FolderIDs {
		query.Add("folder_ids[]", strconv.FormatInt(id, 10))
	}
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	_, err = s.client.Delete(ctx, path)
	return err
}

// ListLicenses returns the list of licenses that can be applied to files.
func (s *FilesService) ListLicenses(ctx context.Context, courseID, groupID, userID int64) ([]ContentLicense, error) {
	ctxSeg, err := filesUsageRightsContextPath(courseID, groupID, userID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/%s/content_licenses", ctxSeg)

	var licenses []ContentLicense
	if err := s.client.GetAllPages(ctx, path, &licenses); err != nil {
		return nil, err
	}

	return licenses, nil
}

// FileDateDetails represents date availability details for a file
type FileDateDetails struct {
	DueAt     *string              `json:"due_at,omitempty"`
	UnlockAt  *string              `json:"unlock_at,omitempty"`
	LockAt    *string              `json:"lock_at,omitempty"`
	Overrides []AssignmentOverride `json:"overrides,omitempty"`
}

// GetCourseDateDetails retrieves date details for a file in a course context
// Canvas path: GET /api/v1/courses/:course_id/files/:attachment_id/date_details
func (s *FilesService) GetCourseDateDetails(ctx context.Context, courseID, attachmentID int64) (*FileDateDetails, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/files/%d/date_details", courseID, attachmentID)

	var result FileDateDetails
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateCourseDateDetails updates date details for a file in a course context
// Canvas path: PUT /api/v1/courses/:course_id/files/:attachment_id/date_details
func (s *FilesService) UpdateCourseDateDetails(ctx context.Context, courseID, attachmentID int64, params *FileDateDetails) (*FileDateDetails, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/files/%d/date_details", courseID, attachmentID)

	var result FileDateDetails
	if err := s.client.PutJSON(ctx, path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// FilePublicURL represents the public URL for a file
type FilePublicURL struct {
	PublicURL string `json:"public_url"`
}

// GetPublicURL retrieves the public URL for a file
// Canvas path: GET /api/v1/files/:id/public_url
func (s *FilesService) GetPublicURL(ctx context.Context, fileID int64) (*FilePublicURL, error) {
	path := fmt.Sprintf("/api/v1/files/%d/public_url", fileID)

	var result FilePublicURL
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetIconMetadata retrieves icon metadata for a file
// Canvas path: GET /api/v1/files/:id/icon_metadata
func (s *FilesService) GetIconMetadata(ctx context.Context, fileID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/files/%d/icon_metadata", fileID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetGroupFile retrieves a file by ID in a group context
// Canvas path: GET /api/v1/groups/:group_id/files/:id
func (s *FilesService) GetGroupFile(ctx context.Context, groupID, fileID int64) (*Attachment, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/files/%d", groupID, fileID)

	var file Attachment
	if err := s.client.GetJSON(ctx, path, &file); err != nil {
		return nil, err
	}

	return &file, nil
}

// GetCourseFile retrieves a file by ID in a course context
// Canvas path: GET /api/v1/courses/:course_id/files/:id
func (s *FilesService) GetCourseFile(ctx context.Context, courseID, fileID int64) (*Attachment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/files/%d", courseID, fileID)

	var file Attachment
	if err := s.client.GetJSON(ctx, path, &file); err != nil {
		return nil, err
	}

	return &file, nil
}

// GetFileRef retrieves a file by migration ID in a course context
// Canvas path: GET /api/v1/courses/:course_id/files/file_ref/:migration_id
func (s *FilesService) GetFileRef(ctx context.Context, courseID int64, migrationID string) (*Attachment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/files/file_ref/%s", courseID, migrationID)

	var file Attachment
	if err := s.client.GetJSON(ctx, path, &file); err != nil {
		return nil, err
	}

	return &file, nil
}
