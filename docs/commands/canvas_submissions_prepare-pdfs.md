---
title: canvas submissions prepare-pdfs
---

## canvas submissions prepare-pdfs

Classify local submission PDFs and prepare page images for review.

This command operates only on PDFs that have already been downloaded to this
Mac. It never contacts Canvas, Docling, Qwen, or another remote service.

### Synopsis

```text
canvas submissions prepare-pdfs --folder DIR --output DIR [flags]
```

### Examples

```bash
canvas submissions prepare-pdfs --folder ./assignment-456-submissions \
  --output ./assignment-456-review
```

### Options

```text
      --folder string   Folder containing already-downloaded submission PDFs (required)
  -h, --help            help for prepare-pdfs
      --output string   Directory for the local review manifest and page images (required)
      --overwrite       Replace an existing local PDF-review manifest and page images
```

The command records local text, font, and image signals in
`submission-pdf-manifest.jsonl`. When a PDF contains exactly one embedded image
per page, it preserves those original images. Other PDFs are rendered locally at
300 DPI. The manifest and images can contain FERPA-sensitive student records, so
keep the output in a restricted local location. An existing manifest is kept
unless `--overwrite` is explicitly supplied.
