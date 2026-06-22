package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// observeesCmd represents the observees command group
var observeesCmd = &cobra.Command{
	Use:   "observees",
	Short: "Manage user observees and observers",
	Long: `Manage observees and observers for Canvas users.

Examples:
  canvas observees list --user-id 123
  canvas observees get --user-id 123 --observee-id 456
  canvas observees remove --user-id 123 --observee-id 456
  canvas observees observers list --user-id 123`,
}

// observeesObserversCmd represents the observees observers subcommand group
var observeesObserversCmd = &cobra.Command{
	Use:   "observers",
	Short: "Manage user observers",
	Long:  `Manage observers for Canvas users.`,
}

func init() {
	rootCmd.AddCommand(observeesCmd)
	observeesCmd.AddCommand(newObserveesListCmd())
	observeesCmd.AddCommand(newObserveesGetCmd())
	observeesCmd.AddCommand(newObserveesRemoveCmd())
	observeesCmd.AddCommand(observeesObserversCmd)
	observeesObserversCmd.AddCommand(newObserversListCmd())
}

func newObserveesListCmd() *cobra.Command {
	opts := &options.ObserveesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List observees for a user",
		Long: `List all observees for a specific user.

Examples:
  canvas observees list --user-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runObserveesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	mustMarkRequired(cmd, "user-id")

	return cmd
}

func newObserveesGetCmd() *cobra.Command {
	opts := &options.ObserveesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a specific observee for a user",
		Long: `Get details of a specific observee for a user.

Examples:
  canvas observees get --user-id 123 --observee-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runObserveesGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.Flags().Int64Var(&opts.ObserveeID, "observee-id", 0, "Observee user ID (required)")
	mustMarkRequired(cmd, "user-id", "observee-id")

	return cmd
}

func newObserveesRemoveCmd() *cobra.Command {
	opts := &options.ObserveesDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove an observee from a user",
		Long: `Remove an observee relationship from a user.

Examples:
  canvas observees remove --user-id 123 --observee-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runObserveesRemove(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.Flags().Int64Var(&opts.ObserveeID, "observee-id", 0, "Observee user ID (required)")
	mustMarkRequired(cmd, "user-id", "observee-id")

	return cmd
}

func newObserversListCmd() *cobra.Command {
	opts := &options.ObserversListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List observers for a user",
		Long: `List all observers for a specific user.

Examples:
  canvas observees observers list --user-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runObserversList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	mustMarkRequired(cmd, "user-id")

	return cmd
}

func runObserveesList(ctx context.Context, client *api.Client, opts *options.ObserveesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "observees.list", map[string]interface{}{
		"user_id": opts.UserID,
	})

	svc := api.NewObserveesService(client)

	observees, err := svc.ListObservees(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "observees.list", err, map[string]interface{}{
			"user_id": opts.UserID,
		})
		return fmt.Errorf("failed to list observees: %w", err)
	}

	logger.LogCommandComplete(ctx, "observees.list", len(observees))
	return formatEmptyOrOutput(observees, "No observees found")
}

func runObserveesGet(ctx context.Context, client *api.Client, opts *options.ObserveesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "observees.get", map[string]interface{}{
		"user_id":     opts.UserID,
		"observee_id": opts.ObserveeID,
	})

	svc := api.NewObserveesService(client)

	observee, err := svc.GetObservee(ctx, opts.UserID, opts.ObserveeID)
	if err != nil {
		logger.LogCommandError(ctx, "observees.get", err, map[string]interface{}{
			"user_id":     opts.UserID,
			"observee_id": opts.ObserveeID,
		})
		return fmt.Errorf("failed to get observee: %w", err)
	}

	logger.LogCommandComplete(ctx, "observees.get", 1)
	return formatOutput(observee, nil)
}

func runObserveesRemove(ctx context.Context, client *api.Client, opts *options.ObserveesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "observees.remove", map[string]interface{}{
		"user_id":     opts.UserID,
		"observee_id": opts.ObserveeID,
	})

	svc := api.NewObserveesService(client)

	_, err := svc.RemoveObservee(ctx, opts.UserID, opts.ObserveeID)
	if err != nil {
		logger.LogCommandError(ctx, "observees.remove", err, map[string]interface{}{
			"user_id":     opts.UserID,
			"observee_id": opts.ObserveeID,
		})
		return fmt.Errorf("failed to remove observee: %w", err)
	}

	logger.LogCommandComplete(ctx, "observees.remove", 1)
	printInfo("Observee %d removed successfully\n", opts.ObserveeID)
	return nil
}

func runObserversList(ctx context.Context, client *api.Client, opts *options.ObserversListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "observers.list", map[string]interface{}{
		"user_id": opts.UserID,
	})

	svc := api.NewObserveesService(client)

	observers, err := svc.ListObservers(ctx, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "observers.list", err, map[string]interface{}{
			"user_id": opts.UserID,
		})
		return fmt.Errorf("failed to list observers: %w", err)
	}

	logger.LogCommandComplete(ctx, "observers.list", len(observers))
	return formatEmptyOrOutput(observers, "No observers found")
}
