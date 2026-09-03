package catalog

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/mytheresa/go-hiring-challenge/domain/list"
	"github.com/mytheresa/go-hiring-challenge/domain/query"
)

type GetProductsInput struct {
	Query *query.Query
}

func (app *App) ListProducts(ctx context.Context, q *query.Query) (list.ListResponse[*catalog.Product], error) {
	var zero list.ListResponse[*catalog.Product]

	if !q.HasAnySort() {
		// if not sorting defined set default for code ascending
		q.AddSort("code", false)
	}

	res, err := app.productsRepo.GetProducts(ctx, q)
	if err != nil {
		return zero, err
	}

	return res, nil
}
