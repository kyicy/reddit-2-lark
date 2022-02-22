package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/kyicy/reddit-2-lark/internal/platform"
	"github.com/kyicy/request"
	"go.uber.org/zap"
)

const (
	redditApiRoot   = "https://www.reddit.com/api"
	redditOauthRoot = "https://oauth.reddit.com"
	accessTokenPath = "/v1/access_token"
)

type RedditProvider struct {
	config          *platform.Config
	cache           platform.Cache
	accessTokenLock sync.Mutex
	httpClient      *request.RequestProvider
	logger          *zap.SugaredLogger
}

func NewRedditProvider(conf *platform.Config) *RedditProvider {
	rp, _ := request.NewRequestProvider(http.DefaultClient)
	z, _ := zap.NewDevelopment()
	return &RedditProvider{
		cache:      platform.NewMemory(),
		config:     conf,
		httpClient: rp,
		logger:     z.Sugar().Named("reddit_provider"),
	}
}

func (rp *RedditProvider) getAccessToken(
	ctx context.Context,
) (accessToken string, err error) {
	rp.accessTokenLock.Lock()
	defer rp.accessTokenLock.Unlock()

	// first try to get from cache
	accessTokenCacheKey := fmt.Sprintf("reddit_%s_access_token", rp.config.Reddit.ClientId)
	val := rp.cache.Get(accessTokenCacheKey)
	if val != nil {
		accessToken = val.(string)
		return
	}

	// get from reddit server
	targetUrl := redditApiRoot + accessTokenPath
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", rp.config.Reddit.Username)
	data.Set("password", rp.config.Reddit.Password)
	if err != nil {
		rp.logger.Error(err)
		return
	}

	req, err := request.NewRequestWithContext(
		ctx,
		http.MethodPost,
		targetUrl,
		struct {
			UserAgent string `rh:"User-Agent"`
		}{"android:com.github.kyicy:v0.0.1 (by /u/kyicy)"},
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		rp.logger.Error(err)
		return
	}
	req.SetBasicAuth(rp.config.Reddit.ClientId, rp.config.Reddit.Secret)
	res, err := rp.httpClient.Do(req)
	if err != nil {
		rp.logger.Error(err)
		return
	}
	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		rp.logger.Error(err)
		return
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf(string(bs))
		rp.logger.Error(err)
		return
	}
	tokenResp := new(accessTokenResp)
	err = json.Unmarshal(bs, tokenResp)
	if err != nil {
		rp.logger.Error(err)
		return
	}
	expires := tokenResp.ExpiresIn - 1000
	accessToken = tokenResp.AccessToken

	if expires <= 0 {
		return
	}

	_ = rp.cache.Set(accessTokenCacheKey, accessToken, time.Duration(expires)*time.Second)
	return
}

func (rp *RedditProvider) GetHeader() string {
	return "top /r/Eldenring posts"
}

func (rp *RedditProvider) GetTopPosts(
	ctx context.Context,
) (items []Broadcastable, err error) {
	token, err := rp.getAccessToken(ctx)
	if err != nil {
		rp.logger.Error(err)
		return
	}

	targetUrl := fmt.Sprintf("%s/r/%s/hot", redditOauthRoot, rp.config.Reddit.Subreddit)

	req, err := request.NewRequestWithContext(
		ctx,
		http.MethodGet,
		targetUrl,
		struct {
			UserAgent     string `rh:"User-Agent"`
			Authorization string `rh:"Authorization"`
			RawJson       string `rq:"raw_json"`
		}{
			"android:com.github.kyicy:v0.0.1 (by /u/kyicy)",
			fmt.Sprintf("Bearer %s", token),
			"1",
		},
		nil,
	)
	if err != nil {
		rp.logger.Error(err)
		return
	}
	res, err := rp.httpClient.Do(req)
	if err != nil {
		rp.logger.Error(err)
		return
	}

	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		rp.logger.Error(err)
		return
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf(string(bs))
		rp.logger.Error(err)
		return
	}
	postListResp := new(PostListResp)
	err = json.Unmarshal(bs, postListResp)
	if err != nil {
		rp.logger.Error(err)
		return
	}
	items = make([]Broadcastable, len(postListResp.Data.Children))
	for i, item := range postListResp.Data.Children {
		items[i] = item
	}

	return
}
