package rest

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/lib/operation"
)

func OperationIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxWithOperationId := operation.AddIdToContext(r.Context())
		*r = *r.WithContext(ctxWithOperationId)
		next.ServeHTTP(w, r)
	})
}
