package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	memTitles := []string{}

	godotenv.Load()
	feedURL := os.Getenv("FEED_URL")

	fp := gofeed.NewParser()

	for {
		feed, _ := fp.ParseURL(feedURL)
		log.Printf("Fetching RSS feed...")

		pushNewItems(&memTitles, feed)
		time.Sleep(5000 * time.Millisecond)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendMessage(message string) {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	telegramChannel, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	failOnError(err, "Failed to parse telegramChannel ENV")
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	failOnError(err, "Failed tgbotapi.NewBotAPI initialize")
	msg := tgbotapi.NewMessage(telegramChannel, message)
	bot.Send(msg)
}

func pushNewItems(memTitles *[]string, feed *gofeed.Feed) {
	newTitles := []string{}
	for _, item := range feed.Items {
		isPresent := memTitlesContains(*memTitles, item.Title)
		if !isPresent {
			log.Printf(strings.Join(*memTitles, ","))
			log.Printf(item.Title)
			log.Printf(item.Link)
			// log.Printf(item.Content)
			sendMessage(item.Title)
			sendMessage(item.Link)
			sendMessage(item.Content)
		}
		newTitles = append(newTitles, item.Title)
	}
	*memTitles = newTitles
}

func memTitlesContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
