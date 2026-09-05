package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/commands/internal/options"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
	"github.com/chiptoe-svg/canvas-cli/internal/progress"
)

// coursesCmd represents the courses command group
var coursesCmd = &cobra.Command{
	Use:   "courses",
	Short: "Manage Canvas courses",
	Long: `Manage Canvas courses including listing, viewing, and updating courses.

Examples:
  canvas courses list
  canvas courses get 123
  canvas courses list --enrollment-type teacher
  canvas courses list --state available`,
}

// newCoursesListCmd creates the courses list command
func newCoursesListCmd() *cobra.Command {
	opts := &options.CoursesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List courses",
		Long: `List the courses you are enrolled in.

Examples:
  canvas courses list
  canvas courses list --enrollment-type teacher
  canvas courses list --enrollment-type student
  canvas courses list --enrollment-state active
  canvas courses list --state available
  canvas courses list --include syllabus_body,term`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCoursesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.EnrollmentType, "enrollment-type", "", "Filter by enrollment type (student, teacher, ta, observer, designer)")
	cmd.Flags().StringVar(&opts.EnrollmentState, "enrollment-state", "", "Filter by enrollment state (active, invited_or_pending, completed)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")
	cmd.Flags().StringSliceVar(&opts.State, "state", []string{}, "Filter by course state (comma-separated: available, completed, unpublished, deleted)")

	return cmd
}

func runCoursesList(ctx context.Context, client *api.Client, opts *options.CoursesListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "courses.list", map[string]interface{}{
		"enrollment_type":  opts.EnrollmentType,
		"enrollment_state": opts.EnrollmentState,
	})

	spin := progress.New("Fetching courses...")
	if !quiet {
		spin.Start()
	}

	coursesService := api.NewCoursesService(client)

	reqOpts := &api.ListCoursesOptions{
		EnrollmentType:  opts.EnrollmentType,
		EnrollmentState: opts.EnrollmentState,
		Include:         opts.Include,
		State:           opts.State,
	}

	courses, err := coursesService.List(ctx, reqOpts)
	spin.Stop()
	if err != nil {
		logger.LogCommandError(ctx, "courses.list", err, map[string]interface{}{
			"enrollment_type": opts.EnrollmentType,
		})
		return fmt.Errorf("failed to list courses: %w", err)
	}

	printVerbose("Found %d enrolled courses:\n\n", len(courses))

	// Format and display courses
	if err := formatEmptyOrOutput(courses, "No courses found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "courses.list", len(courses))
	return nil
}

// newCoursesGetCmd creates the courses get command
func newCoursesGetCmd() *cobra.Command {
	opts := &options.CoursesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <course-id>",
		Short: "Get details of a specific course",
		Long: `Get details of a specific course by ID.

Examples:
  canvas courses get 123
  canvas courses get 123 --include syllabus_body,term`,
		Args: ExactArgsWithUsage(1, "course-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse course ID
			var courseID int64
			if _, err := fmt.Sscanf(args[0], "%d", &courseID); err != nil {
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

			return runCoursesGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")

	return cmd
}

func runCoursesGet(ctx context.Context, client *api.Client, opts *options.CoursesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "courses.get", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	// Create courses service
	coursesService := api.NewCoursesService(client)

	// Get course
	course, err := coursesService.Get(ctx, opts.CourseID, opts.Include)
	if err != nil {
		logger.LogCommandError(ctx, "courses.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to get course: %w", err)
	}

	// Format and display course details
	if err := formatOutput(course, nil); err != nil {
		return fmt.Errorf("failed to print result: %w", err)
	}

	logger.LogCommandComplete(ctx, "courses.get", 1)
	return nil
}

// newCoursesUpdateCmd creates the courses update command
func newCoursesUpdateCmd() *cobra.Command {
	opts := &options.CoursesUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <course-id>",
		Short: "Update a course",
		Long: `Update an existing course.

Examples:
  canvas courses update 123 --name "Updated Course Name"
  canvas courses update 123 --code "NEW101" --start-at "2024-10-01"
  canvas courses update 123 --public
  canvas courses update 123 --offer`,
		Args: ExactArgsWithUsage(1, "course-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse course ID
			var courseID int64
			if _, err := fmt.Sscanf(args[0], "%d", &courseID); err != nil {
				return fmt.Errorf("invalid course ID: %s", args[0])
			}
			opts.CourseID = courseID

			// Handle boolean flags that were explicitly set
			if cmd.Flags().Changed("public") {
				public := cmd.Flags().Lookup("public").Value.String() == "true"
				opts.IsPublic = &public
			}

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCoursesUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Course name")
	cmd.Flags().StringVar(&opts.CourseCode, "code", "", "Course code")
	cmd.Flags().StringVar(&opts.StartAt, "start-at", "", "Start date (ISO 8601)")
	cmd.Flags().StringVar(&opts.EndAt, "end-at", "", "End date (ISO 8601)")
	cmd.Flags().StringVar(&opts.License, "license", "", "Course license")
	cmd.Flags().Bool("public", false, "Make course public")
	cmd.Flags().StringVar(&opts.DefaultView, "default-view", "", "Default view (feed, wiki, modules, syllabus, assignments)")

	return cmd
}

func runCoursesUpdate(ctx context.Context, client *api.Client, opts *options.CoursesUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "courses.update", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	// Create courses service
	coursesService := api.NewCoursesService(client)

	// Build params - only include changed values
	params := &api.UpdateCourseParams{
		Name:        opts.Name,
		CourseCode:  opts.CourseCode,
		StartAt:     opts.StartAt,
		EndAt:       opts.EndAt,
		License:     opts.License,
		DefaultView: opts.DefaultView,
		IsPublic:    opts.IsPublic,
	}

	// Update course
	course, err := coursesService.Update(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "courses.update", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to update course: %w", err)
	}

	printInfo("Course updated successfully (ID: %d)\n", course.ID)
	if err := formatOutput(course, nil); err != nil {
		return fmt.Errorf("failed to print result: %w", err)
	}

	logger.LogCommandComplete(ctx, "courses.update", 1)
	return nil
}

func init() {
	rootCmd.AddCommand(coursesCmd)
	coursesCmd.AddCommand(newCoursesListCmd())
	coursesCmd.AddCommand(newCoursesGetCmd())
	coursesCmd.AddCommand(newCoursesUpdateCmd())
}
