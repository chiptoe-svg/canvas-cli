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

var developerKeysCmd = &cobra.Command{
	Use:   "developer-keys",
	Short: "Manage Canvas developer keys",
	Long: `Manage developer keys and their account bindings in Canvas.

Developer keys are used for OAuth2 integrations and LTI tools.

Examples:
  canvas developer-keys list 1
  canvas developer-keys create 1 --name "My Integration" --email dev@example.com
  canvas developer-keys bind 1 10 --workflow-state on`,
}

func init() {
	rootCmd.AddCommand(developerKeysCmd)
	developerKeysCmd.AddCommand(newDeveloperKeysListCmd())
	developerKeysCmd.AddCommand(newDeveloperKeysCreateCmd())
	developerKeysCmd.AddCommand(newDeveloperKeysBindCmd())
}

func newDeveloperKeysListCmd() *cobra.Command {
	opts := &options.DeveloperKeysListOptions{}

	cmd := &cobra.Command{
		Use:   "list <account-id>",
		Short: "List developer keys for an account",
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

			return runDeveloperKeysList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newDeveloperKeysCreateCmd() *cobra.Command {
	opts := &options.DeveloperKeysCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <account-id>",
		Short: "Create a developer key in an account",
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

			return runDeveloperKeysCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Developer key name")
	cmd.Flags().StringVar(&opts.Email, "email", "", "Contact email")
	cmd.Flags().StringVar(&opts.RedirectURI, "redirect-uri", "", "OAuth redirect URI")
	cmd.Flags().StringVar(&opts.Notes, "notes", "", "Notes about the key")

	return cmd
}

func newDeveloperKeysBindCmd() *cobra.Command {
	opts := &options.DeveloperKeysBindOptions{}

	cmd := &cobra.Command{
		Use:   "bind <account-id> <developer-key-id>",
		Short: "Create an account binding for a developer key",
		Args:  ExactArgsWithUsage(2, "account-id", "developer-key-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			developerKeyID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid developer key ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.DeveloperKeyID = developerKeyID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runDeveloperKeysBind(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.WorkflowState, "workflow-state", "", `Binding state: "on", "off", or "allow" (required)`)
	mustMarkRequired(cmd, "workflow-state")

	return cmd
}

func runDeveloperKeysList(ctx context.Context, client *api.Client, opts *options.DeveloperKeysListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "developer-keys.list", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	service := api.NewDeveloperKeysService(client)

	keys, err := service.List(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "developer-keys.list", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "developer-keys.list", len(keys))

	return formatEmptyOrOutput(keys, fmt.Sprintf("No developer keys found for account %d", opts.AccountID))
}

func runDeveloperKeysCreate(ctx context.Context, client *api.Client, opts *options.DeveloperKeysCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "developer-keys.create", map[string]interface{}{
		"account_id": opts.AccountID,
		"name":       opts.Name,
	})

	service := api.NewDeveloperKeysService(client)

	params := &api.DeveloperKeyParams{
		DeveloperKey: api.DeveloperKeyFields{
			Name:        opts.Name,
			Email:       opts.Email,
			RedirectURI: opts.RedirectURI,
			Notes:       opts.Notes,
		},
	}

	key, err := service.Create(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "developer-keys.create", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "developer-keys.create", 1)

	return formatSuccessOutput(key, fmt.Sprintf("Developer key '%s' created with ID %d", key.Name, key.ID))
}

func runDeveloperKeysBind(ctx context.Context, client *api.Client, opts *options.DeveloperKeysBindOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "developer-keys.bind", map[string]interface{}{
		"account_id":       opts.AccountID,
		"developer_key_id": opts.DeveloperKeyID,
		"workflow_state":   opts.WorkflowState,
	})

	service := api.NewDeveloperKeysService(client)

	params := &api.DeveloperKeyBindingParams{
		DeveloperKeyAccountBinding: api.DeveloperKeyBindingFields{
			WorkflowState: opts.WorkflowState,
		},
	}

	binding, err := service.CreateBinding(ctx, opts.AccountID, opts.DeveloperKeyID, params)
	if err != nil {
		logger.LogCommandError(ctx, "developer-keys.bind", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "developer-keys.bind", 1)

	return formatSuccessOutput(binding, fmt.Sprintf("Developer key %d bound with state '%s'", opts.DeveloperKeyID, binding.WorkflowState))
}
