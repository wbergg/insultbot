package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/textproto"
	"regexp"
	"strings"
	"time"

	"github.com/wbergg/insult-bot/irc"

	parser "gopkg.in/sorcix/irc.v2"
)

func main() {

	// IRC server settings
	bot := irc.NewBot(
		"se.quakenet.org:6667", // Server:port
		"insultbot-dev",        // Nick
		"insultbot-dev",        // User
		"#wberg",               // Channel
		"",                     // Channel password
	)

	// Debug output to stdout
	debug := true

	conn, _ := bot.Connect()
	defer conn.Close()
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	// Preload insults
	insultData, err := bot.ReadFile("files/insults.txt")
	if err != nil {
		panic(err)
	}
	bot.InsultData = insultData

	// Preload compliments
	complData, err := bot.ReadFile("files/compliments.txt")
	if err != nil {
		panic(err)
	}
	bot.ComplData = complData

	for {
		line, _ := tp.ReadLine()
		message := parser.ParseMessage(line)

		if debug {
			fmt.Printf("%v \n", message)
		}

		// Ping function, designed to handle Quakenet pings
		if message.Command == "PING" {
			r := regexp.MustCompile("[^0-9.]")
			s := r.ReplaceAllString(line, "")
			fmt.Println(fmt.Sprint("PONG ", s))
			bot.Send(fmt.Sprint("PONG ", s))
		}

		if message.Command == "PRIVMSG" {
			if message.Params[0] == bot.Nick {
				// temp test for joining channel on query
				bot.Send(fmt.Sprintf("JOIN %s\r\n", bot.Channel))
				bot.Send(fmt.Sprintf("PRIVMSG %s hello\r\n", bot.Channel))
			}
		}
		// Simple reply upon bot nick
		if message.Command == "PRIVMSG" {
			if message.Params[1] == "!"+bot.Nick || message.Params[1] == bot.Nick {
				// test message
				bot.Send(fmt.Sprintf("PRIVMSG %s :You said my name?\r\n", bot.Channel))
			}
		}
		// !insult function
		if message.Command == "PRIVMSG" {
			randTime := rand.NewSource(time.Now().UnixNano())
			r := rand.New(randTime)
			randomIndex := r.Intn(len(bot.InsultData))
			insult := bot.InsultData[randomIndex]

			// Split incoming message
			s := strings.Split(message.Params[1], " ")

			// Check first word in message
			if s[0] == "!insult" || s[0] == ".insult" {
				// Check lenght of message and do stuff
				if len(s) > 1 {
					// Load insults
					insultData, err := bot.ReadFile("files/insults.txt")
					if err != nil {
						panic(err)
					}
					bot.InsultData = insultData

					var nickname string
					if s[0] == "!insult" {
						nickname = strings.Replace(message.Params[1], "!insult ", "", 1)
					}
					if s[0] == ".insult" {
						nickname = strings.Replace(message.Params[1], ".insult ", "", 1)
					}
					message := strings.Replace(insult, "%s", nickname, -1)
					bot.Send(fmt.Sprintf("PRIVMSG %s :%s\r\n", bot.Channel, message))
				}
			}
		}
		// !addinsult function
		if message.Command == "PRIVMSG" {
			// Split incoming message
			s := strings.Split(message.Params[1], " ")

			// Check first word in message
			if s[0] == "!addinsult" || s[0] == ".addinsult" {
				// Check lenght of message and do stuff
				if len(s) > 1 {
					var addinsult string
					if s[0] == "!addinsult" {
						addinsult = strings.Replace(message.Params[1], "!addinsult ", "", 1)
					}
					if s[0] == ".addinsult" {
						addinsult = strings.Replace(message.Params[1], ".addinsult ", "", 1)
					}

					err := bot.WriteFile("files/insults.txt", addinsult)
					if err != nil {
						panic(err)
					}
					bot.Send(fmt.Sprintf("PRIVMSG %s :Added insult: %s\r\n", bot.Channel, addinsult))
				}
			}
		}
		// !compliment function
		if message.Command == "PRIVMSG" {
			randTime := rand.NewSource(time.Now().UnixNano())
			r := rand.New(randTime)
			randomIndex := r.Intn(len(bot.ComplData))
			compliment := bot.ComplData[randomIndex]

			// Split incoming message
			s := strings.Split(message.Params[1], " ")

			// Check first word in message
			if s[0] == "!compliment" || s[0] == ".compliment" {
				// Check lenght of message and do stuff
				if len(s) > 1 {
					// Load compliments
					complData, err := bot.ReadFile("files/compliments.txt")
					if err != nil {
						panic(err)
					}
					bot.ComplData = complData

					var nickname string
					if s[0] == "!compliment" {
						nickname = strings.Replace(message.Params[1], "!compliment ", "", 1)
					}
					if s[0] == ".compliment" {
						nickname = strings.Replace(message.Params[1], ".compliment ", "", 1)
					}
					message := strings.Replace(compliment, "%s", nickname, -1)
					bot.Send(fmt.Sprintf("PRIVMSG %s :%s\r\n", bot.Channel, message))
				}
			}
		}
		// !addcompliment function
		if message.Command == "PRIVMSG" {
			// Split incoming message
			s := strings.Split(message.Params[1], " ")

			// Check first word in message
			if s[0] == "!addcompliment" || s[0] == ".addcompliment" {
				// Check lenght of message and do stuff
				if len(s) > 1 {
					var addcompl string
					if s[0] == "!addcompliment" {
						addcompl = strings.Replace(message.Params[1], "!addcompliment ", "", 1)
					}
					if s[0] == ".addcompliment" {
						addcompl = strings.Replace(message.Params[1], ".addcompliment ", "", 1)
					}
					err := bot.WriteFile("files/compliments.txt", addcompl)
					if err != nil {
						panic(err)
					}
					bot.Send(fmt.Sprintf("PRIVMSG %s :Added compliment: %s\r\n", bot.Channel, addcompl))
				}
			}
		}
		// !help function
		if message.Command == "PRIVMSG" {
			// Print help
			if message.Params[1] == "!help" {
				bot.Send(fmt.Sprintf("PRIVMSG %s :Welcome to insultbot 1.1!\r\n", bot.Channel))
				bot.Send(fmt.Sprintf("PRIVMSG %s :-----\r\n", bot.Channel))
				bot.Send(fmt.Sprintf("PRIVMSG %s :To insult someone, use !insult <nick>\r\n", bot.Channel))
				bot.Send(fmt.Sprintf("PRIVMSG %s :To add your own insult, use !addinsult <insult>, use %%s as a placeholder for <nick>.\r\n", bot.Channel))
				bot.Send(fmt.Sprintf("PRIVMSG %s :For example: !addinsult %%s smells bad.\r\n", bot.Channel))
				bot.Send(fmt.Sprintf("PRIVMSG %s :-----\r\n", bot.Channel))
			}
		}
	}
}
