---
title: canvas submissions download
---

## canvas submissions download

Download all files attached to an assignment's submissions.

Files are organized by Canvas user ID and attachment ID, so identical student
filenames never overwrite one another. The command writes a
`submission-download-manifest.json` file that also records submissions with no
attachment, such as text-entry and URL submissions. The manifest does not copy
student text-entry bodies.

### Synopsis

```text
canvas submissions download [flags]
```

### Examples

```bash
canvas submissions download --course-id 123 --assignment-id 456 --destination ./essay-submissions
canvas submissions download --course-id 123 --assignment-id 456 --destination ./essay-submissions --overwrite
```

### Options

```text
      --assignment-id int     Assignment ID (required)
      --course-id int         Course ID (required)
      --destination string    Directory for downloaded submissions (required)
  -h, --help                  help for download
      --overwrite             Replace files that already exist
```

Existing files are skipped by default, making a rerun safe after an interruption.
Each download writes to a temporary `.partial` file and renames it only after a
successful transfer. The destination and user directories are created mode
`0700`; downloaded files and the manifest are mode `0600`.
