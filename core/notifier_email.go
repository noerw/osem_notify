package core

import (
	"errors"
	"fmt"
	"net/smtp"
	"time"

	"github.com/spf13/viper"
)

// box config required for the EmailNotifier (TransportConfig.Options)
type EmailNotifier struct {
	Recipients []string
}

func (n EmailNotifier) New(config TransportConfig) (AbstractNotifier, error) {
	// validate transport configuration
	// :TransportConfSourceHack @FIXME: dont get these values from viper, as the core package
	// should be agnostic of the source of configuration!
	requiredConf := []string{"email.user", "email.pass", "email.host", "email.port", "email.from"}
	for _, key := range requiredConf {
		if viper.GetString(key) == "" {
			return nil, fmt.Errorf("Missing configuration key %s", key)
		}
	}

	// assign configuration to the notifier after ensuring the correct type.
	// lesson of this project: golang requires us to fuck around with type
	// assertions, instead of providing us with proper inheritance.
	asserted, ok := config.Options.(EmailNotifier)
	if !ok || asserted.Recipients == nil {
		// config did not contain valid options.
		// first try fallback: parse result of viper is a map[string]interface{},
		// which requires a different assertion change
		asserted2, ok := config.Options.(map[string]interface{})
		if ok {
			asserted3, ok := asserted2["recipients"].([]interface{})
			if ok {
				asserted = EmailNotifier{Recipients: []string{}}
				for _, rec := range asserted3 {
					asserted.Recipients = append(asserted.Recipients, rec.(string))
				}
			}
		}

		if asserted.Recipients == nil {
			return nil, errors.New("Invalid EmailNotifier options")
		}
	}

	return EmailNotifier{
		Recipients: asserted.Recipients,
	}, nil
}

func (n EmailNotifier) Submit(notification Notification) error {
	// :TransportConfSourceHack
	auth := smtp.PlainAuth(
		"",
		viper.GetString("email.user"),
		viper.GetString("email.pass"),
		viper.GetString("email.host"),
	)

	from := viper.GetString("email.from")
	body := fmt.Sprintf(
		"From: openSenseMap Notifier <%s>\nDate: %s\nSubject: %s\nContent-Type: text/plain; charset=\"utf-8\"\n\n%s",
		from,
		time.Now().Format(time.RFC1123Z),
		notification.Subject,
		notification.Body)

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
