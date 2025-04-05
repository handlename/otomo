package http

import (
	"context"

	"connectrpc.com/connect"
	servicev1 "github.com/handlename/otomo/internal/proto/service/v1"
	"github.com/handlename/otomo/internal/proto/service/v1/servicev1connect"
)

var _ servicev1connect.SlackHandler = (*SlackHandler)(nil)

type SlackHandler struct {
}

// Challenge implements servicev1connect.SlackHandler.
func (s *SlackHandler) Challenge(context.Context, *connect.Request[servicev1.SlackChallengeRequest]) (*connect.Response[servicev1.SlackChallengeResponse], error) {
	panic("unimplemented")
}
