# insultbot

In memory of Riverbot

This is a project to revive the old bot called "riverbot" which was a bot written in Python(?). It had a simple !insult function where you could insult people with in an IRC-channel. The source code from riverbot is long lost, however, I had a backup of insults.txt. I've taken the liberty to rewrite the entire bot in GO. 

In version 2.0 support for Telegram is added since IRC is slowly dying and being more and more phased out. With the new version, a config file is added to control options and enable/disable different services and settings.

## Setup
Create a config-file called "config.json" in the config/ dir.

```
{
    "Telegram": {
        "enabled": "true",
		"tgAPIkey": "xxx",
		"tgChannel": "xxx"
	},
    "IRC": {
        "enabled": "true",
        "server": "xxx:port",
		"nick": "xxx",
		"user": "xxx",
		"channel": "#xxx",
		"password": ""
    }
}
```

Create a directory called files/ and create a file called insults.txt and compliments.txt in that dir.

## Running

```
go run insultbot.go
```

Or with arguments:

```
  -config-file string
    	Absolute path for config-file (default "./config/config.json")
  -stdout
    	Turns on stdout rather than sending to telegram
  -telegram-debug
    	Turns on debug for telegram
  -telegram-test
    	Sends a test message to specified telegram channel
```

### Debug
Set debug or stdout to true in an argument to insultbot.go to receive debug and stdout.
