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

// bookmarksCmd represents the bookmarks command group
var bookmarksCmd = &cobra.Command{
	Use:   "bookmarks",
	Short: "Manage Canvas bookmarks",
	Long: `Manage Canvas user bookmarks.

Examples:
  canvas bookmarks list
  canvas bookmarks get 123
  canvas bookmarks create --name "My Page" --url "https://example.instructure.com/courses/1"
  canvas bookmarks update 123 --name "Updated Name"
  canvas bookmarks delete 123`,
}

func init() {
	rootCmd.AddCommand(bookmarksCmd)
	bookmarksCmd.AddCommand(newBookmarksListCmd())
	bookmarksCmd.AddCommand(newBookmarksGetCmd())
	bookmarksCmd.AddCommand(newBookmarksCreateCmd())
	bookmarksCmd.AddCommand(newBookmarksUpdateCmd())
	bookmarksCmd.AddCommand(newBookmarksDeleteCmd())
}

func newBookmarksListCmd() *cobra.Command {
	opts := &options.BookmarksListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all bookmarks for the current user",
		Long: `List all bookmarks for the currently authenticated user.

Examples:
  canvas bookmarks list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runBookmarksList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newBookmarksGetCmd() *cobra.Command {
	opts := &options.BookmarksGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a specific bookmark",
		Long: `Get details of a specific bookmark by ID.

Examples:
  canvas bookmarks get 123`,
		Args: ExactArgsWithUsage(1, "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid bookmark ID: %s", args[0])
			}
			opts.ID = id
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runBookmarksGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newBookmarksCreateCmd() *cobra.Command {
	opts := &options.BookmarksCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new bookmark",
		Long: `Create a new bookmark for the current user.

Examples:
  canvas bookmarks create --name "My Page" --url "https://example.instructure.com/courses/1"
  canvas bookmarks create --name "Week 1" --url "https://example.instructure.com/courses/1/modules" --position 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runBookmarksCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Bookmark name (required)")
	cmd.Flags().StringVar(&opts.URL, "url", "", "Bookmark URL (required)")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Position in bookmark list")
	mustMarkRequired(cmd, "name", "url")

	return cmd
}

func newBookmarksUpdateCmd() *cobra.Command {
	opts := &options.BookmarksUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an existing bookmark",
		Long: `Update an existing bookmark by ID.

Examples:
  canvas bookmarks update 123 --name "New Name"
  canvas bookmarks update 123 --url "https://example.instructure.com/new-url"
  canvas bookmarks update 123 --position 2`,
		Args: ExactArgsWithUsage(1, "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid bookmark ID: %s", args[0])
			}
			opts.ID = id
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runBookmarksUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Bookmark name")
	cmd.Flags().StringVar(&opts.URL, "url", "", "Bookmark URL")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Position in bookmark list")

	return cmd
}

func newBookmarksDeleteCmd() *cobra.Command {
	opts := &options.BookmarksDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a bookmark",
		Long: `Delete a bookmark by ID.

Examples:
  canvas bookmarks delete 123`,
		Args: ExactArgsWithUsage(1, "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid bookmark ID: %s", args[0])
			}
			opts.ID = id
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runBookmarksDelete(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func runBookmarksList(ctx context.Context, client *api.Client, opts *options.BookmarksListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "bookmarks.list", map[string]interface{}{})

	svc := api.NewBookmarksService(client)

	bookmarks, err := svc.List(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "bookmarks.list", err, map[string]interface{}{})
		return fmt.Errorf("failed to list bookmarks: %w", err)
	}

	logger.LogCommandComplete(ctx, "bookmarks.list", len(bookmarks))
	return formatEmptyOrOutput(bookmarks, "No bookmarks found")
}

func runBookmarksGet(ctx context.Context, client *api.Client, opts *options.BookmarksGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "bookmarks.get", map[string]interface{}{
		"id": opts.ID,
	})

	svc := api.NewBookmarksService(client)

	bookmark, err := svc.Get(ctx, opts.ID)
	if err != nil {
		logger.LogCommandError(ctx, "bookmarks.get", err, map[string]interface{}{
			"id": opts.ID,
		})
		return fmt.Errorf("failed to get bookmark: %w", err)
	}

	logger.LogCommandComplete(ctx, "bookmarks.get", 1)
	return formatOutput(bookmark, nil)
}

func runBookmarksCreate(ctx context.Context, client *api.Client, opts *options.BookmarksCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "bookmarks.create", map[string]interface{}{
		"name": opts.Name,
		"url":  opts.URL,
	})

	svc := api.NewBookmarksService(client)

	params := api.CreateBookmarkParams{
		Name:     opts.Name,
		URL:      opts.URL,
		Position: opts.Position,
	}

	bookmark, err := svc.Create(ctx, params)
	if err != nil {
		logger.LogCommandError(ctx, "bookmarks.create", err, map[string]interface{}{
			"name": opts.Name,
		})
		return fmt.Errorf("failed to create bookmark: %w", err)
	}

	logger.LogCommandComplete(ctx, "bookmarks.create", 1)
	return formatSuccessOutput(bookmark, "Bookmark created successfully!")
}

func runBookmarksUpdate(ctx context.Context, client *api.Client, opts *options.BookmarksUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "bookmarks.update", map[string]interface{}{
		"id": opts.ID,
	})

	svc := api.NewBookmarksService(client)

	params := api.UpdateBookmarkParams{
		Name:     opts.Name,
		URL:      opts.URL,
		Position: opts.Position,
	}

	bookmark, err := svc.Update(ctx, opts.ID, params)
	if err != nil {
		logger.LogCommandError(ctx, "bookmarks.update", err, map[string]interface{}{
			"id": opts.ID,
		})
		return fmt.Errorf("failed to update bookmark: %w", err)
	}

	logger.LogCommandComplete(ctx, "bookmarks.update", 1)
	return formatSuccessOutput(bookmark, "Bookmark updated successfully!")
}

func runBookmarksDelete(ctx context.Context, client *api.Client, opts *options.BookmarksDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "bookmarks.delete", map[string]interface{}{
		"id": opts.ID,
	})

	svc := api.NewBookmarksService(client)

	if err := svc.Delete(ctx, opts.ID); err != nil {
		logger.LogCommandError(ctx, "bookmarks.delete", err, map[string]interface{}{
			"id": opts.ID,
		})
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}

	logger.LogCommandComplete(ctx, "bookmarks.delete", 1)
	printInfo("Bookmark %d deleted successfully\n", opts.ID)
	return nil
}
