package service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

const (
	BEDROCK_MODEL_ID    = "arn:aws:bedrock:us-east-1:014347307013:inference-profile/us.anthropic.claude-3-5-haiku-20241022-v1:0"
	BEDROCK_MAX_TOKENS  = 200
	BEDROCK_TEMPERATURE = 0.5
	ANTHROPIC_VERSION   = "bedrock-2023-05-31"
)

// https://docs.anthropic.com/en/api/messages
// https://docs.aws.amazon.com/bedrock/latest/userguide/model-parameters-anthropic-claude-messages.html#model-parameters-anthropic-claude-messages-request-response

type ClaudeRequest struct {
	AnthropicVersion string                 `json:"anthropic_version"`
	MaxTokens        int                    `json:"max_tokens"`
	Temperature      float64                `json:"temperature"`
	StopSequences    []string               `json:"stop_sequences"`
	Messages         []ClaudeRequestMessage `json:"messages"`
}

type ClaudeRequestMessage struct {
	Role    string                        `json:"role"`
	Content []ClaudeRequestMessageContent `json:"content"`
}

type ClaudeRequestMessageContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ClaudeResponse struct {
	ID           string                  `json:"id"`
	Model        string                  `json:"model"`
	Type         string                  `json:"type"`
	Role         string                  `json:"role"`
	Content      []ClaudeResponseContent `json:"content"`
	StopReason   string                  `json:"stop_reason"`
	StopSequence string                  `json:"stop_sequence"`
	Usage        ClaudeResponseUsage     `json:"usage"`
}

type ClaudeResponseContent struct {
	Type  string          `json:"type"`
	Text  string          `json:"text,omitempty"`
	Image json.RawMessage `json:"image,omitempty"`
	ID    string          `json:"id,omitempty"`
	Name  string          `json:"name,omitempty"`
	Input json.RawMessage `json:"input,omitempty"`
}

type ClaudeResponseUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type Bedrock struct {
	client *bedrockruntime.Client
}

func NewBedrock(ctx context.Context) (*Bedrock, error) {
	awsConf, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to load aws config"))
	}

	client := bedrockruntime.NewFromConfig(awsConf)

	return &Bedrock{
		client: client,
	}, nil
}

func (b *Bedrock) Invoke(ctx context.Context, prompt string) (string, error) {
	modelID := BEDROCK_MODEL_ID

	body, err := json.Marshal(ClaudeRequest{
		AnthropicVersion: ANTHROPIC_VERSION,
		MaxTokens:        BEDROCK_MAX_TOKENS,
		Temperature:      BEDROCK_TEMPERATURE,
		Messages: []ClaudeRequestMessage{
			{
				Role: "user",
				Content: []ClaudeRequestMessageContent{
					{
						Type: "text",
						Text: prompt,
					},
				},
			},
		},
		StopSequences: []string{},
	})
	if err != nil {
		return "", failure.Wrap(err, failure.Message("failed to marshal request for claude"))
	}

	out, err := b.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelID),
		ContentType: aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return "", wrapBedrockError(err, modelID)
	}

	var res ClaudeResponse
	if err := json.Unmarshal(out.Body, &res); err != nil {
		log.Debug().Bytes("out.Body", out.Body).Msg("failed to unmarshal response from claude")
		return "", failure.Wrap(err, failure.Message("failed to unmarshal response from claude"))
	}

	return res.Content[0].Text, nil
}

func wrapBedrockError(err error, modelID string) error {
	msg := err.Error()

	switch {
	case strings.Contains(msg, "no such host"):
		return failure.Wrap(err, failure.Message(`The Bedrock service is not available in the selected region.
                    Please double-check the service availability for your region at
                    https://aws.amazon.com/about-aws/global-infrastructure/regional-product-services/.\n`))
	case strings.Contains(msg, "Could not resolve the foundation model"):
		return failure.Wrap(err, failure.Messagef(`Could not resolve the foundation model from model identifier: "%s".
                    Please verify that the requested model exists and is accessible
                    within the specified region.\n
                    `, modelID))
	default:
		return failure.Wrap(err, failure.Messagef("Couldn't invoke model: %s", modelID))
	}
}
