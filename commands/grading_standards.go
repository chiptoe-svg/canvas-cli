package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// gradingStandardsCmd is the parent command for grading standard operations.
var gradingStandardsCmd = &cobra.Command{
	Use:     "grading-standards",
	Aliases: []string{"gs"},
	Short:   "Manage course grading standards",
	Long: `Manage Canvas course grading standards (grading schemes).

Examples:
  canvas grading-standards list --course-id 1
  canvas grading-standards get 5 --course-id 1
  canvas grading-standards create --course-id 1 --title "Letter Grades"
  canvas grading-standards delete 5 --course-id 1`,
}

func init() {
	rootCmd.AddCommand(gradingStandardsCmd)
	gradingStandardsCmd.AddCommand(newGradingStandardsListCmd())
	gradingStandardsCmd.AddCommand(newGradingStandardsGetCmd())
	gradingStandardsCmd.AddCommand(newGradingStandardsCreateCmd())
	gradingStandardsCmd.AddCommand(newGradingStandardsDeleteCmd())
}

func newGradingStandardsListCmd() *cobra.Command {
	opts := &options.GradingStandardsListOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List grading standards for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGradingStandardsList(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runGradingStandardsList(ctx context.Context, client *api.Client, opts *options.GradingStandardsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-standards.list", map[string]interface{}{"course_id": opts.CourseID})

	svc := api.NewGradingStandardsService(client)
	standards, err := svc.ListForCourse(ctx, opts.CourseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "grading-standards.list", len(standards))
	return formatEmptyOrOutput(standards, "No grading standards found")
}

func newGradingStandardsGetCmd() *cobra.Command {
	opts := &options.GradingStandardsGetOptions{}
	cmd := &cobra.Command{
		Use:   "get <standard-id>",
		Short: "Get a grading standard",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &opts.StandardID); err != nil {
				return fmt.Errorf("invalid standard ID: %s", args[0])
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGradingStandardsGet(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runGradingStandardsGet(ctx context.Context, client *api.Client, opts *options.GradingStandardsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-standards.get", map[string]interface{}{"course_id": opts.CourseID, "standard_id": opts.StandardID})

	svc := api.NewGradingStandardsService(client)
	standard, err := svc.GetForCourse(ctx, opts.CourseID, opts.StandardID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "grading-standards.get", 1)
	return formatOutput(standard, nil)
}

func newGradingStandardsCreateCmd() *cobra.Command {
	opts := &options.GradingStandardsCreateOptions{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a grading standard for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGradingStandardsCreate(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Grading standard title (required)")
	mustMarkRequired(cmd, "course-id", "title")
	return cmd
}

func runGradingStandardsCreate(ctx context.Context, client *api.Client, opts *options.GradingStandardsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-standards.create", map[string]interface{}{"course_id": opts.CourseID})

	svc := api.NewGradingStandardsService(client)
	standard, err := svc.CreateForCourse(ctx, opts.CourseID, api.GradingStandardParams{Title: opts.Title})
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "grading-standards.create", 1)
	return formatSuccessOutput(standard, fmt.Sprintf("Grading standard %d created", standard.ID))
}

func newGradingStandardsDeleteCmd() *cobra.Command {
	opts := &options.GradingStandardsDeleteOptions{}
	cmd := &cobra.Command{
		Use:   "delete <standard-id>",
		Short: "Delete a grading standard",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &opts.StandardID); err != nil {
				return fmt.Errorf("invalid standard ID: %s", args[0])
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGradingStandardsDelete(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runGradingStandardsDelete(ctx context.Context, client *api.Client, opts *options.GradingStandardsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-standards.delete", map[string]interface{}{"course_id": opts.CourseID, "standard_id": opts.StandardID})

	svc := api.NewGradingStandardsService(client)
	if err := svc.DeleteForCourse(ctx, opts.CourseID, opts.StandardID); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "grading-standards.delete", 1)
	printInfo("Grading standard %d deleted\n", opts.StandardID)
	return nil
}
