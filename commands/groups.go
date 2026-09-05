package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/commands/internal/options"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

// groupsCmd represents the groups command group
var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Manage Canvas groups",
	Long: `Manage Canvas groups and group categories.

Groups allow students and instructors to collaborate on projects and activities.
Groups can be organized into categories with different self-signup options.

Examples:
  canvas groups list --course-id 123
  canvas groups get 456
  canvas groups categories list --course-id 123`,
}

// newGroupsListCmd creates the groups list command
func newGroupsListCmd() *cobra.Command {
	opts := &options.GroupsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List groups",
		Long: `List groups for a course, account, or user.

Examples:
  canvas groups list --course-id 123
  canvas groups list --account-id 1
  canvas groups list --user-id 456
  canvas groups list  # Lists current user's groups`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (0 for self)")
	cmd.Flags().BoolVar(&opts.IncludeUsers, "include-users", false, "Include group users")
	cmd.Flags().BoolVar(&opts.IncludePermissions, "include-permissions", false, "Include permissions")

	return cmd
}

// newGroupsGetCmd creates the groups get command
func newGroupsGetCmd() *cobra.Command {
	opts := &options.GroupsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <group-id>",
		Short: "Get group details",
		Long: `Get details of a specific group.

Examples:
  canvas groups get 456
  canvas groups get 456 --include-users`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.IncludeUsers, "include-users", false, "Include group users")
	cmd.Flags().BoolVar(&opts.IncludePermissions, "include-permissions", false, "Include permissions")

	return cmd
}

// newGroupsCreateCmd creates the groups create command
func newGroupsCreateCmd() *cobra.Command {
	opts := &options.GroupsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new group",
		Long: `Create a new group in a group category.

Examples:
  canvas groups create --category-id 123 --name "Study Group"
  canvas groups create --category-id 123 --name "Project Team" --description "Our project team"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CategoryID, "category-id", 0, "Group category ID (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Group name (required)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Group description")
	cmd.Flags().BoolVar(&opts.IsPublic, "public", false, "Whether the group is public")
	cmd.Flags().StringVar(&opts.JoinLevel, "join-level", "", "Join level (parent_context_auto_join, parent_context_request, invitation_only)")
	cmd.Flags().Int64Var(&opts.StorageQuotaMb, "storage-quota-mb", 0, "Storage quota in MB")
	cmd.Flags().StringVar(&opts.SISGroupID, "sis-group-id", "", "SIS group ID")
	mustMarkRequired(cmd, "category-id", "name")

	return cmd
}

// newGroupsUpdateCmd creates the groups update command
func newGroupsUpdateCmd() *cobra.Command {
	opts := &options.GroupsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <group-id>",
		Short: "Update a group",
		Long: `Update an existing group.

Examples:
  canvas groups update 456 --name "New Name"
  canvas groups update 456 --description "Updated description"`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			// Track which fields were set
			opts.NameSet = cmd.Flags().Changed("name")
			opts.DescriptionSet = cmd.Flags().Changed("description")
			opts.IsPublicSet = cmd.Flags().Changed("public")
			opts.JoinLevelSet = cmd.Flags().Changed("join-level")
			opts.StorageQuotaMbSet = cmd.Flags().Changed("storage-quota-mb")
			opts.SISGroupIDSet = cmd.Flags().Changed("sis-group-id")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Group name")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Group description")
	cmd.Flags().BoolVar(&opts.IsPublic, "public", false, "Whether the group is public")
	cmd.Flags().StringVar(&opts.JoinLevel, "join-level", "", "Join level")
	cmd.Flags().Int64Var(&opts.StorageQuotaMb, "storage-quota-mb", 0, "Storage quota in MB")
	cmd.Flags().StringVar(&opts.SISGroupID, "sis-group-id", "", "SIS group ID")

	return cmd
}

// newGroupsDeleteCmd creates the groups delete command
func newGroupsDeleteCmd() *cobra.Command {
	opts := &options.GroupsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <group-id>",
		Short: "Delete a group",
		Long: `Delete an existing group.

Examples:
  canvas groups delete 456
  canvas groups delete 456 --force`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

// groupsMembersCmd represents the members subcommand
var groupsMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage group members",
	Long:  "Commands for managing group memberships.",
}

// newGroupsMembersListCmd creates the groups members list command
func newGroupsMembersListCmd() *cobra.Command {
	opts := &options.GroupsMembersListOptions{}

	cmd := &cobra.Command{
		Use:   "list <group-id>",
		Short: "List group members",
		Long: `List all members of a group.

Examples:
  canvas groups members list 456`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsMembersList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

// newGroupsMembersAddCmd creates the groups members add command
func newGroupsMembersAddCmd() *cobra.Command {
	opts := &options.GroupsMembersAddOptions{}

	cmd := &cobra.Command{
		Use:   "add <group-id>",
		Short: "Add a member to a group",
		Long: `Add a user to a group.

Examples:
  canvas groups members add 456 --user-id 789`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsMembersAdd(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID to add (required)")
	mustMarkRequired(cmd, "user-id")

	return cmd
}

// newGroupsMembersRemoveCmd creates the groups members remove command
func newGroupsMembersRemoveCmd() *cobra.Command {
	opts := &options.GroupsMembersRemoveOptions{}

	cmd := &cobra.Command{
		Use:   "remove <group-id>",
		Short: "Remove a member from a group",
		Long: `Remove a user from a group by membership ID.

Examples:
  canvas groups members remove 456 --membership-id 123`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsMembersRemove(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.MembershipID, "membership-id", 0, "Membership ID to remove (required)")
	mustMarkRequired(cmd, "membership-id")

	return cmd
}

// groupsCategoriesCmd represents the categories subcommand
var groupsCategoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "Manage group categories",
	Long:  "Commands for managing group categories.",
}

// newGroupsCategoriesListCmd creates the groups categories list command
func newGroupsCategoriesListCmd() *cobra.Command {
	opts := &options.GroupsCategoriesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List group categories",
		Long: `List group categories for a course or account.

Examples:
  canvas groups categories list --course-id 123
  canvas groups categories list --account-id 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsCategoriesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")

	return cmd
}

// newGroupsCategoriesGetCmd creates the groups categories get command
func newGroupsCategoriesGetCmd() *cobra.Command {
	opts := &options.GroupsCategoriesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <category-id>",
		Short: "Get group category details",
		Long: `Get details of a specific group category.

Examples:
  canvas groups categories get 456`,
		Args: ExactArgsWithUsage(1, "category-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			categoryID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid category ID: %w", err)
			}
			opts.CategoryID = categoryID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsCategoriesGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

// newGroupsCategoriesCreateCmd creates the groups categories create command
func newGroupsCategoriesCreateCmd() *cobra.Command {
	opts := &options.GroupsCategoriesCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a group category",
		Long: `Create a new group category in a course or account.

Examples:
  canvas groups categories create --course-id 123 --name "Project Teams"
  canvas groups categories create --account-id 1 --name "Clubs" --self-signup enabled`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsCategoriesCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Category name (required)")
	cmd.Flags().StringVar(&opts.SelfSignup, "self-signup", "", "Self signup (enabled, restricted)")
	cmd.Flags().StringVar(&opts.AutoLeader, "auto-leader", "", "Auto leader (first, random)")
	cmd.Flags().IntVar(&opts.GroupLimit, "group-limit", 0, "Group member limit")
	cmd.Flags().IntVar(&opts.CreateGroupCount, "create-group-count", 0, "Number of groups to create")
	cmd.Flags().IntVar(&opts.SplitGroupCount, "split-group-count", 0, "Number of groups to split students into")
	cmd.Flags().StringVar(&opts.SISCategoryID, "sis-category-id", "", "SIS category ID")
	mustMarkRequired(cmd, "name")

	return cmd
}

// newGroupsCategoriesUpdateCmd creates the groups categories update command
func newGroupsCategoriesUpdateCmd() *cobra.Command {
	opts := &options.GroupsCategoriesUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <category-id>",
		Short: "Update a group category",
		Long: `Update an existing group category.

Examples:
  canvas groups categories update 456 --name "New Name"
  canvas groups categories update 456 --self-signup restricted`,
		Args: ExactArgsWithUsage(1, "category-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			categoryID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid category ID: %w", err)
			}
			opts.CategoryID = categoryID

			// Track which fields were set
			opts.NameSet = cmd.Flags().Changed("name")
			opts.SelfSignupSet = cmd.Flags().Changed("self-signup")
			opts.AutoLeaderSet = cmd.Flags().Changed("auto-leader")
			opts.GroupLimitSet = cmd.Flags().Changed("group-limit")
			opts.SISCategoryIDSet = cmd.Flags().Changed("sis-category-id")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsCategoriesUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Category name")
	cmd.Flags().StringVar(&opts.SelfSignup, "self-signup", "", "Self signup (enabled, restricted)")
	cmd.Flags().StringVar(&opts.AutoLeader, "auto-leader", "", "Auto leader (first, random)")
	cmd.Flags().IntVar(&opts.GroupLimit, "group-limit", 0, "Group member limit")
	cmd.Flags().StringVar(&opts.SISCategoryID, "sis-category-id", "", "SIS category ID")

	return cmd
}

// newGroupsCategoriesDeleteCmd creates the groups categories delete command
func newGroupsCategoriesDeleteCmd() *cobra.Command {
	opts := &options.GroupsCategoriesDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <category-id>",
		Short: "Delete a group category",
		Long: `Delete an existing group category.

Examples:
  canvas groups categories delete 456
  canvas groups categories delete 456 --force`,
		Args: ExactArgsWithUsage(1, "category-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			categoryID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid category ID: %w", err)
			}
			opts.CategoryID = categoryID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsCategoriesDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

// newGroupsCategoriesGroupsCmd creates the groups categories groups command
func newGroupsCategoriesGroupsCmd() *cobra.Command {
	opts := &options.GroupsCategoriesGroupsOptions{}

	cmd := &cobra.Command{
		Use:   "groups <category-id>",
		Short: "List groups in a category",
		Long: `List all groups within a specific group category.

Examples:
  canvas groups categories groups 456`,
		Args: ExactArgsWithUsage(1, "category-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			categoryID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid category ID: %w", err)
			}
			opts.CategoryID = categoryID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsCategoriesGroups(cmd.Context(), client, opts)
		},
	}

	return cmd
}

// groupsMembershipsCmd represents the memberships subcommand
var groupsMembershipsCmd = &cobra.Command{
	Use:   "memberships",
	Short: "Manage group memberships",
	Long:  "Commands for viewing and updating group membership records.",
}

// groupsUsersCmd represents the users subcommand under groups
var groupsUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage group users",
	Long:  "Commands for managing users within a group.",
}

func init() {
	rootCmd.AddCommand(groupsCmd)
	groupsCmd.AddCommand(newGroupsListCmd())
	groupsCmd.AddCommand(newGroupsGetCmd())
	groupsCmd.AddCommand(newGroupsCreateCmd())
	groupsCmd.AddCommand(newGroupsCreateStandaloneCmd())
	groupsCmd.AddCommand(newGroupsUpdateCmd())
	groupsCmd.AddCommand(newGroupsDeleteCmd())
	groupsCmd.AddCommand(newGroupsActivityStreamCmd())
	groupsCmd.AddCommand(newGroupsPermissionsCmd())
	groupsCmd.AddCommand(newGroupsInviteCmd())
	groupsCmd.AddCommand(newGroupsTabsCmd())
	groupsCmd.AddCommand(newGroupsCollaborationsCmd())
	groupsCmd.AddCommand(newGroupsConferencesCmd())
	groupsCmd.AddCommand(newGroupsAssignmentOverrideCmd())
	groupsCmd.AddCommand(groupsMembersCmd)
	groupsCmd.AddCommand(groupsMembershipsCmd)
	groupsCmd.AddCommand(groupsUsersCmd)
	groupsCmd.AddCommand(groupsCategoriesCmd)

	// Members subcommands (existing, based on /users endpoint)
	groupsMembersCmd.AddCommand(newGroupsMembersListCmd())
	groupsMembersCmd.AddCommand(newGroupsMembersAddCmd())
	groupsMembersCmd.AddCommand(newGroupsMembersRemoveCmd())

	// Memberships subcommands (based on /memberships endpoint)
	groupsMembershipsCmd.AddCommand(newGroupsMembershipsListCmd())
	groupsMembershipsCmd.AddCommand(newGroupsMembershipsGetCmd())
	groupsMembershipsCmd.AddCommand(newGroupsMembershipsUpdateCmd())

	// Users subcommands (by user_id)
	groupsUsersCmd.AddCommand(newGroupsUsersGetCmd())
	groupsUsersCmd.AddCommand(newGroupsUsersUpdateCmd())
	groupsUsersCmd.AddCommand(newGroupsUsersRemoveCmd())
	groupsUsersCmd.AddCommand(newGroupsUsersRemoveSelfCmd())

	// External feeds subcommands
	groupsExternalFeedsCmd := &cobra.Command{
		Use:   "external-feeds",
		Short: "Manage group external feeds",
		Long:  "Commands for managing external feed subscriptions on a group.",
	}
	groupsCmd.AddCommand(groupsExternalFeedsCmd)
	groupsExternalFeedsCmd.AddCommand(newGroupsExternalFeedsListCmd())
	groupsExternalFeedsCmd.AddCommand(newGroupsExternalFeedsCreateCmd())
	groupsExternalFeedsCmd.AddCommand(newGroupsExternalFeedsDeleteCmd())

	// Content exports subcommands
	groupsContentExportsCmd := &cobra.Command{
		Use:   "content-exports",
		Short: "Manage group content exports",
		Long:  "Commands for managing content exports for a group.",
	}
	groupsCmd.AddCommand(groupsContentExportsCmd)
	groupsContentExportsCmd.AddCommand(newGroupsContentExportsListCmd())
	groupsContentExportsCmd.AddCommand(newGroupsContentExportsCreateCmd())
	groupsContentExportsCmd.AddCommand(newGroupsContentExportsGetCmd())

	// Categories subcommands
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesListCmd())
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesGetCmd())
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesCreateCmd())
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesUpdateCmd())
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesDeleteCmd())
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesGroupsCmd())
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesAssignMembersCmd())
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesUsersCmd())
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesExportCmd())
}

func runGroupsList(ctx context.Context, client *api.Client, opts *options.GroupsListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.list", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"user_id":    opts.UserID,
	})

	service := api.NewGroupsService(client)

	apiOpts := &api.ListGroupsOptions{}
	if opts.IncludeUsers {
		apiOpts.Include = append(apiOpts.Include, "users")
	}
	if opts.IncludePermissions {
		apiOpts.Include = append(apiOpts.Include, "permissions")
	}

	var groups []api.Group
	var err error

	if opts.CourseID > 0 {
		groups, err = service.ListCourse(ctx, opts.CourseID, apiOpts)
	} else if opts.AccountID > 0 {
		groups, err = service.ListAccount(ctx, opts.AccountID, apiOpts)
	} else {
		// Default to user's groups (userID 0 means "self")
		groups, err = service.ListUser(ctx, opts.UserID, apiOpts)
	}

	if err != nil {
		logger.LogCommandError(ctx, "groups.list", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
			"user_id":    opts.UserID,
		})
		return fmt.Errorf("failed to list groups: %w", err)
	}

	if err := formatEmptyOrOutput(groups, "No groups found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "groups.list", len(groups))
	return nil
}

func runGroupsGet(ctx context.Context, client *api.Client, opts *options.GroupsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.get", map[string]interface{}{
		"group_id": opts.GroupID,
	})

	service := api.NewGroupsService(client)

	var include []string
	if opts.IncludeUsers {
		include = append(include, "users")
	}
	if opts.IncludePermissions {
		include = append(include, "permissions")
	}

	group, err := service.Get(ctx, opts.GroupID, include)
	if err != nil {
		logger.LogCommandError(ctx, "groups.get", err, map[string]interface{}{
			"group_id": opts.GroupID,
		})
		return fmt.Errorf("failed to get group: %w", err)
	}

	if err := formatOutput(group, nil); err != nil {
		return fmt.Errorf("failed to print result: %w", err)
	}

	logger.LogCommandComplete(ctx, "groups.get", 1)
	return nil
}

func runGroupsCreate(ctx context.Context, client *api.Client, opts *options.GroupsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.create", map[string]interface{}{
		"category_id": opts.CategoryID,
		"name":        opts.Name,
	})

	service := api.NewGroupsService(client)

	params := &api.CreateGroupParams{
		Name:           opts.Name,
		Description:    opts.Description,
		IsPublic:       opts.IsPublic,
		JoinLevel:      opts.JoinLevel,
		StorageQuotaMb: opts.StorageQuotaMb,
		SISGroupID:     opts.SISGroupID,
	}

	group, err := service.Create(ctx, opts.CategoryID, params)
	if err != nil {
		logger.LogCommandError(ctx, "groups.create", err, map[string]interface{}{
			"category_id": opts.CategoryID,
			"name":        opts.Name,
		})
		return fmt.Errorf("failed to create group: %w", err)
	}

	printInfo("Group created successfully (ID: %d)\n", group.ID)
	if err := formatOutput(group, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.create", 1)
	return nil
}

func runGroupsUpdate(ctx context.Context, client *api.Client, opts *options.GroupsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.update", map[string]interface{}{
		"group_id": opts.GroupID,
	})

	service := api.NewGroupsService(client)

	params := &api.UpdateGroupParams{}

	if opts.NameSet {
		params.Name = &opts.Name
	}
	if opts.DescriptionSet {
		params.Description = &opts.Description
	}
	if opts.IsPublicSet {
		params.IsPublic = &opts.IsPublic
	}
	if opts.JoinLevelSet {
		params.JoinLevel = &opts.JoinLevel
	}
	if opts.StorageQuotaMbSet {
		params.StorageQuotaMb = &opts.StorageQuotaMb
	}
	if opts.SISGroupIDSet {
		params.SISGroupID = &opts.SISGroupID
	}

	group, err := service.Update(ctx, opts.GroupID, params)
	if err != nil {
		logger.LogCommandError(ctx, "groups.update", err, map[string]interface{}{
			"group_id": opts.GroupID,
		})
		return fmt.Errorf("failed to update group: %w", err)
	}

	printInfo("Group updated successfully (ID: %d)\n", group.ID)
	if err := formatOutput(group, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.update", 1)
	return nil
}

func runGroupsDelete(ctx context.Context, client *api.Client, opts *options.GroupsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.delete", map[string]interface{}{
		"group_id": opts.GroupID,
	})

	// Confirm deletion
	confirmed, err := confirmDelete("group", opts.GroupID, opts.Force)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Cancelled.")
		return nil
	}

	service := api.NewGroupsService(client)

	group, err := service.Delete(ctx, opts.GroupID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.delete", err, map[string]interface{}{
			"group_id": opts.GroupID,
		})
		return fmt.Errorf("failed to delete group: %w", err)
	}

	printInfo("Group %d deleted successfully\n", group.ID)

	logger.LogCommandComplete(ctx, "groups.delete", 1)
	return nil
}

func runGroupsMembersList(ctx context.Context, client *api.Client, opts *options.GroupsMembersListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.members.list", map[string]interface{}{
		"group_id": opts.GroupID,
	})

	service := api.NewGroupsService(client)

	users, err := service.ListMembers(ctx, opts.GroupID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.members.list", err, map[string]interface{}{
			"group_id": opts.GroupID,
		})
		return fmt.Errorf("failed to list group members: %w", err)
	}

	if err := formatEmptyOrOutput(users, "No members found in group"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "groups.members.list", len(users))
	return nil
}

func runGroupsMembersAdd(ctx context.Context, client *api.Client, opts *options.GroupsMembersAddOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.members.add", map[string]interface{}{
		"group_id": opts.GroupID,
		"user_id":  opts.UserID,
	})

	service := api.NewGroupsService(client)

	membership, err := service.AddMember(ctx, opts.GroupID, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.members.add", err, map[string]interface{}{
			"group_id": opts.GroupID,
			"user_id":  opts.UserID,
		})
		return fmt.Errorf("failed to add member: %w", err)
	}

	fmt.Printf("User %d added to group %d\n", opts.UserID, opts.GroupID)
	if err := formatOutput(membership, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.members.add", 1)
	return nil
}

func runGroupsMembersRemove(ctx context.Context, client *api.Client, opts *options.GroupsMembersRemoveOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.members.remove", map[string]interface{}{
		"group_id":      opts.GroupID,
		"membership_id": opts.MembershipID,
	})

	service := api.NewGroupsService(client)

	err := service.RemoveMember(ctx, opts.GroupID, opts.MembershipID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.members.remove", err, map[string]interface{}{
			"group_id":      opts.GroupID,
			"membership_id": opts.MembershipID,
		})
		return fmt.Errorf("failed to remove member: %w", err)
	}

	fmt.Printf("Membership %d removed from group %d\n", opts.MembershipID, opts.GroupID)

	logger.LogCommandComplete(ctx, "groups.members.remove", 1)
	return nil
}

func runGroupsCategoriesList(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.categories.list", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
	})

	service := api.NewGroupsService(client)

	var categories []api.GroupCategory
	var err error

	if opts.CourseID > 0 {
		categories, err = service.ListCategoriesCourse(ctx, opts.CourseID, nil)
	} else {
		categories, err = service.ListCategoriesAccount(ctx, opts.AccountID, nil)
	}

	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.list", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
		})
		return fmt.Errorf("failed to list group categories: %w", err)
	}

	if err := formatEmptyOrOutput(categories, "No group categories found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "groups.categories.list", len(categories))
	return nil
}

func runGroupsCategoriesGet(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.categories.get", map[string]interface{}{
		"category_id": opts.CategoryID,
	})

	service := api.NewGroupsService(client)

	category, err := service.GetCategory(ctx, opts.CategoryID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.get", err, map[string]interface{}{
			"category_id": opts.CategoryID,
		})
		return fmt.Errorf("failed to get category: %w", err)
	}

	if err := formatOutput(category, nil); err != nil {
		return fmt.Errorf("failed to print result: %w", err)
	}

	logger.LogCommandComplete(ctx, "groups.categories.get", 1)
	return nil
}

func runGroupsCategoriesCreate(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.categories.create", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"name":       opts.Name,
	})

	service := api.NewGroupsService(client)

	params := &api.CreateCategoryParams{
		Name:               opts.Name,
		SelfSignup:         opts.SelfSignup,
		AutoLeader:         opts.AutoLeader,
		GroupLimit:         opts.GroupLimit,
		CreateGroupCount:   opts.CreateGroupCount,
		SplitGroupCount:    opts.SplitGroupCount,
		SISGroupCategoryID: opts.SISCategoryID,
	}

	var category *api.GroupCategory
	var err error

	if opts.CourseID > 0 {
		category, err = service.CreateCategoryCourse(ctx, opts.CourseID, params)
	} else {
		category, err = service.CreateCategoryAccount(ctx, opts.AccountID, params)
	}

	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.create", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
			"name":       opts.Name,
		})
		return fmt.Errorf("failed to create category: %w", err)
	}

	printInfo("Category created successfully (ID: %d)\n", category.ID)
	if err := formatOutput(category, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.categories.create", 1)
	return nil
}

func runGroupsCategoriesUpdate(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.categories.update", map[string]interface{}{
		"category_id": opts.CategoryID,
	})

	service := api.NewGroupsService(client)

	params := &api.UpdateCategoryParams{}

	if opts.NameSet {
		params.Name = &opts.Name
	}
	if opts.SelfSignupSet {
		params.SelfSignup = &opts.SelfSignup
	}
	if opts.AutoLeaderSet {
		params.AutoLeader = &opts.AutoLeader
	}
	if opts.GroupLimitSet {
		params.GroupLimit = &opts.GroupLimit
	}
	if opts.SISCategoryIDSet {
		params.SISGroupCategoryID = &opts.SISCategoryID
	}

	category, err := service.UpdateCategory(ctx, opts.CategoryID, params)
	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.update", err, map[string]interface{}{
			"category_id": opts.CategoryID,
		})
		return fmt.Errorf("failed to update category: %w", err)
	}

	printInfo("Category updated successfully (ID: %d)\n", category.ID)
	if err := formatOutput(category, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.categories.update", 1)
	return nil
}

func runGroupsCategoriesDelete(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.categories.delete", map[string]interface{}{
		"category_id": opts.CategoryID,
	})

	// Confirm deletion
	ok, confirmErr := confirmDeleteWithDetails("group category", opts.CategoryID, map[string]interface{}{
		"warning": "This will delete the category and all groups in it",
	}, opts.Force)
	if confirmErr != nil {
		return confirmErr
	}
	if !ok {
		return nil
	}

	service := api.NewGroupsService(client)

	category, err := service.DeleteCategory(ctx, opts.CategoryID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.delete", err, map[string]interface{}{
			"category_id": opts.CategoryID,
		})
		return fmt.Errorf("failed to delete category: %w", err)
	}

	printInfo("Category %d deleted successfully\n", category.ID)

	logger.LogCommandComplete(ctx, "groups.categories.delete", 1)
	return nil
}

func runGroupsCategoriesGroups(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesGroupsOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.categories.groups", map[string]interface{}{
		"category_id": opts.CategoryID,
	})

	service := api.NewGroupsService(client)

	groups, err := service.ListGroupsInCategory(ctx, opts.CategoryID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.groups", err, map[string]interface{}{
			"category_id": opts.CategoryID,
		})
		return fmt.Errorf("failed to list groups in category: %w", err)
	}

	if err := formatEmptyOrOutput(groups, "No groups found in category"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "groups.categories.groups", len(groups))
	return nil
}

// newGroupsCreateStandaloneCmd creates the groups create-standalone command
func newGroupsCreateStandaloneCmd() *cobra.Command {
	opts := &options.GroupsCreateStandaloneOptions{}

	cmd := &cobra.Command{
		Use:   "create-standalone",
		Short: "Create a standalone group (account-level, no category required)",
		Long: `Create a new group without a category (account-level group creation).

Examples:
  canvas groups create-standalone --name "Staff Group"
  canvas groups create-standalone --name "Committee" --description "Planning committee"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsCreateStandalone(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Group name (required)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Group description")
	cmd.Flags().BoolVar(&opts.IsPublic, "public", false, "Whether the group is public")
	cmd.Flags().StringVar(&opts.JoinLevel, "join-level", "", "Join level (parent_context_auto_join, parent_context_request, invitation_only)")
	cmd.Flags().Int64Var(&opts.StorageQuotaMb, "storage-quota-mb", 0, "Storage quota in MB")
	cmd.Flags().StringVar(&opts.SISGroupID, "sis-group-id", "", "SIS group ID")
	mustMarkRequired(cmd, "name")

	return cmd
}

// newGroupsActivityStreamCmd creates the groups activity-stream command
func newGroupsActivityStreamCmd() *cobra.Command {
	opts := &options.GroupsActivityStreamOptions{}

	cmd := &cobra.Command{
		Use:   "activity-stream <group-id>",
		Short: "Get group activity stream",
		Long: `Get the activity stream for a group.

Examples:
  canvas groups activity-stream 456
  canvas groups activity-stream 456 --summary`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}

			summary, _ := cmd.Flags().GetBool("summary")
			if summary {
				return runGroupsActivityStreamSummary(cmd.Context(), client, opts)
			}
			return runGroupsActivityStream(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Bool("summary", false, "Return activity stream summary instead of full stream")
	return cmd
}

// newGroupsPermissionsCmd creates the groups permissions command
func newGroupsPermissionsCmd() *cobra.Command {
	opts := &options.GroupsPermissionsOptions{}

	cmd := &cobra.Command{
		Use:   "permissions <group-id>",
		Short: "Get group permissions for current user",
		Long: `Get the permissions the current user has for a group.

Examples:
  canvas groups permissions 456`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsPermissions(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.Permissions, "permissions", nil, "Permission names to check")
	return cmd
}

// newGroupsInviteCmd creates the groups invite command
func newGroupsInviteCmd() *cobra.Command {
	opts := &options.GroupsInviteOptions{}

	cmd := &cobra.Command{
		Use:   "invite <group-id>",
		Short: "Invite users to a group",
		Long: `Invite one or more users (by email) to join a group.

Examples:
  canvas groups invite 456 --invitees user1@example.com,user2@example.com`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsInvite(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.Invitees, "invitees", nil, "Comma-separated list of email addresses to invite (required)")
	mustMarkRequired(cmd, "invitees")
	return cmd
}

// newGroupsTabsCmd creates the groups tabs command
func newGroupsTabsCmd() *cobra.Command {
	opts := &options.GroupsTabsListOptions{}

	cmd := &cobra.Command{
		Use:   "tabs <group-id>",
		Short: "List group navigation tabs",
		Long: `List the navigation tabs for a group.

Examples:
  canvas groups tabs 456`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsTabs(cmd.Context(), client, opts)
		},
	}

	return cmd
}

// newGroupsCollaborationsCmd creates the groups collaborations command
func newGroupsCollaborationsCmd() *cobra.Command {
	opts := &options.GroupsCollaborationsListOptions{}

	cmd := &cobra.Command{
		Use:   "collaborations <group-id>",
		Short: "List group collaborations",
		Long: `List collaborations for a group.

Examples:
  canvas groups collaborations 456`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsCollaborations(cmd.Context(), client, opts)
		},
	}

	return cmd
}

// newGroupsConferencesCmd creates the groups conferences command
func newGroupsConferencesCmd() *cobra.Command {
	opts := &options.GroupsConferencesListOptions{}

	cmd := &cobra.Command{
		Use:   "conferences <group-id>",
		Short: "List group conferences",
		Long: `List conferences for a group.

Examples:
  canvas groups conferences 456`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsConferences(cmd.Context(), client, opts)
		},
	}

	return cmd
}

// newGroupsAssignmentOverrideCmd creates the groups assignment-override command
func newGroupsAssignmentOverrideCmd() *cobra.Command {
	opts := &options.GroupsAssignmentOverrideOptions{}

	cmd := &cobra.Command{
		Use:   "assignment-override <group-id>",
		Short: "Get assignment override for a group",
		Long: `Get the assignment override for a specific assignment and group.

Examples:
  canvas groups assignment-override 456 --assignment-id 789`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsAssignmentOverride(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	mustMarkRequired(cmd, "assignment-id")
	return cmd
}

// newGroupsMembershipsListCmd creates the groups memberships list command
func newGroupsMembershipsListCmd() *cobra.Command {
	opts := &options.GroupsMembershipsListOptions{}

	cmd := &cobra.Command{
		Use:   "list <group-id>",
		Short: "List group membership records",
		Long: `List all membership records for a group.

Examples:
  canvas groups memberships list 456`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsMembershipsList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

// newGroupsMembershipsGetCmd creates the groups memberships get command
func newGroupsMembershipsGetCmd() *cobra.Command {
	opts := &options.GroupsMembershipsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <group-id>",
		Short: "Get a single membership record",
		Long: `Get a specific membership record by membership ID.

Examples:
  canvas groups memberships get 456 --membership-id 123`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsMembershipsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.MembershipID, "membership-id", 0, "Membership ID (required)")
	mustMarkRequired(cmd, "membership-id")
	return cmd
}

// newGroupsMembershipsUpdateCmd creates the groups memberships update command
func newGroupsMembershipsUpdateCmd() *cobra.Command {
	opts := &options.GroupsMembershipsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <group-id>",
		Short: "Update a membership record",
		Long: `Update a specific membership record (e.g. change workflow state or promote to moderator).

Examples:
  canvas groups memberships update 456 --membership-id 123 --moderator
  canvas groups memberships update 456 --membership-id 123 --workflow-state accepted`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID
			opts.WorkflowStateSet = cmd.Flags().Changed("workflow-state")
			opts.ModeratorSet = cmd.Flags().Changed("moderator")

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsMembershipsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.MembershipID, "membership-id", 0, "Membership ID (required)")
	cmd.Flags().StringVar(&opts.WorkflowState, "workflow-state", "", "Workflow state (accepted, invited, rejected)")
	cmd.Flags().BoolVar(&opts.Moderator, "moderator", false, "Set as moderator")
	mustMarkRequired(cmd, "membership-id")
	return cmd
}

// newGroupsUsersGetCmd creates the groups users get command
func newGroupsUsersGetCmd() *cobra.Command {
	opts := &options.GroupsUsersGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <group-id>",
		Short: "Get a user's membership details in a group",
		Long: `Get a specific user's membership details within a group.

Examples:
  canvas groups users get 456 --user-id 789`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsUsersGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	mustMarkRequired(cmd, "user-id")
	return cmd
}

// newGroupsUsersUpdateCmd creates the groups users update command
func newGroupsUsersUpdateCmd() *cobra.Command {
	opts := &options.GroupsUsersUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <group-id>",
		Short: "Update a user's membership in a group",
		Long: `Update a specific user's membership state or moderator status.

Examples:
  canvas groups users update 456 --user-id 789 --moderator
  canvas groups users update 456 --user-id 789 --workflow-state accepted`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID
			opts.WorkflowStateSet = cmd.Flags().Changed("workflow-state")
			opts.ModeratorSet = cmd.Flags().Changed("moderator")

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsUsersUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.Flags().StringVar(&opts.WorkflowState, "workflow-state", "", "Workflow state (accepted, invited, rejected)")
	cmd.Flags().BoolVar(&opts.Moderator, "moderator", false, "Set as moderator")
	mustMarkRequired(cmd, "user-id")
	return cmd
}

// newGroupsUsersRemoveCmd creates the groups users remove command
func newGroupsUsersRemoveCmd() *cobra.Command {
	opts := &options.GroupsUsersRemoveOptions{}

	cmd := &cobra.Command{
		Use:   "remove <group-id>",
		Short: "Remove a specific user from a group",
		Long: `Remove a user from a group by user ID.

Examples:
  canvas groups users remove 456 --user-id 789`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsUsersRemove(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	mustMarkRequired(cmd, "user-id")
	return cmd
}

// newGroupsUsersRemoveSelfCmd creates the groups users remove-self command
func newGroupsUsersRemoveSelfCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove-self <group-id>",
		Short: "Leave a group (remove self)",
		Long: `Remove the current user from a group.

Examples:
  canvas groups users remove-self 456`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsUsersRemoveSelf(cmd.Context(), client, groupID)
		},
	}
}

// newGroupsExternalFeedsListCmd creates the groups external-feeds list command
func newGroupsExternalFeedsListCmd() *cobra.Command {
	opts := &options.GroupsExternalFeedsListOptions{}

	cmd := &cobra.Command{
		Use:   "list <group-id>",
		Short: "List group external feeds",
		Args:  ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsExternalFeedsList(cmd.Context(), client, opts)
		},
	}
	return cmd
}

// newGroupsExternalFeedsCreateCmd creates the groups external-feeds create command
func newGroupsExternalFeedsCreateCmd() *cobra.Command {
	opts := &options.GroupsExternalFeedsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <group-id>",
		Short: "Add an external feed to a group",
		Args:  ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsExternalFeedsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.URL, "url", "", "Feed URL (required)")
	cmd.Flags().StringVar(&opts.Verbosity, "verbosity", "", "Verbosity (full, truncate, link_only)")
	cmd.Flags().StringVar(&opts.Header, "header", "", "Header match string")
	mustMarkRequired(cmd, "url")
	return cmd
}

// newGroupsExternalFeedsDeleteCmd creates the groups external-feeds delete command
func newGroupsExternalFeedsDeleteCmd() *cobra.Command {
	opts := &options.GroupsExternalFeedsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <group-id>",
		Short: "Delete an external feed from a group",
		Args:  ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsExternalFeedsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.FeedID, "feed-id", 0, "External feed ID (required)")
	mustMarkRequired(cmd, "feed-id")
	return cmd
}

// newGroupsContentExportsListCmd creates the groups content-exports list command
func newGroupsContentExportsListCmd() *cobra.Command {
	opts := &options.GroupsContentExportsListOptions{}

	cmd := &cobra.Command{
		Use:   "list <group-id>",
		Short: "List content exports for a group",
		Args:  ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsContentExportsList(cmd.Context(), client, opts)
		},
	}
	return cmd
}

// newGroupsContentExportsCreateCmd creates the groups content-exports create command
func newGroupsContentExportsCreateCmd() *cobra.Command {
	opts := &options.GroupsContentExportsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <group-id>",
		Short: "Start a content export for a group",
		Args:  ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsContentExportsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.ExportType, "export-type", "", "Export type e.g. common_cartridge (required)")
	mustMarkRequired(cmd, "export-type")
	return cmd
}

// newGroupsContentExportsGetCmd creates the groups content-exports get command
func newGroupsContentExportsGetCmd() *cobra.Command {
	opts := &options.GroupsContentExportsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <group-id>",
		Short: "Get a content export for a group",
		Args:  ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsContentExportsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.ExportID, "export-id", 0, "Export ID (required)")
	mustMarkRequired(cmd, "export-id")
	return cmd
}

// newGroupsCategoriesAssignMembersCmd creates the categories assign-members command
func newGroupsCategoriesAssignMembersCmd() *cobra.Command {
	opts := &options.GroupsCategoriesAssignMembersOptions{}

	cmd := &cobra.Command{
		Use:   "assign-members <category-id>",
		Short: "Assign unassigned members to groups in a category",
		Long: `Assign all unassigned members to groups within a group category.

Examples:
  canvas groups categories assign-members 123
  canvas groups categories assign-members 123 --sync`,
		Args: ExactArgsWithUsage(1, "category-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			categoryID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid category ID: %w", err)
			}
			opts.CategoryID = categoryID

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsCategoriesAssignMembers(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Sync, "sync", false, "Run synchronously (wait for completion)")
	return cmd
}

// newGroupsCategoriesUsersCmd creates the categories users command
func newGroupsCategoriesUsersCmd() *cobra.Command {
	opts := &options.GroupsCategoriesUsersListOptions{}

	cmd := &cobra.Command{
		Use:   "users <category-id>",
		Short: "List users in a group category",
		Long: `List users who are enrolled in groups within a category.

Examples:
  canvas groups categories users 123`,
		Args: ExactArgsWithUsage(1, "category-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			categoryID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid category ID: %w", err)
			}
			opts.CategoryID = categoryID

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsCategoriesUsers(cmd.Context(), client, opts)
		},
	}

	return cmd
}

// newGroupsCategoriesExportCmd creates the categories export command
func newGroupsCategoriesExportCmd() *cobra.Command {
	opts := &options.GroupsCategoriesExportOptions{}

	cmd := &cobra.Command{
		Use:   "export <category-id>",
		Short: "Export a group category",
		Long: `Export the membership data for a group category.

Examples:
  canvas groups categories export 123`,
		Args: ExactArgsWithUsage(1, "category-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			categoryID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid category ID: %w", err)
			}
			opts.CategoryID = categoryID

			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGroupsCategoriesExport(cmd.Context(), client, opts)
		},
	}

	return cmd
}

// Runner functions for new commands

func runGroupsCreateStandalone(ctx context.Context, client *api.Client, opts *options.GroupsCreateStandaloneOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.create-standalone", map[string]interface{}{"name": opts.Name})

	service := api.NewGroupsService(client)
	params := &api.CreateGroupParams{
		Name:           opts.Name,
		Description:    opts.Description,
		IsPublic:       opts.IsPublic,
		JoinLevel:      opts.JoinLevel,
		StorageQuotaMb: opts.StorageQuotaMb,
		SISGroupID:     opts.SISGroupID,
	}

	group, err := service.CreateStandalone(ctx, params)
	if err != nil {
		logger.LogCommandError(ctx, "groups.create-standalone", err, nil)
		return fmt.Errorf("failed to create standalone group: %w", err)
	}

	printInfo("Group created successfully (ID: %d)\n", group.ID)
	if err := formatOutput(group, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.create-standalone", 1)
	return nil
}

func runGroupsActivityStream(ctx context.Context, client *api.Client, opts *options.GroupsActivityStreamOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.activity-stream", map[string]interface{}{"group_id": opts.GroupID})

	service := api.NewGroupsService(client)
	items, err := service.GetActivityStream(ctx, opts.GroupID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.activity-stream", err, nil)
		return fmt.Errorf("failed to get activity stream: %w", err)
	}

	if err := formatEmptyOrOutput(items, "No activity in stream"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.activity-stream", len(items))
	return nil
}

func runGroupsActivityStreamSummary(ctx context.Context, client *api.Client, opts *options.GroupsActivityStreamOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.activity-stream.summary", map[string]interface{}{"group_id": opts.GroupID})

	service := api.NewGroupsService(client)
	items, err := service.GetActivityStreamSummary(ctx, opts.GroupID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.activity-stream.summary", err, nil)
		return fmt.Errorf("failed to get activity stream summary: %w", err)
	}

	if err := formatEmptyOrOutput(items, "No activity summary"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.activity-stream.summary", len(items))
	return nil
}

func runGroupsPermissions(ctx context.Context, client *api.Client, opts *options.GroupsPermissionsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.permissions", map[string]interface{}{"group_id": opts.GroupID})

	service := api.NewGroupsService(client)
	perms, err := service.GetPermissions(ctx, opts.GroupID, opts.Permissions)
	if err != nil {
		logger.LogCommandError(ctx, "groups.permissions", err, nil)
		return fmt.Errorf("failed to get permissions: %w", err)
	}

	if err := formatOutput(perms, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.permissions", 1)
	return nil
}

func runGroupsInvite(ctx context.Context, client *api.Client, opts *options.GroupsInviteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.invite", map[string]interface{}{
		"group_id": opts.GroupID,
		"count":    len(opts.Invitees),
	})

	service := api.NewGroupsService(client)
	memberships, err := service.Invite(ctx, opts.GroupID, opts.Invitees)
	if err != nil {
		logger.LogCommandError(ctx, "groups.invite", err, nil)
		return fmt.Errorf("failed to invite users: %w", err)
	}

	printInfo("Invited %d user(s) to group %d\n", len(memberships), opts.GroupID)
	if err := formatOutput(memberships, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.invite", len(memberships))
	return nil
}

func runGroupsTabs(ctx context.Context, client *api.Client, opts *options.GroupsTabsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.tabs", map[string]interface{}{"group_id": opts.GroupID})

	service := api.NewGroupsService(client)
	tabs, err := service.ListTabs(ctx, opts.GroupID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.tabs", err, nil)
		return fmt.Errorf("failed to list tabs: %w", err)
	}

	if err := formatEmptyOrOutput(tabs, "No tabs found"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.tabs", len(tabs))
	return nil
}

func runGroupsCollaborations(ctx context.Context, client *api.Client, opts *options.GroupsCollaborationsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.collaborations", map[string]interface{}{"group_id": opts.GroupID})

	service := api.NewGroupsService(client)
	items, err := service.ListCollaborations(ctx, opts.GroupID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.collaborations", err, nil)
		return fmt.Errorf("failed to list collaborations: %w", err)
	}

	if err := formatEmptyOrOutput(items, "No collaborations found"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.collaborations", len(items))
	return nil
}

func runGroupsConferences(ctx context.Context, client *api.Client, opts *options.GroupsConferencesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.conferences", map[string]interface{}{"group_id": opts.GroupID})

	service := api.NewGroupsService(client)
	items, err := service.ListConferences(ctx, opts.GroupID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.conferences", err, nil)
		return fmt.Errorf("failed to list conferences: %w", err)
	}

	if err := formatEmptyOrOutput(items, "No conferences found"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.conferences", len(items))
	return nil
}

func runGroupsAssignmentOverride(ctx context.Context, client *api.Client, opts *options.GroupsAssignmentOverrideOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.assignment-override", map[string]interface{}{
		"group_id":      opts.GroupID,
		"assignment_id": opts.AssignmentID,
	})

	service := api.NewGroupsService(client)
	result, err := service.GetAssignmentOverride(ctx, opts.GroupID, opts.AssignmentID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.assignment-override", err, nil)
		return fmt.Errorf("failed to get assignment override: %w", err)
	}

	if err := formatOutput(result, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.assignment-override", 1)
	return nil
}

func runGroupsMembershipsList(ctx context.Context, client *api.Client, opts *options.GroupsMembershipsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.memberships.list", map[string]interface{}{"group_id": opts.GroupID})

	service := api.NewGroupsService(client)
	memberships, err := service.ListMemberships(ctx, opts.GroupID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.memberships.list", err, nil)
		return fmt.Errorf("failed to list memberships: %w", err)
	}

	if err := formatEmptyOrOutput(memberships, "No memberships found"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.memberships.list", len(memberships))
	return nil
}

func runGroupsMembershipsGet(ctx context.Context, client *api.Client, opts *options.GroupsMembershipsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.memberships.get", map[string]interface{}{
		"group_id":      opts.GroupID,
		"membership_id": opts.MembershipID,
	})

	service := api.NewGroupsService(client)
	membership, err := service.GetMembership(ctx, opts.GroupID, opts.MembershipID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.memberships.get", err, nil)
		return fmt.Errorf("failed to get membership: %w", err)
	}

	if err := formatOutput(membership, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.memberships.get", 1)
	return nil
}

func runGroupsMembershipsUpdate(ctx context.Context, client *api.Client, opts *options.GroupsMembershipsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.memberships.update", map[string]interface{}{
		"group_id":      opts.GroupID,
		"membership_id": opts.MembershipID,
	})

	service := api.NewGroupsService(client)
	params := &api.UpdateMembershipParams{}
	if opts.WorkflowStateSet {
		params.WorkflowState = &opts.WorkflowState
	}
	if opts.ModeratorSet {
		params.Moderator = &opts.Moderator
	}

	membership, err := service.UpdateMembership(ctx, opts.GroupID, opts.MembershipID, params)
	if err != nil {
		logger.LogCommandError(ctx, "groups.memberships.update", err, nil)
		return fmt.Errorf("failed to update membership: %w", err)
	}

	printInfo("Membership updated (ID: %d)\n", membership.ID)
	if err := formatOutput(membership, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.memberships.update", 1)
	return nil
}

func runGroupsUsersGet(ctx context.Context, client *api.Client, opts *options.GroupsUsersGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.users.get", map[string]interface{}{
		"group_id": opts.GroupID,
		"user_id":  opts.UserID,
	})

	service := api.NewGroupsService(client)
	user, err := service.GetUser(ctx, opts.GroupID, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.users.get", err, nil)
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := formatOutput(user, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.users.get", 1)
	return nil
}

func runGroupsUsersUpdate(ctx context.Context, client *api.Client, opts *options.GroupsUsersUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.users.update", map[string]interface{}{
		"group_id": opts.GroupID,
		"user_id":  opts.UserID,
	})

	service := api.NewGroupsService(client)
	params := &api.UpdateMembershipParams{}
	if opts.WorkflowStateSet {
		params.WorkflowState = &opts.WorkflowState
	}
	if opts.ModeratorSet {
		params.Moderator = &opts.Moderator
	}

	membership, err := service.UpdateUserMembership(ctx, opts.GroupID, opts.UserID, params)
	if err != nil {
		logger.LogCommandError(ctx, "groups.users.update", err, nil)
		return fmt.Errorf("failed to update user membership: %w", err)
	}

	printInfo("User membership updated (user_id: %d)\n", opts.UserID)
	if err := formatOutput(membership, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.users.update", 1)
	return nil
}

func runGroupsUsersRemove(ctx context.Context, client *api.Client, opts *options.GroupsUsersRemoveOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.users.remove", map[string]interface{}{
		"group_id": opts.GroupID,
		"user_id":  opts.UserID,
	})

	service := api.NewGroupsService(client)
	if err := service.RemoveUser(ctx, opts.GroupID, opts.UserID); err != nil {
		logger.LogCommandError(ctx, "groups.users.remove", err, nil)
		return fmt.Errorf("failed to remove user from group: %w", err)
	}

	fmt.Printf("User %d removed from group %d\n", opts.UserID, opts.GroupID)
	logger.LogCommandComplete(ctx, "groups.users.remove", 1)
	return nil
}

func runGroupsUsersRemoveSelf(ctx context.Context, client *api.Client, groupID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.users.remove-self", map[string]interface{}{"group_id": groupID})

	service := api.NewGroupsService(client)
	if err := service.RemoveUserBySelf(ctx, groupID); err != nil {
		logger.LogCommandError(ctx, "groups.users.remove-self", err, nil)
		return fmt.Errorf("failed to leave group: %w", err)
	}

	fmt.Printf("Left group %d\n", groupID)
	logger.LogCommandComplete(ctx, "groups.users.remove-self", 1)
	return nil
}

func runGroupsExternalFeedsList(ctx context.Context, client *api.Client, opts *options.GroupsExternalFeedsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.external-feeds.list", map[string]interface{}{"group_id": opts.GroupID})

	service := api.NewGroupsService(client)
	feeds, err := service.ListExternalFeeds(ctx, opts.GroupID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.external-feeds.list", err, nil)
		return fmt.Errorf("failed to list external feeds: %w", err)
	}

	if err := formatEmptyOrOutput(feeds, "No external feeds found"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.external-feeds.list", len(feeds))
	return nil
}

func runGroupsExternalFeedsCreate(ctx context.Context, client *api.Client, opts *options.GroupsExternalFeedsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.external-feeds.create", map[string]interface{}{
		"group_id": opts.GroupID,
		"url":      opts.URL,
	})

	service := api.NewGroupsService(client)
	params := &api.CreateExternalFeedParams{
		URL:       opts.URL,
		Verbosity: opts.Verbosity,
		Header:    opts.Header,
	}
	feed, err := service.CreateExternalFeed(ctx, opts.GroupID, params)
	if err != nil {
		logger.LogCommandError(ctx, "groups.external-feeds.create", err, nil)
		return fmt.Errorf("failed to create external feed: %w", err)
	}

	printInfo("External feed created (ID: %d)\n", feed.ID)
	if err := formatOutput(feed, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.external-feeds.create", 1)
	return nil
}

func runGroupsExternalFeedsDelete(ctx context.Context, client *api.Client, opts *options.GroupsExternalFeedsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.external-feeds.delete", map[string]interface{}{
		"group_id": opts.GroupID,
		"feed_id":  opts.FeedID,
	})

	service := api.NewGroupsService(client)
	feed, err := service.DeleteExternalFeed(ctx, opts.GroupID, opts.FeedID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.external-feeds.delete", err, nil)
		return fmt.Errorf("failed to delete external feed: %w", err)
	}

	printInfo("External feed %d deleted\n", feed.ID)
	logger.LogCommandComplete(ctx, "groups.external-feeds.delete", 1)
	return nil
}

func runGroupsContentExportsList(ctx context.Context, client *api.Client, opts *options.GroupsContentExportsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.content-exports.list", map[string]interface{}{"group_id": opts.GroupID})

	service := api.NewGroupsService(client)
	exports, err := service.ListContentExports(ctx, opts.GroupID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.content-exports.list", err, nil)
		return fmt.Errorf("failed to list content exports: %w", err)
	}

	if err := formatEmptyOrOutput(exports, "No content exports found"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.content-exports.list", len(exports))
	return nil
}

func runGroupsContentExportsCreate(ctx context.Context, client *api.Client, opts *options.GroupsContentExportsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.content-exports.create", map[string]interface{}{
		"group_id":    opts.GroupID,
		"export_type": opts.ExportType,
	})

	service := api.NewGroupsService(client)
	export, err := service.CreateContentExport(ctx, opts.GroupID, opts.ExportType)
	if err != nil {
		logger.LogCommandError(ctx, "groups.content-exports.create", err, nil)
		return fmt.Errorf("failed to create content export: %w", err)
	}

	printInfo("Content export started (ID: %d, state: %s)\n", export.ID, export.WorkflowState)
	if err := formatOutput(export, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.content-exports.create", 1)
	return nil
}

func runGroupsContentExportsGet(ctx context.Context, client *api.Client, opts *options.GroupsContentExportsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.content-exports.get", map[string]interface{}{
		"group_id":  opts.GroupID,
		"export_id": opts.ExportID,
	})

	service := api.NewGroupsService(client)
	export, err := service.GetContentExport(ctx, opts.GroupID, opts.ExportID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.content-exports.get", err, nil)
		return fmt.Errorf("failed to get content export: %w", err)
	}

	if err := formatOutput(export, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.content-exports.get", 1)
	return nil
}

func runGroupsCategoriesAssignMembers(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesAssignMembersOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.categories.assign-members", map[string]interface{}{"category_id": opts.CategoryID})

	service := api.NewGroupsService(client)
	result, err := service.AssignUnassignedMembers(ctx, opts.CategoryID, opts.Sync)
	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.assign-members", err, nil)
		return fmt.Errorf("failed to assign unassigned members: %w", err)
	}

	if err := formatOutput(result, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.categories.assign-members", 1)
	return nil
}

func runGroupsCategoriesUsers(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesUsersListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.categories.users", map[string]interface{}{"category_id": opts.CategoryID})

	service := api.NewGroupsService(client)
	users, err := service.ListUsersInCategory(ctx, opts.CategoryID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.users", err, nil)
		return fmt.Errorf("failed to list users in category: %w", err)
	}

	if err := formatEmptyOrOutput(users, "No users found in category"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.categories.users", len(users))
	return nil
}

func runGroupsCategoriesExport(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesExportOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "groups.categories.export", map[string]interface{}{"category_id": opts.CategoryID})

	service := api.NewGroupsService(client)
	result, err := service.ExportCategory(ctx, opts.CategoryID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.export", err, nil)
		return fmt.Errorf("failed to export category: %w", err)
	}

	if err := formatOutput(result, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.categories.export", 1)
	return nil
}
