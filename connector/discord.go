package connector

import (
	"sync"

	dg "github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Discord struct {
	session    *dg.Session
	channel    string
	channelID  string
	botToken   string
	botGuildID string
	mu         sync.Mutex
}

// Create a new instance of the Discord connector
func NewDiscordConnector(channelName string, botToken string, botGuildID string) *Discord {
	conn := &Discord{
		channel:    channelName,
		botToken:   botToken,
		botGuildID: botGuildID,
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
	log.Info("Discord connector starting..." + d.botToken)

	session, err := dg.New("Bot " + d.botToken)
	if err != nil {
		return errors.Wrap(err, "Could not create Discord session")
	}

	session.Identify.Intents = dg.IntentsGuilds | dg.IntentsGuildMessages

	session.AddHandler(ready)
	d.session = session

	// Open a websocket connection to Discord and begin listening
	err = d.session.Open()
	if err != nil {
		return errors.Wrap(err, "Could not open Discord connection")
	}

	log.Info("Discord connector started")

	// Loop through each guild in the session
	foundGuild := false
	foundChannel := false
	for _, guild := range session.State.Guilds {
		log.Infof("Found guild [%s]", guild.ID)

		if guild.ID == d.botGuildID {
			log.Info("Guild matched")
			channels, _ := session.GuildChannels(guild.ID)
			for _, c := range channels {
				// Check if channel is a guild text channel
				// and not a voice or a DM channel
				if c.Type != dg.ChannelTypeGuildText {
					continue
				}
				if c.Name == d.channel {
					log.Info(c.Name)
					log.Info(c.ID)

					d.channelID = c.ID
					foundChannel = true
					break
				}
			}
			foundGuild = true
			break
		}
	}
	if !foundGuild {
		log.Errorf("Could not find matching guild with ID %s", d.botGuildID)
	} else if !foundChannel {
		log.Errorf("Could not find matching channel named ID %s", d.channel)
	}

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

func ready(s *dg.Session, event *dg.Ready) {
	log.Info("Discord connector is ready")
}

// Send a message to Discord
func (d *Discord) Send(msg string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	log.Info("Sending `" + msg + "` to " + d.channel + "(" + d.channelID + ")")
	output, err := d.session.ChannelMessageSend(d.channelID, msg)
	if err != nil {
		return err
	}
	log.Info(output)
	return nil
}
