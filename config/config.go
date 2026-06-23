package config

import (
	"github.com/go-playground/validator/v10"
	config "github.com/kayac/go-config"
	"github.com/morikuni/failure/v2"
)

var Config Root

type Root struct {
	Port  int   `validate:"required"`
	Slack Slack `toml:"slack" validate:"required"`
	LLM   LLM   `toml:"llm" validate:"required"`
	Tool  Tool  `toml:"tool"`
}

type Tool struct {
	WebSearch WebSearch `toml:"web_search"`
	WebFetch  WebFetch  `toml:"web_fetch"`
}

type WebSearch struct {
	TavilyAPIKey string `toml:"tavily_api_key"`
}

type WebFetch struct {
	WhitelistPatterns []string `toml:"whitelist_patterns"`
}

type Slack struct {
	SigningSecret string        `toml:"signing_secret" validate:"required"`
	BotUserID     string        `toml:"bot_user_id" validate:"required"`
	BotToken      string        `toml:"bot_token" validate:"required"`
	AppToken      string        `toml:"app_token" validate:"required"`
	ErrorFeedback ErrorFeedback `toml:"error_feedback"`
}

type ErrorFeedback struct {
	EnableReaction    *bool  `toml:"enable_reaction"`
	ReactionEmoji     string `toml:"reaction_emoji"`
	EnablePostSnippet *bool  `toml:"enable_post_snippet"`
}

func (e ErrorFeedback) GetEnableReaction() bool {
	if e.EnableReaction == nil {
		return true
	}
	return *e.EnableReaction
}

func (e ErrorFeedback) GetReactionEmoji() string {
	if e.ReactionEmoji == "" {
		return "warning"
	}
	return e.ReactionEmoji
}

func (e ErrorFeedback) GetEnablePostSnippet() bool {
	if e.EnablePostSnippet == nil {
		return false
	}
	return *e.EnablePostSnippet
}

type LLM struct {
	ModelType    string `toml:"model_type" validate:"required"`
	ModelID      string `toml:"model_id" validate:"required"`
	SystemPrompt string `toml:"system_prompt"`
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
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(Config); err != nil {
		return failure.Wrap(err, failure.Message("failed to validate Config"))
	}

	return nil
}
