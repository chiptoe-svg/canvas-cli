package commands

import (
	"github.com/njayp/ophis"
)

func init() {
	rootCmd.AddCommand(ophis.Command(&ophis.Config{
		ToolNamePrefix: "canvas",
		Selectors: []ophis.Selector{
			{
				// Exclude sensitive inherited flags from MCP exposure
				// show-token and config are PersistentFlags (inherited)
				InheritedFlagSelector: ophis.ExcludeFlags("show-token", "config"),
			},
		},
	}))
}
