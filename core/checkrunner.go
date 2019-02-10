package core

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type BoxCheckResults map[*Box][]CheckResult

func (results BoxCheckResults) Size(statusToCheck []string) int {
	size := 0
	for _, boxResults := range results {
		for _, result := range boxResults {
			if result.HasStatus(statusToCheck) {
				size++
			}
		}
	}
	return size
}

func (results BoxCheckResults) Log() {
	// collect statistics for summary print
	boxesSkipped := 0
	boxesWithIssues := 0
	boxesWithoutIssues := 0
	failedChecks := 0
	errorsByEvent := map[string]int{}
	for event, _ := range checkers {
		errorsByEvent[event] = 0
	}

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
				errorsByEvent[r.Event]++
			}
		}

		if len(boxResults) == 0 {
			boxLog.Infof("%s: no checks defined", box.Name)
			boxesSkipped++
		} else if countErr == 0 {
			boxLog.Infof("%s: all is fine!", box.Name)
			boxesWithoutIssues++
		} else {
			// we logged the error(s) already
			boxesWithIssues++
			failedChecks += countErr
		}
	}

	// print summary
	boxesChecked := boxesWithIssues + boxesWithoutIssues
	if boxesChecked > 1 {
		summaryLog := log.WithFields(log.Fields{
			"boxesChecked":  boxesChecked,
			"boxesSkipped":  boxesSkipped, // boxes are also skipped when they never submitted any measurements before!
			"boxesOk":       boxesWithoutIssues,
			"boxesErr":      boxesWithIssues,
			"failedChecks":  failedChecks,
			"errorsByEvent": errorsByEvent,
		})
		summaryLog.Infof(
			"check summary: %v of %v checked boxes are fine (%v had no checks)!",
			boxesWithoutIssues,
			boxesChecked,
			boxesSkipped)
	}
}

func CheckBoxes(boxLocalConfs map[string]*NotifyConfig, osem *OsemClient) (BoxCheckResults, error) {
	log.Info("Checking notifications for ", len(boxLocalConfs), " box(es)")

	results := BoxCheckResults{}
	errs := []string{}

	// @TODO: check boxes in parallel, capped at 5 at once. and/or rate limit?
	for boxId, localConf := range boxLocalConfs {
		boxLogger := log.WithField("boxId", boxId)
		boxLogger.Debug("checking box for events")

		box, res, err := checkBox(boxId, localConf, osem)
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

func checkBox(boxId string, defaultConf *NotifyConfig, osem *OsemClient) (*Box, []CheckResult, error) {
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
