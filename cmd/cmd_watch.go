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

	watchAllCmd.PersistentFlags().StringVarP(&date, "date", "", "", "filter boxes by date")
	watchAllCmd.PersistentFlags().StringVarP(&exposure, "exposure", "", "", "filter boxes by exposure")
	watchAllCmd.PersistentFlags().StringVarP(&grouptag, "grouptag", "", "", "filter boxes by grouptag")
	watchAllCmd.PersistentFlags().StringVarP(&model, "model", "", "", "filter boxes by model")
	watchAllCmd.PersistentFlags().StringVarP(&phenomenon, "phenomenon", "", "", "filter boxes by phenomenon")
	watchCmd.PersistentFlags().IntVarP(&watchInterval, "interval", "i", 30, "interval to run checks in minutes")
	viper.BindPFlags(watchCmd.PersistentFlags())

	watchCmd.AddCommand(watchBoxesCmd)
	watchCmd.AddCommand(watchAllCmd)
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

var watchAllCmd = &cobra.Command{
	Use:   "all",
	Short: "watch all boxes registered on the map",
	PreRun: func(cmd *cobra.Command, args []string) {
		interval := viper.GetDuration("interval") * time.Minute
		ticker = time.NewTicker(interval).C
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		filters := parseBoxFilters()
		err := checkAndNotifyAll(filters)
		if err != nil {
			return err
		}
		for {
			<-ticker
			err = checkAndNotifyAll(filters)
			if err != nil {
				// we already did retries, so exiting seems appropriate
				return err
			}
		}
	},
}
