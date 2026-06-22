package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// userFeaturesCmd represents the user-features command group
var userFeaturesCmd = &cobra.Command{
	Use:   "user-features",
	Short: "Manage user feature flags",
	Long: `Manage feature flags for Canvas users.

Examples:
  canvas user-features list --user-id 123
  canvas user-features list-enabled --user-id 123
  canvas user-features get-flag --user-id 123 --feature new_gradebook
  canvas user-features set-flag --user-id 123 --feature new_gradebook --state on
  canvas user-features delete-flag --user-id 123 --feature new_gradebook`,
}

func init() {
	rootCmd.AddCommand(userFeaturesCmd)
	userFeaturesCmd.AddCommand(newUserFeaturesListCmd())
	userFeaturesCmd.AddCommand(newUserFeaturesListEnabledCmd())
	userFeaturesCmd.AddCommand(newUserFeaturesGetFlagCmd())
	userFeaturesCmd.AddCommand(newUserFeaturesSetFlagCmd())
	userFeaturesCmd.AddCommand(newUserFeaturesDeleteFlagCmd())
}

func newUserFeaturesListCmd() *cobra.Command {
	opts := &options.UserFeaturesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all features for a user",
		Long: `List all features available for a specific user.

Examples:
  canvas user-features list --user-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUserFeaturesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	mustMarkRequired(cmd, "user-id")

	return cmd
}

func newUserFeaturesListEnabledCmd() *cobra.Command {
	opts := &options.UserFeaturesListEnabledOptions{}

	cmd := &cobra.Command{
		Use:   "list-enabled",
		Short: "List enabled features for a user",
		Long: `List only enabled features for a specific user.

Examples:
  canvas user-features list-enabled --user-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUserFeaturesListEnabled(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	mustMarkRequired(cmd, "user-id")

	return cmd
}

func newUserFeaturesGetFlagCmd() *cobra.Command {
	opts := &options.UserFeaturesGetFlagOptions{}

	cmd := &cobra.Command{
		Use:   "get-flag",
		Short: "Get a feature flag for a user",
		Long: `Get the state of a specific feature flag for a user.

Examples:
  canvas user-features get-flag --user-id 123 --feature new_gradebook`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUserFeaturesGetFlag(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.Flags().StringVar(&opts.Feature, "feature", "", "Feature name (required)")
	mustMarkRequired(cmd, "user-id", "feature")

	return cmd
}

func newUserFeaturesSetFlagCmd() *cobra.Command {
	opts := &options.UserFeaturesSetFlagOptions{}

	cmd := &cobra.Command{
		Use:   "set-flag",
		Short: "Set a feature flag for a user",
		Long: `Set the state of a feature flag for a user.

Valid states: on, off, allowed

Examples:
  canvas user-features set-flag --user-id 123 --feature new_gradebook --state on
  canvas user-features set-flag --user-id 123 --feature new_gradebook --state off`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUserFeaturesSetFlag(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.Flags().StringVar(&opts.Feature, "feature", "", "Feature name (required)")
	cmd.Flags().StringVar(&opts.State, "state", "", "Feature state: on, off, allowed (required)")
	mustMarkRequired(cmd, "user-id", "feature", "state")

	return cmd
}

func newUserFeaturesDeleteFlagCmd() *cobra.Command {
	opts := &options.UserFeaturesDeleteFlagOptions{}

	cmd := &cobra.Command{
		Use:   "delete-flag",
		Short: "Delete a feature flag for a user",
		Long: `Delete (reset to default) a feature flag for a user.

Examples:
  canvas user-features delete-flag --user-id 123 --feature new_gradebook`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUserFeaturesDeleteFlag(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.Flags().StringVar(&opts.Feature, "feature", "", "Feature name (required)")
	mustMarkRequired(cmd, "user-id", "feature")

	return cmd
}

func runUserFeaturesList(ctx context.Context, client *api.Client, opts *options.UserFeaturesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "user-features.list", map[string]interface{}{
		"user_id": opts.UserID,
	})

	svc := api.NewUserFeaturesService(client)

	features, err := svc.List(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "user-features.list", err, map[string]interface{}{
			"user_id": opts.UserID,
		})
		return fmt.Errorf("failed to list user features: %w", err)
	}

	logger.LogCommandComplete(ctx, "user-features.list", len(features))
	return formatEmptyOrOutput(features, "No features found")
}

func runUserFeaturesListEnabled(ctx context.Context, client *api.Client, opts *options.UserFeaturesListEnabledOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "user-features.list-enabled", map[string]interface{}{
		"user_id": opts.UserID,
	})

	svc := api.NewUserFeaturesService(client)

	features, err := svc.ListEnabled(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "user-features.list-enabled", err, map[string]interface{}{
			"user_id": opts.UserID,
		})
		return fmt.Errorf("failed to list enabled user features: %w", err)
	}

	logger.LogCommandComplete(ctx, "user-features.list-enabled", len(features))
	return formatEmptyOrOutput(features, "No enabled features found")
}

func runUserFeaturesGetFlag(ctx context.Context, client *api.Client, opts *options.UserFeaturesGetFlagOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "user-features.get-flag", map[string]interface{}{
		"user_id": opts.UserID,
		"feature": opts.Feature,
	})

	svc := api.NewUserFeaturesService(client)

	flag, err := svc.GetFlag(ctx, opts.UserID, opts.Feature)
	if err != nil {
		logger.LogCommandError(ctx, "user-features.get-flag", err, map[string]interface{}{
			"user_id": opts.UserID,
			"feature": opts.Feature,
		})
		return fmt.Errorf("failed to get feature flag: %w", err)
	}

	logger.LogCommandComplete(ctx, "user-features.get-flag", 1)
	return formatOutput(flag, nil)
}

func runUserFeaturesSetFlag(ctx context.Context, client *api.Client, opts *options.UserFeaturesSetFlagOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "user-features.set-flag", map[string]interface{}{
		"user_id": opts.UserID,
		"feature": opts.Feature,
		"state":   opts.State,
	})

	svc := api.NewUserFeaturesService(client)

	flag, err := svc.SetFlag(ctx, opts.UserID, opts.Feature, opts.State)
	if err != nil {
		logger.LogCommandError(ctx, "user-features.set-flag", err, map[string]interface{}{
			"user_id": opts.UserID,
			"feature": opts.Feature,
		})
		return fmt.Errorf("failed to set feature flag: %w", err)
	}

	logger.LogCommandComplete(ctx, "user-features.set-flag", 1)
	return formatSuccessOutput(flag, "Feature flag set successfully!")
}

func runUserFeaturesDeleteFlag(ctx context.Context, client *api.Client, opts *options.UserFeaturesDeleteFlagOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "user-features.delete-flag", map[string]interface{}{
		"user_id": opts.UserID,
		"feature": opts.Feature,
	})

	svc := api.NewUserFeaturesService(client)

	if err := svc.DeleteFlag(ctx, opts.UserID, opts.Feature); err != nil {
		logger.LogCommandError(ctx, "user-features.delete-flag", err, map[string]interface{}{
			"user_id": opts.UserID,
			"feature": opts.Feature,
		})
		return fmt.Errorf("failed to delete feature flag: %w", err)
	}

	logger.LogCommandComplete(ctx, "user-features.delete-flag", 1)
	printInfo("Feature flag '%s' deleted successfully for user %d\n", opts.Feature, opts.UserID)
	return nil
}
