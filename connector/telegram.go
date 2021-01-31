package connector

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Telegram struct {
	botToken string
	chatID string
	mu sync.Mutex
}

// Create a new instance of the Telegram connector
func NewTelegramConnector(botToken string, chatID string) *Telegram {
	conn := &Telegram{
		botToken: botToken,
		chatID: chatID,
	}
	return conn
}

// Send a message to Telegram
func (t *Telegram) Send(msg string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	log.Infof("Sending Telegram message `%s` to chat_id: %s", msg, t.chatID)

	telegramApi := "https://api.telegram.org/bot" + t.botToken + "/sendMessage"
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {t.chatID},
			"text":    {msg},
		})
	if err != nil {
		return errors.Wrap(err, "Could not post text to the Telegram chat")
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		return errors.Wrap(errRead, "Could not parse telegram response")
	}
	bodyString := string(bodyBytes)

	log.Infof("Body of Telegram Response: %s", bodyString)
	return nil
}