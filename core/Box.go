package core

import (
	"fmt"
	"strconv"
	"time"
)

const (
	CheckOk                       = "OK"
	CheckErr                      = "ERROR"
	eventMeasurementAge           = "measurement_age"        // errors if age of last measurement is higher than a duration
	eventMeasurementValMin        = "measurement_min"        // errors if value of last measurement is lower than threshold
	eventMeasurementValMax        = "measurement_max"        // errors if value of last measurement is higher than threshold
	eventMeasurementValSuspicious = "measurement_suspicious" // checks value of last measurement against a blacklist of values
	eventTargetAll                = "all"                    // if event.Target is this value, all sensors will be checked
)

type SuspiciousValue struct {
	sensor string
	val    float64
}

var suspiciousVals = map[SuspiciousValue]bool{
	SuspiciousValue{sensor: "BMP280", val: 0.0}:  true,
	SuspiciousValue{sensor: "HDC1008", val: 0.0}: true,
	SuspiciousValue{sensor: "HDC1008", val: -40}: true,
	SuspiciousValue{sensor: "SDS 011", val: 0.0}: true,
}

type CheckResult struct {
	Status string
	Event  string
	Target string
	Value  string
}

func (r CheckResult) String() string {
	return fmt.Sprintf("check %s on sensor %s: %s with value %s\n", r.Event, r.Target, r.Status, r.Value)
}

type NotifyEvent struct {
	Type      string `json:"type"`
	Target    string `json:"target"`
	Threshold string `json:"threshold"`
}

type NotifyConfig struct {
	Events []NotifyEvent `json:"events"`
}

type Box struct {
	Id      string `json:"_id"`
	Sensors []struct {
		Id              string `json:"_id"`
		Type            string `json:"sensorType"`
		LastMeasurement *struct {
			Value string    `json:"value"`
			Date  time.Time `json:"createdAt"`
		} `json:"lastMeasurement"`
	} `json:"sensors"`
	NotifyConf *NotifyConfig `json:"notify"`
}

func (box Box) RunChecks() ([]CheckResult, error) {
	var results = []CheckResult{}

	for _, event := range box.NotifyConf.Events {
		target := event.Target

		for _, s := range box.Sensors {
			// if a sensor never measured anything, thats ok. checks would fail anyway
			if s.LastMeasurement == nil {
				continue
			}

			if target == eventTargetAll || target == s.Id {

				switch event.Type {
				case eventMeasurementAge:
					// check if age of lastMeasurement is within threshold
					status := CheckOk
					thresh, err := time.ParseDuration(event.Threshold)
					if err != nil {
						return nil, err
					}
					if time.Since(s.LastMeasurement.Date) > thresh {
						status = CheckErr
					}

					results = append(results, CheckResult{
						Event:  event.Type,
						Target: s.Id,
						Value:  s.LastMeasurement.Date.String(),
						Status: status,
					})

				case eventMeasurementValMin, eventMeasurementValMax:
					status := CheckOk
					thresh, err := strconv.ParseFloat(event.Threshold, 64)
					if err != nil {
						return nil, err
					}
					val, err2 := strconv.ParseFloat(s.LastMeasurement.Value, 64)
					if err2 != nil {
						return nil, err2
					}
					if event.Type == eventMeasurementValMax && val > thresh ||
						event.Type == eventMeasurementValMin && val < thresh {
						status = CheckErr
					}

					results = append(results, CheckResult{
						Event:  event.Type,
						Target: s.Id,
						Value:  s.LastMeasurement.Value,
						Status: status,
					})

				case eventMeasurementValSuspicious:
					status := CheckOk

					val, err := strconv.ParseFloat(s.LastMeasurement.Value, 64)
					if err != nil {
						return nil, err
					}
					if suspiciousVals[SuspiciousValue{
						sensor: s.Type,
						val:    val,
					}] {
						status = CheckErr
					}

					results = append(results, CheckResult{
						Event:  event.Type,
						Target: s.Id,
						Value:  s.LastMeasurement.Value,
						Status: status,
					})
				}
			}
		}
	}
	// must return ALL events to enable Notifier to clear previous notifications
	return results, nil
}
