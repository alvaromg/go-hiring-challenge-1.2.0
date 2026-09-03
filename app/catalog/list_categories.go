package catalog

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/mytheresa/go-hiring-challenge/domain/list"
	"github.com/mytheresa/go-hiring-challenge/domain/query"
)

func (app *App) ListCategories(ctx context.Context, q *query.Query) (list.ListResponse[*catalog.Category], error) {
	var zero list.ListResponse[*catalog.Category]

	// set default sorting
	q.AddSort("code", false)

	res, err := app.categoriesRepo.GetCategories(ctx, q)
	if err != nil {
		return zero, err
	}

	return res, nil
}
