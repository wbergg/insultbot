# insultbot

In memory of Riverbot

This is a project to revive the old bot called "riverbot" which was a bot written by Python(?). It had a simple !insult function were you could insult people with in an IRC-channel. The soruce code from riverbot is lost, however, I had a backup of insults.txt. I've taken the liberty to rewrite the entire bot in GO. 

## Setup
Edit the top part in insultbot.go with your server settings:

```
// IRC server settings
bot := irc.NewBot(
    "xxx:yyyy",    // Server:port
    "nick",        // Nick
    "user",        // User
    "#channel",    // Channel
    "",            // Channel password
)
```

Create a directory called files/ and create a file called insults.txt in that dir.

## Running

```
go run insultbot.go
```

### Debug
Set debug bool to true or false in insultbot.go to recive debug into to stdout

```
// Debug output to stdout
debug := true
```
