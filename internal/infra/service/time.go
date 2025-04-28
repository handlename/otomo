package service

import (
	"fmt"
)

type timeService struct{}

// UnixNanoToSeconds converts UnixNano time to seconds, its precision is microseconds.
func (t *timeService) UnixNanoToSeconds(from int64) float64 {
	result := float64(from) / 1e9
	return float64(int64(result*1e6)) / 1e6
}

func (t *timeService) UnixNanoToSlackID(from int64) string {
	return fmt.Sprintf("%f", t.UnixNanoToSeconds(from))
}

var Time = timeService{}
