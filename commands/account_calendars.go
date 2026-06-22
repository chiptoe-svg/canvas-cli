package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// accountCalendarsCmd represents the account-calendars command group
var accountCalendarsCmd = &cobra.Command{
	Use:   "account-calendars",
	Short: "Manage Canvas account calendars",
	Long: `Manage Canvas account calendars.

Account calendars allow institutions to share events with students and teachers
across courses within an account.

Examples:
  canvas account-calendars list
  canvas account-calendars list --search "main"
  canvas account-calendars get 1
  canvas account-calendars update 1 --visible
  canvas account-calendars list-for-account 1`,
}

func init() {
	rootCmd.AddCommand(accountCalendarsCmd)
	accountCalendarsCmd.AddCommand(newAccountCalendarsListCmd())
	accountCalendarsCmd.AddCommand(newAccountCalendarsGetCmd())
	accountCalendarsCmd.AddCommand(newAccountCalendarsUpdateCmd())
	accountCalendarsCmd.AddCommand(newAccountCalendarsListForAccountCmd())
}

func newAccountCalendarsListCmd() *cobra.Command {
	opts := &options.AccountCalendarsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all account calendars visible to the current user",
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

	cmd.Flags().StringVar(&opts.Search, "search", "", "Filter calendars by name")

	return cmd
}

func newAccountCalendarsGetCmd() *cobra.Command {
	opts := &options.AccountCalendarsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <account-id>",
		Short: "Get an account calendar by account ID",
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

			return runAccountCalendarsGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountCalendarsUpdateCmd() *cobra.Command {
	opts := &options.AccountCalendarsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <account-id>",
		Short: "Update an account calendar",
		Args:  ExactArgsWithUsage(1, "account-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID
			opts.VisibleSet = cmd.Flags().Changed("visible")
			opts.AutoSubSet = cmd.Flags().Changed("auto-subscribe")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountCalendarsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Visible, "visible", false, "Whether the calendar is visible")
	cmd.Flags().BoolVar(&opts.AutoSubscribe, "auto-subscribe", false, "Whether users are auto-subscribed")

	return cmd
}

func newAccountCalendarsListForAccountCmd() *cobra.Command {
	opts := &options.AccountCalendarsListForAccountOptions{}

	cmd := &cobra.Command{
		Use:   "list-for-account <account-id>",
		Short: "List account calendars for a specific account",
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

			return runAccountCalendarsListForAccount(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Search, "search", "", "Filter calendars by name")

	return cmd
}

func runAccountCalendarsList(ctx context.Context, client *api.Client, opts *options.AccountCalendarsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-calendars.list", map[string]interface{}{
		"search": opts.Search,
	})

	svc := api.NewAccountCalendarsService(client)
	calendars, err := svc.ListAll(ctx, opts.Search)
	if err != nil {
		logger.LogCommandError(ctx, "account-calendars.list", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-calendars.list", len(calendars))

	return formatOutput(calendars, func() {
		fmt.Printf("%-8s %-40s %-8s %-14s\n", "ID", "NAME", "VISIBLE", "AUTO_SUBSCRIBE")
		fmt.Println(strings.Repeat("-", 75))
		for _, c := range calendars {
			name := c.Name
			if len(name) > 38 {
				name = name[:35] + "..."
			}
			fmt.Printf("%-8d %-40s %-8v %-14v\n", c.ID, name, c.Visible, c.AutoSubscribe)
		}
		fmt.Printf("\nTotal: %d calendar(s)\n", len(calendars))
	})
}

func runAccountCalendarsGet(ctx context.Context, client *api.Client, opts *options.AccountCalendarsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-calendars.get", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAccountCalendarsService(client)
	cal, err := svc.Get(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "account-calendars.get", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-calendars.get", 1)

	return formatOutput(cal, func() {
		fmt.Printf("Account Calendar Details\n")
		fmt.Printf("========================\n\n")
		fmt.Printf("ID:               %d\n", cal.ID)
		fmt.Printf("Name:             %s\n", cal.Name)
		fmt.Printf("Visible:          %v\n", cal.Visible)
		fmt.Printf("Auto Subscribe:   %v\n", cal.AutoSubscribe)
		if cal.ParentAccountID != 0 {
			fmt.Printf("Parent Account:   %d\n", cal.ParentAccountID)
		}
		if cal.RootAccountID != 0 {
			fmt.Printf("Root Account:     %d\n", cal.RootAccountID)
		}
	})
}

func runAccountCalendarsUpdate(ctx context.Context, client *api.Client, opts *options.AccountCalendarsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-calendars.update", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	params := &api.AccountCalendarParams{}
	if opts.VisibleSet {
		visible := opts.Visible
		params.Visible = &visible
	}
	if opts.AutoSubSet {
		autoSub := opts.AutoSubscribe
		params.AutoSubscribe = &autoSub
	}

	svc := api.NewAccountCalendarsService(client)
	cal, err := svc.Update(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "account-calendars.update", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-calendars.update", 1)

	return formatSuccessOutput(cal, fmt.Sprintf("Account calendar %d updated", opts.AccountID))
}

func runAccountCalendarsListForAccount(ctx context.Context, client *api.Client, opts *options.AccountCalendarsListForAccountOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-calendars.list-for-account", map[string]interface{}{
		"account_id": opts.AccountID,
		"search":     opts.Search,
	})

	svc := api.NewAccountCalendarsService(client)
	calendars, err := svc.ListForAccount(ctx, opts.AccountID, opts.Search)
	if err != nil {
		logger.LogCommandError(ctx, "account-calendars.list-for-account", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-calendars.list-for-account", len(calendars))

	return formatOutput(calendars, func() {
		fmt.Printf("%-8s %-40s %-8s %-14s\n", "ID", "NAME", "VISIBLE", "AUTO_SUBSCRIBE")
		fmt.Println(strings.Repeat("-", 75))
		for _, c := range calendars {
			name := c.Name
			if len(name) > 38 {
				name = name[:35] + "..."
			}
			fmt.Printf("%-8d %-40s %-8v %-14v\n", c.ID, name, c.Visible, c.AutoSubscribe)
		}
		fmt.Printf("\nTotal: %d calendar(s)\n", len(calendars))
	})
}
