package canvascli

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSkillFS verifies the embedded skill bundle ships SKILL.md plus every
// reference SKILL.md links, and that the folded-away reference is gone.
// SkillFS is rooted at the skill directory, so paths are relative to it.
func TestSkillFS(t *testing.T) {
	require.NotNil(t, SkillFS)
	assert.Equal(t, "canvas-cli", SkillName)

	data, err := fs.ReadFile(SkillFS, "SKILL.md")
	require.NoError(t, err)
	skill := string(data)
	assert.Contains(t, skill, "name: canvas-cli")
	assert.Contains(t, skill, "release/audited/install.sh")

	refs := []string{
		"references/canvas-commands.md",
		"references/auth-and-config.md",
		"references/output-and-filtering.md",
		"references/grading-week.md",
		"references/term-setup.md",
		"references/mid-term-check.md",
		"references/accommodations.md",
	}
	for _, ref := range refs {
		assert.Contains(t, skill, ref, "SKILL.md must link every reference")
		_, err := fs.Stat(SkillFS, ref)
		require.NoError(t, err, ref)
	}
	_, err = fs.Stat(SkillFS, "references/grading-workflows.md")
	assert.Error(t, err, "grading-workflows.md was folded into grading-week.md")
}

// Every `canvas <group>` the skill mentions must be a command the faculty
// build has. Groups the trim removed must not linger in the guidance.
func TestSkillMentionsOnlyFacultyCommands(t *testing.T) {
	allowed := map[string]bool{}
	for _, name := range strings.Fields(`activity agent alias analytics announcements api
		appointment-groups assignment-groups assignments auth cache calendar collaborations
		completion config content-exports content-migrations content-shares context conversations
		course-extensions course-features courses discussions doctor enrollments files folders grades grading-periods
		grading-standards groups help modules outcomes overrides pages peer-reviews quizzes
		rubric-associations rubrics schedule sections skills submissions update users version`) {
		allowed[name] = true
	}
	// Both forms count: a backtick-quoted mention in prose, and a command
	// line at the start of a line inside a ```bash fence — which is where
	// most of the command lines in the bundle actually live.
	re := regexp.MustCompile("(?m)(?:`|^\\s*)canvas ([a-z][a-z-]*)")
	err := fs.WalkDir(SkillFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		data, err := fs.ReadFile(SkillFS, path)
		if err != nil {
			return err
		}
		for _, m := range re.FindAllStringSubmatch(string(data), -1) {
			assert.Truef(t, allowed[m[1]], "%s mentions `canvas %s`, which the faculty build does not have", path, m[1])
		}
		return nil
	})
	require.NoError(t, err)
}
