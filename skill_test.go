package canvascli

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSkillFS verifies the embedded skill bundle is rooted at the skill
// directory and ships the expected assets (exercises the package init/mustSub).
func TestSkillFS(t *testing.T) {
	require.NotNil(t, SkillFS)
	assert.Equal(t, "canvas-cli", SkillName)

	data, err := fs.ReadFile(SkillFS, "SKILL.md")
	require.NoError(t, err)
	assert.Contains(t, string(data), "name: canvas-cli")

	for _, ref := range []string{
		"references/canvas-commands.md",
		"references/auth-and-config.md",
		"references/output-and-filtering.md",
	} {
		_, err = fs.Stat(SkillFS, ref)
		require.NoError(t, err, ref)
	}
}
