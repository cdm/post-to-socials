package connector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/dghubble/oauth1"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Represents a Tweet on Twitter
type Tweet struct {
	Id    int64
	IdStr string `json:"id_str"`
	User  User
	Text  string
}

// Represents a User on Twitter
type User struct {
	Id     int64
	IdStr  string `json:"id_str"`
	Name   string
	Handle string `json:"screen_name"`
}

type Twitter struct {
	mu sync.Mutex
	consumerKey string
	consumerSecret string
	accessTokenKey string
	accessTokenSecret string
}

// Create a new instance of the Twitter connector
func NewTwitterConnector(
		consumerKey string, consumerSecret string,
		accessTokenKey string, accessTokenSecret string) *Twitter {
	conn := &Twitter{
		consumerKey: consumerKey,
		consumerSecret: consumerSecret,
		accessTokenKey: accessTokenKey,
		accessTokenSecret: accessTokenSecret,
	}
	return conn
}

// Send a message to Twitter
func (t *Twitter) Send(msg string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	tweet, err := t.sendTweet(msg)
	if err != nil {
		log.Errorf("Response from Twitter: %s", tweet)
		return errors.Wrap(err, "Could not send message to Twitter")
	}
	return nil
}

// Sends a Tweet by the authenticated account (max 140 chars)
func (t *Twitter) sendTweet(tweet string) (*Tweet, error) {
	params := url.Values{}
	params.Set("status",tweet)

	client := t.oauthClient()
	resp, err := client.PostForm("https://api.twitter.com/1.1/statuses/update.json", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	var res Tweet
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, errors.Wrap(err, "Could not unmarshal response")
	}
	return &res, nil
}

// Creates an OAuth client with specified consumer keys and access tokens
func (t *Twitter)  oauthClient() *http.Client {
	config := oauth1.NewConfig(t.consumerKey, t.consumerSecret)
	token := oauth1.NewToken(t.accessTokenKey, t.accessTokenSecret)
	return config.Client(oauth1.NoContext, token)
}
