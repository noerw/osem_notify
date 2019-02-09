package core

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	xmpp "github.com/mattn/go-xmpp"
)

// box config required for the XmppNotifier
type XmppNotifier struct {
	Recipients []string
}

func (n XmppNotifier) New(config interface{}) (AbstractNotifier, error) {
	// assign configuration to the notifier after ensuring the correct type.
	// lesson of this project: golang requires us to fuck around with type
	// assertions, instead of providing us with proper inheritance.

	asserted, ok := config.(XmppNotifier)
	if !ok || asserted.Recipients == nil {
		// config did not contain valid options.
		// first try fallback: parse result of viper is a map[string]interface{},
		// which requires a different assertion change
		asserted2, ok := config.(map[string]interface{})
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

func (n XmppNotifier) ComposeNotification(box *Box, checks []CheckResult) Notification {
	errTexts := []string{}
	resolvedTexts := []string{}
	for _, check := range checks {
		if check.Status == CheckErr {
			errTexts = append(errTexts, check.String())
		} else {
			resolvedTexts = append(resolvedTexts, check.String())
		}
	}

	var (
		resolved     string
		resolvedList string
		errList      string
	)
	if len(resolvedTexts) != 0 {
		resolvedList = fmt.Sprintf("Resolved issue(s):\n\n%s\n\n", strings.Join(resolvedTexts, "\n"))
	}
	if len(errTexts) != 0 {
		errList = fmt.Sprintf("New issue(s):\n\n%s\n\n", strings.Join(errTexts, "\n"))
	} else {
		resolved = "resolved "
	}

	return Notification{
		Subject: fmt.Sprintf("Issues %swith your box \"%s\" on opensensemap.org!", resolved, box.Name),
		Body: fmt.Sprintf("A check at %s identified the following updates for your box \"%s\":\n\n%s%sYou may visit https://opensensemap.org/explore/%s for more details.\n\n--\nSent automatically by osem_notify (https://github.com/noerw/osem_notify)",
			time.Now().Round(time.Minute), box.Name, errList, resolvedList, box.Id),
	}
}

func (n XmppNotifier) Submit(notification Notification) error {
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

	client, err := xmppOpts.NewClient()
	if err != nil {
		return err
	}

	for _, recipient := range n.Recipients {
		_, err = client.Send(xmpp.Chat{
			Remote: recipient,
			Subject: notification.Subject,
			Text: fmt.Sprintf("%s\n\n%s", notification.Subject, notification.Body),
		})

		if err != nil {
			return err
		}
	}

	return err
}
