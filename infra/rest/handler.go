package rest

import (
	"context"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/shared"
)

// NewHandler creates a generic http.HandlerFunc that handle an HTTP request. The prroces is as follows:
// 1. Execute middlewares
// 2. Extract domain input (DI) from HTTP request. If an error occurs, return an HTTP error.
// 3. Execute application handler, DI -> DO. If an error occurs, return an HTTP error.
// 4. Encode domain output (DO) to REST output (RO) (dataEncoder)
// 5. Write REST output (RO) to HTTP response (wrap proper response type)
//
// Type parameters:
//   - RI: REST input struct type
//   - DI: domain input type
//   - DO: domain output type
//   - RO: REST output struct type
//
// Parameters:
//   - handler: application handler function that accepts a context and an input (DI) and returns an output (DO) and an error.
//   - requestDecoder: function that extracts domain input (DI) from HTTP request.
//   - dataEncoder: function that encodes domain output (DO) to REST output (RO).
//   - okStatusCode: HTTP status code to be returned in case of success.
//   - middlewares: variadic parameter of middlewares to be executed before handler.
//
// Developer MAY provide middlewars as variadic parameter. Middlewars will be
// executed before handler in the same order as they are declared
func NewHandler[DO, DI, RO any](
	monitor shared.Monitor,
	appHandler func(context.Context, DI) (DO, error),
	requestDecoder func(*http.Request) (DI, error),
	dataEncoder func(DO) (RO, error),
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

		var out DO
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
