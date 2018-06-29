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

var (
	clearCache bool
)

func init() {
	debugCmd.AddCommand(debugNotificationsCmd)
	debugCacheCmd.PersistentFlags().BoolVarP(&clearCache, "clear", "", false, "reset the notifications cache")
	debugCmd.AddCommand(debugCacheCmd)
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

var debugCacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Print or clear the notifications cache",
	Long:  "osem_notify debug cache prints the contents of the notifications cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		if clearCache {
			return core.ClearCache()
		}
		core.PrintCache()
		return nil
	},
}


var debugNotificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "Verify that notifications are working",
	Long:  `osem_notify debug notifications sends a test notification according
to healthchecks.default.notifications.options as defined in the config file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultNotifyConf := &core.NotifyConfig{}
		err := viper.UnmarshalKey("healthchecks.default", defaultNotifyConf)
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
				Subject: "Test notification from openSenseMap notifier",
				Body:    fmt.Sprintf("Your notification set up on %s is working fine!", host),
			})
			if err != nil {
				notLog.Warnf("could not submit test notification for %s notifier: %s", transport, err)
				continue
			}
			notLog.Info("Test notification (successfully?) submitted, check the specified inbox")
		}

		return nil
	},
}
