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
	teakInterval    string
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
		sleepDuration, err := time.ParseDuration(env.teakInterval)
		failOnError(err, "Failed to parse TEAK_INTERVAL ENV as a time.Duration")
		time.Sleep(sleepDuration)
	}
}

func initEnvironment() Env {
	feedURL := os.Getenv("FEED_URL")
	teakInterval := os.Getenv("TEAK_INTERVAL")

	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	telegramChannel, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	failOnError(err, "Failed to parse CHAT_ID ENV")

	return Env{feedURL, telegramToken, telegramChannel, teakInterval}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendMessage(message string, bot *tgbotapi.BotAPI, env Env) error {
	msg := tgbotapi.NewMessage(env.telegramChannel, formatMessage(message))
	msg.ParseMode = "HTML"
	_, err := bot.Send(msg)
	return err
}

func pushNewItems(memTitles []string, feed *gofeed.Feed, bot *tgbotapi.BotAPI, env Env) []string {
	newTitles := []string{}
	for _, item := range feed.Items {
		isPresent := memTitlesContains(memTitles, item.Title)
		if !isPresent {
			log.Printf(item.Title)
			log.Printf(item.Link)
			sendMessage(item.Title, bot, env)
			err := sendMessage(item.Content, bot, env)
			if err != nil {
				sendMessage("Content body can't be sended. Use a link >", bot, env)
				sendMessage(item.Link, bot, env)
			}
		}
		newTitles = append(newTitles, item.Title)
	}
	return newTitles
}

func formatMessage(message string) string {
	return strings.Replace(message, "<br />", "\n", -1)
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
