// Package pdfprep classifies local PDFs and produces page images for review.
// It intentionally has no Canvas or network dependency.
package pdfprep

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	TextHeavyCharsPerPage = 500
	ScanCharsPerPage      = 75
)

type Classification string

const (
	TextHeavy          Classification = "text-heavy"
	ScanOrImageHeavy   Classification = "scan-or-image-heavy"
	HybridOrUncertain  Classification = "hybrid-or-uncertain"
	EmbeddedImagePages string         = "embedded"
	RenderedPageImages string         = "rendered"
)

// Signals are local, inexpensive indicators. They are routing hints, not a
// grade or a conclusion about the quality of a student's work.
type Signals struct {
	Pages                 int            `json:"pages,omitempty"`
	TextCharacters        int            `json:"text_characters"`
	TextCharactersPerPage float64        `json:"text_characters_per_page,omitempty"`
	EmbeddedFontCount     int            `json:"embedded_font_count,omitempty"`
	ImageCount            int            `json:"image_count,omitempty"`
	Classification        Classification `json:"classification"`
	LocalExtractionError  string         `json:"local_extraction_error,omitempty"`
}

type PageImages struct {
	Source string   `json:"source"`
	Paths  []string `json:"paths"`
}

// CommandResult preserves stderr so callers can return a useful diagnostic
// when Poppler rejects a malformed PDF.
type CommandResult struct {
	Stdout string
	Stderr string
}

// CommandRunner makes the Poppler boundary testable without requiring the
// utilities during unit tests.
type CommandRunner interface {
	Run(ctx context.Context, name string, args ...string) (CommandResult, error)
}

type execRunner struct{}

func (execRunner) Run(ctx context.Context, name string, args ...string) (CommandResult, error) {
	command := exec.CommandContext(ctx, name, args...)
	stdout, err := command.Output()
	result := CommandResult{Stdout: string(stdout)}
	if exitError, ok := err.(*exec.ExitError); ok {
		result.Stderr = string(exitError.Stderr)
	}
	return result, err
}

type Preparer struct {
	Runner CommandRunner
}

func New() Preparer {
	return Preparer{Runner: execRunner{}}
}

func (p Preparer) runner() CommandRunner {
	if p.Runner == nil {
		return execRunner{}
	}
	return p.Runner
}

func (p Preparer) Classify(ctx context.Context, path string) Signals {
	signals := Signals{Classification: HybridOrUncertain}
	var failures []string

	info, err := p.runner().Run(ctx, "pdfinfo", path)
	if err != nil {
		failures = append(failures, commandFailure("pdfinfo", info, err))
	} else if match := regexp.MustCompile(`(?m)^Pages:\s+(\d+)\s*$`).FindStringSubmatch(info.Stdout); len(match) == 2 {
		_, _ = fmt.Sscanf(match[1], "%d", &signals.Pages)
	}

	text, err := p.runner().Run(ctx, "pdftotext", "-layout", path, "-")
	if err != nil {
		failures = append(failures, commandFailure("pdftotext", text, err))
	} else {
		signals.TextCharacters = len(regexp.MustCompile(`\s+`).ReplaceAllString(text.Stdout, ""))
	}

	fonts, err := p.runner().Run(ctx, "pdffonts", path)
	if err != nil {
		failures = append(failures, commandFailure("pdffonts", fonts, err))
	} else {
		signals.EmbeddedFontCount = countEmbeddedFonts(fonts.Stdout)
	}

	images, err := p.runner().Run(ctx, "pdfimages", "-list", path)
	if err != nil {
		failures = append(failures, commandFailure("pdfimages", images, err))
	} else {
		signals.ImageCount = countRows(images.Stdout, `^\s*\d+\s+\d+\s+`)
	}

	if signals.Pages > 0 {
		signals.TextCharactersPerPage = float64(signals.TextCharacters) / float64(signals.Pages)
	}
	signals.Classification = classify(signals)
	signals.LocalExtractionError = strings.Join(failures, "; ")
	return signals
}

func classify(signals Signals) Classification {
	if signals.Pages == 0 {
		return HybridOrUncertain
	}
	if signals.TextCharactersPerPage >= TextHeavyCharsPerPage && signals.EmbeddedFontCount > 0 {
		return TextHeavy
	}
	if signals.TextCharactersPerPage < ScanCharsPerPage || (signals.EmbeddedFontCount == 0 && signals.ImageCount >= signals.Pages) {
		return ScanOrImageHeavy
	}
	return HybridOrUncertain
}

func (p Preparer) ExtractPageImages(ctx context.Context, pdfPath, destination string, pages int) (PageImages, error) {
	if pages < 1 {
		return PageImages{}, fmt.Errorf("cannot extract page images: PDF page count is unavailable")
	}
	if err := os.MkdirAll(destination, 0o700); err != nil {
		return PageImages{}, fmt.Errorf("create image destination: %w", err)
	}
	if err := os.Chmod(destination, 0o700); err != nil {
		return PageImages{}, fmt.Errorf("protect image destination: %w", err)
	}

	embeddedPrefix := filepath.Join(destination, "embedded")
	_, embeddedErr := p.runner().Run(ctx, "pdfimages", "-all", pdfPath, embeddedPrefix)
	embedded := matchingFiles(embeddedPrefix + "-*")
	// A photographed-notes PDF commonly has exactly one full-resolution image
	// per page. Keep those originals instead of needlessly rasterizing them.
	if embeddedErr == nil && len(embedded) == pages {
		return PageImages{Source: EmbeddedImagePages, Paths: embedded}, nil
	}

	renderedPrefix := filepath.Join(destination, "rendered")
	renderedResult, renderedErr := p.runner().Run(ctx, "pdftoppm", "-png", "-r", "300", pdfPath, renderedPrefix)
	if renderedErr != nil {
		return PageImages{}, fmt.Errorf("render PDF pages: %s", commandFailure("pdftoppm", renderedResult, renderedErr))
	}
	rendered := matchingFiles(renderedPrefix + "-*.png")
	if len(rendered) == 0 {
		return PageImages{}, fmt.Errorf("render PDF pages: pdftoppm produced no images")
	}
	return PageImages{Source: RenderedPageImages, Paths: rendered}, nil
}

// pdffonts prints two header lines, then one row per font whose columns end
// with `emb sub uni objectID generation`:
//
//	name                       type       encoding   emb sub uni object ID
//	-------------------------- ---------- ---------- --- --- --- ---------
//	AAAAAB+Monaco              TrueType   MacRoman   yes yes no        6  0
//
// Rows start with the font NAME, never a number — an earlier `^\s*\d+` match
// therefore counted zero fonts for every real PDF, which silently made the
// text-heavy classification unreachable. Count the rows whose `emb` column
// says yes; if the tail pattern ever stops matching (a pdffonts format
// change), fall back to the number of non-header rows so the signal degrades
// to "some fonts" rather than to a false zero.
func countEmbeddedFonts(stdout string) int {
	lines := strings.Split(stdout, "\n")
	start := 0
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "---") {
			start = i + 1
			break
		}
	}
	tail := regexp.MustCompile(`(?i)\b(yes|no)\s+(yes|no)\s+(yes|no)\s+\d+\s+\d+\s*$`)
	embedded, rows := 0, 0
	for _, line := range lines[start:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		rows++
		if match := tail.FindStringSubmatch(line); match != nil && strings.EqualFold(match[1], "yes") {
			embedded++
		}
	}
	if embedded == 0 && rows > 0 && !tail.MatchString(stdout) {
		return rows
	}
	return embedded
}

func countRows(value, expression string) int {
	pattern := regexp.MustCompile(expression)
	count := 0
	for _, line := range strings.Split(value, "\n") {
		if pattern.MatchString(line) {
			count++
		}
	}
	return count
}

func matchingFiles(pattern string) []string {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil
	}
	sort.Strings(paths)
	return paths
}

func commandFailure(name string, result CommandResult, err error) string {
	detail := strings.TrimSpace(result.Stderr)
	if detail == "" {
		detail = err.Error()
	}
	return fmt.Sprintf("%s: %s", name, detail)
}
