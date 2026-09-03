package catalog

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/mytheresa/go-hiring-challenge/domain/list"
	"github.com/mytheresa/go-hiring-challenge/domain/query"
)

type ProductsRepository interface {
	GetProducts(ctx context.Context, q *query.Query) (list.ListResponse[*catalog.Product], error)
	GetProductByCode(ctx context.Context, code catalog.ProductCode) (*catalog.Product, error)
}

type CategoriesRepository interface {
	GetCategories(ctx context.Context, q *query.Query) (list.ListResponse[*catalog.Category], error)
	CreateCategories(ctx context.Context, categories []*catalog.Category) ([]*catalog.Category, error)
}
