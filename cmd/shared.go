package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/noerw/osem_notify/core"
)

/**
 * shared functionality between watch and check
 */

func isValidBoxId(boxId string) bool {
	// boxIds are UUIDs
	r := regexp.MustCompile("^[0-9a-fA-F]{24}$")
	return r.MatchString(boxId)
}

func BoxIdValidator(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("requires at least 1 argument")
	}
	for _, boxId := range args {
		if isValidBoxId(boxId) == false {
			return fmt.Errorf("invalid boxId specified: %s", boxId)
		}
	}
	return nil
}

func checkAndNotify(boxIds []string) error {
	defaultNotifyConf := &core.NotifyConfig{}
	err := viper.UnmarshalKey("defaultHealthchecks", defaultNotifyConf)
	if err != nil {
		return err
	}

	results, err := core.CheckBoxes(boxIds, defaultNotifyConf)
	if err != nil {
		return err
	}

	results.Log()

	if viper.GetBool("notify") {
		return results.SendNotifications()
	}
	return nil
}
