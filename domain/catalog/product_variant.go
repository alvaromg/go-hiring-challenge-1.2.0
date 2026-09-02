package catalog

import "github.com/shopspring/decimal"

type Variant struct {
	ID        uint
	ProductID uint
	Name      string
	SKU       string
	Price     *decimal.Decimal
}
