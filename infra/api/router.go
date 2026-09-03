package api

import (
	"net/http"
	"time"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/infra/rest"
	"github.com/mytheresa/go-hiring-challenge/shared"
)

type Middleware func(http.Handler) http.Handler

// DefaultHandlerTimeout is used when no timeout is configured (e.g. HTTP_HANDLER_TIMEOUT is unset).
const DefaultHandlerTimeout = 10 * time.Second

func NewApiRouter(monitor shared.Monitor, catalogApp *catalog.App, corsAllowedOrigins []string, handlerTimeout time.Duration) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/catalog", productsListController(monitor, catalogApp))
	mux.HandleFunc("GET /v1/catalog/{code}", productsDetailController(monitor, catalogApp))
	mux.HandleFunc("GET /v1/categories", categoriesListController(monitor, catalogApp))
	mux.HandleFunc("POST /v1/categories", createCategoriesController(monitor, catalogApp))
	mux.HandleFunc("/", rest.DefaultNotFound(monitor.Logger()))

	return chainMiddlewares(mux,
		rest.OperationIdMiddleware,
		rest.NewLoggingMiddleware(monitor.Logger()),
		rest.NewTimeoutMiddleware(handlerTimeout),
		rest.NewCorsMiddleware(corsAllowedOrigins...),
	)
}

// Chain applies middlewares in order, so the first one runs first
func chainMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
