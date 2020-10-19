package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"github.com/mmcdole/gofeed"
	"time"
)

func main() {
	var memTitles []string

	godotenv.Load()
	feedUrl := os.Getenv("FEED_URL")

	fp := gofeed.NewParser()

	for {
		feed, _ := fp.ParseURL(feedUrl)
		log.Printf("Fetching RSS feed...")

		pushNewItems(&memTitles, feed)
		time.Sleep(5000 * time.Millisecond)
	}
}

func pushNewItems(memTitles *[]string, feed *gofeed.Feed) {
	var newTitles []string
	for _, item := range feed.Items {
		isPresent := memTitlesContains(*memTitles, item.Title)
		if !isPresent {
			log.Printf(item.Title)
			log.Printf(item.Link)
			// log.Printf(item.Content)
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
