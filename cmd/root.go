package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "osem_notify",
	Long: "Run healthchecks and send notifications for boxes on opensensemap.org",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetOutput(os.Stdout)
		switch logLevel {
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		}
		switch logFormat {
		case "json":
			log.SetFormatter(&log.JSONFormatter{})
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var (
	shouldNotify  bool
	defaultConfig string
	logLevel      string
	logFormat     string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "", "info", "log level, can be one of debug, info, warn, error")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "log-format", "", "plain", "log format, can be plain or json")
	rootCmd.PersistentFlags().BoolVarP(&shouldNotify, "notify", "n", false, "if set, will send out notifications.\nOtherwise results are printed to stdout only")
	rootCmd.PersistentFlags().StringVarP(&defaultConfig, "conf-default", "c", "", "default JSON config to use for event checking")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
