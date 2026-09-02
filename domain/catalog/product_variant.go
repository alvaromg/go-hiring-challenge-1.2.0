package catalog

import "github.com/mytheresa/go-hiring-challenge/domain/price"

type Variant struct {
	ID        uint
	ProductID uint
	Name      string
	SKU       string
	Price     *price.Price
}
