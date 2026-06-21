package tool_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/handlename/otomo/internal/infra/tool"
	"github.com/morikuni/failure/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebSearchTool_Metadata(t *testing.T) {
	cfg := config.WebSearch{TavilyAPIKey: "dummy"}
	tTool := tool.NewWebSearchTool(cfg)

	expectedName, err := reasoning.NewToolName("web_search")
	require.NoError(t, err)
	assert.Equal(t, expectedName, tTool.Name())
	assert.Contains(t, tTool.Description(), "Search the web")
	assert.Contains(t, tTool.InputSchema(), "query")
}

func TestWebSearchTool_Execute(t *testing.T) {
	tests := []struct {
		name       string
		apiKey     string
		inputJSON  string
		mockStatus int
		mockResp   string
		expectErr  bool
		errCode    errorcode.ErrorCode
	}{
		{
			name:      "error empty query",
			apiKey:    "key",
			inputJSON: `{"query":""}`,
			expectErr: true,
			errCode:   errorcode.ErrInvalidArgument,
		},
		{
			name:      "error missing API key",
			apiKey:    "",
			inputJSON: `{"query":"test"}`,
			expectErr: true,
			errCode:   errorcode.ErrInvalidArgument,
		},
		{
			name:       "success search",
			apiKey:     "valid-key",
			inputJSON:  `{"query":"golang"}`,
			mockStatus: http.StatusOK,
			mockResp: `{
				"results": [
					{"title": "Go Programming Language", "url": "https://go.dev", "content": "The Go home page."}
				]
			}`,
			expectErr: false,
		},
		{
			name:       "error API returns 500",
			apiKey:     "valid-key",
			inputJSON:  `{"query":"golang"}`,
			mockStatus: http.StatusInternalServerError,
			mockResp:   `{"error": "internal error"}`,
			expectErr:  true,
			errCode:    errorcode.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				var body map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&body)
				require.NoError(t, err)
				assert.Equal(t, tt.apiKey, body["api_key"])
				assert.Equal(t, "basic", body["search_depth"])

				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			cfg := config.WebSearch{TavilyAPIKey: tt.apiKey}
			tTool := tool.NewWebSearchToolWithEndpoint(cfg, server.URL)

			res, err := tTool.Execute(context.Background(), tt.inputJSON)
			if tt.expectErr {
				assert.Error(t, err)
				assert.True(t, failure.Is(err, tt.errCode), "expected code %v, got %v", tt.errCode, err)
			} else {
				require.NoError(t, err)
				var results map[string]interface{}
				err = json.Unmarshal([]byte(res), &results)
				require.NoError(t, err)
				resList := results["results"].([]interface{})
				assert.Len(t, resList, 1)
			}
		})
	}
}
