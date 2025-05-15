package tele

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wbergg/insultbot/config"
	"github.com/wbergg/telegram"
)

func Run(cfg string, debugTelegram bool, debugStdout bool, telegramTest bool) error {

	// Load config
	config, err := config.LoadConfig(cfg)
	if err != nil {
		log.Error(err)
		panic("Could not load config, check config/config.json")
	}

	channel, err := strconv.ParseInt(config.Telegram.TgChannel, 10, 64)
	if err != nil {
		log.Error(err)
		panic("Could not convert Telegram channel to int64")
	}

	// Initiate telegram
	tg := telegram.New(config.Telegram.TgAPIKey, channel, debugTelegram, debugStdout)
	tg.Init(debugTelegram)

	if telegramTest {
		tg.SendM("DEBUG: insultbot test message")
		os.Exit(0)
	}

	// Preload insults and compliments from file
	insultData, err := os.ReadFile("files/insults.txt")
	if err != nil {
		panic(err)
	}
	complimentData, err := os.ReadFile("files/compliments.txt")
	if err != nil {
		panic(err)
	}
	// Split into lines
	insults := strings.Split(string(insultData), "\n")
	compliments := strings.Split(string(complimentData), "\n")

	// Read messages from Telegram
	updates, err := tg.ReadM()
	if err != nil {
		log.Error(err)
		panic("Cant read from Telegram")
	}

	// Loop
	for update := range updates {
		if update.Message == nil { // ignore non-message updates
			continue
		}

		// Debug
		if debugStdout {
			log.Infof("Received message from chat %d [%s]: %s", update.Message.Chat.ID, update.Message.Chat.Type, update.Message.Text)
		}

		if update.Message.IsCommand() {
			// Create switch to search for commands
			switch strings.ToLower(update.Message.Command()) {

			// Insult case
			case "insult":
				message := update.Message.CommandArguments()

				if message == "" {
					// If nothings wa inpuuted, return calling userid
					message = update.Message.From.UserName
					if message == "" {
						message = update.Message.From.FirstName
					}
				}

				// Seed RNG and select a random line
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				randomIndex := r.Intn(len(insults))
				insult := insults[randomIndex]

				// Replace and send
				reply := strings.Replace(insult, "%s", message, -1)
				tg.SendTo(update.Message.Chat.ID, reply)

			// Add insult case
			case "addinsult":
				message := update.Message.CommandArguments()
				if message == "" {
					tg.SendM("Message is empty, not adding.")
					break
				}

				// Open file
				f, err := os.OpenFile("files/insults.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Error("Error opening file:", err)
					panic(err)
				}

				// Write the content to file
				if _, err := f.WriteString(message + "\n"); err != nil {
					log.Error("Error writing to file:", err)
					tg.SendTo(update.Message.Chat.ID, "Error")
					panic(err)
				} else {
					successfull := "Successfully added insult: " + message
					tg.SendTo(update.Message.Chat.ID, successfull)
				}
				f.Close()

			// Compliment case
			case "compliment":
				message := update.Message.CommandArguments()

				if message == "" {
					// If nothings wa inputed, return calling userid
					message = update.Message.From.UserName
					if message == "" {
						message = update.Message.From.FirstName
					}
				}

				// Seed RNG and select a random line
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				randomIndex := r.Intn(len(compliments))
				compliment := compliments[randomIndex]

				// Replace and send
				reply := strings.Replace(compliment, "%s", message, -1)
				tg.SendTo(update.Message.Chat.ID, reply)

			// Add compliment case
			case "addcompliment":
				message := update.Message.CommandArguments()
				if message == "" {
					tg.SendM("Message is empty, not adding.")
					break
				}

				// Open file
				f, err := os.OpenFile("files/compliments.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Error("Error opening file:", err)
					panic(err)
				}

				// Write the content to file
				if _, err := f.WriteString(message + "\n"); err != nil {
					log.Error("Error writing to file:", err)
					tg.SendTo(update.Message.Chat.ID, "Error")
					panic(err)
				} else {
					successfull := "Successfully added insult: " + message
					tg.SendTo(update.Message.Chat.ID, successfull)
				}
				f.Close()

			case "help":
				helpm := `Insultbot 2.0

				/insult <nick> to insult someone
				/compliment <nick> to compliment someone

				Adding: Use %s in the add command as a placeholder for a users nickname. 

				For example:
				/addinsult %s does not drink beer
				/addcompliment %s does drink beer`
				tg.SendM(helpm)

			default:
				// Unknown command
				tg.SendM("")
			}
		}
	}

	return err
}
