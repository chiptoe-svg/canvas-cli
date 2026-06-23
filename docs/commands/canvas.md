---
title: canvas
---

## canvas

Canvas LMS CLI - Interact with Canvas from the command line

### Synopsis

canvas-cli is a powerful command-line interface for Canvas LMS.
It provides comprehensive access to Canvas API features including courses,
assignments, users, submissions, and more.

Examples:

```bash
canvas auth login                                              # Authenticate with Canvas
canvas courses list                                            # List all courses
canvas assignments list --course-id 123                        # List assignments for a course
canvas submissions bulk-grade --course-id 123 --csv grades.csv # Bulk grade from CSV
```

### Options

```
      --as-user int       Masquerade as another user (admin feature, requires permission)
      --columns strings   Select specific columns to display (comma-separated)
      --config string     config file (default is $HOME/.canvas-cli/config.yaml)
      --dry-run           Print curl commands instead of executing requests
      --filter string     Filter results by text (case-insensitive substring match)
  -h, --help              help for canvas
      --instance string   Canvas instance URL (overrides config)
      --limit int         Limit number of results for list operations (0 = unlimited)
      --no-cache          Disable caching of API responses
  -o, --output string     Output format: table, json, yaml, csv (default "table")
      --quiet             Suppress informational messages, only output data and errors
      --show-token        Show actual token in dry-run output (default: redacted)
      --sort string       Sort results by field (prefix with - for descending, e.g., -name)
  -v, --verbose           Enable verbose output
```

### SEE ALSO

* [canvas account-analytics](canvas_account-analytics.md)	 - View Canvas account-level analytics
* [canvas account-calendars](canvas_account-calendars.md)	 - Manage Canvas account calendars
* [canvas account-content-migrations](canvas_account-content-migrations.md)	 - Manage account-level content migrations
* [canvas account-ext-tools-favorites](canvas_account-ext-tools-favorites.md)	 - Manage account external tool (LTI) favorites
* [canvas account-features](canvas_account-features.md)	 - Manage Canvas account feature flags and settings
* [canvas account-logins](canvas_account-logins.md)	 - Manage Canvas account logins
* [canvas account-notifications](canvas_account-notifications.md)	 - Manage Canvas account notifications
* [canvas account-reports](canvas_account-reports.md)	 - Manage Canvas account reports
* [canvas accounts](canvas_accounts.md)	 - Manage Canvas accounts
* [canvas admins](canvas_admins.md)	 - Manage account administrators
* [canvas alias](canvas_alias.md)	 - Manage command aliases
* [canvas analytics](canvas_analytics.md)	 - View Canvas analytics
* [canvas announcements](canvas_announcements.md)	 - Manage Canvas announcements
* [canvas api](canvas_api.md)	 - Make raw API requests to Canvas
* [canvas appointment-groups](canvas_appointment-groups.md)	 - Manage Canvas appointment groups
* [canvas assignment-groups](canvas_assignment-groups.md)	 - Manage Canvas assignment groups
* [canvas assignments](canvas_assignments.md)	 - Manage Canvas assignments
* [canvas audit](canvas_audit.md)	 - View Canvas audit logs
* [canvas auth](canvas_auth.md)	 - Manage authentication with Canvas
* [canvas auth-providers](canvas_auth-providers.md)	 - Manage Canvas authentication providers
* [canvas blackout-dates](canvas_blackout-dates.md)	 - Manage course blackout dates
* [canvas blueprint](canvas_blueprint.md)	 - Manage blueprint courses
* [canvas bookmarks](canvas_bookmarks.md)	 - Manage Canvas bookmarks
* [canvas brand](canvas_brand.md)	 - View Canvas brand/theme variables
* [canvas cache](canvas_cache.md)	 - Manage Canvas CLI cache
* [canvas calendar](canvas_calendar.md)	 - Manage Canvas calendar events
* [canvas collaborations](canvas_collaborations.md)	 - Manage Canvas collaborations
* [canvas comm-channels](canvas_comm-channels.md)	 - Manage communication channels
* [canvas comm-messages](canvas_comm-messages.md)	 - List Canvas communication messages
* [canvas completion](canvas_completion.md)	 - Generate shell completion scripts
* [canvas conferences](canvas_conferences.md)	 - List Canvas web conferences
* [canvas config](canvas_config.md)	 - Manage Canvas CLI configuration
* [canvas content-exports](canvas_content-exports.md)	 - Manage course content exports
* [canvas content-migrations](canvas_content-migrations.md)	 - Manage content migrations
* [canvas content-shares](canvas_content-shares.md)	 - Manage content shares
* [canvas context](canvas_context.md)	 - Manage working context (course, assignment, user IDs)
* [canvas conversations](canvas_conversations.md)	 - Manage Canvas conversations (inbox)
* [canvas course-extensions](canvas_course-extensions.md)	 - Manage course extensions (quiz and assignment)
* [canvas course-features](canvas_course-features.md)	 - Manage course feature flags
* [canvas course-nicknames](canvas_course-nicknames.md)	 - Manage course nicknames
* [canvas course-pacing](canvas_course-pacing.md)	 - Manage course pacing
* [canvas course-settings](canvas_course-settings.md)	 - Manage course settings and utilities
* [canvas courses](canvas_courses.md)	 - Manage Canvas courses
* [canvas csp-settings](canvas_csp-settings.md)	 - Manage Canvas Content Security Policy settings
* [canvas developer-keys](canvas_developer-keys.md)	 - Manage Canvas developer keys
* [canvas discussions](canvas_discussions.md)	 - Manage Canvas discussion topics
* [canvas doctor](canvas_doctor.md)	 - Run system diagnostics
* [canvas enrollment-terms](canvas_enrollment-terms.md)	 - Manage Canvas enrollment terms
* [canvas enrollments](canvas_enrollments.md)	 - Manage Canvas enrollments
* [canvas eportfolios](canvas_eportfolios.md)	 - Manage Canvas ePortfolios
* [canvas epub-exports](canvas_epub-exports.md)	 - Manage Canvas ePub exports
* [canvas error-reports](canvas_error-reports.md)	 - Submit Canvas error reports
* [canvas external-tools](canvas_external-tools.md)	 - Manage external tools (LTI)
* [canvas favorites](canvas_favorites.md)	 - Manage Canvas favorites
* [canvas files](canvas_files.md)	 - Manage Canvas files
* [canvas folders](canvas_folders.md)	 - Manage Canvas file folders
* [canvas grades](canvas_grades.md)	 - Manage Canvas gradebook
* [canvas grading-period-sets](canvas_grading-period-sets.md)	 - Manage Canvas grading period sets and grading periods
* [canvas grading-periods](canvas_grading-periods.md)	 - Manage course grading periods
* [canvas grading-standards](canvas_grading-standards.md)	 - Manage course grading standards
* [canvas groups](canvas_groups.md)	 - Manage Canvas groups
* [canvas history](canvas_history.md)	 - View Canvas user page-view history
* [canvas jwts](canvas_jwts.md)	 - Create and refresh Canvas JWTs
* [canvas live-assessments](canvas_live-assessments.md)	 - Manage course live assessments
* [canvas mcp](canvas_mcp.md)	 - MCP server management
* [canvas media](canvas_media.md)	 - Manage Canvas media objects and attachments
* [canvas modules](canvas_modules.md)	 - Manage Canvas course modules
* [canvas observees](canvas_observees.md)	 - Manage user observees and observers
* [canvas outcomes](canvas_outcomes.md)	 - Manage Canvas learning outcomes
* [canvas overrides](canvas_overrides.md)	 - Manage Canvas assignment overrides
* [canvas pages](canvas_pages.md)	 - Manage Canvas wiki pages
* [canvas peer-reviews](canvas_peer-reviews.md)	 - Manage peer reviews
* [canvas planner](canvas_planner.md)	 - Manage Canvas planner items and notes
* [canvas polls](canvas_polls.md)	 - Manage Canvas polls
* [canvas progress](canvas_progress.md)	 - Manage Canvas background job progress
* [canvas quizzes](canvas_quizzes.md)	 - Manage Canvas quizzes
* [canvas repl](canvas_repl.md)	 - Start interactive REPL mode
* [canvas roles](canvas_roles.md)	 - Manage account roles
* [canvas rubric-associations](canvas_rubric-associations.md)	 - Manage rubric associations and assessments
* [canvas rubrics](canvas_rubrics.md)	 - Manage Canvas rubrics
* [canvas sections](canvas_sections.md)	 - Manage Canvas course sections
* [canvas sis-imports](canvas_sis-imports.md)	 - Manage SIS imports
* [canvas skills](canvas_skills.md)	 - Install this CLI's AI-agent skill into Claude, Cursor, and other agents
* [canvas submissions](canvas_submissions.md)	 - Manage Canvas submissions
* [canvas sync](canvas_sync.md)	 - Synchronize resources between Canvas instances
* [canvas telemetry](canvas_telemetry.md)	 - Manage telemetry settings
* [canvas temporary-enrollment-pairings](canvas_temporary-enrollment-pairings.md)	 - Manage Canvas temporary enrollment pairings
* [canvas update](canvas_update.md)	 - Check for and install updates
* [canvas user-features](canvas_user-features.md)	 - Manage user feature flags
* [canvas users](canvas_users.md)	 - Manage Canvas users
* [canvas version](canvas_version.md)	 - Display version information
* [canvas webhook](canvas_webhook.md)	 - Manage Canvas webhook listeners

###### Auto generated by spf13/cobra on 23-Jun-2026
