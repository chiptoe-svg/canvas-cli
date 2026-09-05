package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/commands/internal/options"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

// foldersCmd represents the folders command group
var foldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "Manage Canvas file folders",
	Long: `Manage Canvas file folders for courses, groups, and users.

Folders organize files in Canvas. Each context (course, group, user) has its
own folder hierarchy starting at a root folder.

Examples:
  canvas folders list --course-id 123
  canvas folders list --group-id 456
  canvas folders list --user-id 789
  canvas folders get --folder-id 101
  canvas folders create --course-id 123 --name "Lectures"
  canvas folders delete --folder-id 101`,
}

func init() {
	rootCmd.AddCommand(foldersCmd)
	foldersCmd.AddCommand(newFoldersListCmd())
	foldersCmd.AddCommand(newFoldersGetCmd())
	foldersCmd.AddCommand(newFoldersResolvePathCmd())
	foldersCmd.AddCommand(newFoldersCreateCmd())
	foldersCmd.AddCommand(newFoldersUpdateCmd())
	foldersCmd.AddCommand(newFoldersDeleteCmd())
	foldersCmd.AddCommand(newFoldersMediaCmd())
	foldersCmd.AddCommand(newFoldersCopyCmd())
}

func newFoldersListCmd() *cobra.Command {
	opts := &options.FoldersListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List folders",
		Long: `List folders for a course, group, or user, or sub-folders within a folder.

Specify one of --course-id, --group-id, or --user-id to list top-level folders.
Use --folder-id to list sub-folders within a specific folder.

Examples:
  canvas folders list --course-id 123
  canvas folders list --group-id 456
  canvas folders list --user-id 789
  canvas folders list --folder-id 101`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFoldersList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Group ID")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID")
	cmd.Flags().Int64Var(&opts.FolderID, "folder-id", 0, "List sub-folders of this folder ID")

	return cmd
}

func newFoldersGetCmd() *cobra.Command {
	opts := &options.FoldersGetOptions{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a folder",
		Long: `Get details for a specific folder by ID.

Examples:
  canvas folders get --folder-id 101`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFoldersGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.FolderID, "folder-id", 0, "Folder ID (required)")
	mustMarkRequired(cmd, "folder-id")

	return cmd
}

func newFoldersResolvePathCmd() *cobra.Command {
	opts := &options.FoldersResolvePathOptions{}

	cmd := &cobra.Command{
		Use:   "resolve-path",
		Short: "Resolve a folder path",
		Long: `Resolve a full folder path and return the folder hierarchy.

Returns all folders from root to the specified path, in order.
Omit --path to get the root folder.

Examples:
  canvas folders resolve-path --course-id 123 --path "lectures/week1"
  canvas folders resolve-path --group-id 456 --path "shared"
  canvas folders resolve-path --user-id 789`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFoldersResolvePath(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Group ID")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID")
	cmd.Flags().StringVar(&opts.Path, "path", "", "Full folder path to resolve (e.g. 'lectures/week1')")

	return cmd
}

func newFoldersCreateCmd() *cobra.Command {
	opts := &options.FoldersCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a folder",
		Long: `Create a new folder in a course, group, user context, or inside a parent folder.

Specify a context (--course-id, --group-id, --user-id) or --parent-folder-id.

Examples:
  canvas folders create --course-id 123 --name "Lectures"
  canvas folders create --group-id 456 --name "Resources"
  canvas folders create --parent-folder-id 101 --name "Week 1"
  canvas folders create --course-id 123 --name "Archive" --hidden`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFoldersCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Group ID")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID")
	cmd.Flags().Int64Var(&opts.ParentFolderID, "parent-folder-id", 0, "Parent folder ID")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Folder name (required)")
	cmd.Flags().StringVar(&opts.LockAt, "lock-at", "", "Lock date (ISO 8601)")
	cmd.Flags().StringVar(&opts.UnlockAt, "unlock-at", "", "Unlock date (ISO 8601)")
	cmd.Flags().BoolVar(&opts.Locked, "locked", false, "Lock the folder")
	cmd.Flags().BoolVar(&opts.Hidden, "hidden", false, "Hide the folder from students")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Sort position")
	mustMarkRequired(cmd, "name")

	return cmd
}

func newFoldersUpdateCmd() *cobra.Command {
	opts := &options.FoldersUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a folder",
		Long: `Update an existing folder's properties.

Examples:
  canvas folders update --folder-id 101 --name "New Name"
  canvas folders update --folder-id 101 --hidden
  canvas folders update --folder-id 101 --locked --lock-at "2024-12-31T23:59:59Z"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.NameSet = cmd.Flags().Changed("name")
			opts.ParentFolderIDSet = cmd.Flags().Changed("parent-folder-id")
			opts.LockAtSet = cmd.Flags().Changed("lock-at")
			opts.UnlockAtSet = cmd.Flags().Changed("unlock-at")
			opts.LockedSet = cmd.Flags().Changed("locked")
			opts.HiddenSet = cmd.Flags().Changed("hidden")
			opts.PositionSet = cmd.Flags().Changed("position")
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFoldersUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.FolderID, "folder-id", 0, "Folder ID (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "New folder name")
	cmd.Flags().Int64Var(&opts.ParentFolderID, "parent-folder-id", 0, "Move to this parent folder")
	cmd.Flags().StringVar(&opts.LockAt, "lock-at", "", "Lock date (ISO 8601)")
	cmd.Flags().StringVar(&opts.UnlockAt, "unlock-at", "", "Unlock date (ISO 8601)")
	cmd.Flags().BoolVar(&opts.Locked, "locked", false, "Lock the folder")
	cmd.Flags().BoolVar(&opts.Hidden, "hidden", false, "Hide from students")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Sort position")
	mustMarkRequired(cmd, "folder-id")

	return cmd
}

func newFoldersDeleteCmd() *cobra.Command {
	opts := &options.FoldersDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a folder",
		Long: `Delete a folder. By default only empty folders can be deleted.
Use --force to delete non-empty folders.

Examples:
  canvas folders delete --folder-id 101
  canvas folders delete --folder-id 101 --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFoldersDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.FolderID, "folder-id", 0, "Folder ID (required)")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Delete non-empty folders")
	mustMarkRequired(cmd, "folder-id")

	return cmd
}

func newFoldersMediaCmd() *cobra.Command {
	opts := &options.FoldersMediaOptions{}

	cmd := &cobra.Command{
		Use:   "media",
		Short: "Get the media upload folder",
		Long: `Get the designated upload folder for media files.

Returns the folder that media uploads are placed in. Creates it if it
does not yet exist.

Examples:
  canvas folders media --course-id 123
  canvas folders media --group-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFoldersMedia(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Group ID")

	return cmd
}

func newFoldersCopyCmd() *cobra.Command {
	opts := &options.FoldersCopyOptions{}

	cmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy a folder into a destination folder",
		Long: `Copy a folder and its contents into a destination folder.

Examples:
  canvas folders copy --dest-folder-id 20 --source-folder-id 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFoldersCopy(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.DestFolderID, "dest-folder-id", 0, "Destination folder ID (required)")
	cmd.Flags().Int64Var(&opts.SourceFolderID, "source-folder-id", 0, "Source folder ID to copy (required)")
	mustMarkRequired(cmd, "dest-folder-id", "source-folder-id")

	return cmd
}

func runFoldersList(ctx context.Context, client *api.Client, opts *options.FoldersListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "folders.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
		"user_id":   opts.UserID,
		"folder_id": opts.FolderID,
	})

	foldersService := api.NewFoldersService(client)

	var folders []api.Folder
	var err error

	if opts.FolderID > 0 {
		folders, err = foldersService.ListFolderSubFolders(ctx, opts.FolderID, nil)
	} else {
		folders, err = foldersService.ListContextFolders(ctx, opts.CourseID, opts.GroupID, opts.UserID, nil)
	}

	if err != nil {
		logger.LogCommandError(ctx, "folders.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"group_id":  opts.GroupID,
			"user_id":   opts.UserID,
			"folder_id": opts.FolderID,
		})
		return fmt.Errorf("failed to list folders: %w", err)
	}

	printVerbose("Found %d folders:\n\n", len(folders))
	logger.LogCommandComplete(ctx, "folders.list", len(folders))
	return formatEmptyOrOutput(folders, "No folders found")
}

func runFoldersGet(ctx context.Context, client *api.Client, opts *options.FoldersGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "folders.get", map[string]interface{}{
		"folder_id": opts.FolderID,
	})

	foldersService := api.NewFoldersService(client)

	folder, err := foldersService.GetFolder(ctx, opts.FolderID)
	if err != nil {
		logger.LogCommandError(ctx, "folders.get", err, map[string]interface{}{
			"folder_id": opts.FolderID,
		})
		return fmt.Errorf("failed to get folder: %w", err)
	}

	logger.LogCommandComplete(ctx, "folders.get", 1)
	return formatOutput(folder, nil)
}

func runFoldersResolvePath(ctx context.Context, client *api.Client, opts *options.FoldersResolvePathOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "folders.resolve-path", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
		"user_id":   opts.UserID,
		"path":      opts.Path,
	})

	foldersService := api.NewFoldersService(client)

	folders, err := foldersService.ResolvePath(ctx, opts.CourseID, opts.GroupID, opts.UserID, opts.Path)
	if err != nil {
		logger.LogCommandError(ctx, "folders.resolve-path", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"group_id":  opts.GroupID,
			"user_id":   opts.UserID,
			"path":      opts.Path,
		})
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	logger.LogCommandComplete(ctx, "folders.resolve-path", len(folders))
	return formatEmptyOrOutput(folders, "Path not found")
}

func runFoldersCreate(ctx context.Context, client *api.Client, opts *options.FoldersCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "folders.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
		"user_id":   opts.UserID,
		"name":      opts.Name,
	})

	foldersService := api.NewFoldersService(client)

	params := &api.CreateFolderParams{
		Name:     opts.Name,
		LockAt:   opts.LockAt,
		UnlockAt: opts.UnlockAt,
		Locked:   opts.Locked,
		Hidden:   opts.Hidden,
		Position: opts.Position,
	}

	var folder *api.Folder
	var err error

	if opts.ParentFolderID > 0 {
		folder, err = foldersService.CreateSubFolder(ctx, opts.ParentFolderID, params)
	} else {
		folder, err = foldersService.CreateContextFolder(ctx, opts.CourseID, opts.GroupID, opts.UserID, params)
	}

	if err != nil {
		logger.LogCommandError(ctx, "folders.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"group_id":  opts.GroupID,
			"name":      opts.Name,
		})
		return fmt.Errorf("failed to create folder: %w", err)
	}

	logger.LogCommandComplete(ctx, "folders.create", 1)
	return formatSuccessOutput(folder, "Folder created successfully!")
}

func runFoldersUpdate(ctx context.Context, client *api.Client, opts *options.FoldersUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "folders.update", map[string]interface{}{
		"folder_id": opts.FolderID,
	})

	foldersService := api.NewFoldersService(client)

	params := &api.UpdateFolderParams{}

	if opts.NameSet {
		params.Name = &opts.Name
	}
	if opts.ParentFolderIDSet {
		params.ParentFolderID = &opts.ParentFolderID
	}
	if opts.LockAtSet {
		params.LockAt = &opts.LockAt
	}
	if opts.UnlockAtSet {
		params.UnlockAt = &opts.UnlockAt
	}
	if opts.LockedSet {
		params.Locked = &opts.Locked
	}
	if opts.HiddenSet {
		params.Hidden = &opts.Hidden
	}
	if opts.PositionSet {
		params.Position = &opts.Position
	}

	folder, err := foldersService.UpdateFolder(ctx, opts.FolderID, params)
	if err != nil {
		logger.LogCommandError(ctx, "folders.update", err, map[string]interface{}{
			"folder_id": opts.FolderID,
		})
		return fmt.Errorf("failed to update folder: %w", err)
	}

	logger.LogCommandComplete(ctx, "folders.update", 1)
	return formatSuccessOutput(folder, "Folder updated successfully!")
}

func runFoldersDelete(ctx context.Context, client *api.Client, opts *options.FoldersDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "folders.delete", map[string]interface{}{
		"folder_id": opts.FolderID,
		"force":     opts.Force,
	})

	// Confirm deletion (reuse the confirm helper)
	confirmed, err := confirmDelete("folder", opts.FolderID, opts.Force)
	if err != nil {
		logger.LogCommandError(ctx, "folders.delete", err, map[string]interface{}{})
		return err
	}
	if !confirmed {
		logger.LogCommandComplete(ctx, "folders.delete", 0)
		fmt.Println("Delete cancelled")
		return nil
	}

	foldersService := api.NewFoldersService(client)

	if err := foldersService.DeleteFolder(ctx, opts.FolderID, opts.Force); err != nil {
		logger.LogCommandError(ctx, "folders.delete", err, map[string]interface{}{
			"folder_id": opts.FolderID,
		})
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	logger.LogCommandComplete(ctx, "folders.delete", 1)
	printInfo("Folder %d deleted successfully\n", opts.FolderID)
	return nil
}

func runFoldersMedia(ctx context.Context, client *api.Client, opts *options.FoldersMediaOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "folders.media", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
	})

	foldersService := api.NewFoldersService(client)

	folder, err := foldersService.GetMediaFolder(ctx, opts.CourseID, opts.GroupID)
	if err != nil {
		logger.LogCommandError(ctx, "folders.media", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"group_id":  opts.GroupID,
		})
		return fmt.Errorf("failed to get media folder: %w", err)
	}

	logger.LogCommandComplete(ctx, "folders.media", 1)
	return formatOutput(folder, nil)
}

func runFoldersCopy(ctx context.Context, client *api.Client, opts *options.FoldersCopyOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "folders.copy", map[string]interface{}{
		"dest_folder_id":   opts.DestFolderID,
		"source_folder_id": opts.SourceFolderID,
	})

	foldersService := api.NewFoldersService(client)

	folder, err := foldersService.CopyFolder(ctx, opts.DestFolderID, &api.CopyFolderParams{
		SourceFolderID: opts.SourceFolderID,
	})
	if err != nil {
		logger.LogCommandError(ctx, "folders.copy", err, map[string]interface{}{
			"dest_folder_id":   opts.DestFolderID,
			"source_folder_id": opts.SourceFolderID,
		})
		return fmt.Errorf("failed to copy folder: %w", err)
	}

	logger.LogCommandComplete(ctx, "folders.copy", 1)
	return formatSuccessOutput(folder, "Folder copied successfully!")
}
