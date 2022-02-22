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

	"github.com/kyicy/request"
	"github.com/kyicy/rss-2-lark/internal/platform"
	"go.uber.org/zap"
)

type Broadcastable interface {
	GetTitle() string
	GetLink() string
}

type BroadcastSource interface {
	GetName() string
	GetTopPosts(context.Context) ([]Broadcastable, error)
}

type LarkProvider struct {
	config     *platform.Config
	logger     *zap.SugaredLogger
	httpClient *request.RequestProvider
	sources    []BroadcastSource
}

func NewLarkProvider(conf *platform.Config, sources ...BroadcastSource) *LarkProvider {
	rp, _ := request.NewRequestProvider(http.DefaultClient)
	z, _ := zap.NewDevelopment()
	return &LarkProvider{
		config:     conf,
		logger:     z.Sugar().Named("lark_provider"),
		httpClient: rp,
		sources:    sources,
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
) {
	timestamp := time.Now().UnixNano() / int64(time.Second)
	sign, err := lp.genSign(timestamp)
	if err != nil {
		lp.logger.Error(err)
		return
	}
	botMsgReq := &LarkBotMsgReq{
		MsgType:   "post",
		Timestamp: fmt.Sprintf("%d", timestamp),
		Sign:      sign,
	}
	t := &botMsgReq.Content.Post.ZhCn
	t.Title = lp.config.Lark.Header
	t.Content = make([][]map[string]interface{}, 0)

	for _, src := range lp.sources {
		items, err := src.GetTopPosts(ctx)
		if err != nil {
			lp.logger.Error(err)
			continue
		}

		for i, item := range items {
			t.Content = append(t.Content, []map[string]interface{}{
				{
					"tag":       "text",
					"un_escape": true,
					"lines":     1,
					"text":      fmt.Sprintf("%2d: ", i+1),
				},
				{
					"tag":  "a",
					"text": fmt.Sprintf("[%s]%s", src.GetName(), item.GetTitle()),
					"href": item.GetLink(),
				},
			})
		}

	}

	if len(t.Content) == 0 {
		return
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
}
