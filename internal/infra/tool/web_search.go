package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
	"github.com/samber/lo"
)

var _ reasoning.Tool = (*WebSearchTool)(nil)

type WebSearchTool struct {
	name     reasoning.ToolName
	cfg      config.WebSearch
	endpoint string
	client   *http.Client
}

func NewWebSearchTool(cfg config.WebSearch) *WebSearchTool {
	return NewWebSearchToolWithEndpoint(cfg, "https://api.tavily.com/search")
}

func NewWebSearchToolWithEndpoint(cfg config.WebSearch, endpoint string) *WebSearchTool {
	return &WebSearchTool{
		name:     lo.Must(reasoning.NewToolName("web_search")),
		cfg:      cfg,
		endpoint: endpoint,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (t *WebSearchTool) Name() reasoning.ToolName {
	return t.name
}

func (t *WebSearchTool) Description() string {
	return "Search the web for the given query using Tavily Search API and return results."
}

func (t *WebSearchTool) InputSchema() string {
	return `{
		"type": "object",
		"properties": {
			"query": {
				"type": "string",
				"description": "The search query to look up on the internet."
			}
		},
		"required": ["query"]
	}`
}

func (t *WebSearchTool) Execute(ctx context.Context, inputJSON string) (string, error) {
	var input struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", failure.Wrap(err, failure.WithCode(errorcode.ErrInvalidArgument), failure.Message("failed to parse input parameters"))
	}

	if input.Query == "" {
		return "", failure.New(errorcode.ErrInvalidArgument, failure.Message("query parameter is required"))
	}

	if t.cfg.TavilyAPIKey == "" {
		return "", failure.New(errorcode.ErrInvalidArgument, failure.Message("tavily_api_key is not configured"))
	}

	reqBody, err := json.Marshal(map[string]interface{}{
		"api_key":      t.cfg.TavilyAPIKey,
		"query":        input.Query,
		"search_depth": "basic",
	})
	if err != nil {
		return "", failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to marshal request body"))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to create http request"))
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return "", failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to call tavily search api"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", failure.New(errorcode.ErrInternal, failure.Messagef("tavily search api returned non-ok status: %d", resp.StatusCode))
	}

	var apiResp struct {
		Results []struct {
			Title   string  `json:"title"`
			URL     string  `json:"url"`
			Content string  `json:"content"`
			Score   float64 `json:"score"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to decode tavily search response"))
	}

	respBytes, err := json.Marshal(apiResp)
	if err != nil {
		return "", failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to marshal search response"))
	}

	return string(respBytes), nil
}
