package group

import "github.com/spf13/cobra"

// CmdRun TODO
var CmdRun = &cobra.Command{
	Use:     "run",
	Short:   "run a service.",
	Long:    "",
	Example: "",
	Run:     coreRun,
}

func coreRun(c *cobra.Command, args []string) {

}
