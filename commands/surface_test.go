package commands

import (
	"sort"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The faculty surface. Adding or removing a command means editing this list
// in the same commit: the tool's scope is documented here and cannot drift.
var facultySurface = []string{
	"activity", "agent", "alias", "analytics", "announcements", "api",
	"appointment-groups", "assignment-groups", "assignments", "auth", "cache",
	"calendar", "collaborations", "completion", "config", "content-exports",
	"content-migrations", "content-shares", "context", "conversations",
	"course-extensions", "course-features", "courses", "discussions", "doctor",
	"enrollments",
	"files", "folders", "grades", "grading-periods", "grading-standards",
	"groups", "help", "modules", "outcomes", "overrides", "pages",
	"peer-reviews", "quizzes", "rubric-associations", "rubrics", "schedule",
	"sections", "skills", "submissions", "update", "users", "version",
}

var facultyUsersSubcommands = []string{
	"activity-stream", "get", "list", "me", "missing-submissions", "profile",
	"search", "todo", "upcoming-events",
}

var facultyAPISubcommands = []string{"get"}

func commandNames(cmd *cobra.Command) []string {
	var names []string
	for _, c := range cmd.Commands() {
		names = append(names, c.Name())
	}
	sort.Strings(names)
	return names
}

func findCommand(t *testing.T, parent *cobra.Command, name string) *cobra.Command {
	t.Helper()
	for _, c := range parent.Commands() {
		if c.Name() == name {
			return c
		}
	}
	require.Failf(t, "command missing", "%s has no %q subcommand", parent.Name(), name)
	return nil
}

func TestFacultySurface(t *testing.T) {
	rootCmd.InitDefaultHelpCmd()
	assert.Equal(t, facultySurface, commandNames(rootCmd), "top-level commands")
	assert.Equal(t, facultyUsersSubcommands, commandNames(findCommand(t, rootCmd, "users")), "users subcommands")
	assert.Equal(t, facultyAPISubcommands, commandNames(findCommand(t, rootCmd, "api")), "api subcommands")
}
