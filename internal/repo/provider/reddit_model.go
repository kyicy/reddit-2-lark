package provider

type accessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type RedditItem struct {
	Data struct {
		Title     string `json:"title"`
		PermaLink string `json:"permalink"`
	} `json:"data"`
}

func (ri *RedditItem) GetTitle() string {
	return ri.Data.Title
}
func (ri *RedditItem) GetLink() string {
	return "https://www.reddit.com" + ri.Data.PermaLink
}

// PostListResp only select minimal fields
type PostListResp struct {
	Data struct {
		Children []*RedditItem `json:"children"`
	} `json:"data"`
}
