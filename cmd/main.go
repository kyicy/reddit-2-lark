package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kyicy/reddit-2-lark/internal/agent"
	"github.com/kyicy/reddit-2-lark/internal/platform"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cli := &cli{}

	cmd := &cobra.Command{
		PreRunE: cli.setupConfig,
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
	cfg platform.Config
}

func setupFlags(cmd *cobra.Command) error {
	cmd.Flags().String("conf", "config.toml", "Path to config file.")
	return viper.BindPFlags(cmd.Flags())
}

func (c *cli) setupConfig(cmd *cobra.Command, args []string) error {
	var err error

	configFile, err := cmd.Flags().GetString("conf")
	if err != nil {
		return err
	}
	_, err = os.Stat(configFile)
	if err != nil {
		return err
	}
	// read from config file
	viper.SetConfigFile(configFile)
	if err = viper.ReadInConfig(); err != nil {
		return err
	}
	return viper.Unmarshal(&c.cfg)
}

func (c *cli) run(cmd *cobra.Command, args []string) error {
	_ = agent.NewAgent(&c.cfg)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	return nil
}
