package shared

import "context"

type HandlerFunc[OUT, IN any] func(context.Context, IN) (OUT, error)
