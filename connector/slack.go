package connector

import (
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type Slack struct {
	token     string
	channelID string
	mu        sync.Mutex
}

// Create a new instance of the Slack connector
func NewSlackConnector(channelID string, token string) *Slack {
	conn := &Slack{
		token:     token,
		channelID: channelID,
	}
	return conn
}

// Send a message to Slack
func (s *Slack) Send(msg string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Infof("Sending Slack message `%s` to chat_id: %s", msg, s.channelID)

	api := slack.New(s.token)

	channelID, timestamp, err := api.PostMessage(
		s.channelID,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAsUser(true), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
	)
	if err != nil {
		return errors.Wrap(err, "Could not post text to the Telegram chat")
	}

	log.Infof("Message sent to %s at: %s", channelID, timestamp)
	return nil
}
