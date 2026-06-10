package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveSkillTarget(t *testing.T) {
	home, _ := os.UserHomeDir()

	// project (default)
	p, err := resolveSkillTarget("claude", "", false)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(".claude/skills", "canvas-cli"), p)

	// global
	p, err = resolveSkillTarget("claude", "", true)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(home, ".claude/skills", "canvas-cli"), p)

	// cursor project uses .agents/skills
	p, err = resolveSkillTarget("cursor", "", false)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(".agents/skills", "canvas-cli"), p)

	// explicit dir wins
	p, err = resolveSkillTarget("claude", t.TempDir(), true)
	require.NoError(t, err)
	assert.True(t, strings.HasSuffix(p, string(filepath.Separator)+"canvas-cli"))

	// unknown agent errors
	_, err = resolveSkillTarget("nope", "", false)
	assert.Error(t, err)
}

func TestSkillFiles(t *testing.T) {
	files, err := skillFiles()
	require.NoError(t, err)
	assert.Contains(t, files, "SKILL.md")
	assert.Contains(t, files, "references/canvas-commands.md")
	assert.Contains(t, files, "references/auth-and-config.md")
	assert.Contains(t, files, "references/output-and-filtering.md")
}

func TestSkillsInstallToDir(t *testing.T) {
	dir := t.TempDir()

	cmd := newSkillsInstallCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	require.NoError(t, cmd.Flags().Set("dir", dir))
	require.NoError(t, cmd.RunE(cmd, nil))

	installed := filepath.Join(dir, "canvas-cli", "SKILL.md")
	data, err := os.ReadFile(installed) //nolint:gosec // test-controlled path
	require.NoError(t, err)
	assert.Contains(t, string(data), "name: canvas-cli")
	assert.Contains(t, out.String(), "Installed canvas-cli skill")

	ref := filepath.Join(dir, "canvas-cli", "references", "canvas-commands.md")
	_, err = os.Stat(ref)
	require.NoError(t, err)
}

func TestSkillsPrint(t *testing.T) {
	cmd := newSkillsPrintCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	require.NoError(t, cmd.RunE(cmd, nil))
	assert.Contains(t, out.String(), "# Canvas CLI")
}

func TestSkillsInstallDryRun(t *testing.T) {
	dir := t.TempDir()

	origDryRun := dryRun
	dryRun = true
	defer func() { dryRun = origDryRun }()

	cmd := newSkillsInstallCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	require.NoError(t, cmd.Flags().Set("dir", dir))
	require.NoError(t, cmd.RunE(cmd, nil))

	assert.Contains(t, out.String(), "Would install skill to")
	_, err := os.Stat(filepath.Join(dir, "canvas-cli", "SKILL.md"))
	assert.True(t, os.IsNotExist(err), "dry-run must not write files")
}
