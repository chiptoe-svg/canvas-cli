package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// commMessagesCmd represents the comm-messages command group.
var commMessagesCmd = &cobra.Command{
	Use:   "comm-messages",
	Short: "List Canvas communication messages",
	Long: `List communication messages that have been sent to a user.

Examples:
  canvas comm-messages list
  canvas comm-messages list --user-id 42`,
}

func init() {
	rootCmd.AddCommand(commMessagesCmd)
	commMessagesCmd.AddCommand(newCommMessagesListCmd())
}

func newCommMessagesListCmd() *cobra.Command {
	opts := &options.CommMessagesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List communication messages",
		Long: `List communication messages that Canvas has sent to a user.

Examples:
  canvas comm-messages list
  canvas comm-messages list --user-id 42`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCommMessagesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "Filter messages by user ID")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Number of results per page")

	return cmd
}

func runCommMessagesList(ctx context.Context, client *api.Client, opts *options.CommMessagesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "comm-messages.list", map[string]interface{}{
		"user_id": opts.UserID,
	})

	svc := api.NewCommMessagesService(client)

	msgs, err := svc.List(ctx, &api.ListCommMessagesOptions{
		UserID:  opts.UserID,
		PerPage: opts.PerPage,
	})
	if err != nil {
		logger.LogCommandError(ctx, "comm-messages.list", err, nil)
		return fmt.Errorf("failed to list comm messages: %w", err)
	}

	printVerbose("Found %d communication messages:\n\n", len(msgs))
	logger.LogCommandComplete(ctx, "comm-messages.list", len(msgs))
	return formatEmptyOrOutput(msgs, "No communication messages found")
}
