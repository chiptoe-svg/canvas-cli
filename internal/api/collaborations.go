package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// Collaboration represents a Canvas collaboration tool entry (Google Docs, etc.).
type Collaboration struct {
	ID                int64  `json:"id"`
	CollaborationType string `json:"collaboration_type,omitempty"`
	DocumentID        string `json:"document_id,omitempty"`
	UserID            int64  `json:"user_id,omitempty"`
	ContextID         int64  `json:"context_id,omitempty"`
	ContextType       string `json:"context_type,omitempty"`
	URL               string `json:"url,omitempty"`
	CreatedAt         string `json:"created_at,omitempty"`
	UpdatedAt         string `json:"updated_at,omitempty"`
	Description       string `json:"description,omitempty"`
	Title             string `json:"title,omitempty"`
}

// Collaborator represents a user or group that is a member of a collaboration.
type Collaborator struct {
	ID   int64  `json:"id"`
	Type string `json:"type,omitempty"` // user, group
	Name string `json:"name,omitempty"`
}

// CollaborationsService handles collaboration-related API calls.
type CollaborationsService struct {
	client *Client
}

// NewCollaborationsService creates a new CollaborationsService.
func NewCollaborationsService(client *Client) *CollaborationsService {
	return &CollaborationsService{client: client}
}

// ListCollaborationsOptions holds query parameters for listing collaborations.
type ListCollaborationsOptions struct {
	PerPage int
}

func buildCollaborationsQuery(opts *ListCollaborationsOptions) string {
	if opts == nil {
		return ""
	}
	q := url.Values{}
	if opts.PerPage > 0 {
		q.Set("per_page", strconv.Itoa(opts.PerPage))
	}
	if len(q) > 0 {
		return "?" + q.Encode()
	}
	return ""
}

// ListForCourse retrieves collaborations for a course.
func (s *CollaborationsService) ListForCourse(ctx context.Context, courseID int64, opts *ListCollaborationsOptions) ([]Collaboration, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/collaborations", courseID) + buildCollaborationsQuery(opts)

	var colls []Collaboration
	if err := s.client.GetAllPages(ctx, path, &colls); err != nil {
		return nil, err
	}

	return colls, nil
}

// ListForGroup retrieves collaborations for a group.
func (s *CollaborationsService) ListForGroup(ctx context.Context, groupID int64, opts *ListCollaborationsOptions) ([]Collaboration, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/collaborations", groupID) + buildCollaborationsQuery(opts)

	var colls []Collaboration
	if err := s.client.GetAllPages(ctx, path, &colls); err != nil {
		return nil, err
	}

	return colls, nil
}

// ListMembers retrieves the members (collaborators) of a collaboration.
func (s *CollaborationsService) ListMembers(ctx context.Context, collaborationID int64, opts *ListCollaborationsOptions) ([]Collaborator, error) {
	path := fmt.Sprintf("/api/v1/collaborations/%d/members", collaborationID) + buildCollaborationsQuery(opts)

	var members []Collaborator
	if err := s.client.GetAllPages(ctx, path, &members); err != nil {
		return nil, err
	}

	return members, nil
}
