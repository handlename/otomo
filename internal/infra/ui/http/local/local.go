package local

import (
	"context"

	"github.com/mackee/tanukirpc"
)

func New(ctx context.Context) *tanukirpc.Router[*registry] {
	reg := NewRegistry(ctx)
	r := tanukirpc.NewRouter(reg)
	r.Post("/reply", tanukirpc.NewHandler(replyHandler))
	return r
}
