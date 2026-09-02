package catalog

import (
	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/mytheresa/go-hiring-challenge/domain/list"
	"github.com/mytheresa/go-hiring-challenge/domain/query"
)

type ProductsRepository interface {
	GetAllProducts(q *query.Query) (list.ListResponse[*catalog.Product], error)
	GetProductByCode(code string) (*catalog.Product, error)
}

type CategoriesRepository interface {
	GetCategories(q *query.Query) (list.ListResponse[*catalog.Category], error)
	CreateCategories(categories []*catalog.Category) ([]*catalog.Category, error)
}
