package provider

import (
	"context"
	"testing"

	"github.com/kyicy/reddit-2-lark/internal/platform"
	"github.com/stretchr/testify/require"
)

func TestRedditAccessToken(t *testing.T) {
	conf := platform.GetEnvConfig()
	rp := NewRedditProvider(conf)
	ctx := context.Background()

	_, err := rp.getAccessToken(ctx)
	require.NoError(t, err)
}

func TestRedditTopPosts(t *testing.T) {
	conf := platform.GetEnvConfig()
	rp := NewRedditProvider(conf)
	ctx := context.Background()

	_, err := rp.GetTopPosts(ctx, conf.Reddit.Subreddit)
	require.NoError(t, err)
}
