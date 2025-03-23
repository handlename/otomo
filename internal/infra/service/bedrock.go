package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

// https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/bedrock-runtime/actions/invoke_model.go

type ClaudeRequest struct {
	Prompt            string   `json:"prompt"`
	MaxTokensToSample int      `json:"max_tokens_to_sample"`
	Temperature       float64  `json:"temperature,omitempty"`
	StopSequences     []string `json:"stop_sequences,omitempty"`
}

type ClaudeResponse struct {
	Completion string `json:"completion"`
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
	modelID := "anthropic.claude-3-5-haiku-20241022-v1:0"

	body, err := json.Marshal(ClaudeRequest{
		Prompt:            fmt.Sprintf("Human: %s\n\nAssistant:", prompt),
		MaxTokensToSample: 200,
		Temperature:       0.5,
		StopSequences:     []string{"\n\nHuman:"},
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

	return res.Completion, nil
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
		return failure.Wrap(err, failure.Messagef("Couldn't invoke model: \"%s\". Here's why: %w\n", modelID, err))
	}
}
