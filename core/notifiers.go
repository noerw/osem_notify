package core

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var Notifiers = map[string]AbstractNotifier{
	"email": EmailNotifier{},
}

type AbstractNotifier interface {
	New(config interface{}) (AbstractNotifier, error)
	ComposeNotification(box *Box, checks []CheckResult) Notification
	Submit(notification Notification) error
}

type Notification struct {
	Body    string
	Subject string
}

//////

func (box Box) GetNotifier() (AbstractNotifier, error) {
	transport := box.NotifyConf.Notifications.Transport
	if transport == "" {
		return nil, fmt.Errorf("No notification transport provided")
	}

	notifier := Notifiers[transport]
	if notifier == nil {
		return nil, fmt.Errorf("%s is not a supported notification transport", transport)
	}

	return notifier.New(box.NotifyConf.Notifications.Options)
}

func (results BoxCheckResults) SendNotifications(notifyTypes []string, useCache bool) error {
	if useCache {
		results = results.filterChangedFromCache()
	}

	toCheck := results.Size(notifyTypes)
	if toCheck == 0 {
		log.Info("No notifications due.")
	} else {
		log.Infof("Notifying for %v checks changing state to %v...", toCheck, notifyTypes)
	}

	errs := []string{}
	for box, resultsBox := range results {
		// only submit results which are errors
		resultsDue := []CheckResult{}
		for _, result := range resultsBox {
			if result.HasStatus(notifyTypes) {
				resultsDue = append(resultsDue, result)
			}
		}

		transport := box.NotifyConf.Notifications.Transport
		notifyLog := log.WithFields(log.Fields{
			"boxId":     box.Id,
			"transport": transport,
		})

		if len(resultsDue) != 0 {
			notifier, err := box.GetNotifier()
			if err != nil {
				notifyLog.Error(err)
				errs = append(errs, err.Error())
				continue
			}

			notification := notifier.ComposeNotification(box, resultsDue)

			var submitErr error
			submitErr = notifier.Submit(notification)
			for retry := 1; submitErr != nil && retry < 3; retry++ {
				time.Sleep(10 * time.Second)
				notifyLog.Infof("trying to submit (retry %v)", retry)
			}
			if submitErr != nil {
				notifyLog.Error(submitErr)
				errs = append(errs, submitErr.Error())
				continue
			}
		}

		// update cache (with /all/ changed results to reset status)
		if useCache {
			notifyLog.Debug("updating cache")
			cacheError := updateCache(box, resultsBox)
			if cacheError != nil {
				notifyLog.Error("could not cache notification results: ", cacheError)
				errs = append(errs, cacheError.Error())
			}
		}

		if len(resultsDue) != 0 {
			notifyLog.Infof("Sent notification for %s via %s with %v updated issues", box.Name, transport, len(resultsDue))
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}
	return nil
}
