package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// filesCmd represents the files command group
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Manage Canvas files",
	Long: `Manage Canvas files including listing, uploading, downloading, and deleting files.

Examples:
  canvas files list --course-id 123
  canvas files get 456
  canvas files upload --course-id 123 document.pdf
  canvas files download 456 --destination ./downloaded.pdf
  canvas files delete 456`,
}

func init() {
	rootCmd.AddCommand(filesCmd)
	filesCmd.AddCommand(newFilesListCmd())
	filesCmd.AddCommand(newFilesGetCmd())
	filesCmd.AddCommand(newFilesUploadCmd())
	filesCmd.AddCommand(newFilesDownloadCmd())
	filesCmd.AddCommand(newFilesDeleteCmd())
	filesCmd.AddCommand(newFilesQuotaCmd())
	filesCmd.AddCommand(newFilesResetVerifierCmd())
	filesCmd.AddCommand(newFilesCopyCmd())
	filesCmd.AddCommand(newFilesUsageRightsCmd())
	filesCmd.AddCommand(newFilesRemoveUsageRightsCmd())
	filesCmd.AddCommand(newFilesLicensesCmd())
}

func newFilesListCmd() *cobra.Command {
	opts := &options.FilesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List files",
		Long: `List files in a course, group, folder, or user's files.

You must specify one of --course-id, --group-id, --folder-id, or --user-id.

Examples:
  canvas files list --course-id 123
  canvas files list --group-id 456
  canvas files list --folder-id 789
  canvas files list --user-id 101
  canvas files list --course-id 123 --search "assignment"
  canvas files list --course-id 123 --sort name --order asc`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Group ID")
	cmd.Flags().Int64Var(&opts.FolderID, "folder-id", 0, "Folder ID")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID")
	cmd.Flags().StringSliceVar(&opts.ContentTypes, "content-types", []string{}, "Filter by MIME types (comma-separated)")
	cmd.Flags().StringVar(&opts.SearchTerm, "search", "", "Search by file name")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")
	cmd.Flags().StringVar(&opts.Sort, "sort", "", "Sort by (name, size, created_at, updated_at, content_type)")
	cmd.Flags().StringVar(&opts.Order, "order", "", "Order direction (asc, desc)")

	return cmd
}

func newFilesGetCmd() *cobra.Command {
	opts := &options.FilesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <file-id>",
		Short: "Get file details",
		Long: `Get details of a specific file by ID.

Examples:
  canvas files get 456
  canvas files get 456 --include user`,
		Args: ExactArgsWithUsage(1, "file-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid file ID: %s", args[0])
			}
			opts.FileID = fileID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")

	return cmd
}

func newFilesUploadCmd() *cobra.Command {
	opts := &options.FilesUploadOptions{}

	cmd := &cobra.Command{
		Use:   "upload <file-path>",
		Short: "Upload a file",
		Long: `Upload a file to a course, folder, or user's files.

You must specify one of --course-id, --folder-id, or --user-id.

Examples:
  canvas files upload document.pdf --course-id 123
  canvas files upload image.png --folder-id 456
  canvas files upload data.csv --user-id 789
  canvas files upload file.pdf --course-id 123 --on-duplicate overwrite`,
		Args: ExactArgsWithUsage(1, "file-path"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.FilePath = args[0]

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesUpload(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.FolderID, "folder-id", 0, "Folder ID")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID")
	cmd.Flags().StringVar(&opts.OnDuplicate, "on-duplicate", "rename", "How to handle duplicates (overwrite, rename)")
	cmd.Flags().Int64Var(&opts.ParentFolderID, "parent-folder", 0, "Parent folder ID")
	cmd.Flags().BoolVar(&opts.Hidden, "hidden", false, "Hide from students")
	cmd.Flags().BoolVar(&opts.Locked, "locked", false, "Lock the file")

	return cmd
}

func newFilesDownloadCmd() *cobra.Command {
	opts := &options.FilesDownloadOptions{}

	cmd := &cobra.Command{
		Use:   "download <file-id>",
		Short: "Download a file",
		Long: `Download a file from Canvas to your local system.

Examples:
  canvas files download 456
  canvas files download 456 --destination ./my-file.pdf`,
		Args: ExactArgsWithUsage(1, "file-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid file ID: %s", args[0])
			}
			opts.FileID = fileID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesDownload(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Destination, "destination", "", "Destination file path (default: current directory with original filename)")

	return cmd
}

func newFilesDeleteCmd() *cobra.Command {
	opts := &options.FilesDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <file-id>",
		Short: "Delete a file",
		Long: `Delete a file from Canvas.

Examples:
  canvas files delete 456`,
		Args: ExactArgsWithUsage(1, "file-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid file ID: %s", args[0])
			}
			opts.FileID = fileID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func newFilesQuotaCmd() *cobra.Command {
	opts := &options.FilesQuotaOptions{}

	cmd := &cobra.Command{
		Use:   "quota",
		Short: "Get storage quota information",
		Long: `Get storage quota information for a course, group, or user.

You must specify one of --course-id, --group-id, or --user-id.

Examples:
  canvas files quota --course-id 123
  canvas files quota --group-id 456
  canvas files quota --user-id 789`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesQuota(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Group ID")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID")

	return cmd
}

func newFilesResetVerifierCmd() *cobra.Command {
	opts := &options.FilesResetVerifierOptions{}

	cmd := &cobra.Command{
		Use:   "reset-verifier <file-id>",
		Short: "Reset the link verifier for a file",
		Long: `Reset the link verifier for a file. Any existing links using the
previous verifier will no longer automatically grant access.

Examples:
  canvas files reset-verifier 456`,
		Args: ExactArgsWithUsage(1, "file-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid file ID: %s", args[0])
			}
			opts.FileID = fileID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesResetVerifier(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newFilesCopyCmd() *cobra.Command {
	opts := &options.FilesCopyOptions{}

	cmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy a file into a folder",
		Long: `Copy a file from elsewhere in Canvas into a destination folder.

Examples:
  canvas files copy --dest-folder-id 20 --source-file-id 10
  canvas files copy --dest-folder-id 20 --source-file-id 10 --on-duplicate rename`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesCopy(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.DestFolderID, "dest-folder-id", 0, "Destination folder ID (required)")
	cmd.Flags().Int64Var(&opts.SourceFileID, "source-file-id", 0, "Source file ID to copy (required)")
	cmd.Flags().StringVar(&opts.OnDuplicate, "on-duplicate", "", "Duplicate handling: overwrite, rename")
	mustMarkRequired(cmd, "dest-folder-id", "source-file-id")

	return cmd
}

func newFilesUsageRightsCmd() *cobra.Command {
	opts := &options.FilesUsageRightsOptions{}

	cmd := &cobra.Command{
		Use:   "set-usage-rights",
		Short: "Set copyright and license information on files",
		Long: `Set copyright and license information for one or more files.

use-justification values: own_copyright, used_by_permission, fair_use,
public_domain, creative_commons (license required when using creative_commons).

Examples:
  canvas files set-usage-rights --course-id 123 --file-ids 1,2,3 --use-justification own_copyright
  canvas files set-usage-rights --group-id 456 --file-ids 5 --use-justification creative_commons --license cc_by`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesSetUsageRights(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Group ID")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID")
	cmd.Flags().Int64SliceVar(&opts.FileIDs, "file-ids", []int64{}, "File IDs (comma-separated)")
	cmd.Flags().Int64SliceVar(&opts.FolderIDs, "folder-ids", []int64{}, "Folder IDs to apply rights to all files within")
	cmd.Flags().StringVar(&opts.UseJustification, "use-justification", "", "Use justification (required)")
	cmd.Flags().StringVar(&opts.LegalCopyright, "legal-copyright", "", "Legal copyright line")
	cmd.Flags().StringVar(&opts.License, "license", "", "License (required for creative_commons)")
	cmd.Flags().BoolVar(&opts.Publish, "publish", false, "Publish the files after setting rights")
	mustMarkRequired(cmd, "use-justification")

	return cmd
}

func newFilesRemoveUsageRightsCmd() *cobra.Command {
	opts := &options.FilesRemoveUsageRightsOptions{}

	cmd := &cobra.Command{
		Use:   "remove-usage-rights",
		Short: "Remove copyright and license information from files",
		Long: `Remove copyright and license information from one or more files.

Examples:
  canvas files remove-usage-rights --course-id 123 --file-ids 1,2,3
  canvas files remove-usage-rights --group-id 456 --folder-ids 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesRemoveUsageRights(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Group ID")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID")
	cmd.Flags().Int64SliceVar(&opts.FileIDs, "file-ids", []int64{}, "File IDs (comma-separated)")
	cmd.Flags().Int64SliceVar(&opts.FolderIDs, "folder-ids", []int64{}, "Folder IDs")

	return cmd
}

func newFilesLicensesCmd() *cobra.Command {
	opts := &options.FilesLicensesOptions{}

	cmd := &cobra.Command{
		Use:   "licenses",
		Short: "List available content licenses",
		Long: `List content licenses that can be applied to files.

Examples:
  canvas files licenses --course-id 123
  canvas files licenses --group-id 456
  canvas files licenses --user-id 789`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesLicenses(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Group ID")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID")

	return cmd
}

func runFilesList(ctx context.Context, client *api.Client, opts *options.FilesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
		"folder_id": opts.FolderID,
		"user_id":   opts.UserID,
	})

	filesService := api.NewFilesService(client)

	apiOpts := &api.ListFilesOptions{
		ContentTypes: opts.ContentTypes,
		SearchTerm:   opts.SearchTerm,
		Include:      opts.Include,
		Sort:         opts.Sort,
		Order:        opts.Order,
	}

	var files []api.Attachment
	var err error

	switch {
	case opts.CourseID > 0:
		files, err = filesService.ListCourseFiles(ctx, opts.CourseID, apiOpts)
	case opts.GroupID > 0:
		files, err = filesService.ListGroupFiles(ctx, opts.GroupID, apiOpts)
	case opts.FolderID > 0:
		files, err = filesService.ListFolderFiles(ctx, opts.FolderID, apiOpts)
	default:
		files, err = filesService.ListUserFiles(ctx, opts.UserID, apiOpts)
	}

	if err != nil {
		logger.LogCommandError(ctx, "files.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"group_id":  opts.GroupID,
			"folder_id": opts.FolderID,
			"user_id":   opts.UserID,
		})
		return fmt.Errorf("failed to list files: %w", err)
	}

	printVerbose("Found %d files:\n\n", len(files))
	logger.LogCommandComplete(ctx, "files.list", len(files))
	return formatEmptyOrOutput(files, "No files found")
}

func runFilesGet(ctx context.Context, client *api.Client, opts *options.FilesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.get", map[string]interface{}{
		"file_id": opts.FileID,
	})

	filesService := api.NewFilesService(client)

	file, err := filesService.Get(ctx, opts.FileID, opts.Include)
	if err != nil {
		logger.LogCommandError(ctx, "files.get", err, map[string]interface{}{
			"file_id": opts.FileID,
		})
		return fmt.Errorf("failed to get file: %w", err)
	}

	logger.LogCommandComplete(ctx, "files.get", 1)
	return formatOutput(file, nil)
}

func runFilesUpload(ctx context.Context, client *api.Client, opts *options.FilesUploadOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.upload", map[string]interface{}{
		"file_path": opts.FilePath,
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
		"folder_id": opts.FolderID,
		"user_id":   opts.UserID,
	})

	// Check if file exists
	if _, err := os.Stat(opts.FilePath); os.IsNotExist(err) {
		logger.LogCommandError(ctx, "files.upload", err, map[string]interface{}{
			"file_path": opts.FilePath,
		})
		return fmt.Errorf("file does not exist: %s", opts.FilePath)
	}

	filesService := api.NewFilesService(client)

	params := &api.UploadParams{
		OnDuplicate:    opts.OnDuplicate,
		ParentFolderID: opts.ParentFolderID,
		Hidden:         opts.Hidden,
		Locked:         opts.Locked,
	}

	fmt.Printf("Uploading %s...\n", filepath.Base(opts.FilePath))

	var uploadedFile *api.Attachment
	var err error

	switch {
	case opts.CourseID > 0:
		uploadedFile, err = filesService.UploadToCourse(ctx, opts.CourseID, opts.FilePath, params)
	case opts.GroupID > 0:
		// Canvas groups use the same upload path structure as courses
		uploadedFile, err = filesService.UploadToFolder(ctx, opts.GroupID, opts.FilePath, params)
	case opts.FolderID > 0:
		uploadedFile, err = filesService.UploadToFolder(ctx, opts.FolderID, opts.FilePath, params)
	default:
		uploadedFile, err = filesService.UploadToUser(ctx, opts.UserID, opts.FilePath, params)
	}

	if err != nil {
		logger.LogCommandError(ctx, "files.upload", err, map[string]interface{}{
			"file_path": opts.FilePath,
			"course_id": opts.CourseID,
			"group_id":  opts.GroupID,
			"folder_id": opts.FolderID,
			"user_id":   opts.UserID,
		})
		return fmt.Errorf("failed to upload file: %w", err)
	}

	printInfo("✅ File uploaded successfully\n\n")
	printInfo("   ID: %d\n", uploadedFile.ID)
	printInfo("   Name: %s\n", uploadedFile.DisplayName)
	printInfo("   Size: %s\n", formatFileSize(uploadedFile.Size))
	printInfo("   URL: %s\n", uploadedFile.URL)

	logger.LogCommandComplete(ctx, "files.upload", 1)
	return nil
}

func runFilesDownload(ctx context.Context, client *api.Client, opts *options.FilesDownloadOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.download", map[string]interface{}{
		"file_id":     opts.FileID,
		"destination": opts.Destination,
	})

	filesService := api.NewFilesService(client)

	// Get file info first to get the filename
	file, err := filesService.Get(ctx, opts.FileID, nil)
	if err != nil {
		logger.LogCommandError(ctx, "files.download", err, map[string]interface{}{
			"file_id": opts.FileID,
		})
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Determine destination path
	destPath := opts.Destination
	if destPath == "" {
		// Sanitize the server-controlled filename to prevent path traversal.
		// filepath.Base strips any directory components; the additional checks
		// below reject names that would still be unsafe after cleaning.
		clean := filepath.Base(file.Filename)
		if clean == "" || clean == "." || clean == "/" || strings.ContainsAny(clean, "/\\") {
			return fmt.Errorf("server returned an unsafe filename: %q", file.Filename)
		}
		destPath = clean
	}

	fmt.Printf("Downloading %s...\n", file.DisplayName)

	// Download file
	if err := filesService.Download(ctx, opts.FileID, destPath); err != nil {
		logger.LogCommandError(ctx, "files.download", err, map[string]interface{}{
			"file_id":     opts.FileID,
			"destination": destPath,
		})
		return fmt.Errorf("failed to download file: %w", err)
	}

	fmt.Printf("✅ File downloaded to %s\n", destPath)

	logger.LogCommandComplete(ctx, "files.download", 1)
	return nil
}

func runFilesDelete(ctx context.Context, client *api.Client, opts *options.FilesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.delete", map[string]interface{}{
		"file_id": opts.FileID,
		"force":   opts.Force,
	})

	// Confirm deletion
	confirmed, err := confirmDelete("file", opts.FileID, opts.Force)
	if err != nil {
		logger.LogCommandError(ctx, "files.delete", err, map[string]interface{}{
			"file_id": opts.FileID,
		})
		return err
	}
	if !confirmed {
		fmt.Println("Delete cancelled")
		logger.LogCommandComplete(ctx, "files.delete", 0)
		return nil
	}

	filesService := api.NewFilesService(client)

	if err := filesService.Delete(ctx, opts.FileID); err != nil {
		logger.LogCommandError(ctx, "files.delete", err, map[string]interface{}{
			"file_id": opts.FileID,
		})
		return fmt.Errorf("failed to delete file: %w", err)
	}

	printInfo("✅ File %d deleted successfully\n", opts.FileID)

	logger.LogCommandComplete(ctx, "files.delete", 1)
	return nil
}

func runFilesQuota(ctx context.Context, client *api.Client, opts *options.FilesQuotaOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.quota", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
		"user_id":   opts.UserID,
	})

	filesService := api.NewFilesService(client)

	var quota *api.QuotaInfo
	var err error

	switch {
	case opts.CourseID > 0:
		quota, err = filesService.GetCourseQuota(ctx, opts.CourseID)
	case opts.GroupID > 0:
		quota, err = filesService.GetGroupQuota(ctx, opts.GroupID)
	default:
		quota, err = filesService.GetUserQuota(ctx, opts.UserID)
	}

	if err != nil {
		logger.LogCommandError(ctx, "files.quota", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"group_id":  opts.GroupID,
			"user_id":   opts.UserID,
		})
		return fmt.Errorf("failed to get quota: %w", err)
	}

	fmt.Println("Storage Quota:")
	fmt.Printf("   Used: %s\n", formatFileSize(quota.QuotaUsed))
	fmt.Printf("   Total: %s\n", formatFileSize(quota.Quota))

	if quota.Quota > 0 {
		percentUsed := float64(quota.QuotaUsed) / float64(quota.Quota) * 100
		fmt.Printf("   Usage: %.1f%%\n", percentUsed)
	}

	logger.LogCommandComplete(ctx, "files.quota", 1)
	return nil
}

func runFilesResetVerifier(ctx context.Context, client *api.Client, opts *options.FilesResetVerifierOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.reset-verifier", map[string]interface{}{
		"file_id": opts.FileID,
	})

	filesService := api.NewFilesService(client)

	file, err := filesService.ResetVerifier(ctx, opts.FileID)
	if err != nil {
		logger.LogCommandError(ctx, "files.reset-verifier", err, map[string]interface{}{
			"file_id": opts.FileID,
		})
		return fmt.Errorf("failed to reset verifier: %w", err)
	}

	logger.LogCommandComplete(ctx, "files.reset-verifier", 1)
	return formatSuccessOutput(file, "Link verifier reset successfully!")
}

func runFilesCopy(ctx context.Context, client *api.Client, opts *options.FilesCopyOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.copy", map[string]interface{}{
		"dest_folder_id": opts.DestFolderID,
		"source_file_id": opts.SourceFileID,
	})

	filesService := api.NewFilesService(client)

	file, err := filesService.CopyFile(ctx, opts.DestFolderID, &api.CopyFileParams{
		SourceFileID: opts.SourceFileID,
		OnDuplicate:  opts.OnDuplicate,
	})
	if err != nil {
		logger.LogCommandError(ctx, "files.copy", err, map[string]interface{}{
			"dest_folder_id": opts.DestFolderID,
			"source_file_id": opts.SourceFileID,
		})
		return fmt.Errorf("failed to copy file: %w", err)
	}

	logger.LogCommandComplete(ctx, "files.copy", 1)
	return formatSuccessOutput(file, "File copied successfully!")
}

func runFilesSetUsageRights(ctx context.Context, client *api.Client, opts *options.FilesUsageRightsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.set-usage-rights", map[string]interface{}{
		"course_id":         opts.CourseID,
		"group_id":          opts.GroupID,
		"user_id":           opts.UserID,
		"use_justification": opts.UseJustification,
	})

	filesService := api.NewFilesService(client)

	params := &api.SetUsageRightsParams{
		FileIDs:          opts.FileIDs,
		FolderIDs:        opts.FolderIDs,
		Publish:          opts.Publish,
		UseJustification: opts.UseJustification,
		LegalCopyright:   opts.LegalCopyright,
		License:          opts.License,
	}

	rights, err := filesService.SetUsageRights(ctx, opts.CourseID, opts.GroupID, opts.UserID, params)
	if err != nil {
		logger.LogCommandError(ctx, "files.set-usage-rights", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"group_id":  opts.GroupID,
			"user_id":   opts.UserID,
		})
		return fmt.Errorf("failed to set usage rights: %w", err)
	}

	logger.LogCommandComplete(ctx, "files.set-usage-rights", 1)
	return formatSuccessOutput(rights, "Usage rights set successfully!")
}

func runFilesRemoveUsageRights(ctx context.Context, client *api.Client, opts *options.FilesRemoveUsageRightsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.remove-usage-rights", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
		"user_id":   opts.UserID,
	})

	filesService := api.NewFilesService(client)

	params := &api.RemoveUsageRightsParams{
		FileIDs:   opts.FileIDs,
		FolderIDs: opts.FolderIDs,
	}

	if err := filesService.RemoveUsageRights(ctx, opts.CourseID, opts.GroupID, opts.UserID, params); err != nil {
		logger.LogCommandError(ctx, "files.remove-usage-rights", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"group_id":  opts.GroupID,
			"user_id":   opts.UserID,
		})
		return fmt.Errorf("failed to remove usage rights: %w", err)
	}

	logger.LogCommandComplete(ctx, "files.remove-usage-rights", 1)
	printInfo("Usage rights removed successfully\n")
	return nil
}

func runFilesLicenses(ctx context.Context, client *api.Client, opts *options.FilesLicensesOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.licenses", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
		"user_id":   opts.UserID,
	})

	filesService := api.NewFilesService(client)

	licenses, err := filesService.ListLicenses(ctx, opts.CourseID, opts.GroupID, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "files.licenses", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"group_id":  opts.GroupID,
			"user_id":   opts.UserID,
		})
		return fmt.Errorf("failed to list licenses: %w", err)
	}

	logger.LogCommandComplete(ctx, "files.licenses", len(licenses))
	return formatEmptyOrOutput(licenses, "No licenses found")
}

// formatFileSize formats a file size in bytes to a human-readable string
func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
