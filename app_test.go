package main

import (
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestPushNewItems(t *testing.T) {
	feedData := `<rss version="2.0">
	<channel>
	<webMaster>example@site.com (Example Name)</webMaster>
	<item>
		<title>Foo</title>
	</item>
	<item>
		<title>Bar</title>
	</item>
	</channel>
	</rss>`

	godotenv.Load()
	env := initEnvironment()
	bot, err := tgbotapi.NewBotAPI(env.telegramToken)
	failOnError(err, "Failed tgbotapi.NewBotAPI initialize")

	fp := gofeed.NewParser()
	feed, _ := fp.Parse(strings.NewReader(feedData))
	memTitles := []string{"Bar", "Baz"}
	pushNewItems(memTitles, feed, bot, env)
	newTitle := feed.Items[0].Title
	assert.Equal(t, newTitle, "Foo")
	assert.Equal(t, memTitles, []string{"Foo", "Bar"})
}

func TestFormatMessage(t *testing.T) {
	message := `
		<b>Title</b><br />
		1. First <br />
		2. Second
	`

	result := `
	<b>Title</b>
	> 1. First
	> 2. Second
	`
	assert.Equal(t, formatMessage(message), result)
}
