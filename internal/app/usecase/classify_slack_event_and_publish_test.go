package usecase

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/samber/lo"
	"github.com/slack-go/slack/slackevents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func unixNanoToSlackID(nanos int64) string {
	seconds := float64(nanos) / 1e9
	truncated := float64(int64(seconds*1e6)) / 1e6
	return fmt.Sprintf("%f", truncated)
}

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

	mockPublisher := &mockEventPublisher{}
	uc := NewClassifySlackEventAndPublish(mockPublisher)

	// Act

	output, err := uc.Run(ctx, input)

	// Assert

	require.NoError(t, err)

	expect := &ClassifySlackEventAndPublishOutput{
		Challenge: "challenge_value",
	}
	assert.Equal(t, expect, output)
	assert.Equal(t, len(mockPublisher.Published), 0)
}

func Test_ClassifySlackEventAndPublish_handleAppMention(t *testing.T) {
	ctx := t.Context()

	// Arrange

	now := time.Now()
	tsStr := unixNanoToSlackID(now.UnixNano())

	sev := slackevents.EventsAPIEvent{
		Type: slackevents.CallbackEvent,
		InnerEvent: slackevents.EventsAPIInnerEvent{
			Data: &slackevents.AppMentionEvent{
				Text:            "hello, otomo!",
				Channel:         "C1234",
				ThreadTimeStamp: tsStr,
				EventTimeStamp:  tsStr,
				TimeStamp:       tsStr,
			},
		},
	}
	body := lo.Must(json.Marshal(sev))
	input := ClassifySlackEventAndPublishInput{
		Event:   sev,
		RawBody: body,
	}

	mockPublisher := &mockEventPublisher{}
	uc := NewClassifySlackEventAndPublish(mockPublisher)

	// Act

	output, err := uc.Run(ctx, input)

	// Assert

	require.NoError(t, err)

	expect := &ClassifySlackEventAndPublishOutput{
		Status: "ok",
	}
	assert.Equal(t, expect, output)
	require.Equal(t, 1, len(mockPublisher.Published))

	data, ok := mockPublisher.Published[0].Data().(*chat.InstructionReceivedData)
	require.True(t, ok)

	assert.Equal(t, lo.Must(core.NewMessageID(tsStr)), data.MessageID())
	assert.Equal(t, lo.Must(chat.NewThreadID(tsStr)), data.ThreadID())
	assert.Equal(t, chat.RawInstruction("hello, otomo!"), data.RawInstruction())
	assert.Equal(t, now.Unix(), data.SentAt().Unix())
}

func Test_ClassifySlackEventAndPublish_handleAppMention_InvalidTimestamp(t *testing.T) {
	ctx := t.Context()

	sev := slackevents.EventsAPIEvent{
		Type: slackevents.CallbackEvent,
		InnerEvent: slackevents.EventsAPIInnerEvent{
			Data: &slackevents.AppMentionEvent{
				Text:            "hello, otomo!",
				Channel:         "C1234",
				ThreadTimeStamp: "invalid-timestamp",
				EventTimeStamp:  "invalid-timestamp",
				TimeStamp:       "invalid-timestamp",
			},
		},
	}
	body := lo.Must(json.Marshal(sev))
	input := ClassifySlackEventAndPublishInput{
		Event:   sev,
		RawBody: body,
	}

	mockPublisher := &mockEventPublisher{}
	uc := NewClassifySlackEventAndPublish(mockPublisher)

	_, err := uc.Run(ctx, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid timestamp format")
}
