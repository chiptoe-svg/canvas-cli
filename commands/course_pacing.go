package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// coursePacingCmd is the parent command for course pacing operations.
var coursePacingCmd = &cobra.Command{
	Use:     "course-pacing",
	Aliases: []string{"pacing", "pace"},
	Short:   "Manage course pacing",
	Long: `Manage Canvas course pacing (course pace) objects.

Examples:
  canvas course-pacing get 1 --course-id 5
  canvas course-pacing create --course-id 5 --exclude-weekends
  canvas course-pacing update 1 --course-id 5 --hard-end-dates
  canvas course-pacing delete 1 --course-id 5`,
}

func init() {
	rootCmd.AddCommand(coursePacingCmd)
	coursePacingCmd.AddCommand(newCoursePacingGetCmd())
	coursePacingCmd.AddCommand(newCoursePacingCreateCmd())
	coursePacingCmd.AddCommand(newCoursePacingUpdateCmd())
	coursePacingCmd.AddCommand(newCoursePacingDeleteCmd())
}

func newCoursePacingGetCmd() *cobra.Command {
	opts := &options.CoursePacingGetOptions{}
	cmd := &cobra.Command{
		Use:   "get <pace-id>",
		Short: "Get a course pace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &opts.PaceID); err != nil {
				return fmt.Errorf("invalid pace ID: %s", args[0])
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCoursePacingGet(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runCoursePacingGet(ctx context.Context, client *api.Client, opts *options.CoursePacingGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-pacing.get", map[string]interface{}{"course_id": opts.CourseID, "pace_id": opts.PaceID})

	svc := api.NewCoursePacingService(client)
	pace, err := svc.Get(ctx, opts.CourseID, opts.PaceID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-pacing.get", 1)
	return formatOutput(pace, nil)
}

func newCoursePacingCreateCmd() *cobra.Command {
	opts := &options.CoursePacingCreateOptions{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a course pace",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCoursePacingCreate(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().BoolVar(&opts.ExcludeWeekends, "exclude-weekends", false, "Exclude weekends from pacing")
	cmd.Flags().BoolVar(&opts.HardEndDates, "hard-end-dates", false, "Use hard end dates")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "Pace end date (YYYY-MM-DD)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runCoursePacingCreate(ctx context.Context, client *api.Client, opts *options.CoursePacingCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-pacing.create", map[string]interface{}{"course_id": opts.CourseID})

	svc := api.NewCoursePacingService(client)
	pace, err := svc.Create(ctx, opts.CourseID, api.CoursePaceParams{
		ExcludeWeekends: boolPtr(opts.ExcludeWeekends),
		HardEndDates:    boolPtr(opts.HardEndDates),
		EndDate:         opts.EndDate,
	})
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-pacing.create", 1)
	return formatSuccessOutput(pace, fmt.Sprintf("Course pace %d created", pace.ID))
}

func newCoursePacingUpdateCmd() *cobra.Command {
	opts := &options.CoursePacingUpdateOptions{}
	cmd := &cobra.Command{
		Use:   "update <pace-id>",
		Short: "Update a course pace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &opts.PaceID); err != nil {
				return fmt.Errorf("invalid pace ID: %s", args[0])
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCoursePacingUpdate(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().BoolVar(&opts.ExcludeWeekends, "exclude-weekends", false, "Exclude weekends from pacing")
	cmd.Flags().BoolVar(&opts.HardEndDates, "hard-end-dates", false, "Use hard end dates")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "Pace end date (YYYY-MM-DD)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runCoursePacingUpdate(ctx context.Context, client *api.Client, opts *options.CoursePacingUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-pacing.update", map[string]interface{}{"course_id": opts.CourseID, "pace_id": opts.PaceID})

	svc := api.NewCoursePacingService(client)
	pace, err := svc.Update(ctx, opts.CourseID, opts.PaceID, api.CoursePaceParams{
		ExcludeWeekends: boolPtr(opts.ExcludeWeekends),
		HardEndDates:    boolPtr(opts.HardEndDates),
		EndDate:         opts.EndDate,
	})
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-pacing.update", 1)
	return formatSuccessOutput(pace, fmt.Sprintf("Course pace %d updated", pace.ID))
}

func newCoursePacingDeleteCmd() *cobra.Command {
	opts := &options.CoursePacingDeleteOptions{}
	cmd := &cobra.Command{
		Use:   "delete <pace-id>",
		Short: "Delete a course pace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &opts.PaceID); err != nil {
				return fmt.Errorf("invalid pace ID: %s", args[0])
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCoursePacingDelete(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runCoursePacingDelete(ctx context.Context, client *api.Client, opts *options.CoursePacingDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-pacing.delete", map[string]interface{}{"course_id": opts.CourseID, "pace_id": opts.PaceID})

	svc := api.NewCoursePacingService(client)
	if err := svc.Delete(ctx, opts.CourseID, opts.PaceID); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-pacing.delete", 1)
	printInfo("Course pace %d deleted\n", opts.PaceID)
	return nil
}
