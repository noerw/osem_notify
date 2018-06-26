package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ticker <-chan time.Time

func init() {
	var (
		watchInterval int
	)

	watchCmd.PersistentFlags().IntVarP(&watchInterval, "interval", "i", 30, "interval to run checks in minutes")
	viper.BindPFlags(watchCmd.PersistentFlags())

	watchCmd.AddCommand(watchBoxesCmd)
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:     "watch",
	Aliases: []string{"serve"},
	Short:   "Watch boxes for events at an interval",
	Long:    "Watch boxes for events at an interval",
}

var watchBoxesCmd = &cobra.Command{
	Use:     "boxes <boxId> [...<boxIds>]",
	Aliases: []string{"box"},
	Short:   "watch a list of box IDs for events",
	Long:    "specify box IDs to watch them for events",
	Args:    BoxIdValidator,
	PreRun: func(cmd *cobra.Command, args []string) {
		interval := viper.GetDuration("interval") * time.Minute
		ticker = time.NewTicker(interval).C
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
				// we already did retries, so exiting seems appropriate
				return err
			}
		}
	},
}
