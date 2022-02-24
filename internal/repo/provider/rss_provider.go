package provider

import (
	"context"
	"strings"

	"github.com/kyicy/rss-2-lark/internal/platform"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
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
		err = errors.Wrapf(err, "name: %s, url: %s", rp.name, rp.src)
		return
	}
	rp.logger.Debug("keyword from config", "list", rp.conf.Keywords)
	items = make([]Broadcastable, 0)
	for _, item := range feed.Items {
		for _, keyword := range rp.conf.Keywords {
			if strings.Contains(item.Title, keyword) ||
				strings.Contains(item.Description, keyword) {
				rp.logger.Debugw("item added", "item.title", item.Title)
				items = append(items, (*feedItem)(item))
				break
			}
		}
	}
	return
}
