package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// jwtsCmd represents the jwts command group.
var jwtsCmd = &cobra.Command{
	Use:   "jwts",
	Short: "Create and refresh Canvas JWTs",
	Long: `Create new Canvas JWTs or refresh existing ones.

Canvas JWTs are short-lived tokens used for service-to-service authentication.

Examples:
  canvas jwts create
  canvas jwts refresh --token <existing-jwt>`,
}

func init() {
	rootCmd.AddCommand(jwtsCmd)
	jwtsCmd.AddCommand(newJWTCreateCmd())
	jwtsCmd.AddCommand(newJWTRefreshCmd())
}

func newJWTCreateCmd() *cobra.Command {
	opts := &options.JWTCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Canvas JWT",
		Long: `Create a new Canvas JWT for the current user.

Examples:
  canvas jwts create`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runJWTCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Workflow, "workflow", "", "Workflow state to embed in the JWT")

	return cmd
}

func newJWTRefreshCmd() *cobra.Command {
	opts := &options.JWTRefreshOptions{}

	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "Refresh an existing Canvas JWT",
		Long: `Refresh an existing Canvas JWT and get a new one.

Examples:
  canvas jwts refresh --token eyJ...`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runJWTRefresh(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Token, "token", "", "Existing JWT to refresh")
	mustMarkRequired(cmd, "token")

	return cmd
}

func runJWTCreate(ctx context.Context, client *api.Client, opts *options.JWTCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "jwts.create", map[string]interface{}{
		"workflow": opts.Workflow,
	})

	svc := api.NewJWTsService(client)

	jwt, err := svc.Create(ctx, opts.Workflow)
	if err != nil {
		logger.LogCommandError(ctx, "jwts.create", err, nil)
		return fmt.Errorf("failed to create JWT: %w", err)
	}

	logger.LogCommandComplete(ctx, "jwts.create", 1)
	return formatSuccessOutput(jwt, "JWT created successfully.")
}

func runJWTRefresh(ctx context.Context, client *api.Client, opts *options.JWTRefreshOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "jwts.refresh", map[string]interface{}{})

	svc := api.NewJWTsService(client)

	jwt, err := svc.Refresh(ctx, opts.Token)
	if err != nil {
		logger.LogCommandError(ctx, "jwts.refresh", err, nil)
		return fmt.Errorf("failed to refresh JWT: %w", err)
	}

	logger.LogCommandComplete(ctx, "jwts.refresh", 1)
	return formatSuccessOutput(jwt, "JWT refreshed successfully.")
}
