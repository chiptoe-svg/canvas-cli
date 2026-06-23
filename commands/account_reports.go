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

// accountReportsCmd is the root command for account reports
var accountReportsCmd = &cobra.Command{
	Use:   "account-reports",
	Short: "Manage Canvas account reports",
	Long: `Manage account-level reports in Canvas.

Account reports provide data exports and analytics about usage,
enrollments, grades, and other institutional data.

Examples:
  canvas account-reports list 1
  canvas account-reports runs 1 course_storage_csv
  canvas account-reports start 1 course_storage_csv
  canvas account-reports get 1 course_storage_csv 10
  canvas account-reports abort 1 course_storage_csv 10`,
}

func init() {
	rootCmd.AddCommand(accountReportsCmd)
	accountReportsCmd.AddCommand(newAccountReportsListCmd())
	accountReportsCmd.AddCommand(newAccountReportsRunsCmd())
	accountReportsCmd.AddCommand(newAccountReportsStartCmd())
	accountReportsCmd.AddCommand(newAccountReportsGetCmd())
	accountReportsCmd.AddCommand(newAccountReportsDeleteCmd())
	accountReportsCmd.AddCommand(newAccountReportsAbortCmd())
}

func newAccountReportsListCmd() *cobra.Command {
	opts := &options.AccountReportsListOptions{}

	cmd := &cobra.Command{
		Use:   "list <account-id>",
		Short: "List available report types",
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

			return runAccountReportsList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountReportsRunsCmd() *cobra.Command {
	opts := &options.AccountReportsRunsOptions{}

	cmd := &cobra.Command{
		Use:   "runs <account-id> <report-name>",
		Short: "List runs for a specific report type",
		Args:  ExactArgsWithUsage(2, "account-id", "report-name"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID
			opts.ReportName = args[1]

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountReportsRuns(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountReportsStartCmd() *cobra.Command {
	opts := &options.AccountReportsStartOptions{}

	cmd := &cobra.Command{
		Use:   "start <account-id> <report-name>",
		Short: "Start a new report run",
		Args:  ExactArgsWithUsage(2, "account-id", "report-name"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID
			opts.ReportName = args[1]

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountReportsStart(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountReportsGetCmd() *cobra.Command {
	opts := &options.AccountReportsGetRunOptions{}

	cmd := &cobra.Command{
		Use:   "get <account-id> <report-name> <run-id>",
		Short: "Get a specific report run",
		Args:  ExactArgsWithUsage(3, "account-id", "report-name", "run-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			runID, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid run ID: %s", args[2])
			}
			opts.AccountID = accountID
			opts.ReportName = args[1]
			opts.RunID = runID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountReportsGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountReportsDeleteCmd() *cobra.Command {
	opts := &options.AccountReportsDeleteRunOptions{}

	cmd := &cobra.Command{
		Use:   "delete <account-id> <report-name> <run-id>",
		Short: "Delete a report run",
		Args:  ExactArgsWithUsage(3, "account-id", "report-name", "run-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			runID, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid run ID: %s", args[2])
			}
			opts.AccountID = accountID
			opts.ReportName = args[1]
			opts.RunID = runID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountReportsDelete(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountReportsAbortCmd() *cobra.Command {
	opts := &options.AccountReportsAbortRunOptions{}

	cmd := &cobra.Command{
		Use:   "abort <account-id> <report-name> <run-id>",
		Short: "Abort a running report",
		Args:  ExactArgsWithUsage(3, "account-id", "report-name", "run-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			runID, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid run ID: %s", args[2])
			}
			opts.AccountID = accountID
			opts.ReportName = args[1]
			opts.RunID = runID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountReportsAbort(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func runAccountReportsList(ctx context.Context, client *api.Client, opts *options.AccountReportsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-reports.list", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAccountReportsService(client)
	reports, err := svc.List(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "account-reports.list", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "account-reports.list", len(reports))

	return formatEmptyOrOutput(reports, fmt.Sprintf("No reports available for account %d", opts.AccountID))
}

func runAccountReportsRuns(ctx context.Context, client *api.Client, opts *options.AccountReportsRunsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-reports.runs", map[string]interface{}{
		"account_id":  opts.AccountID,
		"report_name": opts.ReportName,
	})

	svc := api.NewAccountReportsService(client)
	runs, err := svc.ListRuns(ctx, opts.AccountID, opts.ReportName)
	if err != nil {
		logger.LogCommandError(ctx, "account-reports.runs", err, map[string]interface{}{
			"account_id":  opts.AccountID,
			"report_name": opts.ReportName,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "account-reports.runs", len(runs))

	return formatEmptyOrOutput(runs, fmt.Sprintf("No runs found for report '%s'", opts.ReportName))
}

func runAccountReportsStart(ctx context.Context, client *api.Client, opts *options.AccountReportsStartOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-reports.start", map[string]interface{}{
		"account_id":  opts.AccountID,
		"report_name": opts.ReportName,
	})

	svc := api.NewAccountReportsService(client)
	run, err := svc.Start(ctx, opts.AccountID, opts.ReportName, nil)
	if err != nil {
		logger.LogCommandError(ctx, "account-reports.start", err, map[string]interface{}{
			"account_id":  opts.AccountID,
			"report_name": opts.ReportName,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "account-reports.start", 1)

	return formatSuccessOutput(run, fmt.Sprintf("Report '%s' started (Run ID: %d)", opts.ReportName, run.ID))
}

func runAccountReportsGet(ctx context.Context, client *api.Client, opts *options.AccountReportsGetRunOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-reports.get", map[string]interface{}{
		"account_id":  opts.AccountID,
		"report_name": opts.ReportName,
		"run_id":      opts.RunID,
	})

	svc := api.NewAccountReportsService(client)
	run, err := svc.GetRun(ctx, opts.AccountID, opts.ReportName, opts.RunID)
	if err != nil {
		logger.LogCommandError(ctx, "account-reports.get", err, map[string]interface{}{
			"account_id":  opts.AccountID,
			"report_name": opts.ReportName,
			"run_id":      opts.RunID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "account-reports.get", 1)

	return formatOutput(run, func() {
		fmt.Printf("Report Run\n")
		fmt.Printf("==========\n\n")
		fmt.Printf("ID:         %d\n", run.ID)
		fmt.Printf("Report:     %s\n", run.Report)
		fmt.Printf("Status:     %s\n", run.Status)
		fmt.Printf("Progress:   %d%%\n", run.Progress)
		fmt.Printf("Created:    %s\n", run.CreatedAt)
		if run.StartedAt != "" {
			fmt.Printf("Started:    %s\n", run.StartedAt)
		}
		if run.EndedAt != "" {
			fmt.Printf("Ended:      %s\n", run.EndedAt)
		}
		if run.Attachment != nil {
			fmt.Printf("Download:   %s\n", run.Attachment.URL)
		}
		if len(run.Parameters) > 0 {
			params := make([]string, 0, len(run.Parameters))
			for k, v := range run.Parameters {
				params = append(params, fmt.Sprintf("%s=%v", k, v))
			}
			fmt.Printf("Parameters: %s\n", strings.Join(params, ", "))
		}
	})
}

func runAccountReportsDelete(ctx context.Context, client *api.Client, opts *options.AccountReportsDeleteRunOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-reports.delete", map[string]interface{}{
		"account_id":  opts.AccountID,
		"report_name": opts.ReportName,
		"run_id":      opts.RunID,
	})

	svc := api.NewAccountReportsService(client)
	err := svc.DeleteRun(ctx, opts.AccountID, opts.ReportName, opts.RunID)
	if err != nil {
		logger.LogCommandError(ctx, "account-reports.delete", err, map[string]interface{}{
			"account_id":  opts.AccountID,
			"report_name": opts.ReportName,
			"run_id":      opts.RunID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "account-reports.delete", 1)

	printInfo("Report run %d deleted\n", opts.RunID)
	return nil
}

func runAccountReportsAbort(ctx context.Context, client *api.Client, opts *options.AccountReportsAbortRunOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-reports.abort", map[string]interface{}{
		"account_id":  opts.AccountID,
		"report_name": opts.ReportName,
		"run_id":      opts.RunID,
	})

	svc := api.NewAccountReportsService(client)
	run, err := svc.AbortRun(ctx, opts.AccountID, opts.ReportName, opts.RunID)
	if err != nil {
		logger.LogCommandError(ctx, "account-reports.abort", err, map[string]interface{}{
			"account_id":  opts.AccountID,
			"report_name": opts.ReportName,
			"run_id":      opts.RunID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "account-reports.abort", 1)

	return formatSuccessOutput(run, fmt.Sprintf("Report run %d aborted", opts.RunID))
}
