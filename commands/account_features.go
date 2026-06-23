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

// accountFeaturesCmd represents the account-features command group
var accountFeaturesCmd = &cobra.Command{
	Use:   "account-features",
	Short: "Manage Canvas account feature flags and settings",
	Long: `Manage Canvas account feature flags, settings, and permissions.

Feature flags control which Canvas features are enabled for an account.
Administrators can enable, disable, or allow sub-accounts to control features.

Examples:
  canvas account-features list 1
  canvas account-features list-enabled 1
  canvas account-features get-flag 1 analytics_2
  canvas account-features set-flag 1 analytics_2 --state on
  canvas account-features settings 1
  canvas account-features permissions 1 --permissions manage_courses,manage_users`,
}

func init() {
	rootCmd.AddCommand(accountFeaturesCmd)
	accountFeaturesCmd.AddCommand(newAccountFeaturesListCmd())
	accountFeaturesCmd.AddCommand(newAccountFeaturesListEnabledCmd())
	accountFeaturesCmd.AddCommand(newAccountFeaturesGetFlagCmd())
	accountFeaturesCmd.AddCommand(newAccountFeaturesSetFlagCmd())
	accountFeaturesCmd.AddCommand(newAccountFeaturesDeleteFlagCmd())
	accountFeaturesCmd.AddCommand(newAccountFeaturesSettingsCmd())
	accountFeaturesCmd.AddCommand(newAccountFeaturesPermissionsCmd())
}

func newAccountFeaturesListCmd() *cobra.Command {
	opts := &options.AccountFeaturesListOptions{}

	cmd := &cobra.Command{
		Use:   "list <account-id>",
		Short: "List all features for an account",
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

			return runAccountFeaturesList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountFeaturesListEnabledCmd() *cobra.Command {
	opts := &options.AccountFeaturesListEnabledOptions{}

	cmd := &cobra.Command{
		Use:   "list-enabled <account-id>",
		Short: "List enabled features for an account",
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

			return runAccountFeaturesListEnabled(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountFeaturesGetFlagCmd() *cobra.Command {
	opts := &options.AccountFeaturesGetFlagOptions{}

	cmd := &cobra.Command{
		Use:   "get-flag <account-id> <feature>",
		Short: "Get the feature flag for a specific feature",
		Args:  ExactArgsWithUsage(2, "account-id", "feature"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID
			opts.Feature = args[1]

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountFeaturesGetFlag(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountFeaturesSetFlagCmd() *cobra.Command {
	opts := &options.AccountFeaturesSetFlagOptions{}

	cmd := &cobra.Command{
		Use:   "set-flag <account-id> <feature>",
		Short: "Set the state of a feature flag for an account",
		Args:  ExactArgsWithUsage(2, "account-id", "feature"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID
			opts.Feature = args[1]

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountFeaturesSetFlag(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.State, "state", "", "Feature state: on, off, allowed, allowed_on")
	mustMarkRequired(cmd, "state")

	return cmd
}

func newAccountFeaturesDeleteFlagCmd() *cobra.Command {
	opts := &options.AccountFeaturesDeleteFlagOptions{}

	cmd := &cobra.Command{
		Use:   "delete-flag <account-id> <feature>",
		Short: "Remove a feature flag override for an account",
		Args:  ExactArgsWithUsage(2, "account-id", "feature"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID
			opts.Feature = args[1]

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountFeaturesDeleteFlag(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountFeaturesSettingsCmd() *cobra.Command {
	opts := &options.AccountFeaturesSettingsOptions{}

	cmd := &cobra.Command{
		Use:   "settings <account-id>",
		Short: "Get settings for an account",
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

			return runAccountFeaturesSettings(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountFeaturesPermissionsCmd() *cobra.Command {
	opts := &options.AccountFeaturesPermissionsOptions{}

	cmd := &cobra.Command{
		Use:   "permissions <account-id>",
		Short: "Get permissions for the current user in an account",
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

			return runAccountFeaturesPermissions(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.Permissions, "permissions", []string{}, "Comma-separated list of permissions to check")

	return cmd
}

func runAccountFeaturesList(ctx context.Context, client *api.Client, opts *options.AccountFeaturesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-features.list", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAccountFeaturesService(client)
	features, err := svc.List(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "account-features.list", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-features.list", len(features))

	return formatOutput(features, func() {
		fmt.Printf("%-30s %-30s %-12s\n", "FEATURE", "DISPLAY_NAME", "APPLIES_TO")
		fmt.Println(strings.Repeat("-", 75))
		for _, f := range features {
			name := f.DisplayName
			if len(name) > 28 {
				name = name[:25] + "..."
			}
			fmt.Printf("%-30s %-30s %-12s\n", f.Feature, name, f.AppliesTo)
		}
		fmt.Printf("\nTotal: %d feature(s)\n", len(features))
	})
}

func runAccountFeaturesListEnabled(ctx context.Context, client *api.Client, opts *options.AccountFeaturesListEnabledOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-features.list-enabled", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAccountFeaturesService(client)
	features, err := svc.ListEnabled(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "account-features.list-enabled", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-features.list-enabled", len(features))

	return formatOutput(features, func() {
		fmt.Println("ENABLED FEATURES")
		fmt.Println(strings.Repeat("-", 40))
		for _, f := range features {
			fmt.Println(f)
		}
		fmt.Printf("\nTotal: %d feature(s)\n", len(features))
	})
}

func runAccountFeaturesGetFlag(ctx context.Context, client *api.Client, opts *options.AccountFeaturesGetFlagOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-features.get-flag", map[string]interface{}{
		"account_id": opts.AccountID,
		"feature":    opts.Feature,
	})

	svc := api.NewAccountFeaturesService(client)
	flag, err := svc.GetFlag(ctx, opts.AccountID, opts.Feature)
	if err != nil {
		logger.LogCommandError(ctx, "account-features.get-flag", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-features.get-flag", 1)

	return formatOutput(flag, func() {
		fmt.Printf("Feature:          %s\n", flag.Feature)
		fmt.Printf("State:            %s\n", flag.State)
		fmt.Printf("Transition Locked: %v\n", flag.TransitionLocked)
		if flag.LockedAt != "" {
			fmt.Printf("Locked At:        %s\n", flag.LockedAt)
		}
	})
}

func runAccountFeaturesSetFlag(ctx context.Context, client *api.Client, opts *options.AccountFeaturesSetFlagOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-features.set-flag", map[string]interface{}{
		"account_id": opts.AccountID,
		"feature":    opts.Feature,
		"state":      opts.State,
	})

	svc := api.NewAccountFeaturesService(client)
	flag, err := svc.SetFlag(ctx, opts.AccountID, opts.Feature, opts.State)
	if err != nil {
		logger.LogCommandError(ctx, "account-features.set-flag", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-features.set-flag", 1)

	return formatSuccessOutput(flag, fmt.Sprintf("Feature flag %q set to %q", opts.Feature, opts.State))
}

func runAccountFeaturesDeleteFlag(ctx context.Context, client *api.Client, opts *options.AccountFeaturesDeleteFlagOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-features.delete-flag", map[string]interface{}{
		"account_id": opts.AccountID,
		"feature":    opts.Feature,
	})

	svc := api.NewAccountFeaturesService(client)
	if err := svc.DeleteFlag(ctx, opts.AccountID, opts.Feature); err != nil {
		logger.LogCommandError(ctx, "account-features.delete-flag", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-features.delete-flag", 1)
	printInfo("Feature flag %q deleted (reverted to parent value)\n", opts.Feature)

	return nil
}

func runAccountFeaturesSettings(ctx context.Context, client *api.Client, opts *options.AccountFeaturesSettingsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-features.settings", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	svc := api.NewAccountFeaturesService(client)
	settings, err := svc.GetSettings(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "account-features.settings", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-features.settings", 1)

	return formatOutput(settings, func() {
		fmt.Printf("Account Settings\n")
		fmt.Printf("================\n\n")
		fmt.Printf("Restrict Student Past View:   %v\n", settings.RestrictStudentPastView)
		fmt.Printf("Restrict Student Future View: %v\n", settings.RestrictStudentFutureView)
		fmt.Printf("Hide Distribution Graphs:     %v\n", settings.HideDistributionGraphs)
		fmt.Printf("Lock All Announcements:       %v\n", settings.LockAllAnnouncements)
		fmt.Printf("Usage Rights Required:        %v\n", settings.UsageRightsRequired)
		fmt.Printf("Default Due Date Restricted:  %v\n", settings.DefaultDueDateRestricted)
	})
}

func runAccountFeaturesPermissions(ctx context.Context, client *api.Client, opts *options.AccountFeaturesPermissionsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-features.permissions", map[string]interface{}{
		"account_id":  opts.AccountID,
		"permissions": opts.Permissions,
	})

	svc := api.NewAccountFeaturesService(client)
	perms, err := svc.GetPermissions(ctx, opts.AccountID, opts.Permissions)
	if err != nil {
		logger.LogCommandError(ctx, "account-features.permissions", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-features.permissions", len(perms))

	return formatOutput(perms, func() {
		fmt.Printf("%-40s %-10s\n", "PERMISSION", "ALLOWED")
		fmt.Println(strings.Repeat("-", 52))
		for perm, allowed := range perms {
			fmt.Printf("%-40s %-10v\n", perm, allowed)
		}
	})
}
