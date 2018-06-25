package core

import (
	"fmt"
	"strconv"
)

var checkMeasurementMin = checkType{
	name: "measurement_min",
	toString: func(r CheckResult) string {
		return fmt.Sprintf("Sensor %s (%s) reads low value of %s", r.TargetName, r.Target, r.Value)
	},
	checkFunc: validateMeasurementMinMax,
}

var checkMeasurementMax = checkType{
	name: "measurement_min",
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

	thresh, err := strconv.ParseFloat(e.Threshold, 64)
	if err != nil {
		return result, err
	}

	val, err := strconv.ParseFloat(s.LastMeasurement.Value, 64)
	if err != nil {
		return result, err
	}

	if e.Type == eventMeasurementValMax && val > thresh ||
		e.Type == eventMeasurementValMin && val < thresh {
		result.Status = CheckErr
	}

	return result, nil
}
