package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// brandCmd represents the brand command group.
var brandCmd = &cobra.Command{
	Use:   "brand",
	Short: "View Canvas brand/theme variables",
	Long: `Retrieve Canvas branding and theme variables (colors, logos, etc.).

Examples:
  canvas brand get
  canvas brand get --account-id 1
  canvas brand get --course-id 123`,
}

func init() {
	rootCmd.AddCommand(brandCmd)
	brandCmd.AddCommand(newBrandGetCmd())
}

func newBrandGetCmd() *cobra.Command {
	opts := &options.BrandGetOptions{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get brand variables",
		Long: `Retrieve Canvas brand/theme variables such as primary colors and logos.

When no flags are provided, the global brand variables are returned.

Examples:
  canvas brand get
  canvas brand get --account-id 1
  canvas brand get --course-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runBrandGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (get account-specific brand variables)")
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (get course-specific brand variables)")

	return cmd
}

func runBrandGet(ctx context.Context, client *api.Client, opts *options.BrandGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "brand.get", map[string]interface{}{
		"account_id": opts.AccountID,
		"course_id":  opts.CourseID,
	})

	svc := api.NewBrandService(client)

	var (
		bv  *api.BrandVariables
		err error
	)

	switch {
	case opts.AccountID > 0:
		bv, err = svc.GetVariablesForAccount(ctx, opts.AccountID)
	case opts.CourseID > 0:
		bv, err = svc.GetVariablesForCourse(ctx, opts.CourseID)
	default:
		bv, err = svc.GetVariables(ctx)
	}

	if err != nil {
		logger.LogCommandError(ctx, "brand.get", err, nil)
		return fmt.Errorf("failed to get brand variables: %w", err)
	}

	logger.LogCommandComplete(ctx, "brand.get", 1)
	return formatOutput(bv, nil)
}
