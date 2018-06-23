package core

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	// log.SetFormatter(&log.JSONFormatter{})
}

func CheckNotifications(boxIds []string, defaultConf *NotifyConfig) ([]Notification, []error) {
	log.Info("Checking notifications for ", len(boxIds), " box(es)")

	// TODO: return a map of Box: []Notification instead?
	notifications := []Notification{}
	errors := []error{}
	for _, boxId := range boxIds {
		n, err := checkBox(boxId, defaultConf)
		if notifications != nil {
			notifications = append(notifications, n...)
		}
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == 0 {
		errors = nil
	}

	return notifications, errors
}

func checkBox(boxId string, defaultConf *NotifyConfig) ([]Notification, error) {
	boxLogger := log.WithFields(log.Fields{"boxId": boxId})
	boxLogger.Debug("checking box for due notifications")

	// get box data
	box, err := osem.GetBox(boxId)
	if err != nil {
		boxLogger.Error(err)
		return nil, err
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
		return notifications, err2
	}
	if notifications == nil {
		boxLogger.Debug("all is fine")
		return nil, nil
	}

	// store notifications for later submit
	// notifier, err3 := box.getNotifier()
	// if err3 != nil {
	// 	boxLogger.Error("could not get notifier for box: ", err)
	// 	return notifications, err3
	// }
	// notifier.AddNotifications(notifications)

	return notifications, nil
}
