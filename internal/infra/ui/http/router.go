package http

import (
	"context"
	"net/http"
)

func NewMux(ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/local/", http.StripPrefix("/local", setupLocal(ctx)))
	return mux
}
