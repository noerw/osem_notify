package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	CheckOk        = "OK"
	CheckErr       = "FAILED"
	eventTargetAll = "all" // if event.Target is this value, all sensors will be checked
)

type checkType struct {
	name      string                          // name that is used in config
	toString  func(result CheckResult) string // error message when check failed
	checkFunc func(event NotifyEvent, sensor Sensor, context Box) (CheckResult, error)
}

var checkers = map[string]checkType{
	checkMeasurementAge.name:    checkMeasurementAge,
	checkMeasurementMin.name:    checkMeasurementMin,
	checkMeasurementMax.name:    checkMeasurementMax,
	checkMeasurementFaulty.name: checkMeasurementFaulty,
}

type CheckResult struct {
	Status     string // should be CheckOk | CheckErr
	TargetName string
	Value      string
	Target     string

	Event     string // these should be copied from the NotifyEvent
	Threshold string
}

func (r CheckResult) HasStatus(statusToCheck []string) bool {
	for _, status := range statusToCheck {
		if status == r.Status {
			return true
		}
	}
	return false
}

func (r CheckResult) EventID() string {
	s := fmt.Sprintf("%s%s%s", r.Event, r.Target, r.Threshold)
	hasher := sha256.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (r CheckResult) String() string {
	if r.Status == CheckOk {
		return fmt.Sprintf("%s: %s (on sensor %s (%s) with value %s)\n", r.Status, r.Event, r.TargetName, r.Target, r.Value)
	} else {
		return fmt.Sprintf("%s: %s\n", r.Status, checkers[r.Event].toString(r))
	}
}

func (box Box) RunChecks() ([]CheckResult, error) {
	var results = []CheckResult{}
	boxLogger := log.WithField("box", box.Id)

	for _, event := range box.NotifyConf.Events {
		for _, s := range box.Sensors {
			// if a sensor never measured anything, thats ok. checks would fail anyway
			if s.LastMeasurement == nil {
				continue
			}

			if event.Target != s.Id && event.Target != eventTargetAll {
				continue
			}

			checker := checkers[event.Type]
			if checker.checkFunc == nil {
				boxLogger.Warnf("ignoring unknown event type %s", event.Type)
				continue
			}

			result, err := checker.checkFunc(event, s, box)
			if err != nil {
				boxLogger.Errorf("error checking event %s: %v", event.Type, err)
			}

			results = append(results, result)
		}
	}

	return results, nil
}
