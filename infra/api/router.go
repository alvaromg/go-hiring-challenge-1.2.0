package api

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/infra/rest"
	"github.com/mytheresa/go-hiring-challenge/shared"
)

type Middleware func(http.Handler) http.Handler

func NewApiRouter(monitor shared.Monitor, catalogApp *catalog.App) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/catalog", rest.NewListByQueryHandle(monitor, catalogApp.ListProducts, encodeProductResponse, wrapProductsListResponse))
	mux.HandleFunc("GET /v1/catalog/{code}", rest.NewHandler(monitor, catalogApp.ProductDetail, decodeProductCodeFromRequest, encodeProductDetailResponse, http.StatusOK))
	mux.HandleFunc("GET /v1/categories", rest.NewListByQueryHandle(monitor, catalogApp.ListCategories, encodeCategoryResponse, wrapCategoriesResponse))
	mux.HandleFunc("POST /v1/categories", rest.NewHandler(monitor, catalogApp.CreateCategories, decodeCreateCategoriesFromRequest, encodeCreateCategoriesResponse, http.StatusCreated))
	mux.HandleFunc("/", rest.DefaultNotFound(monitor.Logger()))

	return chainMiddlewares(mux, rest.OperationIdMiddleware, rest.NewLoggingMiddleware(monitor.Logger()))
}

// Chain applies middlewares in order, so the first one runs first
func chainMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
