package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/noerw/osem_notify/core"
	"github.com/noerw/osem_notify/utils"
)

func init() {
	debugCmd.AddCommand(debugNotificationsCmd)
	rootCmd.AddCommand(debugCmd)
}

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Run some debugging checks on osem_notify itself",
	Long:  "osem_notify debug <feature> tests the functionality of the given feature",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		utils.PrintConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var debugNotificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "Verify that notifications are working",
	Long:  "osem_notify debug <feature> tests the functionality of the given feature",
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultNotifyConf := &core.NotifyConfig{}
		err := viper.UnmarshalKey("defaultHealthchecks", defaultNotifyConf)
		if err != nil {
			return err
		}

		for transport, notifier := range core.Notifiers {
			notLog := log.WithField("transport", transport)
			opts := defaultNotifyConf.Notifications.Options
			notLog.Infof("testing notifer %s with options %v", transport, opts)
			n, err := notifier.New(opts)
			if err != nil {
				notLog.Warnf("could not initialize %s notifier. configuration might be missing?", transport)
				continue
			}

			host, _ := os.Hostname()
			err = n.Submit(core.Notification{
				Subject: "Test notification from opeSenseMap notifier",
				Body:    fmt.Sprintf("Your notification set up on %s is working fine!", host),
			})
			if err != nil {
				notLog.Warnf("could not submit test notification for %s notifier!", transport)
				continue
			}
			notLog.Info("Test notification (successfully?) submitted, check the specified inbox")
		}

		return nil
	},
}
