package rest

import (
	"context"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/domain/list"
	"github.com/mytheresa/go-hiring-challenge/domain/query"
	"github.com/mytheresa/go-hiring-challenge/shared"
)

func NewListByQueryHandle[DO, RO any](
	monitor shared.Monitor,
	appHandler func(context.Context, *query.Query) (list.ListResponse[DO], error),
	itemEncoder func(DO) RO,
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
		EncodeListResponse(monitor.Logger(), r.Context(), itemEncoder, out, q, w, r)
	}
}
