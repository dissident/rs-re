package tg

import (
	"html"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/dissident/rs-re/support"
)

var tagsPattern = regexp.MustCompile(`<[^>]+>`)

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
	message = html.UnescapeString(message)
	message = strings.Replace(message, "<br />", "\n", -1)
	message = strings.Replace(message, "<br>", "\n", -1)
	return tagsPattern.ReplaceAllString(message, "")
}

func (tg *Telegram) SendMessage(message string) error {
	msg := tgbotapi.NewMessage(tg.channel, formatMessage(message))
	_, err := tg.bot.Send(msg)
	return err
}
