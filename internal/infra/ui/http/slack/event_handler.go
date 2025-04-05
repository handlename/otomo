package slack

import (
	"encoding/json"
	"net/http"

	"github.com/handlename/otomo/internal/infra/ui/http/middleware"
	"github.com/rs/zerolog/log"
)

type eventRequest struct {
}

type eventResponse struct {
	Message string `json:"message"`
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req eventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// TODO: write error
		panic(err)
	}

	reg, err := middleware.GetRegistry[*registry](r.Context())
	if err != nil {
		// TODO: write error
		panic(err)
	}
	log.Debug().Str("Dummy", reg.Dummy).Msg("check registry")

	// TODO: run usecase

	res := eventResponse{
		Message: "hello",
	}
	resb, err := json.Marshal(res)
	if err != nil {
		// TODO: write error
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resb)
}
