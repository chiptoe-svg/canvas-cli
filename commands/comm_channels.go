package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// commChannelsCmd represents the comm-channels command group
var commChannelsCmd = &cobra.Command{
	Use:   "comm-channels",
	Short: "Manage communication channels",
	Long: `Manage communication channels for Canvas users.

Examples:
  canvas comm-channels list --user-id 123
  canvas comm-channels create --user-id 123 --address "user@example.com" --type email
  canvas comm-channels delete --user-id 123 --channel-id 456`,
}

func init() {
	rootCmd.AddCommand(commChannelsCmd)
	commChannelsCmd.AddCommand(newCommChannelsListCmd())
	commChannelsCmd.AddCommand(newCommChannelsCreateCmd())
	commChannelsCmd.AddCommand(newCommChannelsDeleteCmd())
}

func newCommChannelsListCmd() *cobra.Command {
	opts := &options.CommChannelsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List communication channels for a user",
		Long: `List all communication channels for a specific user.

Examples:
  canvas comm-channels list --user-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCommChannelsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	mustMarkRequired(cmd, "user-id")

	return cmd
}

func newCommChannelsCreateCmd() *cobra.Command {
	opts := &options.CommChannelsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a communication channel for a user",
		Long: `Create a new communication channel for a user.

Channel types: email, sms, push, twitter, yo

Examples:
  canvas comm-channels create --user-id 123 --address "user@example.com" --type email
  canvas comm-channels create --user-id 123 --address "+15551234567" --type sms --skip-confirmation`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCommChannelsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.Flags().StringVar(&opts.Address, "address", "", "Channel address, e.g. email address or phone number (required)")
	cmd.Flags().StringVar(&opts.Type, "type", "", "Channel type: email, sms, push, twitter, yo (required)")
	cmd.Flags().BoolVar(&opts.SkipConfirmation, "skip-confirmation", false, "Skip confirmation (admin only)")
	mustMarkRequired(cmd, "user-id", "address", "type")

	return cmd
}

func newCommChannelsDeleteCmd() *cobra.Command {
	opts := &options.CommChannelsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a communication channel",
		Long: `Delete a communication channel for a user.

Examples:
  canvas comm-channels delete --user-id 123 --channel-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCommChannelsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.Flags().Int64Var(&opts.ChannelID, "channel-id", 0, "Channel ID (required)")
	mustMarkRequired(cmd, "user-id", "channel-id")

	return cmd
}

func runCommChannelsList(ctx context.Context, client *api.Client, opts *options.CommChannelsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "comm-channels.list", map[string]interface{}{
		"user_id": opts.UserID,
	})

	svc := api.NewCommunicationChannelsService(client)

	channels, err := svc.List(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "comm-channels.list", err, map[string]interface{}{
			"user_id": opts.UserID,
		})
		return fmt.Errorf("failed to list communication channels: %w", err)
	}

	logger.LogCommandComplete(ctx, "comm-channels.list", len(channels))
	return formatEmptyOrOutput(channels, "No communication channels found")
}

func runCommChannelsCreate(ctx context.Context, client *api.Client, opts *options.CommChannelsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "comm-channels.create", map[string]interface{}{
		"user_id": opts.UserID,
		"address": opts.Address,
		"type":    opts.Type,
	})

	svc := api.NewCommunicationChannelsService(client)

	params := api.CreateCommunicationChannelParams{
		Address:          opts.Address,
		Type:             opts.Type,
		SkipConfirmation: opts.SkipConfirmation,
	}

	channel, err := svc.Create(ctx, opts.UserID, params)
	if err != nil {
		logger.LogCommandError(ctx, "comm-channels.create", err, map[string]interface{}{
			"user_id": opts.UserID,
			"address": opts.Address,
		})
		return fmt.Errorf("failed to create communication channel: %w", err)
	}

	logger.LogCommandComplete(ctx, "comm-channels.create", 1)
	return formatSuccessOutput(channel, "Communication channel created successfully!")
}

func runCommChannelsDelete(ctx context.Context, client *api.Client, opts *options.CommChannelsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "comm-channels.delete", map[string]interface{}{
		"user_id":    opts.UserID,
		"channel_id": opts.ChannelID,
	})

	svc := api.NewCommunicationChannelsService(client)

	_, err := svc.Delete(ctx, opts.UserID, opts.ChannelID)
	if err != nil {
		logger.LogCommandError(ctx, "comm-channels.delete", err, map[string]interface{}{
			"user_id":    opts.UserID,
			"channel_id": opts.ChannelID,
		})
		return fmt.Errorf("failed to delete communication channel: %w", err)
	}

	logger.LogCommandComplete(ctx, "comm-channels.delete", 1)
	printInfo("Communication channel %d deleted successfully\n", opts.ChannelID)
	return nil
}
