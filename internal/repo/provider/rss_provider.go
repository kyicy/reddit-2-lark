package provider

import (
	"context"
	"strings"

	"github.com/kyicy/rss-2-lark/internal/platform"
	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
)

type RssProvider struct {
	name   string
	src    string
	logger *zap.SugaredLogger
	conf   *platform.Config
}

func NewRssProvider(
	name, src string,
	conf *platform.Config,
) *RssProvider {
	z, _ := zap.NewDevelopment()
	return &RssProvider{
		name:   name,
		src:    src,
		logger: z.Sugar().Named(name),
		conf:   conf,
	}
}

func (rp *RssProvider) GetName() string {
	return rp.name
}

func (rp *RssProvider) GetTopPosts(
	ctx context.Context,
) (items []Broadcastable, err error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(rp.src, ctx)
	if err != nil {
		return
	}
	items = make([]Broadcastable, 0)
	for _, item := range feed.Items {
		var matched bool
		for _, keyword := range rp.conf.Keywords {
			if strings.Contains(item.Title, keyword) {
				matched = true
				break
			}
		}
		if matched {
			items = append(items, (*feedItem)(item))
		}
	}
	return
}
