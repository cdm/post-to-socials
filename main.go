package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Logger config
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// Config
	conf, err := ReadConfig("config.yaml")
	if err != nil {
		log.Fatal("Failed to read config: ", err)
	}

	if len(conf.Host) <= 0 {
		log.Printf("Error: missing 'host' config (the address to bind the accounts service to)")
		return
	}
	if len(conf.Port) <= 0 {
		log.Printf("Error: missing 'port' config (the address to bind the accounts service to)")
		return
	}
	if len(conf.CredentialFilePath) <= 0 {
		log.Printf("Error: missing 'csv' config (the credentials file path for API users)")
		return
	}

	creds, err := LoadCredentialsFromCsv(conf.CredentialFilePath)
	if err != nil {
		log.WithError(err).Fatal("Fatal error loading excluded parties from csv")
	}

	startService(conf, creds)
}
