package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	appcatalog "github.com/mytheresa/go-hiring-challenge/app/catalog"
	domaincatalog "github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/mytheresa/go-hiring-challenge/infra/rest"
	"github.com/mytheresa/go-hiring-challenge/shared"
)

// Create categories
// @Summary      Create multiple categoriescategories
// @Tags         Catalog
// @Router       /v1/categories [post]
// @Accept       json
// @Produce      json
// @param		request	body	createCategoriesRequest	true "New categories to create"
// @Success      201  {object}  rest.Response[categoriesList]
// @Failure      500  {object}  rest.ErrorResponse
// @Failure      409  {object}  rest.ErrorResponse
func createCategoriesController(monitor shared.Monitor, app *appcatalog.App) http.HandlerFunc {
	return rest.NewHandler(monitor, app.CreateCategories, decodeCreateCategoriesFromRequest, encodeCreateCategoriesResponse, http.StatusCreated)
}

type createCategoriesRequest struct {
	Categories []category `json:"categories"`
}

func decodeCreateCategoriesFromRequest(r *http.Request) (appcatalog.CreateCategoriesInput, error) {
	var zero appcatalog.CreateCategoriesInput
	var req createCategoriesRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return zero, fmt.Errorf("%w: invalid request body: %s", rest.ErrorBadRequest, err)
	}

	if len(req.Categories) == 0 {
		return zero, fmt.Errorf("%w: categories list must not be empty", rest.ErrorBadRequest)
	}

	requested := make([]appcatalog.RequestedCategory, len(req.Categories))
	for i, c := range req.Categories {

		code, err := domaincatalog.ParseCategoryCode(c.Code)
		if err != nil {
			return zero, err
		}

		requested[i] = appcatalog.RequestedCategory{
			Code: code,
			Name: c.Name,
		}
	}

	return appcatalog.CreateCategoriesInput{Categories: requested}, nil
}

func encodeCreateCategoriesResponse(categories []*domaincatalog.Category) (categoriesList, error) {
	encoded := make([]category, len(categories))
	for i, c := range categories {
		encoded[i] = encodeCategoryResponse(c)
	}
	return wrapCategoriesResponse(encoded), nil
}
