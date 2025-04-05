package http

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type handlerBuilder func(ctx context.Context, prefix string) http.Handler

func handle(ctx context.Context, mux *http.ServeMux, pattern string, builder handlerBuilder) {
	if !strings.HasPrefix(pattern, "/") {
		log.Panic("pattern must begin with '/' but", pattern)
	}

	if !strings.HasSuffix(pattern, "/") {
		log.Panic("pattern must end with '/' but", pattern)
	}

	mux.Handle(pattern, builder(ctx, strings.TrimSuffix(pattern, "/")))
}
