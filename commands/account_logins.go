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

var accountLoginsCmd = &cobra.Command{
	Use:   "account-logins",
	Short: "Manage Canvas account logins",
	Long: `Manage user login identities within a Canvas account.

Examples:
  canvas account-logins list 1
  canvas account-logins list 1 --user-id 100
  canvas account-logins create 1 --user-id 100 --unique-id alice@example.com
  canvas account-logins update 1 42 --unique-id new@example.com`,
}

func init() {
	rootCmd.AddCommand(accountLoginsCmd)
	accountLoginsCmd.AddCommand(newAccountLoginsListCmd())
	accountLoginsCmd.AddCommand(newAccountLoginsCreateCmd())
	accountLoginsCmd.AddCommand(newAccountLoginsUpdateCmd())
}

func newAccountLoginsListCmd() *cobra.Command {
	opts := &options.AccountLoginsListOptions{}

	cmd := &cobra.Command{
		Use:   "list <account-id>",
		Short: "List logins for an account",
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

			return runAccountLoginsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "Filter logins by user ID")

	return cmd
}

func newAccountLoginsCreateCmd() *cobra.Command {
	opts := &options.AccountLoginsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <account-id>",
		Short: "Create a login for a user in an account",
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

			return runAccountLoginsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID to create the login for (required)")
	cmd.Flags().StringVar(&opts.UniqueID, "unique-id", "", "Login unique ID / username (required)")
	cmd.Flags().StringVar(&opts.Password, "password", "", "Login password")
	cmd.Flags().StringVar(&opts.SISUserID, "sis-user-id", "", "SIS user ID")
	mustMarkRequired(cmd, "user-id", "unique-id")

	return cmd
}

func newAccountLoginsUpdateCmd() *cobra.Command {
	opts := &options.AccountLoginsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <account-id> <login-id>",
		Short: "Update a login in an account",
		Args:  ExactArgsWithUsage(2, "account-id", "login-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			loginID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid login ID: %s", args[1])
			}
			opts.AccountID = accountID
			opts.LoginID = loginID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountLoginsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.UniqueID, "unique-id", "", "New unique ID / username")
	cmd.Flags().StringVar(&opts.Password, "password", "", "New password")

	return cmd
}

func runAccountLoginsList(ctx context.Context, client *api.Client, opts *options.AccountLoginsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-logins.list", map[string]interface{}{
		"account_id": opts.AccountID,
		"user_id":    opts.UserID,
	})

	service := api.NewAccountLoginsService(client)

	logins, err := service.List(ctx, opts.AccountID, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "account-logins.list", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-logins.list", len(logins))

	return formatEmptyOrOutput(logins, fmt.Sprintf("No logins found for account %d", opts.AccountID))
}

func runAccountLoginsCreate(ctx context.Context, client *api.Client, opts *options.AccountLoginsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-logins.create", map[string]interface{}{
		"account_id": opts.AccountID,
		"user_id":    opts.UserID,
	})

	service := api.NewAccountLoginsService(client)

	params := &api.LoginParams{
		UserID:    opts.UserID,
		UniqueID:  opts.UniqueID,
		Password:  opts.Password,
		SISUserID: opts.SISUserID,
	}

	login, err := service.Create(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "account-logins.create", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-logins.create", 1)

	return formatSuccessOutput(login, fmt.Sprintf("Login '%s' created with ID %d", login.UniqueID, login.ID))
}

func runAccountLoginsUpdate(ctx context.Context, client *api.Client, opts *options.AccountLoginsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "account-logins.update", map[string]interface{}{
		"account_id": opts.AccountID,
		"login_id":   opts.LoginID,
	})

	service := api.NewAccountLoginsService(client)

	params := &api.LoginParams{
		UniqueID: opts.UniqueID,
		Password: opts.Password,
	}

	login, err := service.Update(ctx, opts.AccountID, opts.LoginID, params)
	if err != nil {
		logger.LogCommandError(ctx, "account-logins.update", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "account-logins.update", 1)

	return formatSuccessOutput(login, fmt.Sprintf("Login %d updated", login.ID))
}
