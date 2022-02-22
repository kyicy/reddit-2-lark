package provider

type LarkBotMsgReq struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Post struct {
			ZhCn struct {
				Title   string                     `json:"title"`
				Content [][]map[string]interface{} `json:"content"`
			} `json:"zh_cn"`
		} `json:"post"`
	} `json:"content"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}
