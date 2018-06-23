package core

import (
	"errors"
	"github.com/dghubble/sling"
	"net/http"
)

type OsemError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type OsemClient struct {
	sling *sling.Sling
}

func NewOsemClient(client *http.Client) *OsemClient {
	return &OsemClient{
		sling: sling.New().Client(client).Base("https://api.opensensemap.org/"),
	}
}

func (client *OsemClient) GetBox(boxId string) (Box, error) {
	box := Box{}
	fail := OsemError{}
	client.sling.New().Path("boxes/").Path(boxId).Receive(&box, &fail)
	if fail.Message != "" {
		return box, errors.New("could not fetch box: " + fail.Message)
	}
	return box, nil
}

var osem = NewOsemClient(&http.Client{}) // default client
