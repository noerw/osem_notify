package core

import (
	"fmt"

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
	// FIXME: don't return on errors, process all boxes first!
	// FIXME: only update cache when notifications sent successfully
	results = results.FilterChangedFromCache(false)

	n := results.Size()
	if n == 0 {
		log.Info("No notifications due.")
		return nil
	} else {
		log.Infof("Notifying for %v checks turned bad in total...", results.Size())
	}

	for box, resultsDue := range results {
		if len(resultsDue) == 0 {
			continue
		}

		transport := box.NotifyConf.Notifications.Transport
		notifyLog := log.WithFields(log.Fields{
			"boxId":     box.Id,
			"transport": transport,
		})

		notifier, err2 := box.GetNotifier()
		if err2 != nil {
			notifyLog.Error(err2)
			return err2
		}
		notification := notifier.ComposeNotification(box, resultsDue)
		err3 := notifier.Submit(notification)
		if err3 != nil {
			notifyLog.Error(err3)
			return err3
		}
		notifyLog.Infof("Sent notification for %s via %s with %v new issues", box.Name, transport, len(resultsDue))
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
