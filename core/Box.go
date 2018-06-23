package core

type NotifyEvent struct {
	Type      string `json:"type"`
	Target    string `json:"target"`
	Threshold string `json:"threshold"`
}

type NotifyConfig struct {
	// Transports interface{} `json:"transports"`
	Events []NotifyEvent `json:"events"`
}

type Box struct {
	Id      string `json:"_id"`
	Sensors []struct {
		Id              string `json:"_id"`
		LastMeasurement *struct {
			Value string `json:"value"`
			Date  string `json:"createdAt"`
		} `json:"lastMeasurement"`
	} `json:"sensors"`
	NotifyConf *NotifyConfig `json:"notify"`
}

func (box Box) runChecks() ([]Notification, error) {
	// must return ALL events to enable Notifier to clear previous notifications
	return nil, nil
}

func (box Box) getNotifier() (AbstractNotifier, error) {
	// validate box.NotifyConf.transport

	// try to get notifier state from persistence

	// return
	var notifier AbstractNotifier
	return notifier, nil
}
