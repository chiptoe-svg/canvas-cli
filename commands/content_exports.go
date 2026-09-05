package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/commands/internal/options"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

// contentExportsCmd is the parent command for content export operations.
var contentExportsCmd = &cobra.Command{
	Use:     "content-exports",
	Aliases: []string{"exports", "ce"},
	Short:   "Manage course content exports",
	Long: `Manage Canvas content exports and epub exports for courses.

Examples:
  canvas content-exports list --course-id 1
  canvas content-exports get 5 --course-id 1
  canvas content-exports create --course-id 1 --export-type common_cartridge
  canvas content-exports epub-create --course-id 1
  canvas content-exports epub-get 3 --course-id 1`,
}

func init() {
	rootCmd.AddCommand(contentExportsCmd)
	contentExportsCmd.AddCommand(newContentExportsListCmd())
	contentExportsCmd.AddCommand(newContentExportsGetCmd())
	contentExportsCmd.AddCommand(newContentExportsCreateCmd())
	contentExportsCmd.AddCommand(newEpubExportsCreateCmd())
	contentExportsCmd.AddCommand(newEpubExportsGetCmd())
}

func newContentExportsListCmd() *cobra.Command {
	opts := &options.ContentExportsListOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List content exports for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runContentExportsList(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runContentExportsList(ctx context.Context, client *api.Client, opts *options.ContentExportsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "content-exports.list", map[string]interface{}{"course_id": opts.CourseID})

	svc := api.NewContentExportsService(client)
	exports, err := svc.List(ctx, opts.CourseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "content-exports.list", len(exports))
	return formatEmptyOrOutput(exports, "No content exports found")
}

func newContentExportsGetCmd() *cobra.Command {
	opts := &options.ContentExportsGetOptions{}
	cmd := &cobra.Command{
		Use:   "get <export-id>",
		Short: "Get a content export",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &opts.ID); err != nil {
				return fmt.Errorf("invalid export ID: %s", args[0])
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runContentExportsGet(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runContentExportsGet(ctx context.Context, client *api.Client, opts *options.ContentExportsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "content-exports.get", map[string]interface{}{"course_id": opts.CourseID, "id": opts.ID})

	svc := api.NewContentExportsService(client)
	export, err := svc.Get(ctx, opts.CourseID, opts.ID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "content-exports.get", 1)
	return formatOutput(export, nil)
}

func newContentExportsCreateCmd() *cobra.Command {
	opts := &options.ContentExportsCreateOptions{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a content export",
		Long: `Start a content export for a course.

Export types: common_cartridge, zip, qti, course_copy

Examples:
  canvas content-exports create --course-id 1 --export-type common_cartridge`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runContentExportsCreate(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.ExportType, "export-type", "", "Export type: common_cartridge, zip, qti, course_copy (required)")
	cmd.Flags().BoolVar(&opts.SkipNotifications, "skip-notifications", false, "Skip email notifications")
	mustMarkRequired(cmd, "course-id", "export-type")
	return cmd
}

func runContentExportsCreate(ctx context.Context, client *api.Client, opts *options.ContentExportsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "content-exports.create", map[string]interface{}{"course_id": opts.CourseID, "type": opts.ExportType})

	svc := api.NewContentExportsService(client)
	export, err := svc.Create(ctx, opts.CourseID, api.CreateContentExportParams{
		ExportType:        opts.ExportType,
		SkipNotifications: opts.SkipNotifications,
	})
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "content-exports.create", 1)
	return formatSuccessOutput(export, fmt.Sprintf("Content export %d created (state: %s)", export.ID, export.WorkflowState))
}

func newEpubExportsCreateCmd() *cobra.Command {
	var courseID int64
	cmd := &cobra.Command{
		Use:   "epub-create",
		Short: "Create an epub export for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runEpubExportsCreate(cmd.Context(), client, courseID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runEpubExportsCreate(ctx context.Context, client *api.Client, courseID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "epub-exports.create", map[string]interface{}{"course_id": courseID})

	svc := api.NewContentExportsService(client)
	export, err := svc.CreateEpub(ctx, courseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "epub-exports.create", 1)
	return formatSuccessOutput(export, fmt.Sprintf("Epub export %d created (state: %s)", export.ID, export.WorkflowState))
}

func newEpubExportsGetCmd() *cobra.Command {
	opts := &options.EpubExportsGetOptions{}
	cmd := &cobra.Command{
		Use:   "epub-get <epub-id>",
		Short: "Get an epub export",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &opts.ID); err != nil {
				return fmt.Errorf("invalid epub ID: %s", args[0])
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runEpubExportsGet(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runEpubExportsGet(ctx context.Context, client *api.Client, opts *options.EpubExportsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "epub-exports.get", map[string]interface{}{"course_id": opts.CourseID, "id": opts.ID})

	svc := api.NewContentExportsService(client)
	export, err := svc.GetEpub(ctx, opts.CourseID, opts.ID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "epub-exports.get", 1)
	return formatOutput(export, nil)
}
