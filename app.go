package otomo

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/fujiwara/ridge"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

type App struct {
	Port     int
	AppToken string
	BotToken string
}

func (a *App) Init() error {
	if a.Port == 0 {
		a.Port = 8080
	}

	if a.AppToken == "" || a.BotToken == "" {
		return failure.New(
			errorcode.ErrInvalidArgument,
			failure.Message("both of Slack App token and Slack Bot token are required"),
		)
	}

	return nil
}

func (a *App) Run(ctx context.Context) error {
	if err := a.Init(); err != nil {
		return failure.Wrap(err, failure.Message("failed to init app"))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", a.Handler)
	ridge.RunWithContext(ctx, fmt.Sprintf(":%d", a.Port), "/", mux)

	return nil
}

func (a *App) Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "echo %s", body)
}
