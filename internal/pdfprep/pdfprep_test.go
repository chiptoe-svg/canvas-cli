package pdfprep

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeRunner struct {
	run func(name string, args ...string) (CommandResult, error)
}

func (f fakeRunner) Run(_ context.Context, name string, args ...string) (CommandResult, error) {
	return f.run(name, args...)
}

func TestClassifyUsesLocalSignals(t *testing.T) {
	preparer := Preparer{Runner: fakeRunner{run: func(name string, args ...string) (CommandResult, error) {
		switch name {
		case "pdfinfo":
			return CommandResult{Stdout: "Pages:          2\n"}, nil
		case "pdftotext":
			return CommandResult{Stdout: strings.Repeat("a", 1200)}, nil
		case "pdffonts":
			return CommandResult{Stdout: realPdfFontsOutput}, nil
		case "pdfimages":
			return CommandResult{Stdout: "   0   1 image\n"}, nil
		default:
			t.Fatalf("unexpected command %q", name)
			return CommandResult{}, nil
		}
	}}}

	signals := preparer.Classify(context.Background(), "essay.pdf")
	if signals.Classification != TextHeavy {
		t.Fatalf("classification = %q, want %q", signals.Classification, TextHeavy)
	}
	// One of the two fonts in the real fixture is embedded (emb=yes); the field
	// counts EMBEDDED fonts, not fonts.
	if signals.Pages != 2 || signals.TextCharactersPerPage != 600 || signals.EmbeddedFontCount != 1 {
		t.Fatalf("unexpected signals: %#v", signals)
	}
}

func TestClassifyRecognizesImageOnlyPDF(t *testing.T) {
	preparer := Preparer{Runner: fakeRunner{run: func(name string, args ...string) (CommandResult, error) {
		switch name {
		case "pdfinfo":
			return CommandResult{Stdout: "Pages:          3\n"}, nil
		case "pdftotext", "pdffonts":
			return CommandResult{}, nil
		case "pdfimages":
			return CommandResult{Stdout: "   0   1 image\n   1   1 image\n   2   1 image\n"}, nil
		}
		return CommandResult{}, nil
	}}}

	if got := preparer.Classify(context.Background(), "notes.pdf").Classification; got != ScanOrImageHeavy {
		t.Fatalf("classification = %q, want %q", got, ScanOrImageHeavy)
	}
}

func TestClassifyRetainsToolFailure(t *testing.T) {
	preparer := Preparer{Runner: fakeRunner{run: func(name string, args ...string) (CommandResult, error) {
		if name == "pdftotext" {
			return CommandResult{Stderr: "malformed PDF"}, errors.New("exit status 1")
		}
		return CommandResult{}, nil
	}}}

	if got := preparer.Classify(context.Background(), "broken.pdf").LocalExtractionError; !strings.Contains(got, "pdftotext: malformed PDF") {
		t.Fatalf("local extraction error = %q", got)
	}
}

func TestExtractPageImagesPreservesOriginalPageImages(t *testing.T) {
	destination := t.TempDir()
	preparer := Preparer{Runner: fakeRunner{run: func(name string, args ...string) (CommandResult, error) {
		if name == "pdfimages" && len(args) == 3 && args[0] == "-all" {
			for _, suffix := range []string{"000.jpg", "001.jpg"} {
				if err := os.WriteFile(args[2]+"-"+suffix, []byte("original"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			return CommandResult{}, nil
		}
		if name == "pdftoppm" {
			t.Fatal("pdftoppm must not run when one original image exists per page")
		}
		return CommandResult{}, nil
	}}}

	images, err := preparer.ExtractPageImages(context.Background(), "notes.pdf", destination, 2)
	if err != nil {
		t.Fatalf("ExtractPageImages: %v", err)
	}
	if images.Source != EmbeddedImagePages || len(images.Paths) != 2 {
		t.Fatalf("images = %#v", images)
	}
}

func TestExtractPageImagesRendersWhenEmbeddedImagesDoNotMatchPages(t *testing.T) {
	destination := t.TempDir()
	preparer := Preparer{Runner: fakeRunner{run: func(name string, args ...string) (CommandResult, error) {
		switch name {
		case "pdfimages":
			if err := os.WriteFile(args[2]+"-000.jpg", []byte("figure"), 0o600); err != nil {
				t.Fatal(err)
			}
			return CommandResult{}, nil
		case "pdftoppm":
			if got, want := strings.Join(args[:3], " "), "-png -r 300"; got != want {
				t.Fatalf("render args = %q, want %q", got, want)
			}
			for _, suffix := range []string{"001.png", "002.png"} {
				if err := os.WriteFile(args[4]+"-"+suffix, []byte("rendered"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			return CommandResult{}, nil
		}
		return CommandResult{}, nil
	}}}

	images, err := preparer.ExtractPageImages(context.Background(), "mixed.pdf", destination, 2)
	if err != nil {
		t.Fatalf("ExtractPageImages: %v", err)
	}
	if images.Source != RenderedPageImages || len(images.Paths) != 2 {
		t.Fatalf("images = %#v", images)
	}
	if filepath.Base(images.Paths[0]) != "rendered-001.png" {
		t.Fatalf("first rendered page = %q", images.Paths[0])
	}
}


// VERBATIM `pdffonts` output from Poppler 25.x. The fixture it replaced was
// invented — rows beginning with a number — and the code matched that
// invention, so the font count was always zero against a real PDF and the
// text-heavy classification could never be reached. Pin the real shape: rows
// begin with the font NAME, and the embedded flag is the `emb` column.
const realPdfFontsOutput = `name                                 type              encoding         emb sub uni object ID
------------------------------------ ----------------- ---------------- --- --- --- ---------
AAAAAB+Monaco                        TrueType          MacRoman         yes yes no        6  0
Helvetica                            Type 1            WinAnsi          no  no  no        9  0
`

func TestCountEmbeddedFonts(t *testing.T) {
	if got := countEmbeddedFonts(realPdfFontsOutput); got != 1 {
		t.Fatalf("countEmbeddedFonts(real output) = %d, want 1 (only the emb=yes row)", got)
	}
	if got := countEmbeddedFonts(""); got != 0 {
		t.Errorf("countEmbeddedFonts(empty) = %d, want 0", got)
	}
	// Degrade to a row count rather than a false zero if the format changes.
	unknown := "name  type\n----  ----\nSomeFont  Type1\nOther  TrueType\n"
	if got := countEmbeddedFonts(unknown); got != 2 {
		t.Errorf("countEmbeddedFonts(unrecognized format) = %d, want 2 rows", got)
	}
}

// The whole point of the signal: a printed, text-rich PDF must reach
// text-heavy, which requires a non-zero font count.
func TestClassify_TextHeavyIsReachable(t *testing.T) {
	got := classify(Signals{Pages: 4, TextCharacters: 16320, TextCharactersPerPage: 4080, EmbeddedFontCount: 1})
	if got != TextHeavy {
		t.Fatalf("classification = %q, want %q", got, TextHeavy)
	}
}
