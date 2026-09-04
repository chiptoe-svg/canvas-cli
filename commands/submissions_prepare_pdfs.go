package commands

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/pdfprep"
)

const submissionPDFManifestName = "submission-pdf-manifest.jsonl"

type submissionPDFRecord struct {
	SourceRelativePath string             `json:"source_relative_path"`
	SourceSHA256       string             `json:"source_sha256"`
	SourceBytes        int64              `json:"source_bytes"`
	Signals            pdfprep.Signals    `json:"signals"`
	PageImages         pdfprep.PageImages `json:"page_images,omitempty"`
	Status             string             `json:"status"`
	Error              string             `json:"error,omitempty"`
}

type submissionPDFManifest struct {
	SchemaVersion int                   `json:"schema_version"`
	GeneratedAt   time.Time             `json:"generated_at"`
	SourceFolder  string                `json:"source_folder"`
	Entries       []submissionPDFRecord `json:"entries"`
}

func newSubmissionsPreparePDFsCmd() *cobra.Command {
	opts := &options.SubmissionsPreparePDFsOptions{}
	cmd := &cobra.Command{
		Use:   "prepare-pdfs --folder DIR --output DIR",
		Short: "Classify local submission PDFs and prepare page images for review",
		Long: `Classify PDFs that have already been downloaded locally, then prepare their pages for a local review workflow.

This command never contacts Canvas or any remote extraction/model service. It
uses local Poppler utilities to measure embedded text, fonts, and images. For
photographed notes with one embedded image per PDF page, it retains the original
page images. Other PDFs are rendered locally at 300 DPI. The manifest and page
images may contain student records; store the output in a restricted location.

Examples:
  canvas submissions prepare-pdfs --folder ./assignment-456-submissions --output ./assignment-456-review
  canvas submissions prepare-pdfs --folder ./assignment-456-submissions --output ./assignment-456-review --overwrite`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return runSubmissionsPreparePDFs(cmd.Context(), opts, pdfprep.New())
		},
	}
	cmd.Flags().StringVar(&opts.Folder, "folder", "", "Folder containing already-downloaded submission PDFs (required)")
	cmd.Flags().StringVar(&opts.Output, "output", "", "Directory for the local review manifest and page images (required)")
	cmd.Flags().BoolVar(&opts.Overwrite, "overwrite", false, "Replace an existing local PDF-review manifest and page images")
	mustMarkRequired(cmd, "folder", "output")
	return cmd
}

func runSubmissionsPreparePDFs(ctx context.Context, opts *options.SubmissionsPreparePDFsOptions, preparer pdfprep.Preparer) error {
	root, err := filepath.Abs(opts.Folder)
	if err != nil {
		return fmt.Errorf("resolve folder: %w", err)
	}
	if info, err := os.Stat(root); err != nil {
		return fmt.Errorf("inspect folder: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("folder is not a directory: %s", root)
	}
	output, err := filepath.Abs(opts.Output)
	if err != nil {
		return fmt.Errorf("resolve output: %w", err)
	}
	if output == root {
		return fmt.Errorf("output must be different from the source folder")
	}
	manifestPath := filepath.Join(output, submissionPDFManifestName)
	if !opts.Overwrite {
		if _, err := os.Stat(manifestPath); err == nil {
			return fmt.Errorf("local review manifest already exists at %s; rerun with --overwrite to replace it", manifestPath)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("inspect existing manifest: %w", err)
		}
	}
	if err := os.MkdirAll(output, 0o700); err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	if err := os.Chmod(output, 0o700); err != nil {
		return fmt.Errorf("protect output: %w", err)
	}

	pdfs, err := findPDFs(root, output)
	if err != nil {
		return err
	}
	manifest := submissionPDFManifest{SchemaVersion: 1, GeneratedAt: time.Now().UTC(), SourceFolder: root}
	prepared, failed := 0, 0
	for _, path := range pdfs {
		record, err := preparePDF(ctx, root, output, path, preparer)
		if err != nil {
			record.Status, record.Error = "failed", err.Error()
			failed++
		} else {
			record.Status = "prepared"
			prepared++
		}
		manifest.Entries = append(manifest.Entries, record)
	}
	if err := writeSubmissionPDFManifest(manifestPath, manifest); err != nil {
		return fmt.Errorf("write local PDF-review manifest: %w", err)
	}
	fmt.Printf("Prepared %d PDF(s), failed %d. Manifest: %s\n", prepared, failed, manifestPath)
	if failed > 0 {
		return fmt.Errorf("%d PDF(s) could not be prepared; see %s", failed, manifestPath)
	}
	return nil
}

func findPDFs(root, output string) ([]string, error) {
	var pdfs []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == output && entry.IsDir() {
			return filepath.SkipDir
		}
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".pdf") {
			return nil
		}
		pdfs = append(pdfs, path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("find PDFs: %w", err)
	}
	sort.Strings(pdfs)
	return pdfs, nil
}

func preparePDF(ctx context.Context, root, output, path string, preparer pdfprep.Preparer) (submissionPDFRecord, error) {
	relativePath, err := filepath.Rel(root, path)
	if err != nil {
		return submissionPDFRecord{}, fmt.Errorf("resolve relative path: %w", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		return submissionPDFRecord{SourceRelativePath: relativePath}, fmt.Errorf("inspect PDF: %w", err)
	}
	digest, err := fileSHA256(path)
	if err != nil {
		return submissionPDFRecord{SourceRelativePath: relativePath}, fmt.Errorf("hash PDF: %w", err)
	}
	record := submissionPDFRecord{
		SourceRelativePath: relativePath,
		SourceSHA256:       digest,
		SourceBytes:        info.Size(),
		Signals:            preparer.Classify(ctx, path),
	}
	images, err := preparer.ExtractPageImages(ctx, path, filepath.Join(output, "pages", digest), record.Signals.Pages)
	if err != nil {
		return record, err
	}
	for index, imagePath := range images.Paths {
		if err := os.Chmod(imagePath, 0o600); err != nil {
			return record, fmt.Errorf("protect page image: %w", err)
		}
		relativeImagePath, err := filepath.Rel(output, imagePath)
		if err != nil {
			return record, fmt.Errorf("resolve page image path: %w", err)
		}
		images.Paths[index] = relativeImagePath
	}
	record.PageImages = images
	return record, nil
}

func fileSHA256(path string) (string, error) {
	file, err := os.Open(path) // #nosec G304 -- the user supplies the local PDF folder.
	if err != nil {
		return "", err
	}
	defer file.Close()
	digest := sha256.New()
	if _, err := io.Copy(digest, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}

func writeSubmissionPDFManifest(path string, manifest submissionPDFManifest) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600) // #nosec G304 -- path is under the user-selected output directory.
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(manifest)
}
