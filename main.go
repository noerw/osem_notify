package main

import (
	"net/http"
	"os"
	"runtime"

	"github.com/carlescere/scheduler"
	log "github.com/sirupsen/logrus"
)

func checkBox(boxId string, defaultConf *NotifyConfig) {
	boxLogger := log.WithFields(log.Fields{"boxId": boxId})
	boxLogger.Debug("checking box for due notifications")

	// get box data
	box, err := osem.GetBox(boxId)
	if err != nil {
		boxLogger.Error("could not fetch box: ", err)
		return
	}

	// if box has no notify config, we use the defaultConf
	if box.NotifyConf == nil {
		box.NotifyConf = defaultConf
	}
	boxLogger.Debug(box.NotifyConf)

	// run checks
	notifications, err2 := box.runChecks()
	if err2 != nil {
		boxLogger.Error("could not run checks on box: ", err)
		return
	}
	if notifications == nil {
		boxLogger.Debug("all is fine")
		return
	}

	// store notifications for later submit
	// notifier, err3 := box.getNotifier()
	// if err3 != nil {
	// 	boxLogger.Error("could not get notifier for box: ", err)
	// 	return
	// }
	// notifier.AddNotifications(notifications)
}

func checkNotifications() {
	log.Info("running job checkNotifications()")
	checkBox("593bcd656ccf3b0011791f5a", defaultConf)
}

var osem = NewOsemClient(&http.Client{})
var defaultConf = &NotifyConfig{
	// Transports: struct {
	// 	Slack: SlackConfig{
	// 		Channel: "asdf"
	// 		Token: "qwer"
	// 	}
	// },
	Events: []NotifyEvent{
		NotifyEvent{
			Type:      "measurementAge",
			Target:    "593bcd656ccf3b0011791f5d",
			Threshold: "5h",
		},
	},
}

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	// log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	scheduler.Every(15).Seconds().Run(checkNotifications)
	// scheduler.Every(30).Seconds().Run(submitNotifications)
	runtime.Goexit() // keep runtime running
}
