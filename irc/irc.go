package irc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os"
	"regexp"
	"time"

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

func (bot *Bot) Send(command string) {
	fmt.Fprintf(bot.conn, "%s\r\n", command)
}

func NewBot(server string, Nick string, User string, Channel string, pass string) *Bot {
	return &Bot{
		Server:  server,
		Nick:    Nick,
		User:    User,
		Channel: Channel,
		Pass:    pass,
		conn:    nil,
	}
}

func (bot *Bot) Connect() (conn net.Conn, err error) {

	bot.conn, err = net.Dial("tcp", bot.Server)
	if err != nil {
		log.Fatal("Unable to connect to IRC server ", err)
	}
	log.Printf("Connected to IRC server %s (%s) \n", bot.Server, bot.conn.RemoteAddr())

	reader := bufio.NewReader(bot.conn)
	tp := textproto.NewReader(reader)
	line, err := tp.ReadLine()
	message := parser.ParseMessage(line)

	if message.Command == "PING" {
		r := regexp.MustCompile("[^0-9.]")
		s := r.ReplaceAllString(line, "")
		//fmt.Println(fmt.Sprint("PONG ", s))
		bot.Send(fmt.Sprint("PONG ", s))
	}

	time.Sleep(2 * time.Second)

	if len(bot.Pass) > 0 {
		bot.Send(fmt.Sprintf("PASS %s \r\n", bot.Pass))
	}
	bot.Send(fmt.Sprintf("USER %s 8 * :%s", bot.User, bot.User))
	bot.Send(fmt.Sprintf("NICK %s\r\n", bot.Nick))
	bot.Send(fmt.Sprintf("JOIN %s\r\n", bot.Channel))

	return bot.conn, err
}

func (bot *Bot) ReadFile(filename string) ([]string, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	fileData := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		fileData = append(fileData, line)
	}
	f.Close()
	return fileData, err
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
