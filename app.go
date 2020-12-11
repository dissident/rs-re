package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dissident/rs-re/support"
	"github.com/dissident/rs-re/tg"
	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Env struct {
	feedURL         string
	telegramToken   string
	telegramChannel int64
	teakInterval    string
}

func main() {
	godotenv.Load()
	env := initEnvironment()
	bot, err := tg.InitBot(env.telegramToken)
	failOnError(err, "Failed tgbotapi.NewBotAPI initialize")

	memTitles := []string{}

	for {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(env.feedURL)
		failOnError(err, "Failed to parse env.feedURL")

		log.Printf("Fetching RSS feed...")

		newTitles := pushNewItems(memTitles, feed, bot, *env)
		memTitles = newTitles
		support.PrintMemUsage()
		sleepDuration, err := time.ParseDuration(env.teakInterval)
		failOnError(err, "Failed to parse TEAK_INTERVAL ENV as a time.Duration")
		time.Sleep(sleepDuration)
	}
}

func initEnvironment() *Env {
	feedURL := os.Getenv("FEED_URL")
	teakInterval := os.Getenv("TEAK_INTERVAL")

	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	telegramChannel, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	failOnError(err, "Failed to parse CHAT_ID ENV")

	return &Env{feedURL, telegramToken, telegramChannel, teakInterval}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func pushNewItems(memTitles []string, feed *gofeed.Feed, bot *tgbotapi.BotAPI, env Env) []string {
	newTitles := []string{}
	for _, item := range feed.Items {
		isPresent := memTitlesContains(memTitles, item.Title)
		if !isPresent {
			log.Printf(item.Title)
			log.Printf(item.Link)
			tg.SendMessage(item.Title, bot, env.telegramChannel)
			err := tg.SendMessage(item.Content, bot, env.telegramChannel)
			if err != nil {
				tg.SendMessage("Content body can't be sended. Use a link >", bot, env.telegramChannel)
				tg.SendMessage(item.Link, bot, env.telegramChannel)
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
