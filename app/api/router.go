package api

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/lib/rest"
	"github.com/mytheresa/go-hiring-challenge/shared"
)

func NewApiRouter(monitor shared.Monitor, catalogApp *catalog.App) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /catalog", rest.NewListByQueryHandle(monitor, catalogApp.GetProducts, encodeProductResponse))

	return mux
}
