package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRoundTripper func(req *http.Request) (*http.Response, error)

func (f mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestBedrock_InvokeWithTools(t *testing.T) {
	var capturedRequestBody []byte
	mockClient := &http.Client{
		Transport: mockRoundTripper(func(req *http.Request) (*http.Response, error) {
			var err error
			capturedRequestBody, err = io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}

			resPayload := bedrockResponse{
				ID:   "msg_123",
				Type: "message",
				Role: "assistant",
				Content: []bedrockResponseContent{
					{
						Type: "text",
						Text: "Here is the tool output: ",
					},
					{
						Type:  "tool_use",
						ID:    "toolu_test_123",
						Name:  "dummy_tool",
						Input: json.RawMessage(`{"text":"hello"}`),
					},
				},
				StopReason: "tool_use",
			}
			resBytes, err := json.Marshal(resPayload)
			if err != nil {
				return nil, err
			}

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewReader(resBytes)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	awsConf := aws.Config{
		HTTPClient:  mockClient,
		Region:      "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider("MOCK_KEY", "MOCK_SECRET", "MOCK_TOKEN"),
	}
	client := bedrockruntime.NewFromConfig(awsConf)
	b := &Bedrock{
		modelID: "anthropic.claude-3-sonnet-v1:0",
		client:  client,
	}

	tName, err := reasoning.NewToolName("dummy_tool")
	require.NoError(t, err)

	testTool := mockTool{
		name:        tName,
		description: "A dummy tool",
		schema:      `{"type":"object"}`,
	}

	msgUser, err := reasoning.NewContextMessage("user", core.UserID{}, "Hello", nil, nil)
	require.NoError(t, err)

	tcID, err := reasoning.NewToolCallID("toolu_test_123")
	require.NoError(t, err)
	tc, err := reasoning.NewToolCall(tcID, tName, `{"text":"hello"}`)
	require.NoError(t, err)
	msgAssistant, err := reasoning.NewContextMessage("assistant", core.UserID{}, "", []reasoning.ToolCall{tc}, nil)
	require.NoError(t, err)

	tr, err := reasoning.NewToolResult(tcID, `{"length": 5}`, false)
	require.NoError(t, err)
	msgResult, err := reasoning.NewContextMessage("user", core.UserID{}, "", nil, []reasoning.ToolResult{tr})
	require.NoError(t, err)

	messages := []*reasoning.ContextMessage{msgUser, msgAssistant, msgResult}

	ctx := context.Background()
	ans, err := b.InvokeWithTools(ctx, "You are a helpful assistant", messages, []reasoning.Tool{testTool})
	require.NoError(t, err)

	assert.Equal(t, reasoning.AnswerBody("Here is the tool output: "), ans.Body())
	require.Len(t, ans.ToolCalls(), 1)
	assert.Equal(t, "toolu_test_123", ans.ToolCalls()[0].ID().Value())
	assert.Equal(t, "dummy_tool", ans.ToolCalls()[0].Name().Value())
	assert.JSONEq(t, `{"text":"hello"}`, ans.ToolCalls()[0].InputJSON())

	var reqPayload bedrockRequest
	err = json.Unmarshal(capturedRequestBody, &reqPayload)
	require.NoError(t, err)

	assert.Equal(t, "bedrock-2023-05-31", reqPayload.AnthropicVersion)
	assert.Equal(t, "You are a helpful assistant", reqPayload.System)
	require.Len(t, reqPayload.Tools, 1)
	assert.Equal(t, "dummy_tool", reqPayload.Tools[0].Name)
	assert.Equal(t, "A dummy tool", reqPayload.Tools[0].Description)
	assert.JSONEq(t, `{"type":"object"}`, string(reqPayload.Tools[0].InputSchema))

	require.Len(t, reqPayload.Messages, 3)

	assert.Equal(t, "user", reqPayload.Messages[0].Role)
	require.Len(t, reqPayload.Messages[0].Content, 1)
	assert.Equal(t, "text", reqPayload.Messages[0].Content[0].Type)
	assert.Equal(t, "Hello", reqPayload.Messages[0].Content[0].Text)

	assert.Equal(t, "assistant", reqPayload.Messages[1].Role)
	require.Len(t, reqPayload.Messages[1].Content, 1)
	assert.Equal(t, "tool_use", reqPayload.Messages[1].Content[0].Type)
	assert.Equal(t, "toolu_test_123", reqPayload.Messages[1].Content[0].ID)
	assert.Equal(t, "dummy_tool", reqPayload.Messages[1].Content[0].Name)
	assert.JSONEq(t, `{"text":"hello"}`, string(reqPayload.Messages[1].Content[0].Input))

	assert.Equal(t, "user", reqPayload.Messages[2].Role)
	require.Len(t, reqPayload.Messages[2].Content, 1)
	assert.Equal(t, "tool_result", reqPayload.Messages[2].Content[0].Type)
	assert.Equal(t, "toolu_test_123", reqPayload.Messages[2].Content[0].ToolUseID)
	assert.Equal(t, `{"length": 5}`, reqPayload.Messages[2].Content[0].Content)
	assert.False(t, reqPayload.Messages[2].Content[0].IsError)
}

type mockTool struct {
	name        reasoning.ToolName
	description string
	schema      string
}

func (m mockTool) Name() reasoning.ToolName { return m.name }
func (m mockTool) Description() string     { return m.description }
func (m mockTool) InputSchema() string     { return m.schema }
func (m mockTool) Execute(ctx context.Context, inputJSON string) (string, error) {
	return "", nil
}

func TestBedrock_Invoke(t *testing.T) {
	var capturedRequestBody []byte
	mockClient := &http.Client{
		Transport: mockRoundTripper(func(req *http.Request) (*http.Response, error) {
			var err error
			capturedRequestBody, err = io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}

			resPayload := bedrockResponse{
				ID:   "msg_abc",
				Type: "message",
				Role: "assistant",
				Content: []bedrockResponseContent{
					{
						Type: "text",
						Text: "Hello there!",
					},
				},
				StopReason: "end_turn",
			}
			resBytes, err := json.Marshal(resPayload)
			if err != nil {
				return nil, err
			}

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewReader(resBytes)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	awsConf := aws.Config{
		HTTPClient:  mockClient,
		Region:      "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider("MOCK_KEY", "MOCK_SECRET", "MOCK_TOKEN"),
	}
	client := bedrockruntime.NewFromConfig(awsConf)
	b := &Bedrock{
		modelID: "anthropic.claude-3-sonnet-v1:0",
		client:  client,
	}

	ctx := context.Background()
	ans, err := b.Invoke(ctx, "Hello")
	require.NoError(t, err)

	assert.Equal(t, "Hello there!", ans)

	var reqPayload bedrockRequest
	err = json.Unmarshal(capturedRequestBody, &reqPayload)
	require.NoError(t, err)

	assert.Equal(t, "bedrock-2023-05-31", reqPayload.AnthropicVersion)
	assert.Empty(t, reqPayload.System)
	assert.Empty(t, reqPayload.Tools)

	require.Len(t, reqPayload.Messages, 1)
	assert.Equal(t, "user", reqPayload.Messages[0].Role)
	require.Len(t, reqPayload.Messages[0].Content, 1)
	assert.Equal(t, "text", reqPayload.Messages[0].Content[0].Type)
	assert.Equal(t, "Hello", reqPayload.Messages[0].Content[0].Text)
}
