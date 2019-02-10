package core

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/spf13/viper"
)

var slackClient = sling.New().Client(&http.Client{})

var notificationColors = map[string]string {
	CheckOk: "#00ff00",
	CheckErr: "#ff0000",
}

// slack Notifier has no configuration
type SlackNotifier struct {
	webhook string
}

type SlackMessage struct {
	Text        string            `json:"text"`
	Username    string            `json:"username,omitempty`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

type SlackAttachment struct {
	Text  string `json:"text"`
	Color string `json:"color,omitempty"`
}

func (n SlackNotifier) New(config TransportConfig) (AbstractNotifier, error) {
	// validate transport configuration
	// :TransportConfSourceHack
	baseUrl := viper.GetString("slack.webhook")
	if baseUrl == "" {
		return nil, fmt.Errorf("Missing configuration key slack.webhook")
	}

	return SlackNotifier{
		webhook: baseUrl,
	}, nil
}

func (n SlackNotifier) Submit(notification Notification) error {
	message := &SlackMessage{
		Username: "osem_notify box healthcheck",
		Text:        notification.Subject,
		Attachments: []SlackAttachment{ { notification.Body, notificationColors[notification.Status] } },
	}

	req, err := slackClient.Post(n.webhook).BodyJSON(message).Request()
	if err != nil {
		return err
	}

	c := http.Client{}
	res, err2 := c.Do(req)
	if err2 != nil {
		return err2
	}

	if res.StatusCode > 200 {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("slack webhook failed: %s", body)
	}

	return nil
}
