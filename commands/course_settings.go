package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// courseSettingsCmd is the parent command for course settings/utility operations.
var courseSettingsCmd = &cobra.Command{
	Use:     "settings",
	Aliases: []string{"csettings"},
	Short:   "Manage course settings and utilities",
	Long: `Manage Canvas course settings, tabs, permissions, due dates, and late policy.

Examples:
  canvas courses settings get --course-id 1
  canvas courses settings todo --course-id 1
  canvas courses settings tabs --course-id 1
  canvas courses settings permissions --course-id 1
  canvas courses settings effective-due-dates --course-id 1
  canvas courses settings late-policy --course-id 1
  canvas courses settings recent-students --course-id 1`,
}

func init() {
	coursesCmd.AddCommand(courseSettingsCmd)
	courseSettingsCmd.AddCommand(newCourseSettingsGetCmd())
	courseSettingsCmd.AddCommand(newCourseSettingsTodoCmd())
	courseSettingsCmd.AddCommand(newCourseSettingsTabsCmd())
	courseSettingsCmd.AddCommand(newCourseSettingsPermissionsCmd())
	courseSettingsCmd.AddCommand(newCourseSettingsEffectiveDueDatesCmd())
	courseSettingsCmd.AddCommand(newCourseSettingsLatePolicyCmd())
	courseSettingsCmd.AddCommand(newCourseSettingsRecentStudentsCmd())
}

func newCourseSettingsGetCmd() *cobra.Command {
	var courseID int64
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get settings for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseSettingsGet(cmd.Context(), client, courseID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runCourseSettingsGet(ctx context.Context, client *api.Client, courseID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-settings.get", map[string]interface{}{"course_id": courseID})

	svc := api.NewCourseSettingsService(client)
	settings, err := svc.GetSettings(ctx, courseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-settings.get", 1)
	return formatOutput(settings, nil)
}

func newCourseSettingsTodoCmd() *cobra.Command {
	var courseID int64
	cmd := &cobra.Command{
		Use:   "todo",
		Short: "List todo items for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseSettingsTodo(cmd.Context(), client, courseID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runCourseSettingsTodo(ctx context.Context, client *api.Client, courseID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-settings.todo", map[string]interface{}{"course_id": courseID})

	svc := api.NewCourseSettingsService(client)
	todos, err := svc.GetTodo(ctx, courseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-settings.todo", len(todos))
	return formatEmptyOrOutput(todos, "No todo items found")
}

func newCourseSettingsTabsCmd() *cobra.Command {
	var courseID int64
	cmd := &cobra.Command{
		Use:   "tabs",
		Short: "List navigation tabs for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseSettingsTabs(cmd.Context(), client, courseID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runCourseSettingsTabs(ctx context.Context, client *api.Client, courseID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-settings.tabs", map[string]interface{}{"course_id": courseID})

	svc := api.NewCourseSettingsService(client)
	tabs, err := svc.ListTabs(ctx, courseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-settings.tabs", len(tabs))
	return formatEmptyOrOutput(tabs, "No tabs found")
}

func newCourseSettingsPermissionsCmd() *cobra.Command {
	var courseID int64
	cmd := &cobra.Command{
		Use:   "permissions",
		Short: "Get permissions for the current user in a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseSettingsPermissions(cmd.Context(), client, courseID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runCourseSettingsPermissions(ctx context.Context, client *api.Client, courseID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-settings.permissions", map[string]interface{}{"course_id": courseID})

	svc := api.NewCourseSettingsService(client)
	perms, err := svc.GetPermissions(ctx, courseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-settings.permissions", len(perms))
	return formatOutput(perms, nil)
}

func newCourseSettingsEffectiveDueDatesCmd() *cobra.Command {
	var courseID int64
	cmd := &cobra.Command{
		Use:   "effective-due-dates",
		Short: "Get effective due dates for all assignments in a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseSettingsEffectiveDueDates(cmd.Context(), client, courseID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runCourseSettingsEffectiveDueDates(ctx context.Context, client *api.Client, courseID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-settings.effective-due-dates", map[string]interface{}{"course_id": courseID})

	svc := api.NewCourseSettingsService(client)
	dueDates, err := svc.GetEffectiveDueDates(ctx, courseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-settings.effective-due-dates", len(dueDates))
	return formatOutput(dueDates, nil)
}

func newCourseSettingsLatePolicyCmd() *cobra.Command {
	var courseID int64
	cmd := &cobra.Command{
		Use:   "late-policy",
		Short: "Get the late policy for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseSettingsLatePolicy(cmd.Context(), client, courseID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runCourseSettingsLatePolicy(ctx context.Context, client *api.Client, courseID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-settings.late-policy", map[string]interface{}{"course_id": courseID})

	svc := api.NewCourseSettingsService(client)
	policy, err := svc.GetLatePolicy(ctx, courseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-settings.late-policy", 1)
	return formatOutput(policy, nil)
}

func newCourseSettingsRecentStudentsCmd() *cobra.Command {
	var courseID int64
	cmd := &cobra.Command{
		Use:   "recent-students",
		Short: "List recently enrolled students for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseSettingsRecentStudents(cmd.Context(), client, courseID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runCourseSettingsRecentStudents(ctx context.Context, client *api.Client, courseID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-settings.recent-students", map[string]interface{}{"course_id": courseID})

	svc := api.NewCourseSettingsService(client)
	students, err := svc.GetRecentStudents(ctx, courseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-settings.recent-students", len(students))
	return formatEmptyOrOutput(students, "No recent students found")
}
