package http

import (
	"context"
	"net/http"

	"github.com/handlename/otomo/internal/infra/repository"
)

func NewMux(ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()
	
	localHandler := &LocalHandler{
		RepoSession: &repository.VolatileSession{},
		RepoBrain:   repository.NewGeneralBrain(ctx),
	}
	mux.HandleFunc("POST /local/reply", localHandler.Reply)
	
	slackHandler := &SlackHandler{}
	mux.HandleFunc("POST /slack/event", slackHandler.Event)
	
	return mux
}
