package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

var _ Middleware = (*Accesslog)(nil)

type Accesslog struct{}

func NewAccesslog() *Accesslog {
	return &Accesslog{}
}

// Wrap implements Middleware.
func (a *Accesslog) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		aw := &accesslogWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(aw, req)

		log.Info().
			Str("method", req.Method).
			Str("path", req.URL.Path).
			Int("status", aw.statusCode).
			Float64("req_time", time.Since(start).Seconds()*1000). // in msec
			Int("res_size", aw.size).
			Msg("request")
	})
}

type accesslogWriter struct {
	http.ResponseWriter

	// statusCode stores HTTP statusCode statusCode of response
	statusCode int

	// size stores response body size
	size int
}

func (w *accesslogWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *accesslogWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}
