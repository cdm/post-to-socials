package connector

import (
	"sync"

	"github.com/pkg/errors"

	dg "github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type Discord struct {
   session *dg.Session
   channel string
   botName string
   mu sync.Mutex
}

// Create a new instance of the Discord connector
func NewDiscordConnector(channel string, botName string) *Discord {
	conn := &Discord{
		channel: channel,
		botName: botName, //"Mixel 0.1",
	}
	return conn
}

// Start the Discord connector
func (d *Discord) Start() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.session != nil {
		log.Warn("Discord service already started")
		return nil
	}
	log.Info("Discord connector starting...")

	session, err := dg.New(d.botName)
	if err != nil {
		return errors.Wrap(err,"Could not create Discord session")
	}
	d.session = session

	// Open a websocket connection to Discord and begin listening
	err = d.session.Open()
	if err != nil {
		return errors.Wrap(err,"Could not open Discord connection")
	}

	log.Info("Discord connector started")
	return nil
}

// Stop the Discord connector
func (d *Discord) Stop() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.session == nil {
		log.Warn("Discord service not started")
		return nil
	}
	log.Info("Discord connector stopping...")

	// Cleanly close down the Discord session
	err := d.session.Close()
	if err == nil {
		log.Info("Discord connector stopped")
	} else {
		log.Warn("Discord connector may not have stopped correctly")
	}
	return err
}

// Send a message to Discord
func (d *Discord) Send(msg string) error {
	output, err := d.session.ChannelMessageSend(d.channel, msg)
	if err != nil {
		return err
	}
	log.Info(output)
	return nil
}