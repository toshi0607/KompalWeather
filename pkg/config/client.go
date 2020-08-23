package config

import "context"

//secret manager使おう

type Config struct {
	// secretManager
	slack   SlackConfig
	twitter TwitterConfig
}

type SlackConfig struct {
	channelNames []string
	token        string
}

type TwitterConfig struct {
}

func New() {

}

func (c Config) Init(ctx context.Context) (Config, error) {
	// Get secrets
	// get env
	return Config{}, nil
}
