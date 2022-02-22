package platform

import "os"

// Config definition
type Config struct {
	Keywords []string `toml:"keywords" json:"keywords"`

	Cron struct {
		Internal string `toml:"internal" json:"internal"`
	} `toml:"cron" json:"cron"`

	Feed map[string]struct {
		Name string `toml:"name" json:"name"`
		Src  string `toml:"src" json:"src"`
	} `toml:"feed" json:"feed"`

	Lark struct {
		Header string `toml:"header" json:"header"`
		Hook   string `toml:"hook" json:"hook"`
		Token  string `toml:"token" json:"token"`
	} `toml:"lark" json:"lark"`
}

func GetEnvConfig() *Config {
	config := &Config{}
	config.Lark.Hook = os.Getenv("lark_hook")
	config.Lark.Token = os.Getenv("lark_token")
	return config
}
