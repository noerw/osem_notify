package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "osem_notify",
	Short: "Root command displaying help",
	Long:  "Run healthchecks and send notifications for boxes on opensensemap.org",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// set up logger
		log.SetOutput(os.Stdout)
		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
			printConfig()
		} else {
			log.SetLevel(log.InfoLevel)
		}
		switch viper.Get("logformat") {
		case "json":
			log.SetFormatter(&log.JSONFormatter{})
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// accessed in initConfig(), as it is initialized before config is loaded (sic)
var cfgFile string

func init() {
	var (
		shouldNotify bool
		debug        bool
		logFormat    string
	)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file (default $HOME/.osem_notify.yml)")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "logformat", "l", "plain", "log format, can be plain or json")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&shouldNotify, "notify", "n", false, "if set, will send out notifications.\nOtherwise results are printed to stdout only")

	viper.BindPFlags(rootCmd.PersistentFlags()) // let flags override config
}

func Execute() {
	// generate documentation
	// err := doc.GenMarkdownTree(rootCmd, "./doc")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
