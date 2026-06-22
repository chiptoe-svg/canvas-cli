package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// courseFeaturesCmd is the parent command for course feature flag operations.
var courseFeaturesCmd = &cobra.Command{
	Use:     "course-features",
	Aliases: []string{"cfeatures"},
	Short:   "Manage course feature flags",
	Long: `Manage Canvas course feature flags.

Examples:
  canvas course-features list --course-id 1
  canvas course-features enabled --course-id 1
  canvas course-features get-flag --course-id 1 --feature new_quizzes
  canvas course-features set-flag --course-id 1 --feature new_quizzes --state on
  canvas course-features delete-flag --course-id 1 --feature new_quizzes`,
}

func init() {
	rootCmd.AddCommand(courseFeaturesCmd)
	courseFeaturesCmd.AddCommand(newCourseFeaturesListCmd())
	courseFeaturesCmd.AddCommand(newCourseFeaturesListEnabledCmd())
	courseFeaturesCmd.AddCommand(newCourseFeaturesGetFlagCmd())
	courseFeaturesCmd.AddCommand(newCourseFeaturesSetFlagCmd())
	courseFeaturesCmd.AddCommand(newCourseFeaturesDeleteFlagCmd())
}

func newCourseFeaturesListCmd() *cobra.Command {
	var courseID int64
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all features for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseFeaturesListAll(cmd.Context(), client, courseID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runCourseFeaturesListAll(ctx context.Context, client *api.Client, courseID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-features.list", map[string]interface{}{"course_id": courseID})

	svc := api.NewCourseFeaturesService(client)
	features, err := svc.List(ctx, courseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-features.list", len(features))
	return formatEmptyOrOutput(features, "No features found")
}

func newCourseFeaturesListEnabledCmd() *cobra.Command {
	var courseID int64
	cmd := &cobra.Command{
		Use:   "enabled",
		Short: "List enabled features for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseFeaturesListEnabled(cmd.Context(), client, courseID)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")
	return cmd
}

func runCourseFeaturesListEnabled(ctx context.Context, client *api.Client, courseID int64) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-features.enabled", map[string]interface{}{"course_id": courseID})

	svc := api.NewCourseFeaturesService(client)
	features, err := svc.ListEnabled(ctx, courseID)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-features.enabled", len(features))
	return formatEmptyOrOutput(features, "No enabled features found")
}

func newCourseFeaturesGetFlagCmd() *cobra.Command {
	var courseID int64
	var feature string
	cmd := &cobra.Command{
		Use:   "get-flag",
		Short: "Get a feature flag for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			if feature == "" {
				return fmt.Errorf("--feature is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseFeaturesGetFlag(cmd.Context(), client, courseID, feature)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&feature, "feature", "", "Feature name (required)")
	mustMarkRequired(cmd, "course-id", "feature")
	return cmd
}

func runCourseFeaturesGetFlag(ctx context.Context, client *api.Client, courseID int64, feature string) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-features.get-flag", map[string]interface{}{"course_id": courseID, "feature": feature})

	svc := api.NewCourseFeaturesService(client)
	flag, err := svc.GetFlag(ctx, courseID, feature)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-features.get-flag", 1)
	return formatOutput(flag, nil)
}

func newCourseFeaturesSetFlagCmd() *cobra.Command {
	var courseID int64
	var feature, state string
	cmd := &cobra.Command{
		Use:   "set-flag",
		Short: "Set a feature flag state for a course",
		Long: `Set a feature flag for a course.

Valid states: off, allowed, on

Examples:
  canvas course-features set-flag --course-id 1 --feature new_quizzes --state on`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			if feature == "" {
				return fmt.Errorf("--feature is required")
			}
			if state == "" {
				return fmt.Errorf("--state is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseFeaturesSetFlag(cmd.Context(), client, courseID, feature, state)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&feature, "feature", "", "Feature name (required)")
	cmd.Flags().StringVar(&state, "state", "", "State: off, allowed, on (required)")
	mustMarkRequired(cmd, "course-id", "feature", "state")
	return cmd
}

func runCourseFeaturesSetFlag(ctx context.Context, client *api.Client, courseID int64, feature, state string) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-features.set-flag", map[string]interface{}{"course_id": courseID, "feature": feature, "state": state})

	svc := api.NewCourseFeaturesService(client)
	flag, err := svc.SetFlag(ctx, courseID, feature, state)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-features.set-flag", 1)
	return formatSuccessOutput(flag, fmt.Sprintf("Feature %q set to %q", feature, state))
}

func newCourseFeaturesDeleteFlagCmd() *cobra.Command {
	var courseID int64
	var feature string
	cmd := &cobra.Command{
		Use:   "delete-flag",
		Short: "Remove a feature flag override for a course",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID <= 0 {
				return fmt.Errorf("--course-id is required")
			}
			if feature == "" {
				return fmt.Errorf("--feature is required")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runCourseFeaturesDeleteFlag(cmd.Context(), client, courseID, feature)
		},
	}
	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&feature, "feature", "", "Feature name (required)")
	mustMarkRequired(cmd, "course-id", "feature")
	return cmd
}

func runCourseFeaturesDeleteFlag(ctx context.Context, client *api.Client, courseID int64, feature string) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "course-features.delete-flag", map[string]interface{}{"course_id": courseID, "feature": feature})

	svc := api.NewCourseFeaturesService(client)
	flag, err := svc.DeleteFlag(ctx, courseID, feature)
	if err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "course-features.delete-flag", 1)
	return formatSuccessOutput(flag, fmt.Sprintf("Feature flag %q removed (reverted to state: %q)", feature, flag.State))
}
