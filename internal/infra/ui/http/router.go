package http

import (
	"context"
	"net/http"

	"github.com/handlename/otomo/internal/infra/ui/http/local"
	"github.com/handlename/otomo/internal/infra/ui/http/slack"
)

func NewMux(ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()
	handle(ctx, mux, "/slack/", slack.New)
	handle(ctx, mux, "/local/", func(ctx context.Context, prefix string) http.Handler {
		return local.New(ctx, prefix)
	})
	return mux
}
