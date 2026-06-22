package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// historyCmd represents the history command group.
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View Canvas user page-view history",
	Long: `View the page-view history for a Canvas user.

Examples:
  canvas history list --user-id 42`,
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.AddCommand(newHistoryListCmd())
}

func newHistoryListCmd() *cobra.Command {
	opts := &options.HistoryListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List page-view history for a user",
		Long: `List the page-view history entries for a Canvas user.

Examples:
  canvas history list --user-id 42`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runHistoryList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	mustMarkRequired(cmd, "user-id")

	return cmd
}

func runHistoryList(ctx context.Context, client *api.Client, opts *options.HistoryListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "history.list", map[string]interface{}{
		"user_id": opts.UserID,
	})

	svc := api.NewHistoryService(client)

	entries, err := svc.List(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "history.list", err, nil)
		return fmt.Errorf("failed to list history: %w", err)
	}

	printVerbose("Found %d history entries:\n\n", len(entries))
	logger.LogCommandComplete(ctx, "history.list", len(entries))
	return formatEmptyOrOutput(entries, "No history entries found")
}
