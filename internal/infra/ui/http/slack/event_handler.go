package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/handlename/otomo/internal/app/usecase"
	"github.com/handlename/otomo/internal/infra/ui/http/middleware"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack/slackevents"
)

type eventRequest struct {
}

type eventResponse struct {
	Message string `json:"message"`
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	reg, err := middleware.GetRegistry[*registry](r.Context())
	if err != nil {
		// TODO: write error
		panic(err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: write error
		panic(err)
	}
	defer r.Body.Close()

	event, err := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionNoVerifyToken(),
	)
	if err != nil {
		log.Info().Err(err).Msg("invalid event")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invilid event")
		return
	}

	uc := usecase.NewClassifySlackEventAndPublish(reg.Publisher)
	out, err := uc.Run(r.Context(), usecase.ClassifySlackEventAndPublishInput{
		Event:   event,
		RawBody: body,
	})

	res, err := json.Marshal(out)
	if err != nil {
		// TODO: write error
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
