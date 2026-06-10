package auth

import (
	"os"
	"testing"

	"github.com/zalando/go-keyring"
)

// TestMain stubs the browser launcher for the entire auth test binary.
//
// Several OAuth-flow tests exercise Authenticate / startLocalServer, which call
// openBrowser to launch the system browser. Left real, that spawns browser tabs
// on the developer's machine on every `go test` (it only stayed quiet in CI
// because xdg-open isn't installed on the Linux runner, so the best-effort
// Start() failed silently). Replace browserOpener with a no-op so no auth test
// opens a real browser; the few tests that specifically verify browserOpener
// save and restore it themselves.
func TestMain(m *testing.M) {
	browserOpener = func(string) {}
	// Use go-keyring's in-memory mock for the whole test binary. Without it,
	// token-store tests hit the real OS keyring: they pollute the developer's
	// keychain, and on macOS CI runners they hang forever — each test binary is
	// a distinct unsigned executable, so reading an item written by another
	// binary triggers a keychain ACL prompt that a headless runner can't answer.
	keyring.MockInit()
	os.Exit(m.Run())
}
