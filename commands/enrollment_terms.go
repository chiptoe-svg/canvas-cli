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

var enrollmentTermsCmd = &cobra.Command{
	Use:   "enrollment-terms",
	Short: "Manage Canvas enrollment terms",
	Long: `Manage enrollment terms within a Canvas account.

Enrollment terms define academic periods (semesters, quarters, etc.) and
control when courses are active.

Examples:
  canvas enrollment-terms list 1
  canvas enrollment-terms get 1 42
  canvas enrollment-terms create 1 --name "Fall 2025"
  canvas enrollment-terms update 1 42 --name "Fall 2025 Updated"
  canvas enrollment-terms delete 1 42`,
}

func init() {
	rootCmd.AddCommand(enrollmentTermsCmd)
	enrollmentTermsCmd.AddCommand(newEnrollmentTermsListCmd())
	enrollmentTermsCmd.AddCommand(newEnrollmentTermsGetCmd())
	enrollmentTermsCmd.AddCommand(newEnrollmentTermsCreateCmd())
	enrollmentTermsCmd.AddCommand(newEnrollmentTermsUpdateCmd())
	enrollmentTermsCmd.AddCommand(newEnrollmentTermsDeleteCmd())
}

func newEnrollmentTermsListCmd() *cobra.Command {
	opts := &options.EnrollmentTermsListOptions{}

	cmd := &cobra.Command{
		Use:   "list <account-id>",
		Short: "List enrollment terms for an account",
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

			return runEnrollmentTermsList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newEnrollmentTermsGetCmd() *cobra.Command {
	opts := &options.EnrollmentTermsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <account-id> <term-id>",
		Short: "Get a single enrollment term",
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

			return runEnrollmentTermsGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newEnrollmentTermsCreateCmd() *cobra.Command {
	opts := &options.EnrollmentTermsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <account-id>",
		Short: "Create an enrollment term",
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

			return runEnrollmentTermsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Name of the enrollment term (required)")
	cmd.Flags().StringVar(&opts.StartAt, "start-at", "", "Start date (ISO 8601)")
	cmd.Flags().StringVar(&opts.EndAt, "end-at", "", "End date (ISO 8601)")
	cmd.Flags().StringVar(&opts.SISTermID, "sis-term-id", "", "SIS term ID")
	mustMarkRequired(cmd, "name")

	return cmd
}

func newEnrollmentTermsUpdateCmd() *cobra.Command {
	opts := &options.EnrollmentTermsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <account-id> <term-id>",
		Short: "Update an enrollment term",
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

			return runEnrollmentTermsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "New name for the enrollment term")
	cmd.Flags().StringVar(&opts.StartAt, "start-at", "", "Start date (ISO 8601)")
	cmd.Flags().StringVar(&opts.EndAt, "end-at", "", "End date (ISO 8601)")

	return cmd
}

func newEnrollmentTermsDeleteCmd() *cobra.Command {
	opts := &options.EnrollmentTermsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <account-id> <term-id>",
		Short: "Delete an enrollment term",
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

			return runEnrollmentTermsDelete(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func runEnrollmentTermsList(ctx context.Context, client *api.Client, opts *options.EnrollmentTermsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "enrollment-terms.list", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	service := api.NewEnrollmentTermsService(client)

	terms, err := service.List(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "enrollment-terms.list", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "enrollment-terms.list", len(terms))

	return formatEmptyOrOutput(terms, fmt.Sprintf("No enrollment terms found for account %d", opts.AccountID))
}

func runEnrollmentTermsGet(ctx context.Context, client *api.Client, opts *options.EnrollmentTermsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "enrollment-terms.get", map[string]interface{}{
		"account_id": opts.AccountID,
		"term_id":    opts.TermID,
	})

	service := api.NewEnrollmentTermsService(client)

	term, err := service.Get(ctx, opts.AccountID, opts.TermID)
	if err != nil {
		logger.LogCommandError(ctx, "enrollment-terms.get", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "enrollment-terms.get", 1)

	return formatOutput(term, nil)
}

func runEnrollmentTermsCreate(ctx context.Context, client *api.Client, opts *options.EnrollmentTermsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "enrollment-terms.create", map[string]interface{}{
		"account_id": opts.AccountID,
		"name":       opts.Name,
	})

	service := api.NewEnrollmentTermsService(client)

	params := &api.EnrollmentTermParams{
		EnrollmentTerm: api.EnrollmentTermFields{
			Name:      opts.Name,
			StartAt:   opts.StartAt,
			EndAt:     opts.EndAt,
			SISTermID: opts.SISTermID,
		},
	}

	term, err := service.Create(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "enrollment-terms.create", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "enrollment-terms.create", 1)

	return formatSuccessOutput(term, fmt.Sprintf("Enrollment term '%s' created with ID %d", term.Name, term.ID))
}

func runEnrollmentTermsUpdate(ctx context.Context, client *api.Client, opts *options.EnrollmentTermsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "enrollment-terms.update", map[string]interface{}{
		"account_id": opts.AccountID,
		"term_id":    opts.TermID,
	})

	service := api.NewEnrollmentTermsService(client)

	params := &api.EnrollmentTermParams{
		EnrollmentTerm: api.EnrollmentTermFields{
			Name:    opts.Name,
			StartAt: opts.StartAt,
			EndAt:   opts.EndAt,
		},
	}

	term, err := service.Update(ctx, opts.AccountID, opts.TermID, params)
	if err != nil {
		logger.LogCommandError(ctx, "enrollment-terms.update", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "enrollment-terms.update", 1)

	return formatSuccessOutput(term, fmt.Sprintf("Enrollment term %d updated", term.ID))
}

func runEnrollmentTermsDelete(ctx context.Context, client *api.Client, opts *options.EnrollmentTermsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "enrollment-terms.delete", map[string]interface{}{
		"account_id": opts.AccountID,
		"term_id":    opts.TermID,
	})

	service := api.NewEnrollmentTermsService(client)

	if err := service.Delete(ctx, opts.AccountID, opts.TermID); err != nil {
		logger.LogCommandError(ctx, "enrollment-terms.delete", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "enrollment-terms.delete", 1)

	printInfo("Enrollment term %d deleted\n", opts.TermID)
	return nil
}
