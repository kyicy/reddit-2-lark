package platform

import "os"

// Config definition
type Config struct {
	Reddit struct {
		Subreddit string `toml:"subreddit" json:"subreddit"`
		ClientId  string `toml:"clientId" json:"clientId"`
		Secret    string `toml:"secret" json:"secret"`
		Username  string `toml:"username" json:"username"`
		Password  string `toml:"password" json:"password"`
	} `toml:"reddit" json:"reddit"`
	Lark struct {
		HookId string `toml:"hookId" json:"hookId"`
		Token  string `toml:"token" json:"token"`
	} `toml:"lark" json:"lark"`
}

func GetEnvConfig() *Config {
	config := &Config{}
	config.Reddit.Subreddit = "Eldenring"
	config.Reddit.ClientId = os.Getenv("reddit_client_id")
	config.Reddit.Secret = os.Getenv("reddit_secret")
	config.Reddit.Username = os.Getenv("reddit_username")
	config.Reddit.Password = os.Getenv("reddit_password")
	config.Lark.HookId = os.Getenv("lark_hook_id")
	config.Lark.Token = os.Getenv("lark_token")
	return config
}
