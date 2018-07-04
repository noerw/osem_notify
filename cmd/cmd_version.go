package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TODO: insert automatically via build step?
const (
	VERSION   = "1.1.1"
	BUILDDATE = "2018-07-04T19:42:00+02"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get build and version information",
	Long:  "osem_notify version returns its build and version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s", VERSION, BUILDDATE)
	},
}
