package http

import (
	"context"
	"net/http"

	"github.com/handlename/otomo/internal/infra/ui/http/internal"
	"github.com/handlename/otomo/internal/infra/ui/http/slack"
)

func NewMux(ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()
	handle(ctx, mux, "/slack/", slack.New)
	handle(ctx, mux, "/internal/", internal.New)
	return mux
}
