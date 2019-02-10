package cmd

import (
	"fmt"
	"regexp"
	"strings"

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
	boxLocalConfig := map[string]*core.NotifyConfig{}
	for _, boxID := range boxIds {
		c, err := getNotifyConf(boxID)
		if err != nil {
			return err
		}
		boxLocalConfig[boxID] = c
	}

	osem := core.NewOsemClient(viper.GetString("api"))
	results, err := core.CheckBoxes(boxLocalConfig, osem)
	if err != nil {
		return err
	}

	results.Log()

	notify := strings.ToLower(viper.GetString("notify"))
	if notify != "" {
		types := []string{}
		switch notify {
		case "all":
			types = []string{core.CheckErr, core.CheckOk}
		case "error", "err":
			types = []string{core.CheckErr}
		case "ok":
			types = []string{core.CheckOk}
		default:
			return fmt.Errorf("invalid value %s for \"notify\"", notify)
		}

		useCache := !viper.GetBool("no-cache")
		return results.SendNotifications(types, useCache)
	}
	return nil
}
