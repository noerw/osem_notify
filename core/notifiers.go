package core

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

func (results BoxCheckResults) SendNotifications() error {
	results = results.FilterChangedFromCache(false)
	errs := []string{}

	n := results.Size()
	if n == 0 {
		log.Info("No notifications due.")
		return nil
	} else {
		log.Infof("Notifying for %v checks turned bad in total...", results.Size())
	}

	// FIXME: only update cache when notifications sent successfully
	for box, resultsDue := range results {
		if len(resultsDue) == 0 {
			continue
		}

		transport := box.NotifyConf.Notifications.Transport
		notifyLog := log.WithFields(log.Fields{
			"boxId":     box.Id,
			"transport": transport,
		})

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
			notifyLog.Debugf("trying to submit (retry %v)", retry)
		}
		if submitErr != nil {
			notifyLog.Error(submitErr)
			errs = append(errs, submitErr.Error())
			continue
		}

		notifyLog.Infof("Sent notification for %s via %s with %v new issues", box.Name, transport, len(resultsDue))
	}

	if len(errs) != 0 {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}
	return nil
}

func (results BoxCheckResults) FilterChangedFromCache(keepOk bool) BoxCheckResults {
	remaining := BoxCheckResults{}

	for box, boxResults := range results {
		// get results from cache. they are indexed by an event ID per boxId
		// filter, so that only changed result.Status remain
		remaining[box] = []CheckResult{}
		for _, result := range boxResults {
			cached := viper.GetStringMap(fmt.Sprintf("watchcache.%s.%s", box.Id, result.EventID()))
			if result.Status != cached["laststatus"] {
				if result.Status != CheckOk || keepOk {
					remaining[box] = append(remaining[box], result)
				}
			}
		}

		// TODO: reminder functionality: extract additional results with Status ERR
		// from cache with time.Since(lastNotifyDate) > remindAfter.
		// would require to serialize the full result..
	}

	// upate cache, setting lastNotifyDate to Now()
	for box, boxResults := range results {
		for _, result := range boxResults {
			// FIXME: somehow this is not persisted?
			key := fmt.Sprintf("watchcache.%s.%s", box.Id, result.EventID())
			viper.Set(key+".laststatus", result.Status)
		}
	}

	return remaining
}
