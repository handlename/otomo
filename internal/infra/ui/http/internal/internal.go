package internal

import (
	"context"
	"net/http"

	"github.com/mackee/tanukirpc"
	"github.com/samber/lo"
)

func New(ctx context.Context, prefix string) http.Handler {
	reg := lo.Must(NewRegistry(ctx))
	r := tanukirpc.NewRouter(reg)
	r.Post(prefix+"/slack/event", tanukirpc.NewHandler(slackEventHandler))
	return r
}
