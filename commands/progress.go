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

// progressCmd represents the progress command group.
var progressCmd = &cobra.Command{
	Use:   "progress",
	Short: "Manage Canvas background job progress",
	Long: `Query or cancel Canvas background jobs.

Progress objects track long-running operations such as course imports,
blueprint migrations, and similar asynchronous tasks.

Examples:
  canvas progress get 42
  canvas progress cancel 42`,
}

func init() {
	rootCmd.AddCommand(progressCmd)
	progressCmd.AddCommand(newProgressGetCmd())
	progressCmd.AddCommand(newProgressCancelCmd())
}

func newProgressGetCmd() *cobra.Command {
	opts := &options.ProgressGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a progress job by ID",
		Long: `Retrieve the status of a Canvas background job by its progress ID.

Examples:
  canvas progress get 42`,
		Args: ExactArgsWithUsage(1, "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid progress ID: %s", args[0])
			}
			opts.ID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runProgressGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newProgressCancelCmd() *cobra.Command {
	opts := &options.ProgressCancelOptions{}

	cmd := &cobra.Command{
		Use:   "cancel <id>",
		Short: "Cancel a background job",
		Long: `Cancel an in-progress Canvas background job.

Examples:
  canvas progress cancel 42`,
		Args: ExactArgsWithUsage(1, "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid progress ID: %s", args[0])
			}
			opts.ID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runProgressCancel(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func runProgressGet(ctx context.Context, client *api.Client, opts *options.ProgressGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "progress.get", map[string]interface{}{
		"id": opts.ID,
	})

	svc := api.NewProgressService(client)

	p, err := svc.Get(ctx, opts.ID)
	if err != nil {
		logger.LogCommandError(ctx, "progress.get", err, map[string]interface{}{"id": opts.ID})
		return fmt.Errorf("failed to get progress job: %w", err)
	}

	logger.LogCommandComplete(ctx, "progress.get", 1)
	return formatOutput(p, nil)
}

func runProgressCancel(ctx context.Context, client *api.Client, opts *options.ProgressCancelOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "progress.cancel", map[string]interface{}{
		"id": opts.ID,
	})

	confirmed, err := confirmDelete("progress job", opts.ID, opts.Force)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Cancel aborted")
		return nil
	}

	svc := api.NewProgressService(client)

	p, err := svc.Cancel(ctx, opts.ID)
	if err != nil {
		logger.LogCommandError(ctx, "progress.cancel", err, map[string]interface{}{"id": opts.ID})
		return fmt.Errorf("failed to cancel progress job: %w", err)
	}

	logger.LogCommandComplete(ctx, "progress.cancel", 1)
	return formatSuccessOutput(p, "Progress job cancelled.")
}
