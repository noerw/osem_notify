package core

import (
	"errors"
	"net/http"

	"github.com/dghubble/sling"
)

type OsemError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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
	client.sling.New().Path("boxes/").Path(boxId).Receive(box, fail)
	if fail.Message != "" {
		return box, errors.New("could not fetch box: " + fail.Message)
	}
	return box, nil
}

