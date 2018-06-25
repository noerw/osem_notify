package core

import (
	"fmt"
	"strconv"
)

var checkMeasurementFaulty = checkType{
	name: "measurement_faulty",
	toString: func(r CheckResult) string {
		return fmt.Sprintf("Sensor %s (%s) reads presumably faulty value of %s", r.TargetName, r.Target, r.Value)
	},
	checkFunc: func(e NotifyEvent, s Sensor, b Box) (CheckResult, error) {
		result := CheckResult{
			Event:      e.Type,
			Target:     s.Id,
			TargetName: s.Phenomenon,
			Threshold:  e.Threshold,
			Value:      s.LastMeasurement.Value,
			Status:     CheckOk,
		}

		val, err := strconv.ParseFloat(s.LastMeasurement.Value, 64)
		if err != nil {
			return result, err
		}

		if faultyVals[faultyValue{
			sensor: s.Type,
			val:    val,
		}] {
			result.Status = CheckErr
		}

		return result, nil
	},
}

type faultyValue struct {
	sensor string
	val    float64
}

var faultyVals = map[faultyValue]bool{
	faultyValue{sensor: "BMP280", val: 0.0}:  true,
	faultyValue{sensor: "HDC1008", val: 0.0}: true,
	faultyValue{sensor: "HDC1008", val: -40}: true,
	faultyValue{sensor: "SDS 011", val: 0.0}: true,
}
