package main

import "github.com/ilyakaznacheev/cleanenv"

type ConfigVars struct {
	Port                     string `yaml:"port" env:"PORT" env-default:"5432"`
	Host                     string `yaml:"host" env:"HOST" env-default:"localhost"`
	CredentialFilePath       string `yaml:"csv" env:"CSV" env-default:"auth.csv"`
	TwitterConsumerKey       string `yaml:"twitterConsumerKey" env:"TWITTER-CK" env-default:""`
	TwitterConsumerSecret    string `yaml:"twitterConsumerSecret" env:"TWITTER-CS" env-default:""`
	TwitterAccessTokenKey    string `yaml:"twitterAccessTokenKey" env:"TWITTER-ATK" env-default:""`
	TwitterAccessTokenSecret string `yaml:"twitterAccessTokenSecret" env:"TWITTER-ATS" env-default:""`
	TelegramChatIdentifier   string `yaml:"telegramChatId" env:"TELEGRAM-CHAT-IDENTIFIER" env-default:""`
	TelegramBotToken         string `yaml:"telegramBotToken" env:"TELEGRAM-TOKEN" env-default:""`
	DiscordGuildID           string `yaml:"discordGuildId" env:"DISCORD-GUILD-ID" env-default:""`
	DiscordChannel           string `yaml:"discordChannel" env:"DISCORD-CHANNEL" env-default:""`
	DiscordBotToken          string `yaml:"discordBotToken" env:"DISCORD-TOKEN" env-default:""`
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
