package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

// courseExtrasCmd is the parent command for miscellaneous course extension operations.
var courseExtrasCmd = &cobra.Command{
	Use:     "course-extensions",
	Aliases: []string{"cext"},
	Short:   "Manage course extensions (quiz and assignment)",
	Long: `Manage Canvas course quiz and assignment extensions for students.

Examples:
  canvas course-extensions quiz --course-id 1 --user-id 5 --extra-time 30
  canvas course-extensions assignment --course-id 1 --assignment-id 3 --user-id 5 --extra-attempts 2`,
}

func init() {
	rootCmd.AddCommand(courseExtrasCmd)
	courseExtrasCmd.AddCommand(newQuizExtensionsCmd())
	courseExtrasCmd.AddCommand(newAssignmentExtensionsCmd())
}

func newQuizExtensionsCmd() *cobra.Command {
	var courseID, userID int64
	var extraTime, extraAttempts int
	var manuallyUnlocked bool

	cmd := &cobra.Command{
		Use:   "quiz",
		Short: "Create a quiz extension for a student",
		Long: `Grant a quiz extension to a student in a course.

Examples:
  canvas course-extensions quiz --course-id 1 --user-id 5 --extra-time 30
  canvas course-extensions quiz --course-id 1 --user-id 5 --extra-attempts 1 --manually-unlocked`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			if userID <= 0 {
				return fmt.Errorf("--user-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runQuizExtensions(cmd.Context(), client, courseID, userID, extraTime, extraAttempts, manuallyUnlocked)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&userID, "user-id", 0, "User ID (required)")
	cmd.Flags().IntVar(&extraTime, "extra-time", 0, "Extra time in minutes")
	cmd.Flags().IntVar(&extraAttempts, "extra-attempts", 0, "Extra attempts")
	cmd.Flags().BoolVar(&manuallyUnlocked, "manually-unlocked", false, "Manually unlock submission")
	mustMarkRequired(cmd, "course-id", "user-id")
	return cmd
}

func runQuizExtensions(ctx context.Context, client *api.Client, courseID, userID int64, extraTime, extraAttempts int, manuallyUnlocked bool) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-extensions.quiz", map[string]interface{}{
		"course_id": courseID,
		"user_id":   userID,
	})

	svc := api.NewCourseExtrasService(client)
	if err := svc.CreateQuizExtensions(ctx, courseID, []api.QuizExtensionParams{{
		UserID:           userID,
		ExtraTime:        extraTime,
		ExtraAttempts:    extraAttempts,
		ManuallyUnlocked: manuallyUnlocked,
	}}); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-extensions.quiz", 1)
	printInfo("Quiz extension created for user %d\n", userID)
	return nil
}

func newAssignmentExtensionsCmd() *cobra.Command {
	var courseID, assignmentID, userID int64
	var extraAttempts int

	cmd := &cobra.Command{
		Use:   "assignment",
		Short: "Create an assignment extension for a student",
		Long: `Grant an assignment extension to a student in a course.

Examples:
  canvas course-extensions assignment --course-id 1 --assignment-id 3 --user-id 5 --extra-attempts 2`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			if assignmentID <= 0 {
				return fmt.Errorf("--assignment-id is required")
			}
			if userID <= 0 {
				return fmt.Errorf("--user-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runAssignmentExtensions(cmd.Context(), client, courseID, assignmentID, userID, extraAttempts)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&assignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.Flags().Int64Var(&userID, "user-id", 0, "User ID (required)")
	cmd.Flags().IntVar(&extraAttempts, "extra-attempts", 0, "Extra submission attempts")
	mustMarkRequired(cmd, "course-id", "assignment-id", "user-id")
	return cmd
}

func runAssignmentExtensions(ctx context.Context, client *api.Client, courseID, assignmentID, userID int64, extraAttempts int) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-extensions.assignment", map[string]interface{}{
		"course_id":     courseID,
		"assignment_id": assignmentID,
		"user_id":       userID,
	})

	svc := api.NewCourseExtrasService(client)
	exts, err := svc.CreateAssignmentExtensions(ctx, courseID, assignmentID, []api.AssignmentExtension{{
		UserID:        userID,
		ExtraAttempts: extraAttempts,
	}})
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-extensions.assignment", len(exts))
	printInfo("Assignment extension created for user %d on assignment %d\n", userID, assignmentID)
	return nil
}
