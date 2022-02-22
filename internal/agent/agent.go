package agent

import (
	"context"
	"time"

	"github.com/kyicy/rss-2-lark/internal/platform"
	"github.com/kyicy/rss-2-lark/internal/repo/provider"
	"github.com/robfig/cron/v3"
)

type Agent struct {
	*cron.Cron
}

func NewAgent(conf *platform.Config, xMark *platform.XMark) *Agent {
	agent := &Agent{
		cron.New(),
	}

	for _, mark := range xMark.Items {
		mark.LastPubDate, _ = time.Parse(time.RFC3339, mark.LastPubDateString)
	}

	providers := make([]provider.BroadcastSource, 0)

	for _, rssSrc := range conf.Feed {
		providers = append(
			providers,
			provider.NewRssProvider(rssSrc.Name, rssSrc.Src, conf),
		)
	}

	lp := provider.NewLarkProvider(conf, xMark, providers...)
	cronFunc := func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Minute*2)
		defer cancel()
		lp.Broadcast(ctx)
	}
	go cronFunc()
	agent.Cron.AddFunc(conf.Cron.Internal, cronFunc)
	agent.Start()
	return agent
}
