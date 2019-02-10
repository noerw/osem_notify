package utils

import (
	"os"
	"path"
	"strconv"
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
func GetConfigFile(name string) string {
	flag := viper.GetString("config")

	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg != "" {
		xdg = path.Join(xdg, name, "config.yml")
	}

	home := os.Getenv("HOME")
	homeyml := ""
	homeyaml := ""

	if home != "" {
		homeyml = path.Join(home, "."+name+".yml")
		homeyaml = path.Join(home, "."+name+".yaml")
	}

	tryFiles := []string{
		flag,
		xdg,
		homeyml,
		homeyaml,
	}

	// find a file that exists, and use that
	for _, file := range tryFiles {
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

func PrintConfig() {
	log.Debug("Using config:")
	printKV("config file", viper.ConfigFileUsed())
	for key, val := range viper.AllSettings() {
		printKV(key, val)
	}
}

func printKV(key, val interface{}) {
	log.Debugf("%20s: %v", key, val)
}

func ParseFloat(val string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(val), 64)
}
