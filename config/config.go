package config

import (
	"github.com/handlename/otomo/internal/errorcode"
	config "github.com/kayac/go-config"
	"github.com/morikuni/failure/v2"
)

var Config Root

type Root struct {
	Port    int
	Slack   Slack   `toml:"slack"`
	Bedrock Bedrock `toml:"bedrock"`
}

type Slack struct {
	AppToken string `toml:"app_token"`
	BotToken string `toml:"bot_token"`
}

type Bedrock struct {
	ModelType string `toml:"model_type"`
	ModelID   string `toml:"model_id"`
}

// Load load config data to otomo.Config from TOML file specified by path.
func Load(path string) error {
	if err := config.LoadWithEnvTOML(&Config, path); err != nil {
		return failure.Wrap(err, failure.Messagef("failed to load config from %s", path))
	}

	if err := Validate(); err != nil {
		return failure.Wrap(err, failure.Message("config value(s) are invalid"))
	}

	return nil
}

func Validate() error {
	if s := Config.Slack; s.AppToken == "" || s.BotToken == "" {
		return failure.New(
			errorcode.ErrInvalidArgument,
			failure.Message("both of Slack App token and Slack Bot token are required"),
		)
	}

	return nil
}
