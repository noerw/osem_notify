package core

import (
	"errors"
	"fmt"

	xmpp "github.com/mattn/go-xmpp"
	"github.com/spf13/viper"
)

var xmppClient = &xmpp.Client{} // @Hacky

// box config required for the XmppNotifier (TransportConfig.Options)
type XmppNotifier struct {
	Recipients []string
}

func (n XmppNotifier) New(config TransportConfig) (AbstractNotifier, error) {
	// validate transport configuration
	// :TransportConfSourceHack
	requiredConf := []string{"xmpp.user", "xmpp.pass", "xmpp.host", "xmpp.starttls"}
	for _, key := range requiredConf {
		if viper.GetString(key) == "" {
			return nil, fmt.Errorf("Missing configuration key %s", key)
		}
	}

	// establish connection with server once, and share it accross instances
	// @Hacky
	if xmppClient.JID() == "" {
		c, err := connectXmpp()
		if err != nil {
			return nil, err
		}
		xmppClient = c
	}

	// assign configuration to the notifier after ensuring the correct type.
	// lesson of this project: golang requires us to fuck around with type
	// assertions, instead of providing us with proper inheritance.
	asserted, ok := config.Options.(XmppNotifier)
	if !ok || asserted.Recipients == nil {
		// config did not contain valid options.
		// first try fallback: parse result of viper is a map[string]interface{},
		// which requires a different assertion change
		asserted2, ok := config.Options.(map[string]interface{})
		if ok {
			asserted3, ok := asserted2["recipients"].([]interface{})
			if ok {
				asserted = XmppNotifier{Recipients: []string{}}
				for _, rec := range asserted3 {
					asserted.Recipients = append(asserted.Recipients, rec.(string))
				}
			}
		}

		if asserted.Recipients == nil {
			return nil, errors.New("Invalid XmppNotifier options")
		}
	}

	return XmppNotifier{
		Recipients: asserted.Recipients,
	}, nil
}

func (n XmppNotifier) Submit(notification Notification) error {
	if xmppClient.JID() == "" {
		return fmt.Errorf("xmpp client not correctly initialized!")
	}

	for _, recipient := range n.Recipients {
		_, err := xmppClient.Send(xmpp.Chat{
			Remote:  recipient,
			Subject: notification.Subject,
			Text:    fmt.Sprintf("%s\n\n%s", notification.Subject, notification.Body),
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func connectXmpp() (*xmpp.Client, error) {
	// :TransportConfSourceHack
	xmppOpts := xmpp.Options{
		Host:     viper.GetString("xmpp.host"),
		User:     viper.GetString("xmpp.user"),
		Password: viper.GetString("xmpp.pass"),
		Resource: "osem_notify",
	}

	if viper.GetBool("xmpp.starttls") {
		xmppOpts.NoTLS = true
		xmppOpts.StartTLS = true
	}

	return xmppOpts.NewClient()
}
