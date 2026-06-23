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

// epubExportsCmd represents the epub-exports command group.
var epubExportsCmd = &cobra.Command{
	Use:   "epub-exports",
	Short: "Manage Canvas ePub exports",
	Long: `Create and check the status of ePub course exports.

Examples:
  canvas epub-exports list
  canvas epub-exports create --course-id 123
  canvas epub-exports get --course-id 123 5`,
}

func init() {
	rootCmd.AddCommand(epubExportsCmd)
	epubExportsCmd.AddCommand(newEpubExportsListCmd())
	epubExportsCmd.AddCommand(newEpubExportCreateCmd())
	epubExportsCmd.AddCommand(newEpubExportGetCmd())
}

func newEpubExportsListCmd() *cobra.Command {
	opts := &options.EpubExportsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List ePub exports",
		Long: `List all ePub exports visible to the current user.

Examples:
  canvas epub-exports list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runEpubExportsList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newEpubExportCreateCmd() *cobra.Command {
	opts := &options.EpubExportCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an ePub export for a course",
		Long: `Start a new ePub export job for a course.

Examples:
  canvas epub-exports create --course-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runEpubExportCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")

	return cmd
}

func newEpubExportGetCmd() *cobra.Command {
	opts := &options.EpubExportGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get an ePub export by ID",
		Long: `Retrieve the status and details of an ePub export.

Examples:
  canvas epub-exports get --course-id 123 5`,
		Args: ExactArgsWithUsage(1, "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid export ID: %s", args[0])
			}
			opts.ID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runEpubExportGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")

	return cmd
}

func runEpubExportsList(ctx context.Context, client *api.Client, opts *options.EpubExportsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "epub-exports.list", map[string]interface{}{})

	svc := api.NewEpubExportsService(client)

	exports, err := svc.List(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "epub-exports.list", err, nil)
		return fmt.Errorf("failed to list ePub exports: %w", err)
	}

	printVerbose("Found %d ePub exports:\n\n", len(exports))
	logger.LogCommandComplete(ctx, "epub-exports.list", len(exports))
	return formatEmptyOrOutput(exports, "No ePub exports found")
}

func runEpubExportCreate(ctx context.Context, client *api.Client, opts *options.EpubExportCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "epub-exports.create", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	svc := api.NewEpubExportsService(client)

	export, err := svc.Create(ctx, opts.CourseID)
	if err != nil {
		logger.LogCommandError(ctx, "epub-exports.create", err, nil)
		return fmt.Errorf("failed to create ePub export: %w", err)
	}

	logger.LogCommandComplete(ctx, "epub-exports.create", 1)
	return formatSuccessOutput(export, "ePub export started successfully.")
}

func runEpubExportGet(ctx context.Context, client *api.Client, opts *options.EpubExportGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "epub-exports.get", map[string]interface{}{
		"course_id": opts.CourseID,
		"id":        opts.ID,
	})

	svc := api.NewEpubExportsService(client)

	export, err := svc.Get(ctx, opts.CourseID, opts.ID)
	if err != nil {
		logger.LogCommandError(ctx, "epub-exports.get", err, nil)
		return fmt.Errorf("failed to get ePub export: %w", err)
	}

	logger.LogCommandComplete(ctx, "epub-exports.get", 1)
	return formatOutput(export, nil)
}
