package core

import (
	"fmt"
	"strconv"
	"time"
)

const (
	CheckOk                       = "OK"
	CheckErr                      = "FAILED"
	eventMeasurementAge           = "measurement_age"
	eventMeasurementValMin        = "measurement_min"
	eventMeasurementValMax        = "measurement_max"
	eventMeasurementValSuspicious = "measurement_suspicious"
	eventTargetAll                = "all" // if event.Target is this value, all sensors will be checked
)

type checkType = struct{ description string }

var checkTypes = map[string]checkType{
	eventMeasurementAge:           checkType{"No measurement from %s since %s"},
	eventMeasurementValMin:        checkType{"Sensor %s reads low value of %s"},
	eventMeasurementValMax:        checkType{"Sensor %s reads high value of %s"},
	eventMeasurementValSuspicious: checkType{"Sensor %s reads presumably faulty value of %s"},
}

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

type NotifyEvent struct {
	Type      string `json:"type"`
	Target    string `json:"target"`
	Threshold string `json:"threshold"`
}

type TransportConfig struct {
	Transport string      `json:"transport"`
	Options   interface{} `json:"options"`
}

type NotifyConfig struct {
	Notifications TransportConfig `json:"notifications"`
	Events        []NotifyEvent   `json:"events"`
}

type Box struct {
	Id      string `json:"_id"`
	Name    string `json:"name"`
	Sensors []struct {
		Id              string `json:"_id"`
		Type            string `json:"sensorType"`
		LastMeasurement *struct {
			Value string    `json:"value"`
			Date  time.Time `json:"createdAt"`
		} `json:"lastMeasurement"`
	} `json:"sensors"`
	NotifyConf *NotifyConfig `json:"healthcheck"`
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
						Threshold: event.Threshold,
						Event:     event.Type,
						Target:    s.Id,
						Value:     s.LastMeasurement.Date.String(),
						Status:    status,
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
						Threshold: event.Threshold,
						Event:     event.Type,
						Target:    s.Id,
						Value:     s.LastMeasurement.Value,
						Status:    status,
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
						Threshold: event.Threshold,
						Event:     event.Type,
						Target:    s.Id,
						Value:     s.LastMeasurement.Value,
						Status:    status,
					})
				}
			}
		}
	}
	// must return ALL events to enable Notifier to clear previous notifications
	return results, nil
}

func (box Box) GetNotifier() (AbstractNotifier, error) {
	transport := box.NotifyConf.Notifications.Transport
	if transport == "" {
		return nil, fmt.Errorf("No notification transport provided")
	}

	notifier := notifiers[transport]
	if notifier == nil {
		return nil, fmt.Errorf("%s is not a supported notification transport", transport)
	}

	return notifier.New(box.NotifyConf.Notifications.Options)
}
