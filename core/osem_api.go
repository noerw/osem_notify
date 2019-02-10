package core

import (
	"errors"
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

type OsemError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BoxFilters struct {
	Date       string `url:"date,omitempty"`
	Exposure   string `url:"exposure,omitempty"`
	Grouptag   string `url:"grouptag,omitempty"`
	Model      string `url:"model,omitempty"`
	Phenomenon string `url:"phenomenon,omitempty"`
}

type OsemClient struct {
	sling *sling.Sling
}

func NewOsemClient(endpoint string) *OsemClient {
	return &OsemClient{
		sling: sling.New().Client(&http.Client{}).Base(endpoint),
	}
}

func (client *OsemClient) GetBox(boxId string) (*Box, error) {
	box := &Box{}
	fail := &OsemError{}
	_, err := client.sling.New().Path("boxes/").Path(boxId).Receive(box, fail)
	if err != nil {
		return nil, err
	}
	if fail.Message != "" {
		return box, errors.New("could not fetch box: " + fail.Message)
	}
	return box, nil
}

func (client *OsemClient) GetAllBoxes(params BoxFilters) (*[]BoxMinimal, error) {
	boxes := &[]BoxMinimal{}
	fail := &OsemError{}
	_, err := client.sling.New().Path("boxes?minimal=true").QueryStruct(params).Receive(boxes, fail)
	if err != nil {
		return nil, err
	}
	if fail.Message != "" {
		return boxes, errors.New("could not fetch boxes: " + fail.Message)
	}
	return boxes, nil
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

type Sensor struct {
	Id              string `json:"_id"`
	Phenomenon      string `json:"title"`
	Type            string `json:"sensorType"`
	LastMeasurement *struct {
		Value string    `json:"value"`
		Date  time.Time `json:"createdAt"`
	} `json:"lastMeasurement"`
}

type Box struct {
	Id         string        `json:"_id"`
	Name       string        `json:"name"`
	Sensors    []Sensor      `json:"sensors"`
	NotifyConf *NotifyConfig `json:"healthcheck"`
}

type BoxMinimal struct {
	Id   string `json:"_id"`
	Name string `json:"name"`
}
