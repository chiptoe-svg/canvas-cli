package commands

import (
	"strings"
	"testing"

	"github.com/chiptoe-svg/canvas-cli/internal/config"
)

// setupConfigTestHome isolates config for config command tests.
func setupConfigTestHome(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir) // Windows: os.UserHomeDir() reads %USERPROFILE%
	config.ResetCache()
	t.Cleanup(config.ResetCache)
}

func TestConfigListCmd_Empty(t *testing.T) {
	setupConfigTestHome(t)

	out := captureRunOutput(func() {
		cmd := newConfigListCmd()
		cmd.SetArgs([]string{})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "No instances") {
		t.Errorf("expected 'No instances' in output, got: %q", out)
	}
}

func TestConfigListCmd_WithInstances(t *testing.T) {
	setupConfigTestHome(t)

	cfg, _ := config.Load()
	_ = cfg.AddInstance(&config.Instance{Name: "prod", URL: "https://canvas.example.com"})
	config.ResetCache()

	out := captureRunOutput(func() {
		cmd := newConfigListCmd()
		cmd.SetArgs([]string{})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "prod") {
		t.Errorf("expected 'prod' in output, got: %q", out)
	}
}

func TestConfigAddCmd_Success(t *testing.T) {
	setupConfigTestHome(t)

	out := captureRunOutput(func() {
		cmd := newConfigAddCmd()
		cmd.SetArgs([]string{"staging", "--url", "https://staging.canvas.example.com"})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "staging") {
		t.Errorf("expected 'staging' in output, got: %q", out)
	}
}

func TestConfigAddCmd_MissingURL(t *testing.T) {
	setupConfigTestHome(t)

	cmd := newConfigAddCmd()
	cmd.SetArgs([]string{"staging"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when --url is missing")
	}
}

func TestConfigAddCmd_InvalidScheme(t *testing.T) {
	setupConfigTestHome(t)

	cmd := newConfigAddCmd()
	cmd.SetArgs([]string{"badscheme", "--url", "ftp://canvas.example.com"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for non-http/https URL scheme")
	}
}

func TestConfigUseCmd_Success(t *testing.T) {
	setupConfigTestHome(t)

	cfg, _ := config.Load()
	_ = cfg.AddInstance(&config.Instance{Name: "myprod", URL: "https://canvas.example.com"})
	config.ResetCache()

	out := captureRunOutput(func() {
		cmd := newConfigUseCmd()
		cmd.SetArgs([]string{"myprod"})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "myprod") {
		t.Errorf("expected 'myprod' in output, got: %q", out)
	}
}

func TestConfigUseCmd_NonExistent(t *testing.T) {
	setupConfigTestHome(t)

	cmd := newConfigUseCmd()
	cmd.SetArgs([]string{"doesnotexist"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for non-existent instance")
	}
}

func TestConfigRemoveCmd_WithForce(t *testing.T) {
	setupConfigTestHome(t)

	cfg, _ := config.Load()
	_ = cfg.AddInstance(&config.Instance{Name: "todelete", URL: "https://del.example.com"})
	config.ResetCache()

	out := captureRunOutput(func() {
		cmd := newConfigRemoveCmd()
		cmd.SetArgs([]string{"todelete", "--force"})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "todelete") {
		t.Errorf("expected 'todelete' in output, got: %q", out)
	}
}

func TestConfigRemoveCmd_NonExistent(t *testing.T) {
	setupConfigTestHome(t)

	cmd := newConfigRemoveCmd()
	cmd.SetArgs([]string{"ghost", "--force"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for non-existent instance")
	}
}

func TestConfigShowCmd_WithInstance(t *testing.T) {
	setupConfigTestHome(t)

	cfg, _ := config.Load()
	_ = cfg.AddInstance(&config.Instance{Name: "demo", URL: "https://demo.canvas.com", Description: "Demo env"})
	config.ResetCache()

	out := captureRunOutput(func() {
		cmd := newConfigShowCmd()
		cmd.SetArgs([]string{})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "demo") {
		t.Errorf("expected 'demo' in output, got: %q", out)
	}
}

func TestValueOrNone(t *testing.T) {
	if got := valueOrNone(""); got != "(none)" {
		t.Errorf("expected '(none)', got %q", got)
	}
	if got := valueOrNone("hello"); got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}
