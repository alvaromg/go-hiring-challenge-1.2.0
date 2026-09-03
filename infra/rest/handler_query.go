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
		q, err := DecodeQueryFromRequest(r)
		if err != nil {
			HandleHTTPError(monitor.Logger(), w, r, err)
			return
		}

		out, err := appHandler(r.Context(), q)
		if err != nil {
			HandleHTTPError(monitor.Logger(), w, r, err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		EncodeListResponse(monitor.Logger(), r.Context(), itemEncoder, dataWrapper, out, q, w, r)
	}
}
