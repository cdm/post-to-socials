package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/cdm/post-to-socials/connector"
)

type Result struct {
	Code string `json:"code"`
}

type SendMessage struct {
	Msg string `json:"message"`
}

func getAuthHeaders(r *http.Request) (key string, secret string) {
	k := r.Header.Get("key")
	s := r.Header.Get("secret")
	return k, s
}

func getNanoTime() string {
	now := time.Now()
	return fmt.Sprintf("%d", now.UnixNano())
}

func writeResult(w http.ResponseWriter, status string) {
	w.Header().Set("Content-Type", "application/json")
	result := Result{
		Code: status,
	}
	payload, err := json.Marshal(result)
	if err != nil {
		log.WithError(err).Error("Error marshaling response")
	} else {
		w.Write(payload)
	}
}

func validate(w http.ResponseWriter, r *http.Request, creds map[string]string) (bool, string) {
	if r.Method != http.MethodPost {
		writeResult(w, "method_error")
		return false, ""
	}

	key, secret := getAuthHeaders(r)
	log.Info("key: " + key)
	log.Info("secret: " + secret)

	if len(key) == 0 || len(secret) == 0 {
		writeResult(w, "auth_error")
		return false, ""
	}

	if _, ok := creds[key]; !ok {
		writeResult(w, "auth_error")
		return false, ""
	}

	if creds[key] != secret {
		writeResult(w, "auth_error")
		return false, ""
	}

	decoder := json.NewDecoder(r.Body)
	var m SendMessage
	err := decoder.Decode(&m)
	if err != nil {
		log.Errorf("Error decoding msg body: %s", err.Error())
		writeResult(w, "json_error")
		return false, ""
	}

	if len(m.Msg) > 140 || len(m.Msg) == 0 {
		log.Errorf("Error msg body is out of range: %d >> '%s'", len(m.Msg), m.Msg)
		writeResult(w, "json_error")
		return false, ""
	}

	return true, m.Msg
}

func postMessage(key string, secret string, msg string, path string) (error, string) {
	var jsonStr = []byte(`{"message":"` + msg + `"}`)
	req, err := http.NewRequest("POST", "http://"+path, bytes.NewBuffer(jsonStr))
	req.Header.Set("key", key)
	req.Header.Set("secret", secret)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, ""
	}
	defer resp.Body.Close()

	body := ""
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		body = string(bodyBytes)
	}
	return nil, body
}

func startService(conf ConfigVars, creds map[string]string) {
	log.Infof("Starting post-to-socials API service (%s:%s)", conf.Host, conf.Port)

	discord := connector.NewDiscordConnector(
		conf.DiscordChannel,
		conf.DiscordBotToken,
		conf.DiscordGuildID,
	)
	twitter := connector.NewTwitterConnector(
		conf.TwitterConsumerKey,
		conf.TwitterConsumerSecret,
		conf.TwitterAccessTokenKey,
		conf.TwitterAccessTokenSecret,
	)
	telegram := connector.NewTelegramConnector(
		conf.TelegramBotToken,
		conf.TelegramChatIdentifier,
	)

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("/send/twitter", func(w http.ResponseWriter, r *http.Request) {
		log.Info("/send/twitter")
		valid, msg := validate(w, r, creds)
		if valid {
			err := twitter.Send(msg)
			if err != nil {
				log.Errorf("Error sending msg to Twitter: %s", err.Error())
				writeResult(w, "send_error")
			} else {
				writeResult(w, "success")
			}
		}
	})
	router.HandleFunc("/send/discord", func(w http.ResponseWriter, r *http.Request) {
		log.Info("/send/discord")
		valid, msg := validate(w, r, creds)
		if valid {
			err := discord.Send(msg)
			if err != nil {
				log.Errorf("Error sending msg to Discord: %s", err.Error())
				writeResult(w, "send_error")
			} else {
				writeResult(w, "success")
			}
		}
	})
	router.HandleFunc("/send/telegram", func(w http.ResponseWriter, r *http.Request) {
		log.Info("/send/telegram")
		valid, msg := validate(w, r, creds)
		if valid {
			err := telegram.Send(msg)
			if err != nil {
				log.Errorf("Error sending msg to Telegram: %s", err.Error())
				writeResult(w, "send_error")
			} else {
				writeResult(w, "success")
			}
		}
	})
	router.HandleFunc("/send/all", func(w http.ResponseWriter, r *http.Request) {
		log.Info("/send/all")
		valid, msg := validate(w, r, creds)
		if valid {
			errors := 0
			err := discord.Send(msg)
			if err != nil {
				log.Errorf("Error sending msg to Discord: %s", err.Error())
				errors++
			}
			err = nil
			err = telegram.Send(msg)
			if err != nil {
				log.Errorf("Error sending msg to Telegram: %s", err.Error())
				errors++
			}
			err = nil
			err = twitter.Send(msg)
			if err != nil {
				log.Errorf("Error sending msg to Twitter: %s", err.Error())
				errors++
			}
			if errors > 0 {
				writeResult(w, "send_error")
			} else {
				writeResult(w, "success")
			}
		}
	})
	router.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			key := r.FormValue("key")
			secret := r.FormValue("secret")
			message := r.FormValue("message")
			discord := r.FormValue("discord")
			twitter := r.FormValue("twitter")
			telegram := r.FormValue("telegram")

			if len(discord) == 0 && len(twitter) == 0 && len(telegram) == 0 {
				writeResult(w, "network_required_error")
				return
			}
			if len(key) == 0 || len(secret) == 0 {
				writeResult(w, "auth_required_error")
				return
			}
			if len(message) == 0 || len(message) > 140 {
				writeResult(w, "message_length_error")
				return
			}

			responseContent := ""
			if len(discord) > 0 && discord == "discord" {
				err, body := postMessage(key, secret, message, conf.Host+":"+conf.Port+"/send/discord")
				if err != nil {
					log.Error(errors.Wrap(err, "Error posting to Discord endpoint on service"))
					writeResult(w, "post_error")
					return
				}
				responseContent = body
			}

			if len(twitter) > 0 && twitter == "twitter" {
				err, body := postMessage(key, secret, message, conf.Host+":"+conf.Port+"/send/twitter")
				if err != nil {
					log.Error(errors.Wrap(err, "Error posting to Twitter endpoint on service"))
					writeResult(w, "post_error")
					return
				}
				responseContent = body
			}

			if len(telegram) > 0 && telegram == "telegram" {
				err, body := postMessage(key, secret, message, conf.Host+":"+conf.Port+"/send/telegram")
				if err != nil {
					log.Error(errors.Wrap(err, "Error posting to telegram endpoint on service"))
					writeResult(w, "post_error")
					return
				}
				responseContent = body
			}

			if strings.Contains(responseContent, "success") {
				writeResult(w, "post_success")
			} else {
				errText := "unknown_error"
				re := regexp.MustCompile("\\{\"code\":\"(.*?)\"\\}")
				res := re.FindAllStringSubmatch(responseContent, 1)
				for i := range res {
					log.Info("Error code from API:", res[i][1])
					errText = res[i][1]
				}
				writeResult(w, errText)
			}
		} else {
			// Serve up html
			http.ServeFile(w, r, "form.html")
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
