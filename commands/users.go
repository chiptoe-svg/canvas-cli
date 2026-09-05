package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/progress"
)

// usersCmd represents the users command group
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage Canvas users",
	Long: `Manage Canvas users including listing, viewing, searching, and managing user information.

Examples:
  canvas users list --account-id 1
  canvas users get 123
  canvas users search "john"
  canvas users me`,
}

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(newUsersListCmd())
	usersCmd.AddCommand(newUsersGetCmd())
	usersCmd.AddCommand(newUsersMeCmd())
	usersCmd.AddCommand(newUsersSearchCmd())
}

func newUsersListCmd() *cobra.Command {
	opts := &options.UsersListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users in an account or course",
		Long: `List users in a Canvas account or course.

Specify --account-id or --course-id, or uses default account if configured.

WARNING: Account-level user lists can be very large. Use --limit or --search
to avoid long wait times.

Account context (admin):
  canvas users list --limit 100           # First 100 users (recommended)
  canvas users list --search "john"       # Search in default account

Course context:
  canvas users list --course-id 123       # All users enrolled in course 123
  canvas users list --course-id 123 --enrollment-type teacher

Examples:
  canvas users list --limit 50
  canvas users list --account-id 1 --limit 100
  canvas users list --course-id 123
  canvas users list --search "john"
  canvas users list --include email,enrollments`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUsersList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (for account users)")
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (for course enrollees)")
	cmd.Flags().StringVar(&opts.SearchTerm, "search", "", "Search by name, login ID, or email")
	cmd.Flags().StringVar(&opts.EnrollmentType, "enrollment-type", "", "Filter by enrollment type (student, teacher, ta, observer, designer)")
	cmd.Flags().StringVar(&opts.EnrollmentState, "enrollment-state", "", "Filter by enrollment state (active, invited, rejected, completed, inactive)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")

	return cmd
}

func newUsersGetCmd() *cobra.Command {
	opts := &options.UsersGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <user-id>",
		Short: "Get details of a specific user",
		Long: `Get details of a specific user by ID.

Examples:
  canvas users get 123
  canvas users get 123 --include email,enrollments,avatar_url`,
		Args: ExactArgsWithUsage(1, "user-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			userID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid user ID: %s", args[0])
			}
			opts.UserID = userID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUsersGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")

	return cmd
}

func newUsersMeCmd() *cobra.Command {
	opts := &options.UsersMeOptions{}

	cmd := &cobra.Command{
		Use:   "me",
		Short: "Get details of the current authenticated user",
		Long: `Get details of the current authenticated user.

Examples:
  canvas users me`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUsersMe(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newUsersSearchCmd() *cobra.Command {
	opts := &options.UsersSearchOptions{}

	cmd := &cobra.Command{
		Use:   "search <search-term>",
		Short: "Search for users",
		Long: `Search for users across the Canvas instance.

Examples:
  canvas users search "john doe"
  canvas users search "john@example.com"`,
		Args: ExactArgsWithUsage(1, "search-term"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.SearchTerm = args[0]
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUsersSearch(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func runUsersList(ctx context.Context, client *api.Client, opts *options.UsersListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.list", map[string]interface{}{
		"account_id":       opts.AccountID,
		"course_id":        opts.CourseID,
		"search_term":      opts.SearchTerm,
		"enrollment_type":  opts.EnrollmentType,
		"enrollment_state": opts.EnrollmentState,
	})

	// Use default account ID if neither account nor course is specified
	accountID := opts.AccountID
	if accountID == 0 && opts.CourseID == 0 {
		defaultID, err := getDefaultAccountID()
		if err != nil || defaultID == 0 {
			logger.LogCommandError(ctx, "users.list", fmt.Errorf("no account or course specified"), map[string]interface{}{
				"error": "must specify --account-id or --course-id (no default account configured)",
			})
			return fmt.Errorf("must specify --account-id or --course-id (no default account configured). Use 'canvas config account --detect' to set one")
		}
		accountID = defaultID
		printVerbose("Using default account ID: %d\n", accountID)
	}

	// Create users service
	usersService := api.NewUsersService(client)

	// Build options
	listOpts := &api.ListUsersOptions{
		SearchTerm:      opts.SearchTerm,
		EnrollmentType:  opts.EnrollmentType,
		EnrollmentState: opts.EnrollmentState,
		Include:         opts.Include,
	}

	// List users based on context
	spin := progress.New("Fetching users...")
	if !quiet {
		spin.Start()
	}

	var users []api.User
	var contextName string
	var err error

	if accountID > 0 {
		users, err = usersService.List(ctx, accountID, listOpts)
		contextName = fmt.Sprintf("account %d", accountID)
	} else {
		users, err = usersService.ListCourseUsers(ctx, opts.CourseID, listOpts)
		contextName = fmt.Sprintf("course %d", opts.CourseID)
	}
	spin.Stop()

	if err != nil {
		logger.LogCommandError(ctx, "users.list", err, map[string]interface{}{
			"context": contextName,
		})
		return fmt.Errorf("failed to list users: %w", err)
	}

	if len(users) == 0 {
		fmt.Printf("No users found in %s\n", contextName)
		logger.LogCommandComplete(ctx, "users.list", 0)
		return nil
	}

	// Format and display users
	printVerbose("Found %d users in %s:\n\n", len(users), contextName)
	logger.LogCommandComplete(ctx, "users.list", len(users))

	return formatOutput(users, nil)
}

func runUsersGet(ctx context.Context, client *api.Client, opts *options.UsersGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.get", map[string]interface{}{
		"user_id": opts.UserID,
		"include": opts.Include,
	})

	// Create users service
	usersService := api.NewUsersService(client)

	// Get user
	user, err := usersService.Get(ctx, opts.UserID, opts.Include)
	if err != nil {
		logger.LogCommandError(ctx, "users.get", err, map[string]interface{}{
			"user_id": opts.UserID,
		})
		return fmt.Errorf("failed to get user: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.get", 1)

	// Format and display user details
	return formatOutput(user, nil)
}

func runUsersMe(ctx context.Context, client *api.Client, opts *options.UsersMeOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.me", map[string]interface{}{})

	// Create users service
	usersService := api.NewUsersService(client)

	// Get current user
	user, err := usersService.GetCurrentUser(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "users.me", err, map[string]interface{}{})
		return fmt.Errorf("failed to get current user: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.me", 1)

	// Format and display user details
	return formatOutput(user, nil)
}

func runUsersSearch(ctx context.Context, client *api.Client, opts *options.UsersSearchOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.search", map[string]interface{}{
		"search_term": opts.SearchTerm,
	})

	// Create users service
	usersService := api.NewUsersService(client)

	// Search users
	users, err := usersService.Search(ctx, opts.SearchTerm)
	if err != nil {
		logger.LogCommandError(ctx, "users.search", err, map[string]interface{}{
			"search_term": opts.SearchTerm,
		})
		return fmt.Errorf("failed to search users: %w", err)
	}

	if len(users) == 0 {
		fmt.Printf("No users found matching '%s'\n", opts.SearchTerm)
		logger.LogCommandComplete(ctx, "users.search", 0)
		return nil
	}

	// Format and display users
	printVerbose("Found %d users matching '%s':\n\n", len(users), opts.SearchTerm)
	logger.LogCommandComplete(ctx, "users.search", len(users))

	return formatOutput(users, nil)
}
