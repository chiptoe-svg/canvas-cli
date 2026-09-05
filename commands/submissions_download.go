package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/commands/internal/options"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

const submissionDownloadManifestName = "submission-download-manifest.json"

type submissionDownloadManifest struct {
	CourseID     int64                      `json:"course_id"`
	AssignmentID int64                      `json:"assignment_id"`
	GeneratedAt  time.Time                  `json:"generated_at"`
	Entries      []submissionDownloadRecord `json:"entries"`
}

type submissionDownloadRecord struct {
	SubmissionID   int64  `json:"submission_id"`
	UserID         int64  `json:"user_id"`
	WorkflowState  string `json:"workflow_state"`
	SubmissionType string `json:"submission_type"`
	Attempt        int    `json:"attempt,omitempty"`
	AttachmentID   int64  `json:"attachment_id,omitempty"`
	Filename       string `json:"filename,omitempty"`
	Path           string `json:"path,omitempty"`
	Status         string `json:"status"`
	Error          string `json:"error,omitempty"`
}

func newSubmissionsDownloadCmd() *cobra.Command {
	opts := &options.SubmissionsDownloadOptions{}

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download all files attached to an assignment's submissions",
		Long: `Download every file attached to submissions for one assignment.

Files are organized by Canvas user ID and attachment ID, preventing student
filenames from overwriting one another. Existing files are skipped unless
--overwrite is supplied. A submission-download-manifest.json file records every
submission, including text- and URL-only submissions that have no file to download.

Examples:
  canvas submissions download --course-id 123 --assignment-id 456 --destination ./essay-submissions
  canvas submissions download --course-id 123 --assignment-id 456 --destination ./essay-submissions --overwrite`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runSubmissionsDownload(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.Flags().StringVar(&opts.Destination, "destination", "", "Directory for downloaded submissions (required)")
	cmd.Flags().BoolVar(&opts.Overwrite, "overwrite", false, "Replace files that already exist")
	mustMarkRequired(cmd, "course-id", "assignment-id", "destination")

	return cmd
}

func runSubmissionsDownload(ctx context.Context, client *api.Client, opts *options.SubmissionsDownloadOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "submissions.download", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
		"destination":   opts.Destination,
	})

	if _, err := validateCourseID(ctx, client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "submissions.download", err, nil)
		return err
	}

	destination, err := filepath.Abs(opts.Destination)
	if err != nil {
		return fmt.Errorf("resolve destination: %w", err)
	}
	if err := os.MkdirAll(destination, 0o700); err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	if err := os.Chmod(destination, 0o700); err != nil {
		return fmt.Errorf("protect destination: %w", err)
	}

	submissionsService := api.NewSubmissionsService(client)
	// include[]=submission_history is load-bearing, not extra detail. A
	// submission object carries only the CURRENT attempt's attachments, and
	// students routinely upload the parts of a multi-part assignment across
	// separate attempts — so listing without the history silently reports
	// earlier files as "no_attachment" and hands the grader an incomplete set.
	submissions, err := submissionsService.List(ctx, opts.CourseID, opts.AssignmentID,
		&api.ListSubmissionsOptions{Include: []string{"submission_history"}})
	if err != nil {
		logger.LogCommandError(ctx, "submissions.download", err, nil)
		return fmt.Errorf("list submissions: %w", err)
	}

	manifest := submissionDownloadManifest{
		CourseID:     opts.CourseID,
		AssignmentID: opts.AssignmentID,
		GeneratedAt:  time.Now().UTC(),
	}
	filesService := api.NewFilesService(client)
	downloaded, skipped, failed := 0, 0, 0

	for _, submission := range submissions {
		attachments := submissionAttachmentsAcrossAttempts(submission)
		if len(attachments) == 0 {
			manifest.Entries = append(manifest.Entries, submissionDownloadRecord{
				SubmissionID: submission.ID, UserID: submission.UserID,
				WorkflowState: submission.WorkflowState, SubmissionType: submission.SubmissionType,
				Status: "no_attachment",
			})
			continue
		}

		for _, item := range attachments {
			attachment := item.attachment
			filename := safeSubmissionAttachmentFilename(attachment)
			userDir := filepath.Join(destination, fmt.Sprintf("user-%d", submission.UserID))
			localPath := filepath.Join(userDir, fmt.Sprintf("attachment-%d-%s", attachment.ID, filename))
			record := submissionDownloadRecord{
				SubmissionID: submission.ID, UserID: submission.UserID,
				WorkflowState: submission.WorkflowState, SubmissionType: submission.SubmissionType,
				Attempt: item.attempt, AttachmentID: attachment.ID, Filename: filename, Path: localPath,
			}

			if !opts.Overwrite {
				if _, err := os.Stat(localPath); err == nil {
					record.Status = "skipped_existing"
					manifest.Entries = append(manifest.Entries, record)
					skipped++
					continue
				} else if !os.IsNotExist(err) {
					record.Status, record.Error = "failed", err.Error()
					manifest.Entries = append(manifest.Entries, record)
					failed++
					continue
				}
			}

			if err := os.MkdirAll(userDir, 0o700); err != nil {
				record.Status, record.Error = "failed", err.Error()
				manifest.Entries = append(manifest.Entries, record)
				failed++
				continue
			}
			if err := os.Chmod(userDir, 0o700); err != nil {
				record.Status, record.Error = "failed", err.Error()
				manifest.Entries = append(manifest.Entries, record)
				failed++
				continue
			}

			partialPath := localPath + ".partial"
			_ = os.Remove(partialPath)
			if err := filesService.Download(ctx, attachment.ID, partialPath); err != nil {
				_ = os.Remove(partialPath)
				record.Status, record.Error = "failed", err.Error()
				manifest.Entries = append(manifest.Entries, record)
				failed++
				continue
			}
			if err := os.Rename(partialPath, localPath); err != nil {
				_ = os.Remove(partialPath)
				record.Status, record.Error = "failed", err.Error()
				manifest.Entries = append(manifest.Entries, record)
				failed++
				continue
			}
			if err := os.Chmod(localPath, 0o600); err != nil {
				record.Status, record.Error = "failed", err.Error()
				manifest.Entries = append(manifest.Entries, record)
				failed++
				continue
			}

			record.Status = "downloaded"
			manifest.Entries = append(manifest.Entries, record)
			downloaded++
		}
	}

	manifestPath := filepath.Join(destination, submissionDownloadManifestName)
	if err := writeSubmissionDownloadManifest(manifestPath, manifest); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}

	fmt.Printf("Downloaded %d file(s), skipped %d, failed %d. Manifest: %s\n", downloaded, skipped, failed, manifestPath)
	logger.LogCommandComplete(ctx, "submissions.download", downloaded)
	if failed > 0 {
		return fmt.Errorf("%d file download(s) failed; see %s", failed, manifestPath)
	}
	return nil
}

// attemptAttachment pairs a file with the attempt that carried it, so the
// manifest can say which submission attempt a downloaded file came from.
type attemptAttachment struct {
	attempt    int
	attachment api.Attachment
}

// submissionAttachmentsAcrossAttempts returns every distinct file a student
// attached to this assignment, from the current attempt AND from each earlier
// attempt in submission_history. Deduplicated by attachment id (Canvas repeats
// a carried-forward file in every later attempt), keeping the earliest attempt
// that introduced it, and ordered by attempt so the manifest reads
// chronologically. A submission with no history behaves exactly as before.
func submissionAttachmentsAcrossAttempts(submission api.Submission) []attemptAttachment {
	var out []attemptAttachment
	seen := make(map[int64]bool)

	add := func(attempt int, attachments []api.Attachment) {
		for _, attachment := range attachments {
			if seen[attachment.ID] {
				continue
			}
			seen[attachment.ID] = true
			out = append(out, attemptAttachment{attempt: attempt, attachment: attachment})
		}
	}

	for _, past := range submission.SubmissionHistory {
		add(past.Attempt, past.Attachments)
	}
	add(submission.Attempt, submission.Attachments)

	return out
}

func safeSubmissionAttachmentFilename(attachment api.Attachment) string {
	filename := attachment.Filename
	if filename == "" {
		filename = attachment.DisplayName
	}
	filename = filepath.Base(strings.ReplaceAll(filename, "\\", "_"))
	filename = strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1
		}
		return r
	}, strings.TrimSpace(filename))
	if filename == "" || filename == "." || filename == ".." {
		return fmt.Sprintf("file-%d", attachment.ID)
	}
	return filename
}

func writeSubmissionDownloadManifest(path string, manifest submissionDownloadManifest) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600) // #nosec G304 -- path is under the user-selected destination.
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(manifest)
}
