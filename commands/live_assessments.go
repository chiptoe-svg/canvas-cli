package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// liveAssessmentsCmd is the parent command for live assessment operations.
var liveAssessmentsCmd = &cobra.Command{
	Use:     "live-assessments",
	Aliases: []string{"live"},
	Short:   "Manage course live assessments",
	Long: `Manage Canvas course live assessments and their results.

Examples:
  canvas live-assessments list --course-id 1
  canvas live-assessments results --course-id 1 --assessment-id abc`,
}

func init() {
	rootCmd.AddCommand(liveAssessmentsCmd)
	liveAssessmentsCmd.AddCommand(newLiveAssessmentsListCmd())
	liveAssessmentsCmd.AddCommand(newLiveAssessmentsResultsCmd())
}

func newLiveAssessmentsListCmd() *cobra.Command {
	var courseID int64
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List live assessments for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runLiveAssessmentsList(cmd.Context(), client, courseID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runLiveAssessmentsList(ctx context.Context, client *api.Client, courseID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "live-assessments.list", map[string]interface{}{"course_id": courseID})

	svc := api.NewLiveAssessmentsService(client)
	assessments, err := svc.List(ctx, courseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "live-assessments.list", len(assessments))
	return formatEmptyOrOutput(assessments, "No live assessments found")
}

func newLiveAssessmentsResultsCmd() *cobra.Command {
	var courseID int64
	var assessmentID string
	cmd := &cobra.Command{
		Use:   "results",
		Short: "List results for a live assessment",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			if assessmentID == "" {
				return fmt.Errorf("--assessment-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runLiveAssessmentsResults(cmd.Context(), client, courseID, assessmentID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&assessmentID, "assessment-id", "", "Assessment ID (required)")
	mustMarkRequired(cmd, "course-id", "assessment-id")
	return cmd
}

func runLiveAssessmentsResults(ctx context.Context, client *api.Client, courseID int64, assessmentID string) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "live-assessments.results", map[string]interface{}{
		"course_id":     courseID,
		"assessment_id": assessmentID,
	})

	svc := api.NewLiveAssessmentsService(client)
	results, err := svc.ListResults(ctx, courseID, assessmentID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "live-assessments.results", len(results))
	return formatEmptyOrOutput(results, "No results found")
}
