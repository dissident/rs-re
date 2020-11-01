package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"strings"
	"github.com/mmcdole/gofeed"
)

func TestMain(t *testing.T) {
	feedData := `<rss version="2.0">
	<channel>
	<webMaster>example@site.com (Example Name)</webMaster>
	<item>
		<title>Foo</title>
	</item>
	<item>
		<title>Bar</title>
	</item>
	</channel>
	</rss>`

	fp := gofeed.NewParser()
	feed, _ := fp.Parse(strings.NewReader(feedData))
	memTitles := []string{"Bar", "Baz"}
	pushNewItems(&memTitles, feed)
	newTitle := feed.Items[0].Title
	assert.Equal(t, newTitle, "Foo")
	assert.Equal(t, memTitles, []string{"Foo", "Bar"})
}
