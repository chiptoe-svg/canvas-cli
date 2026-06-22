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

// temporaryEnrollmentPairingsCmd is the root command for temporary enrollment pairings
var temporaryEnrollmentPairingsCmd = &cobra.Command{
	Use:   "temporary-enrollment-pairings",
	Short: "Manage Canvas temporary enrollment pairings",
	Long: `Manage temporary enrollment pairings for Canvas accounts.

Temporary enrollment pairings allow temporary instructors to be enrolled
in courses for a limited time, inheriting permissions from a pairing configuration.

Examples:
  canvas temporary-enrollment-pairings list 1
  canvas temporary-enrollment-pairings get 1 5
  canvas temporary-enrollment-pairings create 1 --enrollment-state active
  canvas temporary-enrollment-pairings delete 1 5`,
}

func init() {
	rootCmd.AddCommand(temporaryEnrollmentPairingsCmd)
	temporaryEnrollmentPairingsCmd.AddCommand(newTemporaryEnrollmentPairingsListCmd())
	temporaryEnrollmentPairingsCmd.AddCommand(newTemporaryEnrollmentPairingsGetCmd())
	temporaryEnrollmentPairingsCmd.AddCommand(newTemporaryEnrollmentPairingsCreateCmd())
	temporaryEnrollmentPairingsCmd.AddCommand(newTemporaryEnrollmentPairingsDeleteCmd())
}

func newTemporaryEnrollmentPairingsListCmd() *cobra.Command {
	opts := &options.TemporaryEnrollmentPairingsListOptions{}

	cmd := &cobra.Command{
		Use:   "list <account-id>",
		Short: "List temporary enrollment pairings for an account",
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

			return runTemporaryEnrollmentPairingsList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newTemporaryEnrollmentPairingsGetCmd() *cobra.Command {
	opts := &options.TemporaryEnrollmentPairingsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <account-id> <id>",
		Short: "Get a specific temporary enrollment pairing",
		Args:  ExactArgsWithUsage(2, "account-id", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			id, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pairing ID: %s", args[1])
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

			return runTemporaryEnrollmentPairingsGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newTemporaryEnrollmentPairingsCreateCmd() *cobra.Command {
	opts := &options.TemporaryEnrollmentPairingsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <account-id>",
		Short: "Create a temporary enrollment pairing for an account",
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

			return runTemporaryEnrollmentPairingsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.StartingEnrollmentState, "enrollment-state", "", "Starting enrollment state (invited, active)")
	cmd.Flags().Int64Var(&opts.RoleID, "role-id", 0, "Role ID for the temporary enrollment")

	return cmd
}

func newTemporaryEnrollmentPairingsDeleteCmd() *cobra.Command {
	opts := &options.TemporaryEnrollmentPairingsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <account-id> <id>",
		Short: "Delete a temporary enrollment pairing",
		Args:  ExactArgsWithUsage(2, "account-id", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			id, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pairing ID: %s", args[1])
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

			return runTemporaryEnrollmentPairingsDelete(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func runTemporaryEnrollmentPairingsList(ctx context.Context, client *api.Client, opts *options.TemporaryEnrollmentPairingsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "temporary-enrollment-pairings.list", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAccountTemporaryEnrollmentPairingsService(client)
	pairings, err := svc.List(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "temporary-enrollment-pairings.list", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "temporary-enrollment-pairings.list", len(pairings))

	return formatEmptyOrOutput(pairings, fmt.Sprintf("No temporary enrollment pairings found for account %d", opts.AccountID))
}

func runTemporaryEnrollmentPairingsGet(ctx context.Context, client *api.Client, opts *options.TemporaryEnrollmentPairingsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "temporary-enrollment-pairings.get", map[string]interface{}{
		"account_id": opts.AccountID,
		"id":         opts.ID,
	})

	svc := api.NewAccountTemporaryEnrollmentPairingsService(client)
	pairing, err := svc.Get(ctx, opts.AccountID, opts.ID)
	if err != nil {
		logger.LogCommandError(ctx, "temporary-enrollment-pairings.get", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"id":         opts.ID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "temporary-enrollment-pairings.get", 1)

	return formatOutput(pairing, func() {
		fmt.Printf("ID:                      %d\n", pairing.ID)
		fmt.Printf("Root Account ID:         %d\n", pairing.RootAccountID)
		fmt.Printf("Workflow State:          %s\n", pairing.WorkflowState)
		fmt.Printf("Starting Enrollment State: %s\n", pairing.StartingEnrollmentState)
		fmt.Printf("Created At:              %s\n", pairing.CreatedAt)
		fmt.Printf("Updated At:              %s\n", pairing.UpdatedAt)
	})
}

func runTemporaryEnrollmentPairingsCreate(ctx context.Context, client *api.Client, opts *options.TemporaryEnrollmentPairingsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "temporary-enrollment-pairings.create", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAccountTemporaryEnrollmentPairingsService(client)
	params := &api.TemporaryEnrollmentPairingParams{
		StartingEnrollmentState: opts.StartingEnrollmentState,
		RoleID:                  opts.RoleID,
	}

	pairing, err := svc.Create(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "temporary-enrollment-pairings.create", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "temporary-enrollment-pairings.create", 1)

	return formatSuccessOutput(pairing, fmt.Sprintf("Temporary enrollment pairing created (ID: %d)", pairing.ID))
}

func runTemporaryEnrollmentPairingsDelete(ctx context.Context, client *api.Client, opts *options.TemporaryEnrollmentPairingsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "temporary-enrollment-pairings.delete", map[string]interface{}{
		"account_id": opts.AccountID,
		"id":         opts.ID,
	})

	svc := api.NewAccountTemporaryEnrollmentPairingsService(client)
	if err := svc.Delete(ctx, opts.AccountID, opts.ID); err != nil {
		logger.LogCommandError(ctx, "temporary-enrollment-pairings.delete", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"id":         opts.ID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "temporary-enrollment-pairings.delete", 1)

	printInfo("Temporary enrollment pairing %d deleted\n", opts.ID)

	return nil
}
