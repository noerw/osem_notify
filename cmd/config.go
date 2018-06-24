package cmd

import (
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

/**
 * config file handling, as it is kinda broken in spf13/viper
 * mostly copied from https://github.com/TheThingsNetwork/ttn/blob/f623a6a/ttnctl/util/config.go
 */

// GetConfigFile returns the location of the configuration file.
// It checks the following (in this order):
// the --config flag
// $XDG_CONFIG_HOME/osem_notify/config.yml (if $XDG_CONFIG_HOME is set)
// $HOME/.osem_notify.yml
func getConfigFile() string {
	flag := viper.GetString("config")

	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg != "" {
		xdg = path.Join(xdg, "osem_notify", "config.yml")
	}

	home := os.Getenv("HOME")
	homeyml := ""
	homeyaml := ""

	if home != "" {
		homeyml = path.Join(home, ".osem_notify.yml")
		homeyaml = path.Join(home, ".osem_notify.yaml")
	}

	try_files := []string{
		flag,
		xdg,
		homeyml,
		homeyaml,
	}

	// find a file that exists, and use that
	for _, file := range try_files {
		if file != "" {
			if _, err := os.Stat(file); err == nil {
				return file
			}
		}
	}

	// no file found, set up correct fallback
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return xdg
	} else {
		return homeyml
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	theConfig := cfgFile
	if cfgFile == "" {
		theConfig = getConfigFile()
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

func printConfig() {
	log.Debug("Using config:")
	printKV("config file", viper.ConfigFileUsed())
	for key, val := range viper.AllSettings() {
		printKV(key, val)
	}
}

func printKV(key, val interface{}) {
	log.Debugf("%20s: %v", key, val)
}
