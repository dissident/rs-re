package environment

import (
	"os"
	"strconv"
	"github.com/joho/godotenv"

	"github.com/dissident/rs-re/support"
)

type Env struct {
	FeedURL         string
	MongoURL        string
	TelegramToken   string
	TelegramChannel int64
	TeakInterval    string
	Db              string
	Collection      string
}

func InitEnvironment() Env {
	godotenv.Load()
	FeedURL := os.Getenv("FEED_URL")
	MongoURL := os.Getenv("MONGO_URL")
	Db := os.Getenv("DB")
	Collection := os.Getenv("COLLECTION")
	TeakInterval := os.Getenv("TEAK_INTERVAL")

	TelegramToken := os.Getenv("TELEGRAM_TOKEN")
	TelegramChannel, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	support.FailOnError(err, "Failed to parse CHAT_ID ENV")

	return Env{
		FeedURL,
		MongoURL,
		TelegramToken,
		TelegramChannel,
		TeakInterval,
		Db,
		Collection,
	}
}
