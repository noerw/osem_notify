package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"../core"
)

var watchInterval int
var ticker <-chan time.Time

func init() {
	watchCmd.AddCommand(watchBoxesCmd)
	watchCmd.PersistentFlags().IntVarP(&watchInterval, "interval", "i", 15, "interval to run checks in minutes")
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:     "watch",
	Aliases: []string{"serve"},
	Short:   "Watch boxes for events at an interval",
	Long:    "Watch boxes for events at an interval",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ticker = time.NewTicker(time.Duration(watchInterval) * time.Second).C
	},
}

var watchBoxesCmd = &cobra.Command{
	Use:   "boxes <boxId> [...<boxIds>]",
	Short: "watch a list of box IDs for events",
	Long:  "specify box IDs to watch them for events",
	Args:  BoxIdValidator,
	RunE: func(cmd *cobra.Command, args []string) error {
		exec := func() error {
			notifications, err := core.CheckNotifications(args, defaultConf)
			if err != nil {
				return fmt.Errorf("error checking for notifications: ", err)
			}
			fmt.Println(notifications)

			// logNotifications(notifications)
			if shouldNotify {
				// TODO
			}
			return nil
		}

		err := exec()
		if err != nil {
			return err
		}
		for {
			<-ticker
			err = exec()
			if err != nil {
				return err
			}
		}

		return nil
	},
}
