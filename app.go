package main

import (
	"log"
	"time"

	"github.com/dissident/rs-re/support"
	"github.com/dissident/rs-re/tg"
	"github.com/dissident/rs-re/db"
	"github.com/dissident/rs-re/environment"
	"github.com/mmcdole/gofeed"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	env := environment.InitEnvironment()
	bot, err := tg.InitBot(env.TelegramToken)
	support.FailOnError(err, "Failed tgbotapi.NewBotAPI initialize")

	memTitles := []string{}

	database := db.InitDb(env.MongoURL, env.Db, env.Collection)

	for {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(env.FeedURL)
		support.FailOnError(err, "Failed to parse Feed URL")

		log.Printf("Fetching RSS feed...")

		newTitles := pushNewItems(memTitles, feed, bot, env, database)
		memTitles = newTitles
		support.PrintMemUsage()
		sleepDuration, err := time.ParseDuration(env.TeakInterval)
		support.FailOnError(err, "Failed to parse TEAK_INTERVAL ENV as a time.Duration")
		time.Sleep(sleepDuration)
	}
}

func pushNewItems(memTitles []string, feed *gofeed.Feed, bot *tgbotapi.BotAPI, env environment.Env, database db.DB) []string {
	newTitles := []string{}
	for _, item := range feed.Items {
		isPresent := memTitlesContains(memTitles, item.Title)
		if !isPresent {
			log.Printf(item.Title)
			log.Printf(item.Link)
			tg.SendMessage(item.Title, bot, env.TelegramChannel)
			database.Insert(item.Title, item.Content)
			err := tg.SendMessage(item.Content, bot, env.TelegramChannel)
			if err != nil {
				tg.SendMessage("Content body can't be sended. Use a link >", bot, env.TelegramChannel)
				tg.SendMessage(item.Link, bot, env.TelegramChannel)
			}
		}
		newTitles = append(newTitles, item.Title)
	}
	return newTitles
}

func memTitlesContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
