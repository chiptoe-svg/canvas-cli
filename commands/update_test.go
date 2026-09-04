package commands

import (
	"testing"

	"github.com/jjuanrivvera/canvas-cli/internal/update"
)

// Version comparison lives in internal/update; the command must use the same
// rules as the background checker (including the +audited.N build suffix).
func TestUpdateCommandUsesSharedVersionComparison(t *testing.T) {
	if !update.IsNewerVersion("v1.13.0+audited.12", "1.13.0+audited.11") {
		t.Error("expected the newer audited build to be reported as newer")
	}
	if update.IsNewerVersion("v1.13.0", "1.13.0+audited.11") {
		t.Error("the upstream base version is not newer than an audited build of it")
	}
	if update.IsNewerVersion("v9.9.9", "dev") {
		t.Error("dev builds never update")
	}
}
