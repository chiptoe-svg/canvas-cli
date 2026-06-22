package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// conferencesCmd represents the conferences command group.
var conferencesCmd = &cobra.Command{
	Use:   "conferences",
	Short: "List Canvas web conferences",
	Long: `List web conferences (BigBlueButton, Zoom, etc.) for the current user,
a course, or a group.

Examples:
  canvas conferences list
  canvas conferences list --course-id 123
  canvas conferences list --group-id 456`,
}

func init() {
	rootCmd.AddCommand(conferencesCmd)
	conferencesCmd.AddCommand(newConferencesListCmd())
}

func newConferencesListCmd() *cobra.Command {
	opts := &options.ConferencesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List web conferences",
		Long: `List web conferences for the current user, a course, or a group.

When neither --course-id nor --group-id is given, the conferences for the
current user across all contexts are returned.

Examples:
  canvas conferences list
  canvas conferences list --course-id 123
  canvas conferences list --group-id 456
  canvas conferences list --state live`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runConferencesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (list conferences for a course)")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Group ID (list conferences for a group)")
	cmd.Flags().StringVar(&opts.State, "state", "", "Filter by state: live, ended")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Number of results per page")

	return cmd
}

func runConferencesList(ctx context.Context, client *api.Client, opts *options.ConferencesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "conferences.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
		"state":     opts.State,
	})

	svc := api.NewConferencesService(client)

	apiOpts := &api.ListConferencesOptions{
		State:   opts.State,
		PerPage: opts.PerPage,
	}

	var (
		confs []api.Conference
		err   error
	)

	switch {
	case opts.CourseID > 0:
		confs, err = svc.ListForCourse(ctx, opts.CourseID, apiOpts)
	case opts.GroupID > 0:
		confs, err = svc.ListForGroup(ctx, opts.GroupID, apiOpts)
	default:
		confs, err = svc.List(ctx, apiOpts)
	}

	if err != nil {
		logger.LogCommandError(ctx, "conferences.list", err, nil)
		return fmt.Errorf("failed to list conferences: %w", err)
	}

	printVerbose("Found %d conferences:\n\n", len(confs))
	logger.LogCommandComplete(ctx, "conferences.list", len(confs))
	return formatEmptyOrOutput(confs, "No conferences found")
}
