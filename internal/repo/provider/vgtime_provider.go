package provider

import (
	"context"
	"strings"

	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
)

const (
	vgtimeFeed = "https://www.vgtime.com/rss.jhtml"
)

type VgtimeProvider struct {
	logger *zap.SugaredLogger
}

func NewVgtimeProvider() *VgtimeProvider {
	z, _ := zap.NewDevelopment()
	return &VgtimeProvider{
		logger: z.Sugar().Named("reddit_provider"),
	}
}

func (vp *VgtimeProvider) GetHeader() string {
	return "游戏时光"
}

func (vp *VgtimeProvider) GetTopPosts(
	ctx context.Context,
) (items []Broadcastable, err error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(vgtimeFeed, ctx)
	if err != nil {
		return
	}
	items = make([]Broadcastable, 0)
	for _, item := range feed.Items {
		if strings.Contains(item.Title, "艾尔登法环") {
			items = append(items, (*feedItem)(item))
		}
	}
	return
}
