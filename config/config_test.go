package config

import (
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_ErrorFeedbackDefaults(t *testing.T) {
	// Create temporary TOML config without error_feedback section
	content := `
port = 8080
[slack]
signing_secret = "secret"
bot_user_id = "@U123"
bot_token = "xoxb-token"
app_token = "xapp-token"
[llm]
model_type = "claude"
model_id = "anthropic"
`
	tmpFile, err := os.CreateTemp("", "config_test_*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)
	tmpFile.Close()

	err = Load(tmpFile.Name())
	require.NoError(t, err)

	// Verify defaults resolved
	assert.True(t, Config.Slack.ErrorFeedback.GetEnableReaction())
	assert.Equal(t, "warning", Config.Slack.ErrorFeedback.GetReactionEmoji())
	assert.False(t, Config.Slack.ErrorFeedback.GetEnablePostSnippet())
}

func TestConfig_ErrorFeedbackOverrides(t *testing.T) {
	// Create temporary TOML config with error_feedback overrides
	content := `
port = 8080
[slack]
signing_secret = "secret"
bot_user_id = "@U123"
bot_token = "xoxb-token"
app_token = "xapp-token"
[slack.error_feedback]
enable_reaction = false
reaction_emoji = "boom"
enable_post_snippet = true
[llm]
model_type = "claude"
model_id = "anthropic"
`
	tmpFile, err := os.CreateTemp("", "config_test_*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)
	tmpFile.Close()

	err = Load(tmpFile.Name())
	require.NoError(t, err)

	// Verify overrides resolved
	assert.False(t, Config.Slack.ErrorFeedback.GetEnableReaction())
	assert.Equal(t, "boom", Config.Slack.ErrorFeedback.GetReactionEmoji())
	assert.True(t, Config.Slack.ErrorFeedback.GetEnablePostSnippet())
}

func TestConfig_Tool(t *testing.T) {
	// Test mapping from a mock TOML representation to the config structure
	tomlData := `
port = 9000
[slack]
signing_secret = "secret"
bot_user_id = "bot"
bot_token = "token"
app_token = "app"
[llm]
model_type = "claude"
model_id = "model"
[tool]
[tool.web_search]
tavily_api_key = "tavily-key"
[tool.web_fetch]
whitelist_patterns = ["^https://example\\.com/.*"]
`
	var cfg Root
	err := toml.Unmarshal([]byte(tomlData), &cfg)
	require.NoError(t, err)
	assert.Equal(t, "tavily-key", cfg.Tool.WebSearch.TavilyAPIKey)
	assert.Equal(t, []string{"^https://example\\.com/.*"}, cfg.Tool.WebFetch.WhitelistPatterns)
}

