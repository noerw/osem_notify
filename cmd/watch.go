package cmd

import (
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
		cmd.SilenceUsage = true
		exec := func() error {
			results, err := CheckBoxes(args, defaultConf)
			if err != nil {
				return err
			}

			results, err = filterFromCache(results)
			if err != nil {
				return err
			}

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
	},
}

func filterFromCache(results []core.CheckResult) ([]core.CheckResult, error) {
	// get results from cache. they are indexed by ______

	// filter, so that only changed result.Status remain

	// extract additional results with Status ERR from cache with time.Since(lastNotifyDate) > thresh

	// upate cache set lastNotifyDate to Now()

	return results, nil
}
