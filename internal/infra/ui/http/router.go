package http

import (
	"context"
	"net/http"

	"github.com/handlename/otomo/internal/infra/repository"
	"github.com/handlename/otomo/internal/proto/service/v1/servicev1connect"
)

func NewMux(ctx context.Context, botToken, appToken string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(servicev1connect.NewLocalHandler(&LocalHandler{
		RepoSession: &repository.VolatileSession{},
		RepoBrain:   repository.NewGeneralBrain(ctx, botToken, appToken),
	}))
	return mux
}
