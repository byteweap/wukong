package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/byteweap/meta/cmd/meta/internal/group"
)

const release = "v1.0.0"

var root = &cobra.Command{
	Use:     "meta",
	Short:   "meta: An elegant toolkit For meta game framework.",
	Long:    "meta: An elegant toolkit For meta game framework.",
	Version: release,
}

func init() {
	root.AddCommand(group.CmdNew)
	root.AddCommand(group.CmdRun)
}

func main() {
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
