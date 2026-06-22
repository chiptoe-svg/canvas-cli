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

// courseNicknamesCmd represents the course-nicknames command group
var courseNicknamesCmd = &cobra.Command{
	Use:   "course-nicknames",
	Short: "Manage course nicknames",
	Long: `Manage nicknames for Canvas courses for the current user.

Examples:
  canvas course-nicknames list
  canvas course-nicknames get 123
  canvas course-nicknames set 123 --nickname "My Fav Course"
  canvas course-nicknames delete 123
  canvas course-nicknames delete-all`,
}

func init() {
	rootCmd.AddCommand(courseNicknamesCmd)
	courseNicknamesCmd.AddCommand(newCourseNicknamesListCmd())
	courseNicknamesCmd.AddCommand(newCourseNicknamesGetCmd())
	courseNicknamesCmd.AddCommand(newCourseNicknamesSetCmd())
	courseNicknamesCmd.AddCommand(newCourseNicknamesDeleteCmd())
	courseNicknamesCmd.AddCommand(newCourseNicknamesDeleteAllCmd())
}

func newCourseNicknamesListCmd() *cobra.Command {
	opts := &options.CourseNicknamesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all course nicknames",
		Long: `List all course nicknames for the current user.

Examples:
  canvas course-nicknames list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseNicknamesList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newCourseNicknamesGetCmd() *cobra.Command {
	opts := &options.CourseNicknamesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <course-id>",
		Short: "Get the nickname for a course",
		Long: `Get the nickname set for a specific course.

Examples:
  canvas course-nicknames get 123`,
		Args: ExactArgsWithUsage(1, "course-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			courseID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid course ID: %s", args[0])
			}
			opts.CourseID = courseID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseNicknamesGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newCourseNicknamesSetCmd() *cobra.Command {
	opts := &options.CourseNicknamesSetOptions{}

	cmd := &cobra.Command{
		Use:   "set <course-id>",
		Short: "Set a nickname for a course",
		Long: `Set or update the nickname for a specific course.

Examples:
  canvas course-nicknames set 123 --nickname "My Fav Course"`,
		Args: ExactArgsWithUsage(1, "course-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			courseID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid course ID: %s", args[0])
			}
			opts.CourseID = courseID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseNicknamesSet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Nickname, "nickname", "", "Nickname to set for the course (required)")
	mustMarkRequired(cmd, "nickname")

	return cmd
}

func newCourseNicknamesDeleteCmd() *cobra.Command {
	opts := &options.CourseNicknamesDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <course-id>",
		Short: "Delete the nickname for a course",
		Long: `Remove the nickname set for a specific course.

Examples:
  canvas course-nicknames delete 123`,
		Args: ExactArgsWithUsage(1, "course-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			courseID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid course ID: %s", args[0])
			}
			opts.CourseID = courseID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseNicknamesDelete(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newCourseNicknamesDeleteAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-all",
		Short: "Delete all course nicknames",
		Long: `Remove all course nicknames for the current user.

Examples:
  canvas course-nicknames delete-all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseNicknamesDeleteAll(cmd.Context(), client)
		},
	}

	return cmd
}

func runCourseNicknamesList(ctx context.Context, client *api.Client, opts *options.CourseNicknamesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-nicknames.list", map[string]interface{}{})

	svc := api.NewCourseNicknamesService(client)

	nicknames, err := svc.List(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "course-nicknames.list", err, map[string]interface{}{})
		return fmt.Errorf("failed to list course nicknames: %w", err)
	}

	logger.LogCommandComplete(ctx, "course-nicknames.list", len(nicknames))
	return formatEmptyOrOutput(nicknames, "No course nicknames found")
}

func runCourseNicknamesGet(ctx context.Context, client *api.Client, opts *options.CourseNicknamesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-nicknames.get", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	svc := api.NewCourseNicknamesService(client)

	nickname, err := svc.Get(ctx, opts.CourseID)
	if err != nil {
		logger.LogCommandError(ctx, "course-nicknames.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to get course nickname: %w", err)
	}

	logger.LogCommandComplete(ctx, "course-nicknames.get", 1)
	return formatOutput(nickname, nil)
}

func runCourseNicknamesSet(ctx context.Context, client *api.Client, opts *options.CourseNicknamesSetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-nicknames.set", map[string]interface{}{
		"course_id": opts.CourseID,
		"nickname":  opts.Nickname,
	})

	svc := api.NewCourseNicknamesService(client)

	params := api.SetCourseNicknameParams{
		Nickname: opts.Nickname,
	}

	nickname, err := svc.Set(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "course-nicknames.set", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to set course nickname: %w", err)
	}

	logger.LogCommandComplete(ctx, "course-nicknames.set", 1)
	return formatSuccessOutput(nickname, "Course nickname set successfully!")
}

func runCourseNicknamesDelete(ctx context.Context, client *api.Client, opts *options.CourseNicknamesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-nicknames.delete", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	svc := api.NewCourseNicknamesService(client)

	if err := svc.Delete(ctx, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "course-nicknames.delete", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to delete course nickname: %w", err)
	}

	logger.LogCommandComplete(ctx, "course-nicknames.delete", 1)
	printInfo("Course nickname for course %d deleted successfully\n", opts.CourseID)
	return nil
}

func runCourseNicknamesDeleteAll(ctx context.Context, client *api.Client) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-nicknames.delete-all", map[string]interface{}{})

	svc := api.NewCourseNicknamesService(client)

	if err := svc.DeleteAll(ctx); err != nil {
		logger.LogCommandError(ctx, "course-nicknames.delete-all", err, map[string]interface{}{})
		return fmt.Errorf("failed to delete all course nicknames: %w", err)
	}

	logger.LogCommandComplete(ctx, "course-nicknames.delete-all", 0)
	printInfo("All course nicknames deleted successfully\n")
	return nil
}
