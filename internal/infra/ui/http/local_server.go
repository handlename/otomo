package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/handlename/otomo/internal/app/usecase"
	"github.com/handlename/otomo/internal/domain/entity"
	drepo "github.com/handlename/otomo/internal/domain/repository"
	"github.com/handlename/otomo/internal/infra/service"
)

type LocalHandler struct {
	RepoSession drepo.Session
	RepoBrain   drepo.Brain
}

type replyRequest struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}

type replyResponse struct {
	Message string `json:"message"`
}

func (h *LocalHandler) Reply(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: response error
		panic(err)
	}

	var req replyRequest
	if err := json.Unmarshal(body, &req); err != nil {
		// TODO: response error
		panic(err)
	}

	msgr := &service.NopSlack{}

	brain, err := h.RepoBrain.New(ctx)
	if err != nil {
		// TODO: response error
		panic(err)
	}

	otomo := entity.NewOtomo(brain)
	inst := entity.NewInstruction("dummy", req.Message)
	uc := usecase.NewReplyToUser(h.RepoSession, msgr)
	if err := uc.Run(ctx, otomo, inst); err != nil {
		// TODO: response error
		panic(err)
	}

	res := replyResponse{
		Message: msgr.Memory,
	}
	resb, err := json.Marshal(&res)
	if err != nil {
		// TODO: response error
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resb)
}
