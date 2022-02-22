package provider

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kyicy/reddit-2-lark/internal/platform"
	"github.com/kyicy/request"
	"go.uber.org/zap"
)

type LarkProvider struct {
	config     *platform.Config
	logger     *zap.SugaredLogger
	httpClient *request.RequestProvider
}

func NewLarkProvider(conf *platform.Config) *LarkProvider {
	rp, _ := request.NewRequestProvider(http.DefaultClient)
	z, _ := zap.NewDevelopment()
	return &LarkProvider{
		config:     conf,
		logger:     z.Sugar().Named("lark_provider"),
		httpClient: rp,
	}
}

func (lp *LarkProvider) genSign(timestamp int64) (string, error) {
	//timestamp + key sha256, and then base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + lp.config.Lark.Token
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}

func (lp *LarkProvider) Broadcast(
	ctx context.Context,
	postListResp *PostListResp,
) (err error) {
	timestamp := time.Now().UnixNano() / int64(time.Second)
	sign, err := lp.genSign(timestamp)
	if err != nil {
		return
	}
	botMsgReq := &LarkBotMsgReq{
		MsgType:   "post",
		Timestamp: fmt.Sprintf("%d", timestamp),
		Sign:      sign,
	}

	t := &botMsgReq.Content.Post.ZhCn
	t.Title = fmt.Sprintf("Breaking posts from /r/%s", lp.config.Reddit.Subreddit)
	t.Content = make([][]map[string]interface{}, 0)
	for i, child := range postListResp.Data.Children {
		if i >= 14 {
			break
		}
		t.Content = append(t.Content, []map[string]interface{}{
			{
				"tag":       "text",
				"un_escape": true,
				"lines":     1,
				"text":      fmt.Sprintf("%d: ", i+1),
			},
			{
				"tag":  "a",
				"text": child.Data.Title,
				"href": fmt.Sprintf("https://www.reddit.com%s", child.Data.PermaLink),
			},
		})
	}
	body, err := json.Marshal(botMsgReq)
	if err != nil {
		lp.logger.Error(err)
		return
	}

	targetUrl := lp.config.Lark.Hook

	req, err := request.NewRequestWithContext(
		ctx,
		http.MethodPost,
		targetUrl,
		struct {
			ContentType string `rh:"Content-Type"`
		}{
			"application/json",
		},
		bytes.NewReader(body),
	)
	if err != nil {
		lp.logger.Error(err)
		return
	}
	res, err := lp.httpClient.Do(req)
	if err != nil {
		lp.logger.Error(err)
		return
	}
	defer res.Body.Close()
	bs, err := io.ReadAll(res.Body)
	if err != nil {
		lp.logger.Error(err)
		return
	}
	lp.logger.Info(string(bs))
	return
}
