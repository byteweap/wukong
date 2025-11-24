package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/byteweap/wukong/cmd/wukong/internal/group"
)

const release = "v1.0.0"

var root = &cobra.Command{
	Use:     "wukong",
	Short:   "WuKong: An elegant toolkit For wukong game framework.",
	Long:    "WuKong: An elegant toolkit For wukong game framework.",
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
