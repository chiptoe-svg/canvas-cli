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

// eportfoliosCmd represents the eportfolios command group.
var eportfoliosCmd = &cobra.Command{
	Use:   "eportfolios",
	Short: "Manage Canvas ePortfolios",
	Long: `Manage Canvas ePortfolios (view, delete, moderate, list pages).

Examples:
  canvas eportfolios list --user-id 42
  canvas eportfolios get 10
  canvas eportfolios delete 10
  canvas eportfolios pages 10`,
}

func init() {
	rootCmd.AddCommand(eportfoliosCmd)
	eportfoliosCmd.AddCommand(newEportfoliosListCmd())
	eportfoliosCmd.AddCommand(newEportfolioGetCmd())
	eportfoliosCmd.AddCommand(newEportfolioDeleteCmd())
	eportfoliosCmd.AddCommand(newEportfolioPagesCmd())
}

func newEportfoliosListCmd() *cobra.Command {
	opts := &options.EportfoliosListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List ePortfolios for a user",
		Long: `List all ePortfolios belonging to a user.

Examples:
  canvas eportfolios list --user-id 42`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runEportfoliosList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	mustMarkRequired(cmd, "user-id")

	return cmd
}

func newEportfolioGetCmd() *cobra.Command {
	opts := &options.EportfolioGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get an ePortfolio by ID",
		Long: `Retrieve details of a specific ePortfolio.

Examples:
  canvas eportfolios get 10`,
		Args: ExactArgsWithUsage(1, "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid eportfolio ID: %s", args[0])
			}
			opts.ID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runEportfolioGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newEportfolioDeleteCmd() *cobra.Command {
	opts := &options.EportfolioDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an ePortfolio",
		Long: `Delete an ePortfolio by its ID.

Examples:
  canvas eportfolios delete 10
  canvas eportfolios delete 10 --force`,
		Args: ExactArgsWithUsage(1, "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid eportfolio ID: %s", args[0])
			}
			opts.ID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runEportfolioDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func newEportfolioPagesCmd() *cobra.Command {
	opts := &options.EportfolioPagesListOptions{}

	cmd := &cobra.Command{
		Use:   "pages <eportfolio-id>",
		Short: "List pages in an ePortfolio",
		Long: `List all pages within a specific ePortfolio.

Examples:
  canvas eportfolios pages 10`,
		Args: ExactArgsWithUsage(1, "eportfolio-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid eportfolio ID: %s", args[0])
			}
			opts.EportfolioID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runEportfolioPages(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func runEportfoliosList(ctx context.Context, client *api.Client, opts *options.EportfoliosListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "eportfolios.list", map[string]interface{}{
		"user_id": opts.UserID,
	})

	svc := api.NewEportfoliosService(client)

	eps, err := svc.ListForUser(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "eportfolios.list", err, nil)
		return fmt.Errorf("failed to list ePortfolios: %w", err)
	}

	printVerbose("Found %d ePortfolios:\n\n", len(eps))
	logger.LogCommandComplete(ctx, "eportfolios.list", len(eps))
	return formatEmptyOrOutput(eps, "No ePortfolios found")
}

func runEportfolioGet(ctx context.Context, client *api.Client, opts *options.EportfolioGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "eportfolios.get", map[string]interface{}{
		"id": opts.ID,
	})

	svc := api.NewEportfoliosService(client)

	ep, err := svc.Get(ctx, opts.ID)
	if err != nil {
		logger.LogCommandError(ctx, "eportfolios.get", err, nil)
		return fmt.Errorf("failed to get ePortfolio: %w", err)
	}

	logger.LogCommandComplete(ctx, "eportfolios.get", 1)
	return formatOutput(ep, nil)
}

func runEportfolioDelete(ctx context.Context, client *api.Client, opts *options.EportfolioDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "eportfolios.delete", map[string]interface{}{
		"id": opts.ID,
	})

	confirmed, err := confirmDelete("ePortfolio", opts.ID, opts.Force)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Delete cancelled")
		return nil
	}

	svc := api.NewEportfoliosService(client)

	if err := svc.Delete(ctx, opts.ID); err != nil {
		logger.LogCommandError(ctx, "eportfolios.delete", err, nil)
		return fmt.Errorf("failed to delete ePortfolio: %w", err)
	}

	printInfo("ePortfolio %d deleted successfully\n", opts.ID)
	logger.LogCommandComplete(ctx, "eportfolios.delete", 1)
	return nil
}

func runEportfolioPages(ctx context.Context, client *api.Client, opts *options.EportfolioPagesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "eportfolios.pages", map[string]interface{}{
		"eportfolio_id": opts.EportfolioID,
	})

	svc := api.NewEportfoliosService(client)

	pages, err := svc.ListPages(ctx, opts.EportfolioID)
	if err != nil {
		logger.LogCommandError(ctx, "eportfolios.pages", err, nil)
		return fmt.Errorf("failed to list ePortfolio pages: %w", err)
	}

	printVerbose("Found %d pages:\n\n", len(pages))
	logger.LogCommandComplete(ctx, "eportfolios.pages", len(pages))
	return formatEmptyOrOutput(pages, "No pages found")
}
