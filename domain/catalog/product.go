package catalog

import "github.com/shopspring/decimal"

type Product struct {
	ID       uint
	Code     string
	Price    decimal.Decimal
	Variants []Variant
}
