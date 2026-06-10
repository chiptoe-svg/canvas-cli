// Package canvascli embeds the agent-skill assets (SKILL.md + references) into
// the binary so `canvas skills install` can write them into an AI agent's
// skills directory. The same files under skills/canvas-cli/ are what
// `npx skills add jjuanrivvera/canvas-cli` and the Claude Code plugin consume.
package canvascli

import (
	"embed"
	"io/fs"
)

//go:embed skills/canvas-cli
var embedded embed.FS

// SkillFS is rooted at the skill directory, so it contains SKILL.md and
// references/ at its top level.
var SkillFS = mustSub(embedded, "skills/"+SkillName)

// SkillName is the directory the skill installs into within an agent's skills dir.
const SkillName = "canvas-cli"

func mustSub(f embed.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
