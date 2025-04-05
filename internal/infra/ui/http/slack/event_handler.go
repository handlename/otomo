package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/handlename/otomo/internal/errorcode"
	"github.com/handlename/otomo/internal/infra/ui/http/middleware"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack/slackevents"
)

type eventRequest struct {
}

type eventResponse struct {
	Message string `json:"message"`
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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

	// TODO: move to middleware
	if err := reg.Slack.Verify(r.Header, body); err != nil {
		if failure.Is(err, errorcode.ErrInvalidArgument) {
			log.Debug().Err(err).Msg("invalid signing")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invilid signing")
			return
		}

		// TODO: write error
		panic(err)
	}

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

	switch event.Type {
	case slackevents.URLVerification:
		log.Info().Msg("url verification received")

		var challenge *slackevents.ChallengeResponse
		if err := json.Unmarshal(body, &challenge); err != nil {
			// TODO: write error
			panic(err)
		}

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, challenge.Challenge)
		return
	default:
		log.Info().Any("event", event).Msg("event parsed")
	}

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
