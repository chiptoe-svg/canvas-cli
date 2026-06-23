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

var gradingPeriodSetsCmd = &cobra.Command{
	Use:   "grading-period-sets",
	Short: "Manage Canvas grading period sets and grading periods",
	Long: `Manage grading period sets and individual grading periods within a Canvas account.

Examples:
  canvas grading-period-sets list 1
  canvas grading-period-sets create 1 --title "2024-2025"
  canvas grading-period-sets update 1 10 --title "2024-2025 Updated"
  canvas grading-period-sets delete 1 10
  canvas grading-period-sets list-periods 1
  canvas grading-period-sets delete-period 1 5`,
}

func init() {
	rootCmd.AddCommand(gradingPeriodSetsCmd)
	gradingPeriodSetsCmd.AddCommand(newGradingPeriodSetsListCmd())
	gradingPeriodSetsCmd.AddCommand(newGradingPeriodSetsCreateCmd())
	gradingPeriodSetsCmd.AddCommand(newGradingPeriodSetsUpdateCmd())
	gradingPeriodSetsCmd.AddCommand(newGradingPeriodSetsDeleteCmd())
	gradingPeriodSetsCmd.AddCommand(newGradingPeriodSetsListPeriodsCmd())
	gradingPeriodSetsCmd.AddCommand(newGradingPeriodSetsDeletePeriodCmd())
}

func newGradingPeriodSetsListCmd() *cobra.Command {
	opts := &options.GradingPeriodSetsListOptions{}

	cmd := &cobra.Command{
		Use:   "list <account-id>",
		Short: "List grading period sets for an account",
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

			return runGradingPeriodSetsList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newGradingPeriodSetsCreateCmd() *cobra.Command {
	opts := &options.GradingPeriodSetsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <account-id>",
		Short: "Create a grading period set",
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

			return runGradingPeriodSetsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Title, "title", "", "Title of the grading period set (required)")
	cmd.Flags().BoolVar(&opts.WeightedGradingPeriods, "weighted", false, "Use weighted grading periods")
	mustMarkRequired(cmd, "title")

	return cmd
}

func newGradingPeriodSetsUpdateCmd() *cobra.Command {
	opts := &options.GradingPeriodSetsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <account-id> <id>",
		Short: "Update a grading period set",
		Args:  ExactArgsWithUsage(2, "account-id", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			setID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid grading period set ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.SetID = setID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGradingPeriodSetsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Title, "title", "", "New title for the grading period set")

	return cmd
}

func newGradingPeriodSetsDeleteCmd() *cobra.Command {
	opts := &options.GradingPeriodSetsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <account-id> <id>",
		Short: "Delete a grading period set",
		Args:  ExactArgsWithUsage(2, "account-id", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			setID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid grading period set ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.SetID = setID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGradingPeriodSetsDelete(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newGradingPeriodSetsListPeriodsCmd() *cobra.Command {
	opts := &options.GradingPeriodSetsListPeriodsOptions{}

	cmd := &cobra.Command{
		Use:   "list-periods <account-id>",
		Short: "List grading periods for an account",
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

			return runGradingPeriodSetsListPeriods(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newGradingPeriodSetsDeletePeriodCmd() *cobra.Command {
	opts := &options.GradingPeriodSetsDeletePeriodOptions{}

	cmd := &cobra.Command{
		Use:   "delete-period <account-id> <id>",
		Short: "Delete a grading period",
		Args:  ExactArgsWithUsage(2, "account-id", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			periodID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid grading period ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.PeriodID = periodID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGradingPeriodSetsDeletePeriod(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func runGradingPeriodSetsList(ctx context.Context, client *api.Client, opts *options.GradingPeriodSetsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-period-sets.list", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	service := api.NewGradingPeriodSetsService(client)

	sets, err := service.List(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "grading-period-sets.list", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "grading-period-sets.list", len(sets))

	return formatEmptyOrOutput(sets, fmt.Sprintf("No grading period sets found for account %d", opts.AccountID))
}

func runGradingPeriodSetsCreate(ctx context.Context, client *api.Client, opts *options.GradingPeriodSetsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-period-sets.create", map[string]interface{}{
		"account_id": opts.AccountID,
		"title":      opts.Title,
	})

	service := api.NewGradingPeriodSetsService(client)

	params := &api.GradingPeriodSetParams{
		GradingPeriodSet: api.GradingPeriodSetFields{
			Title:                  opts.Title,
			WeightedGradingPeriods: boolPtr(opts.WeightedGradingPeriods),
		},
	}

	set, err := service.Create(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "grading-period-sets.create", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "grading-period-sets.create", 1)

	return formatSuccessOutput(set, fmt.Sprintf("Grading period set '%s' created with ID %d", set.Title, set.ID))
}

func runGradingPeriodSetsUpdate(ctx context.Context, client *api.Client, opts *options.GradingPeriodSetsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-period-sets.update", map[string]interface{}{
		"account_id": opts.AccountID,
		"set_id":     opts.SetID,
	})

	service := api.NewGradingPeriodSetsService(client)

	params := &api.GradingPeriodSetParams{
		GradingPeriodSet: api.GradingPeriodSetFields{
			Title: opts.Title,
		},
	}

	set, err := service.Update(ctx, opts.AccountID, opts.SetID, params)
	if err != nil {
		logger.LogCommandError(ctx, "grading-period-sets.update", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "grading-period-sets.update", 1)

	return formatSuccessOutput(set, fmt.Sprintf("Grading period set %d updated", set.ID))
}

func runGradingPeriodSetsDelete(ctx context.Context, client *api.Client, opts *options.GradingPeriodSetsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-period-sets.delete", map[string]interface{}{
		"account_id": opts.AccountID,
		"set_id":     opts.SetID,
	})

	service := api.NewGradingPeriodSetsService(client)

	if err := service.Delete(ctx, opts.AccountID, opts.SetID); err != nil {
		logger.LogCommandError(ctx, "grading-period-sets.delete", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "grading-period-sets.delete", 1)

	printInfo("Grading period set %d deleted\n", opts.SetID)
	return nil
}

func runGradingPeriodSetsListPeriods(ctx context.Context, client *api.Client, opts *options.GradingPeriodSetsListPeriodsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-period-sets.list-periods", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	service := api.NewGradingPeriodSetsService(client)

	periods, err := service.ListPeriods(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "grading-period-sets.list-periods", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "grading-period-sets.list-periods", len(periods))

	return formatEmptyOrOutput(periods, fmt.Sprintf("No grading periods found for account %d", opts.AccountID))
}

func runGradingPeriodSetsDeletePeriod(ctx context.Context, client *api.Client, opts *options.GradingPeriodSetsDeletePeriodOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-period-sets.delete-period", map[string]interface{}{
		"account_id": opts.AccountID,
		"period_id":  opts.PeriodID,
	})

	service := api.NewGradingPeriodSetsService(client)

	if err := service.DeletePeriod(ctx, opts.AccountID, opts.PeriodID); err != nil {
		logger.LogCommandError(ctx, "grading-period-sets.delete-period", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "grading-period-sets.delete-period", 1)

	printInfo("Grading period %d deleted\n", opts.PeriodID)
	return nil
}
