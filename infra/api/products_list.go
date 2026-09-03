package api

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/infra/rest"
	"github.com/mytheresa/go-hiring-challenge/shared"
)

// ProductsList
// @Summary      List products with filters and pagination
// @Tags         Catalog
// @Router       /v1/catalog [get]
// @Accept       json
// @Produce      json
// @Param        page   			query      int  	false  "Page number"
// @Param        pageSize   		query      int  	false  "Page size"
// @Param        filter_price_lt   	query      string  	false  "Filter by price lower than"
// @Param        filter_category_eq query      string  	false  "Filter by category"
// @Success      200  {object}  rest.Response[productsList]
// @Failure      500  {object}  rest.ErrorResponse
// @Failure      400  {object}  rest.ErrorResponse
func productsListController(monitor shared.Monitor, app *catalog.App) http.HandlerFunc {
	return rest.NewListByQueryHandle(monitor, app.ListProducts, encodeProductResponse, wrapProductsListResponse)
}

type productsList struct {
	Products []product `json:"products"`
}

func wrapProductsListResponse(products []product) productsList {
	return productsList{Products: products}
}
