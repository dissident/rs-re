package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Env struct {
	feedURL         string
	telegramToken   string
	telegramChannel int64
}

func main() {
	godotenv.Load()
	env := initEnvironment()
	bot, err := tgbotapi.NewBotAPI(env.telegramToken)
	failOnError(err, "Failed tgbotapi.NewBotAPI initialize")

	memTitles := []string{}

	for {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(env.feedURL)
		failOnError(err, "Failed to parse env.feedURL")

		log.Printf("Fetching RSS feed...")

		newTitles := pushNewItems(memTitles, feed, bot, env)
		memTitles = newTitles
		PrintMemUsage()
		time.Sleep(5000 * time.Millisecond)
	}
}

func initEnvironment() Env {
	feedURL := os.Getenv("FEED_URL")
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	telegramChannel, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	failOnError(err, "Failed to parse telegramChannel ENV")

	return Env{feedURL, telegramToken, telegramChannel}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendMessage(message string, bot *tgbotapi.BotAPI, env Env) {
	msg := tgbotapi.NewMessage(env.telegramChannel, message)
	bot.Send(msg)
}

func pushNewItems(memTitles []string, feed *gofeed.Feed, bot *tgbotapi.BotAPI, env Env) []string {
	newTitles := []string{}
	for _, item := range feed.Items {
		isPresent := memTitlesContains(memTitles, item.Title)
		if !isPresent {
			log.Printf(item.Title)
			log.Printf(item.Link)
			sendMessage(item.Title, bot, env)
			sendMessage(item.Content, bot, env)
		}
		newTitles = append(newTitles, item.Title)
	}
	return newTitles
}

func formatMessage(message string) string {
	return strings.Replace(message, "<br />", ">", -1)
}

func memTitlesContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
