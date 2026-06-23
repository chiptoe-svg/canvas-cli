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

// accountContentMigrationsCmd is the root command for account-scoped content migrations
var accountContentMigrationsCmd = &cobra.Command{
	Use:   "account-content-migrations",
	Short: "Manage account-level content migrations",
	Long: `Manage content migrations at the account level in Canvas LMS.

Account content migrations allow admins to migrate course content across an account.

Examples:
  canvas account-content-migrations list 1
  canvas account-content-migrations get 1 42
  canvas account-content-migrations create 1 --type course_copy_importer
  canvas account-content-migrations migrators 1
  canvas account-content-migrations issues 1 42`,
}

func init() {
	rootCmd.AddCommand(accountContentMigrationsCmd)
	accountContentMigrationsCmd.AddCommand(newAccountContentMigrationsListCmd())
	accountContentMigrationsCmd.AddCommand(newAccountContentMigrationsGetCmd())
	accountContentMigrationsCmd.AddCommand(newAccountContentMigrationsCreateCmd())
	accountContentMigrationsCmd.AddCommand(newAccountContentMigrationsMigratorsCmd())
	accountContentMigrationsCmd.AddCommand(newAccountContentMigrationsIssuesCmd())
}

func newAccountContentMigrationsListCmd() *cobra.Command {
	opts := &options.AccountContentMigrationsListOptions{}

	cmd := &cobra.Command{
		Use:   "list <account-id>",
		Short: "List content migrations for an account",
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

			return runAccountContentMigrationsList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountContentMigrationsGetCmd() *cobra.Command {
	opts := &options.AccountContentMigrationsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <account-id> <migration-id>",
		Short: "Get a content migration for an account",
		Args:  ExactArgsWithUsage(2, "account-id", "migration-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID

			migrationID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid migration ID: %s", args[1])
			}
			opts.MigrationID = migrationID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountContentMigrationsGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountContentMigrationsCreateCmd() *cobra.Command {
	opts := &options.AccountContentMigrationsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <account-id>",
		Short: "Create a content migration for an account",
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

			return runAccountContentMigrationsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.MigrationType, "type", "", "Migration type (e.g. course_copy_importer, common_cartridge_importer)")
	cmd.Flags().Int64Var(&opts.SourceCourseID, "source-course-id", 0, "Source course ID (for course copy migrations)")
	mustMarkRequired(cmd, "type")

	return cmd
}

func newAccountContentMigrationsMigratorsCmd() *cobra.Command {
	opts := &options.AccountContentMigrationsMigratorsOptions{}

	cmd := &cobra.Command{
		Use:   "migrators <account-id>",
		Short: "List available migrator types for an account",
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

			return runAccountContentMigrationsMigrators(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountContentMigrationsIssuesCmd() *cobra.Command {
	opts := &options.AccountContentMigrationsIssuesOptions{}

	cmd := &cobra.Command{
		Use:   "issues <account-id> <migration-id>",
		Short: "List migration issues for an account content migration",
		Args:  ExactArgsWithUsage(2, "account-id", "migration-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID

			migrationID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid migration ID: %s", args[1])
			}
			opts.MigrationID = migrationID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountContentMigrationsIssues(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func runAccountContentMigrationsList(ctx context.Context, client *api.Client, opts *options.AccountContentMigrationsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-content-migrations.list", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	service := api.NewContentMigrationsService(client)

	migrations, err := service.ListForAccount(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "account-content-migrations.list", err, nil)
		return err
	}

	if len(migrations) == 0 {
		fmt.Printf("No content migrations found for account %d\n", opts.AccountID)
		logger.LogCommandComplete(ctx, "account-content-migrations.list", 0)
		return nil
	}

	logger.LogCommandComplete(ctx, "account-content-migrations.list", len(migrations))

	return formatOutput(migrations, func() {
		fmt.Printf("%-8s %-35s %-12s %-15s\n", "ID", "TYPE", "STATE", "STARTED")
		fmt.Println(strings.Repeat("-", 75))

		for _, m := range migrations {
			migType := m.MigrationType
			if len(migType) > 33 {
				migType = migType[:30] + "..."
			}
			fmt.Printf("%-8d %-35s %-12s %-15s\n", m.ID, migType, m.WorkflowState, m.StartedAt)
		}

		fmt.Printf("\nTotal: %d migration(s)\n", len(migrations))
	})
}

func runAccountContentMigrationsGet(ctx context.Context, client *api.Client, opts *options.AccountContentMigrationsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-content-migrations.get", map[string]interface{}{
		"account_id":   opts.AccountID,
		"migration_id": opts.MigrationID,
	})

	service := api.NewContentMigrationsService(client)

	migration, err := service.GetForAccount(ctx, opts.AccountID, opts.MigrationID)
	if err != nil {
		logger.LogCommandError(ctx, "account-content-migrations.get", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-content-migrations.get", 1)

	return formatOutput(migration, func() {
		fmt.Printf("Content Migration Details\n")
		fmt.Printf("=========================\n\n")
		fmt.Printf("ID:             %d\n", migration.ID)
		fmt.Printf("Type:           %s\n", migration.MigrationType)
		fmt.Printf("State:          %s\n", migration.WorkflowState)
		if migration.StartedAt != "" {
			fmt.Printf("Started At:     %s\n", migration.StartedAt)
		}
		if migration.FinishedAt != "" {
			fmt.Printf("Finished At:    %s\n", migration.FinishedAt)
		}
		if migration.MigrationIssuesCount > 0 {
			fmt.Printf("Issues:         %d\n", migration.MigrationIssuesCount)
		}
	})
}

func runAccountContentMigrationsCreate(ctx context.Context, client *api.Client, opts *options.AccountContentMigrationsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-content-migrations.create", map[string]interface{}{
		"account_id": opts.AccountID,
		"type":       opts.MigrationType,
	})

	service := api.NewContentMigrationsService(client)

	params := &api.CreateContentMigrationParams{
		MigrationType: opts.MigrationType,
	}

	if opts.SourceCourseID > 0 {
		params.SourceCourseID = &opts.SourceCourseID
	}

	migration, err := service.CreateForAccount(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "account-content-migrations.create", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-content-migrations.create", 1)

	return formatOutput(migration, func() {
		fmt.Printf("Content migration created successfully\n\n")
		fmt.Printf("ID:    %d\n", migration.ID)
		fmt.Printf("Type:  %s\n", migration.MigrationType)
		fmt.Printf("State: %s\n", migration.WorkflowState)
	})
}

func runAccountContentMigrationsMigrators(ctx context.Context, client *api.Client, opts *options.AccountContentMigrationsMigratorsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-content-migrations.migrators", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	service := api.NewContentMigrationsService(client)

	migrators, err := service.GetMigratorsForAccount(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "account-content-migrations.migrators", err, nil)
		return err
	}

	if len(migrators) == 0 {
		fmt.Println("No migrators available")
		logger.LogCommandComplete(ctx, "account-content-migrations.migrators", 0)
		return nil
	}

	logger.LogCommandComplete(ctx, "account-content-migrations.migrators", len(migrators))

	return formatOutput(migrators, func() {
		fmt.Printf("%-35s %-40s %-15s\n", "TYPE", "NAME", "NEEDS_FILE")
		fmt.Println(strings.Repeat("-", 95))

		for _, m := range migrators {
			fmt.Printf("%-35s %-40s %-15v\n", m.Type, m.Name, m.RequiresFileUpload)
		}

		fmt.Printf("\nTotal: %d migrator(s)\n", len(migrators))
	})
}

func runAccountContentMigrationsIssues(ctx context.Context, client *api.Client, opts *options.AccountContentMigrationsIssuesOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-content-migrations.issues", map[string]interface{}{
		"account_id":   opts.AccountID,
		"migration_id": opts.MigrationID,
	})

	service := api.NewContentMigrationsService(client)

	issues, err := service.GetMigrationIssuesForAccount(ctx, opts.AccountID, opts.MigrationID)
	if err != nil {
		logger.LogCommandError(ctx, "account-content-migrations.issues", err, nil)
		return err
	}

	if len(issues) == 0 {
		fmt.Printf("No issues found for migration %d\n", opts.MigrationID)
		logger.LogCommandComplete(ctx, "account-content-migrations.issues", 0)
		return nil
	}

	logger.LogCommandComplete(ctx, "account-content-migrations.issues", len(issues))

	return formatOutput(issues, func() {
		fmt.Printf("%-8s %-12s %-12s %-40s\n", "ID", "TYPE", "STATE", "DESCRIPTION")
		fmt.Println(strings.Repeat("-", 78))

		for _, issue := range issues {
			desc := issue.Description
			if len(desc) > 38 {
				desc = desc[:35] + "..."
			}
			fmt.Printf("%-8d %-12s %-12s %-40s\n", issue.ID, issue.IssueType, issue.WorkflowState, desc)
		}

		fmt.Printf("\nTotal: %d issue(s)\n", len(issues))
	})
}
