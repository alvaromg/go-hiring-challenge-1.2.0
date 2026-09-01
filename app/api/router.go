package api

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/lib/rest"
	"github.com/mytheresa/go-hiring-challenge/shared"
)

func NewApiRouter(monitor shared.Monitor, catalogApp *catalog.App) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /catalog", rest.NewHandler(monitor, catalogApp.GetProducts, decodeGetProductsListRequest, encodeProductsListResponse, http.StatusOK))

	return mux
}
