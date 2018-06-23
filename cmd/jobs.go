package cmd

import (
	"../core"
	log "github.com/sirupsen/logrus"
)

func CheckBoxes(boxIds []string, defaultConf *core.NotifyConfig) ([]core.CheckResult, error) {
	log.Debug("Checking notifications for ", len(boxIds), " box(es)")

	// TODO: return a map of Box: []Notification instead?
	results := []core.CheckResult{}
	for _, boxId := range boxIds {
		r, err := checkBox(boxId, defaultConf)
		if err != nil {
			return nil, err
		}

		if r != nil {
			results = append(results, r...)
		}
	}

	return results, nil
}

func checkBox(boxId string, defaultConf *core.NotifyConfig) ([]core.CheckResult, error) {
	boxLogger := log.WithFields(log.Fields{"boxId": boxId})
	boxLogger.Info("checking box for due notifications")

	// get box data
	box, err := core.Osem.GetBox(boxId)
	if err != nil {
		boxLogger.Error(err)
		return nil, err
	}

	// if box has no notify config, we use the defaultConf
	if box.NotifyConf == nil {
		box.NotifyConf = defaultConf
	}

	// run checks
	results, err2 := box.RunChecks()
	if err2 != nil {
		boxLogger.Error("could not run checks on box: ", err2)
		return results, err2
	}

	for _, r := range results {
		resultLog := boxLogger.WithFields(log.Fields{
			"status": r.Status,
			"event":  r.Event,
			"value":  r.Value,
			"target": r.Target,
		})
		if r.Status == core.CheckOk {
			resultLog.Debug(r)
		} else {
			resultLog.Warn(r)
		}
	}

	return results, nil
}
