package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// gradingPeriodsCmd is the parent command for grading period operations.
var gradingPeriodsCmd = &cobra.Command{
	Use:     "grading-periods",
	Aliases: []string{"gp"},
	Short:   "Manage course grading periods",
	Long: `Manage Canvas course grading periods.

Examples:
  canvas grading-periods list --course-id 1
  canvas grading-periods get 5 --course-id 1
  canvas grading-periods update 5 --course-id 1 --title "Q1"
  canvas grading-periods delete 5 --course-id 1
  canvas grading-periods batch-update --course-id 1`,
}

func init() {
	rootCmd.AddCommand(gradingPeriodsCmd)
	gradingPeriodsCmd.AddCommand(newGradingPeriodsListCmd())
	gradingPeriodsCmd.AddCommand(newGradingPeriodsGetCmd())
	gradingPeriodsCmd.AddCommand(newGradingPeriodsUpdateCmd())
	gradingPeriodsCmd.AddCommand(newGradingPeriodsDeleteCmd())
}

func newGradingPeriodsListCmd() *cobra.Command {
	opts := &options.GradingPeriodsListOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List grading periods for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGradingPeriodsList(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runGradingPeriodsList(ctx context.Context, client *api.Client, opts *options.GradingPeriodsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-periods.list", map[string]interface{}{"course_id": opts.CourseID})

	svc := api.NewGradingPeriodsService(client)
	periods, err := svc.List(ctx, opts.CourseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "grading-periods.list", len(periods))
	return formatEmptyOrOutput(periods, "No grading periods found")
}

func newGradingPeriodsGetCmd() *cobra.Command {
	opts := &options.GradingPeriodsGetOptions{}
	cmd := &cobra.Command{
		Use:   "get <grading-period-id>",
		Short: "Get a grading period",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &opts.ID); err != nil {
				return fmt.Errorf("invalid grading period ID: %s", args[0])
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGradingPeriodsGet(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runGradingPeriodsGet(ctx context.Context, client *api.Client, opts *options.GradingPeriodsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-periods.get", map[string]interface{}{"course_id": opts.CourseID, "id": opts.ID})

	svc := api.NewGradingPeriodsService(client)
	period, err := svc.Get(ctx, opts.CourseID, opts.ID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "grading-periods.get", 1)
	return formatOutput(period, nil)
}

func newGradingPeriodsUpdateCmd() *cobra.Command {
	opts := &options.GradingPeriodsUpdateOptions{}
	cmd := &cobra.Command{
		Use:   "update <grading-period-id>",
		Short: "Update a grading period",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &opts.ID); err != nil {
				return fmt.Errorf("invalid grading period ID: %s", args[0])
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGradingPeriodsUpdate(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Grading period title")
	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "End date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.CloseDate, "close-date", "", "Close date (YYYY-MM-DD)")
	cmd.Flags().Float64Var(&opts.Weight, "weight", 0, "Weight percentage")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runGradingPeriodsUpdate(ctx context.Context, client *api.Client, opts *options.GradingPeriodsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-periods.update", map[string]interface{}{"course_id": opts.CourseID, "id": opts.ID})

	svc := api.NewGradingPeriodsService(client)
	period, err := svc.Update(ctx, opts.CourseID, opts.ID, api.GradingPeriodParams{
		Title:     opts.Title,
		StartDate: opts.StartDate,
		EndDate:   opts.EndDate,
		CloseDate: opts.CloseDate,
		Weight:    opts.Weight,
	})
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "grading-periods.update", 1)
	return formatSuccessOutput(period, fmt.Sprintf("Grading period %d updated", period.ID))
}

func newGradingPeriodsDeleteCmd() *cobra.Command {
	opts := &options.GradingPeriodsDeleteOptions{}
	cmd := &cobra.Command{
		Use:   "delete <grading-period-id>",
		Short: "Delete a grading period",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &opts.ID); err != nil {
				return fmt.Errorf("invalid grading period ID: %s", args[0])
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGradingPeriodsDelete(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runGradingPeriodsDelete(ctx context.Context, client *api.Client, opts *options.GradingPeriodsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-periods.delete", map[string]interface{}{"course_id": opts.CourseID, "id": opts.ID})

	svc := api.NewGradingPeriodsService(client)
	if err := svc.Delete(ctx, opts.CourseID, opts.ID); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "grading-periods.delete", 1)
	printInfo("Grading period %d deleted\n", opts.ID)
	return nil
}
