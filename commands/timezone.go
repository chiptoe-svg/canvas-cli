package commands

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/config"
	"github.com/jjuanrivvera/canvas-cli/internal/localtime"
)

// timezoneFlagUsage is the help text every --timezone flag shares.
const timezoneFlagUsage = "IANA time zone for local times, e.g. America/New_York (default: settings.timezone in the config, then $TZ, then the system zone)"

// addTimezoneFlag registers the per-command --timezone flag.
func addTimezoneFlag(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVar(p, "timezone", "", timezoneFlagUsage)
}

// configTimezone returns settings.timezone from the config file, or "" when
// there is no config or no such setting. Config trouble is not an error
// here: env-authenticated runs may have no config file at all.
func configTimezone() string {
	cfg, err := config.Load()
	if err != nil || cfg == nil || cfg.Settings == nil {
		return ""
	}
	return cfg.Settings.Timezone
}

// resolveTimezone picks the zone name for local-time parsing: the
// --timezone flag, else settings.timezone, else "" (which localtime
// resolves to $TZ, then the system zone).
func resolveTimezone(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return configTimezone()
}

// resolveLocation is resolveTimezone followed by localtime.Location.
func resolveLocation(flagValue string) (*time.Location, error) {
	return localtime.Location(resolveTimezone(flagValue))
}
