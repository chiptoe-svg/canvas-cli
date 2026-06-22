package api

import (
	"context"
	"fmt"
)

// EpubExport represents a Canvas ePub export job.
type EpubExport struct {
	ID            int64       `json:"id"`
	CreatedAt     string      `json:"created_at,omitempty"`
	Course        interface{} `json:"course,omitempty"`
	DownloadURL   string      `json:"epub_export,omitempty"`
	WorkflowState string      `json:"workflow_state,omitempty"`
	User          *User       `json:"user,omitempty"`
}

// EpubExportsService handles ePub export API calls.
type EpubExportsService struct {
	client *Client
}

// NewEpubExportsService creates a new EpubExportsService.
func NewEpubExportsService(client *Client) *EpubExportsService {
	return &EpubExportsService{client: client}
}

// List retrieves all ePub exports (with associated courses) for the current user.
func (s *EpubExportsService) List(ctx context.Context) ([]EpubExport, error) {
	path := "/api/v1/epub_exports"

	var exports []EpubExport
	if err := s.client.GetAllPages(ctx, path, &exports); err != nil {
		return nil, err
	}

	return exports, nil
}

// Create starts an ePub export for a course.
func (s *EpubExportsService) Create(ctx context.Context, courseID int64) (*EpubExport, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/epub_exports", courseID)

	var export EpubExport
	if err := s.client.PostJSON(ctx, path, nil, &export); err != nil {
		return nil, err
	}

	return &export, nil
}

// Get retrieves the status of an ePub export.
func (s *EpubExportsService) Get(ctx context.Context, courseID, exportID int64) (*EpubExport, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/epub_exports/%d", courseID, exportID)

	var export EpubExport
	if err := s.client.GetJSON(ctx, path, &export); err != nil {
		return nil, err
	}

	return &export, nil
}
