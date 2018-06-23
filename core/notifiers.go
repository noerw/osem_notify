package core

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
