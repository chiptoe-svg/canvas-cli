package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// blackoutDatesCmd is the parent command for blackout date operations.
var blackoutDatesCmd = &cobra.Command{
	Use:     "blackout-dates",
	Aliases: []string{"blackouts"},
	Short:   "Manage course blackout dates",
	Long: `Manage Canvas course blackout dates (periods excluded from pacing).

Examples:
  canvas blackout-dates list --course-id 1
  canvas blackout-dates get 5 --course-id 1
  canvas blackout-dates create --course-id 1 --start-date 2024-12-24 --end-date 2024-12-26 --title "Winter Break"
  canvas blackout-dates update 5 --course-id 1 --title "Extended Break"
  canvas blackout-dates delete 5 --course-id 1`,
}

func init() {
	rootCmd.AddCommand(blackoutDatesCmd)
	blackoutDatesCmd.AddCommand(newBlackoutDatesListCmd())
	blackoutDatesCmd.AddCommand(newBlackoutDatesGetCmd())
	blackoutDatesCmd.AddCommand(newBlackoutDatesCreateCmd())
	blackoutDatesCmd.AddCommand(newBlackoutDatesUpdateCmd())
	blackoutDatesCmd.AddCommand(newBlackoutDatesDeleteCmd())
}

func newBlackoutDatesListCmd() *cobra.Command {
	opts := &options.BlackoutDatesListOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List blackout dates for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runBlackoutDatesList(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runBlackoutDatesList(ctx context.Context, client *api.Client, opts *options.BlackoutDatesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blackout-dates.list", map[string]interface{}{"course_id": opts.CourseID})

	svc := api.NewBlackoutDatesService(client)
	dates, err := svc.List(ctx, opts.CourseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "blackout-dates.list", len(dates))
	return formatEmptyOrOutput(dates, "No blackout dates found")
}

func newBlackoutDatesGetCmd() *cobra.Command {
	opts := &options.BlackoutDatesGetOptions{}
	cmd := &cobra.Command{
		Use:   "get <blackout-date-id>",
		Short: "Get a blackout date",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &opts.ID); err != nil {
				return fmt.Errorf("invalid blackout date ID: %s", args[0])
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runBlackoutDatesGet(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runBlackoutDatesGet(ctx context.Context, client *api.Client, opts *options.BlackoutDatesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blackout-dates.get", map[string]interface{}{"course_id": opts.CourseID, "id": opts.ID})

	svc := api.NewBlackoutDatesService(client)
	date, err := svc.Get(ctx, opts.CourseID, opts.ID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "blackout-dates.get", 1)
	return formatOutput(date, nil)
}

func newBlackoutDatesCreateCmd() *cobra.Command {
	opts := &options.BlackoutDatesCreateOptions{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a blackout date for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runBlackoutDatesCreate(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Start date (YYYY-MM-DD) (required)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "End date (YYYY-MM-DD) (required)")
	cmd.Flags().StringVar(&opts.EventTitle, "title", "", "Event title (required)")
	mustMarkRequired(cmd, "course-id", "start-date", "end-date", "title")
	return cmd
}

func runBlackoutDatesCreate(ctx context.Context, client *api.Client, opts *options.BlackoutDatesCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blackout-dates.create", map[string]interface{}{"course_id": opts.CourseID})

	svc := api.NewBlackoutDatesService(client)
	date, err := svc.Create(ctx, opts.CourseID, api.BlackoutDateParams{
		StartDate:  opts.StartDate,
		EndDate:    opts.EndDate,
		EventTitle: opts.EventTitle,
	})
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "blackout-dates.create", 1)
	return formatSuccessOutput(date, fmt.Sprintf("Blackout date %d created", date.ID))
}

func newBlackoutDatesUpdateCmd() *cobra.Command {
	opts := &options.BlackoutDatesUpdateOptions{}
	cmd := &cobra.Command{
		Use:   "update <blackout-date-id>",
		Short: "Update a blackout date",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &opts.ID); err != nil {
				return fmt.Errorf("invalid blackout date ID: %s", args[0])
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runBlackoutDatesUpdate(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "End date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EventTitle, "title", "", "Event title")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runBlackoutDatesUpdate(ctx context.Context, client *api.Client, opts *options.BlackoutDatesUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blackout-dates.update", map[string]interface{}{"course_id": opts.CourseID, "id": opts.ID})

	svc := api.NewBlackoutDatesService(client)
	date, err := svc.Update(ctx, opts.CourseID, opts.ID, api.BlackoutDateParams{
		StartDate:  opts.StartDate,
		EndDate:    opts.EndDate,
		EventTitle: opts.EventTitle,
	})
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "blackout-dates.update", 1)
	return formatSuccessOutput(date, fmt.Sprintf("Blackout date %d updated", date.ID))
}

func newBlackoutDatesDeleteCmd() *cobra.Command {
	opts := &options.BlackoutDatesDeleteOptions{}
	cmd := &cobra.Command{
		Use:   "delete <blackout-date-id>",
		Short: "Delete a blackout date",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &opts.ID); err != nil {
				return fmt.Errorf("invalid blackout date ID: %s", args[0])
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runBlackoutDatesDelete(cmd.Context(), client, opts)
		},
	}
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runBlackoutDatesDelete(ctx context.Context, client *api.Client, opts *options.BlackoutDatesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blackout-dates.delete", map[string]interface{}{"course_id": opts.CourseID, "id": opts.ID})

	svc := api.NewBlackoutDatesService(client)
	if err := svc.Delete(ctx, opts.CourseID, opts.ID); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "blackout-dates.delete", 1)
	printInfo("Blackout date %d deleted\n", opts.ID)
	return nil
}
