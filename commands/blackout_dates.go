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

// blackoutDatesCmd is the root command for blackout dates
var blackoutDatesCmd = &cobra.Command{
	Use:   "blackout-dates",
	Short: "Manage Canvas account blackout dates",
	Long: `Manage blackout dates for Canvas accounts.

Blackout dates are date ranges during which no course activities are scheduled.
They are used to block out holidays, breaks, and other non-instructional periods.

Examples:
  canvas blackout-dates list 1
  canvas blackout-dates get 1 5
  canvas blackout-dates create 1 --start-date 2024-01-01 --end-date 2024-01-02 --title "New Year"
  canvas blackout-dates update 1 5 --title "Extended Holiday"
  canvas blackout-dates delete 1 5`,
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
		Use:   "list <account-id>",
		Short: "List blackout dates for an account",
		Args:  ExactArgsWithUsage(1, "account-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID

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

	return cmd
}

func newBlackoutDatesGetCmd() *cobra.Command {
	opts := &options.BlackoutDatesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <account-id> <id>",
		Short: "Get a specific blackout date",
		Args:  ExactArgsWithUsage(2, "account-id", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			id, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid blackout date ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.ID = id

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

	return cmd
}

func newBlackoutDatesCreateCmd() *cobra.Command {
	opts := &options.BlackoutDatesCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <account-id>",
		Short: "Create a blackout date for an account",
		Args:  ExactArgsWithUsage(1, "account-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID

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

	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "End date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EventTitle, "title", "", "Event title for the blackout date")
	mustMarkRequired(cmd, "start-date", "end-date")

	return cmd
}

func newBlackoutDatesUpdateCmd() *cobra.Command {
	opts := &options.BlackoutDatesUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <account-id> <id>",
		Short: "Update a blackout date",
		Args:  ExactArgsWithUsage(2, "account-id", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			id, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid blackout date ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.ID = id

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

	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "End date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EventTitle, "title", "", "Event title for the blackout date")

	return cmd
}

func newBlackoutDatesDeleteCmd() *cobra.Command {
	opts := &options.BlackoutDatesDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <account-id> <id>",
		Short: "Delete a blackout date",
		Args:  ExactArgsWithUsage(2, "account-id", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			id, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid blackout date ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.ID = id

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

	return cmd
}

func runBlackoutDatesList(ctx context.Context, client *api.Client, opts *options.BlackoutDatesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blackout-dates.list", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAccountBlackoutDatesService(client)
	dates, err := svc.List(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "blackout-dates.list", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "blackout-dates.list", len(dates))

	return formatEmptyOrOutput(dates, fmt.Sprintf("No blackout dates found for account %d", opts.AccountID))
}

func runBlackoutDatesGet(ctx context.Context, client *api.Client, opts *options.BlackoutDatesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blackout-dates.get", map[string]interface{}{
		"account_id": opts.AccountID,
		"id":         opts.ID,
	})

	svc := api.NewAccountBlackoutDatesService(client)
	date, err := svc.Get(ctx, opts.AccountID, opts.ID)
	if err != nil {
		logger.LogCommandError(ctx, "blackout-dates.get", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"id":         opts.ID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "blackout-dates.get", 1)

	return formatOutput(date, func() {
		fmt.Printf("ID:          %d\n", date.ID)
		fmt.Printf("Event Title: %s\n", date.EventTitle)
		fmt.Printf("Start Date:  %s\n", date.StartDate)
		fmt.Printf("End Date:    %s\n", date.EndDate)
		fmt.Printf("Context:     %s %d\n", date.ContextType, date.ContextID)
	})
}

func runBlackoutDatesCreate(ctx context.Context, client *api.Client, opts *options.BlackoutDatesCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blackout-dates.create", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAccountBlackoutDatesService(client)
	params := &api.BlackoutDateParams{
		StartDate:  opts.StartDate,
		EndDate:    opts.EndDate,
		EventTitle: opts.EventTitle,
	}

	date, err := svc.Create(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "blackout-dates.create", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "blackout-dates.create", 1)

	return formatSuccessOutput(date, fmt.Sprintf("Blackout date created (ID: %d)", date.ID))
}

func runBlackoutDatesUpdate(ctx context.Context, client *api.Client, opts *options.BlackoutDatesUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blackout-dates.update", map[string]interface{}{
		"account_id": opts.AccountID,
		"id":         opts.ID,
	})

	svc := api.NewAccountBlackoutDatesService(client)
	params := &api.BlackoutDateParams{
		StartDate:  opts.StartDate,
		EndDate:    opts.EndDate,
		EventTitle: opts.EventTitle,
	}

	date, err := svc.Update(ctx, opts.AccountID, opts.ID, params)
	if err != nil {
		logger.LogCommandError(ctx, "blackout-dates.update", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"id":         opts.ID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "blackout-dates.update", 1)

	return formatSuccessOutput(date, fmt.Sprintf("Blackout date updated (ID: %d)", date.ID))
}

func runBlackoutDatesDelete(ctx context.Context, client *api.Client, opts *options.BlackoutDatesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blackout-dates.delete", map[string]interface{}{
		"account_id": opts.AccountID,
		"id":         opts.ID,
	})

	svc := api.NewAccountBlackoutDatesService(client)
	if err := svc.Delete(ctx, opts.AccountID, opts.ID); err != nil {
		logger.LogCommandError(ctx, "blackout-dates.delete", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"id":         opts.ID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "blackout-dates.delete", 1)

	printInfo("Blackout date %d deleted\n", opts.ID)

	return nil
}
