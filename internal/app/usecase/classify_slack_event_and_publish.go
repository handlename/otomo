package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	appservice "github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack/slackevents"
)

type ClassifySlackEventAndPublishInput struct {
	Event   slackevents.EventsAPIEvent
	RawBody []byte
}

type ClassifySlackEventAndPublishOutput struct {
	Status    string `json:"status"`
	Challenge string `json:"challenge,omitempty"`
}

type ClassifySlackEventAndPublish struct {
	publisher appservice.Publisher
}

func NewClassifySlackEventAndPublish(publisher appservice.Publisher) *ClassifySlackEventAndPublish {
	return &ClassifySlackEventAndPublish{
		publisher: publisher,
	}
}

// Run implements Usecase.
func (u *ClassifySlackEventAndPublish) Run(ctx context.Context, input ClassifySlackEventAndPublishInput) (*ClassifySlackEventAndPublishOutput, error) {
	var ev core.Event

	switch input.Event.Type {
	case slackevents.URLVerification:
		log.Info().Msg("url verification received")
		return u.handleURLVerification(ctx, input)
	case slackevents.CallbackEvent:
		log.Info().Msg("callback event received")
		var err error
		ev, err = u.handleCallbackEvent(ctx, input)
		if err != nil {
			return nil, failure.Wrap(err, failure.Message("failed to handle callback event"))
		}
	default:
		log.Info().Any("slackEvent", input.Event).Msg("slack event parsed")
	}

	if ev != nil {
		if err := u.publisher.Publish(ctx, ev); err != nil {
			return nil, failure.Wrap(err, failure.Message("failed to publish event"), failure.Context{
				"event": ev.String(),
			})
		}
		log.Debug().Str("event", ev.String()).Msg("event published")
	}

	return &ClassifySlackEventAndPublishOutput{
		Status: "ok",
	}, nil
}

func (u *ClassifySlackEventAndPublish) handleURLVerification(_ context.Context, input ClassifySlackEventAndPublishInput) (*ClassifySlackEventAndPublishOutput, error) {
	var challenge *slackevents.ChallengeResponse
	if err := json.Unmarshal(input.RawBody, &challenge); err != nil {
		return nil, failure.Wrap(err)
	}

	return &ClassifySlackEventAndPublishOutput{
		Challenge: challenge.Challenge,
	}, nil
}

func (u *ClassifySlackEventAndPublish) handleCallbackEvent(_ context.Context, input ClassifySlackEventAndPublishInput) (core.Event, error) {
	switch iev := input.Event.InnerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		sentAt, err := service.Time.ParseUnixTimestamp(iev.TimeStamp)
		if err != nil {
			return nil, failure.Wrap(err,
				failure.Message("failed to parse timestamp"),
				failure.Context{
					"timestamp": iev.TimeStamp,
				},
			)
		}

		ev, err := chat.NewInstructionReceived(chat.InstructionReceivedData{
			ChannelID:      iev.Channel,
			MessageID:      iev.EventTimeStamp,
			ThreadID:       iev.ThreadTimeStamp,
			RawInstruction: iev.Text,
			SentAt:         *sentAt,
		})
		if err != nil {
			return nil, failure.Wrap(err)
		}
		return ev, nil
	default:
		return nil, failure.New(
			"failed to assert Slack inner event",
			failure.Context{
				"data": fmt.Sprintf("%+v", input.Event.InnerEvent.Data),
			},
		)
	}
}
