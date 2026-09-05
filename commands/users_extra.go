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

func init() {
	usersCmd.AddCommand(newUsersProfileCmd())
	usersCmd.AddCommand(newUsersMissingSubmissionsCmd())
	usersCmd.AddCommand(newUsersActivityStreamCmd())
	usersCmd.AddCommand(newUsersTodoCmd())
	usersCmd.AddCommand(newUsersUpcomingEventsCmd())
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
