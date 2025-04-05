package slack

import (
	"context"
	"net/http"

	"github.com/handlename/otomo/internal/infra/ui/http/middleware"
)

func New(ctx context.Context, prefix string) http.Handler {
	reg := NewRegistry(ctx)
	mids := []middleware.Middleware{
		middleware.NewRegistry(reg),
		middleware.NewAccesslog(),
	}
	mux := http.NewServeMux()
	mux.Handle("POST "+prefix+"/event", middleware.Wrap(eventHandler, mids...))
	return mux
}
