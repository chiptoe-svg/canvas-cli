package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// auditCmd represents the audit command group.
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "View Canvas audit logs",
	Long: `View Canvas authentication, course, and grade-change audit logs.

Examples:
  canvas audit list --type authentication --account-id 1
  canvas audit list --type authentication --user-id 42
  canvas audit list --type authentication --login-id 10
  canvas audit list --type course --account-id 1
  canvas audit list --type course --course-id 123
  canvas audit list --type grade-change
  canvas audit list --type grade-change --course-id 123
  canvas audit list --type grade-change --assignment-id 456
  canvas audit list --type grade-change --grader-id 7
  canvas audit list --type grade-change --student-id 8`,
}

func init() {
	rootCmd.AddCommand(auditCmd)
	auditCmd.AddCommand(newAuditListCmd())
}

func newAuditListCmd() *cobra.Command {
	opts := &options.AuditLogsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List audit log events",
		Long: `List Canvas audit log events. Use --type to select the audit category,
then combine with a context flag to narrow the results.

Audit types:
  authentication  — Login/logout events (use --account-id, --user-id, or --login-id)
  course          — Course-change events (use --account-id or --course-id)
  grade-change    — Grade-change events (use --course-id, --assignment-id, --grader-id, or --student-id)

Examples:
  canvas audit list --type authentication --account-id 1
  canvas audit list --type grade-change --course-id 123 --start-time 2024-01-01`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAuditList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.AuditType, "type", "grade-change", "Audit type: authentication, course, grade-change")
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (authentication/course audits)")
	cmd.Flags().Int64Var(&opts.LoginID, "login-id", 0, "Login ID (authentication audits)")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (authentication audits)")
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (course/grade-change audits)")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (grade-change audits)")
	cmd.Flags().Int64Var(&opts.GraderID, "grader-id", 0, "Grader user ID (grade-change audits)")
	cmd.Flags().Int64Var(&opts.StudentID, "student-id", 0, "Student user ID (grade-change audits)")
	cmd.Flags().StringVar(&opts.StartTime, "start-time", "", "Filter events from this time (ISO 8601)")
	cmd.Flags().StringVar(&opts.EndTime, "end-time", "", "Filter events up to this time (ISO 8601)")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Number of results per page")

	return cmd
}

func runAuditList(ctx context.Context, client *api.Client, opts *options.AuditLogsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "audit.list", map[string]interface{}{
		"type":          opts.AuditType,
		"account_id":    opts.AccountID,
		"user_id":       opts.UserID,
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
	})

	svc := api.NewAuditLogsService(client)

	apiOpts := &api.AuditLogOptions{
		StartTime: opts.StartTime,
		EndTime:   opts.EndTime,
		PerPage:   opts.PerPage,
	}

	var (
		events []api.AuditLogEvent
		err    error
	)

	switch opts.AuditType {
	case "authentication":
		switch {
		case opts.AccountID > 0:
			events, err = svc.ListAuthenticationForAccount(ctx, opts.AccountID, apiOpts)
		case opts.UserID > 0:
			events, err = svc.ListAuthenticationForUser(ctx, opts.UserID, apiOpts)
		case opts.LoginID > 0:
			events, err = svc.ListAuthenticationForLogin(ctx, opts.LoginID, apiOpts)
		default:
			return fmt.Errorf("for --type authentication, one of --account-id, --user-id, or --login-id is required")
		}

	case "course":
		switch {
		case opts.AccountID > 0:
			events, err = svc.ListCourseEventsForAccount(ctx, opts.AccountID, apiOpts)
		case opts.CourseID > 0:
			events, err = svc.ListCourseEventsForCourse(ctx, opts.CourseID, apiOpts)
		default:
			return fmt.Errorf("for --type course, one of --account-id or --course-id is required")
		}

	case "grade-change":
		switch {
		case opts.CourseID > 0:
			events, err = svc.ListGradeChangeForCourse(ctx, opts.CourseID, apiOpts)
		case opts.AssignmentID > 0:
			events, err = svc.ListGradeChangeForAssignment(ctx, opts.AssignmentID, apiOpts)
		case opts.GraderID > 0:
			events, err = svc.ListGradeChangeForGrader(ctx, opts.GraderID, apiOpts)
		case opts.StudentID > 0:
			events, err = svc.ListGradeChangeForStudent(ctx, opts.StudentID, apiOpts)
		default:
			events, err = svc.ListGradeChangeEvents(ctx, apiOpts)
		}

	default:
		return fmt.Errorf("unknown audit type %q; valid types: authentication, course, grade-change", opts.AuditType)
	}

	if err != nil {
		logger.LogCommandError(ctx, "audit.list", err, nil)
		return fmt.Errorf("failed to list audit events: %w", err)
	}

	printVerbose("Found %d audit events:\n\n", len(events))
	logger.LogCommandComplete(ctx, "audit.list", len(events))
	return formatEmptyOrOutput(events, "No audit events found")
}
