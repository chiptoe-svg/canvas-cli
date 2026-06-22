package api

import (
	"context"
	"fmt"
)

// ProgressJob represents a Canvas background job progress object.
// This is distinct from the Progress type in types.go, which represents
// course-completion progress.
type ProgressJob struct {
	ID            int64       `json:"id"`
	Context       string      `json:"context_type,omitempty"`
	ContextID     int64       `json:"context_id,omitempty"`
	UserID        int64       `json:"user_id,omitempty"`
	Tag           string      `json:"tag,omitempty"`
	Completion    float64     `json:"completion,omitempty"`
	WorkflowState string      `json:"workflow_state"`
	CreatedAt     string      `json:"created_at,omitempty"`
	UpdatedAt     string      `json:"updated_at,omitempty"`
	Message       string      `json:"message,omitempty"`
	Results       interface{} `json:"results,omitempty"`
	URL           string      `json:"url,omitempty"`
}

// ProgressService handles progress-related API calls.
type ProgressService struct {
	client *Client
}

// NewProgressService creates a new ProgressService.
func NewProgressService(client *Client) *ProgressService {
	return &ProgressService{client: client}
}

// Get retrieves a single progress job by ID.
func (s *ProgressService) Get(ctx context.Context, id int64) (*ProgressJob, error) {
	path := fmt.Sprintf("/api/v1/progress/%d", id)

	var p ProgressJob
	if err := s.client.GetJSON(ctx, path, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

// Cancel cancels an in-progress background job.
func (s *ProgressService) Cancel(ctx context.Context, id int64) (*ProgressJob, error) {
	path := fmt.Sprintf("/api/v1/progress/%d/cancel", id)

	var p ProgressJob
	if err := s.client.PostJSON(ctx, path, nil, &p); err != nil {
		return nil, err
	}

	return &p, nil
}
