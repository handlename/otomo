package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/morikuni/failure/v2"
)

type timeService struct{}

// ParseUnixTimestamp parses a Unix timestamp string into a time.Time object.
func (t *timeService) ParseUnixTimestamp(s string) (*time.Time, error) {
	parts := strings.Split(s, ".")

	sec, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	msec, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	res := time.Unix(sec, msec*1000)
	return &res, nil
}

// UnixNanoToSeconds converts UnixNano time to seconds, its precision is microseconds.
func (t *timeService) UnixNanoToSeconds(from int64) float64 {
	result := float64(from) / 1e9
	return float64(int64(result*1e6)) / 1e6
}

func (t *timeService) UnixNanoToSlackID(from int64) string {
	return fmt.Sprintf("%f", t.UnixNanoToSeconds(from))
}

var Time = timeService{}
