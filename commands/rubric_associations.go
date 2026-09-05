package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

// rubricAssociationsCmd is the parent command for rubric association operations.
var rubricAssociationsCmd = &cobra.Command{
	Use:     "rubric-associations",
	Aliases: []string{"rassoc"},
	Short:   "Manage rubric associations and assessments",
	Long: `Manage Canvas rubric associations and rubric assessments in a course.

Examples:
  canvas rubric-associations update 5 --course-id 1
  canvas rubric-associations delete 5 --course-id 1
  canvas rubric-associations assess 5 --course-id 1
  canvas rubric-associations delete-assessment 5 --course-id 1 --assessment-id 3`,
}

func init() {
	rootCmd.AddCommand(rubricAssociationsCmd)
	rubricAssociationsCmd.AddCommand(newRubricAssociationsUpdateCmd())
	rubricAssociationsCmd.AddCommand(newRubricAssociationsDeleteCmd())
	rubricAssociationsCmd.AddCommand(newRubricAssessmentsCreateCmd())
	rubricAssociationsCmd.AddCommand(newRubricAssessmentsDeleteCmd())
}

func newRubricAssociationsUpdateCmd() *cobra.Command {
	var courseID int64
	var rubricID, purposeTypeID int64
	var purpose string

	cmd := &cobra.Command{
		Use:   "update <association-id>",
		Short: "Update a rubric association",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var assocID int64
			if _, err := fmt.Sscanf(args[0], "%d", &assocID); err != nil {
				return fmt.Errorf("invalid association ID: %s", args[0])
			}
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runRubricAssociationsUpdate(cmd.Context(), client, courseID, assocID, rubricID, purposeTypeID, purpose)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&rubricID, "rubric-id", 0, "Rubric ID")
	cmd.Flags().Int64Var(&purposeTypeID, "purpose-type-id", 0, "Purpose type ID")
	cmd.Flags().StringVar(&purpose, "purpose", "", "Purpose (grading, bookmark)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runRubricAssociationsUpdate(ctx context.Context, client *api.Client, courseID, assocID, rubricID, purposeTypeID int64, purpose string) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "rubric-associations.update", map[string]interface{}{
		"course_id": courseID,
		"assoc_id":  assocID,
	})

	svc := api.NewRubricAssociationsService(client)
	assoc, err := svc.Update(ctx, courseID, assocID, api.RubricAssociationUpdateParams{
		RubricID:      rubricID,
		AssociationID: purposeTypeID,
		Purpose:       purpose,
	})
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "rubric-associations.update", 1)
	return formatSuccessOutput(assoc, fmt.Sprintf("Rubric association %d updated", assoc.ID))
}

func newRubricAssociationsDeleteCmd() *cobra.Command {
	var courseID int64
	cmd := &cobra.Command{
		Use:   "delete <association-id>",
		Short: "Delete a rubric association",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var assocID int64
			if _, err := fmt.Sscanf(args[0], "%d", &assocID); err != nil {
				return fmt.Errorf("invalid association ID: %s", args[0])
			}
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runRubricAssociationsDelete(cmd.Context(), client, courseID, assocID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runRubricAssociationsDelete(ctx context.Context, client *api.Client, courseID, assocID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "rubric-associations.delete", map[string]interface{}{
		"course_id": courseID,
		"assoc_id":  assocID,
	})

	svc := api.NewRubricAssociationsService(client)
	assoc, err := svc.Delete(ctx, courseID, assocID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "rubric-associations.delete", 1)
	printInfo("Rubric association %d deleted\n", assoc.ID)
	return nil
}

func newRubricAssessmentsCreateCmd() *cobra.Command {
	var courseID, assocID, artifactID int64
	var assessmentType string

	cmd := &cobra.Command{
		Use:   "assess <association-id>",
		Short: "Create a rubric assessment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &assocID); err != nil {
				return fmt.Errorf("invalid association ID: %s", args[0])
			}
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runRubricAssessmentsCreate(cmd.Context(), client, courseID, assocID, artifactID, assessmentType)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&artifactID, "artifact-id", 0, "Artifact ID being assessed")
	cmd.Flags().StringVar(&assessmentType, "assessment-type", "grading", "Assessment type: grading, peer_review, provisional_grade")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runRubricAssessmentsCreate(ctx context.Context, client *api.Client, courseID, assocID, artifactID int64, assessmentType string) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "rubric-assessments.create", map[string]interface{}{
		"course_id": courseID,
		"assoc_id":  assocID,
	})

	svc := api.NewRubricAssociationsService(client)
	record, err := svc.CreateAssessment(ctx, courseID, assocID, api.RubricAssessmentRecord{
		ArtifactID:   artifactID,
		ArtifactType: assessmentType,
	})
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "rubric-assessments.create", 1)
	return formatSuccessOutput(record, fmt.Sprintf("Rubric assessment %d created", record.ID))
}

func newRubricAssessmentsDeleteCmd() *cobra.Command {
	var courseID, assocID, assessmentID int64
	cmd := &cobra.Command{
		Use:   "delete-assessment <association-id>",
		Short: "Delete a rubric assessment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Sscanf(args[0], "%d", &assocID); err != nil {
				return fmt.Errorf("invalid association ID: %s", args[0])
			}
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			if assessmentID <= 0 {
				return fmt.Errorf("--assessment-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runRubricAssessmentsDelete(cmd.Context(), client, courseID, assocID, assessmentID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&assessmentID, "assessment-id", 0, "Assessment ID (required)")
	mustMarkRequired(cmd, "course-id", "assessment-id")
	return cmd
}

func runRubricAssessmentsDelete(ctx context.Context, client *api.Client, courseID, assocID, assessmentID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "rubric-assessments.delete", map[string]interface{}{
		"course_id":     courseID,
		"assoc_id":      assocID,
		"assessment_id": assessmentID,
	})

	svc := api.NewRubricAssociationsService(client)
	if _, err := svc.DeleteAssessment(ctx, courseID, assocID, assessmentID); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "rubric-assessments.delete", 1)
	printInfo("Rubric assessment %d deleted\n", assessmentID)
	return nil
}
