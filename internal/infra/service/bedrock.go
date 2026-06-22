package service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

const (
	BEDROCK_MAX_TOKENS  = 10_000
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
	System           string                 `json:"system,omitempty"` // system prompt
	Messages         []ClaudeRequestMessage `json:"messages"`
	Tools            []ClaudeRequestTool    `json:"tools,omitempty"` // tools list
}

type ClaudeRequestTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema"`
}

type ClaudeRequestMessage struct {
	Role    string                        `json:"role"`
	Content []ClaudeRequestMessageContent `json:"content"`
}

type ClaudeRequestMessageContent struct {
	Type      string          `json:"type"` // "text", "tool_use", "tool_result"
	Text      string          `json:"text,omitempty"`
	ID        string          `json:"id,omitempty"`          // tool use id
	Name      string          `json:"name,omitempty"`        // tool name
	Input     json.RawMessage `json:"input,omitempty"`       // tool input arguments
	ToolUseID string          `json:"tool_use_id,omitempty"` // tool result reference
	Content   string          `json:"content,omitempty"`     // tool result output
	IsError   bool            `json:"is_error,omitempty"`
}

type ClaudeResponse struct {
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Role       string                  `json:"role"`
	Content    []ClaudeResponseContent `json:"content"`
	StopReason string                  `json:"stop_reason"`
}

type ClaudeResponseContent struct {
	Type  string          `json:"type"`
	Text  string          `json:"text,omitempty"`
	ID    string          `json:"id,omitempty"`
	Name  string          `json:"name,omitempty"`
	Input json.RawMessage `json:"input,omitempty"`
}

type Bedrock struct {
	modelID string
	client  *bedrockruntime.Client
}

func NewBedrock(ctx context.Context, modelID string) (*Bedrock, error) {
	awsConf, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to load aws config"))
	}

	client := bedrockruntime.NewFromConfig(awsConf)

	return &Bedrock{
		modelID: modelID,
		client:  client,
	}, nil
}

func (b *Bedrock) Invoke(ctx context.Context, prompt string) (string, error) {
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
		return "", failure.Wrap(err, failure.Message("failed to marshal request for bedrock"))
	}

	out, err := b.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(b.modelID),
		ContentType: aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return "", b.wrapBedrockError(err)
	}

	var res ClaudeResponse
	if err := json.Unmarshal(out.Body, &res); err != nil {
		log.Debug().Bytes("out.Body", out.Body).Msg("failed to unmarshal response from bedrock")
		return "", failure.Wrap(err, failure.Message("failed to unmarshal response from bedrock"))
	}

	if len(res.Content) == 0 {
		return "", failure.New(errorcode.ErrInternal, failure.Message("empty content in bedrock response"))
	}

	return res.Content[0].Text, nil
}

func (b *Bedrock) InvokeWithTools(
	ctx context.Context,
	systemPrompt string,
	messages []*reasoning.ContextMessage,
	tools []reasoning.Tool,
) (*reasoning.Answer, error) {
	reqTools := make([]ClaudeRequestTool, 0, len(tools))
	for _, t := range tools {
		reqTools = append(reqTools, ClaudeRequestTool{
			Name:        t.Name().Value(),
			Description: t.Description(),
			InputSchema: json.RawMessage(t.InputSchema()),
		})
	}

	reqMessages := make([]ClaudeRequestMessage, 0, len(messages))
	for _, m := range messages {
		if m.Role() == "system" {
			continue
		}
		var contents []ClaudeRequestMessageContent

		if m.Content() != "" {
			contents = append(contents, ClaudeRequestMessageContent{
				Type: "text",
				Text: string(m.Content()),
			})
		}

		for _, tc := range m.ToolCalls() {
			inputJSON := tc.InputJSON()
			if inputJSON == "" {
				inputJSON = "{}"
			}
			contents = append(contents, ClaudeRequestMessageContent{
				Type:  "tool_use",
				ID:    tc.ID().Value(),
				Name:  tc.Name().Value(),
				Input: json.RawMessage(inputJSON),
			})
		}

		for _, tr := range m.ToolResults() {
			contents = append(contents, ClaudeRequestMessageContent{
				Type:      "tool_result",
				ToolUseID: tr.ToolUseID().Value(),
				Content:   tr.Output(),
				IsError:   bool(tr.Status() == reasoning.ToolResultError),
			})
		}

		reqMessages = append(reqMessages, ClaudeRequestMessage{
			Role:    m.Role(),
			Content: contents,
		})
	}

	body, err := json.Marshal(ClaudeRequest{
		AnthropicVersion: ANTHROPIC_VERSION,
		MaxTokens:        BEDROCK_MAX_TOKENS,
		Temperature:      BEDROCK_TEMPERATURE,
		System:           systemPrompt,
		Messages:         reqMessages,
		Tools:            reqTools,
		StopSequences:    []string{},
	})
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to marshal request for bedrock"))
	}

	out, err := b.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(b.modelID),
		ContentType: aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return nil, b.wrapBedrockError(err)
	}

	var res ClaudeResponse
	if err := json.Unmarshal(out.Body, &res); err != nil {
		log.Debug().Bytes("out.Body", out.Body).Msg("failed to unmarshal response from bedrock")
		return nil, failure.Wrap(err, failure.Message("failed to unmarshal response from bedrock"))
	}

	if len(res.Content) == 0 {
		return nil, failure.New(errorcode.ErrInternal, failure.Message("empty content in bedrock response"))
	}

	var answerText string
	var toolCalls []reasoning.ToolCall

	for _, c := range res.Content {
		switch c.Type {
		case "text":
			answerText += c.Text
		case "tool_use":
			toolCallID, err := reasoning.NewToolCallID(c.ID)
			if err != nil {
				return nil, failure.Wrap(err, failure.Message("failed to parse tool call ID from bedrock response"))
			}
			toolName, err := reasoning.NewToolName(c.Name)
			if err != nil {
				return nil, failure.Wrap(err, failure.Message("failed to parse tool name from bedrock response"))
			}
			tc, err := reasoning.NewToolCall(toolCallID, toolName, string(c.Input))
			if err != nil {
				return nil, failure.Wrap(err, failure.Message("failed to construct ToolCall"))
			}
			toolCalls = append(toolCalls, tc)
		}
	}

	ans, err := reasoning.NewAnswer(reasoning.AnswerBody(answerText), toolCalls)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to create reasoning answer"))
	}
	return ans, nil
}

func (b *Bedrock) wrapBedrockError(err error) error {
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
                    `, b.modelID))
	default:
		return failure.Wrap(err, failure.Messagef("Couldn't invoke model: %s", b.modelID))
	}
}
