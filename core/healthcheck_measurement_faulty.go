package core

import (
	"fmt"

	"github.com/noerw/osem_notify/utils"
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

		val, err := utils.ParseFloat(s.LastMeasurement.Value)
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
	// @TODO: add UV & light sensor: check for 0 if not sunset based on boxlocation
	// @TODO: add BME280 and other sensors..
	faultyValue{sensor: "BMP280", val: 0.0}:  true,
	faultyValue{sensor: "HDC1008", val: 0.0}: true, // @FIXME: check should be on luftfeuchte only!
	faultyValue{sensor: "HDC1008", val: -40}: true,
	faultyValue{sensor: "SDS 011", val: 0.0}: true, // @FIXME: 0.0 seems to be a correct value, need to check over longer periods
}
