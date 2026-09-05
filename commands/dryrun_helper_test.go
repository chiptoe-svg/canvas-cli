package commands

import (
	"os"
	"testing"
)

// withGlobalDryRun sets the dryRun global (the --dry-run persistent flag)
// and restores it via t.Cleanup. Command constructors under test are not
// attached to rootCmd, so the flag itself is not available to them.
func withGlobalDryRun(t *testing.T, value bool) {
	t.Helper()
	orig := dryRun
	dryRun = value
	t.Cleanup(func() { dryRun = orig })
}

// captureRunOutput captures anything fn writes to os.Stdout and returns it.
// Shared by tests across the package (originally defined in telemetry_test.go).
func captureRunOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	buf := make([]byte, 16384)
	n, _ := r.Read(buf)
	r.Close()
	return string(buf[:n])
}
