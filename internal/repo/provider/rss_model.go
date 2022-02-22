package provider

import "github.com/mmcdole/gofeed"

type feedItem gofeed.Item

func (fi *feedItem) GetTitle() string {
	return fi.Title
}

func (fi *feedItem) GetLink() string {
	return fi.Link
}
