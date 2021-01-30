package main

import "github.com/ilyakaznacheev/cleanenv"

type ConfigVars struct {
	Port                      string `yaml:"port" env:"PORT" env-default:"5432"`
	Host                      string `yaml:"host" env:"HOST" env-default:"localhost"`
	BotName                   string `yaml:"name" env:"BOTNAME" env-default:"Mixel"`
	CredentialFilePath        string `yaml:"csv" env:"CSV" env-default:"auth.csv"`
	TwitterConsumerKey        string `yaml:"twitterConsumerKey" env:"TCK" env-default:""`
	TwitterConsumerSecret     string `yaml:"twitterConsumerSecret" env:"TCS" env-default:""`
	TwitterAccessTokenKey     string `yaml:"twitterAccessTokenKey" env:"TATK" env-default:""`
	TwitterAccessTokenSecret  string `yaml:"twitterAccessTokenSecret" env:"TATS" env-default:""`
	TelegramChatIdentifier    string `yaml:"telegramChatId" env:"TCI" env-default:""`
	TelegramBotToken          string `yaml:"telegramBotToken" env:"TBT" env-default:""`
	DiscordChannel            string `yaml:"discordChannel" env:"DCH" env-default:""`
	DiscordToken              string `yaml:"discordToken" env:"DTOK" env-default:""`
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

