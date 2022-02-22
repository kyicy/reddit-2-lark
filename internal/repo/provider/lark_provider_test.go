package provider

import (
	"context"
	"testing"

	"github.com/kyicy/reddit-2-lark/internal/platform"
)

func TestLarkBroadcast(t *testing.T) {
	conf := platform.GetEnvConfig()
	rp := NewRedditProvider(conf)
	lp := NewLarkProvider(conf, rp)
	ctx := context.Background()
	lp.Broadcast(ctx)
}
