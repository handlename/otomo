package tool

import (
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
	"github.com/samber/lo"
)

var _ reasoning.Tool = (*WebFetchTool)(nil)

type WebFetchTool struct {
	name     reasoning.ToolName
	cfg      config.WebFetch
	compiler *converter.Converter
	client   *http.Client
}

func NewWebFetchTool(cfg config.WebFetch) *WebFetchTool {
	return &WebFetchTool{
		name: lo.Must(reasoning.NewToolName("web_fetch")),
		cfg:  cfg,
		compiler: converter.NewConverter(
			converter.WithPlugins(
				base.NewBasePlugin(),
				commonmark.NewCommonmarkPlugin(),
			),
		),
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (t *WebFetchTool) Name() reasoning.ToolName {
	return t.name
}

func (t *WebFetchTool) Description() string {
	return "Fetch content from a URL and return it as Markdown or plain text."
}

func (t *WebFetchTool) InputSchema() string {
	return `{
		"type": "object",
		"properties": {
			"url": {
				"type": "string",
				"description": "The URL of the web page to fetch contents from."
			}
		},
		"required": ["url"]
	}`
}

func (t *WebFetchTool) Execute(ctx context.Context, inputJSON string) (string, error) {
	var input struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", failure.Wrap(err, failure.WithCode(errorcode.ErrInvalidArgument), failure.Message("failed to parse input parameters"))
	}

	if input.URL == "" {
		return "", failure.New(errorcode.ErrInvalidArgument, failure.Message("url parameter is required"))
	}

	parsedURL, err := url.Parse(input.URL)
	if err != nil {
		return "", failure.Wrap(err, failure.WithCode(errorcode.ErrInvalidArgument), failure.Message("invalid URL format"))
	}
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", failure.New(errorcode.ErrInvalidArgument, failure.Message("invalid URL scheme (http/https only)"))
	}

	// Whitelist check
	if len(t.cfg.WhitelistPatterns) > 0 {
		matched := false
		for _, pattern := range t.cfg.WhitelistPatterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return "", failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Messagef("failed to compile whitelist regex pattern: %s", pattern))
			}
			if re.MatchString(input.URL) {
				matched = true
				break
			}
		}
		if !matched {
			return "", failure.New(errorcode.ErrInvalidArgument, failure.Messagef("URL is not allowed by whitelist: %s", input.URL))
		}
	}

	// Fetch page
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, input.URL, nil)
	if err != nil {
		return "", failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to create http request"))
	}
	req.Header.Set("User-Agent", UserAgent())

	resp, err := t.client.Do(req)
	if err != nil {
		return "", failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to fetch URL content"))
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", failure.New(errorcode.ErrInternal, failure.Messagef("failed to fetch URL, server returned status: %d", resp.StatusCode))
	}

	contentType := resp.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		mediaType = "text/plain" // fallback
	}

	switch mediaType {
	case "text/html":
		mdBytes, err := t.compiler.ConvertReader(resp.Body)
		if err != nil {
			return "", failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to convert HTML content to markdown"))
		}
		return string(mdBytes), nil
	case "text/plain", "application/json":
		// Read body directly
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to read response body"))
		}
		return string(bodyBytes), nil
	default:
		return "", failure.New(errorcode.ErrInvalidArgument, failure.Messagef("unsupported content type: %s", mediaType))
	}
}
