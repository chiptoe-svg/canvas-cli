package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/commands/internal/options"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

// contentSharesCmd represents the content-shares command group
var contentSharesCmd = &cobra.Command{
	Use:   "content-shares",
	Short: "Manage content shares",
	Long: `Manage content shares between Canvas users.

Examples:
  canvas content-shares list-sent --user-id 123
  canvas content-shares list-received --user-id 123
  canvas content-shares get --user-id 123 --id 456
  canvas content-shares delete --user-id 123 --id 456`,
}

func init() {
	rootCmd.AddCommand(contentSharesCmd)
	contentSharesCmd.AddCommand(newContentSharesListSentCmd())
	contentSharesCmd.AddCommand(newContentSharesListReceivedCmd())
	contentSharesCmd.AddCommand(newContentSharesGetCmd())
	contentSharesCmd.AddCommand(newContentSharesDeleteCmd())
}

func newContentSharesListSentCmd() *cobra.Command {
	opts := &options.ContentSharesListSentOptions{}

	cmd := &cobra.Command{
		Use:   "list-sent",
		Short: "List sent content shares for a user",
		Long: `List all content shares sent by a specific user.

Examples:
  canvas content-shares list-sent --user-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runContentSharesListSent(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	mustMarkRequired(cmd, "user-id")

	return cmd
}

func newContentSharesListReceivedCmd() *cobra.Command {
	opts := &options.ContentSharesListReceivedOptions{}

	cmd := &cobra.Command{
		Use:   "list-received",
		Short: "List received content shares for a user",
		Long: `List all content shares received by a specific user.

Examples:
  canvas content-shares list-received --user-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runContentSharesListReceived(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	mustMarkRequired(cmd, "user-id")

	return cmd
}

func newContentSharesGetCmd() *cobra.Command {
	opts := &options.ContentSharesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a specific content share",
		Long: `Get details of a specific content share.

Examples:
  canvas content-shares get --user-id 123 --id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runContentSharesGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.Flags().Int64Var(&opts.ID, "id", 0, "Content share ID (required)")
	mustMarkRequired(cmd, "user-id", "id")

	return cmd
}

func newContentSharesDeleteCmd() *cobra.Command {
	opts := &options.ContentSharesDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a content share",
		Long: `Delete a content share for a user.

Examples:
  canvas content-shares delete --user-id 123 --id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runContentSharesDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.Flags().Int64Var(&opts.ID, "id", 0, "Content share ID (required)")
	mustMarkRequired(cmd, "user-id", "id")

	return cmd
}

func runContentSharesListSent(ctx context.Context, client *api.Client, opts *options.ContentSharesListSentOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "content-shares.list-sent", map[string]interface{}{
		"user_id": opts.UserID,
	})

	svc := api.NewContentSharesService(client)

	shares, err := svc.ListSent(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "content-shares.list-sent", err, map[string]interface{}{
			"user_id": opts.UserID,
		})
		return fmt.Errorf("failed to list sent content shares: %w", err)
	}

	logger.LogCommandComplete(ctx, "content-shares.list-sent", len(shares))
	return formatEmptyOrOutput(shares, "No sent content shares found")
}

func runContentSharesListReceived(ctx context.Context, client *api.Client, opts *options.ContentSharesListReceivedOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "content-shares.list-received", map[string]interface{}{
		"user_id": opts.UserID,
	})

	svc := api.NewContentSharesService(client)

	shares, err := svc.ListReceived(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "content-shares.list-received", err, map[string]interface{}{
			"user_id": opts.UserID,
		})
		return fmt.Errorf("failed to list received content shares: %w", err)
	}

	logger.LogCommandComplete(ctx, "content-shares.list-received", len(shares))
	return formatEmptyOrOutput(shares, "No received content shares found")
}

func runContentSharesGet(ctx context.Context, client *api.Client, opts *options.ContentSharesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "content-shares.get", map[string]interface{}{
		"user_id": opts.UserID,
		"id":      opts.ID,
	})

	svc := api.NewContentSharesService(client)

	share, err := svc.Get(ctx, opts.UserID, opts.ID)
	if err != nil {
		logger.LogCommandError(ctx, "content-shares.get", err, map[string]interface{}{
			"user_id": opts.UserID,
			"id":      opts.ID,
		})
		return fmt.Errorf("failed to get content share: %w", err)
	}

	logger.LogCommandComplete(ctx, "content-shares.get", 1)
	return formatOutput(share, nil)
}

func runContentSharesDelete(ctx context.Context, client *api.Client, opts *options.ContentSharesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "content-shares.delete", map[string]interface{}{
		"user_id": opts.UserID,
		"id":      opts.ID,
	})

	svc := api.NewContentSharesService(client)

	if err := svc.Delete(ctx, opts.UserID, opts.ID); err != nil {
		logger.LogCommandError(ctx, "content-shares.delete", err, map[string]interface{}{
			"user_id": opts.UserID,
			"id":      opts.ID,
		})
		return fmt.Errorf("failed to delete content share: %w", err)
	}

	logger.LogCommandComplete(ctx, "content-shares.delete", 1)
	printInfo("Content share %d deleted successfully\n", opts.ID)
	return nil
}
