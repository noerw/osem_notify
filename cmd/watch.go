package cmd

import (
	"time"

	"github.com/spf13/cobra"
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
}

var watchBoxesCmd = &cobra.Command{
	Use:   "boxes <boxId> [...<boxIds>]",
	Short: "watch a list of box IDs for events",
	Long:  "specify box IDs to watch them for events",
	Args:  BoxIdValidator,
	PreRun: func(cmd *cobra.Command, args []string) {
		ticker = time.NewTicker(time.Duration(watchInterval) * time.Second).C
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		err := checkAndNotify(args)
		if err != nil {
			return err
		}
		for {
			<-ticker
			err = checkAndNotify(args)
			if err != nil {
				return err
			}
		}
	},
}
