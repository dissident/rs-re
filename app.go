package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"github.com/mmcdole/gofeed"
)

func main() {
	godotenv.Load()
	feedUrl := os.Getenv("FEED_URL")
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(feedUrl)
	for _, item := range feed.Items {
		log.Printf(item.Title)
		log.Printf(item.Link)
		log.Printf(item.Content)
	}
}
