package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "osem_notify",
	Long: "Run healthchecks and send notifications for boxes on opensensemap.org",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var shouldNotify bool
var defaultConfig string

func init() {
	rootCmd.PersistentFlags().BoolVarP(&shouldNotify, "notify", "n", false, "if set, will send out notifications.\nOtherwise results are printed to stdout only")
	rootCmd.PersistentFlags().StringVarP(&defaultConfig, "confdefault", "c", "", "default JSON config to use for event checking")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
	}
}
