package rest

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type operationIdKey struct{}

func OperationIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		operationIdStr := ""
		operationId, err := uuid.NewV7()
		if err == nil {
			operationIdStr = operationId.String()
		}
		*r = *r.WithContext(context.WithValue(r.Context(), operationIdKey{}, operationIdStr))
		next.ServeHTTP(w, r)
	})
}

func operationIdFromContext(ctx context.Context) string {
	operationId, ok := ctx.Value(operationIdKey{}).(string)
	if !ok {
		return ""
	}
	return operationId
}
