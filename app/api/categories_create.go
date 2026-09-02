package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	catalogapp "github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/mytheresa/go-hiring-challenge/lib/rest"
)

type createCategoriesRequest struct {
	Categories []category `json:"categories"`
}

func decodeCreateCategoriesFromRequest(r *http.Request) (catalogapp.CreateCategoriesInput, error) {
	var zero catalogapp.CreateCategoriesInput
	var req createCategoriesRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return zero, fmt.Errorf("%w: invalid request body: %s", rest.ErrorBadRequest, err)
	}

	if len(req.Categories) == 0 {
		return zero, fmt.Errorf("%w: categories list must not be empty", rest.ErrorBadRequest)
	}

	requested := make([]catalogapp.RequestedCategory, len(req.Categories))
	for i, c := range req.Categories {
		requested[i] = catalogapp.RequestedCategory{
			Code: c.Code,
			Name: c.Name,
		}
	}

	return catalogapp.CreateCategoriesInput{Categories: requested}, nil
}

func encodeCreateCategoriesResponse(categories []*catalog.Category) (categoriesList, error) {
	encoded := make([]category, len(categories))
	for i, c := range categories {
		encoded[i] = encodeCategoryResponse(c)
	}
	return wrapCategoriesResponse(encoded), nil
}
