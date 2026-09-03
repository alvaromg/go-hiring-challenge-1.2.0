package rest

import (
	"context"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/domain/list"
	"github.com/mytheresa/go-hiring-challenge/domain/query"
	"github.com/mytheresa/go-hiring-challenge/shared"
)

// NewListByQueryHandle creates a generic http.HandlerFunc that handle an HTTP request. Similar to NewHandler but in this case input is a
// query (filters + sort + pagination) and response is a generic list response
func NewListByQueryHandle[DO, RO, WD any](
	monitor shared.Monitor,
	appHandler func(context.Context, *query.Query) (list.ListResponse[DO], error),
	itemEncoder func(DO) RO,
	dataWrapper func([]RO) WD,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := startHandlerSpan(monitor, r)
		statusCode := http.StatusOK
		var err error
		defer func() { endHandlerSpan(span, statusCode, err) }()

		q, err := DecodeQueryFromRequest(r)
		if err != nil {
			statusCode = errToHTTPCode(err)
			HandleHTTPError(monitor.Logger(), w, r, err)
			return
		}

		var out list.ListResponse[DO]
		out, err = appHandler(ctx, q)
		if err != nil {
			statusCode = errToHTTPCode(err)
			HandleHTTPError(monitor.Logger(), w, r, err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		EncodeListResponse(monitor.Logger(), ctx, itemEncoder, dataWrapper, out, q, w, r)
	}
}
