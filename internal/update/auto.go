package update

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

// AutoUpdater runs the periodic background check for a newer release and
// leaves a notification. It never installs: installing replaces the running
// binary with one verified only against a checksum file fetched from the
// same place, so that stays an explicit `canvas update`.
type AutoUpdater struct {
	CurrentVersion string
	CheckInterval  time.Duration
	Enabled        bool

	updater      *Updater
	stateManager *StateManager
	done         chan struct{} // Signals when async update check completes
}

// NewAutoUpdater creates a new AutoUpdater
func NewAutoUpdater(currentVersion string, enabled bool, checkInterval time.Duration) (*AutoUpdater, error) {
	stateManager, err := NewStateManager()
	if err != nil {
		return nil, err
	}

	if checkInterval == 0 {
		checkInterval = DefaultCheckInterval
	}

	return &AutoUpdater{
		CurrentVersion: currentVersion,
		CheckInterval:  checkInterval,
		Enabled:        enabled,
		updater:        NewUpdater(currentVersion),
		stateManager:   stateManager,
	}, nil
}

// RunUpdateCheck performs the background check (no download, no install)
// and records the outcome for PrintNotifications. Called at CLI startup.
func (a *AutoUpdater) RunUpdateCheck(ctx context.Context) {
	if !a.Enabled {
		return
	}

	// Skip dev versions
	if a.CurrentVersion == "dev" || a.CurrentVersion == "" {
		return
	}

	// Check if we should perform a check
	if !a.stateManager.ShouldCheck(a.CheckInterval) {
		return
	}

	checkCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	result, _ := a.updater.Check(checkCtx)

	// Record the check (ignore error - update state is non-critical)
	_ = a.stateManager.RecordCheck(result)
}

// RunUpdateCheckAsync runs the update check in a goroutine
// This prevents blocking the main CLI execution
func (a *AutoUpdater) RunUpdateCheckAsync(ctx context.Context) {
	a.done = make(chan struct{})
	go func() {
		defer close(a.done)
		// Recover from any panics to not crash the CLI
		defer func() {
			_ = recover() // Silently ignore panics in the updater
		}()

		a.RunUpdateCheck(ctx)
	}()
}

// WaitForCompletion waits for the async update check to complete with a timeout
// Returns true if completed, false if timeout was reached
func (a *AutoUpdater) WaitForCompletion(timeout time.Duration) bool {
	if a.done == nil {
		return true // No async operation was started
	}
	select {
	case <-a.done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// PrintNotifications prints any pending update notifications to stderr
// This should be called after command execution
func (a *AutoUpdater) PrintNotifications() {
	// Check for successful update notification
	if fromVersion, toVersion, hasNotification := a.stateManager.GetPendingNotification(); hasNotification {
		fmt.Fprintf(os.Stderr, "\n\033[32m✓ Updated canvas-cli from v%s to v%s\033[0m\n", fromVersion, toVersion)
		fmt.Fprintf(os.Stderr, "  Run 'canvas version' to verify the new version.\n")
		return
	}

	if version, ok := a.stateManager.GetAvailableNotification(); ok {
		fmt.Fprintf(os.Stderr, "\n\033[33mA newer canvas-cli is available: v%s (you have v%s).\033[0m\n", version, strings.TrimPrefix(a.CurrentVersion, "v"))
		fmt.Fprintf(os.Stderr, "  %s\n", InstallHint)
		return
	}

	// Check for error notification
	if errMsg, hasError := a.stateManager.GetLastError(); hasError {
		fmt.Fprintf(os.Stderr, "\n\033[33m⚠ Update check failed: %s\033[0m\n", errMsg)
		fmt.Fprintf(os.Stderr, "  %s\n", InstallHint)
	}
}

// CheckNow performs an immediate update check, ignoring the interval
// This is useful for manual update commands
func (a *AutoUpdater) CheckNow(ctx context.Context) *UpdateResult {
	result := a.updater.CheckAndUpdate(ctx)
	_ = a.stateManager.RecordCheck(result) // Ignore error - state is non-critical
	return result
}

// GetState returns the current update state
func (a *AutoUpdater) GetState() (*State, error) {
	return a.stateManager.Load()
}
