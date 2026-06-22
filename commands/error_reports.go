package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// errorReportsCmd represents the error-reports command group.
var errorReportsCmd = &cobra.Command{
	Use:   "error-reports",
	Short: "Submit Canvas error reports",
	Long: `Submit error reports to Canvas administrators.

Examples:
  canvas error-reports create --subject "Page failed to load" --comments "The grades page throws a 500."`,
}

func init() {
	rootCmd.AddCommand(errorReportsCmd)
	errorReportsCmd.AddCommand(newErrorReportCreateCmd())
}

func newErrorReportCreateCmd() *cobra.Command {
	opts := &options.ErrorReportCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Submit an error report",
		Long: `Submit an error report to Canvas administrators.

Examples:
  canvas error-reports create --subject "Login broken" --comments "Cannot log in since yesterday"
  canvas error-reports create --subject "Bug" --email user@example.com --severity medium`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runErrorReportCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Subject, "subject", "", "Subject of the error report")
	cmd.Flags().StringVar(&opts.Comments, "comments", "", "Detailed description of the error")
	cmd.Flags().StringVar(&opts.URL, "url", "", "URL where the error occurred")
	cmd.Flags().StringVar(&opts.Email, "email", "", "Contact email address")
	cmd.Flags().StringVar(&opts.Severity, "severity", "", "Severity: just_a_comment, not_urgent, workaround_possible, blocks_what_i_need_to_do, extreme_critical_emergency")
	mustMarkRequired(cmd, "subject")

	return cmd
}

func runErrorReportCreate(ctx context.Context, client *api.Client, opts *options.ErrorReportCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "error-reports.create", map[string]interface{}{
		"subject":  opts.Subject,
		"severity": opts.Severity,
	})

	svc := api.NewErrorReportsService(client)

	report := &api.ErrorReport{
		Subject:               opts.Subject,
		Comments:              opts.Comments,
		URL:                   opts.URL,
		Email:                 opts.Email,
		UserPerceivedSeverity: opts.Severity,
	}

	result, err := svc.Create(ctx, report)
	if err != nil {
		logger.LogCommandError(ctx, "error-reports.create", err, nil)
		return fmt.Errorf("failed to create error report: %w", err)
	}

	logger.LogCommandComplete(ctx, "error-reports.create", 1)
	return formatSuccessOutput(result, "Error report submitted successfully.")
}
