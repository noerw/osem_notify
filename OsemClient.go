package main

import (
	"net/http"
	"github.com/dghubble/sling"
)

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
	_, err := client.sling.New().Path("boxes/").Path(boxId).ReceiveSuccess(&box)
	return box, err
}
