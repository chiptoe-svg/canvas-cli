package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var gradingStandardsCmd = &cobra.Command{
	Use:   "grading-standards",
	Short: "Manage Canvas grading standards",
	Long: `Manage Canvas grading standards for accounts.

Grading standards define the letter grade scheme used for courses.

Examples:
  canvas grading-standards list 1
  canvas grading-standards get 1 5
  canvas grading-standards create 1 --title "Standard Grading"
  canvas grading-standards delete 1 5`,
}

func init() {
	rootCmd.AddCommand(gradingStandardsCmd)
	gradingStandardsCmd.AddCommand(newGradingStandardsListCmd())
	gradingStandardsCmd.AddCommand(newGradingStandardsGetCmd())
	gradingStandardsCmd.AddCommand(newGradingStandardsCreateCmd())
	gradingStandardsCmd.AddCommand(newGradingStandardsDeleteCmd())
}

func newGradingStandardsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <account-id>",
		Short: "List grading standards for an account",
		Args:  ExactArgsWithUsage(1, "account-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGradingStandardsList(cmd.Context(), client, accountID)
		},
	}
}

func newGradingStandardsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <account-id> <standard-id>",
		Short: "Get a specific grading standard",
		Args:  ExactArgsWithUsage(2, "account-id", "standard-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			standardID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid standard ID: %s", args[1])
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGradingStandardsGet(cmd.Context(), client, accountID, standardID)
		},
	}
}

func newGradingStandardsCreateCmd() *cobra.Command {
	var title string

	cmd := &cobra.Command{
		Use:   "create <account-id>",
		Short: "Create a grading standard for an account",
		Args:  ExactArgsWithUsage(1, "account-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			params := &api.GradingStandardParams{Title: title}
			return runGradingStandardsCreate(cmd.Context(), client, accountID, params)
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Title of the grading standard (required)")
	mustMarkRequired(cmd, "title")

	return cmd
}

func newGradingStandardsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <account-id> <standard-id>",
		Short: "Delete a grading standard",
		Args:  ExactArgsWithUsage(2, "account-id", "standard-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			standardID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid standard ID: %s", args[1])
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGradingStandardsDelete(cmd.Context(), client, accountID, standardID)
		},
	}
}

func runGradingStandardsList(ctx context.Context, client *api.Client, accountID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-standards.list", map[string]interface{}{"account_id": accountID})

	svc := api.NewGradingStandardsService(client)
	standards, err := svc.List(ctx, accountID)
	if err != nil {
		logger.LogCommandError(ctx, "grading-standards.list", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "grading-standards.list", len(standards))
	return formatEmptyOrOutput(standards, fmt.Sprintf("No grading standards found for account %d", accountID))
}

func runGradingStandardsGet(ctx context.Context, client *api.Client, accountID, standardID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-standards.get", map[string]interface{}{"account_id": accountID, "standard_id": standardID})

	svc := api.NewGradingStandardsService(client)
	standard, err := svc.Get(ctx, accountID, standardID)
	if err != nil {
		logger.LogCommandError(ctx, "grading-standards.get", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "grading-standards.get", 1)
	return formatOutput(standard, nil)
}

func runGradingStandardsCreate(ctx context.Context, client *api.Client, accountID int64, params *api.GradingStandardParams) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-standards.create", map[string]interface{}{"account_id": accountID, "title": params.Title})

	svc := api.NewGradingStandardsService(client)
	standard, err := svc.Create(ctx, accountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "grading-standards.create", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "grading-standards.create", 1)
	return formatSuccessOutput(standard, fmt.Sprintf("Grading standard created (ID: %d)", standard.ID))
}

func runGradingStandardsDelete(ctx context.Context, client *api.Client, accountID, standardID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grading-standards.delete", map[string]interface{}{"account_id": accountID, "standard_id": standardID})

	svc := api.NewGradingStandardsService(client)
	if err := svc.Delete(ctx, accountID, standardID); err != nil {
		logger.LogCommandError(ctx, "grading-standards.delete", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "grading-standards.delete", 1)
	printInfo("Grading standard %d deleted\n", standardID)
	return nil
}
