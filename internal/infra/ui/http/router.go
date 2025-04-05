package http

import (
	"context"
	"net/http"

	"github.com/handlename/otomo/internal/infra/repository"
	"github.com/handlename/otomo/internal/proto/service/v1/servicev1connect"
)

func NewMux(ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(servicev1connect.NewLocalHandler(&LocalHandler{
		RepoSession: &repository.VolatileSession{},
		RepoBrain:   repository.NewGeneralBrain(ctx),
	}))
	mux.Handle(servicev1connect.NewSlackHandler(&servicev1connect.UnimplementedSlackHandler{}))
	return mux
}
