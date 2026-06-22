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

// accountExtToolsFavoritesCmd is the root command for account external tool favorites
var accountExtToolsFavoritesCmd = &cobra.Command{
	Use:   "account-ext-tools-favorites",
	Short: "Manage account external tool (LTI) favorites",
	Long: `Manage RCE and top-nav favorites for external tools (LTI apps) in a Canvas account.

Examples:
  canvas account-ext-tools-favorites add-rce 1 5
  canvas account-ext-tools-favorites remove-rce 1 5
  canvas account-ext-tools-favorites add-topnav 1 7
  canvas account-ext-tools-favorites remove-topnav 1 7`,
}

func init() {
	rootCmd.AddCommand(accountExtToolsFavoritesCmd)
	accountExtToolsFavoritesCmd.AddCommand(newAddRCEFavoriteCmd())
	accountExtToolsFavoritesCmd.AddCommand(newRemoveRCEFavoriteCmd())
	accountExtToolsFavoritesCmd.AddCommand(newAddTopNavFavoriteCmd())
	accountExtToolsFavoritesCmd.AddCommand(newRemoveTopNavFavoriteCmd())
}

func newAddRCEFavoriteCmd() *cobra.Command {
	opts := &options.AccountExtToolsFavoritesOptions{}

	cmd := &cobra.Command{
		Use:   "add-rce <account-id> <tool-id>",
		Short: "Add an external tool as an RCE favorite for an account",
		Args:  ExactArgsWithUsage(2, "account-id", "tool-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID

			toolID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid tool ID: %s", args[1])
			}
			opts.ToolID = toolID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAddRCEFavorite(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newRemoveRCEFavoriteCmd() *cobra.Command {
	opts := &options.AccountExtToolsFavoritesOptions{}

	cmd := &cobra.Command{
		Use:   "remove-rce <account-id> <tool-id>",
		Short: "Remove an external tool from the RCE favorites for an account",
		Args:  ExactArgsWithUsage(2, "account-id", "tool-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID

			toolID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid tool ID: %s", args[1])
			}
			opts.ToolID = toolID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runRemoveRCEFavorite(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAddTopNavFavoriteCmd() *cobra.Command {
	opts := &options.AccountExtToolsFavoritesOptions{}

	cmd := &cobra.Command{
		Use:   "add-topnav <account-id> <tool-id>",
		Short: "Add an external tool as a top-nav favorite for an account",
		Args:  ExactArgsWithUsage(2, "account-id", "tool-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID

			toolID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid tool ID: %s", args[1])
			}
			opts.ToolID = toolID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAddTopNavFavorite(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newRemoveTopNavFavoriteCmd() *cobra.Command {
	opts := &options.AccountExtToolsFavoritesOptions{}

	cmd := &cobra.Command{
		Use:   "remove-topnav <account-id> <tool-id>",
		Short: "Remove an external tool from the top-nav favorites for an account",
		Args:  ExactArgsWithUsage(2, "account-id", "tool-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID

			toolID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid tool ID: %s", args[1])
			}
			opts.ToolID = toolID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runRemoveTopNavFavorite(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func runAddRCEFavorite(ctx context.Context, client *api.Client, opts *options.AccountExtToolsFavoritesOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-ext-tools-favorites.add-rce", map[string]interface{}{
		"account_id": opts.AccountID,
		"tool_id":    opts.ToolID,
	})

	service := api.NewAccountExternalToolsService(client)

	if err := service.AddRCEFavorite(ctx, opts.AccountID, opts.ToolID); err != nil {
		logger.LogCommandError(ctx, "account-ext-tools-favorites.add-rce", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-ext-tools-favorites.add-rce", 1)
	fmt.Printf("Tool %d added to RCE favorites for account %d\n", opts.ToolID, opts.AccountID)
	return nil
}

func runRemoveRCEFavorite(ctx context.Context, client *api.Client, opts *options.AccountExtToolsFavoritesOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-ext-tools-favorites.remove-rce", map[string]interface{}{
		"account_id": opts.AccountID,
		"tool_id":    opts.ToolID,
	})

	service := api.NewAccountExternalToolsService(client)

	if err := service.RemoveRCEFavorite(ctx, opts.AccountID, opts.ToolID); err != nil {
		logger.LogCommandError(ctx, "account-ext-tools-favorites.remove-rce", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-ext-tools-favorites.remove-rce", 1)
	fmt.Printf("Tool %d removed from RCE favorites for account %d\n", opts.ToolID, opts.AccountID)
	return nil
}

func runAddTopNavFavorite(ctx context.Context, client *api.Client, opts *options.AccountExtToolsFavoritesOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-ext-tools-favorites.add-topnav", map[string]interface{}{
		"account_id": opts.AccountID,
		"tool_id":    opts.ToolID,
	})

	service := api.NewAccountExternalToolsService(client)

	if err := service.AddTopNavFavorite(ctx, opts.AccountID, opts.ToolID); err != nil {
		logger.LogCommandError(ctx, "account-ext-tools-favorites.add-topnav", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-ext-tools-favorites.add-topnav", 1)
	fmt.Printf("Tool %d added to top-nav favorites for account %d\n", opts.ToolID, opts.AccountID)
	return nil
}

func runRemoveTopNavFavorite(ctx context.Context, client *api.Client, opts *options.AccountExtToolsFavoritesOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-ext-tools-favorites.remove-topnav", map[string]interface{}{
		"account_id": opts.AccountID,
		"tool_id":    opts.ToolID,
	})

	service := api.NewAccountExternalToolsService(client)

	if err := service.RemoveTopNavFavorite(ctx, opts.AccountID, opts.ToolID); err != nil {
		logger.LogCommandError(ctx, "account-ext-tools-favorites.remove-topnav", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-ext-tools-favorites.remove-topnav", 1)
	fmt.Printf("Tool %d removed from top-nav favorites for account %d\n", opts.ToolID, opts.AccountID)
	return nil
}
