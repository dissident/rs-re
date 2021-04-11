package main

import (
	"log"
	"time"

	"github.com/dissident/rs-re/support"
	"github.com/dissident/rs-re/tg"
	"github.com/dissident/rs-re/db"
	"github.com/dissident/rs-re/environment"
	"github.com/mmcdole/gofeed"
)

func main() {
	env := environment.InitEnvironment()
	telega := tg.Init(env.TelegramToken, env.TelegramChannel)
	memTitles := []string{}

	database := db.InitDb(env.MongoURL, env.Db, env.Collection)

	for {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(env.FeedURL)
		support.FailOnError(err, "Failed to parse Feed URL")

		log.Printf("Fetching RSS feed...")

		newTitles := pushNewItems(memTitles, feed, telega, env, database)
		memTitles = newTitles
		support.PrintMemUsage()
		sleepDuration, err := time.ParseDuration(env.TeakInterval)
		support.FailOnError(err, "Failed to parse TEAK_INTERVAL ENV as a time.Duration")
		time.Sleep(sleepDuration)
	}
}

func pushNewItems(memTitles []string, feed *gofeed.Feed, telega tg.Telegram, env environment.Env, database db.DB) []string {
	newTitles := []string{}
	for _, item := range feed.Items {
		isPresent := memTitlesContains(memTitles, item.Title)
		if !isPresent {
			logItem(item)
			database.Insert(item.Title, item.Content)
			telgramItem(item, telega)
		}
		newTitles = append(newTitles, item.Title)
	}
	return newTitles
}

func telgramItem(item *gofeed.Item, telega tg.Telegram) {
	telega.SendMessage(item.Title)
	err := telega.SendMessage(item.Content)
	if err != nil {
		telega.SendMessage("Content body can't be sended. Use a link >")
		telega.SendMessage(item.Link)
	}
}

func logItem(item *gofeed.Item) {
	log.Printf(item.Title)
	log.Printf(item.Link)
}

func memTitlesContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
