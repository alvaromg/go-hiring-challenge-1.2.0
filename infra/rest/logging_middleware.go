package rest

import (
	"net/http"
	"time"

	"github.com/mytheresa/go-hiring-challenge/shared"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// NewLoggingMiddleware logs method, URL (with query params), response status
// code and execution time for every request.
func NewLoggingMiddleware(logger shared.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rec, r)

			logger.WithContext(r.Context()).WithFields(map[string]any{
				"method": r.Method,
				"uri":    r.URL.RequestURI(),
				"status": rec.statusCode,
				"ms":     time.Since(start).Milliseconds(),
			}).
				Infof("http request")
		})
	}
}
