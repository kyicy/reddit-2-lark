package provider

import (
	"context"
	"testing"

	"github.com/kyicy/reddit-2-lark/internal/platform"
	"github.com/stretchr/testify/require"
)

func TestLarkBroadcast(t *testing.T) {
	conf := platform.GetEnvConfig()
	lp := NewLarkProvider(conf)
	rp := NewRedditProvider(conf)

	ctx := context.Background()

	res, err := rp.GetTopPosts(ctx, conf.Reddit.Subreddit)
	require.NoError(t, err)

	err = lp.Broadcast(ctx, res)
	require.NoError(t, err)
}
