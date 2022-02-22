package agent

import (
	"context"
	"time"

	"github.com/kyicy/reddit-2-lark/internal/platform"
	"github.com/kyicy/reddit-2-lark/internal/repo/provider"
	"github.com/robfig/cron/v3"
)

type Agent struct {
	*cron.Cron
}

func NewAgent(conf *platform.Config) *Agent {
	agent := &Agent{
		cron.New(),
	}

	rp := provider.NewRedditProvider(conf)
	vp := provider.NewVgtimeProvider()

	lp := provider.NewLarkProvider(conf, rp, vp)
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
