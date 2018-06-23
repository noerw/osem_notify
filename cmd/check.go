package cmd

import (
	"../core"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	checkCmd.AddCommand(checkBoxCmd)
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "One-off check for events on boxes",
	Long:  "One-off check for events on boxes",
}

var checkBoxCmd = &cobra.Command{
	Use:   "boxes <boxId> [...<boxIds>]",
	Short: "one-off check on one or more box IDs",
	Long:  "specify box IDs to check them for events",
	Args:  BoxIdValidator,
	RunE: func(cmd *cobra.Command, args []string) error {
		notifications, err := core.CheckNotifications(args, defaultConf)
		if err != nil {
			return fmt.Errorf("error checking for notifications:", err)
		}
		fmt.Println(notifications)

		// logNotifications(notifications)
		if shouldNotify {
			// TODO
		}

		return nil
	},
}
