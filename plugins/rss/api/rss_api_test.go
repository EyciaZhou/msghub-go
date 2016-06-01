package rss_api

import (
	"testing"
)

func TestGetNew(t *testing.T) {
	controller := RssController{"https://eycia.me/blog/index.php/feed/"}
	feed, err := controller.GetFeed()
	if err != nil {
		panic(err)
	}

	t.Error(feed.String())
}
