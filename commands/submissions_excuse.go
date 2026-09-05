package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/commands/internal/options"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
	"github.com/chiptoe-svg/canvas-cli/internal/output"
	"github.com/chiptoe-svg/canvas-cli/internal/resolve"
)

func init() {
	submissionsCmd.AddCommand(newSubmissionsExcuseCmd())
}

func newSubmissionsExcuseCmd() *cobra.Command {
	opts := &options.SubmissionsExcuseOptions{}

	cmd := &cobra.Command{
		Use:   "excuse",
		Short: "Excuse a student from an assignment or quiz by name, with read-back",
		Long: `Excuse one student from one assignment or quiz, naming both the way
people do. The student is resolved among the course's active students by
id, exact name, sortable name ("Lovelace, Ada"), login or SIS id, or a
substring that matches exactly one student. The assignment is resolved
among the course's assignments by assignment id, quiz id, exact name, or a
substring that matches exactly one; a quiz resolves to the assignment
Canvas grades it under. Zero or several matches are refused and the
candidates listed.

The command reads the submission, writes submission[excuse], reads the
submission back and prints excused before → after with a verified line.
It exits non-zero if the read-back does not show the requested state.
--unexcuse clears an excusal. --dry-run prints what was resolved and the
request that would be sent, and writes nothing.

Examples:
  canvas submissions excuse --course-id 123 --student "Ada Lovelace" --assignment "Quiz 3"
  canvas submissions excuse --course-id 123 --student 789 --assignment 456 --force
  canvas submissions excuse --course-id 123 --student "lovelace" --assignment "lineup" --dry-run
  canvas submissions excuse --course-id 123 --student "Ada Lovelace" --assignment "Quiz 3" --unexcuse`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.DryRun = dryRun
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			// Resolution needs real reads even under --dry-run; the command
			// prints its own preview and never writes in that mode.
			client.SetDryRun(false)
			// The before-read and the read-back must be live.
			client.SetCacheEnabled(false)
			return runSubmissionsExcuse(cmd.Context(), client, opts, cmd.OutOrStdout(), cmd.InOrStdin())
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Student, "student", "", "Student: id, name, sortable name, login/SIS id, or a unique part of the name (required)")
	cmd.Flags().StringVar(&opts.Assignment, "assignment", "", "Assignment: assignment id, quiz id, name, or a unique part of the name (required)")
	cmd.Flags().BoolVar(&opts.Unexcuse, "unexcuse", false, "Clear the excusal instead of setting it")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip the confirmation prompt")
	mustMarkRequired(cmd, "course-id", "student", "assignment")

	return cmd
}

// excuseResult is the evidence of one excusal write.
type excuseResult struct {
	CourseID   int64               `json:"course_id" yaml:"course_id"`
	Student    *resolve.Student    `json:"student" yaml:"student"`
	Assignment *resolve.Assignment `json:"assignment" yaml:"assignment"`
	Requested  bool                `json:"requested_excused" yaml:"requested_excused"`
	Before     *bool               `json:"before_excused" yaml:"before_excused"` // nil when the pre-read failed
	After      *bool               `json:"after_excused" yaml:"after_excused"`   // nil when nothing was written
	Written    bool                `json:"written" yaml:"written"`
	Verified   bool                `json:"verified" yaml:"verified"`
	Mismatches []string            `json:"mismatches,omitempty" yaml:"mismatches,omitempty"`
	DryRun     bool                `json:"dry_run" yaml:"dry_run"`
	Skipped    string              `json:"skipped,omitempty" yaml:"skipped,omitempty"`
}

func (r *excuseResult) path() string {
	return fmt.Sprintf("/api/v1/courses/%d/assignments/%d/submissions/%d", r.CourseID, r.Assignment.ID, r.Student.ID)
}

func runSubmissionsExcuse(ctx context.Context, client *api.Client, opts *options.SubmissionsExcuseOptions, out io.Writer, in io.Reader) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "submissions.excuse", map[string]interface{}{
		"course_id":  opts.CourseID,
		"student":    opts.Student,
		"assignment": opts.Assignment,
		"excused":    opts.Excused(),
		"dry_run":    opts.DryRun,
	})

	students, err := resolve.ListStudents(ctx, client, opts.CourseID)
	if err != nil {
		return err
	}
	student, err := resolve.FindStudent(opts.Student, students)
	if err != nil {
		return fmt.Errorf("--student: %w", err)
	}
	items, err := resolve.ListAssignments(ctx, client, opts.CourseID)
	if err != nil {
		return err
	}
	assignment, err := resolve.FindAssignment(opts.Assignment, items)
	if err != nil {
		return fmt.Errorf("--assignment: %w", err)
	}

	result := &excuseResult{CourseID: opts.CourseID, Student: student, Assignment: assignment, Requested: opts.Excused(), DryRun: opts.DryRun}
	structured := isStructuredOutput()
	if !structured {
		fmt.Fprintf(out, "student:    %s\n", resolve.DescribeStudent(student))
		fmt.Fprintf(out, "assignment: %s\n", resolve.DescribeAssignment(assignment))
	}

	submissions := api.NewSubmissionsService(client)
	// Canvas API: GET /api/v1/courses/:course_id/assignments/:assignment_id/submissions/:user_id
	// https://canvas.instructure.com/doc/api/submissions.html#method.submissions_api.show
	if before, err := submissions.Get(ctx, opts.CourseID, assignment.ID, student.ID, nil); err == nil {
		v := before.ExcusedTLN
		result.Before = &v
	} else if !structured {
		fmt.Fprintf(out, "note: could not read the submission before writing: %v\n", err)
	}

	if result.Before != nil && *result.Before == result.Requested {
		result.Skipped = fmt.Sprintf("already %s", excusedWord(result.Requested))
		result.Verified = true
		if err := printExcuseResult(out, result); err != nil {
			return err
		}
		return nil
	}

	if opts.DryRun {
		return printExcuseResult(out, result)
	}

	if !opts.Force {
		verb := "Excuse"
		if !result.Requested {
			verb = "Un-excuse"
		}
		fmt.Fprintf(out, "%s %s from %q (assignment %d)? [y/N]: ", verb, student.Name, assignment.Name, assignment.ID)
		response, err := bufio.NewReader(in).ReadString('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read response: %w", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Fprintln(out, "Cancelled.")
			return nil
		}
	}

	if _, err := submissions.SetExcused(ctx, opts.CourseID, assignment.ID, student.ID, result.Requested); err != nil {
		logger.LogCommandError(ctx, "submissions.excuse", err, nil)
		return fmt.Errorf("failed to update the submission: %w", err)
	}
	result.Written = true

	after, err := submissions.Get(ctx, opts.CourseID, assignment.ID, student.ID, nil)
	if err != nil {
		return fmt.Errorf("written, but the read-back failed: %w", err)
	}
	v := after.ExcusedTLN
	result.After = &v
	if v != result.Requested {
		result.Mismatches = append(result.Mismatches, fmt.Sprintf("excused read back %s, requested %s", excusedWord(v), excusedWord(result.Requested)))
	}
	result.Verified = len(result.Mismatches) == 0

	logger.LogCommandComplete(ctx, "submissions.excuse", 1)
	if err := printExcuseResult(out, result); err != nil {
		return err
	}
	if !result.Verified {
		return fmt.Errorf("submission of %s for assignment %d did not read back as requested: %s", student.Name, assignment.ID, strings.Join(result.Mismatches, "; "))
	}
	return nil
}

func excusedWord(b bool) string {
	if b {
		return "excused"
	}
	return "not excused"
}

func excusedState(p *bool) string {
	if p == nil {
		return "?"
	}
	return excusedWord(*p)
}

func printExcuseResult(out io.Writer, r *excuseResult) error {
	format := output.FormatType(outputFormat)
	switch format {
	case output.FormatJSON, output.FormatYAML:
		return output.WriteWithOptions(out, r, format, verbose)
	case output.FormatCSV:
		return fmt.Errorf("unsupported output format %q for submissions excuse (table, json, yaml)", outputFormat)
	}

	switch {
	case r.Skipped != "":
		fmt.Fprintf(out, "excused:    %s (%s; nothing to do)\n", excusedState(r.Before), r.Skipped)
	case r.DryRun:
		fmt.Fprintf(out, "excused:    %s → %s (planned)\n", excusedState(r.Before), excusedWord(r.Requested))
		fmt.Fprintf(out, "request:    PUT %s {\"submission\":{\"excuse\":%t}}\n", r.path(), r.Requested)
		fmt.Fprintln(out, "DRY RUN: no changes were made.")
	default:
		fmt.Fprintf(out, "excused:    %s → %s\n", excusedState(r.Before), excusedState(r.After))
		if r.Verified {
			fmt.Fprintln(out, "verified:   yes")
		} else {
			fmt.Fprintf(out, "verified:   no — %s\n", strings.Join(r.Mismatches, "; "))
		}
	}
	return nil
}
