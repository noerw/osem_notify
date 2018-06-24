package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"

	"../core"
)

/**
 * shared functionality between watch and check
 */

// TODO: actually to be read from arg / file

var defaultConf = &core.NotifyConfig{
	Notifications: core.TransportConfig{
		Transport: "email",
		Options: core.EmailNotifier{
			[]string{"test@nroo.de"},
			"notify@nroo.de",
		},
	},
	Events: []core.NotifyEvent{
		core.NotifyEvent{
			Type:      "measurement_age",
			Target:    "all",
			Threshold: "15m",
		},
		core.NotifyEvent{
			Type:   "measurement_suspicious",
			Target: "all",
		},
	},
}

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

func checkAndNotify(boxIds []string, defaultNotifyConf *core.NotifyConfig) error {
	results, err := core.CheckBoxes(boxIds, defaultConf)
	if err != nil {
		return err
	}

	results.Log()

	if shouldNotify {
		return results.SendNotifications()
	}
	return nil
}
