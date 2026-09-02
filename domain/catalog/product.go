package catalog

import "github.com/shopspring/decimal"

type Product struct {
	Code     string
	Price    decimal.Decimal
	Variants []Variant
	Category *Category
}
