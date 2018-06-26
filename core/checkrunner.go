package core

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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

func CheckBoxes(boxIds []string, defaultConf *NotifyConfig) (BoxCheckResults, error) {
	log.Debug("Checking notifications for ", len(boxIds), " box(es)")

	results := BoxCheckResults{}
	errs := []string{}

	// TODO: check boxes in parallel, capped at 5 at once
	for _, boxId := range boxIds {
		boxLogger := log.WithField("boxId", boxId)
		boxLogger.Info("checking box for events")

		box, res, err := checkBox(boxId, defaultConf)
		if err != nil {
			boxLogger.Errorf("could not run checks on box %s: %s", boxId, err)
			errs = append(errs, err.Error())
			continue
		}
		results[box] = res
	}

	if len(errs) != 0 {
		return results, fmt.Errorf(strings.Join(errs, "\n"))
	}
	return results, nil
}

func checkBox(boxId string, defaultConf *NotifyConfig) (*Box, []CheckResult, error) {

	osem := NewOsemClient(viper.GetString("api"))

	// get box data
	box, err := osem.GetBox(boxId)
	if err != nil {
		return nil, nil, err
	}

	// if box has no notify config, we use the defaultConf
	if box.NotifyConf == nil {
		box.NotifyConf = defaultConf
	}

	// run checks
	results, err2 := box.RunChecks()
	if err2 != nil {
		return box, results, err2
	}

	return box, results, nil
}
