package tool_test

import (
	"context"
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

func TestWebFetchTool_Metadata(t *testing.T) {
	tTool := tool.NewWebFetchTool(config.WebFetch{})
	expectedName, err := reasoning.NewToolName("web_fetch")
	require.NoError(t, err)
	assert.Equal(t, expectedName, tTool.Name())
	assert.Contains(t, tTool.Description(), "Fetch content")
	assert.Contains(t, tTool.InputSchema(), "url")
}

func TestWebFetchTool_Execute(t *testing.T) {
	tests := []struct {
		name              string
		inputJSON         string
		whitelistPatterns []string
		mockStatus        int
		mockContentType   string
		mockResp          string
		expectedOut       string
		expectErr         bool
		errCode           errorcode.ErrorCode
	}{
		{
			name:      "error invalid empty URL",
			inputJSON: `{"url":""}`,
			expectErr: true,
			errCode:   errorcode.ErrInvalidArgument,
		},
		{
			name:              "error not matching whitelist",
			inputJSON:         `{"url":"https://forbidden.com/page"}`,
			whitelistPatterns: []string{`^https://allowed\.com/.*`},
			expectErr:         true,
			errCode:           errorcode.ErrInvalidArgument,
		},
		{
			name:              "success HTML convert to Markdown",
			inputJSON:         `{"url":"https://allowed.com/page"}`,
			whitelistPatterns: []string{`^https://allowed\.com/.*`},
			mockStatus:        http.StatusOK,
			mockContentType:   "text/html; charset=utf-8",
			mockResp:          `<html><body><h1>Hello World</h1><p>Test description.</p></body></html>`,
			expectedOut:       "# Hello World\n\nTest description.",
			expectErr:         false,
		},
		{
			name:            "success plain text returns directly",
			inputJSON:       `{"url":"https://allowed.com/page"}`,
			mockStatus:      http.StatusOK,
			mockContentType: "text/plain",
			mockResp:        "Raw text file contents.",
			expectedOut:     "Raw text file contents.",
			expectErr:       false,
		},
		{
			name:            "error invalid Content-Type",
			inputJSON:       `{"url":"https://allowed.com/page"}`,
			mockStatus:      http.StatusOK,
			mockContentType: "application/pdf",
			mockResp:        "%PDF-1.4 ...",
			expectErr:       true,
			errCode:         errorcode.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", tt.mockContentType)
				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			// If input is using dummy URL, map it to server URL for HTTP call
			actualInput := tt.inputJSON
			if !tt.expectErr || tt.name == "success HTML convert to Markdown" || tt.name == "success plain text returns directly" || tt.name == "error invalid Content-Type" {
				actualInput = `{"url":"` + server.URL + `"}`
			}

			// Rewrite whitelist to include mock server URL if necessary
			whitelist := tt.whitelistPatterns
			if len(whitelist) > 0 && tt.name != "error not matching whitelist" {
				whitelist = []string{"^" + server.URL + ".*"}
			}

			cfg := config.WebFetch{WhitelistPatterns: whitelist}
			tTool := tool.NewWebFetchTool(cfg)

			res, err := tTool.Execute(context.Background(), actualInput)
			if tt.expectErr {
				assert.Error(t, err)
				assert.True(t, failure.Is(err, tt.errCode), "expected code %v, got %v", tt.errCode, err)
			} else {
				require.NoError(t, err)
				assert.Contains(t, res, tt.expectedOut)
			}
		})
	}
}
