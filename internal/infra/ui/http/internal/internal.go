package internal

import (
	"context"

	"github.com/mackee/tanukirpc"
)

func New(ctx context.Context, prefix string) *tanukirpc.Router[*registry] {
	reg := NewRegistry(ctx)
	r := tanukirpc.NewRouter(reg)
	r.Post(prefix+"/reply", tanukirpc.NewHandler(replyHandler))
	return r
}
