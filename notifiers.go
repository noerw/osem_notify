package main

type AbstractNotifier interface {
	GetName() string
	SetupTransport(config interface{}) error
	AddNotifications(notifications []Notification) error
	SendNotifications() error
}

type Notification struct {
}

// TODO: multiple transports? one transport per event? (??)
type SlackConfig struct{}

type NotifyEvent struct {
	Type      string `json:"type"`
	Target    string `json:"target"`
	Threshold string `json:"threshold"`
}
