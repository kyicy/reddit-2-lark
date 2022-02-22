package provider

type accessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

// PostListResp only select minimal fields
type PostListResp struct {
	Data struct {
		Children []struct {
			Data struct {
				Title     string `json:"title"`
				PermaLink string `json:"permalink"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}
