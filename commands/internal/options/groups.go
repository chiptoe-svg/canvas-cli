package options

import "fmt"

// GroupsListOptions contains options for listing groups
type GroupsListOptions struct {
	CourseID           int64
	AccountID          int64
	UserID             int64
	IncludeUsers       bool
	IncludePermissions bool
}

// Validate validates the options
func (o *GroupsListOptions) Validate() error {
	// No required fields - defaults to user's groups
	return nil
}

// GroupsGetOptions contains options for getting a group
type GroupsGetOptions struct {
	GroupID            int64
	IncludeUsers       bool
	IncludePermissions bool
}

// Validate validates the options
func (o *GroupsGetOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// GroupsCreateOptions contains options for creating a group
type GroupsCreateOptions struct {
	CategoryID     int64
	Name           string
	Description    string
	IsPublic       bool
	JoinLevel      string
	StorageQuotaMb int64
	SISGroupID     string
}

// Validate validates the options
func (o *GroupsCreateOptions) Validate() error {
	if o.CategoryID <= 0 {
		return fmt.Errorf("category-id is required and must be greater than 0")
	}
	if o.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// GroupsUpdateOptions contains options for updating a group
type GroupsUpdateOptions struct {
	GroupID        int64
	Name           string
	Description    string
	IsPublic       bool
	JoinLevel      string
	StorageQuotaMb int64
	SISGroupID     string
	// Track which fields were actually set
	NameSet           bool
	DescriptionSet    bool
	IsPublicSet       bool
	JoinLevelSet      bool
	StorageQuotaMbSet bool
	SISGroupIDSet     bool
}

// Validate validates the options
func (o *GroupsUpdateOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	// At least one field must be set for update
	if !o.NameSet && !o.DescriptionSet && !o.IsPublicSet &&
		!o.JoinLevelSet && !o.StorageQuotaMbSet && !o.SISGroupIDSet {
		return fmt.Errorf("at least one field must be specified for update")
	}
	return nil
}

// GroupsDeleteOptions contains options for deleting a group
type GroupsDeleteOptions struct {
	GroupID int64
	Force   bool
}

// Validate validates the options
func (o *GroupsDeleteOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// GroupsMembersListOptions contains options for listing group members
type GroupsMembersListOptions struct {
	GroupID int64
}

// Validate validates the options
func (o *GroupsMembersListOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// GroupsMembersAddOptions contains options for adding a group member
type GroupsMembersAddOptions struct {
	GroupID int64
	UserID  int64
}

// Validate validates the options
func (o *GroupsMembersAddOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// GroupsMembersRemoveOptions contains options for removing a group member
type GroupsMembersRemoveOptions struct {
	GroupID      int64
	MembershipID int64
}

// Validate validates the options
func (o *GroupsMembersRemoveOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.MembershipID <= 0 {
		return fmt.Errorf("membership-id is required and must be greater than 0")
	}
	return nil
}

// GroupsCategoriesListOptions contains options for listing group categories
type GroupsCategoriesListOptions struct {
	CourseID  int64
	AccountID int64
}

// Validate validates the options
func (o *GroupsCategoriesListOptions) Validate() error {
	if o.CourseID == 0 && o.AccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}
	return nil
}

// GroupsCategoriesGetOptions contains options for getting a group category
type GroupsCategoriesGetOptions struct {
	CategoryID int64
}

// Validate validates the options
func (o *GroupsCategoriesGetOptions) Validate() error {
	if o.CategoryID <= 0 {
		return fmt.Errorf("category-id is required and must be greater than 0")
	}
	return nil
}

// GroupsCategoriesCreateOptions contains options for creating a group category
type GroupsCategoriesCreateOptions struct {
	CourseID         int64
	AccountID        int64
	Name             string
	SelfSignup       string
	AutoLeader       string
	GroupLimit       int
	CreateGroupCount int
	SplitGroupCount  int
	SISCategoryID    string
}

// Validate validates the options
func (o *GroupsCategoriesCreateOptions) Validate() error {
	if o.CourseID == 0 && o.AccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}
	if o.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// GroupsCategoriesUpdateOptions contains options for updating a group category
type GroupsCategoriesUpdateOptions struct {
	CategoryID    int64
	Name          string
	SelfSignup    string
	AutoLeader    string
	GroupLimit    int
	SISCategoryID string
	// Track which fields were actually set
	NameSet          bool
	SelfSignupSet    bool
	AutoLeaderSet    bool
	GroupLimitSet    bool
	SISCategoryIDSet bool
}

// Validate validates the options
func (o *GroupsCategoriesUpdateOptions) Validate() error {
	if o.CategoryID <= 0 {
		return fmt.Errorf("category-id is required and must be greater than 0")
	}
	// At least one field must be set for update
	if !o.NameSet && !o.SelfSignupSet && !o.AutoLeaderSet &&
		!o.GroupLimitSet && !o.SISCategoryIDSet {
		return fmt.Errorf("at least one field must be specified for update")
	}
	return nil
}

// GroupsCategoriesDeleteOptions contains options for deleting a group category
type GroupsCategoriesDeleteOptions struct {
	CategoryID int64
	Force      bool
}

// Validate validates the options
func (o *GroupsCategoriesDeleteOptions) Validate() error {
	if o.CategoryID <= 0 {
		return fmt.Errorf("category-id is required and must be greater than 0")
	}
	return nil
}

// GroupsCategoriesGroupsOptions contains options for listing groups in a category
type GroupsCategoriesGroupsOptions struct {
	CategoryID int64
}

// Validate validates the options
func (o *GroupsCategoriesGroupsOptions) Validate() error {
	if o.CategoryID <= 0 {
		return fmt.Errorf("category-id is required and must be greater than 0")
	}
	return nil
}

// GroupsCreateStandaloneOptions contains options for creating a standalone group
type GroupsCreateStandaloneOptions struct {
	Name           string
	Description    string
	IsPublic       bool
	JoinLevel      string
	StorageQuotaMb int64
	SISGroupID     string
}

// Validate validates the options
func (o *GroupsCreateStandaloneOptions) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// GroupsMembershipsListOptions contains options for listing memberships
type GroupsMembershipsListOptions struct {
	GroupID int64
}

// Validate validates the options
func (o *GroupsMembershipsListOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// GroupsMembershipsGetOptions contains options for getting a single membership
type GroupsMembershipsGetOptions struct {
	GroupID      int64
	MembershipID int64
}

// Validate validates the options
func (o *GroupsMembershipsGetOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.MembershipID <= 0 {
		return fmt.Errorf("membership-id is required and must be greater than 0")
	}
	return nil
}

// GroupsMembershipsUpdateOptions contains options for updating a membership
type GroupsMembershipsUpdateOptions struct {
	GroupID          int64
	MembershipID     int64
	WorkflowState    string
	Moderator        bool
	WorkflowStateSet bool
	ModeratorSet     bool
}

// Validate validates the options
func (o *GroupsMembershipsUpdateOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.MembershipID <= 0 {
		return fmt.Errorf("membership-id is required and must be greater than 0")
	}
	return nil
}

// GroupsUsersGetOptions contains options for getting a group user
type GroupsUsersGetOptions struct {
	GroupID int64
	UserID  int64
}

// Validate validates the options
func (o *GroupsUsersGetOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// GroupsUsersUpdateOptions contains options for updating a group user membership
type GroupsUsersUpdateOptions struct {
	GroupID          int64
	UserID           int64
	WorkflowState    string
	Moderator        bool
	WorkflowStateSet bool
	ModeratorSet     bool
}

// Validate validates the options
func (o *GroupsUsersUpdateOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// GroupsUsersRemoveOptions contains options for removing a user from a group
type GroupsUsersRemoveOptions struct {
	GroupID int64
	UserID  int64
}

// Validate validates the options
func (o *GroupsUsersRemoveOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// GroupsActivityStreamOptions contains options for getting a group activity stream
type GroupsActivityStreamOptions struct {
	GroupID int64
}

// Validate validates the options
func (o *GroupsActivityStreamOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// GroupsPermissionsOptions contains options for getting group permissions
type GroupsPermissionsOptions struct {
	GroupID     int64
	Permissions []string
}

// Validate validates the options
func (o *GroupsPermissionsOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// GroupsInviteOptions contains options for inviting users to a group
type GroupsInviteOptions struct {
	GroupID  int64
	Invitees []string
}

// Validate validates the options
func (o *GroupsInviteOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if len(o.Invitees) == 0 {
		return fmt.Errorf("at least one invitee email is required")
	}
	return nil
}

// GroupsTabsListOptions contains options for listing group tabs
type GroupsTabsListOptions struct {
	GroupID int64
}

// Validate validates the options
func (o *GroupsTabsListOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// GroupsCollaborationsListOptions contains options for listing group collaborations
type GroupsCollaborationsListOptions struct {
	GroupID int64
}

// Validate validates the options
func (o *GroupsCollaborationsListOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// GroupsConferencesListOptions contains options for listing group conferences
type GroupsConferencesListOptions struct {
	GroupID int64
}

// Validate validates the options
func (o *GroupsConferencesListOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// GroupsExternalFeedsListOptions contains options for listing group external feeds
type GroupsExternalFeedsListOptions struct {
	GroupID int64
}

// Validate validates the options
func (o *GroupsExternalFeedsListOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// GroupsExternalFeedsCreateOptions contains options for creating a group external feed
type GroupsExternalFeedsCreateOptions struct {
	GroupID   int64
	URL       string
	Verbosity string
	Header    string
}

// Validate validates the options
func (o *GroupsExternalFeedsCreateOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.URL == "" {
		return fmt.Errorf("url is required")
	}
	return nil
}

// GroupsExternalFeedsDeleteOptions contains options for deleting a group external feed
type GroupsExternalFeedsDeleteOptions struct {
	GroupID int64
	FeedID  int64
}

// Validate validates the options
func (o *GroupsExternalFeedsDeleteOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.FeedID <= 0 {
		return fmt.Errorf("feed-id is required and must be greater than 0")
	}
	return nil
}

// GroupsContentExportsListOptions contains options for listing group content exports
type GroupsContentExportsListOptions struct {
	GroupID int64
}

// Validate validates the options
func (o *GroupsContentExportsListOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// GroupsContentExportsCreateOptions contains options for creating a content export
type GroupsContentExportsCreateOptions struct {
	GroupID    int64
	ExportType string
}

// Validate validates the options
func (o *GroupsContentExportsCreateOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.ExportType == "" {
		return fmt.Errorf("export-type is required")
	}
	return nil
}

// GroupsContentExportsGetOptions contains options for getting a content export
type GroupsContentExportsGetOptions struct {
	GroupID  int64
	ExportID int64
}

// Validate validates the options
func (o *GroupsContentExportsGetOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.ExportID <= 0 {
		return fmt.Errorf("export-id is required and must be greater than 0")
	}
	return nil
}

// GroupsAssignmentOverrideOptions contains options for getting a group assignment override
type GroupsAssignmentOverrideOptions struct {
	GroupID      int64
	AssignmentID int64
}

// Validate validates the options
func (o *GroupsAssignmentOverrideOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	return nil
}

// GroupsCategoriesAssignMembersOptions contains options for assigning unassigned members
type GroupsCategoriesAssignMembersOptions struct {
	CategoryID int64
	Sync       bool
}

// Validate validates the options
func (o *GroupsCategoriesAssignMembersOptions) Validate() error {
	if o.CategoryID <= 0 {
		return fmt.Errorf("category-id is required and must be greater than 0")
	}
	return nil
}

// GroupsCategoriesUsersListOptions contains options for listing users in a category
type GroupsCategoriesUsersListOptions struct {
	CategoryID int64
}

// Validate validates the options
func (o *GroupsCategoriesUsersListOptions) Validate() error {
	if o.CategoryID <= 0 {
		return fmt.Errorf("category-id is required and must be greater than 0")
	}
	return nil
}

// GroupsCategoriesExportOptions contains options for exporting a category
type GroupsCategoriesExportOptions struct {
	CategoryID int64
}

// Validate validates the options
func (o *GroupsCategoriesExportOptions) Validate() error {
	if o.CategoryID <= 0 {
		return fmt.Errorf("category-id is required and must be greater than 0")
	}
	return nil
}
