package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	checkAllCmd.PersistentFlags().StringVarP(&date, "date", "", "", "filter boxes by date AND phenomenon")
	checkAllCmd.PersistentFlags().StringVarP(&exposure, "exposure", "", "", "filter boxes by exposure")
	checkAllCmd.PersistentFlags().StringVarP(&grouptag, "grouptag", "", "", "filter boxes by grouptag")
	checkAllCmd.PersistentFlags().StringVarP(&model, "model", "", "", "filter boxes by model")
	checkAllCmd.PersistentFlags().StringVarP(&phenomenon, "phenomenon", "", "", "filter boxes by phenomenon AND date")
	checkCmd.AddCommand(checkBoxCmd)
	checkCmd.AddCommand(checkAllCmd)
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "One-off check for events on boxes",
	Long:  "One-off check for events on boxes",
}

var checkBoxCmd = &cobra.Command{
	Use:     "boxes <boxId> [...<boxIds>]",
	Aliases: []string{"box"},
	Short:   "one-off check on one or more box IDs",
	Long:    "specify box IDs to check them for events",
	Args:    BoxIdValidator,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return checkAndNotify(args)
	},
}

var checkAllCmd = &cobra.Command{
	Use:   "all",
	Short: "one-off check on all boxes registered on the opensensemap instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		// no flag validation, as the API already does a good job at that

		return checkAndNotifyAll(parseBoxFilters())
	},
}
