package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// GroupsService handles group-related API calls
type GroupsService struct {
	client *Client
}

// NewGroupsService creates a new groups service
func NewGroupsService(client *Client) *GroupsService {
	return &GroupsService{client: client}
}

// Group represents a Canvas group
type Group struct {
	ID              int64             `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description,omitempty"`
	IsPublic        bool              `json:"is_public"`
	FollowedByUser  bool              `json:"followed_by_user,omitempty"`
	JoinLevel       string            `json:"join_level,omitempty"`
	MembersCount    int               `json:"members_count,omitempty"`
	AvatarURL       string            `json:"avatar_url,omitempty"`
	ContextType     string            `json:"context_type,omitempty"`
	CourseID        int64             `json:"course_id,omitempty"`
	AccountID       int64             `json:"account_id,omitempty"`
	Role            string            `json:"role,omitempty"`
	GroupCategoryID int64             `json:"group_category_id,omitempty"`
	SISGroupID      string            `json:"sis_group_id,omitempty"`
	SISImportID     int64             `json:"sis_import_id,omitempty"`
	StorageQuotaMb  int64             `json:"storage_quota_mb,omitempty"`
	Permissions     *GroupPermissions `json:"permissions,omitempty"`
	Users           []User            `json:"users,omitempty"`
}

// GroupPermissions represents permissions on a group
type GroupPermissions struct {
	CreateDiscussionTopic bool `json:"create_discussion_topic"`
	CreateAnnouncement    bool `json:"create_announcement"`
}

// GroupCategory represents a Canvas group category
type GroupCategory struct {
	ID                 int64       `json:"id"`
	Name               string      `json:"name"`
	Role               string      `json:"role,omitempty"`
	SelfSignup         string      `json:"self_signup,omitempty"`
	AutoLeader         string      `json:"auto_leader,omitempty"`
	ContextType        string      `json:"context_type,omitempty"`
	AccountID          int64       `json:"account_id,omitempty"`
	CourseID           int64       `json:"course_id,omitempty"`
	GroupLimit         int         `json:"group_limit,omitempty"`
	SISGroupCategoryID string      `json:"sis_group_category_id,omitempty"`
	SISImportID        int64       `json:"sis_import_id,omitempty"`
	Progress           interface{} `json:"progress,omitempty"`
}

// GroupMembership represents a user's membership in a group
type GroupMembership struct {
	ID            int64  `json:"id"`
	GroupID       int64  `json:"group_id"`
	UserID        int64  `json:"user_id"`
	WorkflowState string `json:"workflow_state"`
	Moderator     bool   `json:"moderator"`
	JustCreated   bool   `json:"just_created,omitempty"`
	SISImportID   int64  `json:"sis_import_id,omitempty"`
}

// ListGroupsOptions holds options for listing groups
type ListGroupsOptions struct {
	Include []string // users, permissions, tabs
	Page    int
	PerPage int
}

// ListCourse retrieves all groups for a course
func (s *GroupsService) ListCourse(ctx context.Context, courseID int64, opts *ListGroupsOptions) ([]Group, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/groups", courseID)

	if opts != nil {
		query := url.Values{}

		for _, include := range opts.Include {
			query.Add("include[]", include)
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

	var groups []Group
	if err := s.client.GetAllPages(ctx, path, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// ListAccount retrieves all groups for an account
func (s *GroupsService) ListAccount(ctx context.Context, accountID int64, opts *ListGroupsOptions) ([]Group, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/groups", accountID)

	if opts != nil {
		query := url.Values{}

		for _, include := range opts.Include {
			query.Add("include[]", include)
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

	var groups []Group
	if err := s.client.GetAllPages(ctx, path, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// ListUser retrieves all groups for a user
func (s *GroupsService) ListUser(ctx context.Context, userID int64, opts *ListGroupsOptions) ([]Group, error) {
	var path string
	if userID > 0 {
		path = fmt.Sprintf("/api/v1/users/%d/groups", userID)
	} else {
		path = "/api/v1/users/self/groups"
	}

	if opts != nil {
		query := url.Values{}

		for _, include := range opts.Include {
			query.Add("include[]", include)
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

	var groups []Group
	if err := s.client.GetAllPages(ctx, path, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// Get retrieves a single group
func (s *GroupsService) Get(ctx context.Context, groupID int64, include []string) (*Group, error) {
	path := fmt.Sprintf("/api/v1/groups/%d", groupID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var group Group
	if err := s.client.GetJSON(ctx, path, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// CreateGroupParams holds parameters for creating a group
type CreateGroupParams struct {
	Name           string
	Description    string
	IsPublic       bool
	JoinLevel      string // parent_context_auto_join, parent_context_request, invitation_only
	StorageQuotaMb int64
	SISGroupID     string
}

// Create creates a new group
func (s *GroupsService) Create(ctx context.Context, categoryID int64, params *CreateGroupParams) (*Group, error) {
	path := fmt.Sprintf("/api/v1/group_categories/%d/groups", categoryID)

	body := make(map[string]interface{})

	if params.Name != "" {
		body["name"] = params.Name
	}

	if params.Description != "" {
		body["description"] = params.Description
	}

	if params.IsPublic {
		body["is_public"] = params.IsPublic
	}

	if params.JoinLevel != "" {
		body["join_level"] = params.JoinLevel
	}

	if params.StorageQuotaMb > 0 {
		body["storage_quota_mb"] = params.StorageQuotaMb
	}

	if params.SISGroupID != "" {
		body["sis_group_id"] = params.SISGroupID
	}

	var group Group
	if err := s.client.PostJSON(ctx, path, body, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// UpdateGroupParams holds parameters for updating a group
type UpdateGroupParams struct {
	Name           *string
	Description    *string
	IsPublic       *bool
	JoinLevel      *string
	AvatarID       *int64
	StorageQuotaMb *int64
	SISGroupID     *string
}

// Update updates an existing group
func (s *GroupsService) Update(ctx context.Context, groupID int64, params *UpdateGroupParams) (*Group, error) {
	path := fmt.Sprintf("/api/v1/groups/%d", groupID)

	body := make(map[string]interface{})

	if params.Name != nil {
		body["name"] = *params.Name
	}

	if params.Description != nil {
		body["description"] = *params.Description
	}

	if params.IsPublic != nil {
		body["is_public"] = *params.IsPublic
	}

	if params.JoinLevel != nil {
		body["join_level"] = *params.JoinLevel
	}

	if params.AvatarID != nil {
		body["avatar_id"] = *params.AvatarID
	}

	if params.StorageQuotaMb != nil {
		body["storage_quota_mb"] = *params.StorageQuotaMb
	}

	if params.SISGroupID != nil {
		body["sis_group_id"] = *params.SISGroupID
	}

	var group Group
	if err := s.client.PutJSON(ctx, path, body, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// Delete deletes a group
func (s *GroupsService) Delete(ctx context.Context, groupID int64) (*Group, error) {
	path := fmt.Sprintf("/api/v1/groups/%d", groupID)

	var group Group
	if err := s.client.DeleteJSON(ctx, path, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// ListMembers retrieves all members of a group
func (s *GroupsService) ListMembers(ctx context.Context, groupID int64) ([]User, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/users", groupID)

	var users []User
	if err := s.client.GetAllPages(ctx, path, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// AddMember adds a user to a group
func (s *GroupsService) AddMember(ctx context.Context, groupID, userID int64) (*GroupMembership, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/memberships", groupID)

	body := map[string]interface{}{
		"user_id": userID,
	}

	var membership GroupMembership
	if err := s.client.PostJSON(ctx, path, body, &membership); err != nil {
		return nil, err
	}

	return &membership, nil
}

// RemoveMember removes a user from a group
func (s *GroupsService) RemoveMember(ctx context.Context, groupID, membershipID int64) error {
	path := fmt.Sprintf("/api/v1/groups/%d/memberships/%d", groupID, membershipID)

	_, err := s.client.Delete(ctx, path)
	return err
}

// ListCategoriesOptions holds options for listing group categories
type ListCategoriesOptions struct {
	Page    int
	PerPage int
}

// ListCategoriesCourse retrieves group categories for a course
func (s *GroupsService) ListCategoriesCourse(ctx context.Context, courseID int64, opts *ListCategoriesOptions) ([]GroupCategory, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/group_categories", courseID)

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

	var categories []GroupCategory
	if err := s.client.GetAllPages(ctx, path, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// ListCategoriesAccount retrieves group categories for an account
func (s *GroupsService) ListCategoriesAccount(ctx context.Context, accountID int64, opts *ListCategoriesOptions) ([]GroupCategory, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/group_categories", accountID)

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

	var categories []GroupCategory
	if err := s.client.GetAllPages(ctx, path, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategory retrieves a single group category
func (s *GroupsService) GetCategory(ctx context.Context, categoryID int64) (*GroupCategory, error) {
	path := fmt.Sprintf("/api/v1/group_categories/%d", categoryID)

	var category GroupCategory
	if err := s.client.GetJSON(ctx, path, &category); err != nil {
		return nil, err
	}

	return &category, nil
}

// CreateCategoryParams holds parameters for creating a group category
type CreateCategoryParams struct {
	Name               string
	SelfSignup         string // enabled, restricted
	AutoLeader         string // first, random
	GroupLimit         int
	CreateGroupCount   int
	SplitGroupCount    int
	SISGroupCategoryID string
}

// CreateCategoryCourse creates a new group category in a course
func (s *GroupsService) CreateCategoryCourse(ctx context.Context, courseID int64, params *CreateCategoryParams) (*GroupCategory, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/group_categories", courseID)

	body := make(map[string]interface{})

	if params.Name != "" {
		body["name"] = params.Name
	}

	if params.SelfSignup != "" {
		body["self_signup"] = params.SelfSignup
	}

	if params.AutoLeader != "" {
		body["auto_leader"] = params.AutoLeader
	}

	if params.GroupLimit > 0 {
		body["group_limit"] = params.GroupLimit
	}

	if params.CreateGroupCount > 0 {
		body["create_group_count"] = params.CreateGroupCount
	}

	if params.SplitGroupCount > 0 {
		body["split_group_count"] = params.SplitGroupCount
	}

	if params.SISGroupCategoryID != "" {
		body["sis_group_category_id"] = params.SISGroupCategoryID
	}

	var category GroupCategory
	if err := s.client.PostJSON(ctx, path, body, &category); err != nil {
		return nil, err
	}

	return &category, nil
}

// CreateCategoryAccount creates a new group category in an account
func (s *GroupsService) CreateCategoryAccount(ctx context.Context, accountID int64, params *CreateCategoryParams) (*GroupCategory, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/group_categories", accountID)

	body := make(map[string]interface{})

	if params.Name != "" {
		body["name"] = params.Name
	}

	if params.SelfSignup != "" {
		body["self_signup"] = params.SelfSignup
	}

	if params.AutoLeader != "" {
		body["auto_leader"] = params.AutoLeader
	}

	if params.GroupLimit > 0 {
		body["group_limit"] = params.GroupLimit
	}

	if params.SISGroupCategoryID != "" {
		body["sis_group_category_id"] = params.SISGroupCategoryID
	}

	var category GroupCategory
	if err := s.client.PostJSON(ctx, path, body, &category); err != nil {
		return nil, err
	}

	return &category, nil
}

// UpdateCategoryParams holds parameters for updating a group category
type UpdateCategoryParams struct {
	Name               *string
	SelfSignup         *string
	AutoLeader         *string
	GroupLimit         *int
	SISGroupCategoryID *string
}

// UpdateCategory updates an existing group category
func (s *GroupsService) UpdateCategory(ctx context.Context, categoryID int64, params *UpdateCategoryParams) (*GroupCategory, error) {
	path := fmt.Sprintf("/api/v1/group_categories/%d", categoryID)

	body := make(map[string]interface{})

	if params.Name != nil {
		body["name"] = *params.Name
	}

	if params.SelfSignup != nil {
		body["self_signup"] = *params.SelfSignup
	}

	if params.AutoLeader != nil {
		body["auto_leader"] = *params.AutoLeader
	}

	if params.GroupLimit != nil {
		body["group_limit"] = *params.GroupLimit
	}

	if params.SISGroupCategoryID != nil {
		body["sis_group_category_id"] = *params.SISGroupCategoryID
	}

	var category GroupCategory
	if err := s.client.PutJSON(ctx, path, body, &category); err != nil {
		return nil, err
	}

	return &category, nil
}

// DeleteCategory deletes a group category
func (s *GroupsService) DeleteCategory(ctx context.Context, categoryID int64) (*GroupCategory, error) {
	path := fmt.Sprintf("/api/v1/group_categories/%d", categoryID)

	var category GroupCategory
	if err := s.client.DeleteJSON(ctx, path, &category); err != nil {
		return nil, err
	}

	return &category, nil
}

// ListGroupsInCategory retrieves all groups in a category
func (s *GroupsService) ListGroupsInCategory(ctx context.Context, categoryID int64) ([]Group, error) {
	path := fmt.Sprintf("/api/v1/group_categories/%d/groups", categoryID)

	var groups []Group
	if err := s.client.GetAllPages(ctx, path, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// CreateStandalone creates a group at the account level (no category required).
// POST /api/v1/groups
func (s *GroupsService) CreateStandalone(ctx context.Context, params *CreateGroupParams) (*Group, error) {
	body := make(map[string]interface{})

	if params.Name != "" {
		body["name"] = params.Name
	}
	if params.Description != "" {
		body["description"] = params.Description
	}
	if params.IsPublic {
		body["is_public"] = params.IsPublic
	}
	if params.JoinLevel != "" {
		body["join_level"] = params.JoinLevel
	}
	if params.StorageQuotaMb > 0 {
		body["storage_quota_mb"] = params.StorageQuotaMb
	}
	if params.SISGroupID != "" {
		body["sis_group_id"] = params.SISGroupID
	}

	var group Group
	if err := s.client.PostJSON(ctx, "/api/v1/groups", body, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// ListMemberships returns the group membership records for a group.
// GET /api/v1/groups/:group_id/memberships
func (s *GroupsService) ListMemberships(ctx context.Context, groupID int64) ([]GroupMembership, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/memberships", groupID)

	var memberships []GroupMembership
	if err := s.client.GetAllPages(ctx, path, &memberships); err != nil {
		return nil, err
	}

	return memberships, nil
}

// GetMembership returns a single membership record.
// GET /api/v1/groups/:group_id/memberships/:membership_id
func (s *GroupsService) GetMembership(ctx context.Context, groupID, membershipID int64) (*GroupMembership, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/memberships/%d", groupID, membershipID)

	var membership GroupMembership
	if err := s.client.GetJSON(ctx, path, &membership); err != nil {
		return nil, err
	}

	return &membership, nil
}

// UpdateMembershipParams holds parameters for updating a group membership.
type UpdateMembershipParams struct {
	WorkflowState *string // accepted, invited, rejected, deleted
	Moderator     *bool
}

// UpdateMembership updates an existing membership (e.g. promote to moderator).
// PUT /api/v1/groups/:group_id/memberships/:membership_id
func (s *GroupsService) UpdateMembership(ctx context.Context, groupID, membershipID int64, params *UpdateMembershipParams) (*GroupMembership, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/memberships/%d", groupID, membershipID)

	body := make(map[string]interface{})
	if params.WorkflowState != nil {
		body["workflow_state"] = *params.WorkflowState
	}
	if params.Moderator != nil {
		body["moderator"] = *params.Moderator
	}

	var membership GroupMembership
	if err := s.client.PutJSON(ctx, path, body, &membership); err != nil {
		return nil, err
	}

	return &membership, nil
}

// RemoveUserBySelf removes the calling user from a group (spec path: DELETE /api/v1/groups/:group_id/users).
// This is the "leave group" endpoint that takes no membership_id.
func (s *GroupsService) RemoveUserBySelf(ctx context.Context, groupID int64) error {
	path := fmt.Sprintf("/api/v1/groups/%d/users", groupID)

	_, err := s.client.Delete(ctx, path)
	return err
}

// GetUser returns a specific user's membership record within a group.
// GET /api/v1/groups/:group_id/users/:user_id
func (s *GroupsService) GetUser(ctx context.Context, groupID, userID int64) (*User, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/users/%d", groupID, userID)

	var user User
	if err := s.client.GetJSON(ctx, path, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserMembership updates a user's membership state or moderator status.
// PUT /api/v1/groups/:group_id/users/:user_id
func (s *GroupsService) UpdateUserMembership(ctx context.Context, groupID, userID int64, params *UpdateMembershipParams) (*GroupMembership, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/users/%d", groupID, userID)

	body := make(map[string]interface{})
	if params.WorkflowState != nil {
		body["workflow_state"] = *params.WorkflowState
	}
	if params.Moderator != nil {
		body["moderator"] = *params.Moderator
	}

	var membership GroupMembership
	if err := s.client.PutJSON(ctx, path, body, &membership); err != nil {
		return nil, err
	}

	return &membership, nil
}

// RemoveUser removes a specific user from a group.
// DELETE /api/v1/groups/:group_id/users/:user_id
func (s *GroupsService) RemoveUser(ctx context.Context, groupID, userID int64) error {
	path := fmt.Sprintf("/api/v1/groups/%d/users/%d", groupID, userID)

	_, err := s.client.Delete(ctx, path)
	return err
}

// GetActivityStream returns the activity stream for a group.
// GET /api/v1/groups/:group_id/activity_stream
func (s *GroupsService) GetActivityStream(ctx context.Context, groupID int64) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/activity_stream", groupID)

	var items []map[string]interface{}
	if err := s.client.GetAllPages(ctx, path, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetActivityStreamSummary returns the activity stream summary for a group.
// GET /api/v1/groups/:group_id/activity_stream/summary
func (s *GroupsService) GetActivityStreamSummary(ctx context.Context, groupID int64) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/activity_stream/summary", groupID)

	var items []map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetPermissions returns permissions the current user has on the group.
// GET /api/v1/groups/:group_id/permissions
func (s *GroupsService) GetPermissions(ctx context.Context, groupID int64, permissions []string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/permissions", groupID)

	if len(permissions) > 0 {
		query := url.Values{}
		for _, p := range permissions {
			query.Add("permissions[]", p)
		}
		path += "?" + query.Encode()
	}

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Invite invites users to join a group.
// POST /api/v1/groups/:group_id/invite
func (s *GroupsService) Invite(ctx context.Context, groupID int64, invitees []string) ([]GroupMembership, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/invite", groupID)

	body := map[string]interface{}{
		"invitees": invitees,
	}

	var memberships []GroupMembership
	if err := s.client.PostJSON(ctx, path, body, &memberships); err != nil {
		return nil, err
	}

	return memberships, nil
}

// GroupTab represents a navigation tab in a group
type GroupTab struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Type     string `json:"type"`
	Hidden   bool   `json:"hidden,omitempty"`
	Position int    `json:"position,omitempty"`
}

// ListTabs returns the navigation tabs for a group.
// GET /api/v1/groups/:group_id/tabs
func (s *GroupsService) ListTabs(ctx context.Context, groupID int64) ([]GroupTab, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/tabs", groupID)

	var tabs []GroupTab
	if err := s.client.GetAllPages(ctx, path, &tabs); err != nil {
		return nil, err
	}

	return tabs, nil
}

// ListCollaborations returns collaborations for a group.
// GET /api/v1/groups/:group_id/collaborations
func (s *GroupsService) ListCollaborations(ctx context.Context, groupID int64) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/collaborations", groupID)

	var items []map[string]interface{}
	if err := s.client.GetAllPages(ctx, path, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// ListConferences returns conferences for a group.
// GET /api/v1/groups/:group_id/conferences
func (s *GroupsService) ListConferences(ctx context.Context, groupID int64) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/conferences", groupID)

	var items []map[string]interface{}
	if err := s.client.GetAllPages(ctx, path, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// GroupExternalFeed represents an external feed subscription for a group.
type GroupExternalFeed struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"display_name,omitempty"`
	URL         string `json:"url"`
	Header      string `json:"header_match,omitempty"`
	Created     string `json:"created_at,omitempty"`
	Verbosity   string `json:"verbosity,omitempty"`
	ContextID   int64  `json:"context_id,omitempty"`
	ContextType string `json:"context_type,omitempty"`
}

// ListExternalFeeds returns the external feeds for a group.
// GET /api/v1/groups/:group_id/external_feeds
func (s *GroupsService) ListExternalFeeds(ctx context.Context, groupID int64) ([]GroupExternalFeed, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/external_feeds", groupID)

	var feeds []GroupExternalFeed
	if err := s.client.GetAllPages(ctx, path, &feeds); err != nil {
		return nil, err
	}

	return feeds, nil
}

// CreateExternalFeedParams holds parameters for creating an external feed.
type CreateExternalFeedParams struct {
	URL       string
	Verbosity string // full, truncate, link_only
	Header    string
}

// CreateExternalFeed adds an external feed to a group.
// POST /api/v1/groups/:group_id/external_feeds
func (s *GroupsService) CreateExternalFeed(ctx context.Context, groupID int64, params *CreateExternalFeedParams) (*GroupExternalFeed, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/external_feeds", groupID)

	body := make(map[string]interface{})
	if params.URL != "" {
		body["url"] = params.URL
	}
	if params.Verbosity != "" {
		body["verbosity"] = params.Verbosity
	}
	if params.Header != "" {
		body["header_match"] = params.Header
	}

	var feed GroupExternalFeed
	if err := s.client.PostJSON(ctx, path, body, &feed); err != nil {
		return nil, err
	}

	return &feed, nil
}

// DeleteExternalFeed removes an external feed from a group.
// DELETE /api/v1/groups/:group_id/external_feeds/:external_feed_id
func (s *GroupsService) DeleteExternalFeed(ctx context.Context, groupID, feedID int64) (*GroupExternalFeed, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/external_feeds/%d", groupID, feedID)

	var feed GroupExternalFeed
	if err := s.client.DeleteJSON(ctx, path, &feed); err != nil {
		return nil, err
	}

	return &feed, nil
}

// ListExternalTools returns LTI external tools for a group.
// GET /api/v1/groups/:group_id/external_tools
func (s *GroupsService) ListExternalTools(ctx context.Context, groupID int64) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/external_tools", groupID)

	var tools []map[string]interface{}
	if err := s.client.GetAllPages(ctx, path, &tools); err != nil {
		return nil, err
	}

	return tools, nil
}

// ContentExport represents a content export object.
type ContentExport struct {
	ID            int64       `json:"id"`
	CreatedAt     string      `json:"created_at,omitempty"`
	ExportType    string      `json:"export_type"`
	Attachment    interface{} `json:"attachment,omitempty"`
	ProgressURL   string      `json:"progress_url,omitempty"`
	UserID        int64       `json:"user_id,omitempty"`
	WorkflowState string      `json:"workflow_state,omitempty"`
}

// ListContentExports lists content exports for a group.
// GET /api/v1/groups/:group_id/content_exports
func (s *GroupsService) ListContentExports(ctx context.Context, groupID int64) ([]ContentExport, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/content_exports", groupID)

	var exports []ContentExport
	if err := s.client.GetAllPages(ctx, path, &exports); err != nil {
		return nil, err
	}

	return exports, nil
}

// CreateContentExport starts a content export for a group.
// POST /api/v1/groups/:group_id/content_exports
func (s *GroupsService) CreateContentExport(ctx context.Context, groupID int64, exportType string) (*ContentExport, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/content_exports", groupID)

	body := map[string]interface{}{
		"export_type": exportType,
	}

	var export ContentExport
	if err := s.client.PostJSON(ctx, path, body, &export); err != nil {
		return nil, err
	}

	return &export, nil
}

// GetContentExport retrieves a single content export for a group.
// GET /api/v1/groups/:group_id/content_exports/:id
func (s *GroupsService) GetContentExport(ctx context.Context, groupID, exportID int64) (*ContentExport, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/content_exports/%d", groupID, exportID)

	var export ContentExport
	if err := s.client.GetJSON(ctx, path, &export); err != nil {
		return nil, err
	}

	return &export, nil
}

// ContentLicenseItem represents an available content license.
type ContentLicenseItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// ListContentLicenses returns the list of content licenses for a group.
// GET /api/v1/groups/:group_id/content_licenses
func (s *GroupsService) ListContentLicenses(ctx context.Context, groupID int64) ([]ContentLicenseItem, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/content_licenses", groupID)

	var licenses []ContentLicenseItem
	if err := s.client.GetAllPages(ctx, path, &licenses); err != nil {
		return nil, err
	}

	return licenses, nil
}

// MediaAttachment represents a media object attached to a group.
type MediaAttachment struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"display_name,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	URL         string `json:"url,omitempty"`
	MediaID     string `json:"media_id,omitempty"`
}

// ListMediaAttachments returns media attachments for a group.
// GET /api/v1/groups/:group_id/media_attachments
func (s *GroupsService) ListMediaAttachments(ctx context.Context, groupID int64) ([]MediaAttachment, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/media_attachments", groupID)

	var items []MediaAttachment
	if err := s.client.GetAllPages(ctx, path, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// MediaObject represents a media object in Canvas.
type MediaObject struct {
	MediaID        string `json:"media_id"`
	MediaType      string `json:"media_type,omitempty"`
	Title          string `json:"title,omitempty"`
	CanAddCaptions bool   `json:"can_add_captions,omitempty"`
	UserEntered    bool   `json:"user_entered_title,omitempty"`
}

// ListMediaObjects returns media objects for a group.
// GET /api/v1/groups/:group_id/media_objects
func (s *GroupsService) ListMediaObjects(ctx context.Context, groupID int64) ([]MediaObject, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/media_objects", groupID)

	var items []MediaObject
	if err := s.client.GetAllPages(ctx, path, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// ListPotentialCollaborators returns users who can be added as collaborators.
// GET /api/v1/groups/:group_id/potential_collaborators
func (s *GroupsService) ListPotentialCollaborators(ctx context.Context, groupID int64) ([]User, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/potential_collaborators", groupID)

	var users []User
	if err := s.client.GetAllPages(ctx, path, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// PreviewHTML returns a preview of HTML content within the group context.
// POST /api/v1/groups/:group_id/preview_html
func (s *GroupsService) PreviewHTML(ctx context.Context, groupID int64, html string) (string, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/preview_html", groupID)

	body := map[string]interface{}{
		"html": html,
	}

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return "", err
	}

	if preview, ok := result["html"].(string); ok {
		return preview, nil
	}

	return "", nil
}

// UsageRightsResult holds the result of a usage rights operation.
type UsageRightsResult struct {
	UseJustification string  `json:"use_justification"`
	License          string  `json:"license,omitempty"`
	LicenseName      string  `json:"license_name,omitempty"`
	Message          string  `json:"message,omitempty"`
	FileIDs          []int64 `json:"file_ids,omitempty"`
}

// DeleteUsageRights removes usage rights from files in a group context.
// DELETE /api/v1/groups/:group_id/usage_rights
func (s *GroupsService) DeleteUsageRights(ctx context.Context, groupID int64, fileIDs []int64) (*UsageRightsResult, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/usage_rights", groupID)

	query := url.Values{}
	for _, id := range fileIDs {
		query.Add("file_ids[]", strconv.FormatInt(id, 10))
	}
	if len(fileIDs) > 0 {
		path += "?" + query.Encode()
	}

	var result UsageRightsResult
	if err := s.client.DeleteJSON(ctx, path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SetUsageRightsParams holds parameters for setting usage rights.
type SetUsageRightsGroupParams struct {
	FileIDs          []int64
	UseJustification string
	License          string
	HolderName       string
}

// SetUsageRights sets usage rights for files in a group context.
// PUT /api/v1/groups/:group_id/usage_rights
func (s *GroupsService) SetUsageRights(ctx context.Context, groupID int64, params *SetUsageRightsGroupParams) (*UsageRightsResult, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/usage_rights", groupID)

	body := make(map[string]interface{})
	if len(params.FileIDs) > 0 {
		body["file_ids"] = params.FileIDs
	}
	if params.UseJustification != "" {
		body["usage_rights[use_justification]"] = params.UseJustification
	}
	if params.License != "" {
		body["usage_rights[license]"] = params.License
	}
	if params.HolderName != "" {
		body["usage_rights[license_name]"] = params.HolderName
	}

	var result UsageRightsResult
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAssignmentOverride returns the assignment override for a group.
// GET /api/v1/groups/:group_id/assignments/:assignment_id/override
func (s *GroupsService) GetAssignmentOverride(ctx context.Context, groupID, assignmentID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/assignments/%d/override", groupID, assignmentID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// AssignUnassignedMembers assigns unassigned members to groups in a category.
// POST /api/v1/group_categories/:group_category_id/assign_unassigned_members
func (s *GroupsService) AssignUnassignedMembers(ctx context.Context, categoryID int64, sync bool) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/group_categories/%d/assign_unassigned_members", categoryID)

	body := map[string]interface{}{
		"sync": sync,
	}

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// ListUsersInCategory returns the users in a group category.
// GET /api/v1/group_categories/:group_category_id/users
func (s *GroupsService) ListUsersInCategory(ctx context.Context, categoryID int64) ([]User, error) {
	path := fmt.Sprintf("/api/v1/group_categories/%d/users", categoryID)

	var users []User
	if err := s.client.GetAllPages(ctx, path, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// ExportCategory exports a group category.
// GET /api/v1/group_categories/:group_category_id/export
func (s *GroupsService) ExportCategory(ctx context.Context, categoryID int64) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/group_categories/%d/export", categoryID)

	var result map[string]interface{}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// ImportCategory imports users into groups in a category.
// POST /api/v1/group_categories/:group_category_id/import
func (s *GroupsService) ImportCategory(ctx context.Context, categoryID int64, attachment interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/group_categories/%d/import", categoryID)

	body := map[string]interface{}{
		"attachment": attachment,
	}

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateGroupInCategory creates a group within a specific group category.
// POST /api/v1/group_categories/:group_category_id/groups
// This is an alias for the existing Create method, exposing the spec path explicitly.
func (s *GroupsService) CreateGroupInCategory(ctx context.Context, categoryID int64, params *CreateGroupParams) (*Group, error) {
	return s.Create(ctx, categoryID, params)
}
