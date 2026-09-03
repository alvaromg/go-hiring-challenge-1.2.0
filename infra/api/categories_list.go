package api

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/infra/rest"
	"github.com/mytheresa/go-hiring-challenge/shared"
)

// CategoriesList
// @Summary      List categories with pagination
// @Router       /v1/categories [get]
// @Accept       json
// @Produce      json
// @Param        page   			query      int  	false  "Page number"
// @Param        pageSize   		query      int  	false  "Page size"
// @Success      200  {object}  rest.ListResponse[categoriesList]
// @Failure      500  {object}  rest.ErrorResponse
// @Failure      400  {object}  rest.ErrorResponse
func categoriesListController(monitor shared.Monitor, app *catalog.App) http.HandlerFunc {
	return rest.NewListByQueryHandle(monitor, app.ListCategories, encodeCategoryResponse, wrapCategoriesResponse)
}

type categoriesList struct {
	Categories []category `json:"categories"`
}

func wrapCategoriesResponse(categories []category) categoriesList {
	return categoriesList{Categories: categories}
}
