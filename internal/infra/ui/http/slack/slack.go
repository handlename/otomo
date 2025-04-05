package slack

import (
	"context"
	"net/http"

	"github.com/handlename/otomo/internal/infra/ui/http/middleware"
)

func New(ctx context.Context) *http.ServeMux {
	reg := NewRegistry(ctx)
	mids := []middleware.Middleware{
		middleware.NewRegistry(reg),
	}
	mux := http.NewServeMux()
	mux.Handle("POST /event", middleware.Wrap(eventHandler, mids...))
	return mux
}
