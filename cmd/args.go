package cmd

import (
	"fmt"
	"regexp"

	"../core"
	"github.com/spf13/cobra"
)

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

// TODO: actually to be read from arg / file
var defaultConf = &core.NotifyConfig{
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
