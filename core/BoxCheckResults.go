package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CheckResult struct {
	Status    string
	Event     string
	Target    string
	Value     string
	Threshold string
}

func (r CheckResult) EventID() string {
	s := fmt.Sprintf("%s%s%s", r.Event, r.Target, r.Threshold)
	hasher := sha256.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (r CheckResult) String() string {
	if r.Status == CheckOk {
		return fmt.Sprintf("%s %s (on sensor %s with value %s)\n", r.Event, r.Status, r.Target, r.Value)
	} else {
		return fmt.Sprintf("%s: "+checkTypes[r.Event].description+"\n", r.Status, r.Target, r.Value)
	}
}

type BoxCheckResults map[*Box][]CheckResult

func (results BoxCheckResults) Size() int {
	size := 0
	for _, boxResults := range results {
		size += len(boxResults)
	}
	return size
}

func (results BoxCheckResults) Log() {
	for box, boxResults := range results {
		boxLog := log.WithFields(log.Fields{
			"boxId": box.Id,
		})
		countErr := 0
		for _, r := range boxResults {
			resultLog := boxLog.WithFields(log.Fields{
				"status": r.Status,
				"event":  r.Event,
				"value":  r.Value,
				"target": r.Target,
			})
			if r.Status == CheckOk {
				resultLog.Debugf("%s: %s", box.Name, r)
			} else {
				resultLog.Warnf("%s: %s", box.Name, r)
				countErr++
			}
		}
		if countErr == 0 {
			boxLog.Infof("%s: all is fine!", box.Name)
		}
	}
}

func (results BoxCheckResults) SendNotifications() error {
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

func CheckBoxes(boxIds []string, defaultConf *NotifyConfig) (BoxCheckResults, error) {
	log.Debug("Checking notifications for ", len(boxIds), " box(es)")

	results := BoxCheckResults{}
	for _, boxId := range boxIds {
		box, res, err := checkBox(boxId, defaultConf)
		if err != nil {
			return nil, err
		}
		results[box] = res
	}

	return results, nil
}

func checkBox(boxId string, defaultConf *NotifyConfig) (*Box, []CheckResult, error) {
	boxLogger := log.WithFields(log.Fields{"boxId": boxId})
	boxLogger.Info("checking box for events")

	// get box data
	box, err := Osem.GetBox(boxId)
	if err != nil {
		boxLogger.Error(err)
		return nil, nil, err
	}

	// if box has no notify config, we use the defaultConf
	if box.NotifyConf == nil {
		box.NotifyConf = defaultConf
	}

	// run checks
	results, err2 := box.RunChecks()
	if err2 != nil {
		boxLogger.Error("could not run checks on box: ", err2)
		return box, results, err2
	}

	return box, results, nil
}
