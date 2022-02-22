package provider

import (
	"time"

	"github.com/mmcdole/gofeed"
)

type feedItem gofeed.Item

func (fi *feedItem) GetTitle() string {
	return fi.Title
}

func (fi *feedItem) GetLink() string {
	return fi.Link
}

func (fi *feedItem) GetPubDate() time.Time {
	return *fi.PublishedParsed
}
