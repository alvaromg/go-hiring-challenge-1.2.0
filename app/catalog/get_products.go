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

func (app *App) GetProducts(ctx context.Context, q *query.Query) (list.ListResponse[*catalog.Product], error) {
	var zero list.ListResponse[*catalog.Product]

	res, err := app.productsRepo.GetAllProducts(q)
	if err != nil {
		return zero, err
	}

	return res, nil
}
