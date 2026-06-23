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

// cspCmd is the root command for CSP settings
var cspCmd = &cobra.Command{
	Use:   "csp-settings",
	Short: "Manage Canvas Content Security Policy settings",
	Long: `Manage Content Security Policy (CSP) settings for Canvas accounts.

CSP settings control which domains are allowed to load content in Canvas,
providing security controls against cross-site scripting attacks.

Examples:
  canvas csp-settings get 1
  canvas csp-settings add-domain 1 --domain example.com
  canvas csp-settings remove-domain 1 --domain example.com
  canvas csp-settings lock 1 --locked`,
}

func init() {
	rootCmd.AddCommand(cspCmd)
	cspCmd.AddCommand(newCSPSettingsGetCmd())
	cspCmd.AddCommand(newCSPSettingsAddDomainCmd())
	cspCmd.AddCommand(newCSPSettingsRemoveDomainCmd())
	cspCmd.AddCommand(newCSPSettingsLockCmd())
}

func newCSPSettingsGetCmd() *cobra.Command {
	opts := &options.CSPSettingsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <account-id>",
		Short: "Get CSP settings for an account",
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

			return runCSPSettingsGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newCSPSettingsAddDomainCmd() *cobra.Command {
	opts := &options.CSPSettingsAddDomainOptions{}

	cmd := &cobra.Command{
		Use:   "add-domain <account-id>",
		Short: "Add a domain to the CSP allowlist",
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

			return runCSPSettingsAddDomain(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Domain, "domain", "", "Domain to add to CSP allowlist")
	mustMarkRequired(cmd, "domain")

	return cmd
}

func newCSPSettingsRemoveDomainCmd() *cobra.Command {
	opts := &options.CSPSettingsRemoveDomainOptions{}

	cmd := &cobra.Command{
		Use:   "remove-domain <account-id>",
		Short: "Remove a domain from the CSP allowlist",
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

			return runCSPSettingsRemoveDomain(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Domain, "domain", "", "Domain to remove from CSP allowlist")
	mustMarkRequired(cmd, "domain")

	return cmd
}

func newCSPSettingsLockCmd() *cobra.Command {
	opts := &options.CSPSettingsLockOptions{}

	cmd := &cobra.Command{
		Use:   "lock <account-id>",
		Short: "Lock or unlock CSP settings for sub-accounts",
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

			return runCSPSettingsLock(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Locked, "locked", false, "Lock CSP settings for sub-accounts")

	return cmd
}

func runCSPSettingsGet(ctx context.Context, client *api.Client, opts *options.CSPSettingsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "csp-settings.get", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewCSPSettingsService(client)
	settings, err := svc.Get(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "csp-settings.get", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "csp-settings.get", 1)

	return formatOutput(settings, func() {
		fmt.Printf("CSP Settings\n")
		fmt.Printf("============\n\n")
		fmt.Printf("Status:  %s\n", settings.Status)
		fmt.Printf("Locked:  %v\n", settings.Locked)
		if settings.LockedBy != "" {
			fmt.Printf("Locked By: %s\n", settings.LockedBy)
		}
		if len(settings.Domains) > 0 {
			fmt.Printf("\nAllowed Domains (%d):\n", len(settings.Domains))
			for _, d := range settings.Domains {
				fmt.Printf("  - %s\n", d)
			}
		} else {
			fmt.Printf("\nNo domains configured\n")
		}
	})
}

func runCSPSettingsAddDomain(ctx context.Context, client *api.Client, opts *options.CSPSettingsAddDomainOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "csp-settings.add-domain", map[string]interface{}{
		"account_id": opts.AccountID,
		"domain":     opts.Domain,
	})

	svc := api.NewCSPSettingsService(client)
	settings, err := svc.AddDomains(ctx, opts.AccountID, []string{opts.Domain})
	if err != nil {
		logger.LogCommandError(ctx, "csp-settings.add-domain", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"domain":     opts.Domain,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "csp-settings.add-domain", 1)

	return formatSuccessOutput(settings, fmt.Sprintf("Domain '%s' added to CSP allowlist", opts.Domain))
}

func runCSPSettingsRemoveDomain(ctx context.Context, client *api.Client, opts *options.CSPSettingsRemoveDomainOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "csp-settings.remove-domain", map[string]interface{}{
		"account_id": opts.AccountID,
		"domain":     opts.Domain,
	})

	svc := api.NewCSPSettingsService(client)
	settings, err := svc.RemoveDomains(ctx, opts.AccountID, []string{opts.Domain})
	if err != nil {
		logger.LogCommandError(ctx, "csp-settings.remove-domain", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"domain":     opts.Domain,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "csp-settings.remove-domain", 1)

	return formatSuccessOutput(settings, fmt.Sprintf("Domain '%s' removed from CSP allowlist", opts.Domain))
}

func runCSPSettingsLock(ctx context.Context, client *api.Client, opts *options.CSPSettingsLockOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "csp-settings.lock", map[string]interface{}{
		"account_id": opts.AccountID,
		"locked":     opts.Locked,
	})

	svc := api.NewCSPSettingsService(client)
	settings, err := svc.Lock(ctx, opts.AccountID, opts.Locked)
	if err != nil {
		logger.LogCommandError(ctx, "csp-settings.lock", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "csp-settings.lock", 1)

	action := "unlocked"
	if opts.Locked {
		action = "locked"
	}
	return formatSuccessOutput(settings, fmt.Sprintf("CSP settings %s for account %d", action, opts.AccountID))
}
