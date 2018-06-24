package cmd

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:  "osem_notify",
	Long: "Run healthchecks and send notifications for boxes on opensensemap.org",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// set up config environment FIXME: cannot open / write file?!
		viper.SetConfigType("json")
		viper.SetConfigFile(".osem_notify")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
		// // If a config file is found, read it in.
		// if _, err := os.Stat(path.Join(os.Getenv("HOME"), ".osem_notify.yml")); err == nil {
		// 	err := viper.ReadInConfig()
		// 	if err != nil {
		// 		fmt.Println("Error when reading config file:", err)
		// 	}
		// }
		viper.SetEnvPrefix("osem_notify")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
		viper.AutomaticEnv()

		// set up logger
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

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var (
	shouldNotify bool
	logLevel     string
	logFormat    string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "", "info", "log level, can be one of debug, info, warn, error")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "log-format", "", "plain", "log format, can be plain or json")
	rootCmd.PersistentFlags().BoolVarP(&shouldNotify, "notify", "n", false, "if set, will send out notifications.\nOtherwise results are printed to stdout only")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
