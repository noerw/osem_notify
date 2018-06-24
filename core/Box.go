package core

import (
	"fmt"
	"strconv"
	"time"
)

const (
	CheckOk                   = "OK"
	CheckErr                  = "FAILED"
	eventMeasurementAge       = "measurement_age"
	eventMeasurementValMin    = "measurement_min"
	eventMeasurementValMax    = "measurement_max"
	eventMeasurementValFaulty = "measurement_faulty"
	eventTargetAll            = "all" // if event.Target is this value, all sensors will be checked
)

type checkType = struct{ description string }

var checkTypes = map[string]checkType{
	eventMeasurementAge:       checkType{"No measurement from %s (%s) since %s"},
	eventMeasurementValMin:    checkType{"Sensor %s (%s) reads low value of %s"},
	eventMeasurementValMax:    checkType{"Sensor %s (%s) reads high value of %s"},
	eventMeasurementValFaulty: checkType{"Sensor %s (%s) reads presumably faulty value of %s"},
}

type FaultyValue struct {
	sensor string
	val    float64
}

var faultyVals = map[FaultyValue]bool{
	FaultyValue{sensor: "BMP280", val: 0.0}:  true,
	FaultyValue{sensor: "HDC1008", val: 0.0}: true,
	FaultyValue{sensor: "HDC1008", val: -40}: true,
	FaultyValue{sensor: "SDS 011", val: 0.0}: true,
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
		Phenomenon      string `json:"title"`
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
		for _, s := range box.Sensors {
			// if a sensor never measured anything, thats ok. checks would fail anyway
			if s.LastMeasurement == nil {
				continue
			}

			// a validator must set these values
			var (
				status     = CheckOk
				target     = s.Id
				targetName = s.Phenomenon
				value      string
			)

			if event.Target == eventTargetAll || event.Target == s.Id {

				switch event.Type {
				case eventMeasurementAge:
					thresh, err := time.ParseDuration(event.Threshold)
					if err != nil {
						return nil, err
					}
					if time.Since(s.LastMeasurement.Date) > thresh {
						status = CheckErr
					}

					value = s.LastMeasurement.Date.String()

				case eventMeasurementValMin, eventMeasurementValMax:
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

					value = s.LastMeasurement.Value

				case eventMeasurementValFaulty:
					val, err := strconv.ParseFloat(s.LastMeasurement.Value, 64)
					if err != nil {
						return nil, err
					}
					if faultyVals[FaultyValue{
						sensor: s.Type,
						val:    val,
					}] {
						status = CheckErr
					}

					value = s.LastMeasurement.Value
				}

				results = append(results, CheckResult{
					Threshold:  event.Threshold,
					Event:      event.Type,
					Target:     target,
					TargetName: targetName,
					Value:      value,
					Status:     status,
				})
			}
		}
	}

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
