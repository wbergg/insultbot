package irc

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"net/textproto"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wbergg/insultbot/config"
	parser "gopkg.in/sorcix/irc.v2"
)

type Bot struct {
	Server     string
	Nick       string
	User       string
	Channel    string
	Pass       string
	conn       net.Conn
	InsultData []string
	ComplData  []string
}

func (bot *Bot) Connect() (net.Conn, error) {
	conn, err := net.Dial("tcp", bot.Server)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %w", err)
	}
	bot.conn = conn
	log.Printf("Connected to IRC server %s (%s) \n", bot.Server, bot.conn.RemoteAddr())

	reader := bufio.NewReader(bot.conn)
	tp := textproto.NewReader(reader)
	line, err := tp.ReadLine()
	if err != nil {
		return nil, fmt.Errorf("initial read failed: %w", err)
	}

	message := parser.ParseMessage(line)
	if message.Command == "PING" {
		bot.Send(fmt.Sprintf("PONG %s", message.Trailing()))
	}

	time.Sleep(2 * time.Second)

	if bot.Pass != "" {
		bot.Send(fmt.Sprintf("PASS %s", bot.Pass))
	}
	bot.Send(fmt.Sprintf("USER %s 8 * :%s", bot.User, bot.User))
	bot.Send(fmt.Sprintf("NICK %s", bot.Nick))
	bot.Send(fmt.Sprintf("JOIN %s", bot.Channel))

	return bot.conn, nil
}

func Run(cfg string, debugTelegram bool, debugStdout bool, telegramTest bool) error {
	// Load config
	config, err := config.LoadConfig(cfg)
	if err != nil {
		log.Error(err)
		panic("Could not load config, check config/config.json")
	}

	bot := &Bot{
		Server:  config.IRC.Server,
		Nick:    config.IRC.Nick,
		User:    config.IRC.User,
		Channel: config.IRC.Channel,
		Pass:    config.IRC.Password,
	}

	conn, err := bot.Connect()
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	// Load insults and compliments
	data, err := os.ReadFile("files/insults.txt")
	if err != nil {
		log.Fatalf("Error reading insults: %v", err)
	}
	bot.InsultData = strings.Split(string(data), "\n")

	data, err = os.ReadFile("files/compliments.txt")
	if err != nil {
		log.Fatalf("Error reading compliments: %v", err)
	}
	bot.ComplData = strings.Split(string(data), "\n")

	for {
		line, err := tp.ReadLine()
		if err != nil {
			log.Printf("Read error: %v", err)
			continue
		}

		message := parser.ParseMessage(line)
		if debugStdout {
			fmt.Printf("-> %s\n", line)
		}

		switch message.Command {
		case "PING":
			bot.Send(fmt.Sprintf("PONG %s", message.Trailing()))

		// Server welcome message, send join channel
		case "001":
			bot.Send(fmt.Sprintf("JOIN %s", bot.Channel))

		case "PRIVMSG":
			PrivMsg(bot, message)
		}
	}
}

func PrivMsg(bot *Bot, msg *parser.Message) {
	params := msg.Params
	if len(params) < 2 {
		return
	}

	//target := params[0]
	text := params[1]
	words := strings.Fields(text)
	if len(words) == 0 {
		return
	}

	cmd := words[0]
	args := strings.Join(words[1:], " ")
	//nick := msg.Prefix.Name

	switch cmd {
	case "!" + bot.Nick, bot.Nick:
		bot.Send(fmt.Sprintf("PRIVMSG %s :You said my name?", bot.Channel))

	case "!insult", ".insult":
		if args != "" {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			insult := bot.InsultData[r.Intn(len(bot.InsultData))]
			msg := strings.Replace(insult, "%s", args, -1)
			msg = strings.ReplaceAll(msg, "  ", " ")
			bot.Send(fmt.Sprintf("PRIVMSG %s :%s", bot.Channel, msg))
		}

	case "!addinsult", ".addinsult":
		if args != "" {
			err := bot.WriteFile("files/insults.txt", args)
			if err == nil {
				bot.InsultData = append(bot.InsultData, args)
				bot.Send(fmt.Sprintf("PRIVMSG %s :Added insult: %s", bot.Channel, args))
			}
		}

	case "!compliment", ".compliment":
		if args != "" {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			compl := bot.ComplData[r.Intn(len(bot.ComplData))]
			msg := strings.Replace(compl, "%s", args, -1)
			msg = strings.ReplaceAll(msg, "  ", " ")
			bot.Send(fmt.Sprintf("PRIVMSG %s :%s", bot.Channel, msg))
		}

	case "!addcompliment", ".addcompliment":
		if args != "" {
			err := bot.WriteFile("files/compliments.txt", args)
			if err == nil {
				bot.ComplData = append(bot.ComplData, args)
				bot.Send(fmt.Sprintf("PRIVMSG %s :Added compliment: %s", bot.Channel, args))
			}
		}

	case "!help", ".help":
		message := []string{
			"Welcome to insultbot!",
			"-----",
			"Use !insult <nick> to insult.",
			"Use !addinsult <insult> (with %s as a nickname placeholder).",
			"Use !compliment <nick> to compliment.",
			"Use !addcompliment <compliment> (with %s).",
		}
		for _, h := range message {
			bot.Send(fmt.Sprintf("PRIVMSG %s :%s", bot.Channel, h))
		}
	}
}

func (bot *Bot) Send(command string) {
	fmt.Fprintf(bot.conn, "%s\r\n", command)
}

func (bot *Bot) WriteFile(filename string, text string) error {

	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	if _, err := f.WriteString(text + "\n"); err != nil {
		log.Println(err)
	}
	f.Close()
	return nil
}
