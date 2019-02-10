package core

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var Notifiers = map[string]AbstractNotifier{
	"email": EmailNotifier{},
	"slack": SlackNotifier{},
	"xmpp":  XmppNotifier{},
}

type AbstractNotifier interface {
	New(config TransportConfig) (AbstractNotifier, error)
	Submit(notification Notification) error
}

type Notification struct {
	Status  string // one of CheckOk | CheckErr
	Body    string
	Subject string
}

//////

func (box Box) GetNotifier() (AbstractNotifier, error) {
	return GetNotifier(&box.NotifyConf.Notifications)
}

func GetNotifier(config *TransportConfig) (AbstractNotifier, error) {
	transport := config.Transport

	if transport == "" {
		return nil, fmt.Errorf("No notification transport provided")
	}

	notifier := Notifiers[transport]
	if notifier == nil {
		return nil, fmt.Errorf("%s is not a supported notification transport", transport)
	}

	return notifier.New(*config)
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

			notification := ComposeNotification(box, resultsDue)

			var submitErr error
			submitErr = notifier.Submit(notification)
			for retry := 0; submitErr != nil && retry < 2; retry++ {
				notifyLog.Warnf("sending notification failed (retry %v): %v", retry, submitErr)
				time.Sleep(10 * time.Second)
				submitErr = notifier.Submit(notification)
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
			updateCache(box, resultsBox)
		}

		if len(resultsDue) != 0 {
			notifyLog.Infof("Sent notification for %s via %s with %v updated issues", box.Name, transport, len(resultsDue))
		}
	}

	// persist changes to cache
	if useCache {
		err := writeCache()
		if err != nil {
			log.Error("could not write cache of notification results: ", err)
			errs = append(errs, err.Error())
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}
	return nil
}

func ComposeNotification(box *Box, checks []CheckResult) Notification {
	errTexts := []string{}
	resolvedTexts := []string{}
	for _, check := range checks {
		if check.Status == CheckErr {
			errTexts = append(errTexts, check.String())
		} else {
			resolvedTexts = append(resolvedTexts, check.String())
		}
	}

	var (
		resolved     string
		resolvedList string
		errList      string
		status       string
	)
	if len(resolvedTexts) != 0 {
		resolvedList = fmt.Sprintf("Resolved issue(s):\n\n%s\n\n", strings.Join(resolvedTexts, "\n"))
	}
	if len(errTexts) != 0 {
		errList = fmt.Sprintf("New issue(s):\n\n%s\n\n", strings.Join(errTexts, "\n"))
		status = CheckErr
	} else {
		resolved = "resolved "
		status = CheckOk
	}

	return Notification{
		Status:  status,
		Subject: fmt.Sprintf("Issues %swith your box \"%s\" on opensensemap.org!", resolved, box.Name),
		Body: fmt.Sprintf("A check at %s identified the following updates for your box \"%s\":\n\n%s%sYou may visit https://opensensemap.org/explore/%s for more details.",
			time.Now().Round(time.Minute), box.Name, errList, resolvedList, box.Id),
	}
}
