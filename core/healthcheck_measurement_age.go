package core

import (
	"fmt"
	"time"
)

var checkMeasurementAge = checkType{
	name: "measurement_age",
	toString: func(r CheckResult) string {
		return fmt.Sprintf("No measurement from %s (%s) since %s", r.TargetName, r.Target, r.Value)
	},
	checkFunc: func(e NotifyEvent, s Sensor, b Box) (CheckResult, error) {
		result := CheckResult{
			Event:      e.Type,
			Target:     s.Id,
			TargetName: s.Phenomenon,
			Threshold:  e.Threshold,
			Value:      s.LastMeasurement.Date.String(),
			Status:     CheckOk,
		}

		thresh, err := time.ParseDuration(e.Threshold)
		if err != nil {
			return CheckResult{}, err
		}

		if time.Since(s.LastMeasurement.Date) > thresh {
			result.Status = CheckErr
		}

		result.Value = s.LastMeasurement.Date.String()

		return result, nil
	},
}
