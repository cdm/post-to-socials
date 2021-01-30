package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/cdm/post-to-socials/connector"
)

type Result struct {
	Code string `json:"code"`
}

func getAuthHeaders(r *http.Request) (key string, secret string) {
	k := r.Header.Get("k")
	s := r.Header.Get("s")
	return k, s
}

func getNanoTime() string {
	now := time.Now()
	return fmt.Sprintf("%d", now.UnixNano())
}

func writeResult(w http.ResponseWriter, status string) {
	w.Header().Set("Content-Type", "application/json")
	result := Result {
		 Code: status,
	}
	payload, err := json.Marshal(result)
	if err != nil {
		log.WithError(err).Error("Error marshaling response")
	} else {
		w.Write(payload)
	}
}

func startService(conf ConfigVars, creds map[string]string) {
	log.Infof("Starting post-to-socials API service (%s:%s)", conf.Host, conf.Port)

	discord := connector.NewDiscordConnector(
		conf.DiscordChannel, conf.DiscordBotToken, conf.DiscordGuildID)
	twitter := connector.NewTwitterConnector(
		conf.TwitterConsumerKey, conf.TwitterConsumerSecret,
		conf.TwitterAccessTokenKey, conf.TwitterAccessTokenSecret)
	telegram := connector.NewTelegramConnector(conf.TelegramBotToken, conf.TelegramChatIdentifier)

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		symbols := []string{
			"AAPL",
			"MSFT",
			"AMZN",
			"GOOG",
			"TSLA",
			"FB",
			"NFLX",
			"PYPL",
			"INTC",
		}
		rand.Seed(time.Now().Unix())
		w.Write([]byte(symbols[rand.Intn(len(symbols))]))
	})
	router.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("/send/twitter", func(w http.ResponseWriter, r *http.Request) {
		//key, secret := getAuthHeaders(r)
		log.Debug("/send/twitter")

		twitter.Send("Hello Twitter " + getNanoTime())
	})
	router.HandleFunc("/send/discord", func(w http.ResponseWriter, r *http.Request) {
		//key, secret := getAuthHeaders(r)
		log.Debug("/send/discord")

		discord.Send("Hello Discord " + getNanoTime())
	})
	router.HandleFunc("/send/telegram", func(w http.ResponseWriter, r *http.Request) {
		//key, secret := getAuthHeaders(r)
		log.Debug("/send/telegram")

		telegram.Send("Hello Telegram " + getNanoTime())
	})
	router.HandleFunc("/send/all", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("/send/all")
		key, secret := getAuthHeaders(r)
		if len(key) == 0 || len(secret) == 0 || creds[key] != secret {
			log.Debug("auth error")
			writeResult(w, "auth_error")
		} else {
			writeResult(w, "success")
		}
	})

	srv := &http.Server{
		Addr:         conf.Host + ":" + conf.Port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router),
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		derr := discord.Start()
		if derr != nil {
			log.Warn(derr)
		}
		if err := srv.ListenAndServe(); err != nil {
			log.Warn(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for (15 seconds).
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// Signal to stop the discord connector
	discord.Stop()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Info("Shutting down post-to-socials API service")
	os.Exit(0)
}
