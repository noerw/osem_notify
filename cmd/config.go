package cmd

import (
	"os"
	"strings"

	"github.com/noerw/osem_notify/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
