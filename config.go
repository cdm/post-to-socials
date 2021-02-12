package main

import "github.com/ilyakaznacheev/cleanenv"

type ConfigVars struct {
	Port                     string `yaml:"port" env:"PORT" env-default:"5432"`
	Host                     string `yaml:"host" env:"HOST" env-default:""`
	CredentialFilePath       string `yaml:"csv" env:"CSV" env-default:"auth.csv"`
	TwitterEnabled           bool   `yaml:"twitterEnabled" env:"TWITTER-ENABLED" env-default:"false"`
	TwitterConsumerKey       string `yaml:"twitterConsumerKey" env:"TWITTER-CK" env-default:""`
	TwitterConsumerSecret    string `yaml:"twitterConsumerSecret" env:"TWITTER-CS" env-default:""`
	TwitterAccessTokenKey    string `yaml:"twitterAccessTokenKey" env:"TWITTER-ATK" env-default:""`
	TwitterAccessTokenSecret string `yaml:"twitterAccessTokenSecret" env:"TWITTER-ATS" env-default:""`
	TelegramEnabled          bool   `yaml:"telegramEnabled" env:"TELEGRAM-ENABLED" env-default:"false"`
	TelegramChatIdentifier   string `yaml:"telegramChatId" env:"TELEGRAM-CHAT-IDENTIFIER" env-default:""`
	TelegramBotToken         string `yaml:"telegramBotToken" env:"TELEGRAM-TOKEN" env-default:""`
	DiscordEnabled           bool   `yaml:"discordEnabled" env:"DISCORD-ENABLED" env-default:"false"`
	DiscordGuildID           string `yaml:"discordGuildId" env:"DISCORD-GUILD-ID" env-default:""`
	DiscordChannel           string `yaml:"discordChannel" env:"DISCORD-CHANNEL" env-default:""`
	DiscordBotToken          string `yaml:"discordBotToken" env:"DISCORD-TOKEN" env-default:""`
	SlackEnabled             bool   `yaml:"slackEnabled" env:"SLACK-ENABLE" env-default:"false"`
	SlackToken               string `yaml:"slackToken" env:"SLACK-TOKEN" env-default:""`
	SlackChannelID           string `yaml:"slackChannelID" env:"SLACK-CHANNEL-ID" env-default:""`
}

//"config.yml"
func ReadConfig(path string) (ConfigVars, error) {
	var cfg ConfigVars
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
