package rest

import (
	"context"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/shared"
)

// NewHandler creates a generic http.HandlerFunc that handle an HTTP request. The prroces is as follows:
// 1. Execute middlewares
// 2. Extract application input (AI) from HTTP request. If an error occurs, return an HTTP error.
// 3. Execute application handler, AI -> AO. If an error occurs, return an HTTP error.
// 4. Encode application output (AO) to REST output (RO) (dataEncoder)
// 5. Write REST output (RO) to HTTP response (wrap proper response type)
//
// Type parameters:
//   - RI: REST input struct type
//   - AI: application input type
//   - AO: application output type
//   - RO: REST output struct type
//
// Parameters:
//   - handler: application handler function that accepts a context and an input (AI) and returns an output (AO) and an error.
//   - requestDecoder: function that extracts application input (AI) from HTTP request.
//   - dataEncoder: function that encodes application output (AO) to REST output (RO).
//   - okStatusCode: HTTP status code to be returned in case of success.
//
// Developer MAY provide middlewars as variadic parameter. Middlewars will be
// executed before handler in the same order as they are declared
func NewHandler[AO, AI, RO any](
	monitor shared.Monitor,
	appHandler func(context.Context, AI) (AO, error),
	requestDecoder func(*http.Request) (AI, error),
	dataEncoder func(AO) (RO, error),
	okStatusCode int,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := startHandlerSpan(monitor, r)
		statusCode := okStatusCode
		var err error
		defer func() { endHandlerSpan(span, statusCode, err) }()

		decodedRequest, err := requestDecoder(r)
		if err != nil {
			statusCode = errToHTTPCode(err)
			HandleHTTPError(monitor.Logger(), w, r, err)
			return
		}

		var out AO
		out, err = appHandler(ctx, decodedRequest)
		if err != nil {
			statusCode = errToHTTPCode(err)
			HandleHTTPError(monitor.Logger(), w, r, err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(okStatusCode)
		encodeHandlerResponse(monitor.Logger(), ctx, dataEncoder, out, w, r)
	}
}
