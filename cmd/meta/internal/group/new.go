package group

import "github.com/spf13/cobra"

// CmdNew TODO
var CmdNew = &cobra.Command{
	Use:     "new",
	Short:   "create a service project by the default template.",
	Example: "wukong new gate",
	Run:     coreNew,
}

func coreNew(c *cobra.Command, args []string) {

}
