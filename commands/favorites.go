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

// favoritesCmd represents the favorites command group
var favoritesCmd = &cobra.Command{
	Use:   "favorites",
	Short: "Manage Canvas favorites",
	Long: `Manage favorite courses and groups for the current user.

Examples:
  canvas favorites courses list
  canvas favorites courses add 123
  canvas favorites courses remove 123
  canvas favorites courses reset
  canvas favorites groups list
  canvas favorites groups add 456
  canvas favorites groups remove 456
  canvas favorites groups reset`,
}

// favoritesCourseCmd represents the favorites courses subcommand group
var favoritesCourseCmd = &cobra.Command{
	Use:   "courses",
	Short: "Manage favorite courses",
	Long:  `Manage favorite courses for the current user.`,
}

// favoritesGroupCmd represents the favorites groups subcommand group
var favoritesGroupCmd = &cobra.Command{
	Use:   "groups",
	Short: "Manage favorite groups",
	Long:  `Manage favorite groups for the current user.`,
}

func init() {
	rootCmd.AddCommand(favoritesCmd)
	favoritesCmd.AddCommand(favoritesCourseCmd)
	favoritesCmd.AddCommand(favoritesGroupCmd)

	favoritesCourseCmd.AddCommand(newFavoritesCoursesListCmd())
	favoritesCourseCmd.AddCommand(newFavoritesCoursesAddCmd())
	favoritesCourseCmd.AddCommand(newFavoritesCoursesRemoveCmd())
	favoritesCourseCmd.AddCommand(newFavoritesCoursesResetCmd())

	favoritesGroupCmd.AddCommand(newFavoritesGroupsListCmd())
	favoritesGroupCmd.AddCommand(newFavoritesGroupsAddCmd())
	favoritesGroupCmd.AddCommand(newFavoritesGroupsRemoveCmd())
	favoritesGroupCmd.AddCommand(newFavoritesGroupsResetCmd())
}

func newFavoritesCoursesListCmd() *cobra.Command {
	opts := &options.FavoritesListCoursesOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List favorite courses",
		Long: `List all favorite courses for the current user.

Examples:
  canvas favorites courses list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFavoritesListCourses(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newFavoritesCoursesAddCmd() *cobra.Command {
	opts := &options.FavoritesAddCourseOptions{}

	cmd := &cobra.Command{
		Use:   "add <id>",
		Short: "Add a course to favorites",
		Long: `Add a course to the current user's favorites.

Examples:
  canvas favorites courses add 123`,
		Args: ExactArgsWithUsage(1, "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid course ID: %s", args[0])
			}
			opts.ID = id
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFavoritesAddCourse(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newFavoritesCoursesRemoveCmd() *cobra.Command {
	opts := &options.FavoritesRemoveCourseOptions{}

	cmd := &cobra.Command{
		Use:   "remove <id>",
		Short: "Remove a course from favorites",
		Long: `Remove a course from the current user's favorites.

Examples:
  canvas favorites courses remove 123`,
		Args: ExactArgsWithUsage(1, "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid course ID: %s", args[0])
			}
			opts.ID = id
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFavoritesRemoveCourse(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newFavoritesCoursesResetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset favorite courses to defaults",
		Long: `Reset all favorite courses to the system defaults.

Examples:
  canvas favorites courses reset`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFavoritesResetCourses(cmd.Context(), client)
		},
	}

	return cmd
}

func newFavoritesGroupsListCmd() *cobra.Command {
	opts := &options.FavoritesListGroupsOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List favorite groups",
		Long: `List all favorite groups for the current user.

Examples:
  canvas favorites groups list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFavoritesListGroups(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newFavoritesGroupsAddCmd() *cobra.Command {
	opts := &options.FavoritesAddGroupOptions{}

	cmd := &cobra.Command{
		Use:   "add <id>",
		Short: "Add a group to favorites",
		Long: `Add a group to the current user's favorites.

Examples:
  canvas favorites groups add 456`,
		Args: ExactArgsWithUsage(1, "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %s", args[0])
			}
			opts.ID = id
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFavoritesAddGroup(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newFavoritesGroupsRemoveCmd() *cobra.Command {
	opts := &options.FavoritesRemoveGroupOptions{}

	cmd := &cobra.Command{
		Use:   "remove <id>",
		Short: "Remove a group from favorites",
		Long: `Remove a group from the current user's favorites.

Examples:
  canvas favorites groups remove 456`,
		Args: ExactArgsWithUsage(1, "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %s", args[0])
			}
			opts.ID = id
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFavoritesRemoveGroup(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newFavoritesGroupsResetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset favorite groups to defaults",
		Long: `Reset all favorite groups to the system defaults.

Examples:
  canvas favorites groups reset`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runFavoritesResetGroups(cmd.Context(), client)
		},
	}

	return cmd
}

func runFavoritesListCourses(ctx context.Context, client *api.Client, opts *options.FavoritesListCoursesOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "favorites.courses.list", map[string]interface{}{})

	svc := api.NewFavoritesService(client)

	courses, err := svc.ListCourses(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "favorites.courses.list", err, map[string]interface{}{})
		return fmt.Errorf("failed to list favorite courses: %w", err)
	}

	logger.LogCommandComplete(ctx, "favorites.courses.list", len(courses))
	return formatEmptyOrOutput(courses, "No favorite courses found")
}

func runFavoritesAddCourse(ctx context.Context, client *api.Client, opts *options.FavoritesAddCourseOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "favorites.courses.add", map[string]interface{}{
		"id": opts.ID,
	})

	svc := api.NewFavoritesService(client)

	course, err := svc.AddCourse(ctx, opts.ID)
	if err != nil {
		logger.LogCommandError(ctx, "favorites.courses.add", err, map[string]interface{}{
			"id": opts.ID,
		})
		return fmt.Errorf("failed to add course to favorites: %w", err)
	}

	logger.LogCommandComplete(ctx, "favorites.courses.add", 1)
	return formatSuccessOutput(course, "Course added to favorites!")
}

func runFavoritesRemoveCourse(ctx context.Context, client *api.Client, opts *options.FavoritesRemoveCourseOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "favorites.courses.remove", map[string]interface{}{
		"id": opts.ID,
	})

	svc := api.NewFavoritesService(client)

	_, err := svc.RemoveCourse(ctx, opts.ID)
	if err != nil {
		logger.LogCommandError(ctx, "favorites.courses.remove", err, map[string]interface{}{
			"id": opts.ID,
		})
		return fmt.Errorf("failed to remove course from favorites: %w", err)
	}

	logger.LogCommandComplete(ctx, "favorites.courses.remove", 1)
	printInfo("Course %d removed from favorites\n", opts.ID)
	return nil
}

func runFavoritesResetCourses(ctx context.Context, client *api.Client) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "favorites.courses.reset", map[string]interface{}{})

	svc := api.NewFavoritesService(client)

	courses, err := svc.ResetCourses(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "favorites.courses.reset", err, map[string]interface{}{})
		return fmt.Errorf("failed to reset favorite courses: %w", err)
	}

	logger.LogCommandComplete(ctx, "favorites.courses.reset", len(courses))
	return formatSuccessOutput(courses, "Favorite courses reset to defaults!")
}

func runFavoritesListGroups(ctx context.Context, client *api.Client, opts *options.FavoritesListGroupsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "favorites.groups.list", map[string]interface{}{})

	svc := api.NewFavoritesService(client)

	groups, err := svc.ListGroups(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "favorites.groups.list", err, map[string]interface{}{})
		return fmt.Errorf("failed to list favorite groups: %w", err)
	}

	logger.LogCommandComplete(ctx, "favorites.groups.list", len(groups))
	return formatEmptyOrOutput(groups, "No favorite groups found")
}

func runFavoritesAddGroup(ctx context.Context, client *api.Client, opts *options.FavoritesAddGroupOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "favorites.groups.add", map[string]interface{}{
		"id": opts.ID,
	})

	svc := api.NewFavoritesService(client)

	group, err := svc.AddGroup(ctx, opts.ID)
	if err != nil {
		logger.LogCommandError(ctx, "favorites.groups.add", err, map[string]interface{}{
			"id": opts.ID,
		})
		return fmt.Errorf("failed to add group to favorites: %w", err)
	}

	logger.LogCommandComplete(ctx, "favorites.groups.add", 1)
	return formatSuccessOutput(group, "Group added to favorites!")
}

func runFavoritesRemoveGroup(ctx context.Context, client *api.Client, opts *options.FavoritesRemoveGroupOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "favorites.groups.remove", map[string]interface{}{
		"id": opts.ID,
	})

	svc := api.NewFavoritesService(client)

	_, err := svc.RemoveGroup(ctx, opts.ID)
	if err != nil {
		logger.LogCommandError(ctx, "favorites.groups.remove", err, map[string]interface{}{
			"id": opts.ID,
		})
		return fmt.Errorf("failed to remove group from favorites: %w", err)
	}

	logger.LogCommandComplete(ctx, "favorites.groups.remove", 1)
	printInfo("Group %d removed from favorites\n", opts.ID)
	return nil
}

func runFavoritesResetGroups(ctx context.Context, client *api.Client) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "favorites.groups.reset", map[string]interface{}{})

	svc := api.NewFavoritesService(client)

	groups, err := svc.ResetGroups(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "favorites.groups.reset", err, map[string]interface{}{})
		return fmt.Errorf("failed to reset favorite groups: %w", err)
	}

	logger.LogCommandComplete(ctx, "favorites.groups.reset", len(groups))
	return formatSuccessOutput(groups, "Favorite groups reset to defaults!")
}
