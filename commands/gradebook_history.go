package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

func init() {
	// Extend the existing gradesCmd group with date-scoped history subcommands.
	gradesCmd.AddCommand(newGradesHistoryDayCmd())
	gradesCmd.AddCommand(newGradesHistorySubmissionsCmd())
}

func newGradesHistoryDayCmd() *cobra.Command {
	var courseID int64
	var date string

	cmd := &cobra.Command{
		Use:   "history-day",
		Short: "List graders active on a specific date in gradebook history",
		Long: `Returns a list of graders who were active on the given date for the specified course.

The date must be in YYYY-MM-DD format.

Examples:
  canvas grades history-day --course-id 1 --date 2024-03-15`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			if date == "" {
				return fmt.Errorf("--date is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGradesHistoryDay(cmd.Context(), client, courseID, date)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&date, "date", "", "Date in YYYY-MM-DD format (required)")
	mustMarkRequired(cmd, "course-id", "date")
	return cmd
}

func runGradesHistoryDay(ctx context.Context, client *api.Client, courseID int64, date string) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grades.history-day", map[string]interface{}{
		"course_id": courseID,
		"date":      date,
	})

	svc := api.NewGradesService(client)
	graders, err := svc.GetGradebookHistoryDay(ctx, courseID, date)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "grades.history-day", len(graders))
	return formatEmptyOrOutput(graders, fmt.Sprintf("No graders found for date %s", date))
}

func newGradesHistorySubmissionsCmd() *cobra.Command {
	var courseID, graderID, assignmentID int64
	var date string

	cmd := &cobra.Command{
		Use:   "history-submissions",
		Short: "List submissions graded on a date by a grader for an assignment",
		Long: `Returns submissions that were graded on a specific date by a specific grader for an assignment.

Examples:
  canvas grades history-submissions --course-id 1 --date 2024-03-15 --grader-id 5 --assignment-id 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			if date == "" {
				return fmt.Errorf("--date is required")
			}
			if graderID <= 0 {
				return fmt.Errorf("--grader-id is required")
			}
			if assignmentID <= 0 {
				return fmt.Errorf("--assignment-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runGradesHistorySubmissions(cmd.Context(), client, courseID, date, graderID, assignmentID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&date, "date", "", "Date in YYYY-MM-DD format (required)")
	cmd.Flags().Int64Var(&graderID, "grader-id", 0, "Grader user ID (required)")
	cmd.Flags().Int64Var(&assignmentID, "assignment-id", 0, "Assignment ID (required)")
	mustMarkRequired(cmd, "course-id", "date", "grader-id", "assignment-id")
	return cmd
}

func runGradesHistorySubmissions(ctx context.Context, client *api.Client, courseID int64, date string, graderID, assignmentID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grades.history-submissions", map[string]interface{}{
		"course_id":     courseID,
		"date":          date,
		"grader_id":     graderID,
		"assignment_id": assignmentID,
	})

	svc := api.NewGradesService(client)
	submissions, err := svc.GetGradebookHistorySubmissions(ctx, courseID, date, graderID, assignmentID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "grades.history-submissions", len(submissions))
	return formatEmptyOrOutput(submissions, "No submissions found")
}
