package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dissident/rs-re/support"
	"github.com/dissident/rs-re/tg"
	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Env struct {
	feedURL         string
	mongoURL        string
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(env.mongoURL))
	failOnError(err, "Failed connect to mongo")
	collection := client.Database("rs-re").Collection("upworks")

	for {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(env.feedURL)
		failOnError(err, "Failed to parse env.feedURL")

		log.Printf("Fetching RSS feed...")

		newTitles := pushNewItems(memTitles, feed, bot, *env, collection, ctx)
		memTitles = newTitles
		support.PrintMemUsage()
		sleepDuration, err := time.ParseDuration(env.teakInterval)
		failOnError(err, "Failed to parse TEAK_INTERVAL ENV as a time.Duration")
		time.Sleep(sleepDuration)
	}
}

func initEnvironment() *Env {
	feedURL := os.Getenv("FEED_URL")
	mongoURL := os.Getenv("MONGO_URL")
	teakInterval := os.Getenv("TEAK_INTERVAL")

	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	telegramChannel, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	failOnError(err, "Failed to parse CHAT_ID ENV")

	return &Env{feedURL, mongoURL, telegramToken, telegramChannel, teakInterval}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func pushNewItems(memTitles []string, feed *gofeed.Feed, bot *tgbotapi.BotAPI, env Env, collection *mongo.Collection, ctx context.Context) []string {
	newTitles := []string{}
	for _, item := range feed.Items {
		isPresent := memTitlesContains(memTitles, item.Title)
		if !isPresent {
			log.Printf(item.Title)
			log.Printf(item.Link)
			tg.SendMessage(item.Title, bot, env.telegramChannel)
			collection.InsertOne(ctx, bson.D{{"title", item.Title}, {"body", item.Content}})
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
