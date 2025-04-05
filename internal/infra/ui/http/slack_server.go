package http

import (
	"net/http"
)

type SlackHandler struct {
}

func (s *SlackHandler) Event(w http.ResponseWriter, req *http.Request) {
	panic("not implemented")
}
