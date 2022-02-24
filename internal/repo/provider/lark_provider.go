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
	"sync"
	"time"

	"github.com/kyicy/request"
	"github.com/kyicy/rss-2-lark/internal/platform"
	"go.uber.org/zap"
)

type Broadcastable interface {
	GetTitle() string
	GetLink() string
	GetPubDate() time.Time
}

type BroadcastSource interface {
	GetName() string
	GetTopPosts(context.Context) ([]Broadcastable, error)
}

type LarkProvider struct {
	config     *platform.Config
	xMark      *platform.XMark
	logger     *zap.SugaredLogger
	httpClient *request.RequestProvider
	sources    []BroadcastSource
}

func NewLarkProvider(
	conf *platform.Config,
	xMark *platform.XMark,
	sources ...BroadcastSource,
) *LarkProvider {
	rp, _ := request.NewRequestProvider(http.DefaultClient)
	z, _ := zap.NewDevelopment()
	return &LarkProvider{
		config:     conf,
		xMark:      xMark,
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
	markItems := lp.xMark.Items
	var wg sync.WaitGroup
	var mutex sync.Mutex

	botMsgs := make([]*LarkBotMsgReq, 0, 16)
	for _, src := range lp.sources {
		wg.Add(1)
		go func(src BroadcastSource) {
			defer wg.Done()
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
			t.Title = src.GetName()
			t.Content = make([][]map[string]interface{}, 0)

			items, err := src.GetTopPosts(ctx)
			if err != nil {
				lp.logger.Error(err)
				return
			}

			mark := markItems[src.GetName()]
			if mark == nil {
				mark = &platform.Mark{}
			}

			count := len(items)

			initMarkDate := mark.LastPubDate
			latestMarkDate := mark.LastPubDate

			for i := range items {
				item := items[len(items)-1-i]
				pubDate := item.GetPubDate()
				if !pubDate.After(initMarkDate) {
					lp.logger.Debugw("skipped item", "item.title", item.GetTitle(), "item.pubDate", item.GetPubDate(), "mark.LastPubDate", mark.LastPubDate)
					continue
				}
				index := count - i
				t.Content = append(t.Content, []map[string]interface{}{
					{
						"tag":       "text",
						"un_escape": true,
						"lines":     1,
						"text":      fmt.Sprintf("%2d: ", index),
					},
					{
						"tag":  "a",
						"text": item.GetTitle(),
						"href": item.GetLink(),
					},
				})
				mark.LastPubDate = item.GetPubDate()
				if pubDate.After(latestMarkDate) {
					latestMarkDate = pubDate
				}
			}

			mark.LastPubDate = latestMarkDate
			mark.LastPubDateString = latestMarkDate.Format(time.RFC3339)

			if len(t.Content) == 0 {
				return
			}

			for i, j := 0, len(t.Content)-1; i < j; i, j = i+1, j-1 {
				fmt.Println(i, j)
				t.Content[i], t.Content[j] = t.Content[j], t.Content[i]
			}

			mutex.Lock()
			markItems[src.GetName()] = mark
			mutex.Unlock()

			mutex.Lock()
			botMsgs = append(botMsgs, botMsgReq)
			mutex.Unlock()
		}(src)
	}
	wg.Wait()

	for _, botMsgReq := range botMsgs {
		body, err := json.Marshal(botMsgReq)
		if err != nil {
			lp.logger.Error(err)
			continue
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
			continue
		}
		res, err := lp.httpClient.Do(req)
		// to avoid too many request error
		time.Sleep(time.Second)
		if err != nil {
			lp.logger.Error(err)
			continue
		}
		defer res.Body.Close()
		bs, err := io.ReadAll(res.Body)
		if err != nil {
			lp.logger.Error(err)
			continue
		}
		lp.logger.Info(string(bs))
	}

}
