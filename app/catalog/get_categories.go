package catalog

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/mytheresa/go-hiring-challenge/domain/list"
	"github.com/mytheresa/go-hiring-challenge/domain/query"
)

func (app *App) GetCategories(ctx context.Context, q *query.Query) (list.ListResponse[*catalog.Category], error) {
	var zero list.ListResponse[*catalog.Category]

	res, err := app.categoriesRepo.GetCategories(q)
	if err != nil {
		return zero, err
	}

	return res, nil
}
