package http

import (
	"context"
	"net/http"

	"github.com/handlename/otomo/internal/infra/ui/http/local"
)

func NewMux(ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/local/", http.StripPrefix("/local", local.New(ctx)))
	return mux
}
