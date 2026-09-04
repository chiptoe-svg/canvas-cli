# PDF review with local preparation and Spark

Use this reference only when an instructor asks to review a batch of downloaded
student PDFs with Docling, Qwen, or another extraction service. It is not a
grading shortcut: the local signals and model output are evidence to review
against the live assignment and rubric.

## Boundary and authorization

`canvas submissions prepare-pdfs` is local-only. It reads the chosen folder,
writes a manifest and staged page images, and makes no Canvas or network call.

Uploading student work to Spark is a separate action. Before any upload, obtain
explicit authorization for the named course/assignment or local folder, the
files or pages in scope, the endpoint, the selected service, and the purpose.
Do not treat a general request to grade as permission to send submissions to a
remote service. Never put Canvas credentials in a request, prompt, URL, or
output. The manifest, raw responses, extracted text, and staged images may all
be FERPA-sensitive records; store them in a restricted local directory and do
not paste their contents into chat unless the instructor asks.

## Prepare first

Download every submission attempt only after confirming the Canvas course,
assignment, and local destination with the instructor. Then prepare the local
folder:

```bash
canvas submissions download --course-id 123 --assignment-id 456 \
  --destination ./assignment-456-submissions --no-cache
canvas submissions prepare-pdfs --folder ./assignment-456-submissions \
  --output ./assignment-456-review
```

Read `submission-pdf-manifest.jsonl`. It records each source file, its hash,
page count, text/font/image signals, classification, and page-image paths.
The paths use forward slashes so a saved manifest is portable across machines.
Use the staged original embedded images when `page_images.source` is
`embedded`; do not downsample photographed notes before visual review. A
`rendered` value is a 300-DPI local fallback for vector, mixed, and other PDFs.

## Route work deliberately

The classifier is a routing hint, not handwriting detection or a quality score.

| Local signal | Normal first pass | Escalate when |
|---|---|---|
| `text-heavy` | Docling / Granite Docling | required text is missing, the rubric needs visual evidence, or the document is not actually printed text |
| `scan-or-image-heavy` | Qwen for handwritten notes, hard scans, or visual evidence; otherwise Docling OCR first | output is sparse, unreadable, truncated, or needs visual interpretation |
| `hybrid-or-uncertain` | inspect the assignment type and use Docling first for printed prose | layout, handwriting, images, charts, or sparse extraction are material to the rubric |

For a visual/design rubric, Qwen may be appropriate even for a text-heavy PDF.
For a handwriting-heavy assignment profile, route each staged original page to
Qwen directly after the upload has been authorized. Do not use model routing to
silently change the scope of files or pages approved by the instructor.

## Clemson Spark services

The approved local-service endpoints for this workflow are:

- Docling: `http://gcspark.clemson.edu:5001/v1/convert/file`
- Qwen: `http://gcspark.clemson.edu:8080/v1/chat/completions`

For Docling, send one PDF with the service's supported multipart conversion
request. Save the returned structured extraction next to the manifest and
record failure or sparse-text conditions; do not claim that an OCR pass read
unreadable material.

For Qwen, send one page image per request, using the original staged image when
available. Set the required `X-Client` header to the calling application name,
use model `qwen3.6-35b-a3b`, temperature `0`, and keep thinking disabled. Use a
faithful instruction in the user message: ask for exactly what is visible,
with no interpretation or inference, and ask it to identify unreadable text.
For full-page handwritten transcription, use an output budget of 3000–4000
tokens and check `finish_reason`; retry only a length-truncated page with a
larger budget. Until the caller correctly parses server-sent events, use a
non-streaming response and preserve the raw response before parsing it.

Qwen is particularly useful for real handwriting and difficult scans. Docling
is generally the cheaper first path for printed documents. Docling and Qwen may
run concurrently on Spark when the instructor has authorized both workloads,
but preserve page order, avoid swapping unrelated heavyweight models through a
batch, bound concurrency, and make each completed result resumable by source
hash and page number.

## Evidence before grades

Treat extracted text and vision output as evidence only. Resolve the live
assignment and rubric before scoring. Produce a local proposal CSV with the
student identifier, per-criterion proposed points, total, concise evidence,
questions or uncertainty, and the source/page reference. Do not post grades,
rubric assessments, or feedback until the instructor reviews the proposal and
explicitly approves the exact rows to apply. After approval, use the CLI's
write-and-read-back workflow and report each verified result.
