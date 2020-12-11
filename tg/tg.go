package tg

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func InitBot(telegramToken string) (*tgbotapi.BotAPI, error) {
	return tgbotapi.NewBotAPI(telegramToken)
}

func formatMessage(message string) string {
	return strings.Replace(message, "<br />", "\n", -1)
}

func SendMessage(message string, bot *tgbotapi.BotAPI, telegramChannel int64) error {
	msg := tgbotapi.NewMessage(telegramChannel, formatMessage(message))
	msg.ParseMode = "HTML"
	_, err := bot.Send(msg)
	return err
}
