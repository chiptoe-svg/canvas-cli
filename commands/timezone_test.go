package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/internal/config"
)

// withConfigTimezone points the config loader at a temp file whose
// settings.timezone is tz ("" for a config without the key).
func withConfigTimezone(t *testing.T, tz string) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	path, err := config.GetConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	body := "default_instance: x\ninstances: {}\nsettings:\n  default_output_format: table\n"
	if tz != "" {
		body += "  timezone: " + tz + "\n"
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	config.ResetCache()
	t.Cleanup(config.ResetCache)
}

func TestResolveTimezone(t *testing.T) {
	t.Run("flag wins over config", func(t *testing.T) {
		withConfigTimezone(t, "Europe/Berlin")
		if got := resolveTimezone("Asia/Tokyo"); got != "Asia/Tokyo" {
			t.Errorf("got %q", got)
		}
	})
	t.Run("config when no flag", func(t *testing.T) {
		withConfigTimezone(t, "Europe/Berlin")
		if got := resolveTimezone(""); got != "Europe/Berlin" {
			t.Errorf("got %q", got)
		}
		loc, err := resolveLocation("")
		if err != nil || loc.String() != "Europe/Berlin" {
			t.Errorf("resolveLocation = %v, %v", loc, err)
		}
	})
	t.Run("empty when config has no timezone", func(t *testing.T) {
		withConfigTimezone(t, "")
		if got := resolveTimezone(""); got != "" {
			t.Errorf("got %q", got)
		}
	})
	t.Run("empty falls through to TZ", func(t *testing.T) {
		withConfigTimezone(t, "")
		t.Setenv("TZ", "Asia/Tokyo")
		loc, err := resolveLocation("")
		if err != nil || loc.String() != "Asia/Tokyo" {
			t.Errorf("resolveLocation = %v, %v", loc, err)
		}
	})
	t.Run("bad zone is an error", func(t *testing.T) {
		withConfigTimezone(t, "")
		if _, err := resolveLocation("Nowhere/Nope"); err == nil {
			t.Error("want an error")
		}
	})
}

func TestAddTimezoneFlag(t *testing.T) {
	var tz string
	cmd := &cobra.Command{Use: "x"}
	addTimezoneFlag(cmd, &tz)
	f := cmd.Flags().Lookup("timezone")
	if f == nil {
		t.Fatal("--timezone not registered")
	}
	if f.Usage != timezoneFlagUsage {
		t.Errorf("usage = %q", f.Usage)
	}
	if err := cmd.Flags().Set("timezone", "America/Chicago"); err != nil || tz != "America/Chicago" {
		t.Errorf("set: %v, tz=%q", err, tz)
	}
}
