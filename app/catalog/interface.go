package catalog

import "github.com/mytheresa/go-hiring-challenge/domain/catalog"

type ProductsRepository interface {
	GetAllProducts() ([]*catalog.Product, error)
}
