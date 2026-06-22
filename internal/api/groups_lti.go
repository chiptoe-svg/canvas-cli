package api

import (
	"context"
	"fmt"
)

// GroupsLTIService handles LTI-specific group API calls.
// These endpoints live under /api/lti/ rather than /api/v1/.
type GroupsLTIService struct {
	client *Client
}

// NewGroupsLTIService creates a new LTI groups service.
func NewGroupsLTIService(client *Client) *GroupsLTIService {
	return &GroupsLTIService{client: client}
}

// LTIGroupMember represents a member entry in the LTI Names and Roles response.
type LTIGroupMember struct {
	Status string   `json:"status"`
	Name   string   `json:"name,omitempty"`
	Email  string   `json:"email,omitempty"`
	UserID string   `json:"user_id,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

// LTINamesAndRolesResponse is the envelope returned by the Names and Roles endpoint.
type LTINamesAndRolesResponse struct {
	ID      string                 `json:"id"`
	Context map[string]interface{} `json:"context,omitempty"`
	Members []LTIGroupMember       `json:"members"`
}

// GetNamesAndRoles returns the LTI Names and Roles Provisioning Service response for a group.
// GET /api/lti/groups/:group_id/names_and_roles
func (s *GroupsLTIService) GetNamesAndRoles(ctx context.Context, groupID int64) (*LTINamesAndRolesResponse, error) {
	path := fmt.Sprintf("/api/lti/groups/%d/names_and_roles", groupID)

	var result LTINamesAndRolesResponse
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// LTIUser represents a user in the LTI users list.
type LTIUser struct {
	ID           string `json:"id"`
	LoginID      string `json:"login_id,omitempty"`
	Name         string `json:"name,omitempty"`
	SortableName string `json:"sortable_name,omitempty"`
}

// ListLTIUsers returns the users in a group via the LTI endpoint.
// GET /api/lti/groups/:group_id/users
func (s *GroupsLTIService) ListLTIUsers(ctx context.Context, groupID int64) ([]LTIUser, error) {
	path := fmt.Sprintf("/api/lti/groups/%d/users", groupID)

	var users []LTIUser
	if err := s.client.GetAllPages(ctx, path, &users); err != nil {
		return nil, err
	}

	return users, nil
}
