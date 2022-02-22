package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kyicy/rss-2-lark/internal/agent"
	"github.com/kyicy/rss-2-lark/internal/platform"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cli := &cli{}

	cmd := &cobra.Command{
		PreRunE: cli.setup,
		RunE:    cli.run,
	}

	if err := setupFlags(cmd); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

type cli struct {
	cfg   platform.Config
	xMark platform.XMark
}

func setupFlags(cmd *cobra.Command) error {
	cmd.Flags().String("conf", "config.toml", "Path to config file.")
	cmd.Flags().String("x-mark", "x-mark.json", "Path to x-mark file.")
	return viper.BindPFlags(cmd.Flags())
}

func (c *cli) setup(cmd *cobra.Command, args []string) error {
	{
		// conf file
		configFile, err := cmd.Flags().GetString("conf")
		if err != nil {
			return err
		}
		if _, err = os.Stat(configFile); err != nil {
			return err
		}
		// read from config file
		viper.SetConfigFile(configFile)
		if err = viper.ReadInConfig(); err != nil {
			return err
		}
		if err = viper.Unmarshal(&c.cfg); err != nil {
			return err
		}
	}
	{
		// x-mark file
		xMarkFile, err := cmd.Flags().GetString("x-mark")
		if err != nil {
			return err
		}
		_, err = os.Stat(xMarkFile)
		if errors.Is(err, os.ErrNotExist) {
			c.xMark.Items = make(map[string]*platform.Mark)
			return nil
		}
		if err != nil {
			return err
		}
		bs, err := os.ReadFile(xMarkFile)
		if err != nil {
			return err
		}
		return json.Unmarshal(bs, &c.xMark)
	}
}

func (c *cli) run(cmd *cobra.Command, args []string) error {
	_ = agent.NewAgent(&c.cfg, &c.xMark)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	bs, err := json.Marshal(c.xMark)
	if err != nil {
		return err
	}
	xMarkFile, err := cmd.Flags().GetString("x-mark")
	if err != nil {
		return err
	}
	return os.WriteFile(xMarkFile, bs, 0755)
}
