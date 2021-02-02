# 📣 Post to socials API service

A simple REST based API service written in Go which relays a message via the API to social messaging networks. The service currently supports posting 140 char messages to Twitter, Telegram (chat) and Discord (channel).

## Getting Started

This application requires Go 1.14 or above installed on your machine. Tested on linux and macosx.

### Building

Check out the repository locally and run:
```
make build
```

### Configuration

Edit the file `config.yaml` and fill in the required configuration values:

```
host => Post to socials API service hostname e.g. "localhost"  
port => Post to socials API service port e.g. "5432"

csv => Path to csv of authentication key/secret pairs e.g. "csv/auth.csv"

twitterConsumerKey       => Twitter API consumer key e.g. "xgtOU2XljrxBBQuWzkWfe4ebp"
twitterConsumerSecret    => Twitter API consumer secret e.g. "HMcXX73htvoe73mS2NR2cY9aag3r9D8CmbJtAlMPEfXRD1Qucp" 
twitterAccessTokenKey    => Twitter API access token e.g. "1355884211619828228-NjPbG8kFxVUxjqMWFsPfS7CKfT3vOM"
twitterAccessTokenSecret => Twitter API access token secret e.g. "e2E1iLEWGqL1JMB0oQIPsH1omEyalXj9hcRHi6jAMyGaC"

telegramChatId   => Telegram chat identifier to post to e.g. "-1991156818728"
telegramBotToken => Telegram bot token e.g. "1441979128:AKKpe6CDKXoChpQ2KJ-fwej6F2qjjhbmdNA"

discordGuildId  => Discord guild/server identifier e.g. "120571334818737489"
discordChannel  => Discord channel to post to e.g. "bot-playground"
discordBotToken => Discord bot token e.g "XD1MTE2MzkyMDE3NDk0MDE2.YFFNEg.14wlVKNs5ITFgDLzT3N-bffOIsE" 
```

To help with deployment on cloud hosting providers, the config also supports Environment (ENV) variables. 
The named ENV vars are shown directly in the [configuration source code](https://github.com/cdm/post-to-socials/blob/d19b00279aa2c88ed5c7ee45a701869799e06ecd/config.go#L6).

## Usage

Once built and configured, run the service binary using the following command:

```
./post-to-socials
```

To send a message, construct a JSON POST message to send with the following HTTP headers:

`key` => Post to socials API authentication key (from auth.csv file)  
`secret` => Post to socials API authentication secret (from auth.csv file)  

The message body should include the JSON body content as follows:
```
{"message":"Your message body here, max 140 characters :)"}
```

Please see the different API methods in the next section for your choice of network to relay to.

### API methods

`/send/discord` `POST` Send a Discord message (max 140 chars) to the guild/channel specified in the config file.

`/send/telegram` `POST` Send a Telegram message (max 140 chars) to the chat ID specified in the config file.

`/send/twitter` `POST` Send a Twitter message (max 140 chars) to the profile specified (API keys) in the config file.

`/send/all` `POST` Send a message to all three social channels in one go!

`/send` `GET` Display a helpful web form to test credentials and send messages manually. See below.

### Test form

Available for convenience of testing, uses the API defined above. Accessible in a web browser at `/send`:

<img src="https://i.postimg.cc/fb0wygJH/Screenshot-2021-02-01-at-20-38-17.png" width="400"/>

### CSV auth file

A very simple unencrypted key store, one line per key/secret pair. For example:

<img src="https://i.postimg.cc/nzyTZrNR/Screenshot-2021-02-01-at-21-00-29.png" width="600"/>

## Contributing

Contributions very welcome, please fork the repo and then submit PR changes back here.

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Contact

Project Link: [https://github.com/cdm/post-to-socials](https://github.com/cdm/post-to-socials)


