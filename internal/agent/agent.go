package agent

import (
	"context"

	"github.com/kyicy/reddit-2-lark/internal/platform"
	"github.com/kyicy/reddit-2-lark/internal/repo/provider"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Agent struct {
	*cron.Cron
}

func NewAgent(conf *platform.Config) *Agent {
	agent := &Agent{
		cron.New(),
	}
	lp := provider.NewLarkProvider(conf)
	rp := provider.NewRedditProvider(conf)
	z, _ := zap.NewDevelopment()
	logger := z.Sugar().Named("agent")

	cronFunc := func() {
		ctx := context.Background()
		res, err := rp.GetTopPosts(ctx, conf.Reddit.Subreddit)
		if err != nil {
			logger.Error(err)
			return
		}
		err = lp.Broadcast(ctx, res)
		if err != nil {
			logger.Error(err)
		}

	}
	agent.Cron.AddFunc(conf.Cron.Internal, cronFunc)
	agent.Start()
	return agent
}
