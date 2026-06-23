package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

func init() {
	usersCmd.AddCommand(newUsersProfileCmd())
	usersCmd.AddCommand(newUsersSettingsCmd())
	usersCmd.AddCommand(newUsersUpdateSettingsCmd())
	usersCmd.AddCommand(newUsersPageViewsCmd())
	usersCmd.AddCommand(newUsersLoginsCmd())
	usersCmd.AddCommand(newUsersCoursesCmd())
	usersCmd.AddCommand(newUsersMissingSubmissionsCmd())
	usersCmd.AddCommand(newUsersActivityStreamCmd())
	usersCmd.AddCommand(newUsersTodoCmd())
	usersCmd.AddCommand(newUsersUpcomingEventsCmd())
	usersCmd.AddCommand(newUsersMergeCmd())
	usersCmd.AddCommand(newUsersSplitCmd())
}

// newUsersProfileCmd returns the 'users profile' command
func newUsersProfileCmd() *cobra.Command {
	opts := &options.UsersProfileOptions{}
	cmd := &cobra.Command{
		Use:   "profile <user-id>",
		Short: "Get a user's profile",
		Long: `Get a user's extended profile information.

Examples:
  canvas users profile 123`,
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
			return runUsersProfile(cmd.Context(), client, opts)
		},
	}
	return cmd
}

func runUsersProfile(ctx context.Context, client *api.Client, opts *options.UsersProfileOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.profile", map[string]interface{}{"user_id": opts.UserID})

	svc := api.NewUsersService(client)
	profile, err := svc.GetProfile(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "users.profile", err, map[string]interface{}{"user_id": opts.UserID})
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.profile", 1)
	return formatOutput(profile, nil)
}

// newUsersSettingsCmd returns the 'users settings' command
func newUsersSettingsCmd() *cobra.Command {
	opts := &options.UsersSettingsOptions{}
	cmd := &cobra.Command{
		Use:   "settings <user-id>",
		Short: "Get user settings",
		Long: `Get settings for a user.

Examples:
  canvas users settings 123`,
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
			return runUsersSettings(cmd.Context(), client, opts)
		},
	}
	return cmd
}

func runUsersSettings(ctx context.Context, client *api.Client, opts *options.UsersSettingsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.settings", map[string]interface{}{"user_id": opts.UserID})

	svc := api.NewUsersService(client)
	settings, err := svc.GetSettings(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "users.settings", err, map[string]interface{}{"user_id": opts.UserID})
		return fmt.Errorf("failed to get user settings: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.settings", 1)
	return formatOutput(settings, nil)
}

// newUsersUpdateSettingsCmd returns the 'users update-settings' command
func newUsersUpdateSettingsCmd() *cobra.Command {
	opts := &options.UsersUpdateSettingsOptions{}
	cmd := &cobra.Command{
		Use:   "update-settings <user-id>",
		Short: "Update user settings",
		Long: `Update settings for a user.

Examples:
  canvas users update-settings 123 --manual-mark-as-read`,
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
			return runUsersUpdateSettings(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().BoolVar(&opts.ManualMarkAsRead, "manual-mark-as-read", false, "Manually mark items as read")
	cmd.Flags().BoolVar(&opts.CollapseGlobalNav, "collapse-global-nav", false, "Collapse global navigation")
	return cmd
}

func runUsersUpdateSettings(ctx context.Context, client *api.Client, opts *options.UsersUpdateSettingsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.update-settings", map[string]interface{}{"user_id": opts.UserID})

	svc := api.NewUsersService(client)
	params := api.UpdateUserSettingsParams{}
	if opts.ManualMarkAsRead {
		v := true
		params.ManualMarkAsRead = &v
	}
	if opts.CollapseGlobalNav {
		v := true
		params.CollapseGlobalNav = &v
	}
	settings, err := svc.UpdateSettings(ctx, opts.UserID, params)
	if err != nil {
		logger.LogCommandError(ctx, "users.update-settings", err, map[string]interface{}{"user_id": opts.UserID})
		return fmt.Errorf("failed to update user settings: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.update-settings", 1)
	return formatOutput(settings, nil)
}

// newUsersPageViewsCmd returns the 'users page-views' command
func newUsersPageViewsCmd() *cobra.Command {
	opts := &options.UsersPageViewsOptions{}
	cmd := &cobra.Command{
		Use:   "page-views <user-id>",
		Short: "Get page views for a user",
		Long: `Get page view history for a user.

Examples:
  canvas users page-views 123`,
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
			return runUsersPageViews(cmd.Context(), client, opts)
		},
	}
	return cmd
}

func runUsersPageViews(ctx context.Context, client *api.Client, opts *options.UsersPageViewsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.page-views", map[string]interface{}{"user_id": opts.UserID})

	svc := api.NewUsersService(client)
	views, err := svc.GetPageViews(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "users.page-views", err, map[string]interface{}{"user_id": opts.UserID})
		return fmt.Errorf("failed to get page views: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.page-views", len(views))
	return formatEmptyOrOutput(views, "No page views found")
}

// newUsersLoginsCmd returns the 'users logins' command
func newUsersLoginsCmd() *cobra.Command {
	opts := &options.UsersLoginsOptions{}
	cmd := &cobra.Command{
		Use:   "logins <user-id>",
		Short: "List logins for a user",
		Long: `List login pseudonyms for a user.

Examples:
  canvas users logins 123`,
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
			return runUsersLogins(cmd.Context(), client, opts)
		},
	}
	return cmd
}

func runUsersLogins(ctx context.Context, client *api.Client, opts *options.UsersLoginsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.logins", map[string]interface{}{"user_id": opts.UserID})

	svc := api.NewUsersService(client)
	logins, err := svc.ListLogins(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "users.logins", err, map[string]interface{}{"user_id": opts.UserID})
		return fmt.Errorf("failed to list logins: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.logins", len(logins))
	return formatEmptyOrOutput(logins, "No logins found")
}

// newUsersCoursesCmd returns the 'users courses' command
func newUsersCoursesCmd() *cobra.Command {
	opts := &options.UsersCoursesOptions{}
	cmd := &cobra.Command{
		Use:   "courses <user-id>",
		Short: "List courses for a user",
		Long: `List courses the user is enrolled in.

Examples:
  canvas users courses 123`,
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
			return runUsersCourses(cmd.Context(), client, opts)
		},
	}
	return cmd
}

func runUsersCourses(ctx context.Context, client *api.Client, opts *options.UsersCoursesOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.courses", map[string]interface{}{"user_id": opts.UserID})

	svc := api.NewUsersService(client)
	courses, err := svc.ListUserCourses(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "users.courses", err, map[string]interface{}{"user_id": opts.UserID})
		return fmt.Errorf("failed to list user courses: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.courses", len(courses))
	return formatEmptyOrOutput(courses, "No courses found")
}

// newUsersMissingSubmissionsCmd returns the 'users missing-submissions' command
func newUsersMissingSubmissionsCmd() *cobra.Command {
	opts := &options.UsersMissingSubmissionsOptions{}
	cmd := &cobra.Command{
		Use:   "missing-submissions <user-id>",
		Short: "Get missing submissions for a user",
		Long: `Get assignments with missing submissions for a user.

Examples:
  canvas users missing-submissions 123`,
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
			return runUsersMissingSubmissions(cmd.Context(), client, opts)
		},
	}
	return cmd
}

func runUsersMissingSubmissions(ctx context.Context, client *api.Client, opts *options.UsersMissingSubmissionsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.missing-submissions", map[string]interface{}{"user_id": opts.UserID})

	svc := api.NewUsersService(client)
	subs, err := svc.GetMissingSubmissions(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "users.missing-submissions", err, map[string]interface{}{"user_id": opts.UserID})
		return fmt.Errorf("failed to get missing submissions: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.missing-submissions", len(subs))
	return formatEmptyOrOutput(subs, "No missing submissions found")
}

// newUsersActivityStreamCmd returns the 'users activity-stream' command
func newUsersActivityStreamCmd() *cobra.Command {
	opts := &options.UsersActivityStreamOptions{}
	cmd := &cobra.Command{
		Use:   "activity-stream",
		Short: "Get the activity stream for the current user",
		Long: `Get the activity stream for the currently authenticated user.

Examples:
  canvas users activity-stream`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUsersActivityStream(cmd.Context(), client, opts)
		},
	}
	return cmd
}

func runUsersActivityStream(ctx context.Context, client *api.Client, opts *options.UsersActivityStreamOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.activity-stream", map[string]interface{}{})

	svc := api.NewUsersService(client)
	items, err := svc.GetActivityStream(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "users.activity-stream", err, map[string]interface{}{})
		return fmt.Errorf("failed to get activity stream: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.activity-stream", len(items))
	return formatEmptyOrOutput(items, "No activity stream items found")
}

// newUsersTodoCmd returns the 'users todo' command
func newUsersTodoCmd() *cobra.Command {
	opts := &options.UsersTodoOptions{}
	cmd := &cobra.Command{
		Use:   "todo",
		Short: "Get todo items for the current user",
		Long: `Get todo items for the currently authenticated user.

Examples:
  canvas users todo`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUsersTodo(cmd.Context(), client, opts)
		},
	}
	return cmd
}

func runUsersTodo(ctx context.Context, client *api.Client, opts *options.UsersTodoOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.todo", map[string]interface{}{})

	svc := api.NewUsersService(client)
	items, err := svc.GetTodo(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "users.todo", err, map[string]interface{}{})
		return fmt.Errorf("failed to get todo items: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.todo", len(items))
	return formatEmptyOrOutput(items, "No todo items found")
}

// newUsersUpcomingEventsCmd returns the 'users upcoming-events' command
func newUsersUpcomingEventsCmd() *cobra.Command {
	opts := &options.UsersUpcomingEventsOptions{}
	cmd := &cobra.Command{
		Use:   "upcoming-events",
		Short: "Get upcoming events for the current user",
		Long: `Get upcoming events for the currently authenticated user.

Examples:
  canvas users upcoming-events`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUsersUpcomingEvents(cmd.Context(), client, opts)
		},
	}
	return cmd
}

func runUsersUpcomingEvents(ctx context.Context, client *api.Client, opts *options.UsersUpcomingEventsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.upcoming-events", map[string]interface{}{})

	svc := api.NewUsersService(client)
	events, err := svc.GetUpcomingEvents(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "users.upcoming-events", err, map[string]interface{}{})
		return fmt.Errorf("failed to get upcoming events: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.upcoming-events", len(events))
	return formatEmptyOrOutput(events, "No upcoming events found")
}

// newUsersMergeCmd returns the 'users merge' command
func newUsersMergeCmd() *cobra.Command {
	opts := &options.UsersMergeOptions{}
	cmd := &cobra.Command{
		Use:   "merge <user-id> <destination-user-id>",
		Short: "Merge a user into another user",
		Long: `Merge a user into a destination user. Requires admin privileges.

Examples:
  canvas users merge 123 456`,
		Args: ExactArgsWithUsage(2, "user-id", "destination-user-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			userID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid user ID: %s", args[0])
			}
			destUserID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid destination user ID: %s", args[1])
			}
			opts.UserID = userID
			opts.DestinationUserID = destUserID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUsersMerge(cmd.Context(), client, opts)
		},
	}
	return cmd
}

func runUsersMerge(ctx context.Context, client *api.Client, opts *options.UsersMergeOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.merge", map[string]interface{}{
		"user_id":             opts.UserID,
		"destination_user_id": opts.DestinationUserID,
	})

	svc := api.NewUsersService(client)
	user, err := svc.MergeInto(ctx, opts.UserID, opts.DestinationUserID)
	if err != nil {
		logger.LogCommandError(ctx, "users.merge", err, map[string]interface{}{
			"user_id":             opts.UserID,
			"destination_user_id": opts.DestinationUserID,
		})
		return fmt.Errorf("failed to merge user: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.merge", 1)
	printInfo("User merged successfully!\n")
	return formatOutput(user, nil)
}

// newUsersSplitCmd returns the 'users split' command
func newUsersSplitCmd() *cobra.Command {
	opts := &options.UsersSplitOptions{}
	cmd := &cobra.Command{
		Use:   "split <user-id>",
		Short: "Split a merged user back into separate users",
		Long: `Split a previously merged user back into separate users. Requires admin privileges.

Examples:
  canvas users split 123`,
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
			return runUsersSplit(cmd.Context(), client, opts)
		},
	}
	return cmd
}

func runUsersSplit(ctx context.Context, client *api.Client, opts *options.UsersSplitOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.split", map[string]interface{}{"user_id": opts.UserID})

	svc := api.NewUsersService(client)
	users, err := svc.Split(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "users.split", err, map[string]interface{}{"user_id": opts.UserID})
		return fmt.Errorf("failed to split user: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.split", len(users))
	printInfo("User split successfully into %d users!\n", len(users))
	return formatEmptyOrOutput(users, "No users returned")
}
