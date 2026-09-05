package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/commands/internal/options"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

// collaborationsCmd represents the collaborations command group.
var collaborationsCmd = &cobra.Command{
	Use:   "collaborations",
	Short: "Manage Canvas collaborations",
	Long: `List collaborations (Google Docs, etc.) for courses and groups, and
view collaboration members.

Examples:
  canvas collaborations list --course-id 123
  canvas collaborations list --group-id 456
  canvas collaborations members 789`,
}

func init() {
	rootCmd.AddCommand(collaborationsCmd)
	collaborationsCmd.AddCommand(newCollaborationsListCmd())
	collaborationsCmd.AddCommand(newCollaborationsMembersCmd())
}

func newCollaborationsListCmd() *cobra.Command {
	opts := &options.CollaborationsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List collaborations",
		Long: `List collaborations for a course or a group.

Examples:
  canvas collaborations list --course-id 123
  canvas collaborations list --group-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCollaborationsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Group ID")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Number of results per page")

	return cmd
}

func newCollaborationsMembersCmd() *cobra.Command {
	opts := &options.CollaborationsMembersOptions{}

	cmd := &cobra.Command{
		Use:   "members <collaboration-id>",
		Short: "List collaboration members",
		Long: `List the members (users/groups) of a specific collaboration.

Examples:
  canvas collaborations members 789`,
		Args: ExactArgsWithUsage(1, "collaboration-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid collaboration ID: %s", args[0])
			}
			opts.CollaborationID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCollaborationsMembers(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Number of results per page")

	return cmd
}

func runCollaborationsList(ctx context.Context, client *api.Client, opts *options.CollaborationsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "collaborations.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
	})

	svc := api.NewCollaborationsService(client)

	apiOpts := &api.ListCollaborationsOptions{PerPage: opts.PerPage}

	var (
		colls []api.Collaboration
		err   error
	)

	switch {
	case opts.CourseID > 0:
		colls, err = svc.ListForCourse(ctx, opts.CourseID, apiOpts)
	case opts.GroupID > 0:
		colls, err = svc.ListForGroup(ctx, opts.GroupID, apiOpts)
	default:
		return fmt.Errorf("either --course-id or --group-id is required")
	}

	if err != nil {
		logger.LogCommandError(ctx, "collaborations.list", err, nil)
		return fmt.Errorf("failed to list collaborations: %w", err)
	}

	printVerbose("Found %d collaborations:\n\n", len(colls))
	logger.LogCommandComplete(ctx, "collaborations.list", len(colls))
	return formatEmptyOrOutput(colls, "No collaborations found")
}

func runCollaborationsMembers(ctx context.Context, client *api.Client, opts *options.CollaborationsMembersOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "collaborations.members", map[string]interface{}{
		"collaboration_id": opts.CollaborationID,
	})

	svc := api.NewCollaborationsService(client)

	apiOpts := &api.ListCollaborationsOptions{PerPage: opts.PerPage}

	members, err := svc.ListMembers(ctx, opts.CollaborationID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "collaborations.members", err, nil)
		return fmt.Errorf("failed to list collaboration members: %w", err)
	}

	printVerbose("Found %d members:\n\n", len(members))
	logger.LogCommandComplete(ctx, "collaborations.members", len(members))
	return formatEmptyOrOutput(members, "No members found")
}
