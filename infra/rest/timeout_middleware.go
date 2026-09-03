package rest

import (
	"net/http"
	"time"
)

// timeoutMessage is a JSON body so clients get a response shaped like every
// other error, even though the deadline fires before a handler can build one.
const timeoutMessage = `{"error":{"message":"request timed out"}}`

// NewTimeoutMiddleware aborts a request and responds with 503 if the wrapped
// handler takes longer than d to complete. The handler's context is
// cancelled when the deadline is hit, so downstream calls (e.g. DB queries)
// that respect ctx are cancelled too.
func NewTimeoutMiddleware(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, d, timeoutMessage)
	}
}
