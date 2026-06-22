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

// accountAnalyticsCmd represents the account-analytics command group
var accountAnalyticsCmd = &cobra.Command{
	Use:   "account-analytics",
	Short: "View Canvas account-level analytics",
	Long: `View analytics data for Canvas accounts including activity, grades, and statistics.

Analytics are organized by term or by course completion status (current/completed).

Examples:
  canvas account-analytics term-activity 1 2
  canvas account-analytics term-grades 1 2
  canvas account-analytics term-statistics 1 2
  canvas account-analytics completed-statistics 1`,
}

func init() {
	rootCmd.AddCommand(accountAnalyticsCmd)
	accountAnalyticsCmd.AddCommand(newAccountAnalyticsTermActivityCmd())
	accountAnalyticsCmd.AddCommand(newAccountAnalyticsTermGradesCmd())
	accountAnalyticsCmd.AddCommand(newAccountAnalyticsTermStatisticsCmd())
	accountAnalyticsCmd.AddCommand(newAccountAnalyticsCompletedStatisticsCmd())
}

func newAccountAnalyticsTermActivityCmd() *cobra.Command {
	opts := &options.AccountAnalyticsTermOptions{}

	cmd := &cobra.Command{
		Use:   "term-activity <account-id> <term-id>",
		Short: "Get activity analytics for an account term",
		Args:  ExactArgsWithUsage(2, "account-id", "term-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			termID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid term ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.TermID = termID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountAnalyticsTermActivity(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountAnalyticsTermGradesCmd() *cobra.Command {
	opts := &options.AccountAnalyticsTermOptions{}

	cmd := &cobra.Command{
		Use:   "term-grades <account-id> <term-id>",
		Short: "Get grade analytics for an account term",
		Args:  ExactArgsWithUsage(2, "account-id", "term-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			termID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid term ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.TermID = termID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountAnalyticsTermGrades(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountAnalyticsTermStatisticsCmd() *cobra.Command {
	opts := &options.AccountAnalyticsTermOptions{}

	cmd := &cobra.Command{
		Use:   "term-statistics <account-id> <term-id>",
		Short: "Get statistics for an account term",
		Args:  ExactArgsWithUsage(2, "account-id", "term-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			termID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid term ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.TermID = termID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountAnalyticsTermStatistics(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountAnalyticsCompletedStatisticsCmd() *cobra.Command {
	opts := &options.AccountAnalyticsCompletedOptions{}

	cmd := &cobra.Command{
		Use:   "completed-statistics <account-id>",
		Short: "Get statistics for completed courses in an account",
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

			return runAccountAnalyticsCompletedStatistics(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func runAccountAnalyticsTermActivity(ctx context.Context, client *api.Client, opts *options.AccountAnalyticsTermOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-analytics.term-activity", map[string]interface{}{
		"account_id": opts.AccountID,
		"term_id":    opts.TermID,
	})

	svc := api.NewAccountAnalyticsService(client)
	result, err := svc.GetTermActivity(ctx, opts.AccountID, opts.TermID)
	if err != nil {
		logger.LogCommandError(ctx, "account-analytics.term-activity", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-analytics.term-activity", len(result))

	return formatOutput(result, func() {
		fmt.Printf("Term Activity (Account %d, Term %d)\n", opts.AccountID, opts.TermID)
		fmt.Printf("Total entries: %d\n", len(result))
	})
}

func runAccountAnalyticsTermGrades(ctx context.Context, client *api.Client, opts *options.AccountAnalyticsTermOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-analytics.term-grades", map[string]interface{}{
		"account_id": opts.AccountID,
		"term_id":    opts.TermID,
	})

	svc := api.NewAccountAnalyticsService(client)
	result, err := svc.GetTermGrades(ctx, opts.AccountID, opts.TermID)
	if err != nil {
		logger.LogCommandError(ctx, "account-analytics.term-grades", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-analytics.term-grades", len(result))

	return formatOutput(result, func() {
		fmt.Printf("Term Grades (Account %d, Term %d)\n", opts.AccountID, opts.TermID)
		fmt.Printf("Total entries: %d\n", len(result))
	})
}

func runAccountAnalyticsTermStatistics(ctx context.Context, client *api.Client, opts *options.AccountAnalyticsTermOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-analytics.term-statistics", map[string]interface{}{
		"account_id": opts.AccountID,
		"term_id":    opts.TermID,
	})

	svc := api.NewAccountAnalyticsService(client)
	result, err := svc.GetTermStatistics(ctx, opts.AccountID, opts.TermID)
	if err != nil {
		logger.LogCommandError(ctx, "account-analytics.term-statistics", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-analytics.term-statistics", 1)

	return formatOutput(result, func() {
		fmt.Printf("Term Statistics (Account %d, Term %d)\n", opts.AccountID, opts.TermID)
		fmt.Println()
		for k, v := range result {
			fmt.Printf("  %-30s %v\n", k+":", v)
		}
	})
}

func runAccountAnalyticsCompletedStatistics(ctx context.Context, client *api.Client, opts *options.AccountAnalyticsCompletedOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-analytics.completed-statistics", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAccountAnalyticsService(client)
	result, err := svc.GetCompletedStatistics(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "account-analytics.completed-statistics", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-analytics.completed-statistics", 1)

	return formatOutput(result, func() {
		fmt.Printf("Completed Course Statistics (Account %d)\n", opts.AccountID)
		fmt.Println()
		for k, v := range result {
			fmt.Printf("  %-30s %v\n", k+":", v)
		}
	})
}
