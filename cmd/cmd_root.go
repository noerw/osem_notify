package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/noerw/osem_notify/utils"
)

var configHelpCmd = &cobra.Command{
	Use:   "config",
	Short: "How to configure osem_notify",
	Long: `osem_notify works out of the box for basic functionality, but uses configuration to
  set up notification transports and healthchecks. Additionally, all command line flags can
  be set to default values through the configuration.

  Configuration can be set either through a YAML file, or through environment variables.
  You can use different configuration files per call by settings the --config flag.


> Example configuration:

	healthchecks:
		# override default health checks for all boxes
		default:
			notifications:
				transport: email
				options:
					recipients:
					- fridolina@example.com
			events:
				- type: "measurement_age"
					target: "all"    # all sensors
					threshold: "15m" # any duration
				- type: "measurement_faulty"
					target: "all"
					threshold: ""

		# override default health checks per box
		593bcd656ccf3b0011791f5a:
			notifications:
				options:
					recipients:
					- ruth.less@example.com
			events:
				- type: "measurement_max"
					target: "593bcd656ccf3b0011791f5b"
					threshold: "40"

  # only needed when sending notifications via email
  email:
    host: smtp.example.com
    port: 25
    user: foo
    pass: bar
    from: hildegunst@example.com


> possible values for healthchecks.*.notifications:

  transport | options
  ----------|-------------------------------------
  email     | recipients: list of email addresses


> possible values for healthchecks.*.events[]:

  type               | description
  -------------------|---------------------------------------------------
  measurement_age    | Alert when sensor target has not submitted measurements within threshold duration.
  measurement_faulty | Alert when sensor target's last reading was a presumably faulty value (e.g. broken / disconnected sensor).
  measurement_min    | Alert when sensor target's last measurement is lower than threshold.
  measurement_max    | Alert when sensor target's last measurement is higher than threshold.

  - target can be either a sensor ID, or "all" to match all sensors of the box.
  - threshold must be a string.

> configuration via environment variables

  Instead of a YAML file, you may configure the tool through environment variables. Keys are the same as in the YAML, but:
  keys are prefixed with "OSEM_NOTIFY_", path separator is not ".", but "_", all upper case

  Example: OSEM_NOTIFY_EMAIL_PASS=supersecret osem_notify check boxes`,
}

var rootCmd = &cobra.Command{
	Use:   "osem_notify",
	Short: "Root command displaying help",
	Long:  "Run healthchecks and send notifications for boxes on opensensemap.org",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// set up logger
		log.SetOutput(os.Stdout)
		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
			utils.PrintConfig()
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
		debug        bool
		noCache      bool
		shouldNotify string
		logFormat    string
		api          string
	)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file (default $HOME/.osem_notify.yml)")
	rootCmd.PersistentFlags().StringVarP(&api, "api", "a", "https://api.opensensemap.org", "openSenseMap API to query against")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "logformat", "l", "plain", "log format, can be plain or json")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable verbose logging")
	rootCmd.PersistentFlags().StringVarP(&shouldNotify, "notify", "n", "", `if set, will send out notifications for the specified type of check result,
Otherwise results are printed to stdout only.
Allowed values are "all", "error", "ok".
You might want to run 'osem_notify debug notifications' first to verify everything works.

Notifications for failing checks are sent only once,
and then cached until the issue got resolved.
To clear the cache, delete the file ~/.osem_notify_cache.yaml.
`)
	rootCmd.PersistentFlags().BoolVarP(&noCache, "no-cache", "", false, "send all notifications, ignoring results from previous runs. also don't update the cache.")

	viper.BindPFlags(rootCmd.PersistentFlags()) // let flags override config

	rootCmd.AddCommand(configHelpCmd)
}

func Execute() {
	// generate documentation
	// err := doc.GenMarkdownTree(rootCmd, "./docs")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
