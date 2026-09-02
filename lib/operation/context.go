package operation

import (
	"context"

	"github.com/google/uuid"
)

type operationIdKey struct{}

func AddIdToContext(ctx context.Context) context.Context {
	operationIdStr := ""
	operationId, err := uuid.NewV7()
	if err == nil {
		operationIdStr = operationId.String()
	}
	return context.WithValue(ctx, operationIdKey{}, operationIdStr)
}

func IdFromContext(ctx context.Context) string {
	operationId, ok := ctx.Value(operationIdKey{}).(string)
	if !ok {
		return ""
	}
	return operationId
}
