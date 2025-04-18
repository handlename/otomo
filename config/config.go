package config

import (
	"fmt"

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
	SigningSecret string `toml:"signing_secret"`
	BotUserID     string `toml:"bot_user_id"`
	BotToken      string `toml:"bot_token"`
	AppToken      string `toml:"app_token"`
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
	// TODO: use github.com/go-playground/validator

	if s := Config.Slack; s.SigningSecret == "" || s.AppToken == "" || s.BotToken == "" || s.BotUserID == "" {
		return failure.New(
			errorcode.ErrInvalidArgument,
			failure.Message("configuration for Slack is not satisfied"),
			failure.Context(map[string]string{
				"slack": fmt.Sprintf("%+v", s),
			}),
		)
	}

	return nil
}
