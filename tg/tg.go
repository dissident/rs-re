package tg

import (
	"strings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/dissident/rs-re/support"
)

type Telegram struct {
	bot *tgbotapi.BotAPI
	channel int64
}

func Init(token string, channel int64) Telegram {
	bot, err := tgbotapi.NewBotAPI(token)
	support.FailOnError(err, "Failed tgbotapi.NewBotAPI initialize")

	return Telegram{ bot, channel }
}

func formatMessage(message string) string {
	return strings.Replace(message, "<br />", "\n", -1)
}

func (tg *Telegram) SendMessage(message string) error {
	msg := tgbotapi.NewMessage(tg.channel, formatMessage(message))
	msg.ParseMode = "HTML"
	_, err := tg.bot.Send(msg)
	return err
}
