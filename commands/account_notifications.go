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

// accountNotificationsCmd is the root command for account notifications
var accountNotificationsCmd = &cobra.Command{
	Use:   "account-notifications",
	Short: "Manage Canvas account notifications",
	Long: `Manage account-wide notifications in Canvas.

Account notifications are announcements displayed to users when they log in.
They can be targeted to specific user roles and scheduled for specific time periods.

Examples:
  canvas account-notifications list 1
  canvas account-notifications get 1 5
  canvas account-notifications create 1 --subject "Maintenance" --message "System down" --start-at 2026-01-01 --end-at 2026-01-02
  canvas account-notifications delete 1 5`,
}

func init() {
	rootCmd.AddCommand(accountNotificationsCmd)
	accountNotificationsCmd.AddCommand(newAccountNotificationsListCmd())
	accountNotificationsCmd.AddCommand(newAccountNotificationsGetCmd())
	accountNotificationsCmd.AddCommand(newAccountNotificationsCreateCmd())
	accountNotificationsCmd.AddCommand(newAccountNotificationsUpdateCmd())
	accountNotificationsCmd.AddCommand(newAccountNotificationsDeleteCmd())
}

func newAccountNotificationsListCmd() *cobra.Command {
	opts := &options.AccountNotificationsListOptions{}

	cmd := &cobra.Command{
		Use:   "list <account-id>",
		Short: "List account notifications",
		Args:  ExactArgsWithUsage(1, "account-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountNotificationsList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountNotificationsGetCmd() *cobra.Command {
	opts := &options.AccountNotificationsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <account-id> <id>",
		Short: "Get a specific account notification",
		Args:  ExactArgsWithUsage(2, "account-id", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			notificationID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid notification ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.NotificationID = notificationID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountNotificationsGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountNotificationsCreateCmd() *cobra.Command {
	opts := &options.AccountNotificationsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <account-id>",
		Short: "Create a new account notification",
		Args:  ExactArgsWithUsage(1, "account-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountNotificationsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Subject, "subject", "", "Notification subject")
	cmd.Flags().StringVar(&opts.Message, "message", "", "Notification message")
	cmd.Flags().StringVar(&opts.StartAt, "start-at", "", "Start date/time (ISO 8601)")
	cmd.Flags().StringVar(&opts.EndAt, "end-at", "", "End date/time (ISO 8601)")
	cmd.Flags().StringVar(&opts.Icon, "icon", "", "Icon type (warning, information, question, error, calendar)")
	cmd.Flags().StringSliceVar(&opts.Roles, "roles", []string{}, "Target roles")
	mustMarkRequired(cmd, "subject", "message", "start-at", "end-at")

	return cmd
}

func newAccountNotificationsUpdateCmd() *cobra.Command {
	opts := &options.AccountNotificationsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <account-id> <id>",
		Short: "Update an account notification",
		Args:  ExactArgsWithUsage(2, "account-id", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			notificationID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid notification ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.NotificationID = notificationID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountNotificationsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Subject, "subject", "", "Notification subject")
	cmd.Flags().StringVar(&opts.Message, "message", "", "Notification message")
	cmd.Flags().StringVar(&opts.StartAt, "start-at", "", "Start date/time (ISO 8601)")
	cmd.Flags().StringVar(&opts.EndAt, "end-at", "", "End date/time (ISO 8601)")
	cmd.Flags().StringVar(&opts.Icon, "icon", "", "Icon type")
	cmd.Flags().StringSliceVar(&opts.Roles, "roles", []string{}, "Target roles")

	return cmd
}

func newAccountNotificationsDeleteCmd() *cobra.Command {
	opts := &options.AccountNotificationsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <account-id> <id>",
		Short: "Delete an account notification",
		Args:  ExactArgsWithUsage(2, "account-id", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			notificationID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid notification ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.NotificationID = notificationID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountNotificationsDelete(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func runAccountNotificationsList(ctx context.Context, client *api.Client, opts *options.AccountNotificationsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-notifications.list", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAccountNotificationsService(client)
	notifications, err := svc.List(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "account-notifications.list", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "account-notifications.list", len(notifications))

	return formatEmptyOrOutput(notifications, fmt.Sprintf("No notifications found for account %d", opts.AccountID))
}

func runAccountNotificationsGet(ctx context.Context, client *api.Client, opts *options.AccountNotificationsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-notifications.get", map[string]interface{}{
		"account_id":      opts.AccountID,
		"notification_id": opts.NotificationID,
	})

	svc := api.NewAccountNotificationsService(client)
	notification, err := svc.Get(ctx, opts.AccountID, opts.NotificationID)
	if err != nil {
		logger.LogCommandError(ctx, "account-notifications.get", err, map[string]interface{}{
			"account_id":      opts.AccountID,
			"notification_id": opts.NotificationID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "account-notifications.get", 1)

	return formatOutput(notification, func() {
		fmt.Printf("Account Notification\n")
		fmt.Printf("====================\n\n")
		fmt.Printf("ID:      %d\n", notification.ID)
		fmt.Printf("Subject: %s\n", notification.Subject)
		fmt.Printf("Message: %s\n", notification.Message)
		fmt.Printf("Start:   %s\n", notification.StartAt)
		fmt.Printf("End:     %s\n", notification.EndAt)
		if notification.Icon != "" {
			fmt.Printf("Icon:    %s\n", notification.Icon)
		}
		if len(notification.Roles) > 0 {
			fmt.Printf("Roles:   %v\n", notification.Roles)
		}
	})
}

func runAccountNotificationsCreate(ctx context.Context, client *api.Client, opts *options.AccountNotificationsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-notifications.create", map[string]interface{}{
		"account_id": opts.AccountID,
		"subject":    opts.Subject,
	})

	svc := api.NewAccountNotificationsService(client)
	params := &api.AccountNotificationParams{
		AccountNotification: api.AccountNotificationFields{
			Subject: opts.Subject,
			Message: opts.Message,
			StartAt: opts.StartAt,
			EndAt:   opts.EndAt,
			Icon:    opts.Icon,
		},
		Roles: opts.Roles,
	}

	notification, err := svc.Create(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "account-notifications.create", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "account-notifications.create", 1)

	return formatSuccessOutput(notification, fmt.Sprintf("Notification created (ID: %d)", notification.ID))
}

func runAccountNotificationsUpdate(ctx context.Context, client *api.Client, opts *options.AccountNotificationsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-notifications.update", map[string]interface{}{
		"account_id":      opts.AccountID,
		"notification_id": opts.NotificationID,
	})

	svc := api.NewAccountNotificationsService(client)
	params := &api.AccountNotificationParams{
		AccountNotification: api.AccountNotificationFields{
			Subject: opts.Subject,
			Message: opts.Message,
			StartAt: opts.StartAt,
			EndAt:   opts.EndAt,
			Icon:    opts.Icon,
		},
		Roles: opts.Roles,
	}

	notification, err := svc.Update(ctx, opts.AccountID, opts.NotificationID, params)
	if err != nil {
		logger.LogCommandError(ctx, "account-notifications.update", err, map[string]interface{}{
			"account_id":      opts.AccountID,
			"notification_id": opts.NotificationID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "account-notifications.update", 1)

	return formatSuccessOutput(notification, fmt.Sprintf("Notification %d updated", notification.ID))
}

func runAccountNotificationsDelete(ctx context.Context, client *api.Client, opts *options.AccountNotificationsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-notifications.delete", map[string]interface{}{
		"account_id":      opts.AccountID,
		"notification_id": opts.NotificationID,
	})

	svc := api.NewAccountNotificationsService(client)
	err := svc.Delete(ctx, opts.AccountID, opts.NotificationID)
	if err != nil {
		logger.LogCommandError(ctx, "account-notifications.delete", err, map[string]interface{}{
			"account_id":      opts.AccountID,
			"notification_id": opts.NotificationID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "account-notifications.delete", 1)

	printInfo("Notification %d deleted\n", opts.NotificationID)
	return nil
}
