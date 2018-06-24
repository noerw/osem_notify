package core

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var notifiers = map[string]AbstractNotifier{
	"email": EmailNotifier{},
}

type AbstractNotifier interface {
	New(config interface{}) (AbstractNotifier, error)
	ComposeNotification(box *Box, checks []CheckResult) Notification
	Submit(notification Notification) error
}

type Notification struct {
	body    string
	subject string
}

// box config required for the EmailNotifier
type EmailNotifier struct {
	Recipients []string
}

func (n EmailNotifier) New(config interface{}) (AbstractNotifier, error) {
	res, ok := config.(EmailNotifier)

	if !ok || res.Recipients == nil {
		return nil, errors.New("Invalid EmailNotifier options")
	}

	return EmailNotifier{
		Recipients: res.Recipients,
	}, nil
}

func (n EmailNotifier) ComposeNotification(box *Box, checks []CheckResult) Notification {
	resultTexts := []string{}
	for _, check := range checks {
		resultTexts = append(resultTexts, check.String())
	}

	return Notification{
		subject: fmt.Sprintf("Issues with your box \"%s\" on opensensemap.org!", box.Name),
		body: fmt.Sprintf("A check at %s identified the following issue(s) with your box %s:\n\n%s\n\nYou may visit https://opensensemap.org/explore/%s for more details.\n\n--\nSent automatically by osem_notify (https://github.com/noerw/osem_notify)",
			time.Now().Round(time.Minute), box.Name, strings.Join(resultTexts, "\n"), box.Id),
	}
}

func (n EmailNotifier) Submit(notification Notification) error {
	auth := smtp.PlainAuth(
		"",
		viper.GetString("email.user"),
		viper.GetString("email.pass"),
		viper.GetString("email.host"),
	)

	from := viper.GetString("email.from")
	body := fmt.Sprintf("From: openSenseMap Notifier <%s>\nSubject: %s\nContent-Type: text/plain; charset=\"utf-8\"\n\n%s", from, notification.subject, notification.body)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", viper.GetString("email.host"), viper.GetString("email.port")),
		auth,
		from,
		n.Recipients,
		[]byte(body),
	)

	return err
}
