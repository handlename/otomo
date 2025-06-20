package usecase

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/handlename/otomo/internal/domain/event"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/handlename/otomo/internal/testutil"
	"github.com/samber/lo"
	"github.com/slack-go/slack/slackevents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ClassifySlackEventAndPublish_handleChallenge(t *testing.T) {
	ctx := t.Context()

	// Arrange

	input := ClassifySlackEventAndPublishInput{
		Event: slackevents.EventsAPIEvent{
			Type: slackevents.URLVerification,
		},
		RawBody: lo.Must(json.Marshal(slackevents.ChallengeResponse{
			Challenge: "challenge_value",
		})),
	}

	mockPublisher := testutil.NewMockPublisher()
	uc := NewClassifySlackEventAndPublish(mockPublisher)

	// Act

	output, err := uc.Run(ctx, input)

	// Assert

	require.NoError(t, err)

	expect := &ClassifySlackEventAndPublishOutput{
		Challenge: "challenge_value",
	}
	assert.Equal(t, expect, output)
	assert.Equal(t, len(mockPublisher.History), 0)
}

func Test_ClassifySlackEventAndPublish_handleAppMention(t *testing.T) {
	ctx := t.Context()

	// Arrange

	now := time.Now()

	sev := slackevents.EventsAPIEvent{
		Type: slackevents.CallbackEvent,
		InnerEvent: slackevents.EventsAPIInnerEvent{
			Data: &slackevents.AppMentionEvent{
				Text:            "hello, otomo!",
				Channel:         "C1234",
				ThreadTimeStamp: service.Time.UnixNanoToSlackID(now.UnixNano()),
				EventTimeStamp:  service.Time.UnixNanoToSlackID(now.UnixNano()),
				TimeStamp:       service.Time.UnixNanoToSlackID(now.UnixNano()),
			},
		},
	}
	body := lo.Must(json.Marshal(sev))
	input := ClassifySlackEventAndPublishInput{
		Event:   sev,
		RawBody: body,
	}

	mockPublisher := testutil.NewMockPublisher()
	uc := NewClassifySlackEventAndPublish(mockPublisher)

	// Act

	output, err := uc.Run(ctx, input)

	// Assert

	require.NoError(t, err)

	expect := &ClassifySlackEventAndPublishOutput{
		Status: "ok",
	}
	assert.Equal(t, expect, output)
	require.Equal(t, 1, len(mockPublisher.History))

	data, ok := mockPublisher.History[0].Data().(event.InstructionReceivedData)
	require.True(t, ok)

	assert.Equal(t, service.Time.UnixNanoToSlackID(now.UnixNano()), data.MessageID)
	assert.Equal(t, service.Time.UnixNanoToSlackID(now.UnixNano()), data.ThreadID)
	assert.Equal(t, "hello, otomo!", data.RawInstruction)
	assert.Equal(t, now.Unix(), data.SentAt.Unix())
}
