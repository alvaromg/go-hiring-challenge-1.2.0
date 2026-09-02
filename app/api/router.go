package api

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/lib/rest"
	"github.com/mytheresa/go-hiring-challenge/shared"
)

type Middleware func(http.Handler) http.Handler

func NewApiRouter(monitor shared.Monitor, catalogApp *catalog.App) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /catalog", rest.NewListByQueryHandle(monitor, catalogApp.GetProducts, encodeProductResponse, wrapProductsListResponse))
	mux.HandleFunc("GET /catalog/{code}", rest.NewHandler(monitor, catalogApp.GetProduct, decodeProductCodeFromRequest, encodeProductDataResponse, http.StatusOK))
	mux.HandleFunc("GET /categories", rest.NewListByQueryHandle(monitor, catalogApp.GetCategories, encodeCategoryResponse, wrapCategoriesResponse))
	mux.HandleFunc("POST /categories", rest.NewHandler(monitor, catalogApp.CreateCategories, decodeCreateCategoriesFromRequest, encodeCreateCategoriesResponse, http.StatusCreated))
	mux.HandleFunc("/", rest.DefaultNotFound)

	return chainMiddlewares(mux, rest.OperationIdMiddleware, rest.NewLoggingMiddleware(monitor.Logger()))
}

// Chain applies middlewares in order, so the first one runs first
func chainMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
