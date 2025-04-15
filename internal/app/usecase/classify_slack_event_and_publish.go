package usecase

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/handlename/otomo/internal/domain/event"
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
	publisher event.Publisher
}

func NewClassifySlackEventAndPublish(publisher event.Publisher) *ClassifySlackEventAndPublish {
	return &ClassifySlackEventAndPublish{
		publisher: publisher,
	}
}

// Run implements Usecase.
func (u *ClassifySlackEventAndPublish) Run(ctx context.Context, input ClassifySlackEventAndPublishInput) (*ClassifySlackEventAndPublishOutput, error) {
	var ev event.Event

	switch input.Event.Type {
	case slackevents.URLVerification:
		log.Info().Msg("url verification received")

		var challenge *slackevents.ChallengeResponse
		if err := json.Unmarshal(input.RawBody, &challenge); err != nil {
			return nil, failure.Wrap(err)
		}

		return &ClassifySlackEventAndPublishOutput{
			Challenge: challenge.Challenge,
		}, nil
	case slackevents.CallbackEvent:
		mev, ok := input.Event.InnerEvent.Data.(*slackevents.AppMentionEvent)
		if !ok {
			log.Warn().Any("data", input.Event.InnerEvent.Data).Msg("failed to assert as AppMentionEvent")
			return &ClassifySlackEventAndPublishOutput{
				Status: "ok",
			}, nil
		}

		sentAt, err := parseUnixTimestamp(mev.TimeStamp)
		if err != nil {
			log.Warn().Err(err).Str("timestamp", mev.TimeStamp).Msg("failed to parse timestamp")
			return &ClassifySlackEventAndPublishOutput{
				Status: "ok",
			}, nil
		}

		ev = event.Event(event.NewInstructionReceived(event.InstructionReceivedData{
			MessageID:      mev.EventTimeStamp,
			ThreadID:       mev.ThreadTimeStamp,
			RawInstruction: mev.Text,
			SentAt:         *sentAt,
		}))
	default:
		log.Info().Any("slackEvent", input.Event).Msg("slack event parsed")
	}

	if ev != nil {
		if err := u.publisher.Publish(ev); err != nil {
			return nil, failure.Wrap(err, failure.Message("failed to publish event"), failure.Context{
				"event": ev.String(),
			})
		}
		log.Debug().Any("event", ev).Msg("event published")
	}

	return &ClassifySlackEventAndPublishOutput{
		Status: "ok",
	}, nil
}

func parseUnixTimestamp(s string) (*time.Time, error) {
	parts := strings.Split(s, ".")

	sec, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	msec, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	t := time.Unix(sec, msec*1000)
	return &t, nil
}
