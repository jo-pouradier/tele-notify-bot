package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	// "strings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	metrics "github.com/jo-pouradier/homelab-bot/metrics"
	"github.com/joho/godotenv"
)

func ReadHostTemp() {
	for {
		log.Println("test")
	}
}

func init() {
	godotenv.Load()
	id, err := strconv.Atoi(os.Getenv("ADMIN_USER"))

	if err != nil {
		return
	}
	authorizedUsers = id
}

var authorizedUsers int

func main() {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOCKEN"))
	if err != nil {
		log.Panic("Telegram key not found")
	}
	if err != nil {
		log.Panic("Error with ADMIN env ID")
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			processMessage(update, bot)
		}
	}

}

func processMessage(u tgbotapi.Update, b *tgbotapi.BotAPI) {
	// check user ID
	fromID := int(u.SentFrom().ID)
	fmt.Printf("authorizedUser: %v", authorizedUsers)
	if fromID != authorizedUsers {
		fmt.Printf("Not authorized user: %s, %v\n", u.Message.From.UserName, u.Message.From.ID)
		return
	}

	// check if command
	if !u.Message.IsCommand() {
		return
	}

	cpu, _ := metrics.GetCPU1()
	mem, _ := metrics.GetMEM1()
	// textResponse := commands[u.Message.Command()](u)
	finaleResponse := "METRICS: \n\t" +
		"CPU: " + fmt.Sprintf("%.2f", cpu) + "\n" +
		"MEM: " + fmt.Sprintf("%.2f", mem)

	msg := tgbotapi.NewMessage(u.Message.From.ID, finaleResponse)
	msg.ReplyToMessageID = u.Message.MessageID

	b.Send(msg)

}

func processMetrics(u tgbotapi.Update) string {
	return "metrics: ..."
}

var commands = map[string]func(tgbotapi.Update) string{
	"metrics": processMetrics,
}
