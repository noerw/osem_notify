package cmd

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/noerw/osem_notify/core"
	"github.com/noerw/osem_notify/utils"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	theConfig := cfgFile
	if cfgFile == "" {
		theConfig = utils.GetConfigFile("osem_notify")
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(theConfig)
	viper.SetEnvPrefix("OSEM_NOTIFY") // keys only work in upper case
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv() // WARNING: OSEM_NOTIFIY_CONFIG will not be considered this way. but why should it..

	// If a config file is found, read it in.
	if _, err := os.Stat(theConfig); err == nil {
		err := viper.ReadInConfig()
		if err != nil {
			log.Error("Error when reading config file:", err)
		}
	} else if cfgFile != "" {
		log.Error("Specified config file not found!")
		os.Exit(1)
	}

	validateConfig()
}

func validateConfig() {
	transport := viper.GetString("defaultHealthchecks.notifications.transport")
	if viper.GetBool("notify") && transport == "email" {
		if len(viper.GetStringSlice("defaultHealthchecks.notifications.options.recipients")) == 0 {
			log.Warn("No recipients set up for transport email")
		}

		emailRequired := []string{
			viper.GetString("email.host"),
			viper.GetString("email.port"),
			viper.GetString("email.user"),
			viper.GetString("email.pass"),
			viper.GetString("email.from"),
		}
		for _, conf := range emailRequired {
			if conf == "" {
				log.Error("Default transport set as email, but missing email config")
				os.Exit(1)
			}
		}
	}
}

func getNotifyConf(boxID string) (*core.NotifyConfig, error) {
	// config used when no configuration is present at all
	conf := &core.NotifyConfig{
		Events: []core.NotifyEvent{
			core.NotifyEvent{
				Type:      "measurement_age",
				Target:    "all",
				Threshold: "15m",
			},
			core.NotifyEvent{
				Type:      "measurement_faulty",
				Target:    "all",
				Threshold: "",
			},
		},
	}

	// override with default configuration from file
	// considering the case that .events may be defined but empty
	// to allow to define no events, and don't leak shorter lists into
	// previous longer ones
	if keyDefined("healthchecks.default.events") {
		conf.Events = []core.NotifyEvent{}
	}
	err := viper.UnmarshalKey("healthchecks.default", conf)
	if err != nil {
		return nil, err
	}

	// override with per box configuration from file
	if keyDefined("healthchecks." + boxID + ".events") {
		conf.Events = []core.NotifyEvent{}
	}
	err = viper.UnmarshalKey("healthchecks."+boxID, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

// implement our own keyCheck, as viper.InConfig() does not work
func keyDefined(key string) bool {
	allConfKeys := viper.AllKeys()
	for _, k := range allConfKeys {
		if k == key {
			return true
		}
	}
	return false
}
