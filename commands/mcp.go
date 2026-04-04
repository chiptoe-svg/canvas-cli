package commands

import (
	"github.com/njayp/ophis"
)

func init() {
	rootCmd.AddCommand(ophis.Command(&ophis.Config{
		ToolNamePrefix: "canvas",
		Selectors: []ophis.Selector{
			{
				// Exclude sensitive flags from all commands
				LocalFlagSelector:     ophis.ExcludeFlags("show-token"),
				InheritedFlagSelector: ophis.ExcludeFlags("config"),
			},
		},
	}))
}
