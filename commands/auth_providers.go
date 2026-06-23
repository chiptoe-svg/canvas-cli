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

// authProvidersCmd is the root command for authentication providers
var authProvidersCmd = &cobra.Command{
	Use:   "auth-providers",
	Short: "Manage Canvas authentication providers",
	Long: `Manage authentication providers for Canvas accounts.

Authentication providers control how users log in to Canvas,
supporting SAML, LDAP, OAuth, and other SSO methods.

Examples:
  canvas auth-providers list 1
  canvas auth-providers get 1 10
  canvas auth-providers create 1 --auth-type saml
  canvas auth-providers delete 1 10
  canvas auth-providers sso-settings 1`,
}

func init() {
	rootCmd.AddCommand(authProvidersCmd)
	authProvidersCmd.AddCommand(newAuthProvidersListCmd())
	authProvidersCmd.AddCommand(newAuthProvidersGetCmd())
	authProvidersCmd.AddCommand(newAuthProvidersCreateCmd())
	authProvidersCmd.AddCommand(newAuthProvidersDeleteCmd())
	authProvidersCmd.AddCommand(newAuthProvidersRestoreCmd())
	authProvidersCmd.AddCommand(newAuthProvidersForcePasswordResetCmd())
	authProvidersCmd.AddCommand(newAuthProvidersSSOSettingsCmd())
}

func newAuthProvidersListCmd() *cobra.Command {
	opts := &options.AuthProvidersListOptions{}

	cmd := &cobra.Command{
		Use:   "list <account-id>",
		Short: "List authentication providers for an account",
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

			return runAuthProvidersList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAuthProvidersGetCmd() *cobra.Command {
	opts := &options.AuthProvidersGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <account-id> <id>",
		Short: "Get a specific authentication provider",
		Args:  ExactArgsWithUsage(2, "account-id", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			providerID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid provider ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.ProviderID = providerID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAuthProvidersGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAuthProvidersCreateCmd() *cobra.Command {
	opts := &options.AuthProvidersCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <account-id>",
		Short: "Create a new authentication provider",
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

			return runAuthProvidersCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.AuthType, "auth-type", "", "Authentication type (saml, ldap, google, etc.)")
	cmd.Flags().StringVar(&opts.ClientID, "client-id", "", "OAuth client ID")
	cmd.Flags().StringVar(&opts.ClientSecret, "client-secret", "", "OAuth client secret")
	cmd.Flags().StringVar(&opts.LoginAttribute, "login-attribute", "", "Login attribute")
	mustMarkRequired(cmd, "auth-type")

	return cmd
}

func newAuthProvidersDeleteCmd() *cobra.Command {
	opts := &options.AuthProvidersDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <account-id> <id>",
		Short: "Delete an authentication provider",
		Args:  ExactArgsWithUsage(2, "account-id", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			providerID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid provider ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.ProviderID = providerID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAuthProvidersDelete(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAuthProvidersRestoreCmd() *cobra.Command {
	opts := &options.AuthProvidersRestoreOptions{}

	cmd := &cobra.Command{
		Use:   "restore <account-id> <id>",
		Short: "Restore a deleted authentication provider",
		Args:  ExactArgsWithUsage(2, "account-id", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			providerID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid provider ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.ProviderID = providerID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAuthProvidersRestore(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAuthProvidersForcePasswordResetCmd() *cobra.Command {
	opts := &options.AuthProvidersForcePasswordResetOptions{}

	cmd := &cobra.Command{
		Use:   "force-password-reset <account-id>",
		Short: "Force password reset for all users",
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

			return runAuthProvidersForcePasswordReset(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAuthProvidersSSOSettingsCmd() *cobra.Command {
	opts := &options.AuthProvidersSSOSettingsOptions{}

	cmd := &cobra.Command{
		Use:   "sso-settings <account-id>",
		Short: "Get SSO settings for an account",
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

			return runAuthProvidersSSOSettings(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func runAuthProvidersList(ctx context.Context, client *api.Client, opts *options.AuthProvidersListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "auth-providers.list", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAuthProvidersService(client)
	providers, err := svc.List(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "auth-providers.list", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "auth-providers.list", len(providers))

	return formatEmptyOrOutput(providers, fmt.Sprintf("No authentication providers found for account %d", opts.AccountID))
}

func runAuthProvidersGet(ctx context.Context, client *api.Client, opts *options.AuthProvidersGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "auth-providers.get", map[string]interface{}{
		"account_id":  opts.AccountID,
		"provider_id": opts.ProviderID,
	})

	svc := api.NewAuthProvidersService(client)
	provider, err := svc.Get(ctx, opts.AccountID, opts.ProviderID)
	if err != nil {
		logger.LogCommandError(ctx, "auth-providers.get", err, map[string]interface{}{
			"account_id":  opts.AccountID,
			"provider_id": opts.ProviderID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "auth-providers.get", 1)

	return formatOutput(provider, func() {
		fmt.Printf("Authentication Provider\n")
		fmt.Printf("=======================\n\n")
		fmt.Printf("ID:             %d\n", provider.ID)
		fmt.Printf("Auth Type:      %s\n", provider.AuthType)
		fmt.Printf("Position:       %d\n", provider.Position)
		fmt.Printf("Workflow State: %s\n", provider.WorkflowState)
		fmt.Printf("JIT Provision:  %v\n", provider.JITProvisioning)
	})
}

func runAuthProvidersCreate(ctx context.Context, client *api.Client, opts *options.AuthProvidersCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "auth-providers.create", map[string]interface{}{
		"account_id": opts.AccountID,
		"auth_type":  opts.AuthType,
	})

	svc := api.NewAuthProvidersService(client)
	params := &api.AuthProviderCreateParams{
		AuthType:       opts.AuthType,
		ClientID:       opts.ClientID,
		ClientSecret:   opts.ClientSecret,
		LoginAttribute: opts.LoginAttribute,
	}

	provider, err := svc.Create(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "auth-providers.create", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "auth-providers.create", 1)

	return formatSuccessOutput(provider, fmt.Sprintf("Authentication provider created (ID: %d)", provider.ID))
}

func runAuthProvidersDelete(ctx context.Context, client *api.Client, opts *options.AuthProvidersDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "auth-providers.delete", map[string]interface{}{
		"account_id":  opts.AccountID,
		"provider_id": opts.ProviderID,
	})

	svc := api.NewAuthProvidersService(client)
	err := svc.Delete(ctx, opts.AccountID, opts.ProviderID)
	if err != nil {
		logger.LogCommandError(ctx, "auth-providers.delete", err, map[string]interface{}{
			"account_id":  opts.AccountID,
			"provider_id": opts.ProviderID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "auth-providers.delete", 1)

	printInfo("Authentication provider %d deleted\n", opts.ProviderID)
	return nil
}

func runAuthProvidersRestore(ctx context.Context, client *api.Client, opts *options.AuthProvidersRestoreOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "auth-providers.restore", map[string]interface{}{
		"account_id":  opts.AccountID,
		"provider_id": opts.ProviderID,
	})

	svc := api.NewAuthProvidersService(client)
	provider, err := svc.Restore(ctx, opts.AccountID, opts.ProviderID)
	if err != nil {
		logger.LogCommandError(ctx, "auth-providers.restore", err, map[string]interface{}{
			"account_id":  opts.AccountID,
			"provider_id": opts.ProviderID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "auth-providers.restore", 1)

	return formatSuccessOutput(provider, fmt.Sprintf("Authentication provider %d restored", provider.ID))
}

func runAuthProvidersForcePasswordReset(ctx context.Context, client *api.Client, opts *options.AuthProvidersForcePasswordResetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "auth-providers.force-password-reset", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAuthProvidersService(client)
	err := svc.ForcePasswordReset(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "auth-providers.force-password-reset", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "auth-providers.force-password-reset", 1)

	printInfo("Password reset has been forced for all users in account %d\n", opts.AccountID)
	return nil
}

func runAuthProvidersSSOSettings(ctx context.Context, client *api.Client, opts *options.AuthProvidersSSOSettingsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "auth-providers.sso-settings", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAuthProvidersService(client)
	settings, err := svc.GetSSOSettings(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "auth-providers.sso-settings", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "auth-providers.sso-settings", 1)

	return formatOutput(settings, func() {
		fmt.Printf("SSO Settings\n")
		fmt.Printf("============\n\n")
		if settings.LoginHandleName != "" {
			fmt.Printf("Login Handle Name:  %s\n", settings.LoginHandleName)
		}
		if settings.ChangePasswordURL != "" {
			fmt.Printf("Change Password URL: %s\n", settings.ChangePasswordURL)
		}
		if settings.AuthDiscoveryURL != "" {
			fmt.Printf("Auth Discovery URL:  %s\n", settings.AuthDiscoveryURL)
		}
		if settings.UnknownUserURL != "" {
			fmt.Printf("Unknown User URL:    %s\n", settings.UnknownUserURL)
		}
	})
}
