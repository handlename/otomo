package service

import "fmt"

type time struct{}

// UnixNanoToSeconds converts UnixNano time to seconds, its precision is microseconds.
func (t *time) UnixNanoToSeconds(from int64) float64 {
	result := float64(from) / 1e9
	return float64(int64(result*1e6)) / 1e6
}

func (t *time) UnixNanoToSlackID(from int64) string {
	return fmt.Sprintf("%f", t.UnixNanoToSeconds(from))
}

var Time = time{}
