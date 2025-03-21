package otomo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fujiwara/ridge"
	"github.com/rs/zerolog/log"
)

type App struct {
	appToken string
	botToken string
}

func New() *App {
	return &App{
		appToken: os.Getenv("OTOMO_SLACK_APP_TOKEN"),
		botToken: os.Getenv("OTOMO_SLACK_BOT_TOKEN"),
	}
}

func (a *App) Run(ctx context.Context) error {
	// Lambda function用の処理
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.Handler)
	ridge.RunWithContext(ctx, ":8080", "/", mux)

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

// func (a *App) runSocketMode(ctx context.Context) error {
// 	if a.appToken == "" || a.botToken == "" {
// 		return failure.New(
// 			errcode.ErrInvalidArgument,
// 			failure.Message("both of Slack App token and Slack Bot token are required"),
// 		)
// 	}

// 	api := slack.New(
// 		a.botToken,
// 		slack.OptionAppLevelToken(a.appToken),
// 		slack.OptionLog(SlackLogger{By: "api"}),
// 	)

// 	client := socketmode.New(api, socketmode.OptionDebug(true), socketmode.OptionLog(SlackLogger{By: "client"}))

// 	go func() {
// 		for event := range client.Events {
// 			switch event.Type {
// 			case socketmode.EventTypeEventsAPI:
// 				evt, ok := event.Data.(slackevents.EventsAPIEvent)
// 				if !ok {
// 					continue
// 				}

// 				client.Ack(*event.Request)

// 				switch evt.Type {
// 				case slackevents.CallbackEvent:
// 					switch e := evt.InnerEvent.Data.(type) {
// 					case *slackevents.AppMentionEvent:
// 						mention := strings.TrimSpace(strings.Replace(e.Text, fmt.Sprintf("<@%s>", e.BotID), "", -1))
// 						reply := fmt.Sprintf("「%s」ってどういう意味？", mention)

// 						_, _, err := api.PostMessage(e.Channel, slack.MsgOptionText(reply, false))
// 						if err != nil {
// 							log.Error().Err(err).Msg("error on post message")
// 						}
// 					default:
// 						log.Debug().Any("event", e).Msg("unknown event data")
// 					}
// 				default:
// 					log.Debug().Str("type", string(evt.Type)).Msg("unknown EventsAPI type")
// 				}
// 			default:
// 				log.Debug().Str("type", string(event.Type)).Msg("unknown Events type")
// 			}
// 		}
// 	}()

// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	log.Info().Msg("launching bot...")
// 	err := client.RunContext(ctx)
// 	if err != nil {
// 		return failure.Wrap(err,
// 			failure.WithCode(errcode.ErrInternal),
// 			failure.Message("failed to launch bot"),
// 		)
// 	}

// 	return nil
// }

type SlackLogger struct {
	By string
}

func (l SlackLogger) Output(calldepth int, s string) error {
	log.Info().Str("by", l.By).Msg(s)
	return nil
}
