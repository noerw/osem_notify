package core

import (
	"fmt"

	"github.com/noerw/osem_notify/utils"
)

const (
	nameMin = "measurement_min"
	nameMax = "measurement_max"
)

var checkMeasurementMin = checkType{
	name: nameMin,
	toString: func(r CheckResult) string {
		return fmt.Sprintf("Sensor %s (%s) reads low value of %s", r.TargetName, r.Target, r.Value)
	},
	checkFunc: validateMeasurementMinMax,
}

var checkMeasurementMax = checkType{
	name: nameMax,
	toString: func(r CheckResult) string {
		return fmt.Sprintf("Sensor %s (%s) reads high value of %s", r.TargetName, r.Target, r.Value)
	},
	checkFunc: validateMeasurementMinMax,
}

func validateMeasurementMinMax(e NotifyEvent, s Sensor, b Box) (CheckResult, error) {
	result := CheckResult{
		Event:      e.Type,
		Target:     s.Id,
		TargetName: s.Phenomenon,
		Threshold:  e.Threshold,
		Value:      s.LastMeasurement.Value,
		Status:     CheckOk,
	}

	thresh, err := utils.ParseFloat(e.Threshold)
	if err != nil {
		return result, err
	}

	val, err := utils.ParseFloat(s.LastMeasurement.Value)
	if err != nil {
		return result, err
	}

	if e.Type == nameMax && val > thresh ||
		e.Type == nameMin && val < thresh {
		result.Status = CheckErr
	}

	return result, nil
}
