package httpmw

import (
	"log/slog"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleware(log *slog.Logger) func(handler http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			recorder := &statusRecorder{
				ResponseWriter: w,
				statusCode: http.StatusOK,
			}
			next.ServeHTTP(recorder, r)
			log.Info(
				"method", r.Method,
				"req", r.URL,
				"duration", time.Since(start),
				"code", recorder.statusCode,
				"err", nil,
			)
		})
	}
}