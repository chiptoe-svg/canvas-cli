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

// mediaCmd represents the media command group.
var mediaCmd = &cobra.Command{
	Use:   "media",
	Short: "Manage Canvas media objects and attachments",
	Long: `List and update Canvas media objects (videos/audio) and media attachments.
Also manage associated media tracks (captions/subtitles).

Examples:
  canvas media objects list
  canvas media objects list --course-id 123
  canvas media objects update m-abc123 --title "Intro Video"
  canvas media objects tracks m-abc123
  canvas media attachments list
  canvas media attachments list --course-id 123
  canvas media attachments tracks 5`,
}

// mediaObjectsCmd is the subgroup for media objects.
var mediaObjectsCmd = &cobra.Command{
	Use:   "objects",
	Short: "Manage Canvas media objects",
}

// mediaAttachmentsCmd is the subgroup for media attachments.
var mediaAttachmentsCmd = &cobra.Command{
	Use:   "attachments",
	Short: "Manage Canvas media attachments",
}

func init() {
	rootCmd.AddCommand(mediaCmd)
	mediaCmd.AddCommand(mediaObjectsCmd)
	mediaCmd.AddCommand(mediaAttachmentsCmd)

	mediaObjectsCmd.AddCommand(newMediaObjectsListCmd())
	mediaObjectsCmd.AddCommand(newMediaObjectUpdateCmd())
	mediaObjectsCmd.AddCommand(newMediaObjectTracksCmd())

	mediaAttachmentsCmd.AddCommand(newMediaAttachmentsListCmd())
	mediaAttachmentsCmd.AddCommand(newMediaAttachmentUpdateCmd())
	mediaAttachmentsCmd.AddCommand(newMediaAttachmentTracksCmd())
}

// --- media objects ---

func newMediaObjectsListCmd() *cobra.Command {
	opts := &options.MediaObjectsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List media objects",
		Long: `List media objects for the current user, a course, or a group.

Examples:
  canvas media objects list
  canvas media objects list --course-id 123
  canvas media objects list --group-id 456 --sort title`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runMediaObjectsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Group ID")
	cmd.Flags().StringVar(&opts.Sort, "sort", "", "Sort field: title, created_at, updated_at, user_name")
	cmd.Flags().StringVar(&opts.Order, "order", "", "Sort order: asc, desc")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Number of results per page")

	return cmd
}

func newMediaObjectUpdateCmd() *cobra.Command {
	opts := &options.MediaObjectUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <media-id>",
		Short: "Update a media object",
		Long: `Update the title of a media object.

Examples:
  canvas media objects update m-abc123 --title "Updated Title"`,
		Args: ExactArgsWithUsage(1, "media-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.MediaID = args[0]

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runMediaObjectUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Title, "title", "", "New title for the media object")
	mustMarkRequired(cmd, "title")

	return cmd
}

func newMediaObjectTracksCmd() *cobra.Command {
	opts := &options.MediaTracksListOptions{}

	cmd := &cobra.Command{
		Use:   "tracks <media-id>",
		Short: "List media tracks for a media object",
		Long: `List the caption/subtitle tracks for a media object.

Examples:
  canvas media objects tracks m-abc123`,
		Args: ExactArgsWithUsage(1, "media-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.MediaID = args[0]

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runMediaObjectTracks(cmd.Context(), client, opts)
		},
	}

	return cmd
}

// --- media attachments ---

func newMediaAttachmentsListCmd() *cobra.Command {
	opts := &options.MediaAttachmentsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List media attachments",
		Long: `List media attachments for the current user, a course, or a group.

Examples:
  canvas media attachments list
  canvas media attachments list --course-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runMediaAttachmentsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Group ID")
	cmd.Flags().StringVar(&opts.Sort, "sort", "", "Sort field")
	cmd.Flags().StringVar(&opts.Order, "order", "", "Sort order: asc, desc")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Number of results per page")

	return cmd
}

func newMediaAttachmentUpdateCmd() *cobra.Command {
	opts := &options.MediaAttachmentUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <attachment-id>",
		Short: "Update a media attachment title",
		Long: `Update the title of a media attachment.

Examples:
  canvas media attachments update 5 --title "New Title"`,
		Args: ExactArgsWithUsage(1, "attachment-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid attachment ID: %s", args[0])
			}
			opts.AttachmentID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runMediaAttachmentUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Title, "title", "", "New title for the media attachment")
	mustMarkRequired(cmd, "title")

	return cmd
}

func newMediaAttachmentTracksCmd() *cobra.Command {
	opts := &options.MediaTracksListOptions{}

	cmd := &cobra.Command{
		Use:   "tracks <attachment-id>",
		Short: "List media tracks for an attachment",
		Long: `List the caption/subtitle tracks for a media attachment.

Examples:
  canvas media attachments tracks 5`,
		Args: ExactArgsWithUsage(1, "attachment-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid attachment ID: %s", args[0])
			}
			opts.AttachmentID = id

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runMediaAttachmentTracks(cmd.Context(), client, opts)
		},
	}

	return cmd
}

// --- run functions ---

func runMediaObjectsList(ctx context.Context, client *api.Client, opts *options.MediaObjectsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "media.objects.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
	})

	svc := api.NewMediaObjectsService(client)
	apiOpts := &api.ListMediaObjectsOptions{
		Sort:    opts.Sort,
		Order:   opts.Order,
		PerPage: opts.PerPage,
	}

	var (
		objs []api.MediaObject
		err  error
	)

	switch {
	case opts.CourseID > 0:
		objs, err = svc.ListForCourse(ctx, opts.CourseID, apiOpts)
	case opts.GroupID > 0:
		objs, err = svc.ListForGroup(ctx, opts.GroupID, apiOpts)
	default:
		objs, err = svc.List(ctx, apiOpts)
	}

	if err != nil {
		logger.LogCommandError(ctx, "media.objects.list", err, nil)
		return fmt.Errorf("failed to list media objects: %w", err)
	}

	printVerbose("Found %d media objects:\n\n", len(objs))
	logger.LogCommandComplete(ctx, "media.objects.list", len(objs))
	return formatEmptyOrOutput(objs, "No media objects found")
}

func runMediaObjectUpdate(ctx context.Context, client *api.Client, opts *options.MediaObjectUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "media.objects.update", map[string]interface{}{
		"media_id": opts.MediaID,
		"title":    opts.Title,
	})

	svc := api.NewMediaObjectsService(client)

	obj, err := svc.Update(ctx, opts.MediaID, opts.Title)
	if err != nil {
		logger.LogCommandError(ctx, "media.objects.update", err, nil)
		return fmt.Errorf("failed to update media object: %w", err)
	}

	logger.LogCommandComplete(ctx, "media.objects.update", 1)
	return formatSuccessOutput(obj, "Media object updated successfully.")
}

func runMediaObjectTracks(ctx context.Context, client *api.Client, opts *options.MediaTracksListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "media.objects.tracks", map[string]interface{}{
		"media_id": opts.MediaID,
	})

	svc := api.NewMediaObjectsService(client)

	tracks, err := svc.GetMediaTracks(ctx, opts.MediaID)
	if err != nil {
		logger.LogCommandError(ctx, "media.objects.tracks", err, nil)
		return fmt.Errorf("failed to get media tracks: %w", err)
	}

	printVerbose("Found %d media tracks:\n\n", len(tracks))
	logger.LogCommandComplete(ctx, "media.objects.tracks", len(tracks))
	return formatEmptyOrOutput(tracks, "No media tracks found")
}

func runMediaAttachmentsList(ctx context.Context, client *api.Client, opts *options.MediaAttachmentsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "media.attachments.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
	})

	svc := api.NewMediaObjectsService(client)
	apiOpts := &api.ListMediaObjectsOptions{
		Sort:    opts.Sort,
		Order:   opts.Order,
		PerPage: opts.PerPage,
	}

	var (
		atts []api.MediaAttachment
		err  error
	)

	switch {
	case opts.CourseID > 0:
		atts, err = svc.ListAttachmentsForCourse(ctx, opts.CourseID, apiOpts)
	case opts.GroupID > 0:
		atts, err = svc.ListAttachmentsForGroup(ctx, opts.GroupID, apiOpts)
	default:
		atts, err = svc.ListAttachments(ctx, apiOpts)
	}

	if err != nil {
		logger.LogCommandError(ctx, "media.attachments.list", err, nil)
		return fmt.Errorf("failed to list media attachments: %w", err)
	}

	printVerbose("Found %d media attachments:\n\n", len(atts))
	logger.LogCommandComplete(ctx, "media.attachments.list", len(atts))
	return formatEmptyOrOutput(atts, "No media attachments found")
}

func runMediaAttachmentUpdate(ctx context.Context, client *api.Client, opts *options.MediaAttachmentUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "media.attachments.update", map[string]interface{}{
		"attachment_id": opts.AttachmentID,
		"title":         opts.Title,
	})

	svc := api.NewMediaObjectsService(client)

	att, err := svc.UpdateAttachment(ctx, opts.AttachmentID, opts.Title)
	if err != nil {
		logger.LogCommandError(ctx, "media.attachments.update", err, nil)
		return fmt.Errorf("failed to update media attachment: %w", err)
	}

	logger.LogCommandComplete(ctx, "media.attachments.update", 1)
	return formatSuccessOutput(att, "Media attachment updated successfully.")
}

func runMediaAttachmentTracks(ctx context.Context, client *api.Client, opts *options.MediaTracksListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "media.attachments.tracks", map[string]interface{}{
		"attachment_id": opts.AttachmentID,
	})

	svc := api.NewMediaObjectsService(client)

	tracks, err := svc.GetAttachmentMediaTracks(ctx, opts.AttachmentID)
	if err != nil {
		logger.LogCommandError(ctx, "media.attachments.tracks", err, nil)
		return fmt.Errorf("failed to get attachment media tracks: %w", err)
	}

	printVerbose("Found %d media tracks:\n\n", len(tracks))
	logger.LogCommandComplete(ctx, "media.attachments.tracks", len(tracks))
	return formatEmptyOrOutput(tracks, "No media tracks found")
}
