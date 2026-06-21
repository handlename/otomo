package tool_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
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
		inputJSON         func(*httptest.Server) string
		whitelistPatterns func(*httptest.Server) []string
		mockStatus        int
		mockContentType   string
		mockResp          string
		expectedOut       string
		expectErr         bool
		errCode           errorcode.ErrorCode
	}{
		{
			name: "error invalid empty URL",
			inputJSON: func(s *httptest.Server) string {
				return `{"url":""}`
			},
			expectErr: true,
			errCode:   errorcode.ErrInvalidArgument,
		},
		{
			name: "error not matching whitelist",
			inputJSON: func(s *httptest.Server) string {
				return `{"url":"https://forbidden.com/page"}`
			},
			whitelistPatterns: func(s *httptest.Server) []string {
				return []string{`^https://allowed\.com/.*`}
			},
			expectErr: true,
			errCode:   errorcode.ErrInvalidArgument,
		},
		{
			name: "success HTML convert to Markdown",
			inputJSON: func(s *httptest.Server) string {
				return `{"url":"` + s.URL + `"}`
			},
			whitelistPatterns: func(s *httptest.Server) []string {
				return []string{"^" + s.URL + ".*"}
			},
			mockStatus:      http.StatusOK,
			mockContentType: "text/html; charset=utf-8",
			mockResp:        `<html><body><h1>Hello World</h1><p>Test description.</p></body></html>`,
			expectedOut:     "# Hello World\n\nTest description.",
			expectErr:       false,
		},
		{
			name: "success plain text returns directly",
			inputJSON: func(s *httptest.Server) string {
				return `{"url":"` + s.URL + `"}`
			},
			mockStatus:      http.StatusOK,
			mockContentType: "text/plain",
			mockResp:        "Raw text file contents.",
			expectedOut:     "Raw text file contents.",
			expectErr:       false,
		},
		{
			name: "success JSON returns directly",
			inputJSON: func(s *httptest.Server) string {
				return `{"url":"` + s.URL + `"}`
			},
			mockStatus:      http.StatusOK,
			mockContentType: "application/json",
			mockResp:        `{"status":"success","data":{"id":123}}`,
			expectedOut:     `{"status":"success","data":{"id":123}}`,
			expectErr:       false,
		},
		{
			name: "error invalid Content-Type",
			inputJSON: func(s *httptest.Server) string {
				return `{"url":"` + s.URL + `"}`
			},
			mockStatus:      http.StatusOK,
			mockContentType: "application/pdf",
			mockResp:        "%PDF-1.4 ...",
			expectErr:       true,
			errCode:         errorcode.ErrInvalidArgument,
		},
		{
			name: "error invalid Content-Type CSS",
			inputJSON: func(s *httptest.Server) string {
				return `{"url":"` + s.URL + `"}`
			},
			mockStatus:      http.StatusOK,
			mockContentType: "text/css",
			mockResp:        "body { color: red; }",
			expectErr:       true,
			errCode:         errorcode.ErrInvalidArgument,
		},
		{
			name: "success case-insensitive HTTP scheme",
			inputJSON: func(s *httptest.Server) string {
				return `{"url":"` + strings.Replace(s.URL, "http://", "HTTP://", 1) + `"}`
			},
			mockStatus:      http.StatusOK,
			mockContentType: "text/plain",
			mockResp:        "Raw text file contents.",
			expectedOut:     "Raw text file contents.",
			expectErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", tt.mockContentType)
				w.WriteHeader(tt.mockStatus)
				_, _ = w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			var actualInput string
			if tt.inputJSON != nil {
				actualInput = tt.inputJSON(server)
			}

			var whitelist []string
			if tt.whitelistPatterns != nil {
				whitelist = tt.whitelistPatterns(server)
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
