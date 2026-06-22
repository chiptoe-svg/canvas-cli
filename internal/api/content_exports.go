package api

import (
	"context"
	"fmt"
)

// ContentExportsService handles content export API calls for courses.
type ContentExportsService struct {
	client *Client
}

// NewContentExportsService creates a new content exports service.
func NewContentExportsService(client *Client) *ContentExportsService {
	return &ContentExportsService{client: client}
}

// ContentExport represents a Canvas content export.
type ContentExport struct {
	ID            int64              `json:"id"`
	CreatedAt     string             `json:"created_at"`
	ExportType    string             `json:"export_type"`
	ProgressURL   string             `json:"progress_url,omitempty"`
	UserID        int64              `json:"user_id,omitempty"`
	Attachment    *ContentExportFile `json:"attachment,omitempty"`
	WorkflowState string             `json:"workflow_state"`
	CourseID      int64              `json:"course_id,omitempty"`
}

// ContentExportFile holds the download link for a finished export.
type ContentExportFile struct {
	URL string `json:"url"`
}

// EpubExport represents a Canvas epub export.
type EpubExport struct {
	ID            int64              `json:"id"`
	CreatedAt     string             `json:"created_at,omitempty"`
	WorkflowState string             `json:"workflow_state"`
	Attachment    *ContentExportFile `json:"attachment,omitempty"`
	CourseID      int64              `json:"course_id,omitempty"`
	UserID        int64              `json:"user_id,omitempty"`
}

// ListContentExports lists content exports for a course.
// GET /api/v1/courses/:course_id/content_exports
func (s *ContentExportsService) List(ctx context.Context, courseID int64) ([]ContentExport, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/content_exports", courseID)
	var out []ContentExport
	if err := s.client.GetAllPages(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("listing content exports for course %d: %w", courseID, err)
	}
	return out, nil
}

// GetContentExport retrieves a single content export.
// GET /api/v1/courses/:course_id/content_exports/:id
func (s *ContentExportsService) Get(ctx context.Context, courseID, exportID int64) (*ContentExport, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/content_exports/%d", courseID, exportID)
	var out ContentExport
	if err := s.client.GetJSON(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("getting content export %d for course %d: %w", exportID, courseID, err)
	}
	return &out, nil
}

// CreateContentExportParams holds parameters for creating a content export.
type CreateContentExportParams struct {
	// ExportType is one of "common_cartridge", "zip", "qti", "course_copy".
	ExportType        string `json:"export_type"`
	SkipNotifications bool   `json:"skip_notifications,omitempty"`
}

// CreateContentExport starts a new content export.
// POST /api/v1/courses/:course_id/content_exports
func (s *ContentExportsService) Create(ctx context.Context, courseID int64, params CreateContentExportParams) (*ContentExport, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/content_exports", courseID)
	var out ContentExport
	if err := s.client.PostJSON(ctx, path, params, &out); err != nil {
		return nil, fmt.Errorf("creating content export for course %d: %w", courseID, err)
	}
	return &out, nil
}

// CreateEpubExport starts an epub export for a course.
// POST /api/v1/courses/:course_id/epub_exports
func (s *ContentExportsService) CreateEpub(ctx context.Context, courseID int64) (*EpubExport, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/epub_exports", courseID)
	var out EpubExport
	if err := s.client.PostJSON(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("creating epub export for course %d: %w", courseID, err)
	}
	return &out, nil
}

// GetEpubExport retrieves an epub export by ID.
// GET /api/v1/courses/:course_id/epub_exports/:id
func (s *ContentExportsService) GetEpub(ctx context.Context, courseID, epubID int64) (*EpubExport, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/epub_exports/%d", courseID, epubID)
	var out EpubExport
	if err := s.client.GetJSON(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("getting epub export %d for course %d: %w", epubID, courseID, err)
	}
	return &out, nil
}
