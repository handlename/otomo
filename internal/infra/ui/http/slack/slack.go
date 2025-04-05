package slack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/handlename/otomo/internal/infra/ui/http/middleware"
)

func New(ctx context.Context, prefix string) *http.ServeMux {
	reg := NewRegistry(ctx)
	mids := []middleware.Middleware{
		middleware.NewRegistry(reg),
		middleware.NewAccesslog(),
	}
	mux := http.NewServeMux()
	mux.Handle(fmt.Sprintf("POST %s/event", prefix), middleware.Wrap(eventHandler, mids...))
	return mux
}
