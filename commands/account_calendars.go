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

// accountCalendarsCmd represents the account-calendars command group.
var accountCalendarsCmd = &cobra.Command{
	Use:   "account-calendars",
	Short: "Manage Canvas account calendars",
	Long: `Manage Canvas account-level calendars and their visibility settings.

Examples:
  canvas account-calendars list
  canvas account-calendars list --account-id 1
  canvas account-calendars get 1
  canvas account-calendars update 1 --visible`,
}

func init() {
	rootCmd.AddCommand(accountCalendarsCmd)
	accountCalendarsCmd.AddCommand(newAccountCalendarsListCmd())
	accountCalendarsCmd.AddCommand(newAccountCalendarGetCmd())
	accountCalendarsCmd.AddCommand(newAccountCalendarUpdateCmd())
}

func newAccountCalendarsListCmd() *cobra.Command {
	opts := &options.AccountCalendarsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List account calendars",
		Long: `List visible account calendars for the current user, or all calendars
under a given account (requires admin permissions).

Examples:
  canvas account-calendars list
  canvas account-calendars list --account-id 1 --filter "Math"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountCalendarsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (list sub-account calendars, requires admin)")
	cmd.Flags().StringVar(&opts.Filter, "filter", "", "Search term to filter calendars")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Number of results per page")

	return cmd
}

func newAccountCalendarGetCmd() *cobra.Command {
	opts := &options.AccountCalendarGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <account-id>",
		Short: "Get a single account calendar",
		Long: `Retrieve details of an account calendar by account ID.

Examples:
  canvas account-calendars get 1`,
		Args: ExactArgsWithUsage(1, "account-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountCalendarGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountCalendarUpdateCmd() *cobra.Command {
	opts := &options.AccountCalendarUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <account-id>",
		Short: "Update account calendar visibility",
		Long: `Update the visibility and auto-subscribe settings of an account calendar.

Examples:
  canvas account-calendars update 1 --visible
  canvas account-calendars update 1 --no-visible --auto-subscribe`,
		Args: ExactArgsWithUsage(1, "account-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = id
			opts.VisibleSet = cmd.Flags().Changed("visible")
			opts.AutoSubSet = cmd.Flags().Changed("auto-subscribe")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountCalendarUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Visible, "visible", false, "Make the calendar visible to users")
	cmd.Flags().BoolVar(&opts.AutoSubscribe, "auto-subscribe", false, "Auto-subscribe users to this calendar")

	return cmd
}

func runAccountCalendarsList(ctx context.Context, client *api.Client, opts *options.AccountCalendarsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-calendars.list", map[string]interface{}{
		"account_id": opts.AccountID,
		"filter":     opts.Filter,
	})

	svc := api.NewAccountCalendarsService(client)

	apiOpts := &api.ListAccountCalendarsOptions{
		Filter:  opts.Filter,
		PerPage: opts.PerPage,
	}

	var (
		cals []api.AccountCalendar
		err  error
	)

	if opts.AccountID > 0 {
		cals, err = svc.ListForAccount(ctx, opts.AccountID, apiOpts)
	} else {
		cals, err = svc.List(ctx, apiOpts)
	}

	if err != nil {
		logger.LogCommandError(ctx, "account-calendars.list", err, nil)
		return fmt.Errorf("failed to list account calendars: %w", err)
	}

	printVerbose("Found %d account calendars:\n\n", len(cals))
	logger.LogCommandComplete(ctx, "account-calendars.list", len(cals))
	return formatEmptyOrOutput(cals, "No account calendars found")
}

func runAccountCalendarGet(ctx context.Context, client *api.Client, opts *options.AccountCalendarGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-calendars.get", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAccountCalendarsService(client)

	cal, err := svc.Get(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "account-calendars.get", err, nil)
		return fmt.Errorf("failed to get account calendar: %w", err)
	}

	logger.LogCommandComplete(ctx, "account-calendars.get", 1)
	return formatOutput(cal, nil)
}

func runAccountCalendarUpdate(ctx context.Context, client *api.Client, opts *options.AccountCalendarUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-calendars.update", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAccountCalendarsService(client)

	params := &api.UpdateAccountCalendarParams{}
	if opts.VisibleSet {
		params.Visible = &opts.Visible
	}
	if opts.AutoSubSet {
		params.AutoSubscribe = &opts.AutoSubscribe
	}

	cal, err := svc.Update(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "account-calendars.update", err, nil)
		return fmt.Errorf("failed to update account calendar: %w", err)
	}

	logger.LogCommandComplete(ctx, "account-calendars.update", 1)
	return formatSuccessOutput(cal, "Account calendar updated successfully.")
}
