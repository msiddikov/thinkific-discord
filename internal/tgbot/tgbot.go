package tgbot

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var chanId int64

func Start() {
	if bot != nil {
		return
	}
	chanId, _ = strconv.ParseInt(os.Getenv("TG_CHAN_ID"), 0, 64)
	b, err := tgbotapi.NewBotAPI(os.Getenv("TG_API"))
	bot = b
	if err != nil {
		panic("unable to connec to tg bot")
	}
	fmt.Println(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))
	bot.Debug = false
	go updates(bot)
}

func SendString(s string) {
	msg := tgbotapi.NewMessage(chanId, escape(s))
	msg.ParseMode = "MarkdownV2"
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
}

func updates(bot *tgbotapi.BotAPI) {
	fmt.Println("Getting tg updates...")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Your chat id is: %v", update.Message.Chat.ID))
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func escape(s string) (r string) {
	var escapeChars = []string{
		"\\", "_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!",
	}

	// // unEscape the string
	// for _, v := range escapeChars {
	// 	s = strings.Replace(s, "\\"+v, v, -1)
	// }

	// escape the string
	for _, v := range escapeChars {
		s = strings.Replace(s, v, "\\"+v, -1)
	}

	return s
}
