package http

import (
	"context"
	"net/http"

	"github.com/handlename/otomo/internal/infra/ui/http/local"
	"github.com/handlename/otomo/internal/infra/ui/http/slack"
)

func NewMux(ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/slack/", http.StripPrefix("/slack", slack.New(ctx)))
	mux.Handle("/local/", http.StripPrefix("/local", local.New(ctx)))
	return mux
}
